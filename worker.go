package brickognize

import (
	"context"
	"sync"
)

type Result struct {
	Path     string
	Response *Response
	Err      error
}

func (c *Client) PredictPartsQueue(ctx context.Context, paths []string) []Result {
	jobs := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup

	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if !IsValidImage(p) {
					results <- Result{Path: p, Err: ErrInvalidImage}
					continue
				}
				r, err := c.PredictParts(ctx, p)
				results <- Result{Path: p, Response: r, Err: err}
			}
		}()
	}

	go func() {
		for _, p := range paths {
			jobs <- p
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []Result
	for r := range results {
		out = append(out, r)
	}
	return out
}
