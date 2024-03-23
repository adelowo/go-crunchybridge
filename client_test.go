package gocrunchybridge

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToReader(t *testing.T) {
	tt := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "nil can be encoded",
			value: nil,
		},
		{
			name:  "noop request can be encoded",
			value: NoopRequestBody{},
		},
		{
			name: "structs can be encoded",
			value: struct {
				name string
			}{
				name: "Crunchybridge sdk",
			},
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			r, err := ToReader(v.value)
			require.NoError(t, err)

			_, err = io.Copy(io.Discard, r)
			require.NoError(t, err)
		})
	}
}
