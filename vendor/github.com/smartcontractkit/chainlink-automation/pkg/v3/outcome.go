package ocr2keepers

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

// NOTE: Any change to these values should keep backwards compatibility in mind
// as different nodes would upgrade at different times and would need to
// adhere to each others' limits
const (
	OutcomeAgreedPerformablesLimit = 100
	OutcomeSurfacedProposalsLimit  = 50
	// NOTE: This should be enough number of rounds to encompass ObservationProcessLimit
	// so that proposals have enough time to be processed before getting coordinated on
	// a new block
	OutcomeSurfacedProposalsRoundHistoryLimit = 20

	// MaxOutcomeLength applies a limit to the total length of bytes in an outcome for
	// a round. NOTE: This is derived from a limit of 10000 on performData
	// which is guaranteed onchain
	MaxOutcomeLength = 2_500_000
	// MaxReportLength limits the total length of bytes for a single report.
	MaxReportLength = 1_000_000
	// MaxReportCount limits the total number of reports allowed to be produced
	// by the OCR protocol in a single round. This should be at least the number
	// of allowed agreed performables in a single round.
	MaxReportCount = OutcomeAgreedPerformablesLimit
)

// AutomationOutcome represents agreed upon state by the network, derived from
// a collection of AutomationObservations with applied quorum thresholds
// A byzQuorem # of nodes must agree that at least f+1 nodes sent the same AgreedPerformable
// or SurfacedProposal to include that item in an outcome
// NOTE: Any change to this structure should keep backwards compatibility in mind
// as different nodes would upgrade at different times and would need to understand
// each others' outcome meanwhile
type AutomationOutcome struct {
	// These are the upkeeps that achieved quorum, meaning should be performed on chain
	// These require quorum of f+1 nodes
	AgreedPerformables []ocr2keepers.CheckResult
	// These are the proposals with a coordinated block that should be run through the
	// check pipeline. The proposals remain valid for a range of rounds where they do
	// not get tied to a new coordinated block in order to give check pipeline enough
	// time to run asynchronously
	// Quorum of f+1 is only applied on the blockNumber and blockHash of the proposal
	// rest of the fields can be manipulated by malicious nodes
	SurfacedProposals [][]ocr2keepers.CoordinatedBlockProposal
}

// ValidateAutomationOutcome validates individual values in an AutomationOutcome
func validateAutomationOutcome(o AutomationOutcome, utg types.UpkeepTypeGetter, wg types.WorkIDGenerator) error {
	// Validate AgreedPerformables
	if (len(o.AgreedPerformables)) > OutcomeAgreedPerformablesLimit {
		return fmt.Errorf("outcome performable length cannot be greater than %d", OutcomeAgreedPerformablesLimit)
	}
	seenPerformables := make(map[string]bool)
	for _, res := range o.AgreedPerformables {
		if err := validateCheckResult(res, utg, wg); err != nil {
			return err
		}
		if seenPerformables[res.WorkID] {
			return fmt.Errorf("agreed performable cannot have duplicate workIDs")
		}
		seenPerformables[res.WorkID] = true
	}

	// Validate SurfacedProposals
	if len(o.SurfacedProposals) >
		OutcomeSurfacedProposalsRoundHistoryLimit {
		return fmt.Errorf("number of rounds for surfaced proposals cannot be greater than %d", OutcomeSurfacedProposalsRoundHistoryLimit)
	}
	seenProposals := make(map[string]bool)
	for _, round := range o.SurfacedProposals {
		if len(round) > OutcomeSurfacedProposalsLimit {
			return fmt.Errorf("number of surfaced proposals in a round cannot be greater than %d", OutcomeSurfacedProposalsLimit)
		}
		for _, proposal := range round {
			if err := validateUpkeepProposal(proposal, utg, wg); err != nil {
				return err
			}
			if seenProposals[proposal.WorkID] {
				return fmt.Errorf("proposals cannot have duplicate workIDs")
			}
			seenProposals[proposal.WorkID] = true
		}
	}
	return nil
}

// Encode produces a json encoded array of bytes. Possible errors come from the
// encoding/json package
func (outcome AutomationOutcome) Encode() ([]byte, error) {
	return json.Marshal(outcome)
}

// DecodeAutomationOutcome decodes an AutomationOutcome from an encoded array
// of bytes. Possible errors come from the encoding/json package
func DecodeAutomationOutcome(data []byte, utg types.UpkeepTypeGetter, wg types.WorkIDGenerator) (AutomationOutcome, error) {
	ao := AutomationOutcome{}
	err := json.Unmarshal(data, &ao)
	if err != nil {
		return AutomationOutcome{}, err
	}
	err = validateAutomationOutcome(ao, utg, wg)
	if err != nil {
		return AutomationOutcome{}, err
	}
	return ao, err
}
