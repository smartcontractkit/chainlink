package launcher

import (
	"golang.org/x/exp/maps"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"
)

type ccipPlugins map[ocrtypes.ConfigDigest]*cctypes.CCIPOracle

// Close is used to shut down an entire don immediately
func (c *ccipPlugins) Close() error {
	if c == nil {
		return nil
	}
	var err error

	for _, oracle := range *c {
		err = multierr.Append(err, (*oracle).Close())
	}

	return err
}

// Transition manages starting and stopping ocr instances
// If there are any new config digests, we need to start those instances
// If any of the previous config digests are no longer present, we need to shut those down
// We don't care about if they're exec/commit or active/candidate, that all happens in the plugin
func (c *ccipPlugins) Transition(prevPlugins *ccipPlugins) error {
	var err error
	for _, digest := range maps.Keys(*prevPlugins) {
		if *c == nil || (*c)[digest] == nil {
			err = multierr.Append(err, (*(*prevPlugins)[digest]).Close())
		}
	}

	for _, digest := range maps.Keys(*c) {
		if *prevPlugins == nil || (*prevPlugins)[digest] == nil {
			err = multierr.Append(err, (*(*c)[digest]).Start())
		}
	}
	return err
}
