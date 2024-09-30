package launcher

import (
	"fmt"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"go.uber.org/multierr"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
)

// blueGreenDeployment represents a blue-green deployment of OCR instances.
type blueGreenDeployment struct {
	// blue is the blue OCR instance.
	// blue must always be present.
	blue cctypes.CCIPOracle

	// green is the green OCR instance.
	// green may or may not be present.
	// green must never be present if blue is not present.
	// TODO: should we enforce this invariant somehow?
	green cctypes.CCIPOracle
}

// ccipDeployment represents blue-green deployments of both commit and exec
// OCR instances.
type ccipDeployment struct {
	commit blueGreenDeployment
	exec   blueGreenDeployment
}

// Close shuts down all OCR instances in the deployment.
func (c *ccipDeployment) Close() error {
	var err error

	// shutdown blue commit instance.
	err = multierr.Append(err, c.commit.blue.Close())

	// shutdown green commit instance.
	if c.commit.green != nil {
		err = multierr.Append(err, c.commit.green.Close())
	}

	// shutdown blue exec instance.
	err = multierr.Append(err, c.exec.blue.Close())

	// shutdown green exec instance.
	if c.exec.green != nil {
		err = multierr.Append(err, c.exec.green.Close())
	}

	return err
}

// StartBlue starts the blue OCR instances.
func (c *ccipDeployment) StartBlue() error {
	var err error

	err = multierr.Append(err, c.commit.blue.Start())
	err = multierr.Append(err, c.exec.blue.Start())

	return err
}

// CloseBlue shuts down the blue OCR instances.
func (c *ccipDeployment) CloseBlue() error {
	var err error

	err = multierr.Append(err, c.commit.blue.Close())
	err = multierr.Append(err, c.exec.blue.Close())

	return err
}

// HandleBlueGreen handles the blue-green deployment transition.
// prevDeployment is the previous deployment state.
// there are two possible cases:
//
// 1. both blue and green are present in prevDeployment, but only blue is present in c.
// this is a promotion of green to blue, so we need to shut down the blue deployment
// and make green the new blue. In this case green is already running, so there's no
// need to start it. However, we need to shut down the blue deployment.
//
// 2. only blue is present in prevDeployment, both blue and green are present in c.
// In this case, blue is already running, so there's no need to start it. We need to
// start green.
func (c *ccipDeployment) HandleBlueGreen(prevDeployment *ccipDeployment) error {
	if prevDeployment == nil {
		return fmt.Errorf("previous deployment is nil")
	}

	var err error
	if prevDeployment.commit.green != nil && c.commit.green == nil {
		err = multierr.Append(err, prevDeployment.commit.blue.Close())
	} else if prevDeployment.commit.green == nil && c.commit.green != nil {
		err = multierr.Append(err, c.commit.green.Start())
	} else {
		return fmt.Errorf("invalid blue-green deployment transition")
	}

	if prevDeployment.exec.green != nil && c.exec.green == nil {
		err = multierr.Append(err, prevDeployment.exec.blue.Close())
	} else if prevDeployment.exec.green == nil && c.exec.green != nil {
		err = multierr.Append(err, c.exec.green.Start())
	} else {
		return fmt.Errorf("invalid blue-green deployment transition")
	}

	return err
}

// HasGreenInstance returns true if the deployment has a green instance for the
// given plugin type.
func (c *ccipDeployment) HasGreenInstance(pluginType cctypes.PluginType) bool {
	switch pluginType {
	case cctypes.PluginTypeCCIPCommit:
		return c.commit.green != nil
	case cctypes.PluginTypeCCIPExec:
		return c.exec.green != nil
	default:
		return false
	}
}

func isNewGreenInstance(pluginType cctypes.PluginType, ocrConfigs []ccipreaderpkg.OCR3ConfigWithMeta, prevDeployment ccipDeployment) bool {
	return len(ocrConfigs) == 2 && !prevDeployment.HasGreenInstance(pluginType)
}

func isPromotion(pluginType cctypes.PluginType, ocrConfigs []ccipreaderpkg.OCR3ConfigWithMeta, prevDeployment ccipDeployment) bool {
	return len(ocrConfigs) == 1 && prevDeployment.HasGreenInstance(pluginType)
}
