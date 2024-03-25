package gocrunchybridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultUserAgent = "go-crunchybridge/sdk"
	defaultBaseURL   = "https://api.crunchybridge.com"
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

	return b, err
}

type service struct {
	client *Client
}

type Client struct {
	httpClient *http.Client
	userAgent  string
	apikey     APIKey

	Account     *AccountService
	AccessToken *AccessTokenService
	Certificate *CertificateService
	Cluster     *ClusterService
}

type APIKey string

func (a APIKey) String() string { return string(a) }

func New(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	srv := &service{client: c}

	c.Account = (*AccountService)(srv)
	c.AccessToken = (*AccessTokenService)(srv)
	c.Certificate = (*CertificateService)(srv)
	c.Cluster = (*ClusterService)(srv)

	return c, nil
}

func (c *Client) newRequest(method, resource string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(resource, "/") {
		return nil, errors.New("resource must contain a / prefix")
	}

	req, err := http.NewRequest(method, defaultBaseURL+resource, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+c.apikey.String())

	return req, nil
}

type Response struct {
	*http.Response
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*Response, error) {
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

	defer func() {
		_ = resp.Body.Close()
	}()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	case *string:
		var s strings.Builder
		_, err = io.Copy(&s, resp.Body)
		if err == nil {
			*v = s.String()
		}
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)

		switch decErr {
		case io.EOF:
			err = nil
		default:
			err = decErr
		}
	}

	if err != nil {
		return nil, err
	}

	return &Response{resp}, err
}
