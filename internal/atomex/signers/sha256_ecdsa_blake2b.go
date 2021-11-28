package signers

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/goat-systems/go-tezos/v4/keys"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

// Blake2bWithEcdsaSecp256k1 -
type Blake2bWithEcdsaSecp256k1 struct{}

// Generate -
func (Blake2bWithEcdsaSecp256k1) Generate() (*Key, error) {
	key, err := keys.Generate(keys.Secp256k1)
	if err != nil {
		return nil, err
	}

	return &Key{
		Private: key.GetBytes(),
		Public:  key.PubKey.GetBytes(),
	}, nil
}

// Sign -
func (s Blake2bWithEcdsaSecp256k1) Sign(key *Key, msg []byte) ([]byte, error) {
	private, err := crypto.ToECDSA(key.Private)
	if err != nil {
		return nil, err
	}
	return s.sign(private, msg)
}

// Verify -
func (Blake2bWithEcdsaSecp256k1) Verify(key *Key, msg, signature []byte) bool {
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
	return secp256k1.VerifySignature(ed25519.PublicKey(key.Public), hash.Sum([]byte{}), signature)
}

func (Blake2bWithEcdsaSecp256k1) sign(private *ecdsa.PrivateKey, msg []byte) ([]byte, error) {
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

	r, ss, err := ecdsa.Sign(rand.Reader, private, hash.Sum([]byte{}))
	if err != nil {
		return nil, err
	}

	if ss.Cmp(maxS()) > 0 {
		ss = big.NewInt(0).Sub(order(), ss)
	}

	signature := append(r.Bytes(), ss.Bytes()...)
	return signature, nil
}
