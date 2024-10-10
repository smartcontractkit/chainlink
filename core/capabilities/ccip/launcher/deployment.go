package launcher

import (
	"fmt"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"go.uber.org/multierr"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
)

// activeCandidateDeployment represents a active-candidate deployment of OCR instances.
type activeCandidateDeployment struct {
	// active is the active OCR instance.
	// active must always be present.
	active cctypes.CCIPOracle

	// candidate is the candidate OCR instance.
	// candidate may or may not be present.
	// candidate must never be present if active is not present.
	candidate cctypes.CCIPOracle
}

// ccipDeployment represents active-candidate deployments of both commit and exec
// OCR instances.
type ccipDeployment struct {
	commit activeCandidateDeployment
	exec   activeCandidateDeployment
}

// Close shuts down all OCR instances in the deployment.
func (c *ccipDeployment) Close() error {
	var err error

	// shutdown active commit instance.
	err = multierr.Append(err, c.commit.active.Close())

	// shutdown candidate commit instance.
	if c.commit.candidate != nil {
		err = multierr.Append(err, c.commit.candidate.Close())
	}

	// shutdown active exec instance.
	err = multierr.Append(err, c.exec.active.Close())

	// shutdown candidate exec instance.
	if c.exec.candidate != nil {
		err = multierr.Append(err, c.exec.candidate.Close())
	}

	return err
}

// StartActive starts the active OCR instances.
func (c *ccipDeployment) StartActive() error {
	var err error

	err = multierr.Append(err, c.commit.active.Start())
	err = multierr.Append(err, c.exec.active.Start())

	return err
}

// CloseActive shuts down the active OCR instances.
func (c *ccipDeployment) CloseActive() error {
	var err error

	err = multierr.Append(err, c.commit.active.Close())
	err = multierr.Append(err, c.exec.active.Close())

	return err
}

// TransitionDeployment handles the active-candidate deployment transition.
// prevDeployment is the previous deployment state.
// there are two possible cases:
//
// 1. both active and candidate are present in prevDeployment, but only active is present in c.
// this is a promotion of candidate to active, so we need to shut down the active deployment
// and make candidate the new active. In this case candidate is already running, so there's no
// need to start it. However, we need to shut down the active deployment.
//
// 2. only active is present in prevDeployment, both active and candidate are present in c.
// In this case, active is already running, so there's no need to start it. We need to
// start candidate.
func (c *ccipDeployment) TransitionDeployment(prevDeployment *ccipDeployment) error {
	if prevDeployment == nil {
		return fmt.Errorf("previous deployment is nil")
	}

	var err error
	if prevDeployment.commit.candidate != nil && c.commit.candidate == nil {
		err = multierr.Append(err, prevDeployment.commit.active.Close())
	} else if prevDeployment.commit.candidate == nil && c.commit.candidate != nil {
		err = multierr.Append(err, c.commit.candidate.Start())
	} else {
		return fmt.Errorf("invalid active-candidate deployment transition")
	}

	if prevDeployment.exec.candidate != nil && c.exec.candidate == nil {
		err = multierr.Append(err, prevDeployment.exec.active.Close())
	} else if prevDeployment.exec.candidate == nil && c.exec.candidate != nil {
		err = multierr.Append(err, c.exec.candidate.Start())
	} else {
		return fmt.Errorf("invalid active-candidate deployment transition")
	}

	return err
}

// HasCandidateInstance returns true if the deployment has a candidate instance for the
// given plugin type.
func (c *ccipDeployment) HasCandidateInstance(pluginType cctypes.PluginType) bool {
	switch pluginType {
	case cctypes.PluginTypeCCIPCommit:
		return c.commit.candidate != nil
	case cctypes.PluginTypeCCIPExec:
		return c.exec.candidate != nil
	default:
		return false
	}
}

func isNewCandidateInstance(pluginType cctypes.PluginType, ocrConfigs []ccipreaderpkg.OCR3ConfigWithMeta, prevDeployment ccipDeployment) bool {
	return len(ocrConfigs) == 2 && !prevDeployment.HasCandidateInstance(pluginType)
}

func isPromotion(pluginType cctypes.PluginType, ocrConfigs []ccipreaderpkg.OCR3ConfigWithMeta, prevDeployment ccipDeployment) bool {
	return len(ocrConfigs) == 1 && prevDeployment.HasCandidateInstance(pluginType)
}
