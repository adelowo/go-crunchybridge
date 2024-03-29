//go:build integration
// +build integration

package gocrunchybridge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func getClient(t *testing.T) *Client {
	c, err := New(WithAPIKey(getAPIKey(t)))
	require.NoError(t, err)

	return c
}

func TestCluster(t *testing.T) {
	client := getClient(t)

	cluster, err := client.Cluster.Create(context.Background(),
		&gocrunchybridge.CreateClusterOptions{
			PlanID:            "hobby-0",
			TeamID:            getTeamID(t),
			RegionID:          "eu-west-2",
			ProviderID:        gocrunchybridge.ClusterProviderAws,
			StorageSize:       10,
			Environment:       gocrunchybridge.ClusterEnvironmentProduction,
			HighlyAvailable:   false,
			PostgresVersionID: 16,
		})

	require.NoError(t, err)

	err = client.Cluster.Delete(context.Background(), cluster.ID)
	require.NoError(t, err)
}
