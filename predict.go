package brickognize

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func (c *Client) doRequest(ctx context.Context, endpoint, path string) (*Response, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Detect actual MIME type so the API accepts the file part
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	mimeType := http.DetectContentType(buf[:n])
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="query_image"; filename="%s"`, filepath.Base(path)))
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, &body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Printf("Failed to parse response: %v\nResponse body: %s", err, string(bodyBytes))
		return nil, err
	}

	if len(result.Items) == 0 {
		log.Printf("No results found in response. Response body: %s", string(bodyBytes))
	}

	return &result, nil
}

func (c *Client) predict(ctx context.Context, endpoint, path string) (*Response, error) {
	var out *Response
	err := retry(c.retries, func() error {
		if c.limiter != nil {
			c.limiter.Wait()
		}
		r, err := c.doRequest(ctx, endpoint, path)
		if err == nil {
			out = r
		}
		return err
	})
	return out, err
}

func (c *Client) PredictParts(ctx context.Context, path string) (*Response, error) {
	return c.predict(ctx, "/predict/parts/", path)
}

func (c *Client) PredictSets(ctx context.Context, path string) (*Response, error) {
	return c.predict(ctx, "/predict/sets/", path)
}

func (c *Client) PredictMinifigs(ctx context.Context, path string) (*Response, error) {
	return c.predict(ctx, "/predict/figs/", path)
}

func (c *Client) PredictAll(ctx context.Context, path string) (*Response, error) {
	return c.predict(ctx, "/predict/", path)
}
