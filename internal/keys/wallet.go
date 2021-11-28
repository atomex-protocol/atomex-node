package keys

import (
	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goat-systems/go-tezos/v4/keys"
)

// Wallet -
type Wallet struct {
	chain.Wallet
}

// NewTezosWallet -
func NewTezosWallet() (*Wallet, error) {
	secret, err := chain.LoadSecret("TEZOS_PRIVATE")
	if err != nil {
		return nil, err
	}

	key, err := keys.FromBase58(secret, keys.Ed25519)
	if err != nil {
		return nil, err
	}

	return &Wallet{chain.Wallet{
		Address:   key.PubKey.GetAddress(),
		PublicKey: key.PubKey.GetBytes(),
		Private:   key.GetBytes(),
	}}, nil
}

// NewEthereumWallet -
func NewEthereumWallet() (*Wallet, error) {
	secret, err := chain.LoadSecret("ETHEREUM_PRIVATE")
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return nil, err
	}

	return &Wallet{chain.Wallet{
		Address:   crypto.PubkeyToAddress(privateKey.PublicKey).String(),
		Private:   crypto.FromECDSA(privateKey),
		PublicKey: crypto.CompressPubkey(&privateKey.PublicKey),
	}}, nil
}

// Get -
func (w *Wallet) Get(filepath string) (*signers.Key, error) {
	return &signers.Key{
		Public:  w.PublicKey,
		Private: w.Private,
	}, nil
}

// Create -
func (w *Wallet) Create(filepath string) (*signers.Key, error) {
	return nil, nil
}
