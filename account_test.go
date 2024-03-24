//go:build integration
// +build integration

package gocrunchybridge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount_User(t *testing.T) {
	client, err := New(
		WithAPIKey(getAPIKey(t)),
	)

	require.NoError(t, err)

	usr, err := client.Account.User(context.Background())
	require.NoError(t, err)

	require.NotEmpty(t, usr.Email)
	require.NotEmpty(t, usr.ID)
}
