package deployment

import (
	"errors"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
)

// Services as input to CI/Async tasks
type ChangesetOutput struct {
	JobSpecs    map[string][]string
	Proposals   []timelock.MCMSWithTimelockProposal
	AddressBook AddressBook
}

type Changeset func(e Environment, config interface{}) (ChangesetOutput, error)

var (
	ErrInvalidConfig = errors.New("invalid config")
)

type ViewState func(e Environment, ab AddressBook) (string, error)
