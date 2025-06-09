package config

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// DeployAptosChainConfig is a configuration for deploying CCIP Package for Aptos chains
type DeployAptosChainConfig struct {
	MCMSDeployConfigPerChain   map[uint64]types.MCMSWithTimelockConfigV2
	MCMSTimelockConfigPerChain map[uint64]proposalutils.TimelockConfig
	ContractParamsPerChain     map[uint64]ChainContractParams
}

func (c DeployAptosChainConfig) Validate() error {
	for cs, args := range c.ContractParamsPerChain {
		if err := cldf.IsValidChainSelector(cs); err != nil {
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
	TokenPriceStalenessThreshold uint64
	FeeTokens                    []aptos.AccountAddress
}

func (f FeeQuoterParams) Validate() error {
	if f.TokenPriceStalenessThreshold == 0 {
		return errors.New("TokenPriceStalenessThreshold can't be 0")
	}
	return nil
}

type OffRampParams struct {
	ChainSelector                    uint64
	PermissionlessExecutionThreshold uint32
	IsRMNVerificationDisabled        []bool
	SourceChainSelectors             []uint64
	SourceChainIsEnabled             []bool
	SourceChainsOnRamp               [][]byte
}

func (o OffRampParams) Validate() error {
	if err := cldf.IsValidChainSelector(o.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", o.ChainSelector, err)
	}
	if o.PermissionlessExecutionThreshold == 0 {
		return errors.New("PermissionlessExecutionThreshold can't be 0")
	}
	if len(o.SourceChainSelectors) != len(o.SourceChainIsEnabled) {
		return errors.New("SourceChainSelectors and SourceChainIsEnabled must have the same length")
	}
	return nil
}

type OnRampParams struct {
	ChainSelector  uint64
	AllowlistAdmin aptos.AccountAddress
	FeeAggregator  aptos.AccountAddress
}

func (o OnRampParams) Validate() error {
	if err := cldf.IsValidChainSelector(o.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", o.ChainSelector, err)
	}
	return nil
}
