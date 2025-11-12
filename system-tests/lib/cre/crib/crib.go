package crib

import (
	"context"
	"fmt"
	"maps"
	"slices"
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
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func Bootstrap(ctx context.Context, infraInput infra.Provider) error {
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
	_, err := plan.Apply(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to apply plan")
	}

	return nil
}

type DeployCribBlockchainInput struct {
	Blockchain     *blockchain.Input
	CribConfigsDir string
	Namespace      string
}

func (d *DeployCribBlockchainInput) Validate() error {
	if d.Blockchain == nil {
		return errors.New("blockchain input not set")
	}
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	if d.Namespace == "" {
		return errors.New("namespace not set")
	}
	return nil
}

func DeployBlockchain(ctx context.Context, input *DeployCribBlockchainInput) (*blockchain.Output, error) {
	err := input.Validate()
	if err != nil {
		return nil, errors.Wrapf(err, "invalid input for deploying blockchain")
	}

	anvil := anvilv1.Component(&anvilv1.Props{
		Namespace: input.Namespace,
		ChainID:   input.Blockchain.ChainID,
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
			Type:    input.Blockchain.Type,
			Family:  "evm",
			ChainID: input.Blockchain.ChainID,
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

type DeployCribDonsInput struct {
	Topology       *cre.Topology
	NodeSet        []*cre.NodeSet
	CribConfigsDir string
	Namespace      string
}

func (d *DeployCribDonsInput) Validate() error {
	if d.Topology == nil {
		return errors.New("topology not set")
	}
	if len(d.Topology.DonsMetadata.List()) == 0 {
		return errors.New("metadata not set")
	}
	if len(d.NodeSet) == 0 {
		return errors.New("node set inputs not set")
	}
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}

func DeployDons(ctx context.Context, input *DeployCribDonsInput) ([]*cre.NodeSet, error) {
	if input == nil {
		return nil, errors.New("DeployCribDonsInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	componentFuncs := make([]crib.ComponentFunc, 0)

	for donIdx, donMetadata := range input.Topology.DonsMetadata.List() {
		imageName, imageTag, err := imageNameAndTag(input, donIdx)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get image name and tag for %s", donMetadata.Name)
		}

		for nodeIdx, nodeMetadata := range donMetadata.NodesMetadata {
			configToml, secrets, confSecretsErr := getConfigAndSecretsForNode(nodeMetadata, donIdx, input)
			if confSecretsErr != nil {
				return nil, confSecretsErr
			}

			cFunc := nodev1.Component(&nodev1.Props{
				Namespace:       input.Namespace,
				Image:           fmt.Sprintf("%s:%s", imageName, imageTag),
				AppInstanceName: fmt.Sprintf("%s-%d", donMetadata.Name, nodeIdx),
				// passing as config not as override
				Config: *configToml,
				SecretsOverrides: map[string]string{
					"overrides": *secrets,
				},
				EnvVars: input.NodeSet[donIdx].NodeSpecs[nodeMetadata.Index].Node.EnvVars,
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

	planState, err := plan.Apply(ctx)
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
	for j := range input.Topology.DonsMetadata.List() {
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
		input.NodeSet[j].Out = out
	}

	return input.NodeSet, nil
}

func getConfigAndSecretsForNode(nodeMetadata *cre.NodeMetadata, donIndex int, input *DeployCribDonsInput) (*string, *string, error) {
	nodeSpec := input.NodeSet[donIndex].NodeSpecs[nodeMetadata.Index]

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

func imageNameAndTag(input *DeployCribDonsInput, j int) (string, string, error) {
	// validate that all nodes in the same node set use the same Docker image
	dockerImage, dockerImagesErr := nodesetDockerImage(input.NodeSet[j])
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
	var data any
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
	var baseConfig map[string]any
	if err := toml.Unmarshal(tomlOne, &baseConfig); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal first TOML")
	}

	// Parse the second TOML
	var overlayConfig map[string]any
	if err := toml.Unmarshal(tomlTwo, &overlayConfig); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal second TOML")
	}

	// Deep merge the maps
	for k, v := range overlayConfig {
		// If both values are maps, merge them recursively
		if baseVal, ok := baseConfig[k]; ok {
			if baseMap, isBaseMap := baseVal.(map[string]any); isBaseMap {
				if overlayMap, isOverlayMap := v.(map[string]any); isOverlayMap {
					// Recursively merge nested maps
					maps.Copy(baseMap, overlayMap)
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

type DeployCribJdInput struct {
	JDInput        jd.Input
	CribConfigsDir string
	Namespace      string
}

func (d *DeployCribJdInput) Validate() error {
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}

func DeployJd(ctx context.Context, input *DeployCribJdInput) (*jd.Output, error) {
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

	planState, err := plan.Apply(ctx)
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

func nodesetDockerImage(nodeSet *cre.NodeSet) (string, error) {
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
