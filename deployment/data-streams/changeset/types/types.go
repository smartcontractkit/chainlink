package types

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// data streams contract types
const (
	ChannelConfigStore cldf.ContractType = "ChannelConfigStore"
	Configurator       cldf.ContractType = "Configurator"
	FeeManager         cldf.ContractType = "FeeManager"
	RewardManager      cldf.ContractType = "RewardManager"
	Verifier           cldf.ContractType = "Verifier"
	VerifierProxy      cldf.ContractType = "VerifierProxy"
)

type (
	MCMSConfig struct {
		MinDelay     time.Duration
		OverrideRoot bool
	}
)

type BillingFeature struct {
	Enabled bool
	Config  *BillingConfig
}

type BillingConfig struct {
	LinkTokenAddress   common.Address
	NativeTokenAddress common.Address
	Surcharge          uint64
}

func (b BillingFeature) Validate() error {
	if b.Enabled && b.Config == nil {
		return errors.New("BillingConfig is required if Billing is enabled")
	}

	if b.Config != nil {
		if b.Config.LinkTokenAddress == (common.Address{}) {
			return errors.New("LinkTokenAddress is required in BillingConfig")
		}
		if b.Config.NativeTokenAddress == (common.Address{}) {
			return errors.New("NativeTokenAddress is required in BillingConfig")
		}
	}

	return nil
}

type OwnershipFeature struct {
	ShouldTransfer     bool // If true,  MCMS takes ownership
	ShouldDeployMCMS   bool
	DeployMCMSConfig   *commontypes.MCMSWithTimelockConfigV2
	MCMSProposalConfig *proposalutils.TimelockConfig
}

func (f OwnershipFeature) Validate() error {
	if f.ShouldTransfer && f.MCMSProposalConfig == nil {
		return errors.New("MCMSProposalConfig is required if ShouldTransfer is true")
	}
	if f.ShouldDeployMCMS && f.DeployMCMSConfig == nil {
		return errors.New("DeployMCMSConfig is required if ShouldDeployMCMS is true")
	}
	return nil
}

func (f OwnershipFeature) AsSettings() OwnershipSettings {
	return OwnershipSettings{
		ShouldTransfer:     f.ShouldTransfer,
		MCMSProposalConfig: f.MCMSProposalConfig,
	}
}

type OwnershipSettings struct {
	ShouldTransfer     bool
	MCMSProposalConfig *proposalutils.TimelockConfig
}
