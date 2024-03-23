package gocrunchybridge

import "net/http"

const (
	defaultUserAgent = "go-crunchybridge/sdk"
)

type Client struct {
	httpClient *http.Client
	userAgent  string
}

type APIKey string

func New(opts ...Option) (*Client, error) {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}
