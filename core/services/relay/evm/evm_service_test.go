package evm

import (
	"errors"
	"math"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	configmocks "github.com/smartcontractkit/chainlink-evm/pkg/config/mocks"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"

	evmmocks "github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

const ExpectedTxHash = "0xabcd"

type Mocks struct {
	Chain         *evmmocks.Chain
	TxManager     *txmmocks.MockEvmTxManager
	Config        *configmocks.ChainScopedConfig
	EVM           *configmocks.EVM
	Workflow      *configmocks.Workflow
	EvmClient     *clienttest.Client
	Poller        *lpmocks.LogPoller
	HeaderTracker *headstest.Tracker[*types.Head, common.Hash]
	Relayer       *Relayer
}

type returnedStatusAndReceipts struct {
	Status   []commontypes.TransactionStatus
	Receipts []receiptResult
}

type receiptResult struct {
	Receipt *txmgr.ChainReceipt
	Error   error
}

func createMockReceipt(t *testing.T) *txmgr.ChainReceipt {
	receipt := NewChainReceipt(common.HexToHash(ExpectedTxHash), t)
	return &receipt
}

func setupMocksAndRelayer(t *testing.T) (*Mocks, *Relayer) {
	chain := evmmocks.NewChain(t)
	txManager := txmmocks.NewMockEvmTxManager(t)
	mockConfig := configmocks.NewChainScopedConfig(t)
	mockEVM := configmocks.NewEVM(t)
	mockWorkflow := configmocks.NewWorkflow(t)
	evmClient := clienttest.NewClient(t)
	poller := lpmocks.NewLogPoller(t)
	ht := headstest.NewTracker[*types.Head](t)

	chain.On("TxManager").Return(txManager).Maybe()
	chain.On("LogPoller").Return(poller).Maybe()
	chain.On("HeadTracker").Return(ht).Maybe()
	chain.On("Client").Return(evmClient).Maybe()
	chain.EXPECT().Config().Return(mockConfig).Maybe()
	mockConfig.EXPECT().EVM().Return(mockEVM).Maybe()
	mockEVM.EXPECT().Workflow().Return(mockWorkflow).Maybe()

	lggr, err := logger.New()
	require.NoError(t, err)
	relayer := &Relayer{
		chain:      chain,
		evmService: evmService{chain: chain, logger: lggr},
	}

	return &Mocks{
		Chain:         chain,
		TxManager:     txManager,
		Config:        mockConfig,
		EVM:           mockEVM,
		Workflow:      mockWorkflow,
		EvmClient:     evmClient,
		Poller:        poller,
		HeaderTracker: ht,
	}, relayer
}

type SubmitTransactionTestCase struct {
	Name           string
	SetupMocks     func(m *Mocks, ctx any)
	ExpectedResult *evm.TransactionResult
	ExpectedError  string
}

func runSubmitTransactionTest(t *testing.T, tc SubmitTransactionTestCase) {
	ctx := t.Context()
	mocks, relayer := setupMocksAndRelayer(t)

	if tc.SetupMocks != nil {
		tc.SetupMocks(mocks, ctx)
	}

	setCommonSubmitTransactionMocks(mocks)

	receiver := createToAddress()
	gasLimit := uint64(1000)
	result, err := relayer.SubmitTransaction(ctx, evm.SubmitTransactionRequest{
		To:   receiver,
		Data: createPayload(),
		GasConfig: &evm.GasConfig{
			GasLimit: &gasLimit,
		},
	})

	if tc.ExpectedError != "" {
		require.Error(t, err)
		require.Contains(t, err.Error(), tc.ExpectedError)
	} else {
		require.NoError(t, err)
		require.Equal(t, tc.ExpectedResult, result)
	}
}

func setCommonSubmitTransactionMocks(m *Mocks) {
	fromAddress := createFromAddress()
	m.Workflow.EXPECT().FromAddress().Return(&fromAddress)
	m.EVM.EXPECT().ConfirmationTimeout().Return(2 * time.Second)
}

func createFromAddress() types.EIP55Address {
	address, _ := types.NewEIP55Address("0x222")
	return address
}

func createToAddress() common.Address {
	return common.HexToAddress("0x555")
}

func createPayload() evm.ABIPayload {
	return evm.ABIPayload("kitties")
}

func TestEVMService(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	t.Run("RegisterLogTracking", func(t *testing.T) {
		mocks, relayer := setupMocksAndRelayer(t)
		filter := evm.LPFilterQuery{
			Name:         "filter-1",
			Retention:    time.Second,
			Addresses:    []evm.Address{common.HexToAddress("0x123")},
			EventSigs:    []evm.Hash{common.HexToHash("0x321")},
			Topic2:       []evm.Hash{common.HexToHash("0x222")},
			Topic3:       []evm.Hash{common.HexToHash("0x543")},
			Topic4:       []evm.Hash{common.HexToHash("0x432")},
			MaxLogsKept:  100,
			LogsPerBlock: 10,
		}

		mocks.Poller.On("HasFilter", mock.MatchedBy(func(fname string) bool {
			return fname == filter.Name
		})).Return(false)
		mocks.Poller.On("RegisterFilter", ctx, mock.MatchedBy(func(f logpoller.Filter) bool {
			return f.LogsPerBlock == filter.LogsPerBlock &&
				f.Retention == filter.Retention &&
				f.Topic2[0] == filter.Topic2[0] &&
				f.Topic3[0] == filter.Topic3[0] &&
				f.Topic4[0] == filter.Topic4[0] &&
				f.EventSigs[0] == filter.EventSigs[0] &&
				f.MaxLogsKept == filter.MaxLogsKept &&
				f.Addresses[0] == filter.Addresses[0] &&
				f.Name == filter.Name
		})).Return(nil)

		err := relayer.RegisterLogTracking(ctx, filter)
		require.NoError(t, err)
	})

	t.Run("GetTransactionByHash", func(t *testing.T) {
		mocks, relayer := setupMocksAndRelayer(t)

		hash := common.HexToHash("0x123")
		nonce := uint64(1)
		to := common.HexToAddress("0x555")
		amount := big.NewInt(1)
		gasLimit := uint64(2)
		gasPrice := big.NewInt(2)
		data := []byte("kitties")

		transaction := gethtypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
		mocks.EvmClient.EXPECT().TransactionByHashWithOpts(ctx, hash, types.TransactionByHashOpts{}).Return(transaction, nil)
		tx, err := relayer.GetTransactionByHash(ctx, evm.GetTransactionByHashRequest{Hash: hash})
		require.NoError(t, err)
		require.Equal(t, transaction.Hash().Bytes(), tx.Hash[:])
		require.Equal(t, transaction.Nonce(), tx.Nonce)
		require.Equal(t, transaction.GasPrice(), tx.GasPrice)
		require.Equal(t, transaction.Data(), tx.Data)
		require.Equal(t, transaction.Gas(), tx.Gas)
		require.Equal(t, transaction.To().Bytes(), tx.To[:])
	})

	t.Run("GetFiltersNames", func(t *testing.T) {
		// TODO PLEX-1465: once code is moved away, remove this test
		mocks, relayer := setupMocksAndRelayer(t)
		filtersMap := map[string]logpoller.Filter{
			"filterA": {},
			"filterB": {},
		}
		mocks.Poller.On("GetFilters").Return(filtersMap)
		names, _ := relayer.GetFiltersNames(ctx)
		require.ElementsMatch(t, []string{"filterA", "filterB"}, names)
	})

	submitTxCases := []SubmitTransactionTestCase{
		{
			Name: "Executes successfully",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTxRequest := txmgr.TxRequest{
					FromAddress:    createFromAddress().Address(),
					ToAddress:      createToAddress(),
					EncodedPayload: createPayload(),
				}
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.MatchedBy(func(txRequest txmgr.TxRequest) bool {
					return txRequest.FromAddress == expectedTxRequest.FromAddress &&
						txRequest.ToAddress == expectedTxRequest.ToAddress &&
						slices.Equal(txRequest.EncodedPayload, expectedTxRequest.EncodedPayload)
				})).Return(expectedTx, nil)
				m.TxManager.EXPECT().GetTransactionStatus(mock.Anything, mock.Anything).Return(commontypes.Unconfirmed, nil)
				txHash := common.HexToHash(ExpectedTxHash)
				mockReceipt := NewChainReceipt(txHash, t)
				m.TxManager.EXPECT().GetTransactionReceipt(mock.Anything, mock.Anything).Return(&mockReceipt, nil)
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.HexToHash(ExpectedTxHash),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Fails creating transaction",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				m.TxManager.EXPECT().GetTransactionStatus(mock.Anything, mock.Anything).Return(commontypes.Unconfirmed, nil)
				expectedMessage := "fail creating transaction"
				m.TxManager.EXPECT().GetTransactionReceipt(mock.Anything, mock.Anything).Return(nil, errors.New(expectedMessage))
			},
			ExpectedError: "getting transaction receipt",
		},
		{
			Name: "Fails getting transaction status",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				expectedMessage := "fail getting transaction status"
				m.TxManager.EXPECT().GetTransactionStatus(mock.Anything, mock.Anything).Return(commontypes.Fatal, errors.New(expectedMessage))
			},
			ExpectedError: "failed getting transaction status",
		},
		{
			Name: "Success with pending status and then finalized status",
			SetupMocks: func(m *Mocks, ctx any) {
				runSubmitTxGettingDifferentStatusAndReceipts(m, ctx, returnedStatusAndReceipts{
					Status:   []commontypes.TransactionStatus{commontypes.Pending, commontypes.Finalized},
					Receipts: []receiptResult{{Receipt: createMockReceipt(t), Error: nil}}})
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.HexToHash(ExpectedTxHash),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Success with unknown status and then finalized status",
			SetupMocks: func(m *Mocks, ctx any) {
				runSubmitTxGettingDifferentStatusAndReceipts(m, ctx, returnedStatusAndReceipts{
					Status:   []commontypes.TransactionStatus{commontypes.Unknown, commontypes.Finalized},
					Receipts: []receiptResult{{Receipt: createMockReceipt(t), Error: nil}}})
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.HexToHash(ExpectedTxHash),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Success with unknown status and then unconfirmed status",
			SetupMocks: func(m *Mocks, ctx any) {
				runSubmitTxGettingDifferentStatusAndReceipts(m, ctx, returnedStatusAndReceipts{
					Status:   []commontypes.TransactionStatus{commontypes.Unknown, commontypes.Unconfirmed},
					Receipts: []receiptResult{{Receipt: createMockReceipt(t), Error: nil}}})
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.HexToHash(ExpectedTxHash),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Success with unknown status and then unconfirmed status and failed get receipt attempt with null receipt",
			SetupMocks: func(m *Mocks, ctx any) {
				runSubmitTxGettingDifferentStatusAndReceipts(m, ctx, returnedStatusAndReceipts{
					Status:   []commontypes.TransactionStatus{commontypes.Unknown, commontypes.Unconfirmed},
					Receipts: []receiptResult{{Receipt: nil, Error: nil}, {Receipt: createMockReceipt(t), Error: nil}}})
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.HexToHash(ExpectedTxHash),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Success with unknown status and then finalized status and failed get receipt attempt with error",
			SetupMocks: func(m *Mocks, ctx any) {
				runSubmitTxGettingDifferentStatusAndReceipts(m, ctx, returnedStatusAndReceipts{
					Status:   []commontypes.TransactionStatus{commontypes.Unknown, commontypes.Finalized},
					Receipts: []receiptResult{{Receipt: nil, Error: errors.New("Some error")}, {Receipt: createMockReceipt(t), Error: nil}}})
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.HexToHash(ExpectedTxHash),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Fails with pending and later on Fatal",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				m.TxManager.EXPECT().GetTransactionStatus(mock.Anything, mock.Anything).Return(commontypes.Pending, nil).Once()
				m.TxManager.EXPECT().GetTransactionStatus(mock.Anything, mock.Anything).Return(commontypes.Fatal, nil).Once()
			},
			ExpectedResult: &evm.TransactionResult{
				TxHash:   common.Hash{},
				TxStatus: evm.TxFatal,
			},
		},
	}

	for _, tc := range submitTxCases {
		t.Run("SubmitTransaction - "+tc.Name, func(t *testing.T) {
			runSubmitTransactionTest(t, tc)
		})
	}
}

func TestEVMService_HeaderByNumber(t *testing.T) {
	testCases := []struct {
		Name           string
		Request        evm.HeaderByNumberRequest
		ExpectedResult evm.HeaderByNumberReply
		PrepareMocks   func(m *Mocks)
		ExpectedError  string
	}{
		{
			Name: "Explicit Latest header",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(rpc.LatestBlockNumber.Int64()),
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestAndFinalizedBlock(mock.Anything).Return(&types.Head{Number: 10}, &types.Head{Number: 8}, nil).Once()
			},
			ExpectedResult: evm.HeaderByNumberReply{Header: &evm.Header{Number: big.NewInt(10)}},
		},
		{
			Name: "Nil BlockNumber - should return latest header",
			Request: evm.HeaderByNumberRequest{
				Number: nil,
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestAndFinalizedBlock(mock.Anything).Return(&types.Head{Number: 10}, &types.Head{Number: 8}, nil).Once()
			},
			ExpectedResult: evm.HeaderByNumberReply{Header: &evm.Header{Number: big.NewInt(10)}},
		},
		{
			Name: "Finalized",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(rpc.FinalizedBlockNumber.Int64()),
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestAndFinalizedBlock(mock.Anything).Return(&types.Head{Number: 10}, &types.Head{Number: 8}, nil).Once()
			},
			ExpectedResult: evm.HeaderByNumberReply{Header: &evm.Header{Number: big.NewInt(8)}},
		},
		{
			Name: "Safe",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(rpc.SafeBlockNumber.Int64()),
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestSafeBlock(mock.Anything).Return(&types.Head{Number: 9}, nil).Once()
			},
			ExpectedResult: evm.HeaderByNumberReply{Header: &evm.Header{Number: big.NewInt(9)}},
		},
		{
			Name: "Unknown special block number",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(-42),
			},
			ExpectedError: "unexpected block number -42",
		},
		{
			Name: "Non-special block number",
			Request: evm.HeaderByNumberRequest{
				Number:          big.NewInt(42),
				ConfidenceLevel: primitives.Finalized,
			},
			PrepareMocks: func(m *Mocks) {
				m.EvmClient.EXPECT().HeaderByNumberWithOpts(
					mock.Anything,
					big.NewInt(42),
					types.HeaderByNumberOpts{ConfidenceLevel: primitives.Finalized},
				).Return(&types.Header{Number: 42}, nil).Once()
			},
			ExpectedResult: evm.HeaderByNumberReply{Header: &evm.Header{Number: big.NewInt(42)}},
		},
		{
			Name: "Large block number",
			Request: evm.HeaderByNumberRequest{
				Number:          big.NewInt(0).SetUint64(math.MaxInt64 + 1),
				ConfidenceLevel: primitives.Finalized,
			},
			ExpectedError: "block number 9223372036854775808 is larger than int64: not found",
		},
		{
			Name: "Failed to get latest",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(rpc.LatestBlockNumber.Int64()),
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestAndFinalizedBlock(mock.Anything).Return(nil, nil, errors.New("failed to get latest")).Once()
			},
			ExpectedError: "failed to get latest",
		},
		{
			Name: "Failed to get finalized",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(rpc.FinalizedBlockNumber.Int64()),
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestAndFinalizedBlock(mock.Anything).Return(nil, nil, errors.New("failed to get finalized")).Once()
			},
			ExpectedError: "failed to get finalized",
		},
		{
			Name: "Safe",
			Request: evm.HeaderByNumberRequest{
				Number: big.NewInt(rpc.SafeBlockNumber.Int64()),
			},
			PrepareMocks: func(m *Mocks) {
				m.HeaderTracker.EXPECT().LatestSafeBlock(mock.Anything).Return(nil, errors.New("failed to get safe")).Once()
			},
			ExpectedError: "failed to get safe",
		},
		{
			Name: "Failed to get non-special block number",
			Request: evm.HeaderByNumberRequest{
				Number:          big.NewInt(42),
				ConfidenceLevel: primitives.Finalized,
			},
			PrepareMocks: func(m *Mocks) {
				m.EvmClient.EXPECT().HeaderByNumberWithOpts(
					mock.Anything,
					big.NewInt(42),
					types.HeaderByNumberOpts{ConfidenceLevel: primitives.Finalized},
				).Return(nil, errors.New("failed to get block 42")).Once()
			},
			ExpectedError: "failed to get block 42",
		},
		{
			Name: "Block not found",
			Request: evm.HeaderByNumberRequest{
				Number:          big.NewInt(404),
				ConfidenceLevel: primitives.Finalized,
			},
			PrepareMocks: func(m *Mocks) {
				m.EvmClient.EXPECT().HeaderByNumberWithOpts(
					mock.Anything,
					big.NewInt(404),
					types.HeaderByNumberOpts{ConfidenceLevel: primitives.Finalized},
				).Return(nil, nil).Once()
			},
			ExpectedError: ethereum.NotFound.Error(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mocks, relayer := setupMocksAndRelayer(t)

			if tc.PrepareMocks != nil {
				tc.PrepareMocks(mocks)
			}

			result, err := relayer.HeaderByNumber(t.Context(), tc.Request)

			if tc.ExpectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.ExpectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedResult.Header.Number, result.Header.Number)
			}
		})
	}
}

func runSubmitTxGettingDifferentStatusAndReceipts(m *Mocks, ctx any, expectedReturns returnedStatusAndReceipts) {
	expectedTx := txmgr.Tx{}
	m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
	for _, status := range expectedReturns.Status {
		m.TxManager.EXPECT().GetTransactionStatus(mock.Anything, mock.Anything).Return(status, nil).Once()
	}
	for _, receiptResult := range expectedReturns.Receipts {
		m.TxManager.EXPECT().GetTransactionReceipt(mock.Anything, mock.Anything).Return(receiptResult.Receipt, receiptResult.Error).Once()
	}
}

func TestConverters(t *testing.T) {
	t.Parallel()

	t.Run("convert head", func(t *testing.T) {
		head := types.Head{
			Timestamp: time.Unix(100000, 100),
			Number:    100,
			Hash:      common.HexToHash("0x123"),
		}
		result := convertHead(&head)
		require.Equal(t, head.Hash.Bytes(), result.Hash[:])
	})

	t.Run("convert transaction", func(t *testing.T) {
		tx := gethtypes.NewTransaction(
			1,
			common.HexToAddress("0xabc123"),
			big.NewInt(1000),
			21000,
			big.NewInt(1e9),
			[]byte{1, 2, 3},
		)

		result := convertTransaction(tx)
		require.NotNil(t, result)
		require.Equal(t, tx.Hash().Bytes(), result.Hash[:])
		require.Equal(t, tx.Nonce(), result.Nonce)
		require.Equal(t, tx.Gas(), result.Gas)
		require.Equal(t, tx.GasPrice(), result.GasPrice)
		require.Equal(t, tx.Value(), result.Value)
		require.Equal(t, tx.To().Bytes(), result.To[:])
		require.Equal(t, tx.Data(), result.Data)
	})
}

func NewChainReceipt(txHash common.Hash, t *testing.T) txmgr.ChainReceipt {
	mock := txmmocks.NewChainReceipt[common.Hash, common.Hash](t)
	mock.EXPECT().GetTxHash().Return(txHash).Maybe()
	return mock
}
