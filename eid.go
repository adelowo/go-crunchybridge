package gocrunchybridge

import (
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
)

var cachedBase32 = base32.StdEncoding.WithPadding(base32.NoPadding)

// EID represents an encoded ID.
type EID string

func (e EID) IsValid() error {
	if len(string(e)) != 26 {
		return fmt.Errorf("EID string should be exactly 26 characters long")
	}

	decoded, err := cachedBase32.DecodeString(strings.ToUpper(string(e)))
	if err != nil {
		return fmt.Errorf("failed to decode EID: %w", err)
	}

	dst := [16]byte{}
	copy(dst[:], decoded)
	if strings.ToLower(cachedBase32.EncodeToString(dst[:])) != e.String() {
		return errors.New("string does not match EID format")
	}

	return nil
}

func (e EID) String() string { return strings.ToLower(string(e)) }
