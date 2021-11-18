package signers

import (
	"crypto/ed25519"

	"github.com/goat-systems/go-tezos/v4/keys"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

// Ed25519Blake2b -
type Ed25519Blake2b struct{}

// Generate -
func (Ed25519Blake2b) Generate() (*Key, error) {
	key, err := keys.Generate(keys.Ed25519)
	if err != nil {
		return nil, err
	}

	return &Key{
		Private: key.GetBytes(),
		Public:  key.PubKey.GetBytes(),
	}, nil
}

// Sign -
func (Ed25519Blake2b) Sign(key *Key, msg []byte) ([]byte, error) {
	hash, err := blake2b.New(32, []byte{})
	if err != nil {
		return nil, err
	}

	i, err := hash.Write(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign operation bytes")
	}
	if i != len(msg) {
		return nil, errors.Errorf("failed to sign operation: generic hash length %d does not match bytes length %d", i, len(msg))
	}

	return ed25519.Sign(key.Private, hash.Sum([]byte{})), nil
}

// Verify -
func (Ed25519Blake2b) Verify(key *Key, msg, signature []byte) bool {
	hash, err := blake2b.New(32, []byte{})
	if err != nil {
		return false
	}
	i, err := hash.Write(msg)
	if err != nil {
		return false
	}
	if i != len(msg) {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(key.Public), hash.Sum([]byte{}), signature)
}
