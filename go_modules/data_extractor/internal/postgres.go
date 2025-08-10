package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type ConflictMode string

const (
	ConflictModeUpdate ConflictMode = "update"
	ConflictModeIgnore ConflictMode = "ignore"
)

type TableSpec struct {
	Name           string
	Columns        []string
	ConflictMode   ConflictMode
	ConflictColumn string
	UpdateColumns  []string
}

type PGOptions struct {
	TxTimeout time.Duration
}

func (o PGOptions) Default() PGOptions {
	if o.TxTimeout == 0 {
		o.TxTimeout = 10 * time.Second
	}
	return o
}

type Postgres[T any] struct {
	pool    *pgxpool.Pool
	spec    TableSpec
	encoder DBEncoder[T]
	options PGOptions
}

func NewPostgresRepository[T any](pool *pgxpool.Pool, spec TableSpec, encoder DBEncoder[T], options PGOptions) Sink[T] {
	return &Postgres[T]{
		pool:    pool,
		spec:    spec,
		encoder: encoder,
		options: options,
	}
}

func (p *Postgres[T]) WriteBatch(ctx context.Context, batch []T) error {
	if len(batch) == 0 {
		log.Printf("Skipping empty batch")
		return nil
	}

	return p.upsertBatch(ctx, batch)
}

func (p *Postgres[T]) upsertBatch(ctx context.Context, batch []T) error {
	cols := p.spec.Columns
	nCols := len(cols)

	seen := make(map[any]bool)
	uniqueBatch := make([]T, 0, len(batch))

	for _, v := range batch {
		values, err := p.encoder.Encode(ctx, v)
		if err != nil {
			return fmt.Errorf("error encoding batch: %w", err)
		}

		if len(values) != nCols {
			return fmt.Errorf("expected %d values, got %d", nCols, len(values))
		}

		if !seen[values[0]] {
			seen[values[0]] = true
			uniqueBatch = append(uniqueBatch, v)
		}
	}

	placeholders := make([]string, 0, len(uniqueBatch))
	args := make([]any, 0, len(uniqueBatch)*nCols)

	for i, v := range uniqueBatch {
		values, err := p.encoder.Encode(ctx, v)
		if err != nil {
			return fmt.Errorf("error encoding batch: %w", err)
		}

		base := i * nCols
		slots := make([]string, nCols)
		for j := 0; j < nCols; j++ {
			slots[j] = fmt.Sprintf("$%d", base+j+1)
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(slots, ", ")))
		args = append(args, values...)
	}

	var conflitClause, updateClause string

	if p.spec.ConflictMode == ConflictModeUpdate && len(p.spec.UpdateColumns) > 0 {
		conflictColumn := p.spec.ConflictColumn
		if conflictColumn == "" {
			conflictColumn = p.spec.Columns[0]
		}
		conflitClause = fmt.Sprintf(" ON CONFLICT (%s) ", conflictColumn)
		if len(p.spec.UpdateColumns) > 0 {
			sets := make([]string, len(p.spec.UpdateColumns))
			for i, col := range p.spec.UpdateColumns {
				sets[i] = fmt.Sprintf("%s = EXCLUDED.%s", col, col)
			}
			updateClause = " DO UPDATE SET " + strings.Join(sets, ", ")
		} else {
			updateClause = " DO NOTHING"
		}
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s%s%s",
		p.spec.Name,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
		conflitClause,
		updateClause,
	)

	ctxTx := ctx
	if p.options.TxTimeout > 0 {
		var cancel context.CancelFunc
		ctxTx, cancel = context.WithTimeout(ctx, p.options.TxTimeout)
		defer cancel()
	}

	tx, err := p.pool.BeginTx(ctxTx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return fmt.Errorf("error executing query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
