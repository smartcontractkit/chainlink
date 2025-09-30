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

type PostDONStartupOpDeps struct {
	TestLogger       zerolog.Logger
	CreEnv           *cre.Environment
	NodeSetOutput    []*cre.WrappedNodeOutput
	ContractVersions map[string]string
}

type PostDONStartupOpInput struct {
	Features []cre.Feature
}

type PostDONStartupOpOutput struct {
	DONCapabilityWithConfigs map[int][]keystone_changeset.DONCapabilityWithConfig `toml:"-" json:"-"`
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
	"Apply features that require DONs to be started",
	func(b operations.Bundle, deps PostDONStartupOpDeps, input PostDONStartupOpInput) (*PostDONStartupOpOutput, error) {
		var postOut *cre.PostDONStartupOutput
		for _, feature := range input.Features {
			var pErr error
			postOut, pErr = feature.PostDONStartup(
				b.GetContext(),
				deps.TestLogger,
				deps.CreEnv,
				deps.NodeSetOutput,
				deps.ContractVersions,
			)

			if pErr != nil {
				return nil, fmt.Errorf("failed to execute PostDONStartup for feature %s: %w", feature.Flag(), pErr)
			}
		}

		return &PostDONStartupOpOutput{
			DONCapabilityWithConfigs: postOut.DONCapabilityWithConfigs,
		}, nil
	},
)

type PreDONStartupOpDeps struct {
	CldfEnv           *cldf.Environment
	Provider          infra.Provider
	NodeSetOutput     []*cre.CapabilitiesAwareNodeSet
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	ContractVersions  map[string]string
	CapabilityConfigs cre.CapabilityConfigs
}

type PreDONStartupOpInput struct {
	RegistryChainSelector uint64
	Features              []cre.Feature
}

type PreDONStartupOpOutput struct{}

var PreDONStartupOp = operations.NewOperation(
	"pre-don-startup-op",
	semver.MustParse("1.0.0"),
	"Apply features that do not require DONs to be started",
	func(b operations.Bundle, deps PreDONStartupOpDeps, input PreDONStartupOpInput) (*PreDONStartupOpOutput, error) {

		for _, feature := range input.Features {
			if err := feature.PreDONStartup(
				input.RegistryChainSelector,
				deps.CldfEnv,
				deps.Provider,
				deps.NodeSetOutput,
				deps.BlockchainOutputs,
				deps.CapabilityConfigs,
			); err != nil {
				return nil, fmt.Errorf("failed to execute PreDONStartup for feature %s: %w", feature.Flag(), err)
			}
		}

		return &PreDONStartupOpOutput{}, nil
	},
)
