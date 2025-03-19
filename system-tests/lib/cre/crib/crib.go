package crib

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
)

func StartNixShell(input *types.StartNixShellInput) (*nix.NixShell, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, valErr
	}

	globalEnvVars := map[string]string{
		"PROVIDER":           input.InfraInput.CRIB.Provider,
		"DEVSPACE_NAMESPACE": input.InfraInput.CRIB.Namespace,
	}

	for key, value := range input.ExtraEnvVars {
		globalEnvVars[key] = value
	}

	if input.InfraInput.CRIB.Provider == "aws" {
		globalEnvVars["CHAINLINK_TEAM"] = input.InfraInput.CRIB.TeamInput.Team
		globalEnvVars["CHAINLINK_PRODUCT"] = input.InfraInput.CRIB.TeamInput.Product
		globalEnvVars["CHAINLINK_COST_CENTER"] = input.InfraInput.CRIB.TeamInput.CostCenter
		globalEnvVars["CHAINLINK_COMPONENT"] = input.InfraInput.CRIB.TeamInput.Component
	}

	cribConfigDirAbs, absErr := filepath.Abs(filepath.Join(".", input.CribConfigsDir))
	if absErr != nil {
		return nil, errors.Wrapf(absErr, "failed to get absolute path to crib configs dir %s", input.CribConfigsDir)
	}

	globalEnvVars["CRE_CONFIG_DIR"] = cribConfigDirAbs

	// this will run `nix develop`, which will login to all ECRs and set up the environment
	// by running `crib init`
	nixShell, err := nix.NewNixShell(input.InfraInput.CRIB.FolderLocation, globalEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Nix shell")
	}

	// we need to run `devspace purge` to clean up the environment, in case our namespace is already used
	_, err = nixShell.RunCommand("devspace purge")
	if err != nil {
		return nil, errors.Wrap(err, "failed to run devspace purge")
	}

	return nixShell, nil
}

func DeployBlockchain(input *types.DeployCribBlockchainInput) (*blockchain.Output, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, valErr
	}

	gethChainEnvVars := map[string]string{
		"CHAIN_ID": input.BlockchainInput.ChainID,
	}
	_, err := input.NixShell.RunCommandWithEnvVars("devspace run deploy-geth-chain", gethChainEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run devspace run deploy-geth-chain")
	}

	// TODO read family from blockchain input
	blockchainOut, err := infra.ReadBlockchainUrl(filepath.Join(".", input.CribConfigsDir), "evm", input.BlockchainInput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read blockchain URLs")
	}

	return blockchainOut, nil
}

func DeployDons(input *types.DeployCribDonsInput) ([]*types.CapabilitiesAwareNodeSet, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, valErr
	}

	for j, donMetadata := range input.Topology.DonsMetadata {
		deployDonEnvVars := map[string]string{}
		cribConfigsDirAbs := filepath.Join(".", input.CribConfigsDir, donMetadata.Name)
		err := os.MkdirAll(cribConfigsDirAbs, os.ModePerm)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create crib configs directory '%s' for %s", cribConfigsDirAbs, donMetadata.Name)
		}

		// validate that all nodes in the same node set use the same Docker image
		dockerImages := []string{}
		for _, nodeSpec := range input.NodeSetInputs[j].NodeSpecs {
			if nodeSpec.Node.DockerContext != "" {
				return nil, fmt.Errorf("docker context is not supported in CRIB. Please remove docker_ctx from the node spec")
			}
			if nodeSpec.Node.DockerFilePath != "" {
				return nil, fmt.Errorf("dockerfile is not supported in CRIB. Please remove docker_file from the node spec")
			}

			// TODO use kubectl cp to copy them?
			if len(nodeSpec.Node.CapabilitiesBinaryPaths) > 0 {
				return nil, fmt.Errorf("capabilities binaries are not supported in CRIB. Please use a Docker image that already contains the capabilities and remove capabilities_binary_paths from the node spec")
			}
			if nodeSpec.Node.CapabilityContainerDir != "" {
				return nil, fmt.Errorf("capabilities binaries are not supported in CRIB. Please use a Docker image that already contains the capabilities and remove capability_container_dir from the node spec")
			}

			if slices.Contains(dockerImages, nodeSpec.Node.Image) {
				continue

			}
			dockerImages = append(dockerImages, nodeSpec.Node.Image)
		}

		if len(dockerImages) != 1 {
			return nil, fmt.Errorf("all nodes in each nodeset must use the same Docker image, but %d different images were found", len(dockerImages))
		}

		imgTagIndex := strings.LastIndex(dockerImages[0], ":")
		if imgTagIndex == -1 {
			return nil, fmt.Errorf("docker image must have an explicit tag, but it was: %s", dockerImages[0])
		}

		deployDonEnvVars["DEVSPACE_IMAGE"] = dockerImages[0][:imgTagIndex]
		deployDonEnvVars["DEVSPACE_IMAGE_TAG"] = dockerImages[0][imgTagIndex+1:] // +1 to exclude the colon

		bootstrapNodes, err := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.BootstrapNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}

		var cleanToml = func(tomlStr string) ([]byte, error) {
			// unmarshall and marshall to conver it into proper multi-line string
			// that will be correctly serliazed to YAML
			var data interface{}
			tomlErr := toml.Unmarshal([]byte(tomlStr), &data)
			if tomlErr != nil {
				return nil, errors.Wrapf(tomlErr, "failed to unmarshal toml: %s", tomlStr)
			}
			newTOMLBytes, err := toml.Marshal(data)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal toml")
			}

			return newTOMLBytes, nil
		}

		for i, btNode := range bootstrapNodes {
			nodeIndexStr, err := libnode.FindLabelValue(btNode, libnode.IndexKey)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to find node index for bootstrap node %d in nodeset %s", i, donMetadata.Name)
			}

			nodeIndex, err := strconv.Atoi(nodeIndexStr)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert node index '%s' to int for bootstrap node %d in nodeset %s", nodeIndexStr, i, donMetadata.Name)
			}

			cleanToml, tomlErr := cleanToml(input.NodeSetInputs[j].NodeSpecs[nodeIndex].Node.TestConfigOverrides)
			if tomlErr != nil {
				return nil, errors.Wrap(tomlErr, "failed to clean TOML")
			}

			writeErr := os.WriteFile(filepath.Join(cribConfigsDirAbs, fmt.Sprintf("config-override-bt-%d.toml", i)), cleanToml, 0644)
			if writeErr != nil {
				return nil, errors.Wrapf(writeErr, "failed to write config override for bootstrap node %d to file", i)
			}

			writeErr = os.WriteFile(filepath.Join(cribConfigsDirAbs, fmt.Sprintf("secrets-override-bt-%d.toml", i)), []byte(input.NodeSetInputs[j].NodeSpecs[nodeIndex].Node.TestSecretsOverrides), 0644)
			if writeErr != nil {
				return nil, errors.Wrapf(writeErr, "failed to write secrets override for bootstrap node %d to file", i)
			}
		}

		workerNodes, err := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for i, workerNode := range workerNodes {
			nodeIndexStr, err := libnode.FindLabelValue(workerNode, libnode.IndexKey)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to find node index for worker node %d in nodeset %s", i, donMetadata.Name)
			}

			nodeIndex, err := strconv.Atoi(nodeIndexStr)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert node index '%s' to int for worker node %d in nodeset %s", nodeIndexStr, i, donMetadata.Name)
			}

			cleanToml, tomlErr := cleanToml(input.NodeSetInputs[j].NodeSpecs[nodeIndex].Node.TestConfigOverrides)
			if tomlErr != nil {
				return nil, errors.Wrap(tomlErr, "failed to clean TOML")
			}

			writeErr := os.WriteFile(filepath.Join(cribConfigsDirAbs, fmt.Sprintf("config-override-%d.toml", i)), cleanToml, 0644)
			if writeErr != nil {
				return nil, errors.Wrapf(writeErr, "failed to write config override for worker node %d to file", i)
			}

			writeErr = os.WriteFile(filepath.Join(cribConfigsDirAbs, fmt.Sprintf("secrets-override-%d.toml", i)), []byte(input.NodeSetInputs[j].NodeSpecs[nodeIndex].Node.TestSecretsOverrides), 0644)
			if writeErr != nil {
				return nil, errors.Wrapf(writeErr, "failed to write secrets override for worker node %d to file", i)
			}
		}

		deployDonEnvVars["DON_BOOT_NODE_COUNT"] = strconv.Itoa(len(bootstrapNodes))
		deployDonEnvVars["DON_NODE_COUNT"] = fmt.Sprint(len(workerNodes))
		// IMPORTANT: CRIB will deploy gateway only if don_type == "gateway"
		deployDonEnvVars["DON_TYPE"] = donMetadata.Name

		_, err = input.NixShell.RunCommandWithEnvVars("devspace run deploy-don", deployDonEnvVars)
		if err != nil {
			return nil, errors.Wrap(err, "failed to run devspace run deploy-don")
		}

		nsOutput, err := infra.ReadNodeSetUrls(filepath.Join(".", input.CribConfigsDir), donMetadata)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read node set URLs from file")
		}

		input.NodeSetInputs[j].Out = nsOutput
	}

	return input.NodeSetInputs, nil
}

func DeployJd(input *types.DeployCribJdInput) (*jd.Output, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, valErr
	}

	imgTagIndex := strings.LastIndex(input.JDInput.Image, ":")
	if imgTagIndex == -1 {
		return nil, fmt.Errorf("docker image must have an explicit tag, but it was: %s", input.JDInput.Image)
	}

	jdEnvVars := map[string]string{
		"JOB_DISTRIBUTOR_IMAGE_TAG": input.JDInput.Image[imgTagIndex+1:], // +1 to exclude the colon
	}
	_, err := input.NixShell.RunCommandWithEnvVars("devspace run deploy-jd", jdEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run devspace run deploy-jd")
	}

	jdOut, err := infra.ReadJdUrl(filepath.Join(".", input.CribConfigsDir))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read JD URL from file")
	}

	return jdOut, nil
}
