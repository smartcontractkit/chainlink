package environment

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	creflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func StartJD(lggr zerolog.Logger, jdInput jd.Input, infraInput infra.Provider) (*jd.Output, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")

	var jdOutput *jd.Output
	if infraInput.Type == infra.CRIB {
		deployCribJdInput := &cre.DeployCribJdInput{
			JDInput:        jdInput,
			CribConfigsDir: cribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		var jdErr error
		jdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
		}
	}

	var jdErr error
	jdOutput, jdErr = CreateJobDistributor(jdInput)
	if jdErr != nil {
		jdErr = fmt.Errorf("failed to start JD container for image %s: %w", jdInput.Image, jdErr)

		// useful end user messages
		if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
			jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
		}
		return nil, jdErr
	}

	lggr.Info().Msgf("Job Distributor started in %.2f seconds", time.Since(startTime).Seconds())

	return jdOutput, nil
}

func CreateJobDistributor(input jd.Input) (*jd.Output, error) {
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(E2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(E2eJobDistributorVersionEnvVarName)
		input.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}

	jdOutput, err := jd.NewJD(&input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create new job distributor")
	}

	return jdOutput, nil
}

func StartDONsAndJD(
	lggr zerolog.Logger,
	jdInput *jd.Input,
	registryChainBlockchainOutput *blockchain.Output,
	topology *cre.Topology,
	provider infra.Provider,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
) (*jd.Output, []*cre.WrappedNodeOutput, error) {
	if jdInput == nil {
		return nil, nil, errors.New("jd input is nil")
	}
	if registryChainBlockchainOutput == nil {
		return nil, nil, errors.New("registry chain blockchain output is nil")
	}
	if topology == nil {
		return nil, nil, errors.New("topology is nil")
	}
	var jdOutput *jd.Output
	jdAndDonsErrGroup := &errgroup.Group{}

	jdAndDonsErrGroup.Go(func() error {
		var startJDErr error
		jdOutput, startJDErr = StartJD(lggr, *jdInput, provider)
		if startJDErr != nil {
			return pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}

		return nil
	})

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
			return nil, nil, pkgerrors.Wrap(executableErr, "failed to make binaries executable")
		}

		var err error
		ns, err := crecapabilities.AppendBinariesPathsNodeSpec(capabilitiesAwareNodeSets[donIdx], donMetadata, customBinariesPaths)
		if err != nil {
			return nil, nil, pkgerrors.Wrapf(err, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
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
			return nil, nil, fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
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

	nodeSetOutput := make([]*cre.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))
	jdAndDonsErrGroup.Go(func() error {
		var startDonsErr error
		nodeSetOutput, startDonsErr = StartDONs(lggr, topology, provider, registryChainBlockchainOutput, capabilitiesAwareNodeSets)
		if startDonsErr != nil {
			return pkgerrors.Wrap(startDonsErr, "failed to start DONs")
		}

		return nil
	})

	if jdAndDonErr := jdAndDonsErrGroup.Wait(); jdAndDonErr != nil {
		return nil, nil, pkgerrors.Wrap(jdAndDonErr, "failed to start Job Distributor or DONs")
	}

	return jdOutput, nodeSetOutput, nil
}
