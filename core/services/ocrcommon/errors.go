package ocrcommon

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type ErrRelayNotFound struct {
	Relay relay.Network
}

func (e ErrRelayNotFound) Error() string {
	return fmt.Sprintf("%v relay does not exist is it enabled?", e.Relay)
}

type ErrUnexpectedJobType struct {
	Expected job.Type
	Actual   job.Type
}

func (e ErrUnexpectedJobType) Error() string {
	return fmt.Sprintf("Delegate expected an %s to be present, got %s", e.Actual, e.Expected)
}
