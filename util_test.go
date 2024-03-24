//go:build integration
// +build integration

package gocrunchybridge

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getAPIKey(t *testing.T) APIKey {
	t.Helper()

	val := os.Getenv("INTEGRATION_API_KEY")
	require.NotEmpty(t, val)

	return APIKey(val)
}

func getTeamID(t *testing.T) string {
	t.Helper()

	val := os.Getenv("TEAM_ID")
	require.NotEmpty(t, val)

	return val
}
