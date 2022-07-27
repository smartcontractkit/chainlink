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
	//event RandomnessRequested(
	//  uint64 indexed nextBeaconOutputHeight,
	//  ConfirmationDelay confDelay
	//);
	args := vrfABI.Events[randomnessRequestedEvent].Inputs

	// unpack the non-indexed args into a map all together,
	// and unpack the indexed args separately.
	m := make(map[string]any)
	err = args.UnpackIntoMap(m, lg.Data)
	if err != nil {
		return r, errors.Wrapf(err, "unpack %s into map (RandomnessRequested)", hexutil.Encode(lg.Data))
	}

	r.ConfDelay = *abi.ConvertType(m["confDelay"], new(*big.Int)).(**big.Int)
	nextBeaconOutputHeightType, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return r, errors.Wrap(err, "abi NewType uint64")
	}
	indexedArgs := abi.Arguments{abi.Argument{
		Name: "nextBeaconOutputHeight",
		Type: nextBeaconOutputHeightType,
	}}
	nextBeaconOutputHeightInterface, err := indexedArgs.Unpack(lg.Topics[1])
	if err != nil {
		return r, errors.Wrap(err, "unpack nextBeaconOutputHeight")
	}
	r.NextBeaconOutputHeight = *abi.ConvertType(nextBeaconOutputHeightInterface[0], new(uint64)).(*uint64)
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
	//event RandomnessFulfillmentRequested(
	//  uint64 nextBeaconOutputHeight,
	//  ConfirmationDelay confDelay,
	//  uint64 subID,
	//  Callback callback
	//);
	args := vrfABI.Events[randomnessFulfillmentRequestedEvent].Inputs

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
	//event RandomWordsFulfilled(
	//  RequestID[] requestIDs,
	//  bytes successfulFulfillment,
	//  bytes[] truncatedErrorData
	//);
	args := vrfABI.Events[randomWordsFulfilledEvent].Inputs

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
	//event NewTransmission(
	//  uint32 indexed aggregatorRoundId,
	//  uint40 indexed epochAndRound,
	//  address transmitter,
	//  uint192 juelsPerFeeCoin,
	//  bytes32 configDigest,
	//  OutputServed[] outputsServed
	//);
	args := vrfABI.Events[newTransmissionEvent].Inputs

	m := make(map[string]any)
	err = args.UnpackIntoMap(m, lg.Data)
	if err != nil {
		return r, errors.Wrapf(err, "unpack %s into map (RandomWordsFulfilled)", hexutil.Encode(lg.Data))
	}

	r.OutputsServed = *abi.ConvertType(m["outputsServed"], new([]vrf_wrapper.VRFBeaconReportOutputServed)).(*[]vrf_wrapper.VRFBeaconReportOutputServed)
	r.ConfigDigest = *abi.ConvertType(m["configDigest"], new([32]byte)).(*[32]byte)
	r.JuelsPerFeeCoin = *abi.ConvertType(m["juelsPerFeeCoin"], new(*big.Int)).(**big.Int)
	r.Transmitter = *abi.ConvertType(m["transmitter"], new(common.Address)).(*common.Address)

	// aggregatorRoundId is indexed
	aggregatorRoundIDType, err := abi.NewType("uint32", "", nil)
	if err != nil {
		return r, errors.Wrap(err, "abi NewType uint32")
	}
	indexedArgs := abi.Arguments{
		{
			Name: "aggregatorRoundId",
			Type: aggregatorRoundIDType,
		},
	}
	aggregatorRoundIDInterface, err := indexedArgs.Unpack(lg.Topics[1])
	if err != nil {
		return r, errors.Wrap(err, "unpack aggregatorRoundId")
	}
	r.AggregatorRoundId = *abi.ConvertType(aggregatorRoundIDInterface[0], new(uint32)).(*uint32)

	// epochAndRound is indexed
	epochAndRoundType, err := abi.NewType("uint40", "", nil)
	if err != nil {
		return r, errors.Wrap(err, "abi NewType uint40")
	}
	indexedArgs = abi.Arguments{
		{
			Name: "epochAndRound",
			Type: epochAndRoundType,
		},
	}
	epochAndRoundInterface, err := indexedArgs.Unpack(lg.Topics[2])
	if err != nil {
		return r, errors.Wrap(err, "unpack epochAndRound")
	}
	r.EpochAndRound = *abi.ConvertType(epochAndRoundInterface[0], new(*big.Int)).(**big.Int)

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
