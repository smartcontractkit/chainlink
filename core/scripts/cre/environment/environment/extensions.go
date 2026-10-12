package environment

import (
	"context"
	"os"
	"os/exec"
	"slices"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const (
	extensionsStackName = "cre-extensions"

	workflowSourceRole = "cre-workflow-source:50052"
	artifactAPIRole    = "cre-artifact-api:50052"
	artifactStorePort  = 4566
)

var artifactStoreRole = "cre-artifact-store:" + strconv.Itoa(artifactStorePort)

// extensionsNodeConfig points workflow DON nodes at the extension roles: workflows come from cre-workflow-source and
// artifacts from cre-artifact-api / cre-artifact-store.
func extensionsNodeConfig(input cre.GenerateConfigsInput, configs cre.NodeIndexToConfigOverride) (cre.NodeIndexToConfigOverride, error) {
	if !input.DonMetadata.HasFlag(cre.WorkflowDON) {
		return configs, nil
	}

	out := make(cre.NodeIndexToConfigOverride, len(configs))
	for idx, raw := range configs {
		var c corechainlink.Config
		if err := toml.Unmarshal([]byte(raw), &c); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal config of node %d in DON %s", idx, input.DonMetadata.Name)
		}

		c.CRE.WorkflowFetcher = &coretoml.WorkflowFetcherConfig{URL: new("")}
		c.Capabilities.WorkflowRegistry.WorkflowStorage = coretoml.WorkflowStorage{
			URL:                 new(artifactAPIRole),
			TLSEnabled:          new(false),
			ArtifactStorageHost: new(artifactStoreRole),
		}
		c.Capabilities.WorkflowRegistry.AdditionalSourcesConfig = append(c.Capabilities.WorkflowRegistry.AdditionalSourcesConfig, coretoml.AdditionalWorkflowSource{
			URL:        new(workflowSourceRole),
			TLSEnabled: new(false),
			Name:       new(extensionsStackName),
		})

		marshalled, err := toml.Marshal(c)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal config of node %d in DON %s", idx, input.DonMetadata.Name)
		}
		out[idx] = string(marshalled)
	}
	return out, nil
}

// startExtensions starts the extension stacks on the environment's network, with the registry chain, its
// CapabilitiesRegistry and the workflow DON family as inputs.
func startExtensions(ctx context.Context, ext *envconfig.Extensions, output *creenv.SetupOutput) error {
	inputs, err := extensionsInputs(output)
	if err != nil {
		return err
	}

	_ = framework.RemoveTestStack(extensionsStackName)

	cmd := exec.CommandContext(ctx, "docker", "compose", "-p", extensionsStackName, "-f", ext.StorageCompose, "up", "-d", "--wait") //nolint:gosec // G204: compose file path from the developer's TOML config
	cmd.Env = append(os.Environ(), inputs...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return errors.Wrapf(cmd.Run(), "failed to start %s", ext.StorageCompose)
}

func stopExtensions() error {
	return framework.RemoveTestStack(extensionsStackName)
}

func extensionsInputs(output *creenv.SetupOutput) ([]string, error) {
	creEnv := output.CreEnvironment
	selector := creEnv.RegistryChainSelector

	var rpcURL string
	for _, bc := range creEnv.Blockchains {
		if bc.ChainSelector() == selector && len(bc.CtfOutput().Nodes) > 0 {
			rpcURL = bc.CtfOutput().Nodes[0].InternalHTTPUrl
		}
	}
	if rpcURL == "" {
		return nil, errors.Errorf("no RPC URL for registry chain %d", selector)
	}

	capReg := contracts.MightGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, selector,
		keystone_changeset.CapabilitiesRegistry.String(), creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	if capReg == nil {
		return nil, errors.Errorf("no CapabilitiesRegistry on registry chain %d", selector)
	}

	var families []string
	for _, don := range output.Dons.DonsWithFlag(cre.WorkflowDON) {
		if f := don.DonFamily(); f != "" && !slices.Contains(families, f) {
			families = append(families, f)
		}
	}
	if len(families) != 1 {
		return nil, errors.Errorf("extensions support exactly one workflow DON family, found %v", families)
	}

	return []string{
		"CRE_NETWORK=" + framework.DefaultNetworkName,
		"CRE_NETWORK_EXTERNAL=true",
		"CRE_LOCAL_REGISTRY_RPC_URL=" + rpcURL,
		"CRE_LOCAL_REGISTRY_CHAIN_SELECTOR=" + strconv.FormatUint(selector, 10),
		"CRE_LOCAL_REGISTRY_ADDRESS=" + capReg.Hex(),
		"CRE_LOCAL_DON_FAMILY=" + families[0],
	}, nil
}
