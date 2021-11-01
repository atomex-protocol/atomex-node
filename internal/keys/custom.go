package keys

import (
	"encoding/json"
	"os"

	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/ebellocchia/go-base58"
)

// Custom -
type Custom struct {
	algo string
}

// NewCustom -
func NewCustom(algo string) *Custom {
	return &Custom{algo}
}

// NewCustomsBlake2bWithEcdsaSecp256k1 -
func NewCustomsBlake2bWithEcdsaSecp256k1() *Custom {
	return NewCustom(signers.AlgorithmBlake2bWithEcdsaSecp256k1)
}

// Get -
func (c *Custom) Get(filepath string) (*signers.Key, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var key fileStructure
	if err := json.NewDecoder(f).Decode(&key); err != nil {
		return nil, err
	}

	decoder := base58.New(base58.AlphabetBitcoin)
	public, err := decoder.Decode(key.Public)
	if err != nil {
		return nil, err
	}
	private, err := decoder.Decode(key.Private)
	if err != nil {
		return nil, err
	}
	return &signers.Key{
		Public:  public,
		Private: private,
	}, nil
}

// Create -
func (c *Custom) Create(filepath string) (*signers.Key, error) {
	generated, err := signers.Generate(c.algo)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	encoder := base58.New(base58.AlphabetBitcoin)
	return generated, json.NewEncoder(f).Encode(fileStructure{
		Public:  encoder.Encode(generated.Public),
		Private: encoder.Encode(generated.Private),
	})
}

type fileStructure struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}
