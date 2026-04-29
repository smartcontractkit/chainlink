package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

var _ cldf.ChangeSetV2[ProposeSolanaJobSpecInput] = ProposeSolanaJobSpec{}

const solanaNetwork = "solana"

// SolanaOverrideDefaultCfg holds optional per-node overrides for the Solana chain capability JSON config.
// JSON field names match capabilities/chain_capabilities/solana/config.Config.
type SolanaOverrideDefaultCfg struct {
	CREForwarderAddress string        `json:"creForwarderAddress,omitempty" yaml:"creForwarderAddress,omitempty"`
	CREForwarderState   string        `json:"creForwarderState,omitempty" yaml:"creForwarderState,omitempty"`
	Transmitter         string        `json:"transmitter,omitempty" yaml:"transmitter,omitempty"`
	IsLocal             bool          `json:"isLocal,omitempty" yaml:"isLocal,omitempty"`
	Network             string        `json:"network,omitempty" yaml:"network,omitempty"`
	ChainID             string        `json:"chainId,omitempty" yaml:"chainId,omitempty"`
	DeltaStage          time.Duration `json:"deltaStage,omitempty" yaml:"deltaStage,omitempty"`
}

type SolanaCapabilityInput struct {
	NodeID             string                   `json:"nodeID" yaml:"nodeID"`
	Transmitter        string                   `json:"transmitter" yaml:"transmitter"`
	OverrideDefaultCfg SolanaOverrideDefaultCfg `json:"overrideDefaultCfg" yaml:"overrideDefaultCfg"`
}

type ProposeSolanaJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Zone        string `json:"zone" yaml:"zone"`
	Domain      string `json:"domain" yaml:"domain"`
	DONName     string `json:"donName" yaml:"donName"`

	ChainSelector uint64        `json:"chainSelector" yaml:"chainSelector"`
	DeltaStage    time.Duration `json:"deltaStage" yaml:"deltaStage,omitempty"`

	CREForwarderAddress string `json:"creForwarderAddress" yaml:"creForwarderAddress"`
	CREForwarderState   string `json:"creForwarderState" yaml:"creForwarderState"`

	SolanaCapabilityInputs []SolanaCapabilityInput `json:"solanaCapabilityInputs" yaml:"solanaCapabilityInputs"`
}

type ProposeSolanaJobSpec struct{}

func validateSolanaJobCommonFields(environment, domain, zone, donName string, chainSelector uint64, deltaStage time.Duration) error {
	if environment == "" {
		return errors.New("environment is required")
	}
	if domain == "" {
		return errors.New("domain is required")
	}
	if zone == "" {
		return errors.New("zone is required")
	}
	if donName == "" {
		return errors.New("donName is required")
	}
	if chainSelector == 0 {
		return errors.New("chain selector is required")
	}
	if deltaStage <= 0 {
		return fmt.Errorf("deltaStage (%s) must be greater than 0", deltaStage)
	}
	return nil
}

func (u ProposeSolanaJobSpec) VerifyPreconditions(_ cldf.Environment, input ProposeSolanaJobSpecInput) error {
	if len(input.SolanaCapabilityInputs) == 0 {
		return errors.New("at least one solana capability input is required")
	}
	for i, solIn := range input.SolanaCapabilityInputs {
		if solIn.NodeID == "" {
			return fmt.Errorf("nodeID is required for solana capability input at index %d", i)
		}
		if solIn.Transmitter == "" && solIn.OverrideDefaultCfg.Transmitter == "" {
			return fmt.Errorf("transmitter is required for solana capability input at index %d (set transmitter or overrideDefaultCfg.transmitter)", i)
		}
	}

	if err := validateSolanaJobCommonFields(input.Environment, input.Domain, input.Zone, input.DONName, input.ChainSelector, input.DeltaStage); err != nil {
		return err
	}
	if input.CREForwarderAddress == "" {
		return errors.New("creForwarderAddress is required")
	}
	if input.CREForwarderState == "" {
		return errors.New("creForwarderState is required")
	}

	family, err := chainselectors.GetSelectorFamily(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get family for chain selector %d: %w", input.ChainSelector, err)
	}
	if family != chainselectors.FamilySolana {
		return fmt.Errorf("chain selector %d belongs to family %q, expected %q", input.ChainSelector, family, chainselectors.FamilySolana)
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chainID from selector: %w", err)
	}

	for _, solIn := range input.SolanaCapabilityInputs {
		ov := solIn.OverrideDefaultCfg
		if ov.ChainID != "" && ov.ChainID != chainIDStr {
			return fmt.Errorf(
				"chainID in override config (%s) does not match chainID from chain selector (%s) for node %s; "+
					"this field is auto-populated and can be omitted",
				ov.ChainID, chainIDStr, solIn.NodeID,
			)
		}
		if err := validateOverrideNetwork(ov.Network, solanaNetwork, solIn.NodeID); err != nil {
			return err
		}
		if err := validateOverrideForwarder(ov.CREForwarderAddress, input.CREForwarderAddress, solIn.NodeID); err != nil {
			return err
		}
		if err := validateOverrideForwarder(ov.CREForwarderState, input.CREForwarderState, solIn.NodeID); err != nil {
			return err
		}
	}

	return nil
}

func (u ProposeSolanaJobSpec) Apply(e cldf.Environment, input ProposeSolanaJobSpecInput) (cldf.ChangesetOutput, error) {
	chainName, err := chainselectors.GetChainNameFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain name from selector: %w", err)
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	jobName := fmt.Sprintf("solana-cap-v2-%s-%s", chainName, input.Zone)
	job := pkg.StandardCapabilityJob{
		JobName:               jobName,
		Command:               "/usr/local/bin/solana",
		GenerateOracleFactory: false,
	}

	nodeIDToConfig := make(map[string]string, len(input.SolanaCapabilityInputs))
	for _, solIn := range input.SolanaCapabilityInputs {
		if _, exists := nodeIDToConfig[solIn.NodeID]; exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("duplicate nodeID %q in solanaCapabilityInputs", solIn.NodeID)
		}

		cfg := solIn.OverrideDefaultCfg
		if cfg.CREForwarderAddress == "" {
			cfg.CREForwarderAddress = input.CREForwarderAddress
		}
		if cfg.CREForwarderState == "" {
			cfg.CREForwarderState = input.CREForwarderState
		}
		if cfg.Transmitter == "" {
			cfg.Transmitter = solIn.Transmitter
		}
		cfg.ChainID = chainIDStr
		cfg.Network = solanaNetwork
		cfg.DeltaStage = input.DeltaStage
		if solIn.OverrideDefaultCfg.DeltaStage != 0 {
			cfg.DeltaStage = solIn.OverrideDefaultCfg.DeltaStage
		}
		cfg.IsLocal = false // always false for this job

		enc, err := json.Marshal(cfg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal solana cap config: %w", err)
		}

		nodeIDToConfig[solIn.NodeID] = string(enc)
	}

	return proposeAndReport(e, job, nodeIDToConfig, input.Domain, input.DONName, input.Zone)
}
