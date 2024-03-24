package gocrunchybridge

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEID(t *testing.T) {
	tt := []struct {
		name     string
		hasError bool
		value    string
	}{
		{
			name:     "valid EID",
			hasError: false,
			value:    "zar5556utjb3hkevt5dkxj2o4i",
		},
		{
			name:     "invalid EID",
			hasError: true,
			value:    "aocqkyenaddjfojpfjpfjpjncbctffvfkmcxk",
		},
		{
			name:     "invalid EID",
			hasError: true,
			value:    "aocqkyenaddjncbctffvfkmcxk",
		},
		{
			name:     "uuid is invalid",
			hasError: true,
			value:    "4b866d87-a12e-4f8b-afd9-14408575f738",
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			eid := EID(v.value)
			err := eid.IsValid()
			if v.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
