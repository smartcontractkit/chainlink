package cre

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
	stellarfeature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/stellar"
	datafeedswrite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/stellar/datafeeds/write/config"
	thelpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const stellarDataFeedsWorkflowFile = "./stellar/datafeeds/write/main.go"

const stellarDataFeedsAnswer int64 = 1234567890

var stellarDataFeedsDataID = [32]byte{0x01, 0x8e, 0x16, 0xc3, 0x9e, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

func executeStellarDataFeedsWriteTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	stellarChain *stellchain.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	t.Helper()
	lggr := framework.L
	ctx := context.Background()

	workflowName := thelpers.UniqueStellarWorkflowName("stellar-df-write-workflow")

	var workflowOwner [20]byte
	copy(workflowOwner[:], workflowRegistryOwnerBytes(t, tenv))
	var truncatedName [10]byte
	copy(truncatedName[:], pkgworkflows.HashTruncateName(workflowName))

	cacheID, err := stellarfeature.DeployAndConfigureStellarDataFeedsCache(ctx, stellarChain, tenv.CreEnvironment, stellarDataFeedsDataID, "CRE-E2E", workflowOwner, truncatedName)
	require.NoError(t, err, "failed to deploy and configure Stellar data feeds cache")
	lggr.Info().Str("cache", cacheID).Msg("Deployed Stellar data feeds cache")

	requiredSignatures := stellarRequiredSignatures(t, tenv)

	workflowConfig := datafeedswrite_config.Config{
		ChainSelector:      stellarChain.ChainSelector(),
		WorkflowName:       workflowName,
		CacheContractID:    cacheID,
		DataIDHex:          hex.EncodeToString(stellarDataFeedsDataID[:]),
		Answer:             stellarDataFeedsAnswer,
		RequiredSignatures: requiredSignatures,
	}
	workflowID := thelpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, stellarDataFeedsWorkflowFile)

	expectedLog := "Stellar DF write succeeded"
	thelpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, thelpers.WorkflowEngineInitErrorLog, expectedLog, thelpers.StellarWorkflowTimeout, thelpers.WithUserLogWorkflowID(workflowID))

	expected := big.NewInt(stellarDataFeedsAnswer)
	require.Eventually(t, func() bool {
		round, rErr := stellarfeature.DataFeedsCacheLatestRound(ctx, stellarChain, cacheID, stellarDataFeedsDataID)
		if rErr != nil {
			lggr.Warn().Err(rErr).Msg("stellar data feeds latest_round query failed; retrying")
			return false
		}
		if round == nil || round.Answer == nil {
			return false
		}
		lggr.Info().Str("on_chain_answer", round.Answer.String()).Uint64("round_id", round.RoundId).Msg("Stellar data feeds round read")
		return round.Answer.Cmp(expected) == 0 && round.RoundId >= 1 && round.Timestamp > 0
	}, 2*time.Minute, 5*time.Second, "Stellar data feeds cache did not record answer %d", stellarDataFeedsAnswer)

	lggr.Info().Str("expected_log", expectedLog).Int64("answer", stellarDataFeedsAnswer).Msg("Stellar data feeds write test passed")
}
