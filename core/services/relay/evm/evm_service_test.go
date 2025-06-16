package evm

import (
	"errors"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
	evmtypes "github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
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

type Mocks struct {
	Chain     *evmmocks.Chain
	TxManager *txmmocks.MockEvmTxManager
	Config    *configmocks.ChainScopedConfig
	EVM       *configmocks.EVM
	Workflow  *configmocks.Workflow
	EvmClient *clienttest.Client
	Poller    *lpmocks.LogPoller
	Relayer   *Relayer
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

	relayer := &Relayer{
		chain: chain,
	}

	return &Mocks{
		Chain:     chain,
		TxManager: txManager,
		Config:    mockConfig,
		EVM:       mockEVM,
		Workflow:  mockWorkflow,
		EvmClient: evmClient,
		Poller:    poller,
	}, relayer
}

type SubmitTransactionTestCase struct {
	Name           string
	SetupMocks     func(m *Mocks, ctx any)
	ExpectedResult *evmtypes.TransactionResult
	ExpectedError  string
}

func runSubmitTransactionTest(t *testing.T, tc SubmitTransactionTestCase) {
	ctx := t.Context()
	mocks, relayer := setupMocksAndRelayer(t)

	if tc.SetupMocks != nil {
		tc.SetupMocks(mocks, ctx)
	}

	setCommonSubmitTransactionMocks(mocks, ctx)

	receiver := createToAddress()
	gasLimit := uint64(1000)
	result, err := relayer.SubmitTransaction(ctx, evmtypes.SubmitTransactionRequest{
		To:   receiver,
		Data: createPayload(),
		GasConfig: &evmtypes.GasConfig{
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

func setCommonSubmitTransactionMocks(m *Mocks, ctx any) {
	fromAddress := createFromAddress()
	m.Workflow.EXPECT().FromAddress().Return(&fromAddress)
	m.EVM.EXPECT().TxMinimumWaitTimeForConfirmation().Return(time.Millisecond)
	m.EVM.EXPECT().TxMaximumWaitTimeForConfirmation().Return(time.Millisecond)
}

func createFromAddress() types.EIP55Address {
	address, _ := types.NewEIP55Address("0x222")
	return address
}

func createToAddress() common.Address {
	return common.HexToAddress("0x555")
}

func createPayload() evm.ABIPayload {
	return evm.ABIPayload([]byte("kitties"))
}

func TestEVMService(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	t.Run("RegisterLogTracking", func(t *testing.T) {
		mocks, relayer := setupMocksAndRelayer(t)
		filter := evmtypes.LPFilterQuery{
			Name:         "filter-1",
			Retention:    time.Second,
			Addresses:    []evmtypes.Address{common.HexToAddress("0x123")},
			EventSigs:    []evmtypes.Hash{common.HexToHash("0x321")},
			Topic2:       []evmtypes.Hash{common.HexToHash("0x222")},
			Topic3:       []evmtypes.Hash{common.HexToHash("0x543")},
			Topic4:       []evmtypes.Hash{common.HexToHash("0x432")},
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
		mocks.EvmClient.On("TransactionByHash", ctx, hash).Return(transaction, nil)
		tx, err := relayer.GetTransactionByHash(ctx, hash)
		require.NoError(t, err)
		require.Equal(t, transaction.Hash().Bytes(), tx.Hash[:])
		require.Equal(t, transaction.Nonce(), tx.Nonce)
		require.Equal(t, transaction.GasPrice(), tx.GasPrice)
		require.Equal(t, transaction.Data(), tx.Data)
		require.Equal(t, transaction.Gas(), tx.Gas)
		require.Equal(t, transaction.To().Bytes(), tx.To[:])
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
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Confirmed, nil)
				txHash := common.HexToHash("0xabcd")
				mockReceipt := NewChainReceipt(txHash, t)
				m.TxManager.EXPECT().GetTransactionReceipt(ctx, mock.Anything).Return(&mockReceipt, nil)
			},
			ExpectedResult: &evmtypes.TransactionResult{
				TxHash:   common.HexToHash("0xabcd"),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Fail creating transaction",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Confirmed, nil)
				expectedMessage := "fail creating transaction"
				m.TxManager.EXPECT().GetTransactionReceipt(ctx, mock.Anything).Return(nil, errors.New(expectedMessage))
			},
			ExpectedError: "fail creating transaction",
		},
		{
			Name: "Fails getting transaction status",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				expectedMessage := "fail getting transaction status"
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Fatal, errors.New(expectedMessage))
			},
			ExpectedError: "fail getting transaction status",
		},
		{
			Name: "Success with unconfirmed status and then finalized status",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Unconfirmed, nil).Once()
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Finalized, nil).Once()
				txHash := common.HexToHash("0xabcd")
				mockReceipt := NewChainReceipt(txHash, t)
				m.TxManager.EXPECT().GetTransactionReceipt(ctx, mock.Anything).Return(&mockReceipt, nil)
			},
			ExpectedResult: &evmtypes.TransactionResult{
				TxHash:   common.HexToHash("0xabcd"),
				TxStatus: evm.TxSuccess,
			},
		},
		{
			Name: "Fails with pending and later on Fatal",
			SetupMocks: func(m *Mocks, ctx any) {
				expectedTx := txmgr.Tx{}
				m.TxManager.EXPECT().CreateTransaction(ctx, mock.Anything).Return(expectedTx, nil)
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Pending, nil).Once()
				m.TxManager.EXPECT().GetTransactionStatus(ctx, mock.Anything).Return(commontypes.Fatal, nil).Once()
			},
			ExpectedResult: &evmtypes.TransactionResult{
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
	mock.EXPECT().GetTxHash().Return(txHash)
	return mock
}
