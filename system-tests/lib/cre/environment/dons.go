package environment

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	creflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type StartedDON struct {
	NodeOutput *cre.WrappedNodeOutput
	DON        *cre.DON
}

type StartedDONs []*StartedDON

func (s *StartedDONs) NodeOutputs() []*cre.WrappedNodeOutput {
	outputs := make([]*cre.WrappedNodeOutput, len(*s))
	for idx, don := range *s {
		outputs[idx] = don.NodeOutput
	}
	return outputs
}

func (s *StartedDONs) DONs() []*cre.DON {
	dons := make([]*cre.DON, len(*s))
	for idx, don := range *s {
		dons[idx] = don.DON
	}
	return dons
}

func StartDONs(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	infraInput infra.Provider,
	registryChainBlockchainOutput *blockchain.Output,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
) (*StartedDONs, error) {
	if infraInput.Type == infra.CRIB {
		lggr.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &cre.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  capabilitiesAwareNodeSets,
			CribConfigsDir: cribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		var devspaceErr error
		capabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with crib-sdk")
		}
	}

	for donIdx, donMetadata := range topology.DonsMetadata.List() {
		if !copyCapabilityBinaries {
			continue
		}

		customBinariesPaths := make(map[cre.CapabilityFlag]string)
		for flag, config := range capabilityConfigs {
			if creflags.HasFlagForAnyChain(donMetadata.Flags, flag) && config.BinaryPath != "" {
				customBinariesPaths[flag] = config.BinaryPath
			}
		}

		executableErr := crecapabilities.MakeBinariesExecutable(customBinariesPaths)
		if executableErr != nil {
			return nil, pkgerrors.Wrap(executableErr, "failed to make binaries executable")
		}

		var err error
		ns, err := crecapabilities.AppendBinariesPathsNodeSpec(capabilitiesAwareNodeSets[donIdx], donMetadata, customBinariesPaths)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
		}
		capabilitiesAwareNodeSets[donIdx] = ns
	}

	// Add env vars, which were provided programmatically, to the node specs
	// or fail, if node specs already had some env vars set in the TOML config
	for donIdx, donMetadata := range topology.DonsMetadata.List() {
		hasEnvVarsInTomlConfig := false
		for nodeIdx, nodeSpec := range capabilitiesAwareNodeSets[donIdx].NodeSpecs {
			if len(nodeSpec.Node.EnvVars) > 0 {
				hasEnvVarsInTomlConfig = true
				break
			}

			capabilitiesAwareNodeSets[donIdx].NodeSpecs[nodeIdx].Node.EnvVars = capabilitiesAwareNodeSets[donIdx].EnvVars
		}

		if hasEnvVarsInTomlConfig && len(capabilitiesAwareNodeSets[donIdx].EnvVars) > 0 {
			return nil, fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
		}
	}

	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range capabilitiesAwareNodeSets {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range capabilitiesAwareNodeSets[i].NodeSpecs {
				capabilitiesAwareNodeSets[i].NodeSpecs[j].Node.Image = image
				// unset docker context and file path, so that we can use the image from the registry
				capabilitiesAwareNodeSets[i].NodeSpecs[j].Node.DockerContext = ""
				capabilitiesAwareNodeSets[i].NodeSpecs[j].Node.DockerFilePath = ""
			}
		}
	}

	errGroup, _ := errgroup.WithContext(ctx)
	var resultMap sync.Map

	for idx, nodeSetInput := range capabilitiesAwareNodeSets {
		errGroup.Go(func() error {
			startTime := time.Now()
			lggr.Info().Msgf("Starting DON named %s", nodeSetInput.Name)
			nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, registryChainBlockchainOutput)
			if nodesetErr != nil {
				return pkgerrors.Wrapf(nodesetErr, "failed to start nodeSet named %s", nodeSetInput.Name)
			}

			don, donErr := cre.NewDON(ctx, topology.DonsMetadata.List()[idx], nodeset.CLNodes)
			if donErr != nil {
				return pkgerrors.Wrapf(donErr, "failed to create DON from node set named %s", nodeSetInput.Name)
			}

			resultMap.Store(idx, &StartedDON{
				NodeOutput: &cre.WrappedNodeOutput{
					Output:       nodeset,
					NodeSetName:  nodeSetInput.Name,
					Capabilities: nodeSetInput.ComputedCapabilities,
				},
				DON: don,
			})

			lggr.Info().Msgf("DON %s started in %.2f seconds", nodeSetInput.Name, time.Since(startTime).Seconds())

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	startedDONs := make(StartedDONs, len(capabilitiesAwareNodeSets))
	resultMap.Range(func(key, value any) bool {
		// key is index in the original slice
		startedDONs[key.(int)] = value.(*StartedDON)
		return true
	})

	return &startedDONs, nil
}
