package signers

import (
	"crypto/ed25519"

	"github.com/pkg/errors"
)

// Ed25519 -
type Ed25519 struct{}

// Generate -
func (Ed25519) Generate() (*Key, error) {
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, errors.Wrap(err, "Ed25519.Generate")
	}
	return &Key{
		Private: private,
		Public:  public,
	}, nil
}

// Sign -
func (Ed25519) Sign(key *Key, msg []byte) ([]byte, error) {
	return ed25519.Sign(key.Private, msg), nil
}

// Verify -
func (Ed25519) Verify(key *Key, msg, signature []byte) bool {
	return ed25519.Verify(key.Public, msg, signature)
}
