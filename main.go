package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
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
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start multiple workers
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Received shutdown signal, cancelling context...")
	cancel()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers shut down gracefully. Exiting.")
}
