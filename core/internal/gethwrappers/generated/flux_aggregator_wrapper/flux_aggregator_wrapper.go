// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package flux_aggregator_wrapper

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FluxAggregatorABI is the input ABI used to generate the binding from.
const FluxAggregatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_description\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AvailableFundsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"OracleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdateRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"OracleRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransfered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"RequesterAuthorizationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint128\",\"name\":\"paymentAmount\",\"type\":\"uint128\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"minAnswerCount\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"maxAnswerCount\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"restartDelay\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timeout\",\"type\":\"uint32\"}],\"name\":\"RoundDetailsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"round\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"SubmissionReceived\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"acceptAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"addOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_admins\",\"type\":\"address[]\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"addOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allocatedFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getOriginatingRoundOfAnswer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getRoundStartedAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimedOutStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"latestSubmission\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswerCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswerCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paymentAmount\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"removeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reportingRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reportingRoundStartedAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"restartDelay\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"roundState\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"_reportableRoundId\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"_eligibleToSubmit\",\"type\":\"bool\"},{\"internalType\":\"int256\",\"name\":\"_latestRoundAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint64\",\"name\":\"_timesOutAt\",\"type\":\"uint64\"},{\"internalType\":\"uint128\",\"name\":\"_availableFunds\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_allowed\",\"type\":\"bool\"}],\"name\":\"setAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startNewRound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeout\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"updateAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateAvailableFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"_newPaymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"}],\"name\":\"updateFutureRounds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"withdrawablePayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FluxAggregatorBin is the compiled bytecode used for deploying new contracts.
var FluxAggregatorBin = "0x60806040523480156200001157600080fd5b5060405162005ab838038062005ab8833981810160405260a08110156200003757600080fd5b810190808051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555084600660086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600360006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550826003601c6101000a81548163ffffffff021916908363ffffffff16021790555081600460006101000a81548160ff021916908360ff160217905550806005819055506200018d8363ffffffff1642620001da60201b620037b81790919060201c565b600860008063ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050505062000264565b60008282111562000253576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b61584480620002746000396000f3fe608060405234801561001057600080fd5b506004361061025e5760003560e01c80638da5cb5b11610146578063c410579e116100c3578063e2e4031711610087578063e2e4031714610d1d578063e6330cf714610d75578063e9ee6eeb14610dad578063eecea00014610e11578063f2fde38b14610e61578063ffa1ad7414610ea55761025e565b8063c410579e14610b82578063ca04f8f014610c69578063d002988c14610c87578063d4cc54e414610cb1578063e052cb0414610cf35761025e565b8063bb07bacd1161010a578063bb07bacd1461098b578063bbf0b7e9146109ea578063bd85948c14610ae8578063c107532914610af2578063c35905c614610b405761025e565b80638da5cb5b14610764578063a4c0ed36146107ae578063a4ce9a2714610893578063b5ab58dc14610907578063b633620c146109495761025e565b806350d25bcd116101df5780636fb4bb4e116101a35780636fb4bb4e1461069457806370dea79a146106b25780637284e416146106dc57806379b38bbb146106fa57806379ba50971461073c5780638205bf6a146107465761025e565b806350d25bcd14610566578063613d8fcc14610584578063628806ef146105ae57806364efb22b146105f2578063668a0f02146106765761025e565b806338aa4c721161022657806338aa4c72146103cd5780633d3d77141461044d57806340884c52146104bb57806346fcff4c1461051a5780634f8fc3b51461055c5761025e565b806309e24ae01461026357806325b6ae00146102a55780632f2f4767146102eb578063313ce5671461037f578063357ebb02146103a3575b600080fd5b61028f6004803603602081101561027957600080fd5b8101908080359060200190929190505050610ec3565b6040518082815260200191505060405180910390f35b6102d1600480360360208110156102bb57600080fd5b8101908080359060200190929190505050610f05565b604051808215151515815260200191505060405180910390f35b61037d600480360360a081101561030157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190505050610f6d565b005b610387610fbb565b604051808260ff1660ff16815260200191505060405180910390f35b6103ab610fce565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b61044b600480360360a08110156103e357600080fd5b8101908080356fffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190505050610fe4565b005b6104b96004803603606081101561046357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061126e565b005b6104c36115bc565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156105065780820151818401526020810190506104eb565b505050509050019250505060405180910390f35b61052261164a565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61056461166c565b005b61056e611833565b6040518082815260200191505060405180910390f35b61058c611842565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b6105f0600480360360208110156105c457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061184f565b005b6106346004803603602081101561060857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611a49565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61067e611ab5565b6040518082815260200191505060405180910390f35b61069c611ad5565b6040518082815260200191505060405180910390f35b6106ba611af5565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b6106e4611b0b565b6040518082815260200191505060405180910390f35b6107266004803603602081101561071057600080fd5b8101908080359060200190929190505050611b11565b6040518082815260200191505060405180910390f35b610744611b5b565b005b61074e611d23565b6040518082815260200191505060405180910390f35b61076c611d32565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610891600480360360608110156107c457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561080b57600080fd5b82018360208201111561081d57600080fd5b8035906020019184600183028401116401000000008311171561083f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050611d57565b005b610905600480360360808110156108a957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190505050611d64565b005b6109336004803603602081101561091d57600080fd5b810190808035906020019092919050505061218e565b6040518082815260200191505060405180910390f35b6109756004803603602081101561095f57600080fd5b81019080803590602001909291905050506121a0565b6040518082815260200191505060405180910390f35b6109cd600480360360208110156109a157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506121b2565b604051808381526020018281526020019250505060405180910390f35b610ae6600480360360a0811015610a0057600080fd5b8101908080359060200190640100000000811115610a1d57600080fd5b820183602082011115610a2f57600080fd5b80359060200191846020830284011164010000000083111715610a5157600080fd5b909192939192939080359060200190640100000000811115610a7257600080fd5b820183602082011115610a8457600080fd5b80359060200191846020830284011164010000000083111715610aa657600080fd5b9091929391929390803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff16906020019092919050505061225d565b005b610af061237a565b005b610b3e60048036036040811015610b0857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050612469565b005b610b48612661565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610bc460048036036020811015610b9857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612683565b604051808763ffffffff1663ffffffff168152602001861515151581526020018581526020018467ffffffffffffffff1667ffffffffffffffff168152602001836fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001965050505050505060405180910390f35b610c716128ec565b6040518082815260200191505060405180910390f35b610c8f612946565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610cb961295c565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610cfb61297e565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610d5f60048036036020811015610d3357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612994565b6040518082815260200191505060405180910390f35b610dab60048036036040811015610d8b57600080fd5b810190808035906020019092919080359060200190929190505050612a0e565b005b610e0f60048036036040811015610dc357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612c6e565b005b610e5f60048036036040811015610e2757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190505050612e3b565b005b610ea360048036036020811015610e7757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050613008565b005b610ead613189565b6040518082815260200191505060405180910390f35b6000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff1663ffffffff169050919050565b6000808290506000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff16905060008163ffffffff16118015610f6457508163ffffffff168163ffffffff1614155b92505050919050565b610f77858561318e565b610fb4600360009054906101000a90046fffffffffffffffffffffffffffffffff168484846003601c9054906101000a900463ffffffff16610fe4565b5050505050565b600460009054906101000a900460ff1681565b600360189054906101000a900463ffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146110a6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b83838360006110b3611842565b90508263ffffffff168163ffffffff1610156110ce57600080fd5b8363ffffffff168363ffffffff1610156110e757600080fd5b60008163ffffffff16148061110757508163ffffffff168163ffffffff16115b61111057600080fd5b88600360006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555087600360146101000a81548163ffffffff021916908363ffffffff16021790555086600360106101000a81548163ffffffff021916908363ffffffff16021790555085600360186101000a81548163ffffffff021916908363ffffffff160217905550846003601c6101000a81548163ffffffff021916908363ffffffff1602179055508663ffffffff168863ffffffff16600360009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff167f56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f8989604051808363ffffffff1663ffffffff1681526020018263ffffffff1663ffffffff1681526020019250505060405180910390a4505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff16600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461130857600080fd5b60008190506000600760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff169050816fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff1610156113a157600080fd5b6113c682826fffffffffffffffffffffffffffffffff1661370b90919063ffffffff16565b600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555061148182600260009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1661370b90919063ffffffff16565b600260006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550600660089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85846fffffffffffffffffffffffffffffffff166040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561157457600080fd5b505af1158015611588573d6000803e3d6000fd5b505050506040513d602081101561159e57600080fd5b81019080805190602001909291905050506115b557fe5b5050505050565b6060600a80548060200260200160405190810160405280929190818152602001828054801561164057602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116115f6575b5050505050905090565b600260109054906101000a90046fffffffffffffffffffffffffffffffff1681565b6000600260109054906101000a90046fffffffffffffffffffffffffffffffff16905060006117ad600260009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16600660089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561176457600080fd5b505afa158015611778573d6000803e3d6000fd5b505050506040513d602081101561178e57600080fd5b81019080805190602001909291905050506137b890919063ffffffff16565b905080600260106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555080826fffffffffffffffffffffffffffffffff161461182f57807ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f60405160405180910390a25b5050565b600061183d613841565b905090565b6000600a80549050905090565b3373ffffffffffffffffffffffffffffffffffffffff16600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146118e957600080fd5b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a350565b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6000600660049054906101000a900463ffffffff1663ffffffff16905090565b6000600660009054906101000a900463ffffffff1663ffffffff16905090565b6003601c9054906101000a900463ffffffff1681565b60055481565b6000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff169050919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611c1e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4d7573742062652070726f706f736564206f776e65720000000000000000000081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f0d18b5fd22306e373229b9439188228edca81207d1667f604daf6cef8aa3ee6760405160405180910390a350565b6000611d2d61387d565b905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b611d5f61166c565b505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611e26576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b8363ffffffff8016600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff1614611e9157600080fd5b600660009054906101000a900463ffffffff16600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff1602179055506000600a611f296001611f15611842565b63ffffffff166138d790919063ffffffff16565b63ffffffff1681548110611f3957fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506000600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900461ffff16905080600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81548161ffff021916908361ffff160217905550600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81549061ffff021916905581600a8261ffff168154811061207d57fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600a8054806120d057fe5b6001900381819060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905590558673ffffffffffffffffffffffffffffffffffffffff167f9c8e7d83025bef8a04c664b2f753f64b8814bdb7e27291d7e50935f18cc3c71260405160405180910390a2612185600360009054906101000a90046fffffffffffffffffffffffffffffffff168787876003601c9054906101000a900463ffffffff16610fe4565b50505050505050565b60006121998261396c565b9050919050565b60006121ab82613998565b9050919050565b600080600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff168063ffffffff16905091509150915091565b8484905087879050146122bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001806157e3602c913960400191505060405180910390fd5b60008090505b85859050811015612333576123268888838181106122db57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1687878481811061230457fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1661318e565b80806001019150506122c1565b50612371600360009054906101000a90046fffffffffffffffffffffffffffffffff168484846003601c9054906101000a900463ffffffff16610fe4565b50505050505050565b600960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166123d057600080fd5b6000600660009054906101000a900463ffffffff1690506000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16118061243b575061243a816139e2565b5b61244457600080fd5b61246661246160018363ffffffff16613ab890919063ffffffff16565b613b4c565b50565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461252b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600260109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16101561256857600080fd5b600660089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561261157600080fd5b505af1158015612625573d6000803e3d6000fd5b505050506040513d602081101561263b57600080fd5b810190808051906020019092919050505061265557600080fd5b61265d61166c565b5050565b600360009054906101000a90046fffffffffffffffffffffffffffffffff1681565b600080600080600080600060086000600660009054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff1660086000600660009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020019081526020016000206002016000018054905010158061273b575061273a600660009054906101000a900463ffffffff166139e2565b5b90508061275a57600660009054906101000a900463ffffffff16612787565b6127866001600660009054906101000a900463ffffffff1663ffffffff16613ab890919063ffffffff16565b5b965086612795898984613f23565b60086000600660049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020600001548361284857600860008b63ffffffff1663ffffffff16815260200190815260200160002060020160010160089054906101000a900463ffffffff1663ffffffff16600860008c63ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff160161284b565b60005b600260109054906101000a90046fffffffffffffffffffffffffffffffff16856128b657600860008d63ffffffff1663ffffffff168152602001908152602001600020600201600101600c9054906101000a90046fffffffffffffffffffffffffffffffff166128d6565b600360009054906101000a90046fffffffffffffffffffffffffffffffff165b9650965096509650965096505091939550919395565b600060086000600660009054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff16905090565b600360149054906101000a900463ffffffff1681565b600260009054906101000a90046fffffffffffffffffffffffffffffffff1681565b600360109054906101000a900463ffffffff1681565b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169050919050565b81600660009054906101000a900463ffffffff1663ffffffff168163ffffffff161480612a715750612a626001600660009054906101000a900463ffffffff1663ffffffff16613ab890919063ffffffff16565b63ffffffff168163ffffffff16145b612a7a57600080fd5b60018163ffffffff161480612aad5750612aac612aa760018363ffffffff166138d790919063ffffffff16565b6141c2565b5b80612ad65750612ad5612ad060018363ffffffff166138d790919063ffffffff16565b6139e2565b5b612adf57600080fd5b826000600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160109054906101000a900463ffffffff16905060008163ffffffff161415612b4b57600080fd5b8163ffffffff168163ffffffff161115612b6457600080fd5b8163ffffffff16600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff161015612bcf57600080fd5b8163ffffffff16600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff1663ffffffff1610612c3957600080fd5b612c4285613b4c565b612c4c848661420e565b612c558561439c565b612c5e8561459e565b612c678561481e565b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff16600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614612d0857600080fd5b80600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff167fb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b81043383604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a25050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612efd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b801515600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615151415612f5a57613004565b80600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f270d0c10dfbbdb6bb7206de0d1854b34e71664636d27af06feda4326a8d2437982604051808215151515815260200191505060405180910390a25b5050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146130ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b600281565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613250576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b8163ffffffff8016600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff1614156132bc57600080fd5b602a6132c6611842565b63ffffffff16106132d657600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561331057600080fd5b600073ffffffffffffffffffffffffffffffffffffffff16600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16148061343a57508173ffffffffffffffffffffffffffffffffffffffff16600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b61344357600080fd5b61344c8361492c565b600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160106101000a81548163ffffffff021916908363ffffffff16021790555063ffffffff600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff160217905550600a839080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061358d6001600a805490506137b890919063ffffffff16565b600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81548161ffff021916908361ffff16021790555081600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff167e47706786c922d17b39285dc59d696bafea72c0b003d3841ae1202076f4c2e460405160405180910390a28173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a3505050565b6000826fffffffffffffffffffffffffffffffff16826fffffffffffffffffffffffffffffffff1611156137a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b600082821115613830576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b600060086000600660049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060000154905090565b600060086000600660049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16905090565b60008263ffffffff168263ffffffff16111561395b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b6000600860008363ffffffff1663ffffffff168152602001908152602001600020600001549050919050565b6000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff169050919050565b600080600860008463ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1690506000600860008563ffffffff1663ffffffff16815260200190815260200160002060020160010160089054906101000a900463ffffffff16905060008267ffffffffffffffff16118015613a78575060008163ffffffff16115b8015613aaf575042613aa38263ffffffff168467ffffffffffffffff166149e990919063ffffffff16565b67ffffffffffffffff16105b92505050919050565b60008082840190508363ffffffff168163ffffffff161015613b42576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b80613b796001600660009054906101000a900463ffffffff1663ffffffff16613ab890919063ffffffff16565b63ffffffff168163ffffffff161415613f1f57816000600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff169050600360189054906101000a900463ffffffff1663ffffffff1681018263ffffffff161180613c185750600081145b15613f1c57613c3f613c3a60018663ffffffff166138d790919063ffffffff16565b614a85565b83600660006101000a81548163ffffffff021916908363ffffffff160217905550600360109054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160006101000a81548163ffffffff021916908363ffffffff160217905550600360149054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160046101000a81548163ffffffff021916908363ffffffff160217905550600360009054906101000a90046fffffffffffffffffffffffffffffffff16600860008663ffffffff1663ffffffff168152602001908152602001600020600201600101600c6101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506003601c9054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160086101000a81548163ffffffff021916908363ffffffff16021790555042600860008663ffffffff1663ffffffff16815260200190815260200160002060010160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c6101000a81548163ffffffff021916908363ffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168463ffffffff167f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271600860008863ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff16604051808267ffffffffffffffff16815260200191505060405180910390a35b50505b5050565b600080600760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160109054906101000a900463ffffffff16905060008163ffffffff161415613f945760009150506141bb565b8363ffffffff168163ffffffff161115613fb25760009150506141bb565b8363ffffffff16600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff1610156140225760009150506141bb565b8363ffffffff16600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff1663ffffffff16106140915760009150506141bb565b8215614166576000600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff169050600360189054906101000a900463ffffffff16810163ffffffff168563ffffffff1611158015614124575060008163ffffffff16115b15614134576000925050506141bb565b6000600360109054906101000a900463ffffffff1663ffffffff161415614160576000925050506141bb565b506141b5565b6000600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff1614156141b45760009150506141bb565b5b60019150505b9392505050565b600080600860008463ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16119050919050565b806000600860008363ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff16141561425857600080fd5b600860008363ffffffff1663ffffffff16815260200190815260200160002060020160000183908060018154018082558091505060019003906000526020600020016000909190919091505581600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160186101000a81548163ffffffff021916908363ffffffff16021790555082600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101819055503373ffffffffffffffffffffffffffffffffffffffff168263ffffffff16847f92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c60405160405180910390a4505050565b80600860008263ffffffff1663ffffffff16815260200190815260200160002060020160010160049054906101000a900463ffffffff1663ffffffff16600860008363ffffffff1663ffffffff168152602001908152602001600020600201600001805490501061459a576000614485600860008563ffffffff1663ffffffff16815260200190815260200160002060020160000180548060200260200160405190810160405280929190818152602001828054801561447b57602002820191906000526020600020905b815481526020019060010190808311614467575b5050505050614cc3565b905080600860008563ffffffff1663ffffffff1681526020019081526020016000206000018190555042600860008563ffffffff1663ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555082600860008563ffffffff1663ffffffff16815260200190815260200160002060010160106101000a81548163ffffffff021916908363ffffffff16021790555082600660046101000a81548163ffffffff021916908363ffffffff1602179055508263ffffffff16817f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a3505b5050565b6000600860008363ffffffff1663ffffffff168152602001908152602001600020600201600101600c9054906101000a90046fffffffffffffffffffffffffffffffff169050600061462982600260109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1661370b90919063ffffffff16565b905080600260106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506146a782600260009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16614db290919063ffffffff16565b600260006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555061476282600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16614db290919063ffffffff16565b600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550806fffffffffffffffffffffffffffffffff167ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f60405160405180910390a2505050565b80600860008263ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff16600860008363ffffffff1663ffffffff16815260200190815260200160002060020160000180549050141561492857600860008363ffffffff1663ffffffff168152602001908152602001600020600201600080820160006148bd919061577b565b6001820160006101000a81549063ffffffff02191690556001820160046101000a81549063ffffffff02191690556001820160086101000a81549063ffffffff021916905560018201600c6101000a8154906fffffffffffffffffffffffffffffffff021916905550505b5050565b600080600660009054906101000a900463ffffffff16905060008163ffffffff16141580156149b85750600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff168163ffffffff16145b156149c657809150506149e4565b6149e060018263ffffffff16613ab890919063ffffffff16565b9150505b919050565b60008082840190508367ffffffffffffffff168167ffffffffffffffff161015614a7b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b80614a8f816139e2565b15614cbf5781600060086000614ab560018563ffffffff166138d790919063ffffffff16565b63ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff161415614afc57600080fd5b6000614b1860018563ffffffff166138d790919063ffffffff16565b9050600860008263ffffffff1663ffffffff16815260200190815260200160002060000154600860008663ffffffff1663ffffffff16815260200190815260200160002060000181905550600860008263ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060010160106101000a81548163ffffffff021916908363ffffffff16021790555042600860008663ffffffff1663ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550600860008563ffffffff1663ffffffff16815260200190815260200160002060020160008082016000614c52919061577b565b6001820160006101000a81549063ffffffff02191690556001820160046101000a81549063ffffffff02191690556001820160086101000a81549063ffffffff021916905560018201600c6101000a8154906fffffffffffffffffffffffffffffffff0219169055505050505b5050565b60008151600010614d3c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6c697374206d757374206e6f7420626520656d7074790000000000000000000081525060200191505060405180910390fd5b600082519050600060028281614d4e57fe5b049050600060028381614d5d57fe5b061415614d9857600080614d7b866000600187036001870387614e5e565b8092508193505050614d8d8282614f4b565b945050505050614dad565b614da88460006001850384614fe8565b925050505b919050565b6000808284019050836fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff161015614e54576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b600080828410614e6d57600080fd5b838611158015614e7d5750848411155b614e8657600080fd5b828611158015614e965750848311155b614e9f57600080fd5b5b600115614f405760078686031015614ec857614ebf8787878787615082565b91509150614f41565b6000614ed58888886155f6565b9050808411614ee657809550614f3a565b84811015614ef957600181019650614f39565b808511158015614f0857508381105b614f0e57fe5b614f1a88888388614fe8565b9250614f2b88600183018887614fe8565b915082829250925050614f41565b5b50614ea0565b5b9550959350505050565b60008083128015614f5c5750600082135b80614f735750600083138015614f725750600082125b5b15614f93576002614f8484846156ed565b81614f8b57fe5b059050614fe2565b60006002808481614fa057fe5b0760028681614fab57fe5b070181614fb457fe5b059050614fde614fd860028681614fc757fe5b0560028681614fd257fe5b056156ed565b826156ed565b9150505b92915050565b600081841115614ff757600080fd5b8282111561500457600080fd5b5b8284101561506357600784840310156150385760006150278686868687615082565b80925081935050508191505061507a565b60006150458686866155f6565b90508083116150565780935061505d565b6001810194505b50615005565b84848151811061506f57fe5b602002602001015190505b949350505050565b6000806000866001870103905060008860008901815181106150a057fe5b602002602001015190506000826001106150da577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6150f2565b8960018a01815181106150e957fe5b60200260200101515b9050600083600210615124577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61513c565b8a60028b018151811061513357fe5b60200260200101515b905060008460031061516e577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff615186565b8b60038c018151811061517d57fe5b60200260200101515b90506000856004106151b8577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6151d0565b8c60048d01815181106151c757fe5b60200260200101515b9050600086600510615202577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61521a565b8d60058e018151811061521157fe5b60200260200101515b905060008760061061524c577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff615264565b8e60068f018151811061525b57fe5b60200260200101515b90508587131561527957858780975081985050505b8385131561528c57838580955081965050505b8183131561529f57818380935081945050505b848713156152b257848780965081985050505b838613156152c557838680955081975050505b808313156152d857808380925081945050505b848613156152eb57848680965081975050505b808213156152fe57808280925081935050505b8287131561531157828780945081985050505b8186131561532457818680935081975050505b8085131561533757808580925081965050505b8286131561534a57828680945081975050505b8084131561535d57808480925081955050505b8285131561537057828580945081965050505b8184131561538357818480935081955050505b8284131561539657828480945081955050505b60008e8d03905060008114156153ae57879a50615488565b60018114156153bf57869a50615487565b60028114156153d057859a50615486565b60038114156153e157849a50615485565b60048114156153f257839a50615484565b600581141561540357829a50615483565b600681141561541457819a50615482565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f6b31206f7574206f6620626f756e64730000000000000000000000000000000081525060200191505060405180910390fd5b5b5b5b5b5b5b60008f8d0390508c8e14156154ac578b8c9b509b50505050505050505050506155ec565b60008114156154ca578b899b509b50505050505050505050506155ec565b60018114156154e8578b889b509b50505050505050505050506155ec565b6002811415615506578b879b509b50505050505050505050506155ec565b6003811415615524578b869b509b50505050505050505050506155ec565b6004811415615542578b859b509b50505050505050505050506155ec565b6005811415615560578b849b509b50505050505050505050506155ec565b600681141561557e578b839b509b50505050505050505050506155ec565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f6b32206f7574206f6620626f756e64730000000000000000000000000000000081525060200191505060405180910390fd5b9550959350505050565b6000808460028486018161560657fe5b048151811061561157fe5b602002602001015190506001840393506001830192505b6001156156e4575b6001840193508085858151811061564357fe5b602002602001015112615630575b6001830392508085848151811061566457fe5b60200260200101511361565157828410156156d65784838151811061568557fe5b602002602001015185858151811061569957fe5b60200260200101518686815181106156ad57fe5b602002602001018786815181106156c057fe5b60200260200101828152508281525050506156df565b829150506156e6565b615628565b505b9392505050565b6000808284019050600083121580156157065750838112155b8061571c575060008312801561571b57508381125b5b615771576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806157c26021913960400191505060405180910390fd5b8091505092915050565b5080546000825590600052602060002090810190615799919061579c565b50565b6157be91905b808211156157ba5760008160009055506001016157a2565b5090565b9056fe5369676e6564536166654d6174683a206164646974696f6e206f766572666c6f776d7573742062652065786163746c79206f6e652061646d696e206164647265737320706572206f7261636c65a2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"

// DeployFluxAggregator deploys a new Ethereum contract, binding an instance of FluxAggregator to it.
func DeployFluxAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _paymentAmount *big.Int, _timeout uint32, _decimals uint8, _description [32]byte) (common.Address, *types.Transaction, *FluxAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(FluxAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FluxAggregatorBin), backend, _link, _paymentAmount, _timeout, _decimals, _description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FluxAggregator{FluxAggregatorCaller: FluxAggregatorCaller{contract: contract}, FluxAggregatorTransactor: FluxAggregatorTransactor{contract: contract}, FluxAggregatorFilterer: FluxAggregatorFilterer{contract: contract}}, nil
}

// FluxAggregator is an auto generated Go binding around an Ethereum contract.
type FluxAggregator struct {
	FluxAggregatorCaller     // Read-only binding to the contract
	FluxAggregatorTransactor // Write-only binding to the contract
	FluxAggregatorFilterer   // Log filterer for contract events
}

// FluxAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type FluxAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FluxAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FluxAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FluxAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FluxAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FluxAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FluxAggregatorSession struct {
	Contract     *FluxAggregator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FluxAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FluxAggregatorCallerSession struct {
	Contract *FluxAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// FluxAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FluxAggregatorTransactorSession struct {
	Contract     *FluxAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// FluxAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type FluxAggregatorRaw struct {
	Contract *FluxAggregator // Generic contract binding to access the raw methods on
}

// FluxAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FluxAggregatorCallerRaw struct {
	Contract *FluxAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// FluxAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FluxAggregatorTransactorRaw struct {
	Contract *FluxAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFluxAggregator creates a new instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregator(address common.Address, backend bind.ContractBackend) (*FluxAggregator, error) {
	contract, err := bindFluxAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FluxAggregator{FluxAggregatorCaller: FluxAggregatorCaller{contract: contract}, FluxAggregatorTransactor: FluxAggregatorTransactor{contract: contract}, FluxAggregatorFilterer: FluxAggregatorFilterer{contract: contract}}, nil
}

// NewFluxAggregatorCaller creates a new read-only instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregatorCaller(address common.Address, caller bind.ContractCaller) (*FluxAggregatorCaller, error) {
	contract, err := bindFluxAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorCaller{contract: contract}, nil
}

// NewFluxAggregatorTransactor creates a new write-only instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*FluxAggregatorTransactor, error) {
	contract, err := bindFluxAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorTransactor{contract: contract}, nil
}

// NewFluxAggregatorFilterer creates a new log filterer instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*FluxAggregatorFilterer, error) {
	contract, err := bindFluxAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorFilterer{contract: contract}, nil
}

// bindFluxAggregator binds a generic wrapper to an already deployed contract.
func bindFluxAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FluxAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FluxAggregator *FluxAggregatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FluxAggregator.Contract.FluxAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FluxAggregator *FluxAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.Contract.FluxAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FluxAggregator *FluxAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FluxAggregator.Contract.FluxAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FluxAggregator *FluxAggregatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FluxAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FluxAggregator *FluxAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FluxAggregator *FluxAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FluxAggregator.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) VERSION(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) VERSION() (*big.Int, error) {
	return _FluxAggregator.Contract.VERSION(&_FluxAggregator.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) VERSION() (*big.Int, error) {
	return _FluxAggregator.Contract.VERSION(&_FluxAggregator.CallOpts)
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCaller) AllocatedFunds(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "allocatedFunds")
	return *ret0, err
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorSession) AllocatedFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AllocatedFunds(&_FluxAggregator.CallOpts)
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCallerSession) AllocatedFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AllocatedFunds(&_FluxAggregator.CallOpts)
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCaller) AvailableFunds(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "availableFunds")
	return *ret0, err
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorSession) AvailableFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AvailableFunds(&_FluxAggregator.CallOpts)
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCallerSession) AvailableFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AvailableFunds(&_FluxAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_FluxAggregator *FluxAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_FluxAggregator *FluxAggregatorSession) Decimals() (uint8, error) {
	return _FluxAggregator.Contract.Decimals(&_FluxAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_FluxAggregator *FluxAggregatorCallerSession) Decimals() (uint8, error) {
	return _FluxAggregator.Contract.Decimals(&_FluxAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(bytes32)
func (_FluxAggregator *FluxAggregatorCaller) Description(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "description")
	return *ret0, err
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(bytes32)
func (_FluxAggregator *FluxAggregatorSession) Description() ([32]byte, error) {
	return _FluxAggregator.Contract.Description(&_FluxAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(bytes32)
func (_FluxAggregator *FluxAggregatorCallerSession) Description() ([32]byte, error) {
	return _FluxAggregator.Contract.Description(&_FluxAggregator.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) constant returns(address)
func (_FluxAggregator *FluxAggregatorCaller) GetAdmin(opts *bind.CallOpts, _oracle common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getAdmin", _oracle)
	return *ret0, err
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) constant returns(address)
func (_FluxAggregator *FluxAggregatorSession) GetAdmin(_oracle common.Address) (common.Address, error) {
	return _FluxAggregator.Contract.GetAdmin(&_FluxAggregator.CallOpts, _oracle)
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) constant returns(address)
func (_FluxAggregator *FluxAggregatorCallerSession) GetAdmin(_oracle common.Address) (common.Address, error) {
	return _FluxAggregator.Contract.GetAdmin(&_FluxAggregator.CallOpts, _oracle)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) constant returns(int256)
func (_FluxAggregator *FluxAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getAnswer", _roundId)
	return *ret0, err
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) constant returns(int256)
func (_FluxAggregator *FluxAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetAnswer(&_FluxAggregator.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) constant returns(int256)
func (_FluxAggregator *FluxAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetAnswer(&_FluxAggregator.CallOpts, _roundId)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() constant returns(address[])
func (_FluxAggregator *FluxAggregatorCaller) GetOracles(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getOracles")
	return *ret0, err
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() constant returns(address[])
func (_FluxAggregator *FluxAggregatorSession) GetOracles() ([]common.Address, error) {
	return _FluxAggregator.Contract.GetOracles(&_FluxAggregator.CallOpts)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() constant returns(address[])
func (_FluxAggregator *FluxAggregatorCallerSession) GetOracles() ([]common.Address, error) {
	return _FluxAggregator.Contract.GetOracles(&_FluxAggregator.CallOpts)
}

// GetOriginatingRoundOfAnswer is a free data retrieval call binding the contract method 0x09e24ae0.
//
// Solidity: function getOriginatingRoundOfAnswer(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) GetOriginatingRoundOfAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getOriginatingRoundOfAnswer", _roundId)
	return *ret0, err
}

// GetOriginatingRoundOfAnswer is a free data retrieval call binding the contract method 0x09e24ae0.
//
// Solidity: function getOriginatingRoundOfAnswer(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) GetOriginatingRoundOfAnswer(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetOriginatingRoundOfAnswer(&_FluxAggregator.CallOpts, _roundId)
}

// GetOriginatingRoundOfAnswer is a free data retrieval call binding the contract method 0x09e24ae0.
//
// Solidity: function getOriginatingRoundOfAnswer(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) GetOriginatingRoundOfAnswer(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetOriginatingRoundOfAnswer(&_FluxAggregator.CallOpts, _roundId)
}

// GetRoundStartedAt is a free data retrieval call binding the contract method 0x79b38bbb.
//
// Solidity: function getRoundStartedAt(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) GetRoundStartedAt(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getRoundStartedAt", _roundId)
	return *ret0, err
}

// GetRoundStartedAt is a free data retrieval call binding the contract method 0x79b38bbb.
//
// Solidity: function getRoundStartedAt(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) GetRoundStartedAt(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetRoundStartedAt(&_FluxAggregator.CallOpts, _roundId)
}

// GetRoundStartedAt is a free data retrieval call binding the contract method 0x79b38bbb.
//
// Solidity: function getRoundStartedAt(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) GetRoundStartedAt(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetRoundStartedAt(&_FluxAggregator.CallOpts, _roundId)
}

// GetTimedOutStatus is a free data retrieval call binding the contract method 0x25b6ae00.
//
// Solidity: function getTimedOutStatus(uint256 _roundId) constant returns(bool)
func (_FluxAggregator *FluxAggregatorCaller) GetTimedOutStatus(opts *bind.CallOpts, _roundId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getTimedOutStatus", _roundId)
	return *ret0, err
}

// GetTimedOutStatus is a free data retrieval call binding the contract method 0x25b6ae00.
//
// Solidity: function getTimedOutStatus(uint256 _roundId) constant returns(bool)
func (_FluxAggregator *FluxAggregatorSession) GetTimedOutStatus(_roundId *big.Int) (bool, error) {
	return _FluxAggregator.Contract.GetTimedOutStatus(&_FluxAggregator.CallOpts, _roundId)
}

// GetTimedOutStatus is a free data retrieval call binding the contract method 0x25b6ae00.
//
// Solidity: function getTimedOutStatus(uint256 _roundId) constant returns(bool)
func (_FluxAggregator *FluxAggregatorCallerSession) GetTimedOutStatus(_roundId *big.Int) (bool, error) {
	return _FluxAggregator.Contract.GetTimedOutStatus(&_FluxAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getTimestamp", _roundId)
	return *ret0, err
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetTimestamp(&_FluxAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetTimestamp(&_FluxAggregator.CallOpts, _roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_FluxAggregator *FluxAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_FluxAggregator *FluxAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestAnswer(&_FluxAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestAnswer(&_FluxAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "latestRound")
	return *ret0, err
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) LatestRound() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestRound(&_FluxAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestRound(&_FluxAggregator.CallOpts)
}

// LatestSubmission is a free data retrieval call binding the contract method 0xbb07bacd.
//
// Solidity: function latestSubmission(address _oracle) constant returns(int256, uint256)
func (_FluxAggregator *FluxAggregatorCaller) LatestSubmission(opts *bind.CallOpts, _oracle common.Address) (*big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _FluxAggregator.contract.Call(opts, out, "latestSubmission", _oracle)
	return *ret0, *ret1, err
}

// LatestSubmission is a free data retrieval call binding the contract method 0xbb07bacd.
//
// Solidity: function latestSubmission(address _oracle) constant returns(int256, uint256)
func (_FluxAggregator *FluxAggregatorSession) LatestSubmission(_oracle common.Address) (*big.Int, *big.Int, error) {
	return _FluxAggregator.Contract.LatestSubmission(&_FluxAggregator.CallOpts, _oracle)
}

// LatestSubmission is a free data retrieval call binding the contract method 0xbb07bacd.
//
// Solidity: function latestSubmission(address _oracle) constant returns(int256, uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestSubmission(_oracle common.Address) (*big.Int, *big.Int, error) {
	return _FluxAggregator.Contract.LatestSubmission(&_FluxAggregator.CallOpts, _oracle)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestTimestamp(&_FluxAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestTimestamp(&_FluxAggregator.CallOpts)
}

// MaxAnswerCount is a free data retrieval call binding the contract method 0xe052cb04.
//
// Solidity: function maxAnswerCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) MaxAnswerCount(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "maxAnswerCount")
	return *ret0, err
}

// MaxAnswerCount is a free data retrieval call binding the contract method 0xe052cb04.
//
// Solidity: function maxAnswerCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) MaxAnswerCount() (uint32, error) {
	return _FluxAggregator.Contract.MaxAnswerCount(&_FluxAggregator.CallOpts)
}

// MaxAnswerCount is a free data retrieval call binding the contract method 0xe052cb04.
//
// Solidity: function maxAnswerCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) MaxAnswerCount() (uint32, error) {
	return _FluxAggregator.Contract.MaxAnswerCount(&_FluxAggregator.CallOpts)
}

// MinAnswerCount is a free data retrieval call binding the contract method 0xd002988c.
//
// Solidity: function minAnswerCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) MinAnswerCount(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "minAnswerCount")
	return *ret0, err
}

// MinAnswerCount is a free data retrieval call binding the contract method 0xd002988c.
//
// Solidity: function minAnswerCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) MinAnswerCount() (uint32, error) {
	return _FluxAggregator.Contract.MinAnswerCount(&_FluxAggregator.CallOpts)
}

// MinAnswerCount is a free data retrieval call binding the contract method 0xd002988c.
//
// Solidity: function minAnswerCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) MinAnswerCount() (uint32, error) {
	return _FluxAggregator.Contract.MinAnswerCount(&_FluxAggregator.CallOpts)
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) OracleCount(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "oracleCount")
	return *ret0, err
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) OracleCount() (uint32, error) {
	return _FluxAggregator.Contract.OracleCount(&_FluxAggregator.CallOpts)
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) OracleCount() (uint32, error) {
	return _FluxAggregator.Contract.OracleCount(&_FluxAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FluxAggregator *FluxAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FluxAggregator *FluxAggregatorSession) Owner() (common.Address, error) {
	return _FluxAggregator.Contract.Owner(&_FluxAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FluxAggregator *FluxAggregatorCallerSession) Owner() (common.Address, error) {
	return _FluxAggregator.Contract.Owner(&_FluxAggregator.CallOpts)
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCaller) PaymentAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "paymentAmount")
	return *ret0, err
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorSession) PaymentAmount() (*big.Int, error) {
	return _FluxAggregator.Contract.PaymentAmount(&_FluxAggregator.CallOpts)
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCallerSession) PaymentAmount() (*big.Int, error) {
	return _FluxAggregator.Contract.PaymentAmount(&_FluxAggregator.CallOpts)
}

// ReportingRound is a free data retrieval call binding the contract method 0x6fb4bb4e.
//
// Solidity: function reportingRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) ReportingRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "reportingRound")
	return *ret0, err
}

// ReportingRound is a free data retrieval call binding the contract method 0x6fb4bb4e.
//
// Solidity: function reportingRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) ReportingRound() (*big.Int, error) {
	return _FluxAggregator.Contract.ReportingRound(&_FluxAggregator.CallOpts)
}

// ReportingRound is a free data retrieval call binding the contract method 0x6fb4bb4e.
//
// Solidity: function reportingRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) ReportingRound() (*big.Int, error) {
	return _FluxAggregator.Contract.ReportingRound(&_FluxAggregator.CallOpts)
}

// ReportingRoundStartedAt is a free data retrieval call binding the contract method 0xca04f8f0.
//
// Solidity: function reportingRoundStartedAt() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) ReportingRoundStartedAt(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "reportingRoundStartedAt")
	return *ret0, err
}

// ReportingRoundStartedAt is a free data retrieval call binding the contract method 0xca04f8f0.
//
// Solidity: function reportingRoundStartedAt() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) ReportingRoundStartedAt() (*big.Int, error) {
	return _FluxAggregator.Contract.ReportingRoundStartedAt(&_FluxAggregator.CallOpts)
}

// ReportingRoundStartedAt is a free data retrieval call binding the contract method 0xca04f8f0.
//
// Solidity: function reportingRoundStartedAt() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) ReportingRoundStartedAt() (*big.Int, error) {
	return _FluxAggregator.Contract.ReportingRoundStartedAt(&_FluxAggregator.CallOpts)
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) RestartDelay(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "restartDelay")
	return *ret0, err
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) RestartDelay() (uint32, error) {
	return _FluxAggregator.Contract.RestartDelay(&_FluxAggregator.CallOpts)
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) RestartDelay() (uint32, error) {
	return _FluxAggregator.Contract.RestartDelay(&_FluxAggregator.CallOpts)
}

// RoundState is a free data retrieval call binding the contract method 0xc410579e.
//
// Solidity: function roundState(address _oracle) constant returns(uint32 _reportableRoundId, bool _eligibleToSubmit, int256 _latestRoundAnswer, uint64 _timesOutAt, uint128 _availableFunds, uint128 _paymentAmount)
func (_FluxAggregator *FluxAggregatorCaller) RoundState(opts *bind.CallOpts, _oracle common.Address) (struct {
	ReportableRoundId uint32
	EligibleToSubmit  bool
	LatestRoundAnswer *big.Int
	TimesOutAt        uint64
	AvailableFunds    *big.Int
	PaymentAmount     *big.Int
}, error) {
	ret := new(struct {
		ReportableRoundId uint32
		EligibleToSubmit  bool
		LatestRoundAnswer *big.Int
		TimesOutAt        uint64
		AvailableFunds    *big.Int
		PaymentAmount     *big.Int
	})
	out := ret
	err := _FluxAggregator.contract.Call(opts, out, "roundState", _oracle)
	return *ret, err
}

// RoundState is a free data retrieval call binding the contract method 0xc410579e.
//
// Solidity: function roundState(address _oracle) constant returns(uint32 _reportableRoundId, bool _eligibleToSubmit, int256 _latestRoundAnswer, uint64 _timesOutAt, uint128 _availableFunds, uint128 _paymentAmount)
func (_FluxAggregator *FluxAggregatorSession) RoundState(_oracle common.Address) (struct {
	ReportableRoundId uint32
	EligibleToSubmit  bool
	LatestRoundAnswer *big.Int
	TimesOutAt        uint64
	AvailableFunds    *big.Int
	PaymentAmount     *big.Int
}, error) {
	return _FluxAggregator.Contract.RoundState(&_FluxAggregator.CallOpts, _oracle)
}

// RoundState is a free data retrieval call binding the contract method 0xc410579e.
//
// Solidity: function roundState(address _oracle) constant returns(uint32 _reportableRoundId, bool _eligibleToSubmit, int256 _latestRoundAnswer, uint64 _timesOutAt, uint128 _availableFunds, uint128 _paymentAmount)
func (_FluxAggregator *FluxAggregatorCallerSession) RoundState(_oracle common.Address) (struct {
	ReportableRoundId uint32
	EligibleToSubmit  bool
	LatestRoundAnswer *big.Int
	TimesOutAt        uint64
	AvailableFunds    *big.Int
	PaymentAmount     *big.Int
}, error) {
	return _FluxAggregator.Contract.RoundState(&_FluxAggregator.CallOpts, _oracle)
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) Timeout(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "timeout")
	return *ret0, err
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) Timeout() (uint32, error) {
	return _FluxAggregator.Contract.Timeout(&_FluxAggregator.CallOpts)
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) Timeout() (uint32, error) {
	return _FluxAggregator.Contract.Timeout(&_FluxAggregator.CallOpts)
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) WithdrawablePayment(opts *bind.CallOpts, _oracle common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "withdrawablePayment", _oracle)
	return *ret0, err
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) WithdrawablePayment(_oracle common.Address) (*big.Int, error) {
	return _FluxAggregator.Contract.WithdrawablePayment(&_FluxAggregator.CallOpts, _oracle)
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) WithdrawablePayment(_oracle common.Address) (*big.Int, error) {
	return _FluxAggregator.Contract.WithdrawablePayment(&_FluxAggregator.CallOpts, _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_FluxAggregator *FluxAggregatorTransactor) AcceptAdmin(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "acceptAdmin", _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_FluxAggregator *FluxAggregatorSession) AcceptAdmin(_oracle common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptAdmin(&_FluxAggregator.TransactOpts, _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AcceptAdmin(_oracle common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptAdmin(&_FluxAggregator.TransactOpts, _oracle)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FluxAggregator *FluxAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FluxAggregator *FluxAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptOwnership(&_FluxAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptOwnership(&_FluxAggregator.TransactOpts)
}

// AddOracle is a paid mutator transaction binding the contract method 0x2f2f4767.
//
// Solidity: function addOracle(address _oracle, address _admin, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactor) AddOracle(opts *bind.TransactOpts, _oracle common.Address, _admin common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "addOracle", _oracle, _admin, _minAnswers, _maxAnswers, _restartDelay)
}

// AddOracle is a paid mutator transaction binding the contract method 0x2f2f4767.
//
// Solidity: function addOracle(address _oracle, address _admin, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorSession) AddOracle(_oracle common.Address, _admin common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AddOracle(&_FluxAggregator.TransactOpts, _oracle, _admin, _minAnswers, _maxAnswers, _restartDelay)
}

// AddOracle is a paid mutator transaction binding the contract method 0x2f2f4767.
//
// Solidity: function addOracle(address _oracle, address _admin, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AddOracle(_oracle common.Address, _admin common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AddOracle(&_FluxAggregator.TransactOpts, _oracle, _admin, _minAnswers, _maxAnswers, _restartDelay)
}

// AddOracles is a paid mutator transaction binding the contract method 0xbbf0b7e9.
//
// Solidity: function addOracles(address[] _oracles, address[] _admins, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactor) AddOracles(opts *bind.TransactOpts, _oracles []common.Address, _admins []common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "addOracles", _oracles, _admins, _minAnswers, _maxAnswers, _restartDelay)
}

// AddOracles is a paid mutator transaction binding the contract method 0xbbf0b7e9.
//
// Solidity: function addOracles(address[] _oracles, address[] _admins, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorSession) AddOracles(_oracles []common.Address, _admins []common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AddOracles(&_FluxAggregator.TransactOpts, _oracles, _admins, _minAnswers, _maxAnswers, _restartDelay)
}

// AddOracles is a paid mutator transaction binding the contract method 0xbbf0b7e9.
//
// Solidity: function addOracles(address[] _oracles, address[] _admins, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AddOracles(_oracles []common.Address, _admins []common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AddOracles(&_FluxAggregator.TransactOpts, _oracles, _admins, _minAnswers, _maxAnswers, _restartDelay)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes ) returns()
func (_FluxAggregator *FluxAggregatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "onTokenTransfer", arg0, arg1, arg2)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes ) returns()
func (_FluxAggregator *FluxAggregatorSession) OnTokenTransfer(arg0 common.Address, arg1 *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _FluxAggregator.Contract.OnTokenTransfer(&_FluxAggregator.TransactOpts, arg0, arg1, arg2)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes ) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) OnTokenTransfer(arg0 common.Address, arg1 *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _FluxAggregator.Contract.OnTokenTransfer(&_FluxAggregator.TransactOpts, arg0, arg1, arg2)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0xa4ce9a27.
//
// Solidity: function removeOracle(address _oracle, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactor) RemoveOracle(opts *bind.TransactOpts, _oracle common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "removeOracle", _oracle, _minAnswers, _maxAnswers, _restartDelay)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0xa4ce9a27.
//
// Solidity: function removeOracle(address _oracle, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorSession) RemoveOracle(_oracle common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.RemoveOracle(&_FluxAggregator.TransactOpts, _oracle, _minAnswers, _maxAnswers, _restartDelay)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0xa4ce9a27.
//
// Solidity: function removeOracle(address _oracle, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) RemoveOracle(_oracle common.Address, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.RemoveOracle(&_FluxAggregator.TransactOpts, _oracle, _minAnswers, _maxAnswers, _restartDelay)
}

// SetAuthorization is a paid mutator transaction binding the contract method 0xeecea000.
//
// Solidity: function setAuthorization(address _requester, bool _allowed) returns()
func (_FluxAggregator *FluxAggregatorTransactor) SetAuthorization(opts *bind.TransactOpts, _requester common.Address, _allowed bool) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "setAuthorization", _requester, _allowed)
}

// SetAuthorization is a paid mutator transaction binding the contract method 0xeecea000.
//
// Solidity: function setAuthorization(address _requester, bool _allowed) returns()
func (_FluxAggregator *FluxAggregatorSession) SetAuthorization(_requester common.Address, _allowed bool) (*types.Transaction, error) {
	return _FluxAggregator.Contract.SetAuthorization(&_FluxAggregator.TransactOpts, _requester, _allowed)
}

// SetAuthorization is a paid mutator transaction binding the contract method 0xeecea000.
//
// Solidity: function setAuthorization(address _requester, bool _allowed) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) SetAuthorization(_requester common.Address, _allowed bool) (*types.Transaction, error) {
	return _FluxAggregator.Contract.SetAuthorization(&_FluxAggregator.TransactOpts, _requester, _allowed)
}

// StartNewRound is a paid mutator transaction binding the contract method 0xbd85948c.
//
// Solidity: function startNewRound() returns()
func (_FluxAggregator *FluxAggregatorTransactor) StartNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "startNewRound")
}

// StartNewRound is a paid mutator transaction binding the contract method 0xbd85948c.
//
// Solidity: function startNewRound() returns()
func (_FluxAggregator *FluxAggregatorSession) StartNewRound() (*types.Transaction, error) {
	return _FluxAggregator.Contract.StartNewRound(&_FluxAggregator.TransactOpts)
}

// StartNewRound is a paid mutator transaction binding the contract method 0xbd85948c.
//
// Solidity: function startNewRound() returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) StartNewRound() (*types.Transaction, error) {
	return _FluxAggregator.Contract.StartNewRound(&_FluxAggregator.TransactOpts)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_FluxAggregator *FluxAggregatorTransactor) TransferAdmin(opts *bind.TransactOpts, _oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "transferAdmin", _oracle, _newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_FluxAggregator *FluxAggregatorSession) TransferAdmin(_oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferAdmin(&_FluxAggregator.TransactOpts, _oracle, _newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) TransferAdmin(_oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferAdmin(&_FluxAggregator.TransactOpts, _oracle, _newAdmin)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_FluxAggregator *FluxAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_FluxAggregator *FluxAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferOwnership(&_FluxAggregator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferOwnership(&_FluxAggregator.TransactOpts, _to)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xe6330cf7.
//
// Solidity: function updateAnswer(uint256 _round, int256 _answer) returns()
func (_FluxAggregator *FluxAggregatorTransactor) UpdateAnswer(opts *bind.TransactOpts, _round *big.Int, _answer *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "updateAnswer", _round, _answer)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xe6330cf7.
//
// Solidity: function updateAnswer(uint256 _round, int256 _answer) returns()
func (_FluxAggregator *FluxAggregatorSession) UpdateAnswer(_round *big.Int, _answer *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateAnswer(&_FluxAggregator.TransactOpts, _round, _answer)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xe6330cf7.
//
// Solidity: function updateAnswer(uint256 _round, int256 _answer) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) UpdateAnswer(_round *big.Int, _answer *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateAnswer(&_FluxAggregator.TransactOpts, _round, _answer)
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_FluxAggregator *FluxAggregatorTransactor) UpdateAvailableFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "updateAvailableFunds")
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_FluxAggregator *FluxAggregatorSession) UpdateAvailableFunds() (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateAvailableFunds(&_FluxAggregator.TransactOpts)
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) UpdateAvailableFunds() (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateAvailableFunds(&_FluxAggregator.TransactOpts)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _newPaymentAmount, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay, uint32 _timeout) returns()
func (_FluxAggregator *FluxAggregatorTransactor) UpdateFutureRounds(opts *bind.TransactOpts, _newPaymentAmount *big.Int, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "updateFutureRounds", _newPaymentAmount, _minAnswers, _maxAnswers, _restartDelay, _timeout)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _newPaymentAmount, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay, uint32 _timeout) returns()
func (_FluxAggregator *FluxAggregatorSession) UpdateFutureRounds(_newPaymentAmount *big.Int, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateFutureRounds(&_FluxAggregator.TransactOpts, _newPaymentAmount, _minAnswers, _maxAnswers, _restartDelay, _timeout)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _newPaymentAmount, uint32 _minAnswers, uint32 _maxAnswers, uint32 _restartDelay, uint32 _timeout) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) UpdateFutureRounds(_newPaymentAmount *big.Int, _minAnswers uint32, _maxAnswers uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateFutureRounds(&_FluxAggregator.TransactOpts, _newPaymentAmount, _minAnswers, _maxAnswers, _restartDelay, _timeout)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawFunds(&_FluxAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawFunds(&_FluxAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "withdrawPayment", _oracle, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorSession) WithdrawPayment(_oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawPayment(&_FluxAggregator.TransactOpts, _oracle, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) WithdrawPayment(_oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawPayment(&_FluxAggregator.TransactOpts, _oracle, _recipient, _amount)
}

// FluxAggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the FluxAggregator contract.
type FluxAggregatorAnswerUpdatedIterator struct {
	Event *FluxAggregatorAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorAnswerUpdated)
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
		it.Event = new(FluxAggregatorAnswerUpdated)
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
func (it *FluxAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorAnswerUpdated represents a AnswerUpdated event raised by the FluxAggregator contract.
type FluxAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_FluxAggregator *FluxAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*FluxAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorAnswerUpdatedIterator{contract: _FluxAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_FluxAggregator *FluxAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorAnswerUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_FluxAggregator *FluxAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*FluxAggregatorAnswerUpdated, error) {
	event := new(FluxAggregatorAnswerUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorAvailableFundsUpdatedIterator is returned from FilterAvailableFundsUpdated and is used to iterate over the raw logs and unpacked data for AvailableFundsUpdated events raised by the FluxAggregator contract.
type FluxAggregatorAvailableFundsUpdatedIterator struct {
	Event *FluxAggregatorAvailableFundsUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorAvailableFundsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorAvailableFundsUpdated)
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
		it.Event = new(FluxAggregatorAvailableFundsUpdated)
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
func (it *FluxAggregatorAvailableFundsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorAvailableFundsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorAvailableFundsUpdated represents a AvailableFundsUpdated event raised by the FluxAggregator contract.
type FluxAggregatorAvailableFundsUpdated struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAvailableFundsUpdated is a free log retrieval operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_FluxAggregator *FluxAggregatorFilterer) FilterAvailableFundsUpdated(opts *bind.FilterOpts, amount []*big.Int) (*FluxAggregatorAvailableFundsUpdatedIterator, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "AvailableFundsUpdated", amountRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorAvailableFundsUpdatedIterator{contract: _FluxAggregator.contract, event: "AvailableFundsUpdated", logs: logs, sub: sub}, nil
}

// WatchAvailableFundsUpdated is a free log subscription operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_FluxAggregator *FluxAggregatorFilterer) WatchAvailableFundsUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorAvailableFundsUpdated, amount []*big.Int) (event.Subscription, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "AvailableFundsUpdated", amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorAvailableFundsUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "AvailableFundsUpdated", log); err != nil {
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

// ParseAvailableFundsUpdated is a log parse operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_FluxAggregator *FluxAggregatorFilterer) ParseAvailableFundsUpdated(log types.Log) (*FluxAggregatorAvailableFundsUpdated, error) {
	event := new(FluxAggregatorAvailableFundsUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "AvailableFundsUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the FluxAggregator contract.
type FluxAggregatorNewRoundIterator struct {
	Event *FluxAggregatorNewRound // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorNewRound)
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
		it.Event = new(FluxAggregatorNewRound)
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
func (it *FluxAggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorNewRound represents a NewRound event raised by the FluxAggregator contract.
type FluxAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_FluxAggregator *FluxAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*FluxAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorNewRoundIterator{contract: _FluxAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_FluxAggregator *FluxAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *FluxAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorNewRound)
				if err := _FluxAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_FluxAggregator *FluxAggregatorFilterer) ParseNewRound(log types.Log) (*FluxAggregatorNewRound, error) {
	event := new(FluxAggregatorNewRound)
	if err := _FluxAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOracleAddedIterator is returned from FilterOracleAdded and is used to iterate over the raw logs and unpacked data for OracleAdded events raised by the FluxAggregator contract.
type FluxAggregatorOracleAddedIterator struct {
	Event *FluxAggregatorOracleAdded // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOracleAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOracleAdded)
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
		it.Event = new(FluxAggregatorOracleAdded)
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
func (it *FluxAggregatorOracleAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOracleAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOracleAdded represents a OracleAdded event raised by the FluxAggregator contract.
type FluxAggregatorOracleAdded struct {
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOracleAdded is a free log retrieval operation binding the contract event 0x0047706786c922d17b39285dc59d696bafea72c0b003d3841ae1202076f4c2e4.
//
// Solidity: event OracleAdded(address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOracleAdded(opts *bind.FilterOpts, oracle []common.Address) (*FluxAggregatorOracleAddedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OracleAdded", oracleRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOracleAddedIterator{contract: _FluxAggregator.contract, event: "OracleAdded", logs: logs, sub: sub}, nil
}

// WatchOracleAdded is a free log subscription operation binding the contract event 0x0047706786c922d17b39285dc59d696bafea72c0b003d3841ae1202076f4c2e4.
//
// Solidity: event OracleAdded(address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOracleAdded(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOracleAdded, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OracleAdded", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOracleAdded)
				if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdded", log); err != nil {
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

// ParseOracleAdded is a log parse operation binding the contract event 0x0047706786c922d17b39285dc59d696bafea72c0b003d3841ae1202076f4c2e4.
//
// Solidity: event OracleAdded(address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOracleAdded(log types.Log) (*FluxAggregatorOracleAdded, error) {
	event := new(FluxAggregatorOracleAdded)
	if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOracleAdminUpdateRequestedIterator is returned from FilterOracleAdminUpdateRequested and is used to iterate over the raw logs and unpacked data for OracleAdminUpdateRequested events raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdateRequestedIterator struct {
	Event *FluxAggregatorOracleAdminUpdateRequested // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOracleAdminUpdateRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOracleAdminUpdateRequested)
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
		it.Event = new(FluxAggregatorOracleAdminUpdateRequested)
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
func (it *FluxAggregatorOracleAdminUpdateRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOracleAdminUpdateRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOracleAdminUpdateRequested represents a OracleAdminUpdateRequested event raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdateRequested struct {
	Oracle   common.Address
	Admin    common.Address
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOracleAdminUpdateRequested is a free log retrieval operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOracleAdminUpdateRequested(opts *bind.FilterOpts, oracle []common.Address) (*FluxAggregatorOracleAdminUpdateRequestedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OracleAdminUpdateRequested", oracleRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOracleAdminUpdateRequestedIterator{contract: _FluxAggregator.contract, event: "OracleAdminUpdateRequested", logs: logs, sub: sub}, nil
}

// WatchOracleAdminUpdateRequested is a free log subscription operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOracleAdminUpdateRequested(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOracleAdminUpdateRequested, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OracleAdminUpdateRequested", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOracleAdminUpdateRequested)
				if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdateRequested", log); err != nil {
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

// ParseOracleAdminUpdateRequested is a log parse operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOracleAdminUpdateRequested(log types.Log) (*FluxAggregatorOracleAdminUpdateRequested, error) {
	event := new(FluxAggregatorOracleAdminUpdateRequested)
	if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdateRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOracleAdminUpdatedIterator is returned from FilterOracleAdminUpdated and is used to iterate over the raw logs and unpacked data for OracleAdminUpdated events raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdatedIterator struct {
	Event *FluxAggregatorOracleAdminUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOracleAdminUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOracleAdminUpdated)
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
		it.Event = new(FluxAggregatorOracleAdminUpdated)
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
func (it *FluxAggregatorOracleAdminUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOracleAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOracleAdminUpdated represents a OracleAdminUpdated event raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdated struct {
	Oracle   common.Address
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOracleAdminUpdated is a free log retrieval operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOracleAdminUpdated(opts *bind.FilterOpts, oracle []common.Address, newAdmin []common.Address) (*FluxAggregatorOracleAdminUpdatedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OracleAdminUpdated", oracleRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOracleAdminUpdatedIterator{contract: _FluxAggregator.contract, event: "OracleAdminUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleAdminUpdated is a free log subscription operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOracleAdminUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOracleAdminUpdated, oracle []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OracleAdminUpdated", oracleRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOracleAdminUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdated", log); err != nil {
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

// ParseOracleAdminUpdated is a log parse operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOracleAdminUpdated(log types.Log) (*FluxAggregatorOracleAdminUpdated, error) {
	event := new(FluxAggregatorOracleAdminUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOracleRemovedIterator is returned from FilterOracleRemoved and is used to iterate over the raw logs and unpacked data for OracleRemoved events raised by the FluxAggregator contract.
type FluxAggregatorOracleRemovedIterator struct {
	Event *FluxAggregatorOracleRemoved // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOracleRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOracleRemoved)
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
		it.Event = new(FluxAggregatorOracleRemoved)
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
func (it *FluxAggregatorOracleRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOracleRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOracleRemoved represents a OracleRemoved event raised by the FluxAggregator contract.
type FluxAggregatorOracleRemoved struct {
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOracleRemoved is a free log retrieval operation binding the contract event 0x9c8e7d83025bef8a04c664b2f753f64b8814bdb7e27291d7e50935f18cc3c712.
//
// Solidity: event OracleRemoved(address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOracleRemoved(opts *bind.FilterOpts, oracle []common.Address) (*FluxAggregatorOracleRemovedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OracleRemoved", oracleRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOracleRemovedIterator{contract: _FluxAggregator.contract, event: "OracleRemoved", logs: logs, sub: sub}, nil
}

// WatchOracleRemoved is a free log subscription operation binding the contract event 0x9c8e7d83025bef8a04c664b2f753f64b8814bdb7e27291d7e50935f18cc3c712.
//
// Solidity: event OracleRemoved(address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOracleRemoved(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOracleRemoved, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OracleRemoved", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOracleRemoved)
				if err := _FluxAggregator.contract.UnpackLog(event, "OracleRemoved", log); err != nil {
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

// ParseOracleRemoved is a log parse operation binding the contract event 0x9c8e7d83025bef8a04c664b2f753f64b8814bdb7e27291d7e50935f18cc3c712.
//
// Solidity: event OracleRemoved(address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOracleRemoved(log types.Log) (*FluxAggregatorOracleRemoved, error) {
	event := new(FluxAggregatorOracleRemoved)
	if err := _FluxAggregator.contract.UnpackLog(event, "OracleRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferRequestedIterator struct {
	Event *FluxAggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOwnershipTransferRequested)
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
		it.Event = new(FluxAggregatorOwnershipTransferRequested)
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
func (it *FluxAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FluxAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOwnershipTransferRequestedIterator{contract: _FluxAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOwnershipTransferRequested)
				if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*FluxAggregatorOwnershipTransferRequested, error) {
	event := new(FluxAggregatorOwnershipTransferRequested)
	if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOwnershipTransferedIterator is returned from FilterOwnershipTransfered and is used to iterate over the raw logs and unpacked data for OwnershipTransfered events raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferedIterator struct {
	Event *FluxAggregatorOwnershipTransfered // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOwnershipTransferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOwnershipTransfered)
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
		it.Event = new(FluxAggregatorOwnershipTransfered)
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
func (it *FluxAggregatorOwnershipTransferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOwnershipTransferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOwnershipTransfered represents a OwnershipTransfered event raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransfered struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransfered is a free log retrieval operation binding the contract event 0x0d18b5fd22306e373229b9439188228edca81207d1667f604daf6cef8aa3ee67.
//
// Solidity: event OwnershipTransfered(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOwnershipTransfered(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FluxAggregatorOwnershipTransferedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OwnershipTransfered", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOwnershipTransferedIterator{contract: _FluxAggregator.contract, event: "OwnershipTransfered", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransfered is a free log subscription operation binding the contract event 0x0d18b5fd22306e373229b9439188228edca81207d1667f604daf6cef8aa3ee67.
//
// Solidity: event OwnershipTransfered(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOwnershipTransfered(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOwnershipTransfered, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OwnershipTransfered", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOwnershipTransfered)
				if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransfered", log); err != nil {
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

// ParseOwnershipTransfered is a log parse operation binding the contract event 0x0d18b5fd22306e373229b9439188228edca81207d1667f604daf6cef8aa3ee67.
//
// Solidity: event OwnershipTransfered(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOwnershipTransfered(log types.Log) (*FluxAggregatorOwnershipTransfered, error) {
	event := new(FluxAggregatorOwnershipTransfered)
	if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransfered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorRequesterAuthorizationSetIterator is returned from FilterRequesterAuthorizationSet and is used to iterate over the raw logs and unpacked data for RequesterAuthorizationSet events raised by the FluxAggregator contract.
type FluxAggregatorRequesterAuthorizationSetIterator struct {
	Event *FluxAggregatorRequesterAuthorizationSet // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorRequesterAuthorizationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorRequesterAuthorizationSet)
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
		it.Event = new(FluxAggregatorRequesterAuthorizationSet)
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
func (it *FluxAggregatorRequesterAuthorizationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorRequesterAuthorizationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorRequesterAuthorizationSet represents a RequesterAuthorizationSet event raised by the FluxAggregator contract.
type FluxAggregatorRequesterAuthorizationSet struct {
	Requester common.Address
	Allowed   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequesterAuthorizationSet is a free log retrieval operation binding the contract event 0x270d0c10dfbbdb6bb7206de0d1854b34e71664636d27af06feda4326a8d24379.
//
// Solidity: event RequesterAuthorizationSet(address indexed requester, bool allowed)
func (_FluxAggregator *FluxAggregatorFilterer) FilterRequesterAuthorizationSet(opts *bind.FilterOpts, requester []common.Address) (*FluxAggregatorRequesterAuthorizationSetIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "RequesterAuthorizationSet", requesterRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorRequesterAuthorizationSetIterator{contract: _FluxAggregator.contract, event: "RequesterAuthorizationSet", logs: logs, sub: sub}, nil
}

// WatchRequesterAuthorizationSet is a free log subscription operation binding the contract event 0x270d0c10dfbbdb6bb7206de0d1854b34e71664636d27af06feda4326a8d24379.
//
// Solidity: event RequesterAuthorizationSet(address indexed requester, bool allowed)
func (_FluxAggregator *FluxAggregatorFilterer) WatchRequesterAuthorizationSet(opts *bind.WatchOpts, sink chan<- *FluxAggregatorRequesterAuthorizationSet, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "RequesterAuthorizationSet", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorRequesterAuthorizationSet)
				if err := _FluxAggregator.contract.UnpackLog(event, "RequesterAuthorizationSet", log); err != nil {
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

// ParseRequesterAuthorizationSet is a log parse operation binding the contract event 0x270d0c10dfbbdb6bb7206de0d1854b34e71664636d27af06feda4326a8d24379.
//
// Solidity: event RequesterAuthorizationSet(address indexed requester, bool allowed)
func (_FluxAggregator *FluxAggregatorFilterer) ParseRequesterAuthorizationSet(log types.Log) (*FluxAggregatorRequesterAuthorizationSet, error) {
	event := new(FluxAggregatorRequesterAuthorizationSet)
	if err := _FluxAggregator.contract.UnpackLog(event, "RequesterAuthorizationSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorRoundDetailsUpdatedIterator is returned from FilterRoundDetailsUpdated and is used to iterate over the raw logs and unpacked data for RoundDetailsUpdated events raised by the FluxAggregator contract.
type FluxAggregatorRoundDetailsUpdatedIterator struct {
	Event *FluxAggregatorRoundDetailsUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorRoundDetailsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorRoundDetailsUpdated)
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
		it.Event = new(FluxAggregatorRoundDetailsUpdated)
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
func (it *FluxAggregatorRoundDetailsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorRoundDetailsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorRoundDetailsUpdated represents a RoundDetailsUpdated event raised by the FluxAggregator contract.
type FluxAggregatorRoundDetailsUpdated struct {
	PaymentAmount  *big.Int
	MinAnswerCount uint32
	MaxAnswerCount uint32
	RestartDelay   uint32
	Timeout        uint32
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRoundDetailsUpdated is a free log retrieval operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minAnswerCount, uint32 indexed maxAnswerCount, uint32 restartDelay, uint32 timeout)
func (_FluxAggregator *FluxAggregatorFilterer) FilterRoundDetailsUpdated(opts *bind.FilterOpts, paymentAmount []*big.Int, minAnswerCount []uint32, maxAnswerCount []uint32) (*FluxAggregatorRoundDetailsUpdatedIterator, error) {

	var paymentAmountRule []interface{}
	for _, paymentAmountItem := range paymentAmount {
		paymentAmountRule = append(paymentAmountRule, paymentAmountItem)
	}
	var minAnswerCountRule []interface{}
	for _, minAnswerCountItem := range minAnswerCount {
		minAnswerCountRule = append(minAnswerCountRule, minAnswerCountItem)
	}
	var maxAnswerCountRule []interface{}
	for _, maxAnswerCountItem := range maxAnswerCount {
		maxAnswerCountRule = append(maxAnswerCountRule, maxAnswerCountItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "RoundDetailsUpdated", paymentAmountRule, minAnswerCountRule, maxAnswerCountRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorRoundDetailsUpdatedIterator{contract: _FluxAggregator.contract, event: "RoundDetailsUpdated", logs: logs, sub: sub}, nil
}

// WatchRoundDetailsUpdated is a free log subscription operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minAnswerCount, uint32 indexed maxAnswerCount, uint32 restartDelay, uint32 timeout)
func (_FluxAggregator *FluxAggregatorFilterer) WatchRoundDetailsUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorRoundDetailsUpdated, paymentAmount []*big.Int, minAnswerCount []uint32, maxAnswerCount []uint32) (event.Subscription, error) {

	var paymentAmountRule []interface{}
	for _, paymentAmountItem := range paymentAmount {
		paymentAmountRule = append(paymentAmountRule, paymentAmountItem)
	}
	var minAnswerCountRule []interface{}
	for _, minAnswerCountItem := range minAnswerCount {
		minAnswerCountRule = append(minAnswerCountRule, minAnswerCountItem)
	}
	var maxAnswerCountRule []interface{}
	for _, maxAnswerCountItem := range maxAnswerCount {
		maxAnswerCountRule = append(maxAnswerCountRule, maxAnswerCountItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "RoundDetailsUpdated", paymentAmountRule, minAnswerCountRule, maxAnswerCountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorRoundDetailsUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "RoundDetailsUpdated", log); err != nil {
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

// ParseRoundDetailsUpdated is a log parse operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minAnswerCount, uint32 indexed maxAnswerCount, uint32 restartDelay, uint32 timeout)
func (_FluxAggregator *FluxAggregatorFilterer) ParseRoundDetailsUpdated(log types.Log) (*FluxAggregatorRoundDetailsUpdated, error) {
	event := new(FluxAggregatorRoundDetailsUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "RoundDetailsUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorSubmissionReceivedIterator is returned from FilterSubmissionReceived and is used to iterate over the raw logs and unpacked data for SubmissionReceived events raised by the FluxAggregator contract.
type FluxAggregatorSubmissionReceivedIterator struct {
	Event *FluxAggregatorSubmissionReceived // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorSubmissionReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorSubmissionReceived)
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
		it.Event = new(FluxAggregatorSubmissionReceived)
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
func (it *FluxAggregatorSubmissionReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorSubmissionReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorSubmissionReceived represents a SubmissionReceived event raised by the FluxAggregator contract.
type FluxAggregatorSubmissionReceived struct {
	Answer *big.Int
	Round  uint32
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSubmissionReceived is a free log retrieval operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed answer, uint32 indexed round, address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) FilterSubmissionReceived(opts *bind.FilterOpts, answer []*big.Int, round []uint32, oracle []common.Address) (*FluxAggregatorSubmissionReceivedIterator, error) {

	var answerRule []interface{}
	for _, answerItem := range answer {
		answerRule = append(answerRule, answerItem)
	}
	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "SubmissionReceived", answerRule, roundRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorSubmissionReceivedIterator{contract: _FluxAggregator.contract, event: "SubmissionReceived", logs: logs, sub: sub}, nil
}

// WatchSubmissionReceived is a free log subscription operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed answer, uint32 indexed round, address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) WatchSubmissionReceived(opts *bind.WatchOpts, sink chan<- *FluxAggregatorSubmissionReceived, answer []*big.Int, round []uint32, oracle []common.Address) (event.Subscription, error) {

	var answerRule []interface{}
	for _, answerItem := range answer {
		answerRule = append(answerRule, answerItem)
	}
	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "SubmissionReceived", answerRule, roundRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorSubmissionReceived)
				if err := _FluxAggregator.contract.UnpackLog(event, "SubmissionReceived", log); err != nil {
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

// ParseSubmissionReceived is a log parse operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed answer, uint32 indexed round, address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) ParseSubmissionReceived(log types.Log) (*FluxAggregatorSubmissionReceived, error) {
	event := new(FluxAggregatorSubmissionReceived)
	if err := _FluxAggregator.contract.UnpackLog(event, "SubmissionReceived", log); err != nil {
		return nil, err
	}
	return event, nil
}
