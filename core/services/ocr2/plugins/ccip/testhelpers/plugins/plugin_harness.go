package plugintesthelpers

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type OCR2TestContract interface {
	SetOCR2Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)
}

func generateAndSetTestOCR2Config(contract OCR2TestContract, owner *bind.TransactOpts, onChainConfig []byte) (*types.Transaction, error) {
	var signers []common.Address
	var transmitters []common.Address

	for i := 0; i < 4; i++ {
		signers = append(signers, utils.RandomAddress())
		transmitters = append(transmitters, utils.RandomAddress())
	}

	return contract.SetOCR2Config(owner, signers, transmitters, 1, onChainConfig, 2, nil)
}

type CCIPPluginTestHarness struct {
	testhelpers.CCIPContracts
	Lggr logger.Logger

	SourceLP     logpoller.LogPollerTest
	DestLP       logpoller.LogPollerTest
	DestClient   evmclient.Client
	SourceClient evmclient.Client

	CommitOnchainConfig ccipconfig.CommitOnchainConfig
	ExecOnchainConfig   ccipconfig.ExecOnchainConfig
}

func (th *CCIPPluginTestHarness) CommitAndPollLogs(t *testing.T) {
	th.Source.Chain.Commit()
	th.SourceLP.PollAndSaveLogs(testutils.Context(t), th.Source.Chain.Blockchain().CurrentBlock().Number.Int64())

	th.Dest.Chain.Commit()
	th.DestLP.PollAndSaveLogs(testutils.Context(t), th.Dest.Chain.Blockchain().CurrentBlock().Number.Int64())
}

func SetupCCIPTestHarness(t *testing.T) CCIPPluginTestHarness {
	c := testhelpers.SetupCCIPContracts(t, testhelpers.SourceChainID, testhelpers.SourceChainSelector, testhelpers.DestChainID, testhelpers.DestChainSelector)

	lggr := logger.TestLogger(t)

	// db, clients and logpollers
	db := pgtest.NewSqlxDB(t)

	sourceORM := logpoller.NewORM(new(big.Int).SetUint64(c.Source.ChainID), db, lggr, pgtest.NewQConfig(true))
	var sourceLP logpoller.LogPollerTest = logpoller.NewLogPoller(
		sourceORM,
		evmclient.NewSimulatedBackendClient(t, c.Source.Chain, new(big.Int).SetUint64(c.Source.ChainID)),
		lggr.Named("sourceLP"),
		1*time.Hour, 1, 10, 10, 1000,
	)

	destORM := logpoller.NewORM(new(big.Int).SetUint64(c.Dest.ChainID), db, lggr, pgtest.NewQConfig(true))
	var destLP logpoller.LogPollerTest = logpoller.NewLogPoller(
		destORM,
		evmclient.NewSimulatedBackendClient(t, c.Dest.Chain, new(big.Int).SetUint64(c.Dest.ChainID)),
		lggr.Named("destLP"),
		1*time.Hour, 1, 10, 10, 1000,
	)

	// onChain configs
	encodedCommitOnchainConfig := c.CreateDefaultCommitOnchainConfig(t)
	commitOnchainConfig, err := abihelpers.DecodeAbiStruct[ccipconfig.CommitOnchainConfig](encodedCommitOnchainConfig)
	require.NoError(t, err)

	_, err = generateAndSetTestOCR2Config(c.Dest.CommitStore, c.Dest.User, encodedCommitOnchainConfig)
	require.NoError(t, err)

	encodedExecOnchainConfig := c.CreateDefaultExecOnchainConfig(t)
	execOnchainConfig, err := abihelpers.DecodeAbiStruct[ccipconfig.ExecOnchainConfig](encodedExecOnchainConfig)
	require.NoError(t, err)

	_, err = generateAndSetTestOCR2Config(c.Dest.OffRamp, c.Dest.User, encodedExecOnchainConfig)
	require.NoError(t, err)
	c.Dest.Chain.Commit()

	// approve router
	_, err = c.Source.LinkToken.Approve(c.Source.User, c.Source.Router.Address(), testhelpers.Link(500))
	require.NoError(t, err)
	c.Source.Chain.Commit()

	_, err = c.Dest.PriceRegistry.UpdatePrices(c.Dest.User, price_registry.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry.InternalTokenPriceUpdate{
			{SourceToken: c.Dest.LinkToken.Address(), UsdPerToken: big.NewInt(5)},
			{SourceToken: c.Dest.WrappedNative.Address(), UsdPerToken: big.NewInt(5)},
		},
		DestChainSelector: c.Dest.ChainSelector,
		UsdPerUnitGas:     big.NewInt(1),
	})
	require.NoError(t, err)

	// register filters in logPoller
	require.NoError(t, sourceLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Commit ccip sends", c.Source.OnRamp.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.SendRequested}, Addresses: []common.Address{c.Source.OnRamp.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Commit price updates", c.Dest.PriceRegistry.Address()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.UsdPerUnitGasUpdated, abihelpers.EventSignatures.UsdPerTokenUpdated}, Addresses: []common.Address{c.Dest.PriceRegistry.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Exec report accepts", c.Dest.CommitStore.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.ReportAccepted}, Addresses: []common.Address{c.Dest.CommitStore.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Exec execution state changes", c.Dest.OffRamp.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.ExecutionStateChanged}, Addresses: []common.Address{c.Dest.OffRamp.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Token pool added", c.Dest.OffRamp.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.PoolAdded}, Addresses: []common.Address{c.Dest.OffRamp.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Token pool removed", c.Dest.OffRamp.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.PoolRemoved}, Addresses: []common.Address{c.Dest.OffRamp.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Fee token added", c.Dest.PriceRegistry.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenAdded}, Addresses: []common.Address{c.Dest.PriceRegistry.Address()},
	}))
	require.NoError(t, destLP.RegisterFilter(logpoller.Filter{
		Name:      logpoller.FilterName("Fee token removed", c.Dest.PriceRegistry.Address().String()),
		EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenRemoved}, Addresses: []common.Address{c.Dest.PriceRegistry.Address()},
	}))

	// start and backfill logpollers
	require.NoError(t, sourceLP.Start(testutils.Context(t)))
	require.NoError(t, destLP.Start(testutils.Context(t)))
	require.NoError(t, sourceLP.Replay(testutils.Context(t), 1))
	require.NoError(t, destLP.Replay(testutils.Context(t), 1))

	th := CCIPPluginTestHarness{
		CCIPContracts: c,
		Lggr:          lggr,

		SourceLP:            sourceLP,
		DestLP:              destLP,
		DestClient:          evmclient.NewSimulatedBackendClient(t, c.Dest.Chain, new(big.Int).SetUint64(c.Dest.ChainID)),
		SourceClient:        evmclient.NewSimulatedBackendClient(t, c.Source.Chain, new(big.Int).SetUint64(c.Source.ChainID)),
		CommitOnchainConfig: commitOnchainConfig,
		ExecOnchainConfig:   execOnchainConfig,
	}

	th.CommitAndPollLogs(t)
	return th
}

type MessageBatch struct {
	TokenData [][][]byte
	Messages  []evm_2_evm_offramp.InternalEVM2EVMMessage
	Interval  commit_store.CommitStoreInterval
	Root      [32]byte
	Proof     merklemulti.Proof[[32]byte]
	ProofBits *big.Int
	Tree      *merklemulti.Tree[[32]byte]
}

func (mb MessageBatch) ToExecutionReport() evm_2_evm_offramp.InternalExecutionReport {
	return evm_2_evm_offramp.InternalExecutionReport{
		Messages:          mb.Messages,
		OffchainTokenData: mb.TokenData,
		Proofs:            mb.Proof.Hashes,
		ProofFlagBits:     mb.ProofBits,
	}
}

func (th *CCIPPluginTestHarness) GenerateAndSendMessageBatch(t *testing.T, nMessages int, payloadSize int, nTokensPerMessage int) MessageBatch {
	mctx := hashlib.NewKeccakCtx()
	leafHasher := hashlib.NewLeafHasher(th.Source.ChainSelector, th.Dest.ChainSelector, th.Source.OnRamp.Address(), mctx)

	maxPayload := make([]byte, payloadSize)
	for i := 0; i < payloadSize; i++ {
		maxPayload[i] = 0xa
	}

	var offchainTokenData [][]byte
	var tokenAmounts []router.ClientEVMTokenAmount
	for i := 0; i < nTokensPerMessage; i++ {
		tokenAmounts = append(tokenAmounts, router.ClientEVMTokenAmount{
			Token:  th.Source.LinkToken.Address(),
			Amount: big.NewInt(int64(1 + i)),
		})
		offchainTokenData = append(offchainTokenData, []byte{})
	}

	th.Source.Chain.Commit()
	startBlock := th.Source.Chain.Blockchain().CurrentBlock().Number
	var lastFlush int
	for i := 0; i < nMessages; i++ {
		routerMsg := router.ClientEVM2AnyMessage{
			Receiver:     testhelpers.MustEncodeAddress(t, th.Dest.Receivers[0].Receiver.Address()),
			FeeToken:     th.Source.LinkToken.Address(),
			TokenAmounts: tokenAmounts,
			Data:         maxPayload,
			ExtraArgs:    []byte{},
		}
		_, err := th.Source.Router.CcipSend(th.Source.User, th.Dest.ChainSelector, routerMsg)
		require.NoError(t, err)
		lastFlush++
		if lastFlush*payloadSize > 700_000 {
			th.CommitAndPollLogs(t)
			lastFlush = 0
		}
	}
	th.CommitAndPollLogs(t)

	leafHashes := make([][32]byte, nMessages)
	tokenData := make([][][]byte, nMessages)
	indices := make([]int, nMessages)
	messages := make([]evm_2_evm_offramp.InternalEVM2EVMMessage, nMessages)
	seqNums := make([]uint64, nMessages)

	sendEvents, err := th.Source.OnRamp.FilterCCIPSendRequested(&bind.FilterOpts{Start: startBlock.Uint64(), Context: testutils.Context(t)})
	require.NoError(t, err)
	var i int
	for ; sendEvents.Next(); i++ {
		indices[i] = i
		tokenData[i] = offchainTokenData
		messages[i] = abihelpers.OnRampMessageToOffRampMessage(sendEvents.Event.Message)
		leafHash, err2 := leafHasher.HashLeaf(sendEvents.Event.Raw)
		require.NoError(t, err2)
		leafHashes[i] = leafHash
		seqNums[i] = sendEvents.Event.Message.SequenceNumber
	}
	require.Equal(t, nMessages, i)

	tree, err := merklemulti.NewTree(mctx, leafHashes)
	require.NoError(t, err)
	proof, err := tree.Prove(indices)
	require.NoError(t, err)
	root := tree.Root()
	rootLocal, err := merklemulti.VerifyComputeRoot(mctx, leafHashes, proof)
	require.NoError(t, err)
	require.Equal(t, root, rootLocal)

	return MessageBatch{
		Messages:  messages,
		Interval:  commit_store.CommitStoreInterval{Min: seqNums[0], Max: seqNums[len(seqNums)-1]},
		TokenData: tokenData,
		Root:      root,
		Proof:     proof,
		ProofBits: abihelpers.ProofFlagsToBits(proof.SourceFlags),
		Tree:      tree,
	}
}
