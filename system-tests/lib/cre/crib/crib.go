package crib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/crib-sdk/crib"
	anvilv1 "github.com/smartcontractkit/crib-sdk/crib/composite/blockchain/anvil/v1"
	jdv1 "github.com/smartcontractkit/crib-sdk/crib/composite/chainlink/jd/v1"
	nodev1 "github.com/smartcontractkit/crib-sdk/crib/composite/chainlink/node/v1"
	telepresencev1 "github.com/smartcontractkit/crib-sdk/crib/composite/cluster-services/telepresence/v1"
	namespacev1 "github.com/smartcontractkit/crib-sdk/crib/scalar/k8s/namespace/v1"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
)

func StartNixShell(input *cre.StartNixShellInput) (*nix.Shell, error) {
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

	if strings.EqualFold(input.InfraInput.CRIB.Provider, infra.AWS) {
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

func Bootstrap(infraInput *infra.Input) error {
	plan := crib.NewPlan(
		"namespace",
		crib.Namespace(infraInput.CRIB.Namespace),
		crib.ComponentSet(
			namespacev1.Component(infraInput.CRIB.Namespace),
			telepresencev1.Component(&telepresencev1.Props{
				Namespace:         infraInput.CRIB.Namespace,
				QuitBeforeRunning: true,
			}),
		),
	)
	_, err := plan.Apply(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to apply plan")
	}

	return nil
}

func DeployBlockchain(input *cre.DeployCribBlockchainInput) (*blockchain.Output, error) {
	err := input.Validate()
	if err != nil {
		return nil, errors.Wrapf(err, "invalid input for deploying blockchain")
	}

	ctx := context.Background()

	anvil := anvilv1.Component(&anvilv1.Props{
		Namespace: input.Namespace,
		ChainID:   input.BlockchainInput.ChainID,
	})

	plan := crib.NewPlan(
		"anvilv1",
		crib.Namespace(input.Namespace),
		crib.ComponentSet(
			anvil,
		),
	)

	result, err := plan.Apply(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply a plan")
	}

	anvilComponents := result.ComponentByName(anvilv1.ComponentName)

	for component := range anvilComponents {
		res := crib.ComponentState[anvilv1.Result](component)

		return &blockchain.Output{
			Type:    input.BlockchainInput.Type,
			Family:  "evm",
			ChainID: input.BlockchainInput.ChainID,
			Nodes: []*blockchain.Node{
				{
					InternalWSUrl:   res.RPCWebsocketURL(),
					ExternalWSUrl:   res.RPCWebsocketURL(),
					InternalHTTPUrl: res.RPCHTTPURL(),
					ExternalHTTPUrl: res.RPCHTTPURL(),
				},
			},
		}, nil
	}

	return nil, errors.New("failed to find a valid component")
}
func DeployDons(input *cre.DeployCribDonsInput) ([]*cre.CapabilitiesAwareNodeSet, error) {
	if input == nil {
		return nil, errors.New("DeployCribDonsInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	componentFuncs := make([]crib.ComponentFunc, 0)

	for j, donMetadata := range input.Topology.DonsMetadata {
		imageName, imageTag, err := imageNameAndTag(input, j)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get image name and tag for %s", donMetadata.Name)
		}

		for i, nodeMetadata := range donMetadata.NodesMetadata {
			configToml, secrets, confSecretsErr := getConfigAndSecretsForNode(nodeMetadata, j, input, donMetadata)
			if confSecretsErr != nil {
				return nil, confSecretsErr
			}
			nodeSpec, confSecretsErr := getNodeSpecForNode(nodeMetadata, j, input, donMetadata)
			if confSecretsErr != nil {
				return nil, errors.Wrapf(confSecretsErr, "failed to get node spec for %s", donMetadata.Name)
			}
			cFunc := nodev1.Component(&nodev1.Props{
				Namespace:       input.Namespace,
				Image:           fmt.Sprintf("%s:%s", imageName, imageTag),
				AppInstanceName: fmt.Sprintf("%s-%d", donMetadata.Name, i),
				// passing as config not as override
				Config: *configToml,
				SecretsOverrides: map[string]string{
					"overrides": *secrets,
				},
				EnvVars: nodeSpec.Node.EnvVars,
			})
			componentFuncs = append(componentFuncs, cFunc)
		}
	}

	plan := crib.NewPlan(
		"nodesets",
		crib.Namespace(input.Namespace),
		crib.ComponentSet(
			componentFuncs...,
		),
	)

	planState, err := plan.Apply(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply plan")
	}

	// Setting outputs based on the Plan Results
	nodeComponents := planState.ComponentByName(nodev1.ComponentName)

	nodeResults := make([]nodev1.Result, 0)

	for component := range nodeComponents {
		res := crib.ComponentState[nodev1.Result](component)
		nodeResults = append(nodeResults, res)
		fmt.Printf("Node API URL: %s\n", res.APIUrl())
		fmt.Printf("API Credentials: username: %s , password: %s\n", res.APICredentials.UserName, res.APICredentials.Password)
	}

	// setting outputs in a similar way as in func ReadNodeSetURL
	for j := range input.Topology.DonsMetadata {
		out := &ns.Output{
			// UseCache: true will disable deploying docker containers via CTF
			UseCache: true,
			CLNodes:  []*clnode.Output{},
		}
		// todo: for now this is hardcoded for a single don, we need to group results for each don
		for _, res := range nodeResults {
			out.CLNodes = append(out.CLNodes, &clnode.Output{
				// UseCache: true will disable deploying docker containers via CTF
				UseCache: true,
				Node: &clnode.NodeOut{
					APIAuthUser:     res.APICredentials.UserName,
					APIAuthPassword: res.APICredentials.Password,
					ExternalURL:     res.APIUrl(),
					InternalURL:     res.APIUrl(),
					// todo: this should be simplified in the CTF types, we should just pass P2P port
					InternalP2PUrl: fmt.Sprintf("http://%s:%d", res.HostName(), res.P2PPort),
					InternalIP:     res.HostName(),
				},
			})
		}
		input.NodeSetInputs[j].Out = out
	}

	return input.NodeSetInputs, nil
}

func getNodeSpecForNode(nodeMetadata *cre.NodeMetadata, donIndex int, input *cre.DeployCribDonsInput, donMetadata *cre.DonMetadata) (*clnode.Input, error) {
	nodeIndexStr, findErr := libnode.FindLabelValue(nodeMetadata, libnode.IndexKey)
	if findErr != nil {
		return nil, errors.Wrapf(findErr, "failed to find node index in nodeset %s", donMetadata.Name)
	}

	nodeIndex, convErr := strconv.Atoi(nodeIndexStr)
	if convErr != nil {
		return nil, errors.Wrapf(convErr, "failed to convert node index '%s' to int in nodeset %s", nodeIndexStr, donMetadata.Name)
	}

	nodeSpec := input.NodeSetInputs[donIndex].NodeSpecs[nodeIndex]
	return nodeSpec, nil
}

func getConfigAndSecretsForNode(nodeMetadata *cre.NodeMetadata, donIndex int, input *cre.DeployCribDonsInput, donMetadata *cre.DonMetadata) (*string, *string, error) {
	nodeSpec, err := getNodeSpecForNode(nodeMetadata, donIndex, input, donMetadata)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get node spec")
	}
	cleanedToml, tomlErr := cleanToml(nodeSpec.Node.TestConfigOverrides)
	if tomlErr != nil {
		return nil, nil, errors.Wrap(tomlErr, "failed to clean TOML")
	}

	// Merge user overrides
	cleanedUserToml, tomlErr := cleanToml(nodeSpec.Node.UserConfigOverrides)
	if tomlErr != nil {
		return nil, nil, errors.Wrap(tomlErr, "failed to clean user TOML")
	}

	finalToml, err := mergeToml(cleanedToml, cleanedUserToml)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to merge TOML")
	}

	secretsFileBytes := []byte(nodeSpec.Node.TestSecretsOverrides)

	tomlString := string(finalToml)
	secretsString := string(secretsFileBytes)
	return &tomlString, &secretsString, nil
}

//nolint:unused // for now we don't need to set capabilities (high complexity, low impact) we'll rely on plugins image which contains all required capabilities
func setCapabilities(input *cre.DeployCribDonsInput, donIndex int, workerNodes []*cre.NodeMetadata) error {
	// validate capabilities-related configuration and copy capabilities to pods
	podNamePattern := input.NodeSetInputs[donIndex].Name + `-\\d+`
	_, regErr := regexp.Compile(podNamePattern)
	if regErr != nil {
		return errors.Wrapf(regErr, "failed to compile regex for pod name pattern %s", podNamePattern)
	}

	capabilitiesFound := map[string]int{}
	capabilitiesDirs := []string{}
	capabilitiesDirsFound := map[string]int{}

	// make sure all worker nodes in DON have the same set of capabilities
	// in the future we might want to allow different capabilities for different nodes
	// but for now we require all worker nodes in the same DON to have the same capabilities
	for _, nodeSpec := range input.NodeSetInputs[donIndex].NodeSpecs {
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
			return fmt.Errorf("capability %s wasn't defined for all worker nodes in nodeset %s. All worker nodes in the same nodeset must have the same capabilities", capability, input.NodeSetInputs[donIndex].Name)
		}
	}

	destinationDir, err := crecapabilities.DefaultContainerDirectory(infra.CRIB)
	if err != nil {
		return errors.Wrap(err, "failed to get default directory for capabilities in CRIB")
	}

	// all of them need to use the same capabilities directory inside the container
	if len(capabilitiesDirs) > 1 {
		for capabilityDir, count := range capabilitiesDirsFound {
			if count != len(workerNodes) {
				return fmt.Errorf("the same capability container dir %s wasn't defined for all worker nodes in nodeset %s. All worker nodes in the same nodeset must have the same capability container dir", capabilityDir, input.NodeSetInputs[donIndex].Name)
			}
		}
		destinationDir = capabilitiesDirs[0]
	}

	for capability := range capabilitiesFound {
		absSource, pathErr := filepath.Abs(capability)
		if pathErr != nil {
			return errors.Wrapf(pathErr, "failed to get absolute path to capability %s", capability)
		}
		// ensure +x chmod in capability binary before copying to pods
		err := os.Chmod(capability, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to chmod capability %s", capability)
		}
		destination := filepath.Join(destinationDir, filepath.Base(capability))
		_, copyErr := input.NixShell.RunCommand(fmt.Sprintf("devspace run copy-to-pods --no-warn --var POD_NAME_PATTERN=%s --var SOURCE=%s --var DESTINATION=%s", podNamePattern, absSource, destination))
		if copyErr != nil {
			return errors.Wrap(copyErr, "failed to copy capability to pods")
		}
	}
	return nil
}

func imageNameAndTag(input *cre.DeployCribDonsInput, j int) (string, string, error) {
	// validate that all nodes in the same node set use the same Docker image
	dockerImage, dockerImagesErr := nodesetDockerImage(input.NodeSetInputs[j])
	if dockerImagesErr != nil {
		return "", "", errors.Wrap(dockerImagesErr, "failed to validate node set Docker images")
	}

	imageName, imageErr := dockerImageName(dockerImage)
	if imageErr != nil {
		return "", "", errors.Wrap(imageErr, "failed to get image name")
	}

	imageTag, imageErr := dockerImageTag(dockerImage)
	if imageErr != nil {
		return "", "", errors.Wrap(imageErr, "failed to get image tag")
	}
	return imageName, imageTag, nil
}

func cleanToml(tomlStr string) ([]byte, error) {
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

// mergeToml merges two TOML configurations.
// It takes base TOML content (tomlOne) and overlay TOML content (tomlTwo) as byte slices,
// and combines them with the overlay values taking precedence over the base values.
func mergeToml(tomlOne []byte, tomlTwo []byte) ([]byte, error) {
	// Parse the first TOML
	var baseConfig map[string]interface{}
	if err := toml.Unmarshal(tomlOne, &baseConfig); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal first TOML")
	}

	// Parse the second TOML
	var overlayConfig map[string]interface{}
	if err := toml.Unmarshal(tomlTwo, &overlayConfig); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal second TOML")
	}

	// Deep merge the maps
	for k, v := range overlayConfig {
		// If both values are maps, merge them recursively
		if baseVal, ok := baseConfig[k]; ok {
			if baseMap, isBaseMap := baseVal.(map[string]interface{}); isBaseMap {
				if overlayMap, isOverlayMap := v.(map[string]interface{}); isOverlayMap {
					// Recursively merge nested maps
					for nestedKey, nestedVal := range overlayMap {
						baseMap[nestedKey] = nestedVal
					}
					continue
				}
			}
		}
		// Otherwise, override the value
		baseConfig[k] = v
	}

	// Marshal back to TOML
	result, err := toml.Marshal(baseConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal merged config")
	}

	return result, nil
}

func DeployJd(input *cre.DeployCribJdInput) (*jd.Output, error) {
	if input == nil {
		return nil, errors.New("DeployCribJdInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	jdComponent := jdv1.Component(&jdv1.Props{
		Namespace: input.Namespace,
		JD: jdv1.JDProps{
			Image:            input.JDInput.Image,
			CSAEncryptionKey: input.JDInput.CSAEncryptionKey,
		},
		WaitForRollout: true,
	})

	plan := crib.NewPlan(
		"jd",
		crib.Namespace(input.Namespace),
		crib.ComponentSet(
			jdComponent,
		),
	)

	planState, err := plan.Apply(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply a plan")
	}

	for component := range planState.ComponentByName(jdv1.ComponentName) {
		jdResult := crib.ComponentState[jdv1.Result](component)

		out := &jd.Output{}
		out.UseCache = true
		out.ExternalGRPCUrl = jdResult.GRPCHostURL()
		out.InternalGRPCUrl = jdResult.GRPCHostURL()
		out.InternalWSRPCUrl = jdResult.WSRPCHostURL()

		return out, nil
	}

	return nil, errors.New("failed to find a valid jd component in results")
}

func nodesetDockerImage(nodeSet *cre.CapabilitiesAwareNodeSet) (string, error) {
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
