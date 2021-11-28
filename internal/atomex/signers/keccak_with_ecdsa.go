package signers

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// KeccakWithEcdsa -
type KeccakWithEcdsa struct{}

// Generate -
func (KeccakWithEcdsa) Generate() (*Key, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &Key{
		Private: key.D.Bytes(),
		Public:  crypto.CompressPubkey(&key.PublicKey),
	}, nil
}

// Sign -
func (KeccakWithEcdsa) Sign(key *Key, msg []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(msg)
	signature, err := secp256k1.Sign(hash.Bytes(), key.Private)
	if err != nil {
		return nil, err
	}
	return signature[:64], nil
}

// Verify -
func (KeccakWithEcdsa) Verify(key *Key, msg, signature []byte) bool {
	hash := crypto.Keccak256Hash(msg)
	public, err := crypto.DecompressPubkey(key.Public)
	if err != nil {
		return false
	}
	return secp256k1.VerifySignature(crypto.FromECDSAPub(public), hash.Bytes(), signature)
}
