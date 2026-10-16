package environment

import (
	"context"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

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

// Role aliases the extension stacks answer to on ctf. They are a contract with those stacks' compose files; keep them in
// step with configs/extensions/README.md and cre-storage-service local-stack/README.md.
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

// validateExtensions rejects, before the environment starts, extensions config that can't work: a non-Docker provider,
// a missing compose file, other than one workflow DON family, or stack images that are neither local nor pullable. It
// returns the workflow DON family the stacks serve.
func validateExtensions(ctx context.Context, in *envconfig.Config) (string, error) {
	if in.Infra == nil || !in.Infra.IsDocker() {
		return "", errors.New("extensions require the docker infra provider")
	}
	if _, err := os.Stat(in.Extensions.StorageCompose); err != nil {
		return "", errors.Wrap(err, "extensions.storage_compose")
	}

	var families []string
	for _, ns := range in.NodeSets {
		if slices.Contains(ns.DONTypes, cre.WorkflowDON) && len(ns.DonFamilies) > 0 && !slices.Contains(families, ns.DonFamilies[0]) {
			families = append(families, ns.DonFamilies[0])
		}
	}
	if len(families) != 1 {
		return "", errors.Errorf("extensions support exactly one workflow DON family, found %v", families)
	}

	return families[0], ensureExtensionsImages(ctx, in.Extensions)
}

// ensureExtensionsImages pulls the stacks' images that aren't local, and names the ones that can't be pulled (images
// built from the stack's repo) instead of failing later in docker compose up.
func ensureExtensionsImages(ctx context.Context, ext *envconfig.Extensions) error {
	cmd := docker(ctx, "compose", "-f", ext.StorageCompose, "config", "--images")
	cmd.Env = append(os.Environ(), extensionsEnv("http://placeholder:8545", 0, "0x0000000000000000000000000000000000000000", "placeholder")...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return errors.Wrapf(err, "failed to list images in %s", ext.StorageCompose)
	}

	images := strings.Fields(string(out))
	slices.Sort(images)
	var missing []string
	for _, image := range slices.Compact(images) {
		if docker(ctx, "image", "inspect", image).Run() != nil && docker(ctx, "pull", "-q", image).Run() != nil {
			missing = append(missing, image)
		}
	}
	if len(missing) > 0 {
		return errors.Errorf("images %v are neither local nor pullable; build them as described in the README next to %s", missing, ext.StorageCompose)
	}
	return nil
}

// startExtensions starts the extension stacks on the environment's network, with the registry chain, its
// CapabilitiesRegistry and the workflow DON family as inputs.
func startExtensions(ctx context.Context, ext *envconfig.Extensions, donFamily string, output *creenv.SetupOutput) error {
	inputs, err := extensionsInputs(output, donFamily)
	if err != nil {
		return err
	}

	_ = framework.RemoveTestStack(extensionsStackName)

	cmd := docker(ctx, "compose", "-p", extensionsStackName, "-f", ext.StorageCompose, "up", "-d", "--wait")
	cmd.Env = append(os.Environ(), inputs...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return errors.Wrapf(cmd.Run(), "failed to start %s", ext.StorageCompose)
}

func stopExtensions() error {
	return framework.RemoveTestStack(extensionsStackName)
}

func extensionsInputs(output *creenv.SetupOutput, donFamily string) ([]string, error) {
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

	return extensionsEnv(rpcURL, selector, capReg.Hex(), donFamily), nil
}

func extensionsEnv(rpcURL string, selector uint64, registry, donFamily string) []string {
	return []string{
		"CRE_NETWORK=" + framework.DefaultNetworkName,
		"CRE_NETWORK_EXTERNAL=true",
		"CRE_LOCAL_REGISTRY_RPC_URL=" + rpcURL,
		"CRE_LOCAL_REGISTRY_CHAIN_SELECTOR=" + strconv.FormatUint(selector, 10),
		"CRE_LOCAL_REGISTRY_ADDRESS=" + registry,
		"CRE_LOCAL_DON_FAMILY=" + donFamily,
	}
}

func docker(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "docker", args...)
}
