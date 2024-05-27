// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package StateOverride

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
	_ = abi.ConvertType
)

// StateOverrideMetaData contains all meta data concerning the StateOverride contract.
var StateOverrideMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"addrBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getText\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6001600081905560c060405260046080819052631d195e1d60e21b60a0908152610029929161003c565b5034801561003657600080fd5b50610110565b828054610048906100d5565b90600052602060002090601f01602090048101928261006a57600085556100b0565b82601f1061008357805160ff19168380011785556100b0565b828001600101855582156100b0579182015b828111156100b0578251825591602001919060010190610095565b506100bc9291506100c0565b5090565b5b808211156100bc57600081556001016100c1565b600181811c908216806100e957607f821691505b6020821081141561010a57634e487b7160e01b600052602260045260246000fd5b50919050565b6102198061011f6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80638262fc7d14610046578063e00fe2eb14610074578063f2c9ecd814610089575b600080fd5b610061610054366004610123565b6001600160a01b03163190565b6040519081526020015b60405180910390f35b61007c610091565b60405161006b9190610153565b600054610061565b6060600180546100a0906101a8565b80601f01602080910402602001604051908101604052809291908181526020018280546100cc906101a8565b80156101195780601f106100ee57610100808354040283529160200191610119565b820191906000526020600020905b8154815290600101906020018083116100fc57829003601f168201915b5050505050905090565b60006020828403121561013557600080fd5b81356001600160a01b038116811461014c57600080fd5b9392505050565b600060208083528351808285015260005b8181101561018057858101830151858201604001528201610164565b81811115610192576000604083870101525b50601f01601f1916929092016040019392505050565b600181811c908216806101bc57607f821691505b602082108114156101dd57634e487b7160e01b600052602260045260246000fd5b5091905056fea26469706673582212202b71250e01a702a9026cbfa6ea39d8f3324e85249751ef5ce0955bc30e413aaa64736f6c634300080c0033",
}

// StateOverrideABI is the input ABI used to generate the binding from.
// Deprecated: Use StateOverrideMetaData.ABI instead.
var StateOverrideABI = StateOverrideMetaData.ABI

// StateOverrideBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StateOverrideMetaData.Bin instead.
var StateOverrideBin = StateOverrideMetaData.Bin

// DeployStateOverride deploys a new Ethereum contract, binding an instance of StateOverride to it.
func DeployStateOverride(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StateOverride, error) {
	parsed, err := StateOverrideMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StateOverrideBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StateOverride{StateOverrideCaller: StateOverrideCaller{contract: contract}, StateOverrideTransactor: StateOverrideTransactor{contract: contract}, StateOverrideFilterer: StateOverrideFilterer{contract: contract}}, nil
}

// StateOverride is an auto generated Go binding around an Ethereum contract.
type StateOverride struct {
	StateOverrideCaller     // Read-only binding to the contract
	StateOverrideTransactor // Write-only binding to the contract
	StateOverrideFilterer   // Log filterer for contract events
}

// StateOverrideCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateOverrideCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateOverrideTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateOverrideTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateOverrideFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateOverrideFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateOverrideSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateOverrideSession struct {
	Contract     *StateOverride    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateOverrideCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateOverrideCallerSession struct {
	Contract *StateOverrideCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StateOverrideTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateOverrideTransactorSession struct {
	Contract     *StateOverrideTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StateOverrideRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateOverrideRaw struct {
	Contract *StateOverride // Generic contract binding to access the raw methods on
}

// StateOverrideCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateOverrideCallerRaw struct {
	Contract *StateOverrideCaller // Generic read-only contract binding to access the raw methods on
}

// StateOverrideTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateOverrideTransactorRaw struct {
	Contract *StateOverrideTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStateOverride creates a new instance of StateOverride, bound to a specific deployed contract.
func NewStateOverride(address common.Address, backend bind.ContractBackend) (*StateOverride, error) {
	contract, err := bindStateOverride(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StateOverride{StateOverrideCaller: StateOverrideCaller{contract: contract}, StateOverrideTransactor: StateOverrideTransactor{contract: contract}, StateOverrideFilterer: StateOverrideFilterer{contract: contract}}, nil
}

// NewStateOverrideCaller creates a new read-only instance of StateOverride, bound to a specific deployed contract.
func NewStateOverrideCaller(address common.Address, caller bind.ContractCaller) (*StateOverrideCaller, error) {
	contract, err := bindStateOverride(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateOverrideCaller{contract: contract}, nil
}

// NewStateOverrideTransactor creates a new write-only instance of StateOverride, bound to a specific deployed contract.
func NewStateOverrideTransactor(address common.Address, transactor bind.ContractTransactor) (*StateOverrideTransactor, error) {
	contract, err := bindStateOverride(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateOverrideTransactor{contract: contract}, nil
}

// NewStateOverrideFilterer creates a new log filterer instance of StateOverride, bound to a specific deployed contract.
func NewStateOverrideFilterer(address common.Address, filterer bind.ContractFilterer) (*StateOverrideFilterer, error) {
	contract, err := bindStateOverride(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateOverrideFilterer{contract: contract}, nil
}

// bindStateOverride binds a generic wrapper to an already deployed contract.
func bindStateOverride(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StateOverrideMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateOverride *StateOverrideRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateOverride.Contract.StateOverrideCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateOverride *StateOverrideRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateOverride.Contract.StateOverrideTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateOverride *StateOverrideRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateOverride.Contract.StateOverrideTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateOverride *StateOverrideCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateOverride.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateOverride *StateOverrideTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateOverride.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateOverride *StateOverrideTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateOverride.Contract.contract.Transact(opts, method, params...)
}

// AddrBalance is a free data retrieval call binding the contract method 0x8262fc7d.
//
// Solidity: function addrBalance(address a) view returns(uint256)
func (_StateOverride *StateOverrideCaller) AddrBalance(opts *bind.CallOpts, a common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StateOverride.contract.Call(opts, &out, "addrBalance", a)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AddrBalance is a free data retrieval call binding the contract method 0x8262fc7d.
//
// Solidity: function addrBalance(address a) view returns(uint256)
func (_StateOverride *StateOverrideSession) AddrBalance(a common.Address) (*big.Int, error) {
	return _StateOverride.Contract.AddrBalance(&_StateOverride.CallOpts, a)
}

// AddrBalance is a free data retrieval call binding the contract method 0x8262fc7d.
//
// Solidity: function addrBalance(address a) view returns(uint256)
func (_StateOverride *StateOverrideCallerSession) AddrBalance(a common.Address) (*big.Int, error) {
	return _StateOverride.Contract.AddrBalance(&_StateOverride.CallOpts, a)
}

// GetNumber is a free data retrieval call binding the contract method 0xf2c9ecd8.
//
// Solidity: function getNumber() view returns(uint256)
func (_StateOverride *StateOverrideCaller) GetNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateOverride.contract.Call(opts, &out, "getNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumber is a free data retrieval call binding the contract method 0xf2c9ecd8.
//
// Solidity: function getNumber() view returns(uint256)
func (_StateOverride *StateOverrideSession) GetNumber() (*big.Int, error) {
	return _StateOverride.Contract.GetNumber(&_StateOverride.CallOpts)
}

// GetNumber is a free data retrieval call binding the contract method 0xf2c9ecd8.
//
// Solidity: function getNumber() view returns(uint256)
func (_StateOverride *StateOverrideCallerSession) GetNumber() (*big.Int, error) {
	return _StateOverride.Contract.GetNumber(&_StateOverride.CallOpts)
}

// GetText is a free data retrieval call binding the contract method 0xe00fe2eb.
//
// Solidity: function getText() view returns(string)
func (_StateOverride *StateOverrideCaller) GetText(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StateOverride.contract.Call(opts, &out, "getText")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetText is a free data retrieval call binding the contract method 0xe00fe2eb.
//
// Solidity: function getText() view returns(string)
func (_StateOverride *StateOverrideSession) GetText() (string, error) {
	return _StateOverride.Contract.GetText(&_StateOverride.CallOpts)
}

// GetText is a free data retrieval call binding the contract method 0xe00fe2eb.
//
// Solidity: function getText() view returns(string)
func (_StateOverride *StateOverrideCallerSession) GetText() (string, error) {
	return _StateOverride.Contract.GetText(&_StateOverride.CallOpts)
}
