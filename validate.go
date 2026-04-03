package brickognize

import (
	"net/http"
	"os"
)

func IsValidImage(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil {
		return false
	}

	switch http.DetectContentType(buf) {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}
