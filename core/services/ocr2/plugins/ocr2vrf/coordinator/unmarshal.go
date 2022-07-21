package coordinator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	vrf_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/ocr2vrf/generated/vrf_beacon_coordinator"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
)

func unmarshalRandomnessRequested(lg logpoller.Log) (r vrf_wrapper.VRFBeaconCoordinatorRandomnessRequested, err error) {
	args := unindexedArgs(vrfABI, randomnessRequestedEvent)

	m := make(map[string]any)
	err = args.UnpackIntoMap(m, lg.Data)
	if err != nil {
		return r, errors.Wrapf(err, "unpack %s into map (RandomnessRequested)", hexutil.Encode(lg.Data))
	}

	r.ConfDelay = *abi.ConvertType(m["confDelay"], new(*big.Int)).(**big.Int)
	r.NextBeaconOutputHeight = *abi.ConvertType(m["nextBeaconOutputHeight"], new(uint64)).(*uint64)
	r.Raw = types.Log{
		Data:        lg.Data,
		Address:     lg.Address,
		BlockHash:   lg.BlockHash,
		BlockNumber: uint64(lg.BlockNumber),
		TxHash:      lg.TxHash,
		Index:       uint(lg.LogIndex),
		Removed:     false,
	}

	return
}

func unmarshalRandomnessFulfillmentRequested(lg logpoller.Log) (r vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested, err error) {
	args := unindexedArgs(vrfABI, randomnessFulfillmentRequestedEvent)

	m := make(map[string]any)
	err = args.UnpackIntoMap(m, lg.Data)
	if err != nil {
		return r, errors.Wrapf(err, "unpack %s into map (RandomnessFulfillmentRequested)", hexutil.Encode(lg.Data))
	}

	r.ConfDelay = *abi.ConvertType(m["confDelay"], new(*big.Int)).(**big.Int)
	r.NextBeaconOutputHeight = *abi.ConvertType(m["nextBeaconOutputHeight"], new(uint64)).(*uint64)
	r.SubID = *abi.ConvertType(m["subID"], new(uint64)).(*uint64)
	r.Callback = *abi.ConvertType(m["callback"], new(vrf_wrapper.VRFBeaconTypesCallback)).(*vrf_wrapper.VRFBeaconTypesCallback)
	r.Raw = types.Log{
		Data:        lg.Data,
		Address:     lg.Address,
		BlockHash:   lg.BlockHash,
		BlockNumber: uint64(lg.BlockNumber),
		TxHash:      lg.TxHash,
		Index:       uint(lg.LogIndex),
		Removed:     false,
	}

	return
}

func unmarshalRandomWordsFulfilled(lg logpoller.Log) (r vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled, err error) {
	args := unindexedArgs(vrfABI, randomWordsFulfilledEvent)

	m := make(map[string]any)
	err = args.UnpackIntoMap(m, lg.Data)
	if err != nil {
		return r, errors.Wrapf(err, "unpack %s into map (RandomWordsFulfilled)", hexutil.Encode(lg.Data))
	}

	r.SuccessfulFulfillment = *abi.ConvertType(m["successfulFulfillment"], new([]byte)).(*[]byte)
	r.TruncatedErrorData = *abi.ConvertType(m["truncatedErrorData"], new([][]byte)).(*[][]byte)
	r.RequestIDs = *abi.ConvertType(m["requestIDs"], new([]*big.Int)).(*[]*big.Int)
	r.Raw = types.Log{
		Data:        lg.Data,
		Address:     lg.Address,
		BlockHash:   lg.BlockHash,
		BlockNumber: uint64(lg.BlockNumber),
		TxHash:      lg.TxHash,
		Index:       uint(lg.LogIndex),
		Removed:     false,
	}

	return
}

func unmarshalNewTransmission(lg logpoller.Log) (r vrf_wrapper.VRFBeaconCoordinatorNewTransmission, err error) {
	args := unindexedArgs(vrfABI, newTransmissionEvent)

	m := make(map[string]any)
	err = args.UnpackIntoMap(m, lg.Data)
	if err != nil {
		return r, errors.Wrapf(err, "unpack %s into map (RandomWordsFulfilled)", hexutil.Encode(lg.Data))
	}

	r.EpochAndRound = *abi.ConvertType(m["epochAndRound"], new(*big.Int)).(**big.Int)
	r.OutputsServed = *abi.ConvertType(m["outputsServed"], new([]vrf_wrapper.VRFBeaconReportOutputServed)).(*[]vrf_wrapper.VRFBeaconReportOutputServed)
	r.ConfigDigest = *abi.ConvertType(m["configDigest"], new([32]byte)).(*[32]byte)
	r.JuelsPerFeeCoin = *abi.ConvertType(m["juelsPerFeeCoin"], new(*big.Int)).(**big.Int)
	r.Transmitter = *abi.ConvertType(m["transmitter"], new(common.Address)).(*common.Address)
	r.AggregatorRoundId = *abi.ConvertType(m["aggregatorRoundId"], new(uint32)).(*uint32)
	r.Raw = types.Log{
		Data:        lg.Data,
		Address:     lg.Address,
		BlockHash:   lg.BlockHash,
		BlockNumber: uint64(lg.BlockNumber),
		TxHash:      lg.TxHash,
		Index:       uint(lg.LogIndex),
		Removed:     false,
	}

	return
}

func unindexedArgs(tabi abi.ABI, eventName string) (u abi.Arguments) {
	for _, a := range tabi.Events[eventName].Inputs {
		u = append(u, abi.Argument{
			Name:    a.Name,
			Type:    a.Type,
			Indexed: false,
		})
	}
	return
}
