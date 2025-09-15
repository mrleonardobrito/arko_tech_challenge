package internal

import (
	"context"
	"fmt"
	"math"
	"sync"
)

type FileSizeBatcher[T any] struct {
	maxBatchBytes       int
	estimatedTotalBytes int64

	sizeOf func(T) int

	mu                sync.Mutex
	currentBatch      []T
	currentBatchBytes int

	totalBatches    int
	upsertedBatches int
}

func NewFileSizeBatcher[T any](maxBatchMb int, estimatedTotalMb int64, sizeOf func(T) int) Batcher[T] {
	if sizeOf == nil {
		panic("sizeOf function must not be nil")
	}

	const MiB = 1024 * 1024
	var maxBytes int
	if maxBatchMb <= 0 {
		maxBytes = math.MaxInt / 2
	} else {
		maxBytes = maxBatchMb * MiB
	}

	totalBatches := 0
	if estimatedTotalMb > 0 {
		totalBatches = int(int64(math.Ceil(float64(estimatedTotalMb) / float64(maxBytes))))
	}

	return &FileSizeBatcher[T]{
		maxBatchBytes:       maxBytes,
		estimatedTotalBytes: estimatedTotalMb,
		sizeOf:              sizeOf,
		currentBatch:        make([]T, 0, 64),
		totalBatches:        totalBatches,
	}
}

func (b *FileSizeBatcher[T]) Push(ctx context.Context, item T) (bool, []T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	select {
	case <-ctx.Done():
		return false, nil, ctx.Err()
	default:
	}

	itemBytes := b.sizeOf(item)
	if itemBytes < 0 {
		return false, nil, fmt.Errorf("sizeOf returned negative size: %d", itemBytes)
	}

	if len(b.currentBatch) == 0 && itemBytes >= b.maxBatchBytes {
		return true, []T{item}, nil
	}

	if b.currentBatchBytes+itemBytes > b.maxBatchBytes && len(b.currentBatch) > 0 {
		out := make([]T, len(b.currentBatch))
		copy(out, b.currentBatch)
		b.currentBatch = b.currentBatch[:0]
		b.currentBatch = append(b.currentBatch, item)
		b.currentBatchBytes = itemBytes
		return true, out, nil
	}

	b.currentBatch = append(b.currentBatch, item)
	b.currentBatchBytes += itemBytes
	return false, nil, nil
}

func (b *FileSizeBatcher[T]) Flush(ctx context.Context) ([]T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(b.currentBatch) == 0 {
		return nil, nil
	}
	out := make([]T, len(b.currentBatch))
	copy(out, b.currentBatch)
	b.currentBatch = b.currentBatch[:0]
	b.currentBatchBytes = 0
	return out, nil
}

func (b *FileSizeBatcher[T]) Progress() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	total := b.totalBatches
	if total == 0 {
		total = b.upsertedBatches
		if total == 0 {
			total = 1
		}
	}
	progress := float64(b.upsertedBatches) / float64(total) * 100
	return fmt.Sprintf("%.2f%% (%d/%d batches)", progress, b.upsertedBatches, total)
}

func (b *FileSizeBatcher[T]) AddUpsertedBatch() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.upsertedBatches++
}
