package gocrunchybridge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/adelowo/go-crunchybridge/internal/util"
)

type AccessTokenService service

type CreateAccessTokenOptions struct {
	clientSecret string
	ExpiresAt    time.Time
}

func (c CreateAccessTokenOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ExpiresAt    string `json:"expires_at"`
		ClientSecret string `json:"client_secret"`
	}{
		ClientSecret: c.clientSecret,
		ExpiresAt:    c.ExpiresAt.Format(time.RFC3339),
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

	if time.Until(opts.ExpiresAt) <= time.Minute {
		return resp, errors.New("the expiration date must be more than 1 minute")
	}

	opts.clientSecret = a.client.apikey.String()

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

func (a *AccessTokenService) Delete(ctx context.Context,
	token AccessToken,
) error {
	if util.IsStringEmpty(token.ID) {
		return errors.New("the id of the token to delete must be provided")
	}

	if util.IsStringEmpty(token.AccessToken) {
		return errors.New(`access token to be deleted must be provided as crunchybridge requires a token can only be invalidated by itself`)
	}

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	req, err := a.client.newRequest(http.MethodDelete,
		"/access-tokens/"+token.ID, body)
	if err != nil {
		return err
	}

	// update the Authorization token as it is a core requirement to
	// make sure the api call is authenticated by the token to be invalidated
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	_, err = a.client.Do(ctx, req, nil)
	return err
}
