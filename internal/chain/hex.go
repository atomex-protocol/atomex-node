package chain

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

// Hex -
type Hex string

// NewHexFromBytes -
func NewHexFromBytes(data []byte) Hex {
	return Hex(hex.EncodeToString(data))
}

// NewHexFromBytes32 -
func NewHexFromBytes32(data [32]byte) Hex {
	return Hex(hex.EncodeToString(data[:]))
}

// Bytes -
func (h Hex) Bytes() ([]byte, error) {
	return hex.DecodeString(string(h))
}

// Bytes32 -
func (h Hex) Bytes32() ([32]byte, error) {
	if len(h) != 64 {
		return [32]byte{}, errors.Errorf("invalid hex length %d for string %s", len(h), h)
	}
	data, err := h.Bytes()
	if err != nil {
		return [32]byte{}, err
	}
	var bytes32 [32]byte
	copy(bytes32[:], data)
	return bytes32, nil
}

// String -
func (h Hex) String() string {
	return string(h)
}
