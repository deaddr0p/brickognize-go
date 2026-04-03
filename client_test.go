package brickognize

import (
	"testing"
	"time"
)

func TestNewClientDefaults(t *testing.T) {
	c := NewClient()

	if c.baseURL == "" {
		t.Fatal("expected baseURL to be set")
	}

	if c.workers != 4 {
		t.Fatalf("expected default workers=4, got %d", c.workers)
	}

	if c.retries != 3 {
		t.Fatalf("expected retries=3, got %d", c.retries)
	}

	if c.httpClient.Timeout != 30*time.Second {
		t.Fatal("unexpected timeout")
	}
}

func TestWithOptions(t *testing.T) {
	c := NewClient(
		WithWorkers(10),
		WithRetries(5),
		WithRateLimit(20),
	)

	if c.workers != 10 {
		t.Fatal("workers not applied")
	}

	if c.retries != 5 {
		t.Fatal("retries not applied")
	}

	if c.limiter == nil {
		t.Fatal("rate limiter not set")
	}
}

func TestInvalidImage(t *testing.T) {
	ok := IsValidImage("nonexistent.jpg")
	if ok {
		t.Fatal("expected invalid image")
	}
}
