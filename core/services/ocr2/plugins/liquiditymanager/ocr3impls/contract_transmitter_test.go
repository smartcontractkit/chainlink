package ocr3impls_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/ocr3impls"
)

type testUniverse[RI ocr3impls.MultichainMeta] struct {
	simClient       *client.SimulatedBackendClient
	backend         *backends.SimulatedBackend
	deployer        *bind.TransactOpts
	transmitters    []*bind.TransactOpts
	signers         []common.Address
	wrapper         *no_op_ocr3.NoOpOCR3
	ocr3Transmitter ocr3types.ContractTransmitter[RI]
	keyrings        []ocr3types.OnchainKeyring[RI]
	f               uint8
}

type keyringsAndSigners[RI ocr3impls.MultichainMeta] struct {
	keyrings []ocr3types.OnchainKeyring[RI]
	signers  []common.Address
}

func newTestUniverse[RI ocr3impls.MultichainMeta](
	t *testing.T,
	ks *keyringsAndSigners[RI]) testUniverse[RI] {
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
		keyrings []ocr3types.OnchainKeyring[RI]
		signers  []common.Address
	)
	if ks != nil {
		keyrings = ks.keyrings
		signers = ks.signers
	} else {
		for i := 0; i < 4; i++ {
			kb, err2 := ocr2key.New(chaintype.EVM)
			require.NoError(t, err2, "failed to create key")
			kr := ocr3impls.NewOnchainKeyring[RI](kb, logger.TestLogger(t))
			signers = append(signers, common.BytesToAddress(kr.PublicKey()))
			keyrings = append(keyrings, kr)
		}
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

	// print config set event
	iter, err := wrapper.FilterConfigSet(&bind.FilterOpts{
		Start: 1,
	})
	require.NoError(t, err, "failed to create filter iterator")
	for iter.Next() {
		event := iter.Event
		t.Log("onchain signers:", event.Signers, "set signers:", signers)
		t.Log("transmitters:", event.Transmitters, "set transmitters:", transmitters)
	}

	contractABI, err := no_op_ocr3.NoOpOCR3MetaData.GetAbi()
	require.NoError(t, err, "failed to get abi")
	tImpl := &transmitterImpl{
		backend: backend,
		from:    transmitters[0],
		t:       t,
	}
	ocr3Transmitter, err := ocr3impls.NewOCR3ContractTransmitter[RI](
		addr,
		*contractABI,
		tImpl,
		logger.TestLogger(t),
		nil, // reportToEvmTxMeta, unused
	)
	require.NoError(t, err, "failed to create OCR3ContractTransmitter")

	return testUniverse[RI]{
		backend:         backend,
		deployer:        deployer,
		transmitters:    transmitters,
		signers:         signers,
		wrapper:         wrapper,
		keyrings:        keyrings,
		ocr3Transmitter: ocr3Transmitter,
		f:               f,
		simClient:       client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID),
	}
}

func (uni testUniverse[RI]) SignReport(t *testing.T, configDigest ocrtypes.ConfigDigest, rwi ocr3types.ReportWithInfo[RI], seqNum uint64) []ocrtypes.AttributedOnchainSignature {
	var attributedSigs []ocrtypes.AttributedOnchainSignature
	for i := uint8(0); i < uni.f+1; i++ {
		t.Log("signing report with", hexutil.Encode(uni.keyrings[i].PublicKey()))
		sig, err := uni.keyrings[i].Sign(configDigest, seqNum, rwi)
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

// TODO: might be useful to fuzz test this
func TestContractTransmitter(t *testing.T) {
	t.Parallel()

	t.Run("empty report", func(t *testing.T) {
		uni := newTestUniverse[multichainMeta](t, nil)

		c, err := uni.wrapper.LatestConfigDetails(nil)
		require.NoError(t, err, "failed to get latest config digest and epoch")
		configDigest := c.ConfigDigest

		// create the attributed signatures
		// only need f+1 which is 2 in this case
		rwi := ocr3types.ReportWithInfo[multichainMeta]{
			Report: []byte{},
			Info: multichainMeta{
				configDigest: configDigest,
			},
		}
		seqNum := uint64(1)
		attributedSigs := uni.SignReport(t, configDigest, rwi, seqNum)

		account, err := uni.ocr3Transmitter.FromAccount()
		require.NoError(t, err, "failed to get from account")
		require.Equal(t, account, ocrtypes.Account(uni.transmitters[0].From.Hex()), "unexpected from account")
		err = uni.ocr3Transmitter.Transmit(context.Background(), configDigest, seqNum, rwi, attributedSigs)
		require.NoError(t, err, "failed to transmit report")

		lcde, err := uni.wrapper.LatestConfigDetails(nil)
		require.NoError(t, err, "failed to get latest config details")
		sn, err := uni.wrapper.LatestSequenceNumber(nil)
		require.NoError(t, err, "failed to get latest sequence number")
		require.Equal(t, configDigest, lcde.ConfigDigest, "config digest mismatch")
		require.Equal(t, seqNum, sn, "seq number mismatch")

		// check for transmitted event
		events := uni.TransmittedEvents(t)
		require.Len(t, events, 1, "expected one transmitted event")
		event := events[0]
		require.Equal(t, configDigest, event.ConfigDigest, "unexpected config digest")
		require.Equal(t, seqNum, event.SequenceNumber, "unexpected sequence number")
	})

	t.Run("non-empty report", func(t *testing.T) {
		uni := newTestUniverse[multichainMeta](t, nil)

		c, err := uni.wrapper.LatestConfigDetails(nil)
		require.NoError(t, err, "failed to get latest config digest and epoch")
		configDigest := c.ConfigDigest

		// create the attributed signatures
		// only need f+1 which is 2 in this case
		rep := testutils.Random32Byte()
		rwi := ocr3types.ReportWithInfo[multichainMeta]{
			// Report bytes must always be aligned to 32 byte boundaries otherwise the on-chain
			// length check will fail.
			Report: rep[:],
			Info: multichainMeta{
				configDigest: configDigest,
			},
		}
		seqNum := uint64(1)
		attributedSigs := uni.SignReport(t, configDigest, rwi, seqNum)

		account, err := uni.ocr3Transmitter.FromAccount()
		require.NoError(t, err, "failed to get from account")
		require.Equal(t, account, ocrtypes.Account(uni.transmitters[0].From.Hex()), "unexpected from account")
		err = uni.ocr3Transmitter.Transmit(context.Background(), configDigest, seqNum, rwi, attributedSigs)
		require.NoError(t, err, "failed to transmit report")

		lcde, err := uni.wrapper.LatestConfigDetails(nil)
		require.NoError(t, err, "failed to get latest config details")
		sn, err := uni.wrapper.LatestSequenceNumber(nil)
		require.NoError(t, err, "failed to get latest sequence number")
		require.Equal(t, configDigest, lcde.ConfigDigest, "config digest mismatch")
		require.Equal(t, seqNum, sn, "seq number mismatch")

		// check for transmitted event
		events := uni.TransmittedEvents(t)
		require.Len(t, events, 1, "expected one transmitted event")
		event := events[0]
		require.Equal(t, configDigest, event.ConfigDigest, "unexpected config digest")
		require.Equal(t, seqNum, event.SequenceNumber, "unexpected sequence number")
	})
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
		Gas:      1e6,
		To:       &to,
		Data:     data,
	})
	signedTx, err := t.from.Signer(t.from.From, rawTx)
	require.NoError(t.t, err, "failed to sign tx")
	g, err := t.backend.EstimateGas(ctx, ethereum.CallMsg{
		From:     t.from.From,
		To:       &to,
		Gas:      1e6,
		GasPrice: gp,
		Data:     data,
	})
	require.NoError(t.t, err, "failed to estimate gas")
	t.t.Log("estimated gas:", g)
	err = t.backend.SendTransaction(ctx, signedTx)
	require.NoError(t.t, err, "failed to send tx")
	t.backend.Commit()
	receipt, err := t.backend.TransactionReceipt(ctx, signedTx.Hash())
	require.NoError(t.t, err, "failed to get tx receipt")
	require.Equal(t.t, uint64(1), receipt.Status, "tx failed")
	return nil
}
