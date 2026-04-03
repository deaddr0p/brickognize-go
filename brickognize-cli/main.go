package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deaddr0p/brickognize-go"
)

func main() {
	input := flag.String("input", "", "file or dir")
	flag.Parse()

	if *input == "" {
		fmt.Println("missing input")
		return
	}

	var paths []string

	info, _ := os.Stat(*input)
	if info.IsDir() {
		filepath.Walk(*input, func(p string, i os.FileInfo, e error) error {
			if !i.IsDir() {
				paths = append(paths, p)
			}
			return nil
		})
	} else {
		paths = []string{*input}
	}

	client := brickognize.NewClient()

	results := client.PredictPartsQueue(context.Background(), paths)

	for _, r := range results {
		if r.Err != nil {
			fmt.Println("ERR:", r.Path, r.Err)
			continue
		}
		fmt.Println("OK:", r.Path, r.Response.Items)
	}
}
