package launcher

import (
	"fmt"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

type ccipPlugins map[ocrtypes.ConfigDigest]cctypes.CCIPOracle

// StartAll will call Oracle.Start on an entire don
func (c ccipPlugins) StartAll() error {
	nilPlugins := make(ccipPlugins)
	return c.Transition(nilPlugins)
}

// CloseAll is used to shut down an entire don immediately
func (c ccipPlugins) CloseAll() error {
	nilPlugins := make(ccipPlugins)
	return nilPlugins.Transition(c)
}

// Transition manages starting and stopping ocr instances
// If there are any new config digests, we need to start those instances
// If any of the previous config digests are no longer present, we need to shut those down
// We don't care about if they're exec/commit or active/candidate, that all happens in the plugin
func (c ccipPlugins) Transition(prevPlugins ccipPlugins) error {
	var err error

	if len(c) > 4 || len(prevPlugins) > 4 {
		return fmt.Errorf("have more than 4 instances somehow")
	}
	// This shuts down instances that were present previously, but are no longer needed
	for digest, oracle := range prevPlugins {
		if _, ok := c[digest]; !ok {
			err = multierr.Append(err, oracle.Close())
		}
	}

	// This will start the instances that were not previously present, but are in the new config
	for digest, oracle := range c {
		if _, ok := prevPlugins[digest]; !ok {
			err = multierr.Append(err, oracle.Start())
		}
	}
	return err
}
