package flows

import (
	"log"
	"time"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/service"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	common "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func ConditionalTriggerFlows(
	coord ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	ratio types.Ratio,
	getter common.ConditionalUpkeepProvider,
	subscriber common.BlockSubscriber,
	builder common.PayloadBuilder,
	resultStore types.ResultStore,
	metadataStore types.MetadataStore,
	runner ocr2keepersv3.Runner,
	proposalQ types.ProposalQueue,
	retryQ types.RetryQueue,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) []service.Recoverable {
	preprocessors := []ocr2keepersv3.PreProcessor[common.UpkeepPayload]{coord}

	// runs full check pipeline on a coordinated block with coordinated upkeeps
	conditionalFinal := newFinalConditionalFlow(preprocessors, resultStore, runner, FinalConditionalInterval, proposalQ, builder, retryQ, stateUpdater, logger)

	// the sampling proposal flow takes random samples of active upkeeps, checks
	// them and surfaces the ids if the items are eligible
	conditionalProposal := newSampleProposalFlow(preprocessors, ratio, getter, metadataStore, runner, SamplingConditionInterval, logger)

	return []service.Recoverable{conditionalFinal, conditionalProposal}
}

func LogTriggerFlows(
	coord ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	resultStore types.ResultStore,
	metadataStore types.MetadataStore,
	runner ocr2keepersv3.Runner,
	logProvider common.LogEventProvider,
	rp common.RecoverableProvider,
	builder common.PayloadBuilder,
	logInterval time.Duration,
	recoveryProposalInterval time.Duration,
	recoveryFinalInterval time.Duration,
	retryQ types.RetryQueue,
	proposals types.ProposalQueue,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) []service.Recoverable {
	// all flows use the same preprocessor based on the coordinator
	// each flow can add preprocessors to this provided slice
	preprocessors := []ocr2keepersv3.PreProcessor[common.UpkeepPayload]{coord}

	// the recovery proposal flow is for nodes to surface payloads that should
	// be recovered. these values are passed to the network and the network
	// votes on the proposed values
	rcvProposal := newRecoveryProposalFlow(preprocessors, runner, metadataStore, rp, recoveryProposalInterval, stateUpdater, logger)

	// the final recovery flow takes recoverable payloads merged with the latest
	// blocks and runs the pipeline for them. these values to run are derived
	// from node coordination and it can be assumed that all values should be
	// run.
	rcvFinal := newFinalRecoveryFlow(preprocessors, resultStore, runner, retryQ, recoveryFinalInterval, proposals, builder, stateUpdater, logger)

	// the log trigger flow is the happy path for log trigger payloads. all
	// retryables that are encountered in this flow are elevated to the retry
	// flow
	logTrigger := newLogTriggerFlow(preprocessors, resultStore, runner, logProvider, logInterval, retryQ, stateUpdater, logger)

	return []service.Recoverable{
		rcvProposal,
		rcvFinal,
		logTrigger,
	}
}
