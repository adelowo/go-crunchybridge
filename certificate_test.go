//go:build integration
// +build integration

package gocrunchybridge

import (
	"context"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCertificate_Get(t *testing.T) {
	client, err := New(
		WithAPIKey(getAPIKey(t)),
	)

	require.NoError(t, err)

	pemContent, err := client.Certificate.Get(context.Background(), FetchTeamCertificateOption{
		TeamID: getTeamID(t),
	})
	require.NoError(t, err)

	block, _ := pem.Decode([]byte(pemContent))
	require.NotNil(t, block)
	require.Equal(t, block.Type, "CERTIFICATE")
}
