// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package metamodel

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

// ModelContext is an auto generated low-level Go binding around an user-defined struct.
type ModelContext struct {
	Sequence    *big.Int
	State       []*big.Int
	Places      []ModelPlace
	Transitions []ModelTransition
}

// ModelPlace is an auto generated low-level Go binding around an user-defined struct.
type ModelPlace struct {
	Label    string
	Offset   uint8
	Position ModelPosition
	Initial  *big.Int
	Capacity *big.Int
}

// ModelPosition is an auto generated low-level Go binding around an user-defined struct.
type ModelPosition struct {
	X uint8
	Y uint8
}

// ModelTransition is an auto generated low-level Go binding around an user-defined struct.
type ModelTransition struct {
	Label    string
	Offset   uint8
	Position ModelPosition
	Role     uint8
	Delta    []*big.Int
	Guard    []*big.Int
}

// MetamodelMetaData contains all meta data concerning the Metamodel contract.
var MetamodelMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"role\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"actionId\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"}],\"name\":\"SignaledEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"context\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"int256[]\",\"name\":\"state\",\"type\":\"int256[]\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"offset\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"x\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"y\",\"type\":\"uint8\"}],\"internalType\":\"structModel.Position\",\"name\":\"position\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"initial\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"}],\"internalType\":\"structModel.Place[]\",\"name\":\"places\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"offset\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"x\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"y\",\"type\":\"uint8\"}],\"internalType\":\"structModel.Position\",\"name\":\"position\",\"type\":\"tuple\"},{\"internalType\":\"uint8\",\"name\":\"role\",\"type\":\"uint8\"},{\"internalType\":\"int256[]\",\"name\":\"delta\",\"type\":\"int256[]\"},{\"internalType\":\"int256[]\",\"name\":\"guard\",\"type\":\"int256[]\"}],\"internalType\":\"structModel.Transition[]\",\"name\":\"transitions\",\"type\":\"tuple[]\"}],\"internalType\":\"structModel.Context\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sequence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"action\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"}],\"name\":\"signal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8[]\",\"name\":\"actions\",\"type\":\"uint8[]\"},{\"internalType\":\"uint256[]\",\"name\":\"scalars\",\"type\":\"uint256[]\"}],\"name\":\"signalMany\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"state\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// MetamodelABI is the input ABI used to generate the binding from.
// Deprecated: Use MetamodelMetaData.ABI instead.
var MetamodelABI = MetamodelMetaData.ABI

// Metamodel is an auto generated Go binding around an Ethereum contract.
type Metamodel struct {
	MetamodelCaller     // Read-only binding to the contract
	MetamodelTransactor // Write-only binding to the contract
	MetamodelFilterer   // Log filterer for contract events
}

// MetamodelCaller is an auto generated read-only Go binding around an Ethereum contract.
type MetamodelCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetamodelTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MetamodelTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetamodelFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MetamodelFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetamodelSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MetamodelSession struct {
	Contract     *Metamodel        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MetamodelCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MetamodelCallerSession struct {
	Contract *MetamodelCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MetamodelTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MetamodelTransactorSession struct {
	Contract     *MetamodelTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MetamodelRaw is an auto generated low-level Go binding around an Ethereum contract.
type MetamodelRaw struct {
	Contract *Metamodel // Generic contract binding to access the raw methods on
}

// MetamodelCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MetamodelCallerRaw struct {
	Contract *MetamodelCaller // Generic read-only contract binding to access the raw methods on
}

// MetamodelTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MetamodelTransactorRaw struct {
	Contract *MetamodelTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMetamodel creates a new instance of Metamodel, bound to a specific deployed contract.
func NewMetamodel(address common.Address, backend bind.ContractBackend) (*Metamodel, error) {
	contract, err := bindMetamodel(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Metamodel{MetamodelCaller: MetamodelCaller{contract: contract}, MetamodelTransactor: MetamodelTransactor{contract: contract}, MetamodelFilterer: MetamodelFilterer{contract: contract}}, nil
}

// NewMetamodelCaller creates a new read-only instance of Metamodel, bound to a specific deployed contract.
func NewMetamodelCaller(address common.Address, caller bind.ContractCaller) (*MetamodelCaller, error) {
	contract, err := bindMetamodel(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MetamodelCaller{contract: contract}, nil
}

// NewMetamodelTransactor creates a new write-only instance of Metamodel, bound to a specific deployed contract.
func NewMetamodelTransactor(address common.Address, transactor bind.ContractTransactor) (*MetamodelTransactor, error) {
	contract, err := bindMetamodel(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MetamodelTransactor{contract: contract}, nil
}

// NewMetamodelFilterer creates a new log filterer instance of Metamodel, bound to a specific deployed contract.
func NewMetamodelFilterer(address common.Address, filterer bind.ContractFilterer) (*MetamodelFilterer, error) {
	contract, err := bindMetamodel(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MetamodelFilterer{contract: contract}, nil
}

// bindMetamodel binds a generic wrapper to an already deployed contract.
func bindMetamodel(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MetamodelMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Metamodel *MetamodelRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Metamodel.Contract.MetamodelCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Metamodel *MetamodelRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Metamodel.Contract.MetamodelTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Metamodel *MetamodelRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Metamodel.Contract.MetamodelTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Metamodel *MetamodelCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Metamodel.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Metamodel *MetamodelTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Metamodel.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Metamodel *MetamodelTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Metamodel.Contract.contract.Transact(opts, method, params...)
}

// Context is a free data retrieval call binding the contract method 0xd0496d6a.
//
// Solidity: function context() view returns((uint256,int256[],(string,uint8,(uint8,uint8),uint256,uint256)[],(string,uint8,(uint8,uint8),uint8,int256[],int256[])[]))
func (_Metamodel *MetamodelCaller) Context(opts *bind.CallOpts) (ModelContext, error) {
	var out []interface{}
	err := _Metamodel.contract.Call(opts, &out, "context")

	if err != nil {
		return *new(ModelContext), err
	}

	out0 := *abi.ConvertType(out[0], new(ModelContext)).(*ModelContext)

	return out0, err

}

// Context is a free data retrieval call binding the contract method 0xd0496d6a.
//
// Solidity: function context() view returns((uint256,int256[],(string,uint8,(uint8,uint8),uint256,uint256)[],(string,uint8,(uint8,uint8),uint8,int256[],int256[])[]))
func (_Metamodel *MetamodelSession) Context() (ModelContext, error) {
	return _Metamodel.Contract.Context(&_Metamodel.CallOpts)
}

// Context is a free data retrieval call binding the contract method 0xd0496d6a.
//
// Solidity: function context() view returns((uint256,int256[],(string,uint8,(uint8,uint8),uint256,uint256)[],(string,uint8,(uint8,uint8),uint8,int256[],int256[])[]))
func (_Metamodel *MetamodelCallerSession) Context() (ModelContext, error) {
	return _Metamodel.Contract.Context(&_Metamodel.CallOpts)
}

// Sequence is a free data retrieval call binding the contract method 0x529d15cc.
//
// Solidity: function sequence() view returns(uint256)
func (_Metamodel *MetamodelCaller) Sequence(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Metamodel.contract.Call(opts, &out, "sequence")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Sequence is a free data retrieval call binding the contract method 0x529d15cc.
//
// Solidity: function sequence() view returns(uint256)
func (_Metamodel *MetamodelSession) Sequence() (*big.Int, error) {
	return _Metamodel.Contract.Sequence(&_Metamodel.CallOpts)
}

// Sequence is a free data retrieval call binding the contract method 0x529d15cc.
//
// Solidity: function sequence() view returns(uint256)
func (_Metamodel *MetamodelCallerSession) Sequence() (*big.Int, error) {
	return _Metamodel.Contract.Sequence(&_Metamodel.CallOpts)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 ) view returns(int256)
func (_Metamodel *MetamodelCaller) State(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Metamodel.contract.Call(opts, &out, "state", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 ) view returns(int256)
func (_Metamodel *MetamodelSession) State(arg0 *big.Int) (*big.Int, error) {
	return _Metamodel.Contract.State(&_Metamodel.CallOpts, arg0)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 ) view returns(int256)
func (_Metamodel *MetamodelCallerSession) State(arg0 *big.Int) (*big.Int, error) {
	return _Metamodel.Contract.State(&_Metamodel.CallOpts, arg0)
}

// Signal is a paid mutator transaction binding the contract method 0xddc3b187.
//
// Solidity: function signal(uint8 action, uint256 scalar) returns()
func (_Metamodel *MetamodelTransactor) Signal(opts *bind.TransactOpts, action uint8, scalar *big.Int) (*types.Transaction, error) {
	return _Metamodel.contract.Transact(opts, "signal", action, scalar)
}

// Signal is a paid mutator transaction binding the contract method 0xddc3b187.
//
// Solidity: function signal(uint8 action, uint256 scalar) returns()
func (_Metamodel *MetamodelSession) Signal(action uint8, scalar *big.Int) (*types.Transaction, error) {
	return _Metamodel.Contract.Signal(&_Metamodel.TransactOpts, action, scalar)
}

// Signal is a paid mutator transaction binding the contract method 0xddc3b187.
//
// Solidity: function signal(uint8 action, uint256 scalar) returns()
func (_Metamodel *MetamodelTransactorSession) Signal(action uint8, scalar *big.Int) (*types.Transaction, error) {
	return _Metamodel.Contract.Signal(&_Metamodel.TransactOpts, action, scalar)
}

// SignalMany is a paid mutator transaction binding the contract method 0xfff01fe2.
//
// Solidity: function signalMany(uint8[] actions, uint256[] scalars) returns()
func (_Metamodel *MetamodelTransactor) SignalMany(opts *bind.TransactOpts, actions []uint8, scalars []*big.Int) (*types.Transaction, error) {
	return _Metamodel.contract.Transact(opts, "signalMany", actions, scalars)
}

// SignalMany is a paid mutator transaction binding the contract method 0xfff01fe2.
//
// Solidity: function signalMany(uint8[] actions, uint256[] scalars) returns()
func (_Metamodel *MetamodelSession) SignalMany(actions []uint8, scalars []*big.Int) (*types.Transaction, error) {
	return _Metamodel.Contract.SignalMany(&_Metamodel.TransactOpts, actions, scalars)
}

// SignalMany is a paid mutator transaction binding the contract method 0xfff01fe2.
//
// Solidity: function signalMany(uint8[] actions, uint256[] scalars) returns()
func (_Metamodel *MetamodelTransactorSession) SignalMany(actions []uint8, scalars []*big.Int) (*types.Transaction, error) {
	return _Metamodel.Contract.SignalMany(&_Metamodel.TransactOpts, actions, scalars)
}

// MetamodelSignaledEventIterator is returned from FilterSignaledEvent and is used to iterate over the raw logs and unpacked data for SignaledEvent events raised by the Metamodel contract.
type MetamodelSignaledEventIterator struct {
	Event *MetamodelSignaledEvent // Event containing the contract specifics and raw log

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
func (it *MetamodelSignaledEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetamodelSignaledEvent)
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
		it.Event = new(MetamodelSignaledEvent)
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
func (it *MetamodelSignaledEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetamodelSignaledEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetamodelSignaledEvent represents a SignaledEvent event raised by the Metamodel contract.
type MetamodelSignaledEvent struct {
	Role     uint8
	ActionId uint8
	Scalar   *big.Int
	Sequence *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSignaledEvent is a free log retrieval operation binding the contract event 0x50e4a52b0772bed9f06a7d3f7dfa66d73640067a5cc77a73c24556ccc900fa08.
//
// Solidity: event SignaledEvent(uint8 indexed role, uint8 indexed actionId, uint256 indexed scalar, uint256 sequence)
func (_Metamodel *MetamodelFilterer) FilterSignaledEvent(opts *bind.FilterOpts, role []uint8, actionId []uint8, scalar []*big.Int) (*MetamodelSignaledEventIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var actionIdRule []interface{}
	for _, actionIdItem := range actionId {
		actionIdRule = append(actionIdRule, actionIdItem)
	}
	var scalarRule []interface{}
	for _, scalarItem := range scalar {
		scalarRule = append(scalarRule, scalarItem)
	}

	logs, sub, err := _Metamodel.contract.FilterLogs(opts, "SignaledEvent", roleRule, actionIdRule, scalarRule)
	if err != nil {
		return nil, err
	}
	return &MetamodelSignaledEventIterator{contract: _Metamodel.contract, event: "SignaledEvent", logs: logs, sub: sub}, nil
}

// WatchSignaledEvent is a free log subscription operation binding the contract event 0x50e4a52b0772bed9f06a7d3f7dfa66d73640067a5cc77a73c24556ccc900fa08.
//
// Solidity: event SignaledEvent(uint8 indexed role, uint8 indexed actionId, uint256 indexed scalar, uint256 sequence)
func (_Metamodel *MetamodelFilterer) WatchSignaledEvent(opts *bind.WatchOpts, sink chan<- *MetamodelSignaledEvent, role []uint8, actionId []uint8, scalar []*big.Int) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var actionIdRule []interface{}
	for _, actionIdItem := range actionId {
		actionIdRule = append(actionIdRule, actionIdItem)
	}
	var scalarRule []interface{}
	for _, scalarItem := range scalar {
		scalarRule = append(scalarRule, scalarItem)
	}

	logs, sub, err := _Metamodel.contract.WatchLogs(opts, "SignaledEvent", roleRule, actionIdRule, scalarRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetamodelSignaledEvent)
				if err := _Metamodel.contract.UnpackLog(event, "SignaledEvent", log); err != nil {
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

// ParseSignaledEvent is a log parse operation binding the contract event 0x50e4a52b0772bed9f06a7d3f7dfa66d73640067a5cc77a73c24556ccc900fa08.
//
// Solidity: event SignaledEvent(uint8 indexed role, uint8 indexed actionId, uint256 indexed scalar, uint256 sequence)
func (_Metamodel *MetamodelFilterer) ParseSignaledEvent(log types.Log) (*MetamodelSignaledEvent, error) {
	event := new(MetamodelSignaledEvent)
	if err := _Metamodel.contract.UnpackLog(event, "SignaledEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
