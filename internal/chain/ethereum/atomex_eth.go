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

// AtomexEthMetaData contains all meta data concerning the AtomexEth contract.
var AtomexEthMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"Activated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Added\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_initiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_refundTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_payoff\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_secret\",\"type\":\"bytes32\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swaps\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"addresspayable\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"participant\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"refundTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"countdown\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"payoff\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"enumAtomexEthVault.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"addresspayable\",\"name\":\"_participant\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_refundTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_countdown\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_payoff\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_active\",\"type\":\"bool\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"add\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"activate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_secret\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hashedSecret\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AtomexEthABI is the input ABI used to generate the binding from.
// Deprecated: Use AtomexEthMetaData.ABI instead.
var AtomexEthABI = AtomexEthMetaData.ABI

// AtomexEth is an auto generated Go binding around an Ethereum contract.
type AtomexEth struct {
	AtomexEthCaller     // Read-only binding to the contract
	AtomexEthTransactor // Write-only binding to the contract
	AtomexEthFilterer   // Log filterer for contract events
}

// AtomexEthCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomexEthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomexEthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomexEthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomexEthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomexEthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomexEthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomexEthSession struct {
	Contract     *AtomexEth        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomexEthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomexEthCallerSession struct {
	Contract *AtomexEthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// AtomexEthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomexEthTransactorSession struct {
	Contract     *AtomexEthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AtomexEthRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomexEthRaw struct {
	Contract *AtomexEth // Generic contract binding to access the raw methods on
}

// AtomexEthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomexEthCallerRaw struct {
	Contract *AtomexEthCaller // Generic read-only contract binding to access the raw methods on
}

// AtomexEthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomexEthTransactorRaw struct {
	Contract *AtomexEthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomexEth creates a new instance of AtomexEth, bound to a specific deployed contract.
func NewAtomexEth(address common.Address, backend bind.ContractBackend) (*AtomexEth, error) {
	contract, err := bindAtomexEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomexEth{AtomexEthCaller: AtomexEthCaller{contract: contract}, AtomexEthTransactor: AtomexEthTransactor{contract: contract}, AtomexEthFilterer: AtomexEthFilterer{contract: contract}}, nil
}

// NewAtomexEthCaller creates a new read-only instance of AtomexEth, bound to a specific deployed contract.
func NewAtomexEthCaller(address common.Address, caller bind.ContractCaller) (*AtomexEthCaller, error) {
	contract, err := bindAtomexEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomexEthCaller{contract: contract}, nil
}

// NewAtomexEthTransactor creates a new write-only instance of AtomexEth, bound to a specific deployed contract.
func NewAtomexEthTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomexEthTransactor, error) {
	contract, err := bindAtomexEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomexEthTransactor{contract: contract}, nil
}

// NewAtomexEthFilterer creates a new log filterer instance of AtomexEth, bound to a specific deployed contract.
func NewAtomexEthFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomexEthFilterer, error) {
	contract, err := bindAtomexEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomexEthFilterer{contract: contract}, nil
}

// bindAtomexEth binds a generic wrapper to an already deployed contract.
func bindAtomexEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomexEthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomexEth *AtomexEthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomexEth.Contract.AtomexEthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomexEth *AtomexEthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomexEth.Contract.AtomexEthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomexEth *AtomexEthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomexEth.Contract.AtomexEthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomexEth *AtomexEthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomexEth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomexEth *AtomexEthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomexEth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomexEth *AtomexEthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomexEth.Contract.contract.Transact(opts, method, params...)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(bytes32 hashedSecret, address initiator, address participant, uint256 refundTimestamp, uint256 countdown, uint256 value, uint256 payoff, bool active, uint8 state)
func (_AtomexEth *AtomexEthCaller) Swaps(opts *bind.CallOpts, arg0 [32]byte) (struct {
	HashedSecret    [32]byte
	Initiator       common.Address
	Participant     common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	State           uint8
}, error) {
	var out []interface{}
	err := _AtomexEth.contract.Call(opts, &out, "swaps", arg0)

	outstruct := new(struct {
		HashedSecret    [32]byte
		Initiator       common.Address
		Participant     common.Address
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
	outstruct.Initiator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Participant = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.RefundTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Countdown = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Value = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Payoff = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[7], new(bool)).(*bool)
	outstruct.State = *abi.ConvertType(out[8], new(uint8)).(*uint8)

	return *outstruct, err

}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(bytes32 hashedSecret, address initiator, address participant, uint256 refundTimestamp, uint256 countdown, uint256 value, uint256 payoff, bool active, uint8 state)
func (_AtomexEth *AtomexEthSession) Swaps(arg0 [32]byte) (struct {
	HashedSecret    [32]byte
	Initiator       common.Address
	Participant     common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	State           uint8
}, error) {
	return _AtomexEth.Contract.Swaps(&_AtomexEth.CallOpts, arg0)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(bytes32 hashedSecret, address initiator, address participant, uint256 refundTimestamp, uint256 countdown, uint256 value, uint256 payoff, bool active, uint8 state)
func (_AtomexEth *AtomexEthCallerSession) Swaps(arg0 [32]byte) (struct {
	HashedSecret    [32]byte
	Initiator       common.Address
	Participant     common.Address
	RefundTimestamp *big.Int
	Countdown       *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Active          bool
	State           uint8
}, error) {
	return _AtomexEth.Contract.Swaps(&_AtomexEth.CallOpts, arg0)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _hashedSecret) returns()
func (_AtomexEth *AtomexEthTransactor) Activate(opts *bind.TransactOpts, _hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.contract.Transact(opts, "activate", _hashedSecret)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _hashedSecret) returns()
func (_AtomexEth *AtomexEthSession) Activate(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Activate(&_AtomexEth.TransactOpts, _hashedSecret)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _hashedSecret) returns()
func (_AtomexEth *AtomexEthTransactorSession) Activate(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Activate(&_AtomexEth.TransactOpts, _hashedSecret)
}

// Add is a paid mutator transaction binding the contract method 0x446bffba.
//
// Solidity: function add(bytes32 _hashedSecret) payable returns()
func (_AtomexEth *AtomexEthTransactor) Add(opts *bind.TransactOpts, _hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.contract.Transact(opts, "add", _hashedSecret)
}

// Add is a paid mutator transaction binding the contract method 0x446bffba.
//
// Solidity: function add(bytes32 _hashedSecret) payable returns()
func (_AtomexEth *AtomexEthSession) Add(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Add(&_AtomexEth.TransactOpts, _hashedSecret)
}

// Add is a paid mutator transaction binding the contract method 0x446bffba.
//
// Solidity: function add(bytes32 _hashedSecret) payable returns()
func (_AtomexEth *AtomexEthTransactorSession) Add(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Add(&_AtomexEth.TransactOpts, _hashedSecret)
}

// Initiate is a paid mutator transaction binding the contract method 0xc5f1ee15.
//
// Solidity: function initiate(bytes32 _hashedSecret, address _participant, uint256 _refundTimestamp, uint256 _countdown, uint256 _payoff, bool _active) payable returns()
func (_AtomexEth *AtomexEthTransactor) Initiate(opts *bind.TransactOpts, _hashedSecret [32]byte, _participant common.Address, _refundTimestamp *big.Int, _countdown *big.Int, _payoff *big.Int, _active bool) (*types.Transaction, error) {
	return _AtomexEth.contract.Transact(opts, "initiate", _hashedSecret, _participant, _refundTimestamp, _countdown, _payoff, _active)
}

// Initiate is a paid mutator transaction binding the contract method 0xc5f1ee15.
//
// Solidity: function initiate(bytes32 _hashedSecret, address _participant, uint256 _refundTimestamp, uint256 _countdown, uint256 _payoff, bool _active) payable returns()
func (_AtomexEth *AtomexEthSession) Initiate(_hashedSecret [32]byte, _participant common.Address, _refundTimestamp *big.Int, _countdown *big.Int, _payoff *big.Int, _active bool) (*types.Transaction, error) {
	return _AtomexEth.Contract.Initiate(&_AtomexEth.TransactOpts, _hashedSecret, _participant, _refundTimestamp, _countdown, _payoff, _active)
}

// Initiate is a paid mutator transaction binding the contract method 0xc5f1ee15.
//
// Solidity: function initiate(bytes32 _hashedSecret, address _participant, uint256 _refundTimestamp, uint256 _countdown, uint256 _payoff, bool _active) payable returns()
func (_AtomexEth *AtomexEthTransactorSession) Initiate(_hashedSecret [32]byte, _participant common.Address, _refundTimestamp *big.Int, _countdown *big.Int, _payoff *big.Int, _active bool) (*types.Transaction, error) {
	return _AtomexEth.Contract.Initiate(&_AtomexEth.TransactOpts, _hashedSecret, _participant, _refundTimestamp, _countdown, _payoff, _active)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 _hashedSecret, bytes32 _secret) returns()
func (_AtomexEth *AtomexEthTransactor) Redeem(opts *bind.TransactOpts, _hashedSecret [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.contract.Transact(opts, "redeem", _hashedSecret, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 _hashedSecret, bytes32 _secret) returns()
func (_AtomexEth *AtomexEthSession) Redeem(_hashedSecret [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Redeem(&_AtomexEth.TransactOpts, _hashedSecret, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 _hashedSecret, bytes32 _secret) returns()
func (_AtomexEth *AtomexEthTransactorSession) Redeem(_hashedSecret [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Redeem(&_AtomexEth.TransactOpts, _hashedSecret, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _hashedSecret) returns()
func (_AtomexEth *AtomexEthTransactor) Refund(opts *bind.TransactOpts, _hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.contract.Transact(opts, "refund", _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _hashedSecret) returns()
func (_AtomexEth *AtomexEthSession) Refund(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Refund(&_AtomexEth.TransactOpts, _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _hashedSecret) returns()
func (_AtomexEth *AtomexEthTransactorSession) Refund(_hashedSecret [32]byte) (*types.Transaction, error) {
	return _AtomexEth.Contract.Refund(&_AtomexEth.TransactOpts, _hashedSecret)
}

// AtomexEthActivatedIterator is returned from FilterActivated and is used to iterate over the raw logs and unpacked data for Activated events raised by the AtomexEth contract.
type AtomexEthActivatedIterator struct {
	Event *AtomexEthActivated // Event containing the contract specifics and raw log

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
func (it *AtomexEthActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexEthActivated)
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
		it.Event = new(AtomexEthActivated)
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
func (it *AtomexEthActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexEthActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexEthActivated represents a Activated event raised by the AtomexEth contract.
type AtomexEthActivated struct {
	HashedSecret [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterActivated is a free log retrieval operation binding the contract event 0xe1abfe35306def8dbc83e3cb0bc76ffd144cee4ab7707b4e888afd4d24c2d6ca.
//
// Solidity: event Activated(bytes32 indexed _hashedSecret)
func (_AtomexEth *AtomexEthFilterer) FilterActivated(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexEthActivatedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.FilterLogs(opts, "Activated", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexEthActivatedIterator{contract: _AtomexEth.contract, event: "Activated", logs: logs, sub: sub}, nil
}

// WatchActivated is a free log subscription operation binding the contract event 0xe1abfe35306def8dbc83e3cb0bc76ffd144cee4ab7707b4e888afd4d24c2d6ca.
//
// Solidity: event Activated(bytes32 indexed _hashedSecret)
func (_AtomexEth *AtomexEthFilterer) WatchActivated(opts *bind.WatchOpts, sink chan<- *AtomexEthActivated, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.WatchLogs(opts, "Activated", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexEthActivated)
				if err := _AtomexEth.contract.UnpackLog(event, "Activated", log); err != nil {
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
func (_AtomexEth *AtomexEthFilterer) ParseActivated(log types.Log) (*AtomexEthActivated, error) {
	event := new(AtomexEthActivated)
	if err := _AtomexEth.contract.UnpackLog(event, "Activated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexEthAddedIterator is returned from FilterAdded and is used to iterate over the raw logs and unpacked data for Added events raised by the AtomexEth contract.
type AtomexEthAddedIterator struct {
	Event *AtomexEthAdded // Event containing the contract specifics and raw log

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
func (it *AtomexEthAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexEthAdded)
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
		it.Event = new(AtomexEthAdded)
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
func (it *AtomexEthAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexEthAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexEthAdded represents a Added event raised by the AtomexEth contract.
type AtomexEthAdded struct {
	HashedSecret [32]byte
	Sender       common.Address
	Value        *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAdded is a free log retrieval operation binding the contract event 0xd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423.
//
// Solidity: event Added(bytes32 indexed _hashedSecret, address _sender, uint256 _value)
func (_AtomexEth *AtomexEthFilterer) FilterAdded(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexEthAddedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.FilterLogs(opts, "Added", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexEthAddedIterator{contract: _AtomexEth.contract, event: "Added", logs: logs, sub: sub}, nil
}

// WatchAdded is a free log subscription operation binding the contract event 0xd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423.
//
// Solidity: event Added(bytes32 indexed _hashedSecret, address _sender, uint256 _value)
func (_AtomexEth *AtomexEthFilterer) WatchAdded(opts *bind.WatchOpts, sink chan<- *AtomexEthAdded, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.WatchLogs(opts, "Added", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexEthAdded)
				if err := _AtomexEth.contract.UnpackLog(event, "Added", log); err != nil {
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
func (_AtomexEth *AtomexEthFilterer) ParseAdded(log types.Log) (*AtomexEthAdded, error) {
	event := new(AtomexEthAdded)
	if err := _AtomexEth.contract.UnpackLog(event, "Added", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexEthInitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the AtomexEth contract.
type AtomexEthInitiatedIterator struct {
	Event *AtomexEthInitiated // Event containing the contract specifics and raw log

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
func (it *AtomexEthInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexEthInitiated)
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
		it.Event = new(AtomexEthInitiated)
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
func (it *AtomexEthInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexEthInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexEthInitiated represents a Initiated event raised by the AtomexEth contract.
type AtomexEthInitiated struct {
	HashedSecret    [32]byte
	Participant     common.Address
	Initiator       common.Address
	RefundTimestamp *big.Int
	Value           *big.Int
	Payoff          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x5e919055312829285818d366d1cfe50a1ba27ce2c752b655cb2faa0179e14227.
//
// Solidity: event Initiated(bytes32 indexed _hashedSecret, address indexed _participant, address _initiator, uint256 _refundTimestamp, uint256 _value, uint256 _payoff)
func (_AtomexEth *AtomexEthFilterer) FilterInitiated(opts *bind.FilterOpts, _hashedSecret [][32]byte, _participant []common.Address) (*AtomexEthInitiatedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}
	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _AtomexEth.contract.FilterLogs(opts, "Initiated", _hashedSecretRule, _participantRule)
	if err != nil {
		return nil, err
	}
	return &AtomexEthInitiatedIterator{contract: _AtomexEth.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x5e919055312829285818d366d1cfe50a1ba27ce2c752b655cb2faa0179e14227.
//
// Solidity: event Initiated(bytes32 indexed _hashedSecret, address indexed _participant, address _initiator, uint256 _refundTimestamp, uint256 _value, uint256 _payoff)
func (_AtomexEth *AtomexEthFilterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *AtomexEthInitiated, _hashedSecret [][32]byte, _participant []common.Address) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}
	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _AtomexEth.contract.WatchLogs(opts, "Initiated", _hashedSecretRule, _participantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexEthInitiated)
				if err := _AtomexEth.contract.UnpackLog(event, "Initiated", log); err != nil {
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

// ParseInitiated is a log parse operation binding the contract event 0x5e919055312829285818d366d1cfe50a1ba27ce2c752b655cb2faa0179e14227.
//
// Solidity: event Initiated(bytes32 indexed _hashedSecret, address indexed _participant, address _initiator, uint256 _refundTimestamp, uint256 _value, uint256 _payoff)
func (_AtomexEth *AtomexEthFilterer) ParseInitiated(log types.Log) (*AtomexEthInitiated, error) {
	event := new(AtomexEthInitiated)
	if err := _AtomexEth.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexEthRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the AtomexEth contract.
type AtomexEthRedeemedIterator struct {
	Event *AtomexEthRedeemed // Event containing the contract specifics and raw log

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
func (it *AtomexEthRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexEthRedeemed)
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
		it.Event = new(AtomexEthRedeemed)
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
func (it *AtomexEthRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexEthRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexEthRedeemed represents a Redeemed event raised by the AtomexEth contract.
type AtomexEthRedeemed struct {
	HashedSecret [32]byte
	Secret       [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab.
//
// Solidity: event Redeemed(bytes32 indexed _hashedSecret, bytes32 _secret)
func (_AtomexEth *AtomexEthFilterer) FilterRedeemed(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexEthRedeemedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.FilterLogs(opts, "Redeemed", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexEthRedeemedIterator{contract: _AtomexEth.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab.
//
// Solidity: event Redeemed(bytes32 indexed _hashedSecret, bytes32 _secret)
func (_AtomexEth *AtomexEthFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *AtomexEthRedeemed, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.WatchLogs(opts, "Redeemed", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexEthRedeemed)
				if err := _AtomexEth.contract.UnpackLog(event, "Redeemed", log); err != nil {
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
func (_AtomexEth *AtomexEthFilterer) ParseRedeemed(log types.Log) (*AtomexEthRedeemed, error) {
	event := new(AtomexEthRedeemed)
	if err := _AtomexEth.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomexEthRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the AtomexEth contract.
type AtomexEthRefundedIterator struct {
	Event *AtomexEthRefunded // Event containing the contract specifics and raw log

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
func (it *AtomexEthRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomexEthRefunded)
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
		it.Event = new(AtomexEthRefunded)
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
func (it *AtomexEthRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomexEthRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomexEthRefunded represents a Refunded event raised by the AtomexEth contract.
type AtomexEthRefunded struct {
	HashedSecret [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed _hashedSecret)
func (_AtomexEth *AtomexEthFilterer) FilterRefunded(opts *bind.FilterOpts, _hashedSecret [][32]byte) (*AtomexEthRefundedIterator, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.FilterLogs(opts, "Refunded", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return &AtomexEthRefundedIterator{contract: _AtomexEth.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed _hashedSecret)
func (_AtomexEth *AtomexEthFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *AtomexEthRefunded, _hashedSecret [][32]byte) (event.Subscription, error) {

	var _hashedSecretRule []interface{}
	for _, _hashedSecretItem := range _hashedSecret {
		_hashedSecretRule = append(_hashedSecretRule, _hashedSecretItem)
	}

	logs, sub, err := _AtomexEth.contract.WatchLogs(opts, "Refunded", _hashedSecretRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomexEthRefunded)
				if err := _AtomexEth.contract.UnpackLog(event, "Refunded", log); err != nil {
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
func (_AtomexEth *AtomexEthFilterer) ParseRefunded(log types.Log) (*AtomexEthRefunded, error) {
	event := new(AtomexEthRefunded)
	if err := _AtomexEth.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
