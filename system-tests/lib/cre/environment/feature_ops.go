package environment

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type PreDONStartupOpDeps struct {
	TestLogger        zerolog.Logger
	CldfEnv           *cldf.Environment
	Provider          infra.Provider
	Topology          *cre.Topology
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	ContractVersions  map[string]string
	CapabilityConfigs cre.CapabilityConfigs
}

type PreDONStartupOpInput struct {
	RegistryChainSelector uint64
	Features              cre.Features `toml:"-" json:"-"` // do not serialize to avoid following error: data cannot be safely written to disk without data lost, avoid type that can't be serialized
}

type PreDONStartupOpOutput struct{}

var PreDONStartupOp = operations.NewOperation(
	"pre-don-startup-op",
	semver.MustParse("1.0.0"),
	"Apply features' logic that needs to be executed before DONs are started",
	func(b operations.Bundle, deps PreDONStartupOpDeps, input PreDONStartupOpInput) (*PreDONStartupOpOutput, error) {
		for _, feature := range input.Features.List() {
			deps.TestLogger.Info().Msgf("Executing PreDONStartup for feature %s", feature.Flag())
			if err := feature.PreDONStartup(
				deps.TestLogger,
				input.RegistryChainSelector,
				deps.CldfEnv,
				deps.Provider,
				deps.Topology,
				deps.BlockchainOutputs,
				deps.CapabilityConfigs,
				deps.ContractVersions,
			); err != nil {
				return nil, fmt.Errorf("failed to execute PreDONStartup for feature %s: %w", feature.Flag(), err)
			}
			deps.TestLogger.Info().Msgf("PreDONStartup for feature %s executed successfully", feature.Flag())
		}

		return &PreDONStartupOpOutput{}, nil
	},
)

type PostDONStartupOpDeps struct {
	TestLogger        zerolog.Logger
	CreEnv            *cre.Environment
	NodeSetOutput     []*cre.WrappedNodeOutput
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	ContractVersions  map[string]string
}

type PostDONStartupOpInput struct {
	Features cre.Features `toml:"-" json:"-"` // do not serialize to avoid following error: data cannot be safely written to disk without data lost, avoid type that can't be serialized
}

type PostDONStartupOpOutput struct {
	DONCapabilityWithConfigs map[int][]keystone_changeset.DONCapabilityWithConfig `toml:"-" json:"-"` // do not serialize to avoid following error: data cannot be safely written to disk without data lost, avoid type that can't be serialized
}

func (c *PostDONStartupOpOutput) MergeWithConfigureInput(configureKeystoneInput *cre.ConfigureKeystoneInput) {
	for k, v := range c.DONCapabilityWithConfigs {
		if configureKeystoneInput.DONCapabilityWithConfigs == nil {
			configureKeystoneInput.DONCapabilityWithConfigs = make(map[int][]keystone_changeset.DONCapabilityWithConfig)
		}

		if configureKeystoneInput.DONCapabilityWithConfigs[k] == nil {
			configureKeystoneInput.DONCapabilityWithConfigs[k] = make([]keystone_changeset.DONCapabilityWithConfig, 0)
		}

		configureKeystoneInput.DONCapabilityWithConfigs[k] = append(configureKeystoneInput.DONCapabilityWithConfigs[k], v...)
	}
}

var PostDONStartupOp = operations.NewOperation(
	"post-don-startup-op",
	semver.MustParse("1.0.0"),
	"Apply features' logic that can be executed after DONs are started",
	func(b operations.Bundle, deps PostDONStartupOpDeps, input PostDONStartupOpInput) (*PostDONStartupOpOutput, error) {
		var postOut *cre.PostDONStartupOutput
		for _, feature := range input.Features.List() {
			var pErr error

			deps.TestLogger.Info().Msgf("Executing PostDONStartup for feature %s", feature.Flag())

			postOut, pErr = feature.PostDONStartup(
				b.GetContext(),
				deps.TestLogger,
				deps.CreEnv,
				deps.NodeSetOutput,
				deps.BlockchainOutputs,
				deps.ContractVersions,
			)

			if pErr != nil {
				return nil, fmt.Errorf("failed to execute PostDONStartup for feature %s: %w", feature.Flag(), pErr)
			}

			deps.TestLogger.Info().Msgf("PostDONStartup for feature %s executed successfully", feature.Flag())
		}

		return &PostDONStartupOpOutput{
			DONCapabilityWithConfigs: postOut.DONCapabilityWithConfigs,
		}, nil
	},
)
