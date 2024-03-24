//go:build integration
// +build integration

package gocrunchybridge

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAccessToken_Create(t *testing.T) {
	client, err := New(
		WithAPIKey(getAPIKey(t)),
	)

	require.NoError(t, err)

	token, err := client.AccessToken.Create(context.Background(), &CreateAccessTokenOptions{
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	require.NoError(t, err)
	require.NotEmpty(t, token.AccessToken)
}
