package brickognize

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

type integrationEndpointCase struct {
	name        string
	imagesEnv   string
	expectedEnv string
	call        func(*Client, context.Context, string) (*Response, error)
}

type integrationImageCase struct {
	path          string
	expectedCount int
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

func parseExpectedCounts(raw string) ([]int, error) {
	parts := parseCSV(raw)
	counts := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid expected count %q", p)
		}
		counts = append(counts, n)
	}
	return counts, nil
}

func buildImageCases(imageDir string, files []string, counts []int) ([]integrationImageCase, []string) {
	runs := make([]integrationImageCase, 0, len(files))
	missing := make([]string, 0)

	// if counts is shorter than files, treat missing counts as 1
	if len(files) < len(counts) {
		for range len(files) - len(counts) {
			counts = append(counts, 1)
		}
	}

	for i, file := range files {
		path := file
		if !filepath.IsAbs(path) {
			path = filepath.Join(imageDir, path)
		}

		if _, err := os.Stat(path); err != nil {
			missing = append(missing, path)
			continue
		}

		runs = append(runs, integrationImageCase{path: path, expectedCount: counts[i]})
	}

	return runs, missing
}

func TestBuildImageCasesSkipsMissingImage(t *testing.T) {
	tempDir := t.TempDir()
	existing := filepath.Join(tempDir, "present.jpg")
	if err := os.WriteFile(existing, []byte("fixture"), 0o600); err != nil {
		t.Fatalf("failed to create fixture: %v", err)
	}

	runs, missing := buildImageCases(tempDir, []string{"present.jpg", "missing.jpg"}, []int{2, 1})

	if len(runs) != 1 {
		t.Fatalf("expected 1 runnable case, got %d", len(runs))
	}
	if runs[0].path != existing || runs[0].expectedCount != 2 {
		t.Fatalf("unexpected runnable case: %+v", runs[0])
	}
	if len(missing) != 1 || !strings.HasSuffix(missing[0], "missing.jpg") {
		t.Fatalf("expected missing missing.jpg, got %+v", missing)
	}
}

func TestIntegrationEndpointMatrix(t *testing.T) {
	_ = godotenv.Load(".env")

	imageDir := envOrDefault("TEST_IMAGE_DIR", "./testdata")

	cases := []integrationEndpointCase{
		{
			name:        "parts",
			imagesEnv:   "TEST_IMAGE_PARTS",
			expectedEnv: "TEST_EXPECTED_PARTS",
			call:        (*Client).PredictParts,
		},
		{
			name:        "sets",
			imagesEnv:   "TEST_IMAGE_SETS",
			expectedEnv: "TEST_EXPECTED_SETS",
			call:        (*Client).PredictSets,
		},
		{
			name:        "minifigs",
			imagesEnv:   "TEST_IMAGE_MINIFIGS",
			expectedEnv: "TEST_EXPECTED_MINIFIGS",
			call:        (*Client).PredictMinifigs,
		},
	}

	client := NewClient()
	runCount := 0

	for _, tc := range cases {
		tc := tc

		images := parseCSV(os.Getenv(tc.imagesEnv))
		if len(images) == 0 {
			t.Logf("skipping %s: %s is empty", tc.name, tc.imagesEnv)
			continue
		}
		expected, err := parseExpectedCounts(os.Getenv(tc.expectedEnv))
		if err != nil {
			t.Fatalf("%s: %v", tc.expectedEnv, err)
		}

		if len(expected) != len(images) {
			t.Fatalf("%s: image list and expected counts must have same length", tc.name)
		}

		runs, missing := buildImageCases(imageDir, images, expected)
		for _, missingPath := range missing {
			t.Logf("skipping %s: missing file %q", tc.name, missingPath)
		}

		for _, run := range runs {
			run := run
			runCount++

			t.Run(tc.name+"_specific_"+filepath.Base(run.path), func(t *testing.T) {
				resp, err := tc.call(client, context.Background(), run.path)
				if err != nil {
					t.Fatalf("specific endpoint error: %v", err)
				}
				if resp == nil {
					t.Fatal("expected response from specific endpoint")
				}
				if len(resp.Items) != run.expectedCount {
					t.Fatalf("expected %d items, got %d", run.expectedCount, len(resp.Items))
				}
			})

			t.Run(tc.name+"_generic_"+filepath.Base(run.path), func(t *testing.T) {
				resp, err := client.PredictAll(context.Background(), run.path)
				if err != nil {
					t.Fatalf("generic endpoint error: %v", err)
				}
				if resp == nil {
					t.Fatal("expected response from generic endpoint")
				}
				if len(resp.Items) != run.expectedCount {
					t.Fatalf("expected %d items, got %d", run.expectedCount, len(resp.Items))
				}
			})
		}
	}

	if runCount == 0 {
		t.Skip("no integration fixtures found; configure .env and add test images")
	}
}
