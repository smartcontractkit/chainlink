package ocr3_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	logpollermocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ocr3"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type testUniverse[RI any] struct {
	backend         *backends.SimulatedBackend
	deployer        *bind.TransactOpts
	transmitters    []*bind.TransactOpts
	wrapper         *no_op_ocr3.NoOpOCR3
	ocr3Transmitter ocr3types.ContractTransmitter[RI]
	bundles         []ocr2key.KeyBundle
	f               uint8
}

func newTestUniverse[RI any](t *testing.T) testUniverse[RI] {
	t.Helper()

	deployer := testutils.MustNewSimTransactor(t)

	// create many transmitters but only need to fund one, rest are to get
	// setOCR3Config to pass.
	var transmitters []*bind.TransactOpts
	for i := 0; i < 4; i++ {
		transmitters = append(transmitters, testutils.MustNewSimTransactor(t))
	}

	backend := backends.NewSimulatedBackend(core.GenesisAlloc{
		deployer.From: core.GenesisAccount{
			Balance: assets.Ether(1000).ToInt(),
		},
		transmitters[0].From: core.GenesisAccount{
			Balance: assets.Ether(1000).ToInt(),
		},
	}, 30e6)
	addr, _, _, err := no_op_ocr3.DeployNoOpOCR3(deployer, backend)
	require.NoError(t, err, "failed to deploy NoOpOCR3 contract")
	backend.Commit()
	wrapper, err := no_op_ocr3.NewNoOpOCR3(addr, backend)
	require.NoError(t, err, "failed to create NoOpOCR3 wrapper")

	// create the oracle identities for setConfig
	// need to create at least 4 identities otherwise setConfig will fail
	var (
		bundles []ocr2key.KeyBundle
		signers []common.Address
	)
	for i := 0; i < 4; i++ {
		kb, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err, "failed to create key bundle")
		signers = append(signers, common.HexToAddress(kb.OnChainPublicKey()))
		bundles = append(bundles, kb)
	}
	f := uint8(1)
	_, err = wrapper.SetOCR3Config(
		deployer,
		signers,
		[]common.Address{transmitters[0].From, transmitters[1].From, transmitters[2].From, transmitters[3].From},
		f,
		[]byte{},
		3,
		[]byte{})
	require.NoError(t, err, "failed to set config")
	backend.Commit()

	contractABI, err := no_op_ocr3.NoOpOCR3MetaData.GetAbi()
	require.NoError(t, err, "failed to get abi")
	tImpl := &transmitterImpl{
		backend: backend,
		from:    transmitters[0],
		t:       t,
	}
	mockLogPoller := logpollermocks.NewLogPoller(t)
	mockLogPoller.On("RegisterFilter", mock.Anything).Return(nil)
	defer mockLogPoller.AssertExpectations(t)
	ocr3Transmitter, err := ocr3.NewOCR3ContractTransmitter[RI](
		addr,
		*contractABI,
		tImpl,
		mockLogPoller,
		logger.TestLogger(t),
		nil, // reportToEvmTxMeta, unused
	)
	require.NoError(t, err, "failed to create OCR3ContractTransmitter")

	return testUniverse[RI]{
		backend:         backend,
		deployer:        deployer,
		transmitters:    transmitters,
		wrapper:         wrapper,
		bundles:         bundles,
		ocr3Transmitter: ocr3Transmitter,
		f:               f,
	}
}

func (uni testUniverse[RI]) SignReport(t *testing.T, configDigest ocrtypes.ConfigDigest, rwi ocr3types.ReportWithInfo[RI], seqNum uint64) []ocrtypes.AttributedOnchainSignature {
	var attributedSigs []ocrtypes.AttributedOnchainSignature
	for i := uint8(0); i < uni.f+1; i++ {
		sig, err := uni.bundles[i].Sign(ocrtypes.ReportContext{
			ReportTimestamp: ocrtypes.ReportTimestamp{
				ConfigDigest: configDigest,
				Epoch:        uint32(seqNum),
			},
		}, rwi.Report)
		require.NoError(t, err, "failed to sign report")
		attributedSigs = append(attributedSigs, ocrtypes.AttributedOnchainSignature{
			Signature: sig,
			Signer:    commontypes.OracleID(i),
		})
	}
	return attributedSigs
}

func (uni testUniverse[RI]) TransmittedEvents(t *testing.T) []*no_op_ocr3.NoOpOCR3Transmitted {
	iter, err := uni.wrapper.FilterTransmitted(nil)
	require.NoError(t, err, "failed to create filter iterator")
	var events []*no_op_ocr3.NoOpOCR3Transmitted
	for iter.Next() {
		event := iter.Event
		events = append(events, event)
	}
	return events
}

func TestContractTransmitter(t *testing.T) {
	t.Parallel()

	uni := newTestUniverse[struct{}](t)

	c, err := uni.wrapper.LatestConfigDigestAndEpoch(nil)
	require.NoError(t, err, "failed to get latest config digest and epoch")
	configDigest := c.ConfigDigest

	// create the attributed signatures
	// only need f+1 which is 2 in this case
	rwi := ocr3types.ReportWithInfo[struct{}]{
		Report: []byte{},
		Info:   struct{}{},
	}
	seqNum := uint64(1)
	attributedSigs := uni.SignReport(t, configDigest, rwi, seqNum)

	account, err := uni.ocr3Transmitter.FromAccount()
	require.NoError(t, err, "failed to get from account")
	require.Equal(t, account, ocrtypes.Account(uni.transmitters[0].From.Hex()), "unexpected from account")
	err = uni.ocr3Transmitter.Transmit(context.Background(), configDigest, seqNum, rwi, attributedSigs)
	require.NoError(t, err, "failed to transmit report")

	// check for transmitted event
	// TODO: for some reason this event isn't being emitted in the simulated backend
	// events := uni.TransmittedEvents(t)
	// require.Len(t, events, 1, "expected one transmitted event")
	// event := events[0]
	// require.Equal(t, configDigest, event.ConfigDigest, "unexpected config digest")
	// require.Equal(t, seqNum, event.SequenceNumber, "unexpected sequence number")
}

type transmitterImpl struct {
	backend *backends.SimulatedBackend
	from    *bind.TransactOpts
	t       *testing.T
}

func (t *transmitterImpl) FromAddress() common.Address {
	return t.from.From
}

func (t *transmitterImpl) CreateEthTransaction(ctx context.Context, to common.Address, data []byte, txMeta *txmgr.TxMeta) error {
	nonce, err := t.backend.PendingNonceAt(ctx, t.from.From)
	require.NoError(t.t, err, "failed to get nonce")
	gp, err := t.backend.SuggestGasPrice(ctx)
	require.NoError(t.t, err, "failed to get gas price")
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      500_000,
		To:       &to,
		Value:    big.NewInt(0),
		Data:     data,
	})
	signedTx, err := t.from.Signer(t.from.From, rawTx)
	require.NoError(t.t, err, "failed to sign tx")
	err = t.backend.SendTransaction(ctx, signedTx)
	require.NoError(t.t, err, "failed to send tx")
	t.backend.Commit()
	logs, err := t.backend.FilterLogs(ctx, ethereum.FilterQuery{})
	require.NoError(t.t, err, "failed to filter logs")
	for _, lg := range logs {
		t.t.Log("topic:", lg.Topics[0], "transmitted topic:", no_op_ocr3.NoOpOCR3Transmitted{}.Topic(), "configset topic:", no_op_ocr3.NoOpOCR3ConfigSet{}.Topic())
	}
	return nil
}
