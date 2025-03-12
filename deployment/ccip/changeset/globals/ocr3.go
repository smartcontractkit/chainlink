package globals

import (
	"fmt"
	"time"

	"dario.cat/mergo"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// Intention of this file is to be a single source of the truth for OCR3 parameters used by CCIP plugins.
//
// Assumptions:
//   - Although, some values are similar between Commit and Execute, we should keep them separate, because
//     these plugins have different requirements and characteristics. This way we can avoid misconfiguration
//     by accidentally changing parameter for one plugin while adjusting it for the other
//   - OCR3 parameters are chain agnostic and should be reused across different chains. There might be some use cases
//     for overrides to accommodate specific chain characteristics (e.g. Ethereum).
//     However, for most of the cases we should strive to rely on defaults under CommitOCRParams and ExecOCRParams.
//     This makes the testing process much easier and increase our confidence that the configuration is safe to use.
//   - The fewer overrides the better. Introducing new overrides should be done with caution and only if there's a strong
//     justification for it. Moreover, it requires detailed chaos / load testing to ensure that the new parameters are safe to use
//     and meet CCIP SLOs
//   - Single params must not be stored under const or exposed outside of this file to limit the risk of
//     accidental configuration or partial configuration
//   - MaxDurations should be set on the latencies observed on various environments using p99 OCR3 latencies
//     These values should be specific to the plugin type and should not depend on the chain family
//     or the environment in which plugin runs
var (
	// CommitOCRParams represents the default OCR3 parameters for all chains (beside Ethereum, see CommitOCRParamsForEthereum).
	// Most of the intervals here should be generic enough (and chain agnostic) to be reused across different chains.
	CommitOCRParams = types.OCRParameters{
		DeltaProgress:               120 * time.Second,
		DeltaResend:                 30 * time.Second,
		DeltaInitial:                20 * time.Second,
		DeltaRound:                  15 * time.Second,
		DeltaGrace:                  5 * time.Second,
		DeltaCertifiedCommitRequest: 10 * time.Second,
		// TransmissionDelayMultiplier overrides DeltaStage
		DeltaStage:                              25 * time.Second,
		Rmax:                                    3,
		MaxDurationQuery:                        7 * time.Second,
		MaxDurationObservation:                  13 * time.Second,
		MaxDurationShouldAcceptAttestedReport:   5 * time.Second,
		MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
	}

	// CommitOCRParamsForEthereum represents a dedicated set of OCR3 parameters for Ethereum.
	// It's driven by the fact that Ethereum block time is slow (12 seconds) and chain is considered
	// more expensive to other EVM compatible chains
	CommitOCRParamsForEthereum = withOverrides(
		CommitOCRParams,
		types.OCRParameters{
			DeltaRound: 90 * time.Second,
			DeltaStage: 60 * time.Second,
		},
	)
)

var (
	// ExecOCRParams represents the default OCR3 parameters for all chains (beside Ethereum, see ExecOCRParamsForEthereum).
	ExecOCRParams = types.OCRParameters{
		DeltaProgress:               100 * time.Second,
		DeltaResend:                 30 * time.Second,
		DeltaInitial:                20 * time.Second,
		DeltaRound:                  15 * time.Second,
		DeltaGrace:                  5 * time.Second,
		DeltaCertifiedCommitRequest: 10 * time.Second,
		// TransmissionDelayMultiplier overrides DeltaStage
		DeltaStage: 25 * time.Second,
		Rmax:       3,
		// MaxDurationQuery is set to very low value, because Execution plugin doesn't use Query
		MaxDurationQuery:                        200 * time.Millisecond,
		MaxDurationObservation:                  13 * time.Second,
		MaxDurationShouldAcceptAttestedReport:   5 * time.Second,
		MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
	}

	// ExecOCRParamsForEthereum represents a dedicated set of OCR3 parameters for Ethereum.
	// Similarly to Commit, it's here to accommodate Ethereum specific characteristics
	ExecOCRParamsForEthereum = withOverrides(
		ExecOCRParams,
		types.OCRParameters{
			DeltaRound: 90 * time.Second,
			DeltaStage: 60 * time.Second,
		},
	)
)

func withOverrides(base types.OCRParameters, overrides types.OCRParameters) types.OCRParameters {
	outcome := base
	if err := mergo.Merge(&outcome, overrides, mergo.WithOverride); err != nil {
		panic(fmt.Sprintf("error while building an OCR config %v", err))
	}
	return outcome
}
