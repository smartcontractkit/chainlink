package environment

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	storage "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/cre_storage_service"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

// Role addresses the extension stacks answer to on ctf, as defined by their CTF components; configs/extensions/README.md
// lists them.
const (
	extensionsSourceName = "cre-extensions"

	workflowSourceRole = storage.WorkflowSourceRole + ":" + storage.NodeServicePort
	artifactAPIRole    = storage.ArtifactAPIRole + ":" + storage.NodeServicePort
	artifactStoreRole  = storage.ArtifactStoreRole + ":" + storage.ArtifactStorePort
)

var artifactStorePort = mustPort(storage.ArtifactStorePort)

// extensionsNodeOverrides points workflow DON nodes at the extension roles: workflows come from cre-workflow-source and
// artifacts from cre-artifact-api / cre-artifact-store. It's appended to user_config_overrides, which the node loads
// last, so the TOML merges with the user's own AdditionalSources instead of being replaced by them.
var extensionsNodeOverrides = fmt.Sprintf(`
[CRE.WorkflowFetcher]
URL = ''

[Capabilities.WorkflowRegistry.WorkflowStorage]
URL = '%s'
TLSEnabled = false
ArtifactStorageHost = '%s'

[[Capabilities.WorkflowRegistry.AdditionalSources]]
URL = '%s'
TLSEnabled = false
Name = '%s'
`, artifactAPIRole, artifactStoreRole, workflowSourceRole, extensionsSourceName)

// extensionsManagedSections are the node config tables extensionsNodeOverrides sets. Like frameworkManagedSections,
// setting them in a workflow nodeset's user_config_overrides would desynchronize the nodes from the extensions.
var extensionsManagedSections = [][]string{
	{"CRE", "WorkflowFetcher"},
	{"Capabilities", "WorkflowRegistry", "WorkflowStorage"},
}

func addExtensionsNodeOverrides(nodeSets []*cre.NodeSet) {
	for _, ns := range nodeSets {
		if !slices.Contains(ns.DONTypes, cre.WorkflowDON) {
			continue
		}
		for _, spec := range ns.NodeSpecs {
			spec.Node.UserConfigOverrides += "\n" + extensionsNodeOverrides
		}
	}
}

func validateExtensionsNodeOverrides(nodeSets []*cre.NodeSet) error {
	for _, ns := range nodeSets {
		if !slices.Contains(ns.DONTypes, cre.WorkflowDON) {
			continue
		}
		for i, spec := range ns.NodeSpecs {
			var parsed map[string]any
			if err := toml.Unmarshal([]byte(spec.Node.UserConfigOverrides), &parsed); err != nil {
				return errors.Wrapf(err, "user_config_overrides of nodeset %q node %d is not valid TOML", ns.Name, i)
			}
			for _, section := range extensionsManagedSections {
				if hasTable(parsed, section) {
					return errors.Errorf("[%s] is set by extensions and cannot be set in user_config_overrides of workflow nodeset %q node %d",
						strings.Join(section, "."), ns.Name, i)
				}
			}
		}
	}
	return nil
}

func hasTable(m map[string]any, path []string) bool {
	for _, key := range path {
		v, ok := m[key]
		if !ok {
			return false
		}
		if m, ok = v.(map[string]any); !ok {
			return true
		}
	}
	return true
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
	if err := validateExtensionsNodeOverrides(in.NodeSets); err != nil {
		return "", err
	}

	return families[0], ensureExtensionsImages(ctx, in.Extensions)
}

// ensureExtensionsImages pulls the stacks' images that aren't local, and names the ones that can't be pulled (images
// built from the stack's repo) instead of failing later in docker compose up.
func ensureExtensionsImages(ctx context.Context, ext *envconfig.Extensions) error {
	cmd := docker(ctx, "compose", "-f", ext.StorageCompose, "config", "--images")
	// the stack's required inputs; their values don't affect the images
	cmd.Env = append(os.Environ(),
		"CRE_LOCAL_REGISTRY_RPC_URL=placeholder",
		"CRE_LOCAL_REGISTRY_CHAIN_SELECTOR=0",
		"CRE_LOCAL_REGISTRY_ADDRESS=placeholder",
		"CRE_LOCAL_DON_FAMILY=placeholder",
	)
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
	in, err := storageInput(ext, output, donFamily)
	if err != nil {
		return err
	}

	_ = storage.Stop()

	out, err := storage.NewWithContext(ctx, in)
	if err != nil {
		return errors.Wrapf(err, "failed to start %s", ext.StorageCompose)
	}
	framework.L.Info().
		Str("cre-deploy-api", out.UserService.ExternalURL).
		Str("cre-workflow-source", out.NodeService.ExternalURL).
		Str("cre-artifact-store", out.ArtifactStore.ExternalURL).
		Msg("Extensions started; reachable from the host at")
	return nil
}

func stopExtensions() error {
	return storage.Stop()
}

func storageInput(ext *envconfig.Extensions, output *creenv.SetupOutput, donFamily string) (*storage.Input, error) {
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

	return &storage.Input{
		ComposeFile:           "file://" + ext.StorageCompose,
		Network:               framework.DefaultNetworkName,
		RegistryRPCURL:        rpcURL,
		RegistryChainSelector: selector,
		RegistryAddress:       capReg.Hex(),
		DonFamily:             donFamily,
	}, nil
}

func mustPort(port string) int {
	p, err := strconv.Atoi(port)
	if err != nil {
		panic(errors.Wrapf(err, "invalid port %q", port))
	}
	return p
}

func docker(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "docker", args...)
}
