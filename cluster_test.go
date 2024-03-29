//go:build integration
// +build integration

package gocrunchybridge

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getClient(t *testing.T) *Client {
	c, err := New(WithAPIKey(getAPIKey(t)))
	require.NoError(t, err)

	return c
}

func TestCluster_Create(t *testing.T) {
	client := getClient(t)
}
