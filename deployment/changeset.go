package deployment

import (
	"encoding/json"
	"errors"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/mcms"
)

var (
	ErrInvalidConfig = errors.New("invalid changeset config")
)

// ChangeSet represents a set of changes to be made to an environment.
// The configuration contains environment specific inputs for a specific changeset.
// The configuration might contain for example the chainSelectors to apply the change to
// or existing environment specific contract addresses.
// Its recommended that changesets operate on a small number of chains (e.g. 1-3)
// to reduce the risk of partial failures.
// If the configuration is unexpected type or format, the changeset should return ErrInvalidConfig.
type ChangeSet[C any] func(e Environment, config C) (ChangesetOutput, error)

// ProposedJob represents a job spec which has been proposed to a node, with the JobID returned by the
// Job Distributor.
type ProposedJob struct {
	JobID string
	Node  string
	Spec  string
}

// ChangesetOutput is the output of a Changeset function.
// Think of it like a state transition output.
// The address book here should contain only new addresses created in
// this changeset.
type ChangesetOutput struct {
	// Deprecated: Prefer Jobs instead.
	JobSpecs map[string][]string `deprecated:"true"`
	Jobs     []ProposedJob
	// Deprecated: Prefer MCMSTimelockProposals instead, will be removed in future
	Proposals             []timelock.MCMSWithTimelockProposal
	MCMSTimelockProposals []mcms.TimelockProposal
	MCMSProposals         []mcms.Proposal
	AddressBook           AddressBook
}

// StateRepresentation is a type which represents a given piece of business domain state, e.g. CCIPView. A
// A StateRepresentation can render into Json (conforms to `json.Marshaler`), and may have other representations in
// the future, such as a protobuf representation.
//
// StateRepresentation is used in storing and transmitting domain state through artifacts (local or transmitted).
// In the chainlink/deployment module, it is only used in tests. ChangeSet execution functions have access to state
// representations via GetState(env, StateKey[SomeStateStruct]{}) calls.
type StateRepresentation interface {
	json.Marshaler
}

// ViewState produces a product specific JSON representation of
// the on and offchain state of the environment.
//
// deprecated: use StateRenderer
type ViewState func(e Environment) (json.Marshaler, error)

// StateRenderer is a function which transforms elements of the environment into a StateRepresentation, so that the
// business domain state may be transmitted or persisted. Every domain should have a state representation, and a
// function to transform environmental elements (such as the address book, and any other artifacts) into a coherent

type StateRenderer[S StateRepresentation] func(e Environment) (*S, error)
