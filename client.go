package brickognize

import (
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	workers    int
	retries    int
	limiter    *RateLimiter
}

type Option func(*Client)

func WithWorkers(n int) Option {
	return func(c *Client) { c.workers = n }
}

func WithRateLimit(rps int) Option {
	return func(c *Client) { c.limiter = NewRateLimiter(rps) }
}

func WithRetries(n int) Option {
	return func(c *Client) { c.retries = n }
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    "https://api.brickognize.com",
		httpClient: &http.Client{Timeout: 30 * time.Second},
		workers:    4,
		retries:    3,
		limiter:    NewRateLimiter(5),
	}

	for _, o := range opts {
		o(c)
	}
	return c
}
