package pipelines

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/arko_tech_challenge/go_modules/data_extractor/internal"
)

func RunPipeline[T any](ctx context.Context, src internal.Source[T], batcher internal.Batcher[T], db internal.Sink[T]) error {
	defer func() { _ = src.Close() }()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Println(batcher.Progress())
			}
		}
	}()

	batchChan := make(chan []T, runtime.NumCPU())
	errChan := make(chan error, 1)
	doneChan := make(chan struct{})

	numWorkers := 5

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batch := range batchChan {
				select {
				case <-ctx.Done():
					return
				default:
					if err := db.WriteBatch(ctx, batch); err != nil {
						select {
						case errChan <- fmt.Errorf("worker %d error: %w", workerID, err):
						default:
						}
						return
					}
					batcher.AddUpsertedBatch()
				}
			}
		}(i)
	}

	go func() {
		defer close(batchChan)

		for src.HasNext(ctx) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			in, err := src.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				errChan <- fmt.Errorf("error reading from source: %w", err)
				return
			}

			ready, batch, err := batcher.Push(ctx, in)
			if err != nil {
				errChan <- fmt.Errorf("error pushing to batcher: %w", err)
				return
			}
			if ready {
				select {
				case <-ctx.Done():
					return
				case batchChan <- batch:
				}
			}
		}

		remaining, err := batcher.Flush(ctx)
		if err != nil {
			errChan <- fmt.Errorf("error flushing batcher: %w", err)
			return
		}
		if len(remaining) > 0 {
			select {
			case <-ctx.Done():
				return
			case batchChan <- remaining:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		log.Println(batcher.Progress())
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
