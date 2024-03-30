//go:build integration
// +build integration

package gocrunchybridge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	client := getClient(t)

	plan, err := client.Provider.Get(context.Background(), FetchProviderOptions{
		Provider: ClusterProviderGcp,
	})

	require.NoError(t, err)
	require.True(t, len(plan.Plans) > 0)
}
