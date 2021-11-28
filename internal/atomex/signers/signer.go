package signers

import "github.com/pkg/errors"

// Signer -
type Signer interface {
	Generate() (*Key, error)
	Sign(key *Key, msg []byte) ([]byte, error)
	Verify(key *Key, msg, signature []byte) bool
}

// Key -
type Key struct {
	Private []byte
	Public  []byte
}

// errors
var (
	ErrNotImplemented     = errors.New("not implemented")
	ErrVerificationFailed = errors.New("signature verification failed")
)

// algorithms
const (
	AlgorithmEd25519                    = "Ed25519"
	AlgorithmEd25519Blake2b             = "Ed25519:Blake2b"
	AlgorithmSha256WithEcdsaSecp256k1   = "Sha256WithEcdsa:Secp256k1"
	AlgorithmBlake2bWithEcdsaSecp256k1  = "Blake2bWithEcdsa:Secp256k1"
	AlgorithmBlake2bWithEcdsaSecp256r1  = "Blake2bWithEcdsa:Secp256r1"
	AlgorithmBlake2bWithEcdsaBtcMsg     = "Blake2bWithEcdsa:BtcMsg"
	AlgorithmKeccak256WithEcdsaGeth2940 = "Keccak256WithEcdsa:Geth2940"
)

// Generate -
func Generate(algo string) (*Key, error) {
	signer, err := Get(algo)
	if err != nil {
		return nil, err
	}
	return signer.Generate()
}

// Sign -
func Sign(algo string, key *Key, msg []byte, verify bool) ([]byte, error) {
	signer, err := Get(algo)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(key, msg)
	if err != nil {
		return nil, err
	}

	if !verify {
		return signature, nil
	}

	ok, err := Verify(algo, key, msg, signature)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = ErrVerificationFailed
	}
	return signature, err
}

// Sign -
func Verify(algo string, key *Key, msg, signature []byte) (bool, error) {
	signer, err := Get(algo)
	if err != nil {
		return false, err
	}

	return signer.Verify(key, msg, signature), nil
}

// Get -
func Get(algo string) (Signer, error) {
	switch algo {
	case AlgorithmEd25519:
		return Ed25519{}, nil
	case AlgorithmEd25519Blake2b:
		return Ed25519Blake2b{}, nil
	case AlgorithmKeccak256WithEcdsaGeth2940:
		return KeccakWithEcdsa{}, nil
	case AlgorithmBlake2bWithEcdsaSecp256k1:
		return Blake2bWithEcdsaSecp256k1{}, nil
	default:
		return nil, errors.Wrap(ErrNotImplemented, algo)
	}
}
