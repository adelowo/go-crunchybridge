//go:build integration
// +build integration

package gocrunchybridge

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAccessToken(t *testing.T) {
	client, err := New(
		WithAPIKey(getAPIKey(t)),
	)

	require.NoError(t, err)

	token, err := client.AccessToken.Create(context.Background(), &CreateAccessTokenOptions{
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	require.NoError(t, err)
	require.NotEmpty(t, token.AccessToken)

	// clean up the access token
	err = client.AccessToken.Delete(context.Background(), token)
	require.NoError(t, err)
}

func TestAccessToken_CreateValidation(t *testing.T) {
	client, err := New(
		WithAPIKey(getAPIKey(t)),
	)

	require.NoError(t, err)

	_, err = client.AccessToken.Create(context.Background(), &CreateAccessTokenOptions{
		ExpiresAt: time.Now(),
	})
	require.Error(t, err)
}

func TestAccessToken_DeleteValidation(t *testing.T) {
	client, err := New(
		WithAPIKey(getAPIKey(t)),
	)

	require.NoError(t, err)

	tt := []struct {
		name  string
		token AccessToken
	}{
		{
			name:  "token not provided",
			token: AccessToken{},
		},
		{
			name: "access token not provided",
			token: AccessToken{
				ID: "aghg4eyqx2kcjdquzenwvviaqi",
			},
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			err := client.AccessToken.Delete(context.Background(), v.token)
			require.Error(t, err)
			t.Log(err)
		})
	}
}
