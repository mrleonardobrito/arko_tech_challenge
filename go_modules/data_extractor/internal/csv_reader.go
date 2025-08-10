package internal

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"context"
)

type CSVSource[T any] struct {
	reader     *csv.Reader
	mapFn      func([]string) (T, error)
	closer     io.Closer
	nextRecord []string
	nextErr    error
	itemCount  int
}

func (s *CSVSource[T]) SkipHeader() error {
	header, err := s.reader.Read()
	if err != nil {
		return fmt.Errorf("error skipping header: %w", err)
	}
	log.Printf("Header skipped: %v", header)
	return nil
}

func NewCSVSource[T any](r io.ReadCloser, comma rune, mapFn func([]string) (T, error), itemCount int) Source[T] {
	bufferedReader := bufio.NewReaderSize(r, 1024*1024)

	csvReader := csv.NewReader(bufferedReader)
	csvReader.Comma = comma
	csvReader.ReuseRecord = true
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	src := &CSVSource[T]{
		reader:    csvReader,
		mapFn:     mapFn,
		closer:    r,
		itemCount: itemCount,
	}

	if err := src.SkipHeader(); err != nil {
		log.Printf("Error skipping header: %v", err)
		return src
	}

	src.nextRecord, src.nextErr = src.reader.Read()
	if src.nextErr != nil {
		log.Printf("Error reading first record: %v", src.nextErr)
	}

	return src
}

func (s *CSVSource[T]) ItemCount() int {
	return s.itemCount
}

func (s *CSVSource[T]) HasNext(ctx context.Context) bool {
	return s.nextErr != io.EOF
}

func (s *CSVSource[T]) Next(ctx context.Context) (T, error) {
	var zero T

	if s.nextErr != nil && s.nextErr != io.EOF {
		return zero, s.nextErr
	}

	record := s.nextRecord

	s.nextRecord, s.nextErr = s.reader.Read()

	v, err := s.mapFn(record)
	if err != nil {
		return zero, fmt.Errorf("error mapping record %v: %w", record, err)
	}

	return v, nil
}

func (s *CSVSource[T]) Close() error {
	return s.closer.Close()
}
