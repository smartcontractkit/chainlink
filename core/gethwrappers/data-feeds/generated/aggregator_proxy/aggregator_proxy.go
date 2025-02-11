// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package aggregator_proxy

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

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

var AggregatorProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_accessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"aggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"confirmAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"name\":\"phaseAggregators\",\"outputs\":[{\"internalType\":\"contractAggregatorV2V3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"phaseId\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"proposeAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposedAggregator\",\"outputs\":[{\"internalType\":\"contractAggregatorV2V3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"proposedGetRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposedLatestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_accessController\",\"type\":\"address\"}],\"name\":\"setController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620026e8380380620026e8833981810160405260408110156200003757600080fd5b508051602090910151600080546001600160a01b031916331790558162000067816001600160e01b036200008416565b506200007c816001600160e01b03620000f316565b505062000175565b60028054604080518082018252600161ffff80851691909101168082526001600160a01b0395909516602091820181905261ffff19909316851762010000600160b01b0319166201000084021790935560009384526004909252912080546001600160a01b0319169091179055565b6000546001600160a01b0316331462000153576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600580546001600160a01b0319166001600160a01b0392909216919091179055565b61256380620001856000396000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c80638f6b4d91116100e3578063bc43cbaf1161008c578063f2fde38b11610066578063f2fde38b1461042b578063f8a2abd31461045e578063feaf968c146104915761018d565b8063bc43cbaf146103fa578063c159730414610402578063e8c4be30146104235761018d565b8063a928c096116100bd578063a928c0961461038d578063b5ab58dc146103c0578063b633620c146103dd5761018d565b80638f6b4d911461032957806392eefe9b146103315780639a6fc8f5146103645761018d565b80636001ac531161014557806379ba50971161011f57806379ba50971461030f5780638205bf6a146103195780638da5cb5b146103215761018d565b80636001ac5314610222578063668a0f021461028a5780637284e416146102925761018d565b806350d25bcd1161017657806350d25bcd146101e157806354fd4d50146101fb57806358303b10146102035761018d565b8063245a7bfc14610192578063313ce567146101c3575b600080fd5b61019a610499565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101cb6104bb565b6040805160ff9092168252519081900360200190f35b6101e9610559565b60408051918252519081900360200190f35b6101e96106e0565b61020b61074d565b6040805161ffff9092168252519081900360200190f35b61024b6004803603602081101561023857600080fd5b503569ffffffffffffffffffff16610757565b6040805169ffffffffffffffffffff96871681526020810195909552848101939093526060840191909152909216608082015290519081900360a00190f35b6101e9610978565b61029a610af9565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102d45781810151838201526020016102bc565b50505050905090810190601f1680156103015780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610317610c76565b005b6101e9610d78565b61019a610ef9565b61024b610f15565b6103176004803603602081101561034757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611134565b61024b6004803603602081101561037a57600080fd5b503569ffffffffffffffffffff16611201565b610317600480360360208110156103a357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661138b565b6101e9600480360360208110156103d657600080fd5b50356114ce565b6101e9600480360360208110156103f357600080fd5b5035611657565b61019a6117d9565b61019a6004803603602081101561041857600080fd5b503561ffff166117f5565b61019a61181d565b6103176004803603602081101561044157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611839565b6103176004803603602081101561047457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611935565b61024b611a02565b60025462010000900473ffffffffffffffffffffffffffffffffffffffff1690565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561052857600080fd5b505afa15801561053c573d6000803e3d6000fd5b505050506040513d602081101561055257600080fd5b5051905090565b60055460009073ffffffffffffffffffffffffffffffffffffffff168015806106675750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561063a57600080fd5b505afa15801561064e573d6000803e3d6000fd5b505050506040513d602081101561066457600080fd5b50515b6106d257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6106da611b8b565b91505090565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166354fd4d506040518163ffffffff1660e01b815260040160206040518083038186803b15801561052857600080fd5b60025461ffff1690565b600554600090819081908190819073ffffffffffffffffffffffffffffffffffffffff1680158061086d5750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561084057600080fd5b505afa158015610854573d6000803e3d6000fd5b505050506040513d602081101561086a57600080fd5b50515b6108d857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff1661095c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4e6f2070726f706f7365642061676772656761746f722070726573656e740000604482015290519081900360640190fd5b61096587611bf8565b939b929a50909850965090945092505050565b60055460009073ffffffffffffffffffffffffffffffffffffffff16801580610a865750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610a5957600080fd5b505afa158015610a6d573d6000803e3d6000fd5b505050506040513d6020811015610a8357600080fd5b50515b610af157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6106da611d57565b6060600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16637284e4166040518163ffffffff1660e01b815260040160006040518083038186803b158015610b6657600080fd5b505afa158015610b7a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526020811015610bc157600080fd5b8101908080516040519392919084640100000000821115610be157600080fd5b908301906020820185811115610bf657600080fd5b8251640100000000811182820188101715610c1057600080fd5b82525081516020918201929091019080838360005b83811015610c3d578181015183820152602001610c25565b50505050905090810190601f168015610c6a5780820380516001836020036101000a031916815260200191505b50604052505050905090565b60015473ffffffffffffffffffffffffffffffffffffffff163314610cfc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60055460009073ffffffffffffffffffffffffffffffffffffffff16801580610e865750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610e5957600080fd5b505afa158015610e6d573d6000803e3d6000fd5b505050506040513d6020811015610e8357600080fd5b50515b610ef157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6106da611e2e565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b600554600090819081908190819073ffffffffffffffffffffffffffffffffffffffff1680158061102b5750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610ffe57600080fd5b505afa158015611012573d6000803e3d6000fd5b505050506040513d602081101561102857600080fd5b50515b61109657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff1661111a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4e6f2070726f706f7365642061676772656761746f722070726573656e740000604482015290519081900360640190fd5b611122611e9b565b95509550955095509550509091929394565b60005473ffffffffffffffffffffffffffffffffffffffff1633146111ba57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600554600090819081908190819073ffffffffffffffffffffffffffffffffffffffff168015806113175750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b1580156112ea57600080fd5b505afa1580156112fe573d6000803e3d6000fd5b505050506040513d602081101561131457600080fd5b50515b61138257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61096587611fe4565b60005473ffffffffffffffffffffffffffffffffffffffff16331461141157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff82811691161461149a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f496e76616c69642070726f706f7365642061676772656761746f720000000000604482015290519081900360640190fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556114cb81612117565b50565b60055460009073ffffffffffffffffffffffffffffffffffffffff168015806115dc5750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b1580156115af57600080fd5b505afa1580156115c3573d6000803e3d6000fd5b505050506040513d60208110156115d957600080fd5b50515b61164757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611650836121de565b9392505050565b60055460009073ffffffffffffffffffffffffffffffffffffffff168015806117655750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561173857600080fd5b505afa15801561174c573d6000803e3d6000fd5b505050506040513d602081101561176257600080fd5b50515b6117d057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611650836122d8565b60055473ffffffffffffffffffffffffffffffffffffffff1681565b60046020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b60035473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146118bf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633146119bb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600554600090819081908190819073ffffffffffffffffffffffffffffffffffffffff16801580611b185750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015611aeb57600080fd5b505afa158015611aff573d6000803e3d6000fd5b505050506040513d6020811015611b1557600080fd5b50515b611b8357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61112261239b565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b815260040160206040518083038186803b15801561052857600080fd5b600354600090819081908190819073ffffffffffffffffffffffffffffffffffffffff16611c8757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4e6f2070726f706f7365642061676772656761746f722070726573656e740000604482015290519081900360640190fd5b600354604080517f9a6fc8f500000000000000000000000000000000000000000000000000000000815269ffffffffffffffffffff89166004820152905173ffffffffffffffffffffffffffffffffffffffff90921691639a6fc8f59160248082019260a092909190829003018186803b158015611d0457600080fd5b505afa158015611d18573d6000803e3d6000fd5b505050506040513d60a0811015611d2e57600080fd5b508051602082015160408301516060840151608090940151929a91995097509195509350915050565b6000611d61612516565b5060408051808201825260025461ffff81168083526201000090910473ffffffffffffffffffffffffffffffffffffffff16602080840182905284517f668a0f0200000000000000000000000000000000000000000000000000000000815294519394611e1c9463668a0f0292600480840193919291829003018186803b158015611deb57600080fd5b505afa158015611dff573d6000803e3d6000fd5b505050506040513d6020811015611e1557600080fd5b50516124b8565b69ffffffffffffffffffff1691505090565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638205bf6a6040518163ffffffff1660e01b815260040160206040518083038186803b15801561052857600080fd5b600354600090819081908190819073ffffffffffffffffffffffffffffffffffffffff16611f2a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4e6f2070726f706f7365642061676772656761746f722070726573656e740000604482015290519081900360640190fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015611f9257600080fd5b505afa158015611fa6573d6000803e3d6000fd5b505050506040513d60a0811015611fbc57600080fd5b5080516020820151604083015160608401516080909401519299919850965091945092509050565b60008060008060008060006120048869ffffffffffffffffffff166124d8565b61ffff821660009081526004602081905260408083205481517f9a6fc8f500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86169381019390935290519496509294509092839283928392839273ffffffffffffffffffffffffffffffffffffffff1691639a6fc8f59160248083019260a0929190829003018186803b1580156120a057600080fd5b505afa1580156120b4573d6000803e3d6000fd5b505050506040513d60a08110156120ca57600080fd5b508051602082015160408301516060840151608090940151929850909650945090925090506120fd85858585858c6124e0565b9b509b509b509b509b505050505050505091939590929450565b60028054604080518082018252600161ffff808516919091011680825273ffffffffffffffffffffffffffffffffffffffff9590951660209182018190527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090931685177fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000084021790935560009384526004909252912080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169091179055565b600069ffffffffffffffffffff8211156121fa575060006122d3565b600080612206846124d8565b61ffff8216600090815260046020526040902054919350915073ffffffffffffffffffffffffffffffffffffffff168061224657600093505050506122d3565b8073ffffffffffffffffffffffffffffffffffffffff1663b5ab58dc836040518263ffffffff1660e01b8152600401808267ffffffffffffffff16815260200191505060206040518083038186803b1580156122a157600080fd5b505afa1580156122b5573d6000803e3d6000fd5b505050506040513d60208110156122cb57600080fd5b505193505050505b919050565b600069ffffffffffffffffffff8211156122f4575060006122d3565b600080612300846124d8565b61ffff8216600090815260046020526040902054919350915073ffffffffffffffffffffffffffffffffffffffff168061234057600093505050506122d3565b8073ffffffffffffffffffffffffffffffffffffffff1663b633620c836040518263ffffffff1660e01b8152600401808267ffffffffffffffff16815260200191505060206040518083038186803b1580156122a157600080fd5b60008060008060006123ab612516565b5060408051808201825260025461ffff8116825262010000900473ffffffffffffffffffffffffffffffffffffffff166020820181905282517ffeaf968c0000000000000000000000000000000000000000000000000000000081529251919260009283928392839283929163feaf968c9160048083019260a0929190829003018186803b15801561243c57600080fd5b505afa158015612450573d6000803e3d6000fd5b505050506040513d60a081101561246657600080fd5b5080516020820151604083015160608401516080909401518a5193995091975095509193509091506124a190869086908690869086906124e0565b9a509a509a509a509a505050505050509091929394565b67ffffffffffffffff1660409190911b69ffff0000000000000000161790565b604081901c91565b60008060008060006124f2868c6124b8565b8a8a8a6124ff8a8c6124b8565b939f929e50909c509a509098509650505050505050565b60408051808201909152600080825260208201529056fea2646970667358221220c6148a0e63011d3b8b4f67078be31115256b163e26351db6fe3b70d7faf433f964736f6c63430006060033",
}

var AggregatorProxyABI = AggregatorProxyMetaData.ABI

var AggregatorProxyBin = AggregatorProxyMetaData.Bin

func DeployAggregatorProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _aggregator common.Address, _accessController common.Address) (common.Address, *types.Transaction, *AggregatorProxy, error) {
	parsed, err := AggregatorProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AggregatorProxyBin), backend, _aggregator, _accessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AggregatorProxy{address: address, abi: *parsed, AggregatorProxyCaller: AggregatorProxyCaller{contract: contract}, AggregatorProxyTransactor: AggregatorProxyTransactor{contract: contract}, AggregatorProxyFilterer: AggregatorProxyFilterer{contract: contract}}, nil
}

type AggregatorProxy struct {
	address common.Address
	abi     abi.ABI
	AggregatorProxyCaller
	AggregatorProxyTransactor
	AggregatorProxyFilterer
}

type AggregatorProxyCaller struct {
	contract *bind.BoundContract
}

type AggregatorProxyTransactor struct {
	contract *bind.BoundContract
}

type AggregatorProxyFilterer struct {
	contract *bind.BoundContract
}

type AggregatorProxySession struct {
	Contract     *AggregatorProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AggregatorProxyCallerSession struct {
	Contract *AggregatorProxyCaller
	CallOpts bind.CallOpts
}

type AggregatorProxyTransactorSession struct {
	Contract     *AggregatorProxyTransactor
	TransactOpts bind.TransactOpts
}

type AggregatorProxyRaw struct {
	Contract *AggregatorProxy
}

type AggregatorProxyCallerRaw struct {
	Contract *AggregatorProxyCaller
}

type AggregatorProxyTransactorRaw struct {
	Contract *AggregatorProxyTransactor
}

func NewAggregatorProxy(address common.Address, backend bind.ContractBackend) (*AggregatorProxy, error) {
	abi, err := abi.JSON(strings.NewReader(AggregatorProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAggregatorProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxy{address: address, abi: abi, AggregatorProxyCaller: AggregatorProxyCaller{contract: contract}, AggregatorProxyTransactor: AggregatorProxyTransactor{contract: contract}, AggregatorProxyFilterer: AggregatorProxyFilterer{contract: contract}}, nil
}

func NewAggregatorProxyCaller(address common.Address, caller bind.ContractCaller) (*AggregatorProxyCaller, error) {
	contract, err := bindAggregatorProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyCaller{contract: contract}, nil
}

func NewAggregatorProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorProxyTransactor, error) {
	contract, err := bindAggregatorProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyTransactor{contract: contract}, nil
}

func NewAggregatorProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorProxyFilterer, error) {
	contract, err := bindAggregatorProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyFilterer{contract: contract}, nil
}

func bindAggregatorProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AggregatorProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AggregatorProxy *AggregatorProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorProxy.Contract.AggregatorProxyCaller.contract.Call(opts, result, method, params...)
}

func (_AggregatorProxy *AggregatorProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AggregatorProxyTransactor.contract.Transfer(opts)
}

func (_AggregatorProxy *AggregatorProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AggregatorProxyTransactor.contract.Transact(opts, method, params...)
}

func (_AggregatorProxy *AggregatorProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_AggregatorProxy *AggregatorProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.contract.Transfer(opts)
}

func (_AggregatorProxy *AggregatorProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.contract.Transact(opts, method, params...)
}

func (_AggregatorProxy *AggregatorProxyCaller) AccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "accessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) AccessController() (common.Address, error) {
	return _AggregatorProxy.Contract.AccessController(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) AccessController() (common.Address, error) {
	return _AggregatorProxy.Contract.AccessController(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) Aggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "aggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) Aggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.Aggregator(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) Aggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.Aggregator(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) Decimals() (uint8, error) {
	return _AggregatorProxy.Contract.Decimals(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) Decimals() (uint8, error) {
	return _AggregatorProxy.Contract.Decimals(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) Description() (string, error) {
	return _AggregatorProxy.Contract.Description(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) Description() (string, error) {
	return _AggregatorProxy.Contract.Description(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetAnswer(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetAnswer(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(GetRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_AggregatorProxy *AggregatorProxySession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _AggregatorProxy.Contract.GetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _AggregatorProxy.Contract.GetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetTimestamp(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetTimestamp(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) LatestAnswer() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestAnswer(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestAnswer(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) LatestRound() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestRound(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestRound(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_AggregatorProxy *AggregatorProxySession) LatestRoundData() (LatestRoundData,

	error) {
	return _AggregatorProxy.Contract.LatestRoundData(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _AggregatorProxy.Contract.LatestRoundData(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestTimestamp(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestTimestamp(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) Owner() (common.Address, error) {
	return _AggregatorProxy.Contract.Owner(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) Owner() (common.Address, error) {
	return _AggregatorProxy.Contract.Owner(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) PhaseAggregators(opts *bind.CallOpts, arg0 uint16) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "phaseAggregators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) PhaseAggregators(arg0 uint16) (common.Address, error) {
	return _AggregatorProxy.Contract.PhaseAggregators(&_AggregatorProxy.CallOpts, arg0)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) PhaseAggregators(arg0 uint16) (common.Address, error) {
	return _AggregatorProxy.Contract.PhaseAggregators(&_AggregatorProxy.CallOpts, arg0)
}

func (_AggregatorProxy *AggregatorProxyCaller) PhaseId(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "phaseId")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) PhaseId() (uint16, error) {
	return _AggregatorProxy.Contract.PhaseId(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) PhaseId() (uint16, error) {
	return _AggregatorProxy.Contract.PhaseId(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) ProposedAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "proposedAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) ProposedAggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.ProposedAggregator(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) ProposedAggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.ProposedAggregator(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) ProposedGetRoundData(opts *bind.CallOpts, _roundId *big.Int) (ProposedGetRoundData,

	error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "proposedGetRoundData", _roundId)

	outstruct := new(ProposedGetRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_AggregatorProxy *AggregatorProxySession) ProposedGetRoundData(_roundId *big.Int) (ProposedGetRoundData,

	error) {
	return _AggregatorProxy.Contract.ProposedGetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) ProposedGetRoundData(_roundId *big.Int) (ProposedGetRoundData,

	error) {
	return _AggregatorProxy.Contract.ProposedGetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

func (_AggregatorProxy *AggregatorProxyCaller) ProposedLatestRoundData(opts *bind.CallOpts) (ProposedLatestRoundData,

	error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "proposedLatestRoundData")

	outstruct := new(ProposedLatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_AggregatorProxy *AggregatorProxySession) ProposedLatestRoundData() (ProposedLatestRoundData,

	error) {
	return _AggregatorProxy.Contract.ProposedLatestRoundData(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) ProposedLatestRoundData() (ProposedLatestRoundData,

	error) {
	return _AggregatorProxy.Contract.ProposedLatestRoundData(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AggregatorProxy *AggregatorProxySession) Version() (*big.Int, error) {
	return _AggregatorProxy.Contract.Version(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyCallerSession) Version() (*big.Int, error) {
	return _AggregatorProxy.Contract.Version(&_AggregatorProxy.CallOpts)
}

func (_AggregatorProxy *AggregatorProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "acceptOwnership")
}

func (_AggregatorProxy *AggregatorProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AcceptOwnership(&_AggregatorProxy.TransactOpts)
}

func (_AggregatorProxy *AggregatorProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AcceptOwnership(&_AggregatorProxy.TransactOpts)
}

func (_AggregatorProxy *AggregatorProxyTransactor) ConfirmAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "confirmAggregator", _aggregator)
}

func (_AggregatorProxy *AggregatorProxySession) ConfirmAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ConfirmAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

func (_AggregatorProxy *AggregatorProxyTransactorSession) ConfirmAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ConfirmAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

func (_AggregatorProxy *AggregatorProxyTransactor) ProposeAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "proposeAggregator", _aggregator)
}

func (_AggregatorProxy *AggregatorProxySession) ProposeAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ProposeAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

func (_AggregatorProxy *AggregatorProxyTransactorSession) ProposeAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ProposeAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

func (_AggregatorProxy *AggregatorProxyTransactor) SetController(opts *bind.TransactOpts, _accessController common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "setController", _accessController)
}

func (_AggregatorProxy *AggregatorProxySession) SetController(_accessController common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.SetController(&_AggregatorProxy.TransactOpts, _accessController)
}

func (_AggregatorProxy *AggregatorProxyTransactorSession) SetController(_accessController common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.SetController(&_AggregatorProxy.TransactOpts, _accessController)
}

func (_AggregatorProxy *AggregatorProxyTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "transferOwnership", _to)
}

func (_AggregatorProxy *AggregatorProxySession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.TransferOwnership(&_AggregatorProxy.TransactOpts, _to)
}

func (_AggregatorProxy *AggregatorProxyTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.TransferOwnership(&_AggregatorProxy.TransactOpts, _to)
}

type AggregatorProxyAnswerUpdatedIterator struct {
	Event *AggregatorProxyAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AggregatorProxyAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyAnswerUpdated)
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

	select {
	case log := <-it.logs:
		it.Event = new(AggregatorProxyAnswerUpdated)
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

func (it *AggregatorProxyAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *AggregatorProxyAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AggregatorProxyAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_AggregatorProxy *AggregatorProxyFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorProxyAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyAnswerUpdatedIterator{contract: _AggregatorProxy.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_AggregatorProxy *AggregatorProxyFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorProxyAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AggregatorProxyAnswerUpdated)
				if err := _AggregatorProxy.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_AggregatorProxy *AggregatorProxyFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorProxyAnswerUpdated, error) {
	event := new(AggregatorProxyAnswerUpdated)
	if err := _AggregatorProxy.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AggregatorProxyNewRoundIterator struct {
	Event *AggregatorProxyNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AggregatorProxyNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyNewRound)
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

	select {
	case log := <-it.logs:
		it.Event = new(AggregatorProxyNewRound)
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

func (it *AggregatorProxyNewRoundIterator) Error() error {
	return it.fail
}

func (it *AggregatorProxyNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AggregatorProxyNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_AggregatorProxy *AggregatorProxyFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorProxyNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyNewRoundIterator{contract: _AggregatorProxy.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_AggregatorProxy *AggregatorProxyFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorProxyNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AggregatorProxyNewRound)
				if err := _AggregatorProxy.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_AggregatorProxy *AggregatorProxyFilterer) ParseNewRound(log types.Log) (*AggregatorProxyNewRound, error) {
	event := new(AggregatorProxyNewRound)
	if err := _AggregatorProxy.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AggregatorProxyOwnershipTransferRequestedIterator struct {
	Event *AggregatorProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AggregatorProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyOwnershipTransferRequested)
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

	select {
	case log := <-it.logs:
		it.Event = new(AggregatorProxyOwnershipTransferRequested)
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

func (it *AggregatorProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AggregatorProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AggregatorProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AggregatorProxy *AggregatorProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AggregatorProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyOwnershipTransferRequestedIterator{contract: _AggregatorProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AggregatorProxy *AggregatorProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AggregatorProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AggregatorProxyOwnershipTransferRequested)
				if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AggregatorProxy *AggregatorProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*AggregatorProxyOwnershipTransferRequested, error) {
	event := new(AggregatorProxyOwnershipTransferRequested)
	if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AggregatorProxyOwnershipTransferredIterator struct {
	Event *AggregatorProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AggregatorProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyOwnershipTransferred)
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

	select {
	case log := <-it.logs:
		it.Event = new(AggregatorProxyOwnershipTransferred)
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

func (it *AggregatorProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AggregatorProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AggregatorProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AggregatorProxy *AggregatorProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AggregatorProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyOwnershipTransferredIterator{contract: _AggregatorProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AggregatorProxy *AggregatorProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AggregatorProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AggregatorProxyOwnershipTransferred)
				if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AggregatorProxy *AggregatorProxyFilterer) ParseOwnershipTransferred(log types.Log) (*AggregatorProxyOwnershipTransferred, error) {
	event := new(AggregatorProxyOwnershipTransferred)
	if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type ProposedGetRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type ProposedLatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_AggregatorProxy *AggregatorProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AggregatorProxy.abi.Events["AnswerUpdated"].ID:
		return _AggregatorProxy.ParseAnswerUpdated(log)
	case _AggregatorProxy.abi.Events["NewRound"].ID:
		return _AggregatorProxy.ParseNewRound(log)
	case _AggregatorProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _AggregatorProxy.ParseOwnershipTransferRequested(log)
	case _AggregatorProxy.abi.Events["OwnershipTransferred"].ID:
		return _AggregatorProxy.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AggregatorProxyAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (AggregatorProxyNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (AggregatorProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AggregatorProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_AggregatorProxy *AggregatorProxy) Address() common.Address {
	return _AggregatorProxy.address
}

type AggregatorProxyInterface interface {
	AccessController(opts *bind.CallOpts) (common.Address, error)

	Aggregator(opts *bind.CallOpts) (common.Address, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error)

	GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PhaseAggregators(opts *bind.CallOpts, arg0 uint16) (common.Address, error)

	PhaseId(opts *bind.CallOpts) (uint16, error)

	ProposedAggregator(opts *bind.CallOpts) (common.Address, error)

	ProposedGetRoundData(opts *bind.CallOpts, _roundId *big.Int) (ProposedGetRoundData,

		error)

	ProposedLatestRoundData(opts *bind.CallOpts) (ProposedLatestRoundData,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ConfirmAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error)

	ProposeAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error)

	SetController(opts *bind.TransactOpts, _accessController common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorProxyAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorProxyAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*AggregatorProxyAnswerUpdated, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorProxyNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorProxyNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*AggregatorProxyNewRound, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AggregatorProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AggregatorProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AggregatorProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AggregatorProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AggregatorProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AggregatorProxyOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
