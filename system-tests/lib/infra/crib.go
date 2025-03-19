package infra

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func ReadBlockchainUrl(cribConfigsDir, chainType, chainID string) (*blockchain.Output, error) {
	fileName := filepath.Join(cribConfigsDir, fmt.Sprintf("chain-%s-urls.json", chainID))
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", fileName)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", fileName)
	}
	chainURLs := types.ChainURLs{}

	err = json.Unmarshal(byteValue, &chainURLs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal JSON from file %s", fileName)
	}

	out := &blockchain.Output{}
	out.UseCache = true
	out.ChainID = chainID
	out.Family = chainType
	out.Nodes = []*blockchain.Node{
		{
			HostWSUrl:             chainURLs.WSHostURL,
			HostHTTPUrl:           chainURLs.HTTPHostURL,
			DockerInternalWSUrl:   chainURLs.WSInternalURL,
			DockerInternalHTTPUrl: chainURLs.HTTPInternalURL,
		},
	}

	return out, nil
}

func ReadJdUrl(cribConfigsDir string) (*jd.Output, error) {
	fileName := filepath.Join(cribConfigsDir, "jd-urls.json")
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", fileName)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", fileName)
	}
	jdURLs := types.JdURLs{}

	err = json.Unmarshal(byteValue, &jdURLs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal JSON from file %s", fileName)
	}

	out := &jd.Output{}
	out.UseCache = true
	out.HostGRPCUrl = jdURLs.GRPCHostURL
	out.HostWSRPCUrl = jdURLs.WSHostURL
	out.DockerGRPCUrl = jdURLs.GRCPInternalUrl
	out.DockerWSRPCUrl = jdURLs.WSInternalURL

	return out, nil
}

func ReadNodeSetUrls(cribConfigsDir string, donMetadata *cretypes.DonMetadata) (*ns.Output, error) {
	donFileName := filepath.Join(cribConfigsDir, fmt.Sprintf("don-%s-urls.json", donMetadata.Name))

	// read DON URLs
	donFile, err := os.Open(donFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", donFileName)
	}
	defer donFile.Close()

	donByteValue, err := io.ReadAll(donFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", donFileName)
	}
	donURLs := types.DonURLs{}

	err = json.Unmarshal(donByteValue, &donURLs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal JSON from file %s", donFileName)
	}

	// read API credentials
	credsFileName := filepath.Join(".", "crib-configs", "don-api-credentials.json")
	credsFile, err := os.Open(credsFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", credsFileName)
	}
	defer credsFile.Close()

	credsByteValue, err := io.ReadAll(credsFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", credsFileName)
	}
	apiCredentials := types.DonAPICredentials{}

	err = json.Unmarshal(credsByteValue, &apiCredentials)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal JSON from file %s", credsFileName)
	}

	bootstrapNodes, err := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &cretypes.Label{Key: libnode.NodeTypeKey, Value: cretypes.BootstrapNode}, libnode.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap nodes")
	}

	workerNodes, err := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &cretypes.Label{Key: libnode.NodeTypeKey, Value: cretypes.WorkerNode}, libnode.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	if len(bootstrapNodes) != len(donURLs.BootstrapNodes) {
		return nil, errors.Errorf("number of bootstrap nodes in JSON file must match the number of bootstrap nodes in the DON, but got %d bootstrap nodes in JSON file and %d bootstrap nodes in the DON", len(donURLs.BootstrapNodes), len(bootstrapNodes))
	}

	if len(workerNodes) != len(donURLs.WorkerNodes) {
		return nil, errors.Errorf("number of worker nodes in JSON file must match the number of worker nodes in the DON, but got %d worker nodes in JSON file and %d worker nodes in the DON", len(donURLs.WorkerNodes), len(workerNodes))
	}

	out := &ns.Output{}
	out.UseCache = true
	out.CLNodes = []*clnode.Output{}

	for i := range bootstrapNodes {
		out.CLNodes = append(out.CLNodes, &clnode.Output{
			UseCache: true,
			Node: &clnode.NodeOut{
				APIAuthUser:     apiCredentials.Username,
				APIAuthPassword: apiCredentials.Password,
				HostURL:         donURLs.BootstrapNodes[i].HostURL,
				DockerURL:       donURLs.BootstrapNodes[i].InternalURL,
				DockerP2PUrl:    donURLs.BootstrapNodes[i].P2PInternalURL,
				InternalIP:      donURLs.BootstrapNodes[i].InternalIP,
			},
		})
	}

	for i := range workerNodes {
		out.CLNodes = append(out.CLNodes, &clnode.Output{
			UseCache: true,
			Node: &clnode.NodeOut{
				APIAuthUser:     apiCredentials.Username,
				APIAuthPassword: apiCredentials.Password,
				HostURL:         donURLs.WorkerNodes[i].HostURL,
				DockerURL:       donURLs.WorkerNodes[i].InternalURL,
				DockerP2PUrl:    donURLs.WorkerNodes[i].P2PInternalURL,
				InternalIP:      donURLs.WorkerNodes[i].InternalIP,
			},
		})
	}

	return out, nil
}
