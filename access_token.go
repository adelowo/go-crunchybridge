package gocrunchybridge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type AccessTokenService service

type CreateAccessTokenOptions struct {
	ClientSecret string    `json:"client_secret,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

func (c CreateAccessTokenOptions) MarshalJSON() ([]byte, error) {
	type Alias CreateAccessTokenOptions

	return json.Marshal(&struct {
		ExpiresAt string `json:"expires_at"`
		Alias
	}{
		Alias:     (Alias)(c),
		ExpiresAt: c.ExpiresAt.Format(time.RFC3339),
	})
}

type AccessToken struct {
	ID          string    `json:"id,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	AccountID   string    `json:"account_id,omitempty"`
	APIKeyID    string    `json:"api_key_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
}

func (a *AccessTokenService) Create(ctx context.Context,
	opts *CreateAccessTokenOptions,
) (AccessToken, error) {
	var resp AccessToken

	if opts.ExpiresAt.Sub(time.Now()) <= time.Minute {
		return resp, errors.New("the expiration date must be more than 1 minute")
	}

	opts.ClientSecret = a.client.apikey.String()

	body, err := ToReader(opts)
	if err != nil {
		return resp, err
	}

	req, err := a.client.newRequest(http.MethodPost, "/access-tokens", body)
	if err != nil {
		return resp, err
	}

	_, err = a.client.Do(ctx, req, &resp)
	return resp, err
}
