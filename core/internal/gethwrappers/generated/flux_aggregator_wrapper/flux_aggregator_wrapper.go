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
const FluxAggregatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_description\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AvailableFundsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"OracleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdateRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"OracleRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransfered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"RequesterAuthorizationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint128\",\"name\":\"paymentAmount\",\"type\":\"uint128\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"minAnswerCount\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"maxAnswerCount\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"restartDelay\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timeout\",\"type\":\"uint32\"}],\"name\":\"RoundDetailsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"round\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"SubmissionReceived\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"acceptAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"addOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allocatedFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getOriginatingRoundOfAnswer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getRoundStartedAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimedOutStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"latestSubmission\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswerCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswerCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paymentAmount\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"removeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reportingRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reportingRoundStartedAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"restartDelay\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"roundState\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"_reportableRoundId\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"_eligibleToSubmit\",\"type\":\"bool\"},{\"internalType\":\"int256\",\"name\":\"_latestRoundAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint64\",\"name\":\"_timesOutAt\",\"type\":\"uint64\"},{\"internalType\":\"uint128\",\"name\":\"_availableFunds\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_allowed\",\"type\":\"bool\"}],\"name\":\"setAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startNewRound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeout\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"updateAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateAvailableFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"_newPaymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_minAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxAnswers\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"}],\"name\":\"updateFutureRounds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"withdrawablePayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FluxAggregatorBin is the compiled bytecode used for deploying new contracts.
var FluxAggregatorBin = "0x60806040523480156200001157600080fd5b506040516200585838038062005858833981810160405260a08110156200003757600080fd5b810190808051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555084600660086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600360006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550826003601c6101000a81548163ffffffff021916908363ffffffff16021790555081600460006101000a81548160ff021916908360ff160217905550806005819055506200018d8363ffffffff1642620001da60201b620035941790919060201c565b600860008063ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050505062000264565b60008282111562000253576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b6155e480620002746000396000f3fe608060405234801561001057600080fd5b50600436106102535760003560e01c80638205bf6a11610146578063c410579e116100c3578063e2e4031711610087578063e2e4031714610c14578063e6330cf714610c6c578063e9ee6eeb14610ca4578063eecea00014610d08578063f2fde38b14610d58578063ffa1ad7414610d9c57610253565b8063c410579e14610a79578063ca04f8f014610b60578063d002988c14610b7e578063d4cc54e414610ba8578063e052cb0414610bea57610253565b8063b633620c1161010a578063b633620c1461093e578063bb07bacd14610980578063bd85948c146109df578063c1075329146109e9578063c35905c614610a3757610253565b80638205bf6a1461073b5780638da5cb5b14610759578063a4c0ed36146107a3578063a4ce9a2714610888578063b5ab58dc146108fc57610253565b806350d25bcd116101d45780636fb4bb4e116101985780636fb4bb4e1461068957806370dea79a146106a75780637284e416146106d157806379b38bbb146106ef57806379ba50971461073157610253565b806350d25bcd1461055b578063613d8fcc14610579578063628806ef146105a357806364efb22b146105e7578063668a0f021461066b57610253565b806338aa4c721161021b57806338aa4c72146103c25780633d3d77141461044257806340884c52146104b057806346fcff4c1461050f5780634f8fc3b51461055157610253565b806309e24ae01461025857806325b6ae001461029a5780632f2f4767146102e0578063313ce56714610374578063357ebb0214610398575b600080fd5b6102846004803603602081101561026e57600080fd5b8101908080359060200190929190505050610dba565b6040518082815260200191505060405180910390f35b6102c6600480360360208110156102b057600080fd5b8101908080359060200190929190505050610dfc565b604051808215151515815260200191505060405180910390f35b610372600480360360a08110156102f657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190505050610e64565b005b61037c611421565b604051808260ff1660ff16815260200191505060405180910390f35b6103a0611434565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610440600480360360a08110156103d857600080fd5b8101908080356fffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff16906020019092919050505061144a565b005b6104ae6004803603606081101561045857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506116d4565b005b6104b8611a22565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156104fb5780820151818401526020810190506104e0565b505050509050019250505060405180910390f35b610517611ab0565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610559611ad2565b005b610563611c99565b6040518082815260200191505060405180910390f35b610581611ca8565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b6105e5600480360360208110156105b957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611cb5565b005b610629600480360360208110156105fd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611eaf565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610673611f1b565b6040518082815260200191505060405180910390f35b610691611f3b565b6040518082815260200191505060405180910390f35b6106af611f5b565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b6106d9611f71565b6040518082815260200191505060405180910390f35b61071b6004803603602081101561070557600080fd5b8101908080359060200190929190505050611f77565b6040518082815260200191505060405180910390f35b610739611fc1565b005b610743612189565b6040518082815260200191505060405180910390f35b610761612198565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610886600480360360608110156107b957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561080057600080fd5b82018360208201111561081257600080fd5b8035906020019184600183028401116401000000008311171561083457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506121bd565b005b6108fa6004803603608081101561089e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff1690602001909291905050506121ca565b005b6109286004803603602081101561091257600080fd5b81019080803590602001909291905050506125f4565b6040518082815260200191505060405180910390f35b61096a6004803603602081101561095457600080fd5b8101908080359060200190929190505050612606565b6040518082815260200191505060405180910390f35b6109c26004803603602081101561099657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612618565b604051808381526020018281526020019250505060405180910390f35b6109e76126c3565b005b610a35600480360360408110156109ff57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506127b2565b005b610a3f6129aa565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610abb60048036036020811015610a8f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506129cc565b604051808763ffffffff1663ffffffff168152602001861515151581526020018581526020018467ffffffffffffffff1667ffffffffffffffff168152602001836fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001965050505050505060405180910390f35b610b68612c35565b6040518082815260200191505060405180910390f35b610b86612c8f565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610bb0612ca5565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610bf2612cc7565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610c5660048036036020811015610c2a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612cdd565b6040518082815260200191505060405180910390f35b610ca260048036036040811015610c8257600080fd5b810190808035906020019092919080359060200190929190505050612d57565b005b610d0660048036036040811015610cba57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612fb7565b005b610d5660048036036040811015610d1e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190505050613184565b005b610d9a60048036036020811015610d6e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050613351565b005b610da46134d2565b6040518082815260200191505060405180910390f35b6000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff1663ffffffff169050919050565b6000808290506000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff16905060008163ffffffff16118015610e5b57508163ffffffff168163ffffffff1614155b92505050919050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610f26576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b8463ffffffff8016600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff161415610f9257600080fd5b602a610f9c611ca8565b63ffffffff1610610fac57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161415610fe657600080fd5b600073ffffffffffffffffffffffffffffffffffffffff16600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16148061111057508473ffffffffffffffffffffffffffffffffffffffff16600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b61111957600080fd5b611122866134d7565b600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160106101000a81548163ffffffff021916908363ffffffff16021790555063ffffffff600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff160217905550600a869080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506112636001600a8054905061359490919063ffffffff16565b600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81548161ffff021916908361ffff16021790555084600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508573ffffffffffffffffffffffffffffffffffffffff167e47706786c922d17b39285dc59d696bafea72c0b003d3841ae1202076f4c2e460405160405180910390a28473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a3611419600360009054906101000a90046fffffffffffffffffffffffffffffffff168585856003601c9054906101000a900463ffffffff1661144a565b505050505050565b600460009054906101000a900460ff1681565b600360189054906101000a900463ffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461150c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b8383836000611519611ca8565b90508263ffffffff168163ffffffff16101561153457600080fd5b8363ffffffff168363ffffffff16101561154d57600080fd5b60008163ffffffff16148061156d57508163ffffffff168163ffffffff16115b61157657600080fd5b88600360006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555087600360146101000a81548163ffffffff021916908363ffffffff16021790555086600360106101000a81548163ffffffff021916908363ffffffff16021790555085600360186101000a81548163ffffffff021916908363ffffffff160217905550846003601c6101000a81548163ffffffff021916908363ffffffff1602179055508663ffffffff168863ffffffff16600360009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff167f56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f8989604051808363ffffffff1663ffffffff1681526020018263ffffffff1663ffffffff1681526020019250505060405180910390a4505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff16600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461176e57600080fd5b60008190506000600760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff169050816fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff16101561180757600080fd5b61182c82826fffffffffffffffffffffffffffffffff1661361d90919063ffffffff16565b600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506118e782600260009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1661361d90919063ffffffff16565b600260006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550600660089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85846fffffffffffffffffffffffffffffffff166040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156119da57600080fd5b505af11580156119ee573d6000803e3d6000fd5b505050506040513d6020811015611a0457600080fd5b8101908080519060200190929190505050611a1b57fe5b5050505050565b6060600a805480602002602001604051908101604052809291908181526020018280548015611aa657602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311611a5c575b5050505050905090565b600260109054906101000a90046fffffffffffffffffffffffffffffffff1681565b6000600260109054906101000a90046fffffffffffffffffffffffffffffffff1690506000611c13600260009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16600660089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015611bca57600080fd5b505afa158015611bde573d6000803e3d6000fd5b505050506040513d6020811015611bf457600080fd5b810190808051906020019092919050505061359490919063ffffffff16565b905080600260106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555080826fffffffffffffffffffffffffffffffff1614611c9557807ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f60405160405180910390a25b5050565b6000611ca36136ca565b905090565b6000600a80549050905090565b3373ffffffffffffffffffffffffffffffffffffffff16600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614611d4f57600080fd5b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a350565b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6000600660049054906101000a900463ffffffff1663ffffffff16905090565b6000600660009054906101000a900463ffffffff1663ffffffff16905090565b6003601c9054906101000a900463ffffffff1681565b60055481565b6000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff169050919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612084576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4d7573742062652070726f706f736564206f776e65720000000000000000000081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f0d18b5fd22306e373229b9439188228edca81207d1667f604daf6cef8aa3ee6760405160405180910390a350565b6000612193613706565b905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6121c5611ad2565b505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461228c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b8363ffffffff8016600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff16146122f757600080fd5b600660009054906101000a900463ffffffff16600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff1602179055506000600a61238f600161237b611ca8565b63ffffffff1661376090919063ffffffff16565b63ffffffff168154811061239f57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506000600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900461ffff16905080600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81548161ffff021916908361ffff160217905550600760008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81549061ffff021916905581600a8261ffff16815481106124e357fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600a80548061253657fe5b6001900381819060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905590558673ffffffffffffffffffffffffffffffffffffffff167f9c8e7d83025bef8a04c664b2f753f64b8814bdb7e27291d7e50935f18cc3c71260405160405180910390a26125eb600360009054906101000a90046fffffffffffffffffffffffffffffffff168787876003601c9054906101000a900463ffffffff1661144a565b50505050505050565b60006125ff826137f5565b9050919050565b600061261182613821565b9050919050565b600080600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff168063ffffffff16905091509150915091565b600960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661271957600080fd5b6000600660009054906101000a900463ffffffff1690506000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16118061278457506127838161386b565b5b61278d57600080fd5b6127af6127aa60018363ffffffff1661394190919063ffffffff16565b6139d5565b50565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612874576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600260109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1610156128b157600080fd5b600660089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561295a57600080fd5b505af115801561296e573d6000803e3d6000fd5b505050506040513d602081101561298457600080fd5b810190808051906020019092919050505061299e57600080fd5b6129a6611ad2565b5050565b600360009054906101000a90046fffffffffffffffffffffffffffffffff1681565b600080600080600080600060086000600660009054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff1660086000600660009054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060020160000180549050101580612a845750612a83600660009054906101000a900463ffffffff1661386b565b5b905080612aa357600660009054906101000a900463ffffffff16612ad0565b612acf6001600660009054906101000a900463ffffffff1663ffffffff1661394190919063ffffffff16565b5b965086612ade898984613dac565b60086000600660049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020019081526020016000206000015483612b9157600860008b63ffffffff1663ffffffff16815260200190815260200160002060020160010160089054906101000a900463ffffffff1663ffffffff16600860008c63ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1601612b94565b60005b600260109054906101000a90046fffffffffffffffffffffffffffffffff1685612bff57600860008d63ffffffff1663ffffffff168152602001908152602001600020600201600101600c9054906101000a90046fffffffffffffffffffffffffffffffff16612c1f565b600360009054906101000a90046fffffffffffffffffffffffffffffffff165b9650965096509650965096505091939550919395565b600060086000600660009054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff16905090565b600360149054906101000a900463ffffffff1681565b600260009054906101000a90046fffffffffffffffffffffffffffffffff1681565b600360109054906101000a900463ffffffff1681565b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169050919050565b81600660009054906101000a900463ffffffff1663ffffffff168163ffffffff161480612dba5750612dab6001600660009054906101000a900463ffffffff1663ffffffff1661394190919063ffffffff16565b63ffffffff168163ffffffff16145b612dc357600080fd5b60018163ffffffff161480612df65750612df5612df060018363ffffffff1661376090919063ffffffff16565b61404b565b5b80612e1f5750612e1e612e1960018363ffffffff1661376090919063ffffffff16565b61386b565b5b612e2857600080fd5b826000600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160109054906101000a900463ffffffff16905060008163ffffffff161415612e9457600080fd5b8163ffffffff168163ffffffff161115612ead57600080fd5b8163ffffffff16600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff161015612f1857600080fd5b8163ffffffff16600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff1663ffffffff1610612f8257600080fd5b612f8b856139d5565b612f958486614097565b612f9e85614225565b612fa785614427565b612fb0856146a7565b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff16600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461305157600080fd5b80600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff167fb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b81043383604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a25050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613246576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b801515600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514156132a35761334d565b80600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f270d0c10dfbbdb6bb7206de0d1854b34e71664636d27af06feda4326a8d2437982604051808215151515815260200191505060405180910390a25b5050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613413576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b600281565b600080600660009054906101000a900463ffffffff16905060008163ffffffff16141580156135635750600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff168163ffffffff16145b15613571578091505061358f565b61358b60018263ffffffff1661394190919063ffffffff16565b9150505b919050565b60008282111561360c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b6000826fffffffffffffffffffffffffffffffff16826fffffffffffffffffffffffffffffffff1611156136b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b600060086000600660049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060000154905090565b600060086000600660049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16905090565b60008263ffffffff168263ffffffff1611156137e4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b6000600860008363ffffffff1663ffffffff168152602001908152602001600020600001549050919050565b6000600860008363ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff169050919050565b600080600860008463ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1690506000600860008563ffffffff1663ffffffff16815260200190815260200160002060020160010160089054906101000a900463ffffffff16905060008267ffffffffffffffff16118015613901575060008163ffffffff16115b801561393857504261392c8263ffffffff168467ffffffffffffffff166147b590919063ffffffff16565b67ffffffffffffffff16105b92505050919050565b60008082840190508363ffffffff168163ffffffff1610156139cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b80613a026001600660009054906101000a900463ffffffff1663ffffffff1661394190919063ffffffff16565b63ffffffff168163ffffffff161415613da857816000600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff169050600360189054906101000a900463ffffffff1663ffffffff1681018263ffffffff161180613aa15750600081145b15613da557613ac8613ac360018663ffffffff1661376090919063ffffffff16565b614851565b83600660006101000a81548163ffffffff021916908363ffffffff160217905550600360109054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160006101000a81548163ffffffff021916908363ffffffff160217905550600360149054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160046101000a81548163ffffffff021916908363ffffffff160217905550600360009054906101000a90046fffffffffffffffffffffffffffffffff16600860008663ffffffff1663ffffffff168152602001908152602001600020600201600101600c6101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506003601c9054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160086101000a81548163ffffffff021916908363ffffffff16021790555042600860008663ffffffff1663ffffffff16815260200190815260200160002060010160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c6101000a81548163ffffffff021916908363ffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168463ffffffff167f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271600860008863ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff16604051808267ffffffffffffffff16815260200191505060405180910390a35b50505b5050565b600080600760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160109054906101000a900463ffffffff16905060008163ffffffff161415613e1d576000915050614044565b8363ffffffff168163ffffffff161115613e3b576000915050614044565b8363ffffffff16600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff161015613eab576000915050614044565b8363ffffffff16600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff1663ffffffff1610613f1a576000915050614044565b8215613fef576000600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff169050600360189054906101000a900463ffffffff16810163ffffffff168563ffffffff1611158015613fad575060008163ffffffff16115b15613fbd57600092505050614044565b6000600360109054906101000a900463ffffffff1663ffffffff161415613fe957600092505050614044565b5061403e565b6000600860008663ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff16141561403d576000915050614044565b5b60019150505b9392505050565b600080600860008463ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16119050919050565b806000600860008363ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff1614156140e157600080fd5b600860008363ffffffff1663ffffffff16815260200190815260200160002060020160000183908060018154018082558091505060019003906000526020600020016000909190919091505581600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160186101000a81548163ffffffff021916908363ffffffff16021790555082600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101819055503373ffffffffffffffffffffffffffffffffffffffff168263ffffffff16847f92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c60405160405180910390a4505050565b80600860008263ffffffff1663ffffffff16815260200190815260200160002060020160010160049054906101000a900463ffffffff1663ffffffff16600860008363ffffffff1663ffffffff168152602001908152602001600020600201600001805490501061442357600061430e600860008563ffffffff1663ffffffff16815260200190815260200160002060020160000180548060200260200160405190810160405280929190818152602001828054801561430457602002820191906000526020600020905b8154815260200190600101908083116142f0575b5050505050614a8f565b905080600860008563ffffffff1663ffffffff1681526020019081526020016000206000018190555042600860008563ffffffff1663ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555082600860008563ffffffff1663ffffffff16815260200190815260200160002060010160106101000a81548163ffffffff021916908363ffffffff16021790555082600660046101000a81548163ffffffff021916908363ffffffff1602179055508263ffffffff16817f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a3505b5050565b6000600860008363ffffffff1663ffffffff168152602001908152602001600020600201600101600c9054906101000a90046fffffffffffffffffffffffffffffffff16905060006144b282600260109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1661361d90919063ffffffff16565b905080600260106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555061453082600260009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16614b7e90919063ffffffff16565b600260006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506145eb82600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16614b7e90919063ffffffff16565b600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550806fffffffffffffffffffffffffffffffff167ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f60405160405180910390a2505050565b80600860008263ffffffff1663ffffffff16815260200190815260200160002060020160010160009054906101000a900463ffffffff1663ffffffff16600860008363ffffffff1663ffffffff1681526020019081526020016000206002016000018054905014156147b157600860008363ffffffff1663ffffffff168152602001908152602001600020600201600080820160006147469190615547565b6001820160006101000a81549063ffffffff02191690556001820160046101000a81549063ffffffff02191690556001820160086101000a81549063ffffffff021916905560018201600c6101000a8154906fffffffffffffffffffffffffffffffff021916905550505b5050565b60008082840190508367ffffffffffffffff168167ffffffffffffffff161015614847576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b8061485b8161386b565b15614a8b578160006008600061488160018563ffffffff1661376090919063ffffffff16565b63ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff1614156148c857600080fd5b60006148e460018563ffffffff1661376090919063ffffffff16565b9050600860008263ffffffff1663ffffffff16815260200190815260200160002060000154600860008663ffffffff1663ffffffff16815260200190815260200160002060000181905550600860008263ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff16600860008663ffffffff1663ffffffff16815260200190815260200160002060010160106101000a81548163ffffffff021916908363ffffffff16021790555042600860008663ffffffff1663ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550600860008563ffffffff1663ffffffff16815260200190815260200160002060020160008082016000614a1e9190615547565b6001820160006101000a81549063ffffffff02191690556001820160046101000a81549063ffffffff02191690556001820160086101000a81549063ffffffff021916905560018201600c6101000a8154906fffffffffffffffffffffffffffffffff0219169055505050505b5050565b60008151600010614b08576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6c697374206d757374206e6f7420626520656d7074790000000000000000000081525060200191505060405180910390fd5b600082519050600060028281614b1a57fe5b049050600060028381614b2957fe5b061415614b6457600080614b47866000600187036001870387614c2a565b8092508193505050614b598282614d17565b945050505050614b79565b614b748460006001850384614db4565b925050505b919050565b6000808284019050836fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff161015614c20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b600080828410614c3957600080fd5b838611158015614c495750848411155b614c5257600080fd5b828611158015614c625750848311155b614c6b57600080fd5b5b600115614d0c5760078686031015614c9457614c8b8787878787614e4e565b91509150614d0d565b6000614ca18888886153c2565b9050808411614cb257809550614d06565b84811015614cc557600181019650614d05565b808511158015614cd457508381105b614cda57fe5b614ce688888388614db4565b9250614cf788600183018887614db4565b915082829250925050614d0d565b5b50614c6c565b5b9550959350505050565b60008083128015614d285750600082135b80614d3f5750600083138015614d3e5750600082125b5b15614d5f576002614d5084846154b9565b81614d5757fe5b059050614dae565b60006002808481614d6c57fe5b0760028681614d7757fe5b070181614d8057fe5b059050614daa614da460028681614d9357fe5b0560028681614d9e57fe5b056154b9565b826154b9565b9150505b92915050565b600081841115614dc357600080fd5b82821115614dd057600080fd5b5b82841015614e2f5760078484031015614e04576000614df38686868687614e4e565b809250819350505081915050614e46565b6000614e118686866153c2565b9050808311614e2257809350614e29565b6001810194505b50614dd1565b848481518110614e3b57fe5b602002602001015190505b949350505050565b600080600086600187010390506000886000890181518110614e6c57fe5b60200260200101519050600082600110614ea6577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff614ebe565b8960018a0181518110614eb557fe5b60200260200101515b9050600083600210614ef0577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff614f08565b8a60028b0181518110614eff57fe5b60200260200101515b9050600084600310614f3a577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff614f52565b8b60038c0181518110614f4957fe5b60200260200101515b9050600085600410614f84577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff614f9c565b8c60048d0181518110614f9357fe5b60200260200101515b9050600086600510614fce577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff614fe6565b8d60058e0181518110614fdd57fe5b60200260200101515b9050600087600610615018577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff615030565b8e60068f018151811061502757fe5b60200260200101515b90508587131561504557858780975081985050505b8385131561505857838580955081965050505b8183131561506b57818380935081945050505b8487131561507e57848780965081985050505b8386131561509157838680955081975050505b808313156150a457808380925081945050505b848613156150b757848680965081975050505b808213156150ca57808280925081935050505b828713156150dd57828780945081985050505b818613156150f057818680935081975050505b8085131561510357808580925081965050505b8286131561511657828680945081975050505b8084131561512957808480925081955050505b8285131561513c57828580945081965050505b8184131561514f57818480935081955050505b8284131561516257828480945081955050505b60008e8d039050600081141561517a57879a50615254565b600181141561518b57869a50615253565b600281141561519c57859a50615252565b60038114156151ad57849a50615251565b60048114156151be57839a50615250565b60058114156151cf57829a5061524f565b60068114156151e057819a5061524e565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f6b31206f7574206f6620626f756e64730000000000000000000000000000000081525060200191505060405180910390fd5b5b5b5b5b5b5b60008f8d0390508c8e1415615278578b8c9b509b50505050505050505050506153b8565b6000811415615296578b899b509b50505050505050505050506153b8565b60018114156152b4578b889b509b50505050505050505050506153b8565b60028114156152d2578b879b509b50505050505050505050506153b8565b60038114156152f0578b869b509b50505050505050505050506153b8565b600481141561530e578b859b509b50505050505050505050506153b8565b600581141561532c578b849b509b50505050505050505050506153b8565b600681141561534a578b839b509b50505050505050505050506153b8565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f6b32206f7574206f6620626f756e64730000000000000000000000000000000081525060200191505060405180910390fd5b9550959350505050565b600080846002848601816153d257fe5b04815181106153dd57fe5b602002602001015190506001840393506001830192505b6001156154b0575b6001840193508085858151811061540f57fe5b6020026020010151126153fc575b6001830392508085848151811061543057fe5b60200260200101511361541d57828410156154a25784838151811061545157fe5b602002602001015185858151811061546557fe5b602002602001015186868151811061547957fe5b6020026020010187868151811061548c57fe5b60200260200101828152508281525050506154ab565b829150506154b2565b6153f4565b505b9392505050565b6000808284019050600083121580156154d25750838112155b806154e857506000831280156154e757508381125b5b61553d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018061558e6021913960400191505060405180910390fd5b8091505092915050565b50805460008255906000526020600020908101906155659190615568565b50565b61558a91905b8082111561558657600081600090555060010161556e565b5090565b9056fe5369676e6564536166654d6174683a206164646974696f6e206f766572666c6f77a2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"

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
