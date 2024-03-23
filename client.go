package gocrunchybridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultUserAgent = "go-crunchybridge/sdk"
)

var errNonNilContext = errors.New("context must be non-nil")

// NoopRequestBody is used to ide
type NoopRequestBody struct{}

// ToReader converts any struct into a io#Reader that can be used
func ToReader[T NoopRequestBody | any](t T) (io.Reader, error) {
	b := bytes.NewBuffer(nil)

	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)
	if err != nil {
		return nil, err
	}

	return b, nil
}

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

func (c *Client) newRequest(method, resource string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, resource, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			return nil, e
		}

		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		// TODO(adelowo): allow users to be able to make sense of the error message instead
		return nil, errors.New("unexpected status code")
	}

	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)

		switch decErr {
		case io.EOF:
			decErr = nil
		default:
			err = decErr
		}
	}

	if err != nil {
		return nil, err
	}

	return resp, err
}
