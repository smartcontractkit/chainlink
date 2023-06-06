package blockhashstore

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocklp "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// VRF-only events.
	randomWordsRequestedV2              string = "RandomWordsRequested"
	randomWordsFulfilledV2              string = "RandomWordsFulfilled"
	randomWordsRequestedV1              string = "RandomnessRequest"
	randomWordsFulfilledV1              string = "RandomnessRequestFulfilled"
	randomnessFulfillmentRequestedEvent string = "RandomnessFulfillmentRequested"
	randomWordsFulfilledEvent           string = "RandomWordsFulfilled"
	newTransmissionEvent                string = "NewTransmission"
	outputsServedEvent                  string = "OutputsServed"
)

var (
	vrfCoordinatorV2ABI = evmtypes.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2MetaData.ABI)
	vrfCoordinatorV1ABI = evmtypes.MustGetABI(solidity_vrf_coordinator_interface.VRFCoordinatorMetaData.ABI)

	_     Coordinator = &TestCoordinator{}
	_     BHS         = &TestBHS{}
	tests             = []struct {
		name           string
		requests       []Event
		fulfillments   []Event
		wait           int
		lookback       int
		latest         uint64
		bhs            TestBHS
		expectedStored []uint64
		expectedErrMsg string
	}{
		{
			name:           "single unfulfilled request",
			requests:       []Event{{Block: 150, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150},
		},
		{
			name:           "single fulfilled request",
			requests:       []Event{{Block: 150, ID: "1000"}},
			fulfillments:   []Event{{Block: 155, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "single already fulfilled",
			requests:       []Event{{Block: 150, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            TestBHS{Stored: []uint64{150}},
			expectedStored: []uint64{150},
		},
		{
			name:           "error checking if stored, store anyway",
			requests:       []Event{{Block: 150, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            TestBHS{ErrorsIsStored: []uint64{150}},
			expectedStored: []uint64{150},
			expectedErrMsg: "checking if stored: error checking if stored",
		},
		{
			name:           "error storing, continue to next block anyway",
			requests:       []Event{{Block: 150, ID: "1000"}, {Block: 151, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            TestBHS{ErrorsStore: []uint64{150}},
			expectedStored: []uint64{151},
			expectedErrMsg: "storing block: error storing",
		},
		{
			name: "multiple requests same block, some fulfilled",
			requests: []Event{
				{Block: 150, ID: "10001"},
				{Block: 150, ID: "10002"},
				{Block: 150, ID: "10003"}},
			fulfillments: []Event{
				{Block: 150, ID: "10001"},
				{Block: 150, ID: "10003"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150},
		},
		{
			name: "multiple requests same block, all fulfilled",
			requests: []Event{
				{Block: 150, ID: "10001"},
				{Block: 150, ID: "10002"},
				{Block: 150, ID: "10003"}},
			fulfillments: []Event{
				{Block: 150, ID: "10001"},
				{Block: 150, ID: "10002"},
				{Block: 150, ID: "10003"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "fulfillment no matching request no error",
			requests:       []Event{{Block: 150, ID: "1000"}},
			fulfillments:   []Event{{Block: 199, ID: "10002"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150},
		},
		{
			name:           "multiple unfulfilled requests",
			requests:       []Event{{Block: 150, ID: "10001"}, {Block: 151, ID: "10002"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150, 151},
		},
		{
			name:           "multiple fulfilled requests",
			requests:       []Event{{Block: 150, ID: "10001"}, {Block: 151, ID: "10002"}},
			fulfillments:   []Event{{Block: 150, ID: "10001"}, {Block: 151, ID: "10002"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "recent unfulfilled request do not store",
			requests:       []Event{{Block: 185, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "old unfulfilled request do not store",
			requests:       []Event{{Block: 99, ID: "1000"}, {Block: 57, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name: "mixed",
			requests: []Event{
				// Block 150
				{Block: 150, ID: "10001"},
				{Block: 150, ID: "10002"},
				{Block: 150, ID: "10003"},

				// Block 151
				{Block: 151, ID: "10004"},
				{Block: 151, ID: "10005"},

				// Block 153
				{Block: 153, ID: "10006"},

				// Block 154
				{Block: 154, ID: "10007"}},
			fulfillments: []Event{
				// Block 150
				{Block: 150, ID: "10001"},
				{Block: 150, ID: "10002"},
				// request3 no fulfillment

				// Block 151
				{Block: 151, ID: "10004"},
				{Block: 151, ID: "10005"},

				// Block 153 - no fulfillment

				// Block 154
				{Block: 154, ID: "10007"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150, 153},
		},
		{
			name:           "lookback before 0th block",
			requests:       []Event{{Block: 20, ID: "1000"}},
			wait:           25,
			lookback:       100,
			latest:         50,
			expectedStored: []uint64{20},
		},
	}
)

func TestFeeder(t *testing.T) {

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			coordinator := &TestCoordinator{
				RequestEvents:     test.requests,
				FulfillmentEvents: test.fulfillments,
			}

			feeder := NewFeeder(
				logger.TestLogger(t),
				coordinator,
				&test.bhs,
				test.wait,
				test.lookback,
				func(ctx context.Context) (uint64, error) {
					return test.latest, nil
				})

			err := feeder.Run(testutils.Context(t))
			if test.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedErrMsg)
			}

			require.ElementsMatch(t, test.expectedStored, test.bhs.Stored)
		})
	}
}

func TestFeederWithLogPollerVRFv1(t *testing.T) {

	var coordinatorAddress = common.HexToAddress("0x514910771AF9Ca656af840dff83E8264EcF986CA")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Instantiate log poller & coordinator.
			lp := &mocklp.LogPoller{}
			lp.On("RegisterFilter", mock.Anything).Return(nil)
			c, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(coordinatorAddress, nil)
			require.NoError(t, err)
			coordinator := &V1Coordinator{
				c:  c,
				lp: lp,
			}

			// Assert search window.
			latest := int64(test.latest)
			fromBlock := mathutil.Max(latest-int64(test.lookback), 0)
			toBlock := mathutil.Max(latest-int64(test.wait), 0)

			// Construct request logs.
			var requestLogs []logpoller.Log
			for _, r := range test.requests {
				if r.Block < uint64(fromBlock) || r.Block > uint64(toBlock) {
					continue // do not include blocks outside our search window
				}
				requestLogs = append(
					requestLogs,
					newRandomnessRequestedLogV1(t, r.Block, r.ID, coordinatorAddress),
				)
			}

			// Construct fulfillment logs.
			var fulfillmentLogs []logpoller.Log
			for _, r := range test.fulfillments {
				fulfillmentLogs = append(
					fulfillmentLogs,
					newRandomnessFulfilledLogV1(t, r.Block, r.ID, coordinatorAddress),
				)
			}

			// Mock log poller.
			lp.On("LatestBlock", mock.Anything).
				Return(latest, nil)
			lp.On(
				"LogsWithSigs",
				fromBlock,
				toBlock,
				[]common.Hash{
					solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic(),
				},
				coordinatorAddress,
				mock.Anything,
			).Return(requestLogs, nil)
			lp.On(
				"LogsWithSigs",
				fromBlock,
				latest,
				[]common.Hash{
					solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{}.Topic(),
				},
				coordinatorAddress,
				mock.Anything,
			).Return(fulfillmentLogs, nil)

			// Instantiate feeder.
			feeder := NewFeeder(
				logger.TestLogger(t),
				coordinator,
				&test.bhs,
				test.wait,
				test.lookback,
				func(ctx context.Context) (uint64, error) {
					return test.latest, nil
				})

			// Run feeder and assert correct results.
			err = feeder.Run(testutils.Context(t))
			if test.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedErrMsg)
			}
			require.ElementsMatch(t, test.expectedStored, test.bhs.Stored)
		})
	}
}

func TestFeederWithLogPollerVRFv2(t *testing.T) {

	var coordinatorAddress = common.HexToAddress("0x514910771AF9Ca656af840dff83E8264EcF986CA")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Instantiate log poller & coordinator.
			lp := &mocklp.LogPoller{}
			lp.On("RegisterFilter", mock.Anything).Return(nil)
			c, err := vrf_coordinator_v2.NewVRFCoordinatorV2(coordinatorAddress, nil)
			require.NoError(t, err)
			coordinator := &V2Coordinator{
				c:  c,
				lp: lp,
			}

			// Assert search window.
			latest := int64(test.latest)
			fromBlock := mathutil.Max(latest-int64(test.lookback), 0)
			toBlock := mathutil.Max(latest-int64(test.wait), 0)

			// Construct request logs.
			var requestLogs []logpoller.Log
			for _, r := range test.requests {
				if r.Block < uint64(fromBlock) || r.Block > uint64(toBlock) {
					continue // do not include blocks outside our search window
				}
				reqId, ok := big.NewInt(0).SetString(r.ID, 10)
				require.True(t, ok)
				requestLogs = append(
					requestLogs,
					newRandomnessRequestedLogV2(t, r.Block, reqId, coordinatorAddress),
				)
			}

			// Construct fulfillment logs.
			var fulfillmentLogs []logpoller.Log
			for _, r := range test.fulfillments {
				reqId, ok := big.NewInt(0).SetString(r.ID, 10)
				require.True(t, ok)
				fulfillmentLogs = append(
					fulfillmentLogs,
					newRandomnessFulfilledLogV2(t, r.Block, reqId, coordinatorAddress),
				)
			}

			// Mock log poller.
			lp.On("LatestBlock", mock.Anything).
				Return(latest, nil)
			lp.On(
				"LogsWithSigs",
				fromBlock,
				toBlock,
				[]common.Hash{
					vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(),
				},
				coordinatorAddress,
				mock.Anything,
			).Return(requestLogs, nil)
			lp.On(
				"LogsWithSigs",
				fromBlock,
				latest,
				[]common.Hash{
					vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled{}.Topic(),
				},
				coordinatorAddress,
				mock.Anything,
			).Return(fulfillmentLogs, nil)

			// Instantiate feeder.
			feeder := NewFeeder(
				logger.TestLogger(t),
				coordinator,
				&test.bhs,
				test.wait,
				test.lookback,
				func(ctx context.Context) (uint64, error) {
					return test.latest, nil
				})

			// Run feeder and assert correct results.
			err = feeder.Run(testutils.Context(t))
			if test.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedErrMsg)
			}
			require.ElementsMatch(t, test.expectedStored, test.bhs.Stored)
		})
	}
}

func TestFeeder_CachesStoredBlocks(t *testing.T) {
	coordinator := &TestCoordinator{
		RequestEvents: []Event{{Block: 100, ID: "1000"}},
	}

	bhs := &TestBHS{}

	feeder := NewFeeder(
		logger.TestLogger(t),
		coordinator,
		bhs,
		100,
		200,
		func(ctx context.Context) (uint64, error) {
			return 250, nil
		})

	// Should store block 100
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.ElementsMatch(t, []uint64{100}, bhs.Stored)

	// Remove 100 from the BHS and try again, it should not be stored since it's cached in the
	// feeder
	bhs.Stored = nil
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, bhs.Stored)

	// Run the feeder on a later block and make sure the cache is pruned
	feeder.latestBlock = func(ctx context.Context) (uint64, error) {
		return 500, nil
	}
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, feeder.stored)
}

func newRandomnessRequestedLogV1(
	t *testing.T,
	requestBlock uint64,
	requestID string,
	coordinatorAddress common.Address,
) logpoller.Log {
	e := solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
		KeyHash:   common.HexToHash("keyhash"),
		Seed:      big.NewInt(0),
		Sender:    common.Address{},
		JobID:     common.HexToHash("job"),
		Fee:       big.NewInt(0),
		RequestID: common.HexToHash(requestID),
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorV1ABI.Events[randomWordsRequestedV1].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(
		e.KeyHash,
		e.Seed,
		e.Sender,
		e.Fee,
		e.RequestID,
	)
	require.NoError(t, err)

	jobIDType, err := abi.NewType("bytes32", "", nil)
	require.NoError(t, err)

	jobIDArg := abi.Arguments{abi.Argument{
		Name:    "jobID",
		Type:    jobIDType,
		Indexed: true,
	}}

	topic1, err := jobIDArg.Pack(e.JobID)
	require.NoError(t, err)

	topic0 := vrfCoordinatorV1ABI.Events[randomWordsRequestedV1].ID
	lg := logpoller.Log{
		Address: coordinatorAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			// first topic is the event signature
			topic0.Bytes(),
			// second topic is JobID since it's indexed
			topic1,
		},
		BlockNumber: int64(requestBlock),
		EventSig:    topic0,
	}
	return lg
}

func newRandomnessFulfilledLogV1(
	t *testing.T,
	requestBlock uint64,
	requestID string,
	coordinatorAddress common.Address,
) logpoller.Log {
	e := solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{
		RequestId: common.HexToHash(requestID),
		Output:    big.NewInt(0),
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorV1ABI.Events[randomWordsFulfilledV1].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(
		e.RequestId,
		e.Output,
	)
	require.NoError(t, err)

	topic0 := vrfCoordinatorV1ABI.Events[randomWordsFulfilledV1].ID
	lg := logpoller.Log{
		Address: coordinatorAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			// first topic is the event signature
			topic0.Bytes(),
		},
		BlockNumber: int64(requestBlock),
		EventSig:    topic0,
	}
	return lg
}

func newRandomnessRequestedLogV2(
	t *testing.T,
	requestBlock uint64,
	requestID *big.Int,
	coordinatorAddress common.Address,
) logpoller.Log {
	e := vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   requestID,
		PreSeed:                     big.NewInt(0),
		MinimumRequestConfirmations: 0,
		CallbackGasLimit:            0,
		NumWords:                    0,
		Sender:                      common.HexToAddress("0xeFF41C8725be95e66F6B10489B6bF34b08055853"),
		Raw: types.Log{
			BlockNumber: requestBlock,
		},
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorV2ABI.Events[randomWordsRequestedV2].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(
		e.RequestId,
		e.PreSeed,
		e.MinimumRequestConfirmations,
		e.CallbackGasLimit,
		e.NumWords,
	)
	require.NoError(t, err)

	keyHashType, err := abi.NewType("bytes32", "", nil)
	require.NoError(t, err)

	subIdType, err := abi.NewType("uint64", "", nil)
	require.NoError(t, err)

	senderType, err := abi.NewType("address", "", nil)
	require.NoError(t, err)

	keyHashArg := abi.Arguments{abi.Argument{
		Name:    "keyHash",
		Type:    keyHashType,
		Indexed: true,
	}}
	subIdArg := abi.Arguments{abi.Argument{
		Name:    "subId",
		Type:    subIdType,
		Indexed: true,
	}}

	senderArg := abi.Arguments{abi.Argument{
		Name:    "sender",
		Type:    senderType,
		Indexed: true,
	}}

	topic1, err := keyHashArg.Pack(e.KeyHash)
	require.NoError(t, err)
	topic2, err := subIdArg.Pack(e.SubId)
	require.NoError(t, err)
	topic3, err := senderArg.Pack(e.Sender)
	require.NoError(t, err)

	topic0 := vrfCoordinatorV2ABI.Events[randomWordsRequestedV2].ID
	lg := logpoller.Log{
		Address: coordinatorAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			// first topic is the event signature
			topic0.Bytes(),
			// second topic is keyHash since it's indexed
			topic1,
			// third topic is subId since it's indexed
			topic2,
			// third topic is sender since it's indexed
			topic3,
		},
		BlockNumber: int64(requestBlock),
		EventSig:    topic0,
	}
	return lg
}

func newRandomnessFulfilledLogV2(
	t *testing.T,
	requestBlock uint64,
	requestID *big.Int,
	coordinatorAddress common.Address,
) logpoller.Log {
	e := vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled{
		RequestId:  requestID,
		OutputSeed: big.NewInt(0),
		Payment:    big.NewInt(0),
		Success:    true,
		Raw: types.Log{
			BlockNumber: requestBlock,
		},
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorV2ABI.Events[randomWordsFulfilledV2].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(
		e.OutputSeed,
		e.Payment,
		e.Success,
	)
	require.NoError(t, err)

	requestIdType, err := abi.NewType("uint256", "", nil)
	require.NoError(t, err)

	requestIdArg := abi.Arguments{abi.Argument{
		Name:    "requestId",
		Type:    requestIdType,
		Indexed: true,
	}}

	topic1, err := requestIdArg.Pack(e.RequestId)
	require.NoError(t, err)

	topic0 := vrfCoordinatorV2ABI.Events[randomWordsFulfilledV2].ID
	lg := logpoller.Log{
		Address: coordinatorAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			// first topic is the event signature
			topic0.Bytes(),
			// second topic is requestId since it's indexed
			topic1,
		},
		BlockNumber: int64(requestBlock),
		EventSig:    topic0,
	}
	return lg
}
