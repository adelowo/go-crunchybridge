package gocrunchybridge

import "net/http"

type Option func(*Client)

func WithAPIKey(a APIKey) Option {
	return func(c *Client) {
		c.apikey = a
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithUserAgent(s string) Option {
	return func(c *Client) {
		c.userAgent = s
	}
}
