package v2

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	clnull "github.com/smartcontractkit/chainlink-common/pkg/utils/null"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func makeTestTxm(t *testing.T, txStore txmgr.TestEvmTxStore, keyStore keystore.Master) txmgrcommon.TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee] {
	_, _, evmConfig := txmgr.MakeTestConfigs(t)
	ec := evmtest.NewEthClientMockWithDefaultChain(t)
	txmConfig := txmgr.NewEvmTxmConfig(evmConfig)
	txm := txmgr.NewEvmTxm(ec.ConfiguredChainID(), txmConfig, evmConfig.Transactions(), keyStore.Eth(), logger.TestLogger(t), nil, nil,
		nil, txStore, nil, nil, nil, nil)

	return txm
}

func MakeTestListenerV2(chain legacyevm.Chain) *listenerV2 {
	return &listenerV2{chainID: chain.Client().ConfiguredChainID(), chain: chain}
}

func txMetaSubIDs(t *testing.T, vrfVersion vrfcommon.Version, subID *big.Int) (*uint64, *string) {
	var (
		txMetaSubID       *uint64
		txMetaGlobalSubID *string
	)
	if vrfVersion == vrfcommon.V2Plus {
		txMetaGlobalSubID = ptr(subID.String())
	} else if vrfVersion == vrfcommon.V2 {
		txMetaSubID = ptr(subID.Uint64())
	} else {
		t.Errorf("unsupported vrf version: %s", vrfVersion)
	}
	return txMetaSubID, txMetaGlobalSubID
}

func addEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, from common.Address, state txmgrtypes.TxState, maxLink string, subID *big.Int, reqTxHash common.Hash, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	b, err := json.Marshal(txmgr.TxMeta{
		MaxLink:       &maxLink,
		SubID:         txMetaSubID,
		GlobalSubID:   txMetaGlobalSubID,
		RequestTxHash: &reqTxHash,
	})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	tx := &txmgr.Tx{
		FromAddress:       from,
		ToAddress:         from,
		EncodedPayload:    []byte(`blah`),
		Value:             *big.NewInt(0),
		FeeLimit:          0,
		State:             state,
		Meta:              &meta,
		Subject:           uuid.NullUUID{},
		ChainID:           testutils.SimulatedChainID,
		MinConfirmations:  clnull.Uint32{Uint32: 0},
		PipelineTaskRunID: uuid.NullUUID{},
	}
	err = txStore.InsertTx(testutils.Context(t), tx)
	require.NoError(t, err)
}

func addConfirmedEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, from common.Address, maxLink string, subID *big.Int, nonce evmtypes.Nonce, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	b, err := json.Marshal(txmgr.TxMeta{
		MaxLink:     &maxLink,
		SubID:       txMetaSubID,
		GlobalSubID: txMetaGlobalSubID,
	})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	now := time.Now()

	tx := &txmgr.Tx{
		Sequence:           &nonce,
		FromAddress:        from,
		ToAddress:          from,
		EncodedPayload:     []byte(`blah`),
		Value:              *big.NewInt(0),
		FeeLimit:           0,
		State:              txmgrcommon.TxConfirmed,
		Meta:               &meta,
		Subject:            uuid.NullUUID{},
		ChainID:            testutils.SimulatedChainID,
		MinConfirmations:   clnull.Uint32{Uint32: 0},
		PipelineTaskRunID:  uuid.NullUUID{},
		BroadcastAt:        &now,
		InitialBroadcastAt: &now,
	}
	err = txStore.InsertTx(testutils.Context(t), tx)
	require.NoError(t, err)
}

func addEthTxNativePayment(t *testing.T, txStore txmgr.TestEvmTxStore, from common.Address, state txmgrtypes.TxState, maxNative string, subID *big.Int, reqTxHash common.Hash, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	b, err := json.Marshal(txmgr.TxMeta{
		MaxEth:        &maxNative,
		SubID:         txMetaSubID,
		GlobalSubID:   txMetaGlobalSubID,
		RequestTxHash: &reqTxHash,
	})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	tx := &txmgr.Tx{
		FromAddress:       from,
		ToAddress:         from,
		EncodedPayload:    []byte(`blah`),
		Value:             *big.NewInt(0),
		FeeLimit:          0,
		State:             state,
		Meta:              &meta,
		Subject:           uuid.NullUUID{},
		ChainID:           testutils.SimulatedChainID,
		MinConfirmations:  clnull.Uint32{Uint32: 0},
		PipelineTaskRunID: uuid.NullUUID{},
	}
	err = txStore.InsertTx(testutils.Context(t), tx)
	require.NoError(t, err)
}

func addConfirmedEthTxNativePayment(t *testing.T, txStore txmgr.TestEvmTxStore, from common.Address, maxNative string, subID *big.Int, nonce evmtypes.Nonce, vrfVersion vrfcommon.Version) {
	txMetaSubID, txMetaGlobalSubID := txMetaSubIDs(t, vrfVersion, subID)
	b, err := json.Marshal(txmgr.TxMeta{
		MaxEth:      &maxNative,
		SubID:       txMetaSubID,
		GlobalSubID: txMetaGlobalSubID,
	})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	now := time.Now()
	tx := &txmgr.Tx{
		Sequence:           &nonce,
		FromAddress:        from,
		ToAddress:          from,
		EncodedPayload:     []byte(`blah`),
		Value:              *big.NewInt(0),
		FeeLimit:           0,
		State:              txmgrcommon.TxConfirmed,
		Meta:               &meta,
		Subject:            uuid.NullUUID{},
		ChainID:            testutils.SimulatedChainID,
		MinConfirmations:   clnull.Uint32{Uint32: 0},
		PipelineTaskRunID:  uuid.NullUUID{},
		BroadcastAt:        &now,
		InitialBroadcastAt: &now,
	}
	err = txStore.InsertTx(testutils.Context(t), tx)
	require.NoError(t, err)
}

func testMaybeSubtractReservedLink(t *testing.T, vrfVersion vrfcommon.Version) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock(ctx, "blah"))
	chainID := testutils.SimulatedChainID
	k, err := ks.Eth().Create(testutils.Context(t), chainID)
	require.NoError(t, err)

	subID := new(big.Int).SetUint64(1)
	reqTxHash := common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")

	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
	}).Toml())
	require.NoError(t, err)
	txstore := txmgr.NewTxStore(db, lggr)
	txm := makeTestTxm(t, txstore, ks)
	chain := evmmocks.NewChain(t)
	chain.On("TxManager").Return(txm)
	listener := &listenerV2{
		respCount: map[string]uint64{},
		job:       j,
		chain:     chain,
	}

	// Insert an unstarted eth tx with link metadata
	addEthTx(t, txstore, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err := listener.MaybeSubtractReservedLink(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)

	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTx(t, txstore, k.Address, "10000", subID, 1, vrfVersion)
	start, err = listener.MaybeSubtractReservedLink(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTx(t, txstore, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedLink(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved link should not affect other subscribers prospective balance.
	otherSubID := new(big.Int).SetUint64(2)
	require.NoError(t, err)
	addEthTx(t, txstore, k.Address, txmgrcommon.TxUnstarted, "10000", otherSubID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedLink(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)

	anotherSubID := new(big.Int).SetUint64(3)
	addEthTx(t, txstore, k2.Address, txmgrcommon.TxUnstarted, "10000", anotherSubID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedLink(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the link reserved across multiple keys,
	// i.e, gas lanes.
	addEthTx(t, txstore, k2.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedLink(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "70000", start.String())
}

func TestMaybeSubtractReservedLinkV2(t *testing.T) {
	testMaybeSubtractReservedLink(t, vrfcommon.V2)
}

func TestMaybeSubtractReservedLinkV2Plus(t *testing.T) {
	testMaybeSubtractReservedLink(t, vrfcommon.V2Plus)
}

func testMaybeSubtractReservedNative(t *testing.T, vrfVersion vrfcommon.Version) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock(ctx, "blah"))
	chainID := testutils.SimulatedChainID
	k, err := ks.Eth().Create(testutils.Context(t), chainID)
	require.NoError(t, err)

	subID := new(big.Int).SetUint64(1)
	reqTxHash := common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")

	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
	}).Toml())
	require.NoError(t, err)
	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	txm := makeTestTxm(t, txstore, ks)
	require.NoError(t, err)
	chain := evmmocks.NewChain(t)
	chain.On("TxManager").Return(txm)
	listener := &listenerV2{
		respCount: map[string]uint64{},
		job:       j,
		chain:     chain,
	}

	// Insert an unstarted eth tx with native metadata
	addEthTxNativePayment(t, txstore, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err := listener.MaybeSubtractReservedEth(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)

	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// A confirmed tx should not affect the starting balance
	addConfirmedEthTxNativePayment(t, txstore, k.Address, "10000", subID, 1, vrfVersion)
	start, err = listener.MaybeSubtractReservedEth(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "90000", start.String())

	// An unconfirmed tx _should_ affect the starting balance.
	addEthTxNativePayment(t, txstore, k.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedEth(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	assert.Equal(t, "80000", start.String())

	// One subscriber's reserved native should not affect other subscribers prospective balance.
	otherSubID := new(big.Int).SetUint64(2)
	require.NoError(t, err)
	addEthTxNativePayment(t, txstore, k.Address, txmgrcommon.TxUnstarted, "10000", otherSubID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedEth(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// One key's data should not affect other keys' data in the case of different subscribers.
	k2, err := ks.Eth().Create(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)

	anotherSubID := new(big.Int).SetUint64(3)
	addEthTxNativePayment(t, txstore, k2.Address, txmgrcommon.TxUnstarted, "10000", anotherSubID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedEth(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "80000", start.String())

	// A subscriber's balance is deducted with the native reserved across multiple keys,
	// i.e, gas lanes.
	addEthTxNativePayment(t, txstore, k2.Address, txmgrcommon.TxUnstarted, "10000", subID, reqTxHash, vrfVersion)
	start, err = listener.MaybeSubtractReservedEth(ctx, big.NewInt(100_000), chainID, subID, vrfVersion)
	require.NoError(t, err)
	require.Equal(t, "70000", start.String())
}

func TestMaybeSubtractReservedNativeV2Plus(t *testing.T) {
	testMaybeSubtractReservedNative(t, vrfcommon.V2Plus)
}

func TestMaybeSubtractReservedNativeV2(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock(ctx, "blah"))
	chainID := testutils.SimulatedChainID
	subID := new(big.Int).SetUint64(1)

	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
	}).Toml())
	require.NoError(t, err)
	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	txm := makeTestTxm(t, txstore, ks)
	chain := evmmocks.NewChain(t)
	chain.On("TxManager").Return(txm).Maybe()
	listener := &listenerV2{
		respCount: map[string]uint64{},
		job:       j,
		chain:     chain,
	}
	// returns error because native payment is not supported for V2
	start, err := listener.MaybeSubtractReservedEth(testutils.Context(t), big.NewInt(100_000), chainID, subID, vrfcommon.V2)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), start)
}

func TestListener_GetConfirmedAt(t *testing.T) {
	j, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 10,
	}).Toml())
	require.NoError(t, err)

	listener := &listenerV2{
		respCount: map[string]uint64{},
		job:       j,
	}

	// Requester asks for 100 confirmations, we have a delay of 10,
	// so we should wait for max(nodeMinConfs, requestedConfs + requestedConfsDelay) = 110 confirmations
	nodeMinConfs := 10
	confirmedAt := listener.getConfirmedAt(NewV2RandomWordsRequested(&vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   big.NewInt(1),
		MinimumRequestConfirmations: 100,
		Raw: types.Log{
			BlockNumber: 100,
		},
	}), uint32(nodeMinConfs))
	require.Equal(t, uint64(210), confirmedAt) // log block number + # of confirmations

	// Requester asks for 100 confirmations, we have a delay of 0,
	// so we should wait for max(nodeMinConfs, requestedConfs + requestedConfsDelay) = 100 confirmations
	j, err = vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		RequestedConfsDelay: 0,
	}).Toml())
	require.NoError(t, err)
	listener.job = j
	confirmedAt = listener.getConfirmedAt(NewV2RandomWordsRequested(&vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
		RequestId:                   big.NewInt(1),
		MinimumRequestConfirmations: 100,
		Raw: types.Log{
			BlockNumber: 100,
		},
	}), uint32(nodeMinConfs))
	require.Equal(t, uint64(200), confirmedAt) // log block number + # of confirmations
}

func TestListener_Backoff(t *testing.T) {
	var tests = []struct {
		name     string
		initial  time.Duration
		max      time.Duration
		last     time.Duration
		retries  int
		expected bool
	}{
		{
			name:     "Backoff disabled, ready",
			expected: true,
		},
		{
			name:     "First try, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     0,
			retries:  0,
			expected: true,
		},
		{
			name:     "Second try, not ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     59 * time.Second,
			retries:  1,
			expected: false,
		},
		{
			name:     "Second try, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     61 * time.Second, // Last try was over a minute ago
			retries:  1,
			expected: true,
		},
		{
			name:     "Third try, not ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     77 * time.Second, // Slightly less than backoffFactor * initial
			retries:  2,
			expected: false,
		},
		{
			name:     "Third try, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     79 * time.Second, // Slightly more than backoffFactor * initial
			retries:  2,
			expected: true,
		},
		{
			name:     "Max, not ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     59 * time.Minute, // Slightly less than max
			retries:  900,
			expected: false,
		},
		{
			name:     "Max, ready",
			initial:  time.Minute,
			max:      time.Hour,
			last:     61 * time.Minute, // Slightly more than max
			retries:  900,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lsn := &listenerV2{job: job.Job{
				VRFSpec: &job.VRFSpec{
					BackoffInitialDelay: test.initial,
					BackoffMaxDelay:     test.max,
				},
			}}

			req := pendingRequest{
				confirmedAtBlock: 5,
				attempts:         test.retries,
				lastTry:          time.Now().Add(-test.last),
			}

			require.Equal(t, test.expected, lsn.ready(req, 10))
		})
	}
}
