package deployment

import (
	"encoding/json"
	"errors"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
)

// Changeset represents a set of changes to be made to an environment.
// The configuration contains environment specific inputs for a specific changeset.
// The configuration might contain for example the chainSelectors to apply the change to
// or existing environment specific contract addresses.
// Its recommended that changesets operate on a small number of chains (e.g. 1-3)
// to reduce the risk of partial failures.
type Changeset func(e Environment, config interface{}) (ChangesetOutput, error)

// ChangesetOutput is the output of a Changeset function.
type ChangesetOutput struct {
	JobSpecs    map[string][]string
	Proposals   []timelock.MCMSWithTimelockProposal
	AddressBook AddressBook
}

var (
	ErrInvalidConfig = errors.New("invalid config")
)

// ViewState produces a product specific JSON representation of
// the on and offchain state of the environment.
type ViewState func(e Environment, ab AddressBook) (json.Marshaler, error)
