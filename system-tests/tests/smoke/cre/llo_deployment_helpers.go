// Package cre provides helpers for CRE (Chainlink Runtime Environment) testing.
//
// This file provides LLO deployment helpers using data-streams-deploy changesets.
// See https://github.com/smartcontractkit/data-streams-deploy for the changesets.
package cre

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	// Import data-streams-deploy changesets for LLO contract deployment
	channel_config_store "github.com/smartcontractkit/data-streams-deploy/changeset/channel-config-store"
	configurator_v0_5_0 "github.com/smartcontractkit/data-streams-deploy/changeset/configurator/v0_5_0"
	dsdeploy_jobs "github.com/smartcontractkit/data-streams-deploy/jobs"
)

// LLODeploymentResult holds the results of LLO contract deployment
type LLODeploymentResult struct {
	ConfiguratorAddress       common.Address
	ChannelConfigStoreAddress common.Address
	ConfigDigest              [32]byte
	FromBlock                 uint64
}

// DeployLLOContracts deploys Configurator and ChannelConfigStore contracts using data-streams-deploy changesets
func DeployLLOContracts(
	logger zerolog.Logger,
	env cldf.Environment,
	chainSelector uint64,
) (*LLODeploymentResult, error) {
	logger.Info().Msg("Deploying LLO contracts using data-streams-deploy changesets...")

	// Deploy Configurator contract
	logger.Info().Msg("Deploying Configurator contract...")
	configuratorOutput, err := configurator_v0_5_0.DeployConfiguratorChangeset.Apply(env, configurator_v0_5_0.DeployConfiguratorConfig{
		ChainsToDeploy: []uint64{chainSelector},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Configurator: %w", err)
	}

	// Extract Configurator address from datastore
	configuratorAddr, err := extractContractAddress(configuratorOutput.DataStore, chainSelector, "Configurator")
	if err != nil {
		return nil, fmt.Errorf("failed to get Configurator address: %w", err)
	}
	logger.Info().Str("address", configuratorAddr.Hex()).Msg("Configurator deployed")

	// Deploy ChannelConfigStore contract
	logger.Info().Msg("Deploying ChannelConfigStore contract...")
	ccsOutput, err := channel_config_store.DeployChannelConfigStoreChangeset.Apply(env, channel_config_store.DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{chainSelector},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChannelConfigStore: %w", err)
	}

	// Extract ChannelConfigStore address from datastore
	ccsAddr, err := extractContractAddress(ccsOutput.DataStore, chainSelector, "ChannelConfigStore")
	if err != nil {
		return nil, fmt.Errorf("failed to get ChannelConfigStore address: %w", err)
	}
	logger.Info().Str("address", ccsAddr.Hex()).Msg("ChannelConfigStore deployed")

	return &LLODeploymentResult{
		ConfiguratorAddress:       configuratorAddr,
		ChannelConfigStoreAddress: ccsAddr,
		FromBlock:                 0, // Set to actual block number if needed
	}, nil
}

// SetOCRConfig sets the OCR configuration on the Configurator contract
func SetOCRConfig(
	logger zerolog.Logger,
	env cldf.Environment,
	chainSelector uint64,
	configuratorAddr common.Address,
	donConfigID string,
	oracles []confighelper.OracleIdentityExtra,
	f int,
) ([32]byte, error) {
	logger.Info().
		Str("configurator", configuratorAddr.Hex()).
		Int("oracles", len(oracles)).
		Int("f", f).
		Msg("Setting OCR configuration using data-streams-deploy changeset...")

	// Build configuration using the builder pattern from data-streams-deploy
	configParams := configurator_v0_5_0.NewConfiguratorConfig(configurator_v0_5_0.ConfiguratorSetParamsOptions{
		DONConfigID:         &donConfigID,
		ConfiguratorAddress: &configuratorAddr,
		OCROptions: &configurator_v0_5_0.OCR3DataStreamsOptions{
			Oracles: oracles,
			F:       &f,
			S:       []int{len(oracles)},
			// Use sensible defaults for timing
			OnchainConfigOptions: &datastreamsllo.OnchainConfig{
				Version:                 1,
				PredecessorConfigDigest: nil,
			},
			OffchainConfigOptions: datastreamsllo.OffchainConfig{
				ProtocolVersion:                     1,
				DefaultMinReportIntervalNanoseconds: uint64(50 * time.Millisecond),
				EnableObservationCompression:        true,
			},
		},
	})

	_, err := configurator_v0_5_0.SetProductionConfigChangeset.Apply(env, configurator_v0_5_0.SetProductionConfig{
		ConfigurationsByChain: map[uint64][]configurator_v0_5_0.ConfiguratorSetParams{
			chainSelector: {*configParams},
		},
	})
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to set production config: %w", err)
	}

	// Return the config digest (would need to be computed or retrieved from event)
	// For now, return empty - the LLO plugin will compute it from the event
	return [32]byte{}, nil
}

// SetChannelDefinitions sets channel definitions on the ChannelConfigStore
func SetChannelDefinitions(
	logger zerolog.Logger,
	env cldf.Environment,
	chainSelector uint64,
	ccsAddr common.Address,
	donID uint32,
	s3URL string,
	hash [32]byte,
) error {
	logger.Info().
		Str("channelConfigStore", ccsAddr.Hex()).
		Uint32("donID", donID).
		Str("s3URL", s3URL).
		Msg("Setting channel definitions using data-streams-deploy changeset...")

	_, err := channel_config_store.SetChannelDefinitionChangeset.Apply(env, channel_config_store.SetChannelDefinitionsConfig{
		DefinitionsByChain: map[uint64][]channel_config_store.ChannelDefinition{
			chainSelector: {{
				ChannelConfigStore: ccsAddr,
				DonID:              donID,
				S3URL:              s3URL,
				Hash:               hash,
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to set channel definitions: %w", err)
	}

	logger.Info().Msg("Channel definitions set successfully")
	return nil
}

// BuildLLOJobSpec creates an LLO job spec using data-streams-deploy jobs package
func BuildLLOJobSpec(
	name string,
	contractID string,
	relay string,
	chainID string,
	donID uint64,
	configuratorAddr string,
	ccsAddr string,
	bootstrappers []string,
	triggerCapabilityName string,
	triggerCapabilityVersion string,
) *dsdeploy_jobs.LLOJobSpec {
	return &dsdeploy_jobs.LLOJobSpec{
		BaseJobSpec: dsdeploy_jobs.BaseJobSpec{
			Type: "offchainreporting2",
			Name: name,
		},
		ContractID:         contractID,
		P2PV2Bootstrappers: bootstrappers,
		Relay:              dsdeploy_jobs.RelayType(relay),
		PluginType:         "llo",
		RelayConfig: dsdeploy_jobs.RelayConfigLLO{
			ChainID:       chainID,
			LLOConfigMode: "bluegreen",
			LLODonID:      donID,
		},
		PluginConfig: dsdeploy_jobs.PluginConfigLLO{
			ChannelDefinitionsContractAddress: ccsAddr,
			DonID:                             donID,
			Transmitters: []dsdeploy_jobs.Transmitter{{
				Type: "cre",
				Opts: dsdeploy_jobs.TransmitterOpts{
					TriggerCapabilityName:    triggerCapabilityName,
					TriggerCapabilityVersion: triggerCapabilityVersion,
				},
			}},
		},
	}
}

// BuildStreamJobSpec creates a stream job spec with hardcoded values
// Note: streamID must be a top-level field, not under [streamSpec]
// Uses hardcoded values to avoid bridge connectivity issues:
// - Stream 1 (TEST/USD): 424242 for Format 5 magic number
// - Stream 2 (NATIVE/USD): 3000
// - Stream 3 (LINK/USD): 15
// - Stream 4 (DATA/USD): 111111 (base value, multiplied by 5 via calculated stream to get 555555)
func BuildStreamJobSpec(
	name string,
	streamID uint32,
	bridgeName string,
	externalJobID string,
) string {
	// Hardcoded values based on stream ID - using raw magic numbers
	var hardcodedValue int64
	switch streamID {
	case 1: // TEST/USD - Format 5 magic number
		hardcodedValue = 424242
	case 2: // NATIVE/USD
		hardcodedValue = 3000
	case 3: // LINK/USD
		hardcodedValue = 15
	case 4: // DATA/USD - base value 111111 (multiplied by 5 via calculated stream to get 555555)
		hardcodedValue = 111111
	default:
		// Default fallback value
		hardcodedValue = 1000
	}

	// Use memo task to output a hardcoded numeric value, then multiply task to convert to decimal
	// This bypasses the bridge entirely and returns a constant value
	// The top-level streamID field identifies which stream this job belongs to
	// The multiply task is just for format conversion (times=1 means no change, but ensures decimal type)
	return fmt.Sprintf(`
type = "stream"
schemaVersion = 1
name = "%s"
streamID = %d
externalJobID = "%s"
observationSource = """
    result    [type=memo value="%d"];
    multiply  [type=multiply times=1 index=0];
    result    -> multiply;
"""
`, name, streamID, externalJobID, hardcodedValue)
}

// extractContractAddress extracts a contract address from the deployment datastore
func extractContractAddress(dataStore ds.MutableDataStore, chainSelector uint64, contractType string) (common.Address, error) {
	if dataStore == nil {
		return common.Address{}, fmt.Errorf("datastore is nil")
	}

	addresses, err := dataStore.Addresses().Fetch()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	for _, addr := range addresses {
		if addr.ChainSelector == chainSelector && string(addr.Type) == contractType {
			return common.HexToAddress(addr.Address), nil
		}
	}

	return common.Address{}, fmt.Errorf("contract %s not found for chain %d", contractType, chainSelector)
}
