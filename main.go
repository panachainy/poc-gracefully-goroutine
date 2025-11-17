package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"poc-gracefully-goroutine/graceful"

	"github.com/gin-gonic/gin"
)

func worker(ctx context.Context, id int) {
	fmt.Printf("Worker %d started\n", id)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return
		default:
			fmt.Printf("Worker %d working...\n", id)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// Create Gin router
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start background workers
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		graceful.RunCancellableGoroutine(ctx, &wg, fmt.Sprintf("worker-%d", i), func(ctx context.Context) {
			worker(ctx, i)
		})
	}

	// Start server in a goroutine
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// Cancel context to stop workers
	cancel()

	// Gracefully shutdown server
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers shut down gracefully. Exiting.")
}
