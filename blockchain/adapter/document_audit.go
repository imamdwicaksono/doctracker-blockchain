// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package adapter

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

// DocumentAuditMetaData contains all meta data concerning the DocumentAudit contract.
var DocumentAuditMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"docId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"action\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"actor\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"DocumentEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"docId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"action\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"actor\",\"type\":\"string\"}],\"name\":\"record\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DocumentAuditABI is the input ABI used to generate the binding from.
// Deprecated: Use DocumentAuditMetaData.ABI instead.
var DocumentAuditABI = DocumentAuditMetaData.ABI

// DocumentAudit is an auto generated Go binding around an Ethereum contract.
type DocumentAudit struct {
	DocumentAuditCaller     // Read-only binding to the contract
	DocumentAuditTransactor // Write-only binding to the contract
	DocumentAuditFilterer   // Log filterer for contract events
}

// DocumentAuditCaller is an auto generated read-only Go binding around an Ethereum contract.
type DocumentAuditCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DocumentAuditTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DocumentAuditTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DocumentAuditFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DocumentAuditFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DocumentAuditSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DocumentAuditSession struct {
	Contract     *DocumentAudit    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DocumentAuditCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DocumentAuditCallerSession struct {
	Contract *DocumentAuditCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// DocumentAuditTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DocumentAuditTransactorSession struct {
	Contract     *DocumentAuditTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// DocumentAuditRaw is an auto generated low-level Go binding around an Ethereum contract.
type DocumentAuditRaw struct {
	Contract *DocumentAudit // Generic contract binding to access the raw methods on
}

// DocumentAuditCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DocumentAuditCallerRaw struct {
	Contract *DocumentAuditCaller // Generic read-only contract binding to access the raw methods on
}

// DocumentAuditTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DocumentAuditTransactorRaw struct {
	Contract *DocumentAuditTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDocumentAudit creates a new instance of DocumentAudit, bound to a specific deployed contract.
func NewDocumentAudit(address common.Address, backend bind.ContractBackend) (*DocumentAudit, error) {
	contract, err := bindDocumentAudit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DocumentAudit{DocumentAuditCaller: DocumentAuditCaller{contract: contract}, DocumentAuditTransactor: DocumentAuditTransactor{contract: contract}, DocumentAuditFilterer: DocumentAuditFilterer{contract: contract}}, nil
}

// NewDocumentAuditCaller creates a new read-only instance of DocumentAudit, bound to a specific deployed contract.
func NewDocumentAuditCaller(address common.Address, caller bind.ContractCaller) (*DocumentAuditCaller, error) {
	contract, err := bindDocumentAudit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DocumentAuditCaller{contract: contract}, nil
}

// NewDocumentAuditTransactor creates a new write-only instance of DocumentAudit, bound to a specific deployed contract.
func NewDocumentAuditTransactor(address common.Address, transactor bind.ContractTransactor) (*DocumentAuditTransactor, error) {
	contract, err := bindDocumentAudit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DocumentAuditTransactor{contract: contract}, nil
}

// NewDocumentAuditFilterer creates a new log filterer instance of DocumentAudit, bound to a specific deployed contract.
func NewDocumentAuditFilterer(address common.Address, filterer bind.ContractFilterer) (*DocumentAuditFilterer, error) {
	contract, err := bindDocumentAudit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DocumentAuditFilterer{contract: contract}, nil
}

// bindDocumentAudit binds a generic wrapper to an already deployed contract.
func bindDocumentAudit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DocumentAuditMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DocumentAudit *DocumentAuditRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DocumentAudit.Contract.DocumentAuditCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DocumentAudit *DocumentAuditRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DocumentAudit.Contract.DocumentAuditTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DocumentAudit *DocumentAuditRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DocumentAudit.Contract.DocumentAuditTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DocumentAudit *DocumentAuditCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DocumentAudit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DocumentAudit *DocumentAuditTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DocumentAudit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DocumentAudit *DocumentAuditTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DocumentAudit.Contract.contract.Transact(opts, method, params...)
}

// Record is a paid mutator transaction binding the contract method 0x2871bb40.
//
// Solidity: function record(string docId, bytes32 hash, string action, string actor) returns()
func (_DocumentAudit *DocumentAuditTransactor) Record(opts *bind.TransactOpts, docId string, hash [32]byte, action string, actor string) (*types.Transaction, error) {
	return _DocumentAudit.contract.Transact(opts, "record", docId, hash, action, actor)
}

// Record is a paid mutator transaction binding the contract method 0x2871bb40.
//
// Solidity: function record(string docId, bytes32 hash, string action, string actor) returns()
func (_DocumentAudit *DocumentAuditSession) Record(docId string, hash [32]byte, action string, actor string) (*types.Transaction, error) {
	return _DocumentAudit.Contract.Record(&_DocumentAudit.TransactOpts, docId, hash, action, actor)
}

// Record is a paid mutator transaction binding the contract method 0x2871bb40.
//
// Solidity: function record(string docId, bytes32 hash, string action, string actor) returns()
func (_DocumentAudit *DocumentAuditTransactorSession) Record(docId string, hash [32]byte, action string, actor string) (*types.Transaction, error) {
	return _DocumentAudit.Contract.Record(&_DocumentAudit.TransactOpts, docId, hash, action, actor)
}

// DocumentAuditDocumentEventIterator is returned from FilterDocumentEvent and is used to iterate over the raw logs and unpacked data for DocumentEvent events raised by the DocumentAudit contract.
type DocumentAuditDocumentEventIterator struct {
	Event *DocumentAuditDocumentEvent // Event containing the contract specifics and raw log

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
func (it *DocumentAuditDocumentEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DocumentAuditDocumentEvent)
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
		it.Event = new(DocumentAuditDocumentEvent)
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
func (it *DocumentAuditDocumentEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DocumentAuditDocumentEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DocumentAuditDocumentEvent represents a DocumentEvent event raised by the DocumentAudit contract.
type DocumentAuditDocumentEvent struct {
	DocId     common.Hash
	Hash      [32]byte
	Action    string
	Actor     string
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDocumentEvent is a free log retrieval operation binding the contract event 0x5a11a87e39b86202f02d630592e143faa6074c7ddd09477321d0588a37fe1b94.
//
// Solidity: event DocumentEvent(string indexed docId, bytes32 hash, string action, string actor, uint256 timestamp)
func (_DocumentAudit *DocumentAuditFilterer) FilterDocumentEvent(opts *bind.FilterOpts, docId []string) (*DocumentAuditDocumentEventIterator, error) {

	var docIdRule []interface{}
	for _, docIdItem := range docId {
		docIdRule = append(docIdRule, docIdItem)
	}

	logs, sub, err := _DocumentAudit.contract.FilterLogs(opts, "DocumentEvent", docIdRule)
	if err != nil {
		return nil, err
	}
	return &DocumentAuditDocumentEventIterator{contract: _DocumentAudit.contract, event: "DocumentEvent", logs: logs, sub: sub}, nil
}

// WatchDocumentEvent is a free log subscription operation binding the contract event 0x5a11a87e39b86202f02d630592e143faa6074c7ddd09477321d0588a37fe1b94.
//
// Solidity: event DocumentEvent(string indexed docId, bytes32 hash, string action, string actor, uint256 timestamp)
func (_DocumentAudit *DocumentAuditFilterer) WatchDocumentEvent(opts *bind.WatchOpts, sink chan<- *DocumentAuditDocumentEvent, docId []string) (event.Subscription, error) {

	var docIdRule []interface{}
	for _, docIdItem := range docId {
		docIdRule = append(docIdRule, docIdItem)
	}

	logs, sub, err := _DocumentAudit.contract.WatchLogs(opts, "DocumentEvent", docIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DocumentAuditDocumentEvent)
				if err := _DocumentAudit.contract.UnpackLog(event, "DocumentEvent", log); err != nil {
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

// ParseDocumentEvent is a log parse operation binding the contract event 0x5a11a87e39b86202f02d630592e143faa6074c7ddd09477321d0588a37fe1b94.
//
// Solidity: event DocumentEvent(string indexed docId, bytes32 hash, string action, string actor, uint256 timestamp)
func (_DocumentAudit *DocumentAuditFilterer) ParseDocumentEvent(log types.Log) (*DocumentAuditDocumentEvent, error) {
	event := new(DocumentAuditDocumentEvent)
	if err := _DocumentAudit.contract.UnpackLog(event, "DocumentEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
