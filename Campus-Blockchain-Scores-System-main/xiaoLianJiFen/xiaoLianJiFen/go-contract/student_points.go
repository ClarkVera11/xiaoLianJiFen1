// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// StudentPointsMetaData contains all meta data concerning the StudentPoints contract.
var StudentPointsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"username\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"phone\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roleName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"points\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"name\":\"updateUserInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"username\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"phone\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"roleName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"points\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"name\":\"UserInfoUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"}],\"name\":\"getUserInfo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"users\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"username\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"phone\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roleName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"points\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b50610dc68061001c5f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c80636386c1c714610043578063a87430ba14610077578063adfcdc30146100ab575b5f80fd5b61005d6004803603810190610058919061076a565b6100c7565b60405161006e95949392919061081d565b60405180910390f35b610091600480360381019061008c919061076a565b61038e565b6040516100a295949392919061081d565b60405180910390f35b6100c560048036038101906100c091906109e0565b6105d7565b005b60608060605f60605f805f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060a00160405290815f8201805461012690610af4565b80601f016020809104026020016040519081016040528092919081815260200182805461015290610af4565b801561019d5780601f106101745761010080835404028352916020019161019d565b820191905f5260205f20905b81548152906001019060200180831161018057829003601f168201915b505050505081526020016001820180546101b690610af4565b80601f01602080910402602001604051908101604052809291908181526020018280546101e290610af4565b801561022d5780601f106102045761010080835404028352916020019161022d565b820191905f5260205f20905b81548152906001019060200180831161021057829003601f168201915b5050505050815260200160028201805461024690610af4565b80601f016020809104026020016040519081016040528092919081815260200182805461027290610af4565b80156102bd5780601f10610294576101008083540402835291602001916102bd565b820191905f5260205f20905b8154815290600101906020018083116102a057829003601f168201915b50505050508152602001600382015481526020016004820180546102e090610af4565b80601f016020809104026020016040519081016040528092919081815260200182805461030c90610af4565b80156103575780601f1061032e57610100808354040283529160200191610357565b820191905f5260205f20905b81548152906001019060200180831161033a57829003601f168201915b5050505050815250509050805f01518160200151826040015183606001518460800151955095509550955095505091939590929450565b5f602052805f5260405f205f91509050805f0180546103ac90610af4565b80601f01602080910402602001604051908101604052809291908181526020018280546103d890610af4565b80156104235780601f106103fa57610100808354040283529160200191610423565b820191905f5260205f20905b81548152906001019060200180831161040657829003601f168201915b50505050509080600101805461043890610af4565b80601f016020809104026020016040519081016040528092919081815260200182805461046490610af4565b80156104af5780601f10610486576101008083540402835291602001916104af565b820191905f5260205f20905b81548152906001019060200180831161049257829003601f168201915b5050505050908060020180546104c490610af4565b80601f01602080910402602001604051908101604052809291908181526020018280546104f090610af4565b801561053b5780601f106105125761010080835404028352916020019161053b565b820191905f5260205f20905b81548152906001019060200180831161051e57829003601f168201915b50505050509080600301549080600401805461055690610af4565b80601f016020809104026020016040519081016040528092919081815260200182805461058290610af4565b80156105cd5780601f106105a4576101008083540402835291602001916105cd565b820191905f5260205f20905b8154815290600101906020018083116105b057829003601f168201915b5050505050905085565b5f3390506040518060a00160405280878152602001868152602001858152602001848152602001838152505f808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f820151815f0190816106519190610cc1565b5060208201518160010190816106679190610cc1565b50604082015181600201908161067d9190610cc1565b5060608201518160030155608082015181600401908161069d9190610cc1565b509050508073ffffffffffffffffffffffffffffffffffffffff167faf667c3c6725b843a118988c12bce6de047d69da021c8fef35514da78585257b87878787876040516106ef95949392919061081d565b60405180910390a2505050505050565b5f604051905090565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61073982610710565b9050919050565b6107498161072f565b8114610753575f80fd5b50565b5f8135905061076481610740565b92915050565b5f6020828403121561077f5761077e610708565b5b5f61078c84828501610756565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6107d782610795565b6107e1818561079f565b93506107f18185602086016107af565b6107fa816107bd565b840191505092915050565b5f819050919050565b61081781610805565b82525050565b5f60a0820190508181035f83015261083581886107cd565b9050818103602083015261084981876107cd565b9050818103604083015261085d81866107cd565b905061086c606083018561080e565b818103608083015261087e81846107cd565b90509695505050505050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6108c8826107bd565b810181811067ffffffffffffffff821117156108e7576108e6610892565b5b80604052505050565b5f6108f96106ff565b905061090582826108bf565b919050565b5f67ffffffffffffffff82111561092457610923610892565b5b61092d826107bd565b9050602081019050919050565b828183375f83830152505050565b5f61095a6109558461090a565b6108f0565b9050828152602081018484840111156109765761097561088e565b5b61098184828561093a565b509392505050565b5f82601f83011261099d5761099c61088a565b5b81356109ad848260208601610948565b91505092915050565b6109bf81610805565b81146109c9575f80fd5b50565b5f813590506109da816109b6565b92915050565b5f805f805f60a086880312156109f9576109f8610708565b5b5f86013567ffffffffffffffff811115610a1657610a1561070c565b5b610a2288828901610989565b955050602086013567ffffffffffffffff811115610a4357610a4261070c565b5b610a4f88828901610989565b945050604086013567ffffffffffffffff811115610a7057610a6f61070c565b5b610a7c88828901610989565b9350506060610a8d888289016109cc565b925050608086013567ffffffffffffffff811115610aae57610aad61070c565b5b610aba88828901610989565b9150509295509295909350565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610b0b57607f821691505b602082108103610b1e57610b1d610ac7565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302610b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610b45565b610b8a8683610b45565b95508019841693508086168417925050509392505050565b5f819050919050565b5f610bc5610bc0610bbb84610805565b610ba2565b610805565b9050919050565b5f819050919050565b610bde83610bab565b610bf2610bea82610bcc565b848454610b51565b825550505050565b5f90565b610c06610bfa565b610c11818484610bd5565b505050565b5b81811015610c3457610c295f82610bfe565b600181019050610c17565b5050565b601f821115610c7957610c4a81610b24565b610c5384610b36565b81016020851015610c62578190505b610c76610c6e85610b36565b830182610c16565b50505b505050565b5f82821c905092915050565b5f610c995f1984600802610c7e565b1980831691505092915050565b5f610cb18383610c8a565b9150826002028217905092915050565b610cca82610795565b67ffffffffffffffff811115610ce357610ce2610892565b5b610ced8254610af4565b610cf8828285610c38565b5f60209050601f831160018114610d29575f8415610d17578287015190505b610d218582610ca6565b865550610d88565b601f198416610d3786610b24565b5f5b82811015610d5e57848901518255600182019150602085019450602081019050610d39565b86831015610d7b5784890151610d77601f891682610c8a565b8355505b6001600288020188555050505b50505050505056fea2646970667358221220daf2522d1dc052211e8b33dc9fb58a252e9177ec9957a80ed0af702c13ca1b3364736f6c634300081a0033",
}

// StudentPointsABI is the input ABI used to generate the binding from.
// Deprecated: Use StudentPointsMetaData.ABI instead.
var StudentPointsABI = StudentPointsMetaData.ABI

// StudentPointsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StudentPointsMetaData.Bin instead.
var StudentPointsBin = StudentPointsMetaData.Bin

// DeployStudentPoints deploys a new Ethereum contract, binding an instance of StudentPoints to it.
func DeployStudentPoints(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StudentPoints, error) {
	parsed, err := StudentPointsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StudentPointsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StudentPoints{StudentPointsCaller: StudentPointsCaller{contract: contract}, StudentPointsTransactor: StudentPointsTransactor{contract: contract}, StudentPointsFilterer: StudentPointsFilterer{contract: contract}}, nil
}

// StudentPoints is an auto generated Go binding around an Ethereum contract.
type StudentPoints struct {
	StudentPointsCaller     // Read-only binding to the contract
	StudentPointsTransactor // Write-only binding to the contract
	StudentPointsFilterer   // Log filterer for contract events
}

// StudentPointsCaller is an auto generated read-only Go binding around an Ethereum contract.
type StudentPointsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StudentPointsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StudentPointsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StudentPointsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StudentPointsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StudentPointsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StudentPointsSession struct {
	Contract     *StudentPoints    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StudentPointsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StudentPointsCallerSession struct {
	Contract *StudentPointsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StudentPointsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StudentPointsTransactorSession struct {
	Contract     *StudentPointsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StudentPointsRaw is an auto generated low-level Go binding around an Ethereum contract.
type StudentPointsRaw struct {
	Contract *StudentPoints // Generic contract binding to access the raw methods on
}

// StudentPointsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StudentPointsCallerRaw struct {
	Contract *StudentPointsCaller // Generic read-only contract binding to access the raw methods on
}

// StudentPointsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StudentPointsTransactorRaw struct {
	Contract *StudentPointsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStudentPoints creates a new instance of StudentPoints, bound to a specific deployed contract.
func NewStudentPoints(address common.Address, backend bind.ContractBackend) (*StudentPoints, error) {
	contract, err := bindStudentPoints(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StudentPoints{StudentPointsCaller: StudentPointsCaller{contract: contract}, StudentPointsTransactor: StudentPointsTransactor{contract: contract}, StudentPointsFilterer: StudentPointsFilterer{contract: contract}}, nil
}

// NewStudentPointsCaller creates a new read-only instance of StudentPoints, bound to a specific deployed contract.
func NewStudentPointsCaller(address common.Address, caller bind.ContractCaller) (*StudentPointsCaller, error) {
	contract, err := bindStudentPoints(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StudentPointsCaller{contract: contract}, nil
}

// NewStudentPointsTransactor creates a new write-only instance of StudentPoints, bound to a specific deployed contract.
func NewStudentPointsTransactor(address common.Address, transactor bind.ContractTransactor) (*StudentPointsTransactor, error) {
	contract, err := bindStudentPoints(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StudentPointsTransactor{contract: contract}, nil
}

// NewStudentPointsFilterer creates a new log filterer instance of StudentPoints, bound to a specific deployed contract.
func NewStudentPointsFilterer(address common.Address, filterer bind.ContractFilterer) (*StudentPointsFilterer, error) {
	contract, err := bindStudentPoints(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StudentPointsFilterer{contract: contract}, nil
}

// bindStudentPoints binds a generic wrapper to an already deployed contract.
func bindStudentPoints(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StudentPointsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StudentPoints *StudentPointsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StudentPoints.Contract.StudentPointsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StudentPoints *StudentPointsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StudentPoints.Contract.StudentPointsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StudentPoints *StudentPointsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StudentPoints.Contract.StudentPointsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StudentPoints *StudentPointsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StudentPoints.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StudentPoints *StudentPointsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StudentPoints.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StudentPoints *StudentPointsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StudentPoints.Contract.contract.Transact(opts, method, params...)
}

// GetUserInfo is a free data retrieval call binding the contract method 0x6386c1c7.
//
// Solidity: function getUserInfo(address userAddr) view returns(string, string, string, uint256, string)
func (_StudentPoints *StudentPointsCaller) GetUserInfo(opts *bind.CallOpts, userAddr common.Address) (string, string, string, *big.Int, string, error) {
	var out []interface{}
	err := _StudentPoints.contract.Call(opts, &out, "getUserInfo", userAddr)

	if err != nil {
		return *new(string), *new(string), *new(string), *new(*big.Int), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(string)).(*string)

	return out0, out1, out2, out3, out4, err

}

// GetUserInfo is a free data retrieval call binding the contract method 0x6386c1c7.
//
// Solidity: function getUserInfo(address userAddr) view returns(string, string, string, uint256, string)
func (_StudentPoints *StudentPointsSession) GetUserInfo(userAddr common.Address) (string, string, string, *big.Int, string, error) {
	return _StudentPoints.Contract.GetUserInfo(&_StudentPoints.CallOpts, userAddr)
}

// GetUserInfo is a free data retrieval call binding the contract method 0x6386c1c7.
//
// Solidity: function getUserInfo(address userAddr) view returns(string, string, string, uint256, string)
func (_StudentPoints *StudentPointsCallerSession) GetUserInfo(userAddr common.Address) (string, string, string, *big.Int, string, error) {
	return _StudentPoints.Contract.GetUserInfo(&_StudentPoints.CallOpts, userAddr)
}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(string username, string phone, string roleName, uint256 points, string title)
func (_StudentPoints *StudentPointsCaller) Users(opts *bind.CallOpts, arg0 common.Address) (struct {
	Username string
	Phone    string
	RoleName string
	Points   *big.Int
	Title    string
}, error) {
	var out []interface{}
	err := _StudentPoints.contract.Call(opts, &out, "users", arg0)

	outstruct := new(struct {
		Username string
		Phone    string
		RoleName string
		Points   *big.Int
		Title    string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Username = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Phone = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.RoleName = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Points = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Title = *abi.ConvertType(out[4], new(string)).(*string)

	return *outstruct, err

}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(string username, string phone, string roleName, uint256 points, string title)
func (_StudentPoints *StudentPointsSession) Users(arg0 common.Address) (struct {
	Username string
	Phone    string
	RoleName string
	Points   *big.Int
	Title    string
}, error) {
	return _StudentPoints.Contract.Users(&_StudentPoints.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(string username, string phone, string roleName, uint256 points, string title)
func (_StudentPoints *StudentPointsCallerSession) Users(arg0 common.Address) (struct {
	Username string
	Phone    string
	RoleName string
	Points   *big.Int
	Title    string
}, error) {
	return _StudentPoints.Contract.Users(&_StudentPoints.CallOpts, arg0)
}

// UpdateUserInfo is a paid mutator transaction binding the contract method 0xadfcdc30.
//
// Solidity: function updateUserInfo(string username, string phone, string roleName, uint256 points, string title) returns()
func (_StudentPoints *StudentPointsTransactor) UpdateUserInfo(opts *bind.TransactOpts, username string, phone string, roleName string, points *big.Int, title string) (*types.Transaction, error) {
	return _StudentPoints.contract.Transact(opts, "updateUserInfo", username, phone, roleName, points, title)
}

// UpdateUserInfo is a paid mutator transaction binding the contract method 0xadfcdc30.
//
// Solidity: function updateUserInfo(string username, string phone, string roleName, uint256 points, string title) returns()
func (_StudentPoints *StudentPointsSession) UpdateUserInfo(username string, phone string, roleName string, points *big.Int, title string) (*types.Transaction, error) {
	return _StudentPoints.Contract.UpdateUserInfo(&_StudentPoints.TransactOpts, username, phone, roleName, points, title)
}

// UpdateUserInfo is a paid mutator transaction binding the contract method 0xadfcdc30.
//
// Solidity: function updateUserInfo(string username, string phone, string roleName, uint256 points, string title) returns()
func (_StudentPoints *StudentPointsTransactorSession) UpdateUserInfo(username string, phone string, roleName string, points *big.Int, title string) (*types.Transaction, error) {
	return _StudentPoints.Contract.UpdateUserInfo(&_StudentPoints.TransactOpts, username, phone, roleName, points, title)
}

// StudentPointsUserInfoUpdatedIterator is returned from FilterUserInfoUpdated and is used to iterate over the raw logs and unpacked data for UserInfoUpdated events raised by the StudentPoints contract.
type StudentPointsUserInfoUpdatedIterator struct {
	Event *StudentPointsUserInfoUpdated // Event containing the contract specifics and raw log

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
func (it *StudentPointsUserInfoUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StudentPointsUserInfoUpdated)
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
		it.Event = new(StudentPointsUserInfoUpdated)
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
func (it *StudentPointsUserInfoUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StudentPointsUserInfoUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StudentPointsUserInfoUpdated represents a UserInfoUpdated event raised by the StudentPoints contract.
type StudentPointsUserInfoUpdated struct {
	User     common.Address
	Username string
	Phone    string
	RoleName string
	Points   *big.Int
	Title    string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUserInfoUpdated is a free log retrieval operation binding the contract event 0xaf667c3c6725b843a118988c12bce6de047d69da021c8fef35514da78585257b.
//
// Solidity: event UserInfoUpdated(address indexed user, string username, string phone, string roleName, uint256 points, string title)
func (_StudentPoints *StudentPointsFilterer) FilterUserInfoUpdated(opts *bind.FilterOpts, user []common.Address) (*StudentPointsUserInfoUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StudentPoints.contract.FilterLogs(opts, "UserInfoUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return &StudentPointsUserInfoUpdatedIterator{contract: _StudentPoints.contract, event: "UserInfoUpdated", logs: logs, sub: sub}, nil
}

// WatchUserInfoUpdated is a free log subscription operation binding the contract event 0xaf667c3c6725b843a118988c12bce6de047d69da021c8fef35514da78585257b.
//
// Solidity: event UserInfoUpdated(address indexed user, string username, string phone, string roleName, uint256 points, string title)
func (_StudentPoints *StudentPointsFilterer) WatchUserInfoUpdated(opts *bind.WatchOpts, sink chan<- *StudentPointsUserInfoUpdated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StudentPoints.contract.WatchLogs(opts, "UserInfoUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StudentPointsUserInfoUpdated)
				if err := _StudentPoints.contract.UnpackLog(event, "UserInfoUpdated", log); err != nil {
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

// ParseUserInfoUpdated is a log parse operation binding the contract event 0xaf667c3c6725b843a118988c12bce6de047d69da021c8fef35514da78585257b.
//
// Solidity: event UserInfoUpdated(address indexed user, string username, string phone, string roleName, uint256 points, string title)
func (_StudentPoints *StudentPointsFilterer) ParseUserInfoUpdated(log types.Log) (*StudentPointsUserInfoUpdated, error) {
	event := new(StudentPointsUserInfoUpdated)
	if err := _StudentPoints.contract.UnpackLog(event, "UserInfoUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
