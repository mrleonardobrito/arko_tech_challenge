package internal

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/net/html/charset"
)

type APISource[T any] struct {
	apiUrl        string
	mapResponseFn func([]byte) ([]T, bool, error)
	cache         []T
	hasNext       bool
}

func NewAPISource[T any](apiUrl string, mapResponseFn func([]byte) ([]T, bool, error)) Source[T] {
	source := &APISource[T]{apiUrl: apiUrl, mapResponseFn: mapResponseFn, cache: make([]T, 0), hasNext: true}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return source
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return source
	}
	defer resp.Body.Close()

	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return source
	}

	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		return source
	}

	records, _, err := mapResponseFn(body)
	if err != nil {
		return source
	}

	source.cache = append(source.cache, records...)
	return source
}

func (s *APISource[T]) HasNext(ctx context.Context) bool {
	return s.hasNext || len(s.cache) > 0
}

func (s *APISource[T]) ItemCount() int {
	return len(s.cache)
}

func (s *APISource[T]) Next(ctx context.Context) (T, error) {
	if len(s.cache) == 0 {
		return s.lazyRequest(ctx)
	}

	v := s.cache[0]
	s.cache = s.cache[1:]

	return v, nil
}

func (s *APISource[T]) lazyRequest(ctx context.Context) (T, error) {
	var empty T

	req, err := http.NewRequestWithContext(ctx, "GET", s.apiUrl, nil)
	if err != nil {
		return empty, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return empty, err
	}
	defer resp.Body.Close()

	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return empty, err
	}

	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		return empty, err
	}

	records, hasNext, err := s.mapResponseFn(body)
	if err != nil {
		return empty, err
	}
	s.hasNext = hasNext

	if len(records) > 0 {
		s.cache = append(s.cache, records...)
	}
	return s.Next(ctx)
}

func (s *APISource[T]) Close() error {
	return nil
}
