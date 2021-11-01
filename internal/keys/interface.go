package keys

import (
	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/pkg/errors"
)

// Storage -
type Storage interface {
	Get(filepath string) (*signers.Key, error)
	Create(filepath string) (*signers.Key, error)
}

// New -
func New(kind StorageKind) (Storage, error) {
	switch kind {
	case StorageKindCustom:
		return NewCustomsBlake2bWithEcdsaSecp256k1(), nil
	default:
		return nil, errors.Errorf("unknown key storage kind: %s", kind)
	}
}

// StorageKind -
type StorageKind string

// kinds
const (
	StorageKindCustom = "custom"
)
