package v2

import (
	"context"
	stderrors "errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/vrfkey"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestSkipRevertedTxnFetchFatal(t *testing.T) {
	t.Parallel()

	ctxAlive := context.Background()
	ctxCanceled, cancel := context.WithCancel(context.Background())
	cancel()

	dbErr := stderrors.New("connection reset")

	tests := []struct {
		name  string
		ctxFn func() context.Context
		err   error
		want  bool
	}{
		{
			name: "nil error",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  nil,
			want: false,
		},
		{
			name: "context.Canceled on alive ctx",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  context.Canceled,
			want: true,
		},
		{
			name: "wrapped context.Canceled",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  errors.Wrap(context.Canceled, "pq"),
			want: true,
		},
		{
			name: "context.DeadlineExceeded on alive ctx inner query timeout",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "wrapped DeadlineExceeded on alive ctx",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  errors.Wrap(context.DeadlineExceeded, "timeout"),
			want: false,
		},
		{
			name: "DeadlineExceeded while outer ctx canceled",
			ctxFn: func() context.Context {
				return ctxCanceled
			},
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			name: "generic DB error on alive ctx",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  dbErr,
			want: false,
		},
		{
			name: "generic DB error while outer ctx canceled",
			ctxFn: func() context.Context {
				return ctxCanceled
			},
			err:  dbErr,
			want: true,
		},
		{
			name: "context.Canceled while outer ctx also canceled",
			ctxFn: func() context.Context {
				return ctxCanceled
			},
			err:  context.Canceled,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := skipRevertedTxnFetchFatal(tt.ctxFn(), tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilterBatchRevertedTxnSkipsLogsMissingTopics(t *testing.T) {
	ctx := context.Background()

	encodedPayload, err := batchCoordinatorV2ABI.Pack("fulfillRandomWords",
		[]vrf_coordinator_v2.VRFProof{},
		[]vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment{},
	)
	assert.NoError(t, err)

	// A consumer callback emitting log0/log1 from inside a batch fulfillment
	// can put logs with fewer than two topics into the batch receipt. Such logs
	// cannot carry a RawErrorReturned event or be matched to a request ID, so
	// they must be skipped.
	receipt := evmtypes.Receipt{
		Status: 1,
		Logs: []*evmtypes.Log{
			{Address: common.HexToAddress("0xdead"), Topics: []common.Hash{}, Data: []byte{}},
			{Address: common.HexToAddress("0xdead"), Topics: []common.Hash{
				batchCoordinatorV2ABI.Events["RawErrorReturned"].ID,
			}, Data: []byte{}},
			{Address: common.HexToAddress("0xdead"), Topics: []common.Hash{
				batchCoordinatorV2ABI.Events["RandomWordsRequested"].ID,
				common.HexToHash("0x01"),
			}, Data: []byte{}},
		},
	}

	vrfKey := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	lsn := &listenerV2{
		l:      commonlogger.TestSugared(t),
		job:    job.Job{VRFSpec: &job.VRFSpec{PublicKey: vrfKey.PublicKey}},
		chStop: make(chan struct{}),
	}

	revertedTxns, err := lsn.filterBatchRevertedTxn(ctx, TxnReceiptDB{
		TxHash:         common.HexToHash("0x01"),
		EVMReceipt:     receipt,
		EncodedPayload: encodedPayload,
		SubID:          1,
	})
	assert.NoError(t, err)
	assert.Empty(t, revertedTxns)
}

func TestFilterSingleRevertedTxnMalformedRevertData(t *testing.T) {
	ctx := context.Background()

	vrfKey := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	lsn := &listenerV2{
		l:           commonlogger.TestSugared(t),
		coordinator: &stubCoordinator{commitment: [32]byte{1}},
		chain:       &stubChain{client: &stubRPCClient{}},
		job:         job.Job{VRFSpec: &job.VRFSpec{PublicKey: vrfKey.PublicKey}},
		chStop:      make(chan struct{}),
	}

	txn := TxnReceiptDB{
		TxHash:      common.HexToHash("0x02"),
		RequestID:   "0x01",
		FromAddress: common.HexToAddress("0x1111111111111111111111111111111111111111"),
		ToAddress:   common.HexToAddress("0x2222222222222222222222222222222222222222"),
		GasLimit:    100000,
		EncodedPayload: func() []byte {
			one := big.NewInt(1)
			proof := vrf_coordinator_v2.VRFProof{
				Pk:            [2]*big.Int{one, one},
				Gamma:         [2]*big.Int{one, one},
				C:             one,
				S:             one,
				Seed:          one,
				CGammaWitness: [2]*big.Int{one, one},
				SHashWitness:  [2]*big.Int{one, one},
				ZInv:          one,
			}
			commitment := vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment{
				BlockNum:         1,
				SubId:            1,
				CallbackGasLimit: 100000,
				NumWords:         1,
				Sender:           common.HexToAddress("0x1111111111111111111111111111111111111111"),
			}
			payload, err := coordinatorV2ABI.Pack("fulfillRandomWords", proof, commitment)
			assert.NoError(t, err)
			return payload
		}(),
		EVMReceipt: evmtypes.Receipt{BlockNumber: big.NewInt(1)},
	}

	tests := []struct {
		name    string
		data    any
		wantNil bool
	}{
		{
			// Some RPC nodes return the data field as a non-string value; the
			// revert reason simply cannot be determined in that case, so the
			// txn is treated like any other undetermined revert reason.
			name:    "non-string data field",
			data:    42,
			wantNil: false,
		},
		{
			// Data that decodes to fewer than 4 bytes cannot match the
			// InsufficientBalance selector, so the txn is skipped.
			name:    "data shorter than a selector",
			data:    "0x01",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lsn.chain.(*stubChain).client.(*stubRPCClient).err = evmclient.JsonError{
				Code:    -32000,
				Message: "execution reverted",
				Data:    tt.data,
			}
			revertedTxn, err := lsn.filterSingleRevertedTxn(ctx, txn)
			assert.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, revertedTxn)
			} else {
				assert.NotNil(t, revertedTxn)
			}
		})
	}
}

type stubCoordinator struct {
	CoordinatorV2_X
	commitment [32]byte
}

func (s *stubCoordinator) GetCommitment(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
	return s.commitment, nil
}

type stubChain struct {
	legacyevm.Chain
	client evmclient.Client
}

func (s *stubChain) Client() evmclient.Client { return s.client }

type stubRPCClient struct {
	evmclient.Client
	err error
}

func (s *stubRPCClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*gethtypes.Transaction, error) {
	return gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce: 0, GasPrice: big.NewInt(1), Gas: 100000, Value: big.NewInt(0),
		Data: []byte{0xde, 0xad, 0xbe, 0xef},
	}), nil
}

func (s *stubRPCClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, s.err
}
