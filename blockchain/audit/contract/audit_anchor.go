// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"BlockAnchored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"anchorBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"blockHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506103a78061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80631d97ce311461004357806334cdf78d14610073578063d73f0685146100a3575b5f5ffd5b61005d600480360381019061005891906101ff565b6100bf565b60405161006a9190610257565b60405180910390f35b61008d60048036038101906100889190610270565b6100db565b60405161009a91906102aa565b60405180910390f35b6100bd60048036038101906100b891906101ff565b6100ef565b005b5f815f5f8581526020019081526020015f205414905092915050565b5f602052805f5260405f205f915090505481565b5f5f1b5f5f8481526020019081526020015f205414610143576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013a9061031d565b60405180910390fd5b805f5f8481526020019081526020015f20819055507ff730790d1ef3dad1766c4b7af29240af4e726dcff5b71db598c7b98efa9f8146828260405161018992919061034a565b60405180910390a15050565b5f5ffd5b5f819050919050565b6101ab81610199565b81146101b5575f5ffd5b50565b5f813590506101c6816101a2565b92915050565b5f819050919050565b6101de816101cc565b81146101e8575f5ffd5b50565b5f813590506101f9816101d5565b92915050565b5f5f6040838503121561021557610214610195565b5b5f610222858286016101b8565b9250506020610233858286016101eb565b9150509250929050565b5f8115159050919050565b6102518161023d565b82525050565b5f60208201905061026a5f830184610248565b92915050565b5f6020828403121561028557610284610195565b5b5f610292848285016101b8565b91505092915050565b6102a4816101cc565b82525050565b5f6020820190506102bd5f83018461029b565b92915050565b5f82825260208201905092915050565b7f414c52454144595f414e43484f524544000000000000000000000000000000005f82015250565b5f6103076010836102c3565b9150610312826102d3565b602082019050919050565b5f6020820190508181035f830152610334816102fb565b9050919050565b61034481610199565b82525050565b5f60408201905061035d5f83018561033b565b61036a602083018461029b565b939250505056fea2646970667358221220a96197ae324a7250a0428175f7174074b3d3806b76e964ca056b8925979affb464736f6c63430008210033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// BlockHashes is a free data retrieval call binding the contract method 0x34cdf78d.
//
// Solidity: function blockHashes(uint256 ) view returns(bytes32)
func (_Contract *ContractCaller) BlockHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "blockHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockHashes is a free data retrieval call binding the contract method 0x34cdf78d.
//
// Solidity: function blockHashes(uint256 ) view returns(bytes32)
func (_Contract *ContractSession) BlockHashes(arg0 *big.Int) ([32]byte, error) {
	return _Contract.Contract.BlockHashes(&_Contract.CallOpts, arg0)
}

// BlockHashes is a free data retrieval call binding the contract method 0x34cdf78d.
//
// Solidity: function blockHashes(uint256 ) view returns(bytes32)
func (_Contract *ContractCallerSession) BlockHashes(arg0 *big.Int) ([32]byte, error) {
	return _Contract.Contract.BlockHashes(&_Contract.CallOpts, arg0)
}

// Verify is a free data retrieval call binding the contract method 0x1d97ce31.
//
// Solidity: function verify(uint256 height, bytes32 hash) view returns(bool)
func (_Contract *ContractCaller) Verify(opts *bind.CallOpts, height *big.Int, hash [32]byte) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "verify", height, hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0x1d97ce31.
//
// Solidity: function verify(uint256 height, bytes32 hash) view returns(bool)
func (_Contract *ContractSession) Verify(height *big.Int, hash [32]byte) (bool, error) {
	return _Contract.Contract.Verify(&_Contract.CallOpts, height, hash)
}

// Verify is a free data retrieval call binding the contract method 0x1d97ce31.
//
// Solidity: function verify(uint256 height, bytes32 hash) view returns(bool)
func (_Contract *ContractCallerSession) Verify(height *big.Int, hash [32]byte) (bool, error) {
	return _Contract.Contract.Verify(&_Contract.CallOpts, height, hash)
}

// AnchorBlock is a paid mutator transaction binding the contract method 0xd73f0685.
//
// Solidity: function anchorBlock(uint256 height, bytes32 hash) returns()
func (_Contract *ContractTransactor) AnchorBlock(opts *bind.TransactOpts, height *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "anchorBlock", height, hash)
}

// AnchorBlock is a paid mutator transaction binding the contract method 0xd73f0685.
//
// Solidity: function anchorBlock(uint256 height, bytes32 hash) returns()
func (_Contract *ContractSession) AnchorBlock(height *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _Contract.Contract.AnchorBlock(&_Contract.TransactOpts, height, hash)
}

// AnchorBlock is a paid mutator transaction binding the contract method 0xd73f0685.
//
// Solidity: function anchorBlock(uint256 height, bytes32 hash) returns()
func (_Contract *ContractTransactorSession) AnchorBlock(height *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _Contract.Contract.AnchorBlock(&_Contract.TransactOpts, height, hash)
}

// ContractBlockAnchoredIterator is returned from FilterBlockAnchored and is used to iterate over the raw logs and unpacked data for BlockAnchored events raised by the Contract contract.
type ContractBlockAnchoredIterator struct {
	Event *ContractBlockAnchored // Event containing the contract specifics and raw log

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
func (it *ContractBlockAnchoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBlockAnchored)
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
		it.Event = new(ContractBlockAnchored)
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
func (it *ContractBlockAnchoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBlockAnchoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBlockAnchored represents a BlockAnchored event raised by the Contract contract.
type ContractBlockAnchored struct {
	Height *big.Int
	Hash   [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBlockAnchored is a free log retrieval operation binding the contract event 0xf730790d1ef3dad1766c4b7af29240af4e726dcff5b71db598c7b98efa9f8146.
//
// Solidity: event BlockAnchored(uint256 height, bytes32 hash)
func (_Contract *ContractFilterer) FilterBlockAnchored(opts *bind.FilterOpts) (*ContractBlockAnchoredIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BlockAnchored")
	if err != nil {
		return nil, err
	}
	return &ContractBlockAnchoredIterator{contract: _Contract.contract, event: "BlockAnchored", logs: logs, sub: sub}, nil
}

// WatchBlockAnchored is a free log subscription operation binding the contract event 0xf730790d1ef3dad1766c4b7af29240af4e726dcff5b71db598c7b98efa9f8146.
//
// Solidity: event BlockAnchored(uint256 height, bytes32 hash)
func (_Contract *ContractFilterer) WatchBlockAnchored(opts *bind.WatchOpts, sink chan<- *ContractBlockAnchored) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BlockAnchored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBlockAnchored)
				if err := _Contract.contract.UnpackLog(event, "BlockAnchored", log); err != nil {
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

// ParseBlockAnchored is a log parse operation binding the contract event 0xf730790d1ef3dad1766c4b7af29240af4e726dcff5b71db598c7b98efa9f8146.
//
// Solidity: event BlockAnchored(uint256 height, bytes32 hash)
func (_Contract *ContractFilterer) ParseBlockAnchored(log types.Log) (*ContractBlockAnchored, error) {
	event := new(ContractBlockAnchored)
	if err := _Contract.contract.UnpackLog(event, "BlockAnchored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
