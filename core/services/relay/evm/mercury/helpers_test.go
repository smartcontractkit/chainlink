package mercury

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	reportcodecv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	reportcodecv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2/reportcodec"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

var (
	sampleFeedID       = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	sampleClientPubKey = hexutil.MustDecode("0x724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93")
)

var sampleReports [][]byte

var (
	sampleV1Report      = buildSampleV1Report(242)
	sampleV2Report      = buildSampleV2Report(242)
	sampleV3Report      = buildSampleV3Report(242)
	sig2                = ocrtypes.AttributedOnchainSignature{Signature: testutils.MustDecodeBase64("kbeuRczizOJCxBzj7MUAFpz3yl2WRM6K/f0ieEBvA+oTFUaKslbQey10krumVjzAvlvKxMfyZo0WkOgNyfF6xwE="), Signer: 2}
	sig3                = ocrtypes.AttributedOnchainSignature{Signature: testutils.MustDecodeBase64("9jz4b6Dh2WhXxQ97a6/S9UNjSfrEi9016XKTrfN0mLQFDiNuws23x7Z4n+6g0sqKH/hnxx1VukWUH/ohtw83/wE="), Signer: 3}
	sampleSigs          = []ocrtypes.AttributedOnchainSignature{sig2, sig3}
	sampleReportContext = ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: MustHexToConfigDigest("0x0006fc30092226b37f6924b464e16a54a7978a9a524519a73403af64d487dc45"),
			Epoch:        6,
			Round:        28,
		},
		ExtraHash: [32]uint8{27, 144, 106, 73, 166, 228, 123, 166, 179, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114},
	}
)

func init() {
	sampleReports = make([][]byte, 4)
	for i := 0; i < len(sampleReports); i++ {
		sampleReports[i] = buildSampleV1Report(int64(i))
	}
}

func buildSampleV1Report(p int64) []byte {
	feedID := sampleFeedID
	timestamp := uint32(42)
	bp := big.NewInt(p)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(143)
	currentBlockHash := utils.NewHash()
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(142)

	b, err := reportcodecv1.ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, currentBlockTimestamp, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}

func buildSampleV2Report(ts int64) []byte {
	feedID := sampleFeedID
	timestamp := uint32(ts)
	bp := big.NewInt(242)
	validFromTimestamp := uint32(123)
	expiresAt := uint32(456)
	linkFee := big.NewInt(3334455)
	nativeFee := big.NewInt(556677)

	b, err := reportcodecv2.ReportTypes.Pack(feedID, validFromTimestamp, timestamp, nativeFee, linkFee, expiresAt, bp)
	if err != nil {
		panic(err)
	}
	return b
}

func buildSampleV3Report(ts int64) []byte {
	feedID := sampleFeedID
	timestamp := uint32(ts)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	validFromTimestamp := uint32(123)
	expiresAt := uint32(456)
	linkFee := big.NewInt(3334455)
	nativeFee := big.NewInt(556677)

	b, err := reportcodecv3.ReportTypes.Pack(feedID, validFromTimestamp, timestamp, nativeFee, linkFee, expiresAt, bp, bid, ask)
	if err != nil {
		panic(err)
	}
	return b
}

func buildSamplePayload(report []byte) []byte {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range sampleSigs {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := evmutil.RawReportContext(sampleReportContext)
	payload, err := PayloadTypes.Pack(rawReportCtx, report, rs, ss, vs)
	if err != nil {
		panic(err)
	}
	return payload
}

type TestHarness struct {
	configPoller     *ConfigPoller
	user             *bind.TransactOpts
	backend          *backends.SimulatedBackend
	verifierAddress  common.Address
	verifierContract *verifier.Verifier
	logPoller        logpoller.LogPoller
}

func SetupTH(t *testing.T, feedID common.Hash) TestHarness {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)

	proxyAddress, _, verifierProxy, err := verifier_proxy.DeployVerifierProxy(user, b, common.Address{})
	require.NoError(t, err, "failed to deploy test mercury verifier proxy contract")
	verifierAddress, _, verifierContract, err := verifier.DeployVerifier(user, b, proxyAddress)
	require.NoError(t, err, "failed to deploy test mercury verifier contract")
	_, err = verifierProxy.InitializeVerifier(user, verifierAddress)
	require.NoError(t, err)
	b.Commit()

	db := pgtest.NewSqlxDB(t)
	ethClient := evmclient.NewSimulatedBackendClient(t, b, big.NewInt(1337))
	lggr := logger.TestLogger(t)
	lorm := logpoller.NewORM(big.NewInt(1337), db, lggr)

	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        2,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	ht := headtracker.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(lorm, ethClient, lggr, ht, lpOpts)
	servicetest.Run(t, lp)

	configPoller, err := NewConfigPoller(testutils.Context(t), lggr, lp, verifierAddress, feedID)
	require.NoError(t, err)

	configPoller.Start()

	return TestHarness{
		configPoller:     configPoller,
		user:             user,
		backend:          b,
		verifierAddress:  verifierAddress,
		verifierContract: verifierContract,
		logPoller:        lp,
	}
}
