package signers

import (
	"crypto/ed25519"

	"github.com/goat-systems/go-tezos/v4/keys"
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
	private, err := keys.FromBytes(key.Private, keys.Ed25519)
	if err != nil {
		return nil, err
	}
	signature, err := private.SignBytes(msg)
	if err != nil {
		return nil, err
	}
	return signature.ToBytes(), nil
}

// Verify -
func (Ed25519Blake2b) Verify(key *Key, msg, signature []byte) bool {
	if msg != nil {
		if msg[0] != byte(3) {
			msg = append([]byte{3}, msg...)
		}
	}
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
