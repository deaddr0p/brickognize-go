package brickognize

import (
	"context"
	"testing"
)

func TestQueueHandlesInvalidFiles(t *testing.T) {
	client := NewClient()

	paths := []string{"fake.txt"}

	results := client.PredictPartsQueue(context.Background(), paths)

	if len(results) != 1 {
		t.Fatal("expected 1 result")
	}

	if results[0].Err == nil {
		t.Fatal("expected error for invalid image")
	}
}
