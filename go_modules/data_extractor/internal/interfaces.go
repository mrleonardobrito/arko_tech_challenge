package internal

import "context"

type Batcher[T any] interface {
	Push(ctx context.Context, item T) (read bool, batch []T, err error)
	Flush(ctx context.Context) (batch []T, err error)
	Progress() string
	AddUpsertedBatch()
}

type DBEncoder[T any] interface {
	Encode(ctx context.Context, v T) ([]any, error)
}

type Source[T any] interface {
	ItemCount() int
	Next(ctx context.Context) (T, error)
	HasNext(ctx context.Context) bool
	Close() error
}

type Identifiable interface {
	ID() string
}

type Sink[T any] interface {
	WriteBatch(ctx context.Context, batch []T) error
}
