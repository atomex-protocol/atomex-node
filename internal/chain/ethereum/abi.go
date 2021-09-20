package ethereum

import (
	"math/big"
	"os"

	abi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var abiAtomexEth *abi.ABI
var abiAtomexErc20 *abi.ABI

// LoadAbi -
func LoadAbi() error {
	if abiAtomexEth == nil {
		ethABI, err := os.Open("abi/AtomexEthVault.json")
		if err != nil {
			return err
		}
		a, err := abi.JSON(ethABI)
		if err != nil {
			return err
		}
		abiAtomexEth = &a
	}

	if abiAtomexErc20 == nil {
		erc20ABI, err := os.Open("abi/AtomexErc20Vault.json")
		if err != nil {
			return err
		}
		a, err := abi.JSON(erc20ABI)
		if err != nil {
			return err
		}
		abiAtomexErc20 = &a
	}

	return nil
}

type initiatedEthArgs struct {
	Initiator       common.Address `abi:"_initiator"`
	RefundTimestamp *big.Int       `abi:"_refundTimestamp"`
	Value           *big.Int       `abi:"_value"`
	PayOff          *big.Int       `abi:"_payoff"`
}

type initiatedErc20Args struct {
	Initiator       common.Address `abi:"_initiator"`
	RefundTimestamp *big.Int       `abi:"_refundTimestamp"`
	Countdown       *big.Int       `abi:"_countdown"`
	Value           *big.Int       `abi:"_value"`
	PayOff          *big.Int       `abi:"_payoff"`
	Active          bool           `abi:"_active"`
}

type redeemedArgs struct {
	Secret [32]byte `abi:"_secret"`
}

// events
const (
	EventActivated = "Activated"
	EventAdded     = "Added"
	EventInitiated = "Initiated"
	EventRedeemed  = "Redeemed"
	EventRefunded  = "Refunded"
)

const (
	ContractTypeEth   = "eth"
	ContractTypeErc20 = "erc20"
)
