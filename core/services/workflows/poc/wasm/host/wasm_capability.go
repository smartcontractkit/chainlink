package host

import (
	"fmt"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type wasmCapability struct {
	runner         *wasmGuestRunner
	capabilityType commoncap.CapabilityType
}

func (c *wasmCapability) CapabilityType() commoncap.CapabilityType {
	return c.capabilityType
}

func (c *wasmCapability) Run(ref string, value values.Value) (values.Value, bool, error) {
	returnCh := make(chan runResult)
	c.runner.inputs <- computeRequest{
		stepRef: ref,
		input:   value,
		retCh:   returnCh,
	}
	select {
	case retVal := <-returnCh:
		return retVal.retVal, retVal.cont, retVal.err
	// TODO what do we actually do on an exit?
	case exit := <-c.runner.exitCh:
		// send back to exit channel so that we don't freeze up for now, need better solution
		c.runner.exitCh <- exit
		return nil, false, fmt.Errorf("program already exited with error %w", exit.err)
	}
}
