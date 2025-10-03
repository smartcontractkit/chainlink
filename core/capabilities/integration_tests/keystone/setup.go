package keystone

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

var (
	workflowName    = "abcdef0123"
	workflowOwnerID = "0100000000000000000000000000000000000001"
)

func setupKeystoneDons(ctx context.Context, t *testing.T, lggr logger.Logger,
	workflowDonInfo framework.DonConfiguration,
	triggerDonInfo framework.DonConfiguration,
	targetDonInfo framework.DonConfiguration,
	trigger framework.TriggerFactory) (workflowDon *framework.DON, consumer *feeds_consumer.KeystoneFeedsConsumer) {
	donContext := framework.CreateDonContext(ctx, t)

	workflowDon = createKeystoneWorkflowDon(ctx, t, lggr, workflowDonInfo, triggerDonInfo, targetDonInfo, donContext)

	forwarderAddr, _ := SetupForwarderContract(t, workflowDon, donContext.EthBlockchain)
	_, consumer = SetupConsumerContract(t, donContext.EthBlockchain, forwarderAddr, workflowOwnerID, workflowName)

	writeTargetDon := createKeystoneWriteTargetDon(ctx, t, lggr, targetDonInfo, donContext, forwarderAddr)

	triggerDon := createKeystoneTriggerDon(ctx, t, lggr, triggerDonInfo, donContext, trigger)

	servicetest.Run(t, workflowDon)
	servicetest.Run(t, triggerDon)
	servicetest.Run(t, writeTargetDon)

	donContext.WaitForCapabilitiesToBeExposed(t, workflowDon, triggerDon, writeTargetDon)

	return workflowDon, consumer
}

func createKeystoneTriggerDon(ctx context.Context, t *testing.T, lggr logger.Logger, triggerDonInfo framework.DonConfiguration,
	donContext framework.DonContext, trigger framework.TriggerFactory) *framework.DON {
	triggerDon := framework.NewDON(ctx, t, lggr, triggerDonInfo,
		[]commoncap.DON{}, donContext, false, 1*time.Second)

	triggerDon.AddExternalTriggerCapability(trigger)
	triggerDon.Initialise()
	return triggerDon
}

func createKeystoneWriteTargetDon(ctx context.Context, t *testing.T, lggr logger.Logger, targetDonInfo framework.DonConfiguration, donContext framework.DonContext, forwarderAddr common.Address) *framework.DON {
	writeTargetDon := framework.NewDON(ctx, t, lggr, targetDonInfo,
		[]commoncap.DON{}, donContext, false, 1*time.Second)
	_, err := writeTargetDon.AddPublishedEthereumWriteTargetNonStandardCapability(forwarderAddr)
	require.NoError(t, err)
	writeTargetDon.Initialise()
	return writeTargetDon
}

func createKeystoneWorkflowDon(ctx context.Context, t *testing.T, lggr logger.Logger, workflowDonInfo framework.DonConfiguration,
	triggerDonInfo framework.DonConfiguration, targetDonInfo framework.DonConfiguration, donContext framework.DonContext) *framework.DON {
	workflowDon := framework.NewDON(ctx, t, lggr, workflowDonInfo,
		[]commoncap.DON{triggerDonInfo.DON, targetDonInfo.DON},
		donContext, true, 1*time.Second)

	workflowDon.AddOCR3NonStandardCapability()
	workflowDon.Initialise()
	return workflowDon
}

func createFeedReport(t *testing.T, price *big.Int, observationTimestamp int64,
	feedIDString string,
	keyBundles []ocr2key.KeyBundle) *datastreams.FeedReport {
	reportCtx := ocrTypes.ReportContext{}
	rawCtx := RawReportContext(reportCtx)

	bytes, err := hex.DecodeString(feedIDString[2:])
	require.NoError(t, err)
	var feedIDBytes [32]byte
	copy(feedIDBytes[:], bytes)

	report := &datastreams.FeedReport{
		FeedID:               feedIDString,
		FullReport:           newReport(t, feedIDBytes, price, observationTimestamp),
		BenchmarkPrice:       price.Bytes(),
		ObservationTimestamp: observationTimestamp,
		Signatures:           [][]byte{},
		ReportContext:        rawCtx,
	}

	for _, key := range keyBundles {
		sig, err := key.Sign(reportCtx, report.FullReport)
		require.NoError(t, err)
		report.Signatures = append(report.Signatures, sig)
	}

	return report
}

func RawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

func newReport(t *testing.T, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	v3Codec := reportcodec.NewReportCodec(feedID, logger.Test(t))
	raw, err := v3Codec.BuildReport(t.Context(), v3.ReportFields{
		BenchmarkPrice: price,

		Timestamp: uint32(timestamp), //nolint:gosec // G115
		Bid:       big.NewInt(0),
		Ask:       big.NewInt(0),
		LinkFee:   big.NewInt(0),
		NativeFee: big.NewInt(0),
	})
	require.NoError(t, err)
	return raw
}
