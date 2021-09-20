// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AtomexErc20MetaData contains all meta data concerning the AtomexErc20 contract.
var AtomexErc20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"Activated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Added\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_initiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_refundTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_countdown\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_payoff\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_active\",\"type\":\"bool\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_secret\",\"type\":\"bytes32\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swaps\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contractAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"refundTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"countdown\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"payoff\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"enumAtomexErc20Vault.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_refundTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_countdown\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_payoff\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_active\",\"type\":\"bool\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"add\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"activate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_secret\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AtomexErc20ABI is the input ABI used to generate the binding from.
// Deprecated: Use AtomexErc20MetaData.ABI instead.
var AtomexErc20ABI = AtomexErc20MetaData.ABI

// AtomexErc20 is an auto generated Go binding around an Ethereum contract.
type AtomexErc20 struct {
	AtomexErc20Caller     // Read-only binding to the contract
	AtomexErc20Transactor // Write-only binding to the contract
	AtomexErc20Filterer   // Log filterer for contract events
}

// AtomexErc20Caller is an auto generated read-only Go binding around an Ethereum contract.
type AtomexErc20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomexErc20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomexErc20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomexErc20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomexErc20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomexErc20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomexErc20Session struct {
	Contract     *AtomexErc20      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomexErc20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomexErc20CallerSession struct {
	Contract *AtomexErc20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AtomexErc20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomexErc20TransactorSession struct {
	Contract     *AtomexErc20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AtomexErc20Raw is an auto generated low-level Go binding around an Ethereum contract.
type AtomexErc20Raw struct {
	Contract *AtomexErc20 // Generic contract binding to access the raw methods on
}

// AtomexErc20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomexErc20CallerRaw struct {
	Contract *AtomexErc20Caller // Generic read-only contract binding to access the raw methods on
}

// AtomexErc20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomexErc20TransactorRaw struct {
	Contract *AtomexErc20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomexErc20 creates a new instance of AtomexErc20, bound to a specific deployed contract.
func NewAtomexErc20(address common.Address, backend bind.ContractBackend) (*AtomexErc20, error) {
	contract, err := bindAtomexErc20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20{AtomexErc20Caller: AtomexErc20Caller{contract: contract}, AtomexErc20Transactor: AtomexErc20Transactor{contract: contract}, AtomexErc20Filterer: AtomexErc20Filterer{contract: contract}}, nil
}

// NewAtomexErc20Caller creates a new read-only instance of AtomexErc20, bound to a specific deployed contract.
func NewAtomexErc20Caller(address common.Address, caller bind.ContractCaller) (*AtomexErc20Caller, error) {
	contract, err := bindAtomexErc20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20Caller{contract: contract}, nil
}

// NewAtomexErc20Transactor creates a new write-only instance of AtomexErc20, bound to a specific deployed contract.
func NewAtomexErc20Transactor(address common.Address, transactor bind.ContractTransactor) (*AtomexErc20Transactor, error) {
	contract, err := bindAtomexErc20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20Transactor{contract: contract}, nil
}

// NewAtomexErc20Filterer creates a new log filterer instance of AtomexErc20, bound to a specific deployed contract.
func NewAtomexErc20Filterer(address common.Address, filterer bind.ContractFilterer) (*AtomexErc20Filterer, error) {
	contract, err := bindAtomexErc20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20Filterer{contract: contract}, nil
}

// bindAtomexErc20 binds a generic wrapper to an already deployed contract.
func bindAtomexErc20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomexErc20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomexErc20 *AtomexErc20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomexErc20.Contract.AtomexErc20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomexErc20 *AtomexErc20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomexErc20.Contract.AtomexErc20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomexErc20 *AtomexErc20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomexErc20.Contract.AtomexErc20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomexErc20 *AtomexErc20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomexErc20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomexErc20 *AtomexErc20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomexErc20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomexErc20 *AtomexErc20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomexErc20.Contract.contract.Transact(opts, method, params...)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(bytes32 hashedSecret, address contractAddr, address participant, address initiator, uint256 refundTimestamp, uint256 countdown, uint256 value, uint256 payoff, bool active, uint8 state)
func (_AtomexErc20 *AtomexErc20Caller) Swaps(opts *bind.CallOpts, arg0 [32]byte) (struct {
	HashedSecret    [32]byte
	ContractAddr    common.Address
	Participant     common.Address
	Initiator       common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	State           uint8
}, error) {
	var out []interface{}
	err := _AtomexErc20.contract.Call(opts, &out, "swaps", arg0)

	outstruct := new(struct {
		HashedSecret    [32]byte
		ContractAddr    common.Address
		Participant     common.Address
		Initiator       common.Address
		RefundTimestamp *big.Int
		Countdown       *big.Int
		Value           *big.Int
		Payoff          *big.Int
		Active          bool
		State           uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.HashedSecret = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.ContractAddr = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Participant = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Initiator = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.RefundTimestamp = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Countdown = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Value = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Payoff = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[8], new(bool)).(*bool)
	outstruct.State = *abi.ConvertType(out[9], new(uint8)).(*uint8)

	return *outstruct, err

}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(bytes32 hashedSecret, address contractAddr, address participant, address initiator, uint256 refundTimestamp, uint256 countdown, uint256 value, uint256 payoff, bool active, uint8 state)
func (_AtomexErc20 *AtomexErc20Session) Swaps(arg0 [32]byte) (struct {
	HashedSecret    [32]byte
	ContractAddr    common.Address
	Participant     common.Address
	Initiator       common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	State           uint8
}, error) {
	return _AtomexErc20.Contract.Swaps(&_AtomexErc20.CallOpts, arg0)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(bytes32 hashedSecret, address contractAddr, address participant, address initiator, uint256 refundTimestamp, uint256 countdown, uint256 value, uint256 payoff, bool active, uint8 state)
func (_AtomexErc20 *AtomexErc20CallerSession) Swaps(arg0 [32]byte) (struct {
	HashedSecret    [32]byte
	ContractAddr    common.Address
	Participant     common.Address
	Initiator       common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	State           uint8
}, error) {
	return _AtomexErc20.Contract.Swaps(&_AtomexErc20.CallOpts, arg0)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _hashedSecret) returns()
func (_AtomexErc20 *AtomexErc20Transactor) Activate(opts *bind.TransactOpts, _hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.contract.Transact(opts, "activate", _hashedSecret)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _hashedSecret) returns()
func (_AtomexErc20 *AtomexErc20Session) Activate(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Activate(&_AtomexErc20.TransactOpts, _hashedSecret)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _hashedSecret) returns()
func (_AtomexErc20 *AtomexErc20TransactorSession) Activate(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Activate(&_AtomexErc20.TransactOpts, _hashedSecret)
}

// Add is a paid mutator transaction binding the contract method 0x5ffa4bce.
//
// Solidity: function add(bytes32 _hashedSecret, uint256 _value) returns()
func (_AtomexErc20 *AtomexErc20Transactor) Add(opts *bind.TransactOpts, _hashedSecret [32]byte, _value *big.Int) (*types.Transaction, error) {
	return _AtomexErc20.contract.Transact(opts, "add", _hashedSecret, _value)
}

// Add is a paid mutator transaction binding the contract method 0x5ffa4bce.
//
// Solidity: function add(bytes32 _hashedSecret, uint256 _value) returns()
func (_AtomexErc20 *AtomexErc20Session) Add(_hashedSecret [32]byte, _value *big.Int) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Add(&_AtomexErc20.TransactOpts, _hashedSecret, _value)
}

// Add is a paid mutator transaction binding the contract method 0x5ffa4bce.
//
// Solidity: function add(bytes32 _hashedSecret, uint256 _value) returns()
func (_AtomexErc20 *AtomexErc20TransactorSession) Add(_hashedSecret [32]byte, _value *big.Int) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Add(&_AtomexErc20.TransactOpts, _hashedSecret, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x6170a610.
//
// Solidity: function initiate(bytes32 _hashedSecret, address _contract, address _participant, uint256 _refundTimestamp, uint256 _countdown, uint256 _value, uint256 _payoff, bool _active) returns()
func (_AtomexErc20 *AtomexErc20Transactor) Initiate(opts *bind.TransactOpts, _hashedSecret [32]byte, _contract common.Address, _participant common.Address, _refundTimestamp *big.Int, _countdown *big.Int, _value *big.Int, _payoff *big.Int, _active bool) (*types.Transaction, error) {
	return _AtomexErc20.contract.Transact(opts, "initiate", _hashedSecret, _contract, _participant, _refundTimestamp, _countdown, _value, _payoff, _active)
}

// Initiate is a paid mutator transaction binding the contract method 0x6170a610.
//
// Solidity: function initiate(bytes32 _hashedSecret, address _contract, address _participant, uint256 _refundTimestamp, uint256 _countdown, uint256 _value, uint256 _payoff, bool _active) returns()
func (_AtomexErc20 *AtomexErc20Session) Initiate(_hashedSecret [32]byte, _contract common.Address, _participant common.Address, _refundTimestamp *big.Int, _countdown *big.Int, _value *big.Int, _payoff *big.Int, _active bool) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Initiate(&_AtomexErc20.TransactOpts, _hashedSecret, _contract, _participant, _refundTimestamp, _countdown, _value, _payoff, _active)
}

// Initiate is a paid mutator transaction binding the contract method 0x6170a610.
//
// Solidity: function initiate(bytes32 _hashedSecret, address _contract, address _participant, uint256 _refundTimestamp, uint256 _countdown, uint256 _value, uint256 _payoff, bool _active) returns()
func (_AtomexErc20 *AtomexErc20TransactorSession) Initiate(_hashedSecret [32]byte, _contract common.Address, _participant common.Address, _refundTimestamp *big.Int, _countdown *big.Int, _value *big.Int, _payoff *big.Int, _active bool) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Initiate(&_AtomexErc20.TransactOpts, _hashedSecret, _contract, _participant, _refundTimestamp, _countdown, _value, _payoff, _active)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 _hashedSecret, bytes32 _secret) returns()
func (_AtomexErc20 *AtomexErc20Transactor) Redeem(opts *bind.TransactOpts, _hashedSecret [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.contract.Transact(opts, "redeem", _hashedSecret, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 _hashedSecret, bytes32 _secret) returns()
func (_AtomexErc20 *AtomexErc20Session) Redeem(_hashedSecret [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Redeem(&_AtomexErc20.TransactOpts, _hashedSecret, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 _hashedSecret, bytes32 _secret) returns()
func (_AtomexErc20 *AtomexErc20TransactorSession) Redeem(_hashedSecret [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Redeem(&_AtomexErc20.TransactOpts, _hashedSecret, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _hashedSecret) returns()
func (_AtomexErc20 *AtomexErc20Transactor) Refund(opts *bind.TransactOpts, _hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.contract.Transact(opts, "refund", _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _hashedSecret) returns()
func (_AtomexErc20 *AtomexErc20Session) Refund(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Refund(&_AtomexErc20.TransactOpts, _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _hashedSecret) returns()
func (_AtomexErc20 *AtomexErc20TransactorSession) Refund(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexErc20.Contract.Refund(&_AtomexErc20.TransactOpts, _hashedSecret)
}

// AtomexErc20ActivatedIterator is returned from FilterActivated and is used to iterate over the raw logs and unpacked data for Activated events raised by the AtomexErc20 contract.
type AtomexErc20ActivatedIterator struct {
	Event *AtomexErc20Activated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AtomexErc20ActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexErc20Activated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AtomexErc20Activated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AtomexErc20ActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexErc20ActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexErc20Activated represents a Activated event raised by the AtomexErc20 contract.
type AtomexErc20Activated struct {
	HashedSecret [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterActivated is a free log retrieval operation binding the contract event 0xe1abfe35306def8dbc83e3cb0bc76ffd144cee4ab7707b4e888afd4d24c2d6ca.
//
// Solidity: event Activated(bytes32 indexed _hashedSecret)
func (_AtomexErc20 *AtomexErc20Filterer) FilterActivated(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexErc20ActivatedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.FilterLogs(opts, "Activated", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20ActivatedIterator{contract: _AtomexErc20.contract, event: "Activated", logs: logs, sub: sub}, nil
}

// WatchActivated is a free log subscription operation binding the contract event 0xe1abfe35306def8dbc83e3cb0bc76ffd144cee4ab7707b4e888afd4d24c2d6ca.
//
// Solidity: event Activated(bytes32 indexed _hashedSecret)
func (_AtomexErc20 *AtomexErc20Filterer) WatchActivated(opts *bind.WatchOpts, sink chan<- *AtomexErc20Activated, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.WatchLogs(opts, "Activated", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexErc20Activated)
				if err := _AtomexErc20.contract.UnpackLog(event, "Activated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseActivated is a log parse operation binding the contract event 0xe1abfe35306def8dbc83e3cb0bc76ffd144cee4ab7707b4e888afd4d24c2d6ca.
//
// Solidity: event Activated(bytes32 indexed _hashedSecret)
func (_AtomexErc20 *AtomexErc20Filterer) ParseActivated(log types.Log) (*AtomexErc20Activated, error) {
	event := new(AtomexErc20Activated)
	if err := _AtomexErc20.contract.UnpackLog(event, "Activated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexErc20AddedIterator is returned from FilterAdded and is used to iterate over the raw logs and unpacked data for Added events raised by the AtomexErc20 contract.
type AtomexErc20AddedIterator struct {
	Event *AtomexErc20Added // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AtomexErc20AddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexErc20Added)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AtomexErc20Added)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AtomexErc20AddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexErc20AddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexErc20Added represents a Added event raised by the AtomexErc20 contract.
type AtomexErc20Added struct {
	HashedSecret [32]byte
	Sender       common.Address
	Value        *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAdded is a free log retrieval operation binding the contract event 0xd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423.
//
// Solidity: event Added(bytes32 indexed _hashedSecret, address _sender, uint256 _value)
func (_AtomexErc20 *AtomexErc20Filterer) FilterAdded(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexErc20AddedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.FilterLogs(opts, "Added", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20AddedIterator{contract: _AtomexErc20.contract, event: "Added", logs: logs, sub: sub}, nil
}

// WatchAdded is a free log subscription operation binding the contract event 0xd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423.
//
// Solidity: event Added(bytes32 indexed _hashedSecret, address _sender, uint256 _value)
func (_AtomexErc20 *AtomexErc20Filterer) WatchAdded(opts *bind.WatchOpts, sink chan<- *AtomexErc20Added, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.WatchLogs(opts, "Added", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexErc20Added)
				if err := _AtomexErc20.contract.UnpackLog(event, "Added", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdded is a log parse operation binding the contract event 0xd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423.
//
// Solidity: event Added(bytes32 indexed _hashedSecret, address _sender, uint256 _value)
func (_AtomexErc20 *AtomexErc20Filterer) ParseAdded(log types.Log) (*AtomexErc20Added, error) {
	event := new(AtomexErc20Added)
	if err := _AtomexErc20.contract.UnpackLog(event, "Added", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexErc20InitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the AtomexErc20 contract.
type AtomexErc20InitiatedIterator struct {
	Event *AtomexErc20Initiated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AtomexErc20InitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexErc20Initiated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AtomexErc20Initiated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AtomexErc20InitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexErc20InitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexErc20Initiated represents a Initiated event raised by the AtomexErc20 contract.
type AtomexErc20Initiated struct {
	HashedSecret    [32]byte
	Contract        common.Address
	Participant     common.Address
	Initiator       common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x99cdc76be187c2919cca1f8a27dac6a651692095c8902fadcf1fc75539d28146.
//
// Solidity: event Initiated(bytes32 indexed _hashedSecret, address indexed _contract, address indexed _participant, address _initiator, uint256 _refundTimestamp, uint256 _countdown, uint256 _value, uint256 _payoff, bool _active)
func (_AtomexErc20 *AtomexErc20Filterer) FilterInitiated(opts *bind.FilterOpts, _hashedSecret [][32]byte, _contract []common.Address, _participant []common.Address) (*AtomexErc20InitiatedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}
	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}
	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _AtomexErc20.contract.FilterLogs(opts, "Initiated", _hashedSecretRule, _contractRule, _participantRule)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20InitiatedIterator{contract: _AtomexErc20.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x99cdc76be187c2919cca1f8a27dac6a651692095c8902fadcf1fc75539d28146.
//
// Solidity: event Initiated(bytes32 indexed _hashedSecret, address indexed _contract, address indexed _participant, address _initiator, uint256 _refundTimestamp, uint256 _countdown, uint256 _value, uint256 _payoff, bool _active)
func (_AtomexErc20 *AtomexErc20Filterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *AtomexErc20Initiated, _hashedSecret [][32]byte, _contract []common.Address, _participant []common.Address) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}
	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}
	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _AtomexErc20.contract.WatchLogs(opts, "Initiated", _hashedSecretRule, _contractRule, _participantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexErc20Initiated)
				if err := _AtomexErc20.contract.UnpackLog(event, "Initiated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitiated is a log parse operation binding the contract event 0x99cdc76be187c2919cca1f8a27dac6a651692095c8902fadcf1fc75539d28146.
//
// Solidity: event Initiated(bytes32 indexed _hashedSecret, address indexed _contract, address indexed _participant, address _initiator, uint256 _refundTimestamp, uint256 _countdown, uint256 _value, uint256 _payoff, bool _active)
func (_AtomexErc20 *AtomexErc20Filterer) ParseInitiated(log types.Log) (*AtomexErc20Initiated, error) {
	event := new(AtomexErc20Initiated)
	if err := _AtomexErc20.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexErc20RedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the AtomexErc20 contract.
type AtomexErc20RedeemedIterator struct {
	Event *AtomexErc20Redeemed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AtomexErc20RedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexErc20Redeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AtomexErc20Redeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AtomexErc20RedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexErc20RedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexErc20Redeemed represents a Redeemed event raised by the AtomexErc20 contract.
type AtomexErc20Redeemed struct {
	HashedSecret [32]byte
	Secret       [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab.
//
// Solidity: event Redeemed(bytes32 indexed _hashedSecret, bytes32 _secret)
func (_AtomexErc20 *AtomexErc20Filterer) FilterRedeemed(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexErc20RedeemedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.FilterLogs(opts, "Redeemed", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20RedeemedIterator{contract: _AtomexErc20.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab.
//
// Solidity: event Redeemed(bytes32 indexed _hashedSecret, bytes32 _secret)
func (_AtomexErc20 *AtomexErc20Filterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *AtomexErc20Redeemed, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.WatchLogs(opts, "Redeemed", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexErc20Redeemed)
				if err := _AtomexErc20.contract.UnpackLog(event, "Redeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemed is a log parse operation binding the contract event 0x489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab.
//
// Solidity: event Redeemed(bytes32 indexed _hashedSecret, bytes32 _secret)
func (_AtomexErc20 *AtomexErc20Filterer) ParseRedeemed(log types.Log) (*AtomexErc20Redeemed, error) {
	event := new(AtomexErc20Redeemed)
	if err := _AtomexErc20.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexErc20RefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the AtomexErc20 contract.
type AtomexErc20RefundedIterator struct {
	Event *AtomexErc20Refunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AtomexErc20RefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexErc20Refunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AtomexErc20Refunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AtomexErc20RefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexErc20RefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexErc20Refunded represents a Refunded event raised by the AtomexErc20 contract.
type AtomexErc20Refunded struct {
	HashedSecret [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed _hashedSecret)
func (_AtomexErc20 *AtomexErc20Filterer) FilterRefunded(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexErc20RefundedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.FilterLogs(opts, "Refunded", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexErc20RefundedIterator{contract: _AtomexErc20.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed _hashedSecret)
func (_AtomexErc20 *AtomexErc20Filterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *AtomexErc20Refunded, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexErc20.contract.WatchLogs(opts, "Refunded", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexErc20Refunded)
				if err := _AtomexErc20.contract.UnpackLog(event, "Refunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefunded is a log parse operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed _hashedSecret)
func (_AtomexErc20 *AtomexErc20Filterer) ParseRefunded(log types.Log) (*AtomexErc20Refunded, error) {
	event := new(AtomexErc20Refunded)
	if err := _AtomexErc20.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
