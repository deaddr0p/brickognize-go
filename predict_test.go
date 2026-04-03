package brickognize

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeTestPNG(t *testing.T) string {
	t.Helper()

	png := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0x60, 0x60, 0x60, 0xf8,
		0x0f, 0x00, 0x01, 0x04, 0x01, 0x00, 0x5f, 0x4b,
		0xb4, 0x7d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}

	filePath := filepath.Join(t.TempDir(), "test.png")
	if err := os.WriteFile(filePath, png, 0o600); err != nil {
		t.Fatalf("failed writing png: %v", err)
	}
	return filePath
}

func TestPredictEndpoints(t *testing.T) {
	tests := []struct {
		name         string
		expectedPath string
		call         func(*Client, context.Context, string) (*Response, error)
	}{
		{name: "parts", expectedPath: "/predict/parts/", call: (*Client).PredictParts},
		{name: "sets", expectedPath: "/predict/sets/", call: (*Client).PredictSets},
		{name: "minifigs", expectedPath: "/predict/figs/", call: (*Client).PredictMinifigs},
		{name: "all", expectedPath: "/predict/", call: (*Client).PredictAll},
	}

	imagePath := writeTestPNG(t)

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != tc.expectedPath {
					t.Fatalf("expected path %q, got %q", tc.expectedPath, r.URL.Path)
				}
				if got := r.Header.Get("Accept"); got != "application/json" {
					t.Fatalf("expected Accept header application/json, got %q", got)
				}
				if ct := r.Header.Get("Content-Type"); ct == "" {
					t.Fatal("expected multipart content-type")
				}

				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(Response{Items: []Item{{ID: "1", Name: "test", Score: 0.99}}})
			}))
			defer server.Close()

			client := NewClient(WithRetries(1))
			client.baseURL = server.URL

			resp, err := tc.call(client, context.Background(), imagePath)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if resp == nil || len(resp.Items) != 1 {
				t.Fatal("expected one item in response")
			}
		})
	}
}

