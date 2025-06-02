package crib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	crecaps "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func StartNixShell(input *types.StartNixShellInput) (*nix.Shell, error) {
	if input == nil {
		return nil, errors.New("StartNixShellInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	globalEnvVars := map[string]string{
		"PROVIDER":           input.InfraInput.CRIB.Provider,
		"DEVSPACE_NAMESPACE": input.InfraInput.CRIB.Namespace,
	}

	for key, value := range input.ExtraEnvVars {
		globalEnvVars[key] = value
	}

	if strings.EqualFold(input.InfraInput.CRIB.Provider, libtypes.AWS) {
		globalEnvVars["CHAINLINK_TEAM"] = input.InfraInput.CRIB.TeamInput.Team
		globalEnvVars["CHAINLINK_PRODUCT"] = input.InfraInput.CRIB.TeamInput.Product
		globalEnvVars["CHAINLINK_COST_CENTER"] = input.InfraInput.CRIB.TeamInput.CostCenter
		globalEnvVars["CHAINLINK_COMPONENT"] = input.InfraInput.CRIB.TeamInput.Component
	}

	cribConfigDirAbs, absErr := filepath.Abs(filepath.Join(".", input.CribConfigsDir))
	if absErr != nil {
		return nil, errors.Wrapf(absErr, "failed to get absolute path to crib configs dir %s", input.CribConfigsDir)
	}

	globalEnvVars["CONFIG_OVERRIDES_DIR"] = cribConfigDirAbs

	// this will run `nix develop`, which will login to all ECRs and set up the environment
	// by running `crib init`
	nixShell, err := nix.NewNixShell(input.InfraInput.CRIB.FolderLocation, globalEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Nix shell")
	}

	if input.PurgeNamespace {
		// we run `devspace purge` to clean up the environment, in case our namespace is already used
		_, err = nixShell.RunCommand("devspace purge --no-warn")
		if err != nil {
			return nil, errors.Wrap(err, "failed to run devspace purge")
		}
	}

	return nixShell, nil
}

func DeployBlockchain(input *types.DeployCribBlockchainInput) (*blockchain.Output, error) {
	if input == nil {
		return nil, errors.New("DeployCribBlockchainInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	gethChainEnvVars := map[string]string{
		"CHAIN_ID": input.BlockchainInput.ChainID,
	}

	// crib init is required to create necessary Kubernetes resources (like configmaps) and set up proper permissions
	_, err := input.NixShell.RunCommand("crib init")
	if err != nil {
		return nil, errors.Wrap(err, "failed to run crib init")
	}

	_, err = input.NixShell.RunCommandWithEnvVars("devspace run deploy-custom-geth-chain --no-warn", gethChainEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run devspace run deploy-custom-geth-chain --no-warn")
	}

	// TODO chain family should be dynamic, but currently we don't have in the input (it's set in the output depending on blockchain type)
	blockchainOut, err := infra.ReadBlockchainURL(filepath.Join(".", input.CribConfigsDir), "evm", input.BlockchainInput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read blockchain URLs")
	}

	return blockchainOut, nil
}

func DeployDons(input *types.DeployCribDonsInput) ([]*types.CapabilitiesAwareNodeSet, error) {
	if input == nil {
		return nil, errors.New("DeployCribDonsInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	for j, donMetadata := range input.Topology.DonsMetadata {
		deployDonEnvVars := map[string]string{}
		cribConfigsDirAbs := filepath.Join(".", input.CribConfigsDir, donMetadata.Name)
		err := os.MkdirAll(cribConfigsDirAbs, os.ModePerm)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create crib configs directory '%s' for %s", cribConfigsDirAbs, donMetadata.Name)
		}

		// validate that all nodes in the same node set use the same Docker image
		dockerImage, dockerImagesErr := nodesetDockerImage(input.NodeSetInputs[j])
		if dockerImagesErr != nil {
			return nil, errors.Wrap(dockerImagesErr, "failed to validate node set Docker images")
		}

		imageName, imageErr := dockerImageName(dockerImage)
		if imageErr != nil {
			return nil, errors.Wrap(imageErr, "failed to get image name")
		}

		imageTag, imageErr := dockerImageTag(dockerImage)
		if imageErr != nil {
			return nil, errors.Wrap(imageErr, "failed to get image tag")
		}

		deployDonEnvVars["DEVSPACE_IMAGE"] = imageName
		deployDonEnvVars["DEVSPACE_IMAGE_TAG"] = imageTag

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
			newTOMLBytes, marshallErr := toml.Marshal(data)
			if marshallErr != nil {
				return nil, errors.Wrap(marshallErr, "failed to marshal toml")
			}

			return newTOMLBytes, nil
		}

		var writeOverrides = func(nodeMetadata *types.NodeMetadata, i int, nodeType types.NodeType) error {
			nodeIndexStr, findErr := libnode.FindLabelValue(nodeMetadata, libnode.IndexKey)
			if findErr != nil {
				return errors.Wrapf(findErr, "failed to find node index for %s node %d in nodeset %s", nodeType, i, donMetadata.Name)
			}

			nodeIndex, convErr := strconv.Atoi(nodeIndexStr)
			if convErr != nil {
				return errors.Wrapf(convErr, "failed to convert node index '%s' to int for %s node %d in nodeset %s", nodeIndexStr, nodeType, i, donMetadata.Name)
			}

			cleanToml, tomlErr := cleanToml(input.NodeSetInputs[j].NodeSpecs[nodeIndex].Node.TestConfigOverrides)
			if tomlErr != nil {
				return errors.Wrap(tomlErr, "failed to clean TOML")
			}

			configFileMask := "config-override-bt-%d.toml"
			secretsFileMask := "secrets-override-bt-%d.toml"

			if nodeType != types.BootstrapNode {
				configFileMask = "config-override-%d.toml"
				secretsFileMask = "secrets-override-%d.toml"
			}

			writeErr := os.WriteFile(filepath.Join(cribConfigsDirAbs, fmt.Sprintf(configFileMask, i)), cleanToml, 0600)
			if writeErr != nil {
				return errors.Wrapf(writeErr, "failed to write config override for bootstrap node %d to file", i)
			}

			writeErr = os.WriteFile(filepath.Join(cribConfigsDirAbs, fmt.Sprintf(secretsFileMask, i)), []byte(input.NodeSetInputs[j].NodeSpecs[nodeIndex].Node.TestSecretsOverrides), 0600)
			if writeErr != nil {
				return errors.Wrapf(writeErr, "failed to write secrets override for bootstrap node %d to file", i)
			}

			return nil
		}

		for i, btNode := range bootstrapNodes {
			writeErr := writeOverrides(btNode, i, types.BootstrapNode)
			if writeErr != nil {
				return nil, writeErr
			}
		}

		workerNodes, err := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for i, workerNode := range workerNodes {
			writeErr := writeOverrides(workerNode, i, types.WorkerNode)
			if writeErr != nil {
				return nil, writeErr
			}
		}

		deployDonEnvVars["DON_BOOT_NODE_COUNT"] = strconv.Itoa(len(bootstrapNodes))
		deployDonEnvVars["DON_NODE_COUNT"] = strconv.Itoa(len(workerNodes))
		// IMPORTANT: CRIB will deploy gateway only if don_type == "gateway", in other cases the DON_TYPE value has no other impact than being uses in release/service/etc names
		deployDonEnvVars["DON_TYPE"] = donMetadata.Name

		_, deployErr := input.NixShell.RunCommandWithEnvVars("devspace run deploy-don --no-warn", deployDonEnvVars)
		if deployErr != nil {
			return nil, errors.Wrap(deployErr, "failed to run devspace run deploy-don")
		}

		// validate capabilities-related configuration and copy capabilities to pods
		podNamePattern := input.NodeSetInputs[j].Name + `-\\d+`
		_, regErr := regexp.Compile(podNamePattern)
		if regErr != nil {
			return nil, errors.Wrapf(regErr, "failed to compile regex for pod name pattern %s", podNamePattern)
		}
		capabilitiesFound := map[string]int{}
		capabilitiesDirs := []string{}
		capabilitiesDirsFound := map[string]int{}

		// make sure all worker nodes in DON have the same set of capabilities
		// in the future we might want to allow different capabilities for different nodes
		// but for now we require all worker nodes in the same DON to have the same capabilities
		for _, nodeSpec := range input.NodeSetInputs[j].NodeSpecs {
			for _, capabilityBinaryPath := range nodeSpec.Node.CapabilitiesBinaryPaths {
				capabilitiesFound[capabilityBinaryPath]++
			}

			if nodeSpec.Node.CapabilityContainerDir != "" {
				capabilitiesDirs = append(capabilitiesDirs, nodeSpec.Node.CapabilityContainerDir)
				capabilitiesDirsFound[nodeSpec.Node.CapabilityContainerDir]++
			}
		}

		for capability, count := range capabilitiesFound {
			// we only care about worker nodes, because bootstrap nodes cannot execute any workflows, so they don't need capabilities
			if count != len(workerNodes) {
				return nil, fmt.Errorf("capability %s wasn't defined for all worker nodes in nodeset %s. All worker nodes in the same nodeset must have the same capabilities", capability, input.NodeSetInputs[j].Name)
			}
		}

		destinationDir, err := crecaps.DefaultContainerDirectory(libtypes.CRIB)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get default directory for capabilities in CRIB")
		}

		// all of them need to use the same capabilities directory inside the container
		if len(capabilitiesDirs) > 1 {
			for capabilityDir, count := range capabilitiesDirsFound {
				if count != len(workerNodes) {
					return nil, fmt.Errorf("the same capability container dir %s wasn't defined for all worker nodes in nodeset %s. All worker nodes in the same nodeset must have the same capability container dir", capabilityDir, input.NodeSetInputs[j].Name)
				}
			}
			destinationDir = capabilitiesDirs[0]
		}

		for capability := range capabilitiesFound {
			absSource, pathErr := filepath.Abs(capability)
			if err != nil {
				return nil, errors.Wrapf(pathErr, "failed to get absolute path to capability %s", capability)
			}

			destination := filepath.Join(destinationDir, filepath.Base(capability))
			_, copyErr := input.NixShell.RunCommand(fmt.Sprintf("devspace run copy-to-pods --no-warn --var POD_NAME_PATTERN=%s --var SOURCE=%s --var DESTINATION=%s", podNamePattern, absSource, destination))
			if copyErr != nil {
				return nil, errors.Wrap(copyErr, "failed to copy capability to pods")
			}
		}

		nsOutput, err := infra.ReadNodeSetURL(filepath.Join(".", input.CribConfigsDir), donMetadata)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read node set URLs from file")
		}

		input.NodeSetInputs[j].Out = nsOutput
	}

	return input.NodeSetInputs, nil
}

func DeployJd(input *types.DeployCribJdInput) (*jd.Output, error) {
	if input == nil {
		return nil, errors.New("DeployCribJdInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	imgTagIndex, err := dockerImageTag(input.JDInput.Image)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image tag")
	}

	jdEnvVars := map[string]string{
		"JOB_DISTRIBUTOR_IMAGE_TAG": imgTagIndex,
	}
	_, err = input.NixShell.RunCommandWithEnvVars("devspace run deploy-jd --no-warn", jdEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run devspace run deploy-jd")
	}

	jdOut, err := infra.ReadJdURL(filepath.Join(".", input.CribConfigsDir))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read JD URL from file")
	}

	return jdOut, nil
}

func nodesetDockerImage(nodeSet *types.CapabilitiesAwareNodeSet) (string, error) {
	dockerImages := []string{}
	for nodeIdx, nodeSpec := range nodeSet.NodeSpecs {
		if nodeSpec.Node.DockerContext != "" {
			return "", fmt.Errorf("docker context is not supported in CRIB. Please remove docker_ctx from the node at index %d in nodeSet %s", nodeIdx, nodeSet.Name)
		}
		if nodeSpec.Node.DockerFilePath != "" {
			return "", fmt.Errorf("dockerfile is not supported in CRIB. Please remove docker_file from the node spec at index %d in nodeSet %s", nodeIdx, nodeSet.Name)
		}

		if slices.Contains(dockerImages, nodeSpec.Node.Image) {
			continue
		}
		dockerImages = append(dockerImages, nodeSpec.Node.Image)
	}

	if len(dockerImages) != 1 {
		return "", fmt.Errorf("all nodes in each nodeSet %s must use the same Docker image, but %d different images were found: %s", nodeSet.Name, len(dockerImages), strings.Join(dockerImages, ", "))
	}

	return dockerImages[0], nil
}

func dockerImageName(image string) (string, error) {
	imgTagIndex := strings.LastIndex(image, ":")
	if imgTagIndex == -1 {
		return "", fmt.Errorf("docker image must have an explicit tag, but it was: %s", image)
	}

	return image[:imgTagIndex], nil
}

func dockerImageTag(image string) (string, error) {
	imgTagIndex := strings.LastIndex(image, ":")
	if imgTagIndex == -1 {
		return "", fmt.Errorf("docker image must have an explicit tag, but it was: %s", image)
	}

	return image[imgTagIndex+1:], nil // +1 to exclude the colon
}
