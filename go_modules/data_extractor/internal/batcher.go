package internal

import (
	"context"
	"fmt"
	"sync"
)

type FixedSizeBatcher[T any] struct {
	batchSize       int
	batch           []T
	totalBatches    int
	upsertedBatches int
	mu              sync.Mutex
}

func (b *FixedSizeBatcher[T]) Progress() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	progress := float64(b.upsertedBatches) / float64(b.totalBatches) * 100
	return fmt.Sprintf("%.2f%% (%d/%d batches)", progress, b.upsertedBatches, b.totalBatches)
}

func (b *FixedSizeBatcher[T]) AddUpsertedBatch() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.upsertedBatches++
}

func NewFixedSizeBatcher[T any](batchSize int, totalItems int) Batcher[T] {
	return &FixedSizeBatcher[T]{
		batchSize:       batchSize,
		totalBatches:    totalItems / batchSize,
		upsertedBatches: 0,
	}
}

func (b *FixedSizeBatcher[T]) Push(ctx context.Context, item T) (read bool, batch []T, err error) {
	b.batch = append(b.batch, item)

	if len(b.batch) >= b.batchSize {
		out := make([]T, len(b.batch))
		copy(out, b.batch)
		b.batch = b.batch[:0]
		return true, out, nil
	}
	return false, nil, nil
}

func (b *FixedSizeBatcher[T]) Flush(ctx context.Context) (batch []T, err error) {
	if len(b.batch) == 0 {
		return nil, nil
	}
	batch = make([]T, len(b.batch))
	copy(batch, b.batch)
	b.batch = b.batch[:0]
	return batch, nil
}
