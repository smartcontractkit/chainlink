package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink/deployment"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// DeployAptosChainConfig is a configuration for deploying CCIP Package for Aptos chains
type DeployAptosChainConfig struct {
	ContractParamsPerChain map[uint64]ChainContractParams
}

func (c DeployAptosChainConfig) Validate() error {
	for cs, args := range c.ContractParamsPerChain {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
		if err := args.Validate(); err != nil {
			return fmt.Errorf("invalid contract args for chain %d: %w", cs, err)
		}
	}
	return nil
}

// ChainContractParams stores configuration to call initialize in CCIP contracts
type ChainContractParams struct {
	FeeQuoterParams FeeQuoterParams
	OffRampParams   OffRampParams
	OnRampParams    OnRampParams
}

func (c ChainContractParams) Validate() error {
	// Validate every field
	if err := c.FeeQuoterParams.Validate(); err != nil {
		return fmt.Errorf("invalid FeeQuoterParams: %w", err)
	}
	if err := c.OffRampParams.Validate(); err != nil {
		return fmt.Errorf("invalid OffRampParams: %w", err)
	}
	if err := c.OnRampParams.Validate(); err != nil {
		return fmt.Errorf("invalid OnRampParams: %w", err)
	}
	return nil
}

type FeeQuoterParams struct {
	MaxFeeJuelsPerMsg            uint64
	LinkToken                    aptos.AccountAddress
	TokenPriceStalenessThreshold uint64
	FeeTokens                    []aptos.AccountAddress
}

func (f FeeQuoterParams) Validate() error {
	if f.MaxFeeJuelsPerMsg < 0 {
		return fmt.Errorf("MaxFeeJuelsPerMsg must be positive")
	}
	if (f.LinkToken == aptos.AccountAddress{}) {
		return fmt.Errorf("LinkToken is required")
	}
	if f.TokenPriceStalenessThreshold == 0 {
		return fmt.Errorf("TokenPriceStalenessThreshold can't be 0")
	}
	if len(f.FeeTokens) == 0 {
		return fmt.Errorf("at least one FeeTokens is required")
	}
	return nil
}

type OffRampParams struct {
	ChainSelector                    uint64
	PermissionlessExecutionThreshold uint32
	IsRMNVerificationDisabled        bool
	SourceChainSelectors             []uint64
	SourceChainIsEnabled             []bool
}

func (o OffRampParams) Validate() error {
	if err := deployment.IsValidChainSelector(o.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", o.ChainSelector, err)
	}
	if o.PermissionlessExecutionThreshold == 0 {
		return fmt.Errorf("PermissionlessExecutionThreshold can't be 0")
	}
	if len(o.SourceChainSelectors) != len(o.SourceChainIsEnabled) {
		return fmt.Errorf("SourceChainSelectors and SourceChainIsEnabled must have the same length")
	}
	return nil
}

type OnRampParams struct {
	ChainSelector             uint64
	AllowlistAdmin            aptos.AccountAddress
	DestChainSelectors        []uint64
	DestChainEnabled          []bool
	DestChainAllowlistEnabled []bool
}

func (o OnRampParams) Validate() error {
	if err := deployment.IsValidChainSelector(o.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", o.ChainSelector, err)
	}
	if (o.AllowlistAdmin == aptos.AccountAddress{}) {
		return fmt.Errorf("AllowlistAdmin is required")
	}
	if len(o.DestChainSelectors) != len(o.DestChainEnabled) {
		return fmt.Errorf("DestChainSelectors and DestChainEnabled must have the same length")
	}
	if len(o.DestChainSelectors) != len(o.DestChainAllowlistEnabled) {
		return fmt.Errorf("DestChainSelectors and DestChainAllowlistEnabled must have the same length")
	}
	return nil
}

// DeployAptosMCMSConfig is a configuration for deploying MCMS contracts for Aptos chains
type DeployAptosMCMSConfig struct {
	MCMSConfigPerChain map[uint64]mcmstypes.Config
}
