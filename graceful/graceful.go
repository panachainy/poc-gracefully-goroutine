package graceful

import (
	"context"
	"log"
	"sync"
)

// RunCancellableGoroutine runs a function in a goroutine that can be cancelled via context.
// It uses a WaitGroup to allow the caller to wait for completion.
func RunCancellableGoroutine(ctx context.Context, wg *sync.WaitGroup, prefix string, fn func(ctx context.Context)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				// log.Printf("panic recovered: %v", r)
				log.Printf("%s goroutine recovered from panic: %v\n", prefix, r)
			}
		}()
		fn(ctx)
	}()
}
