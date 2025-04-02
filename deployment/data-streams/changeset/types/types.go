package types

import (
	"time"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// data streams contract types
const (
	ChannelConfigStore deployment.ContractType = "ChannelConfigStore"
	Configurator       deployment.ContractType = "Configurator"
	FeeManager         deployment.ContractType = "FeeManager"
	RewardManager      deployment.ContractType = "RewardManager"
	Verifier           deployment.ContractType = "Verifier"
	VerifierProxy      deployment.ContractType = "VerifierProxy"
)

type (
	MCMSConfig struct {
		MinDelay     time.Duration
		OverrideRoot bool
	}
)

type OwnershipFeature struct {
	Transfer           bool // If true, MCMS takes ownership
	DeployMCMS         bool
	DeployMCMSConfig   commontypes.MCMSWithTimelockConfigV2
	MCMSProposalConfig *proposalutils.TimelockConfig
}

func (f OwnershipFeature) AsSettings() OwnershipSettings {
	if f.Transfer && f.MCMSProposalConfig == nil {
		panic("MCMSConfig is required if Transfer is true")
	}
	return OwnershipSettings{
		Transfer:           f.Transfer,
		MCMSProposalConfig: f.MCMSProposalConfig,
	}
}

type OwnershipSettings struct {
	Transfer           bool
	MCMSProposalConfig *proposalutils.TimelockConfig
}
