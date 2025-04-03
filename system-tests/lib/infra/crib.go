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

func ReadBlockchainURL(cribConfigsDir, chainType, chainID string) (*blockchain.Output, error) {
	fileName := filepath.Join(cribConfigsDir, fmt.Sprintf("chain-%s-urls.json", chainID))
	chainURLs := types.ChainURLs{}
	err := readAndUnmarshalJSON(fileName, &chainURLs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read and unmarshal chain URLs JSON")
	}

	out := &blockchain.Output{}
	out.UseCache = true
	out.ChainID = chainID
	out.Family = chainType
	out.Nodes = []*blockchain.Node{
		{
			ExternalWSUrl:   chainURLs.WSExternalURL,
			ExternalHTTPUrl: chainURLs.HTTPExternalURL,
			InternalWSUrl:   chainURLs.WSInternalURL,
			InternalHTTPUrl: chainURLs.HTTPInternalURL,
		},
	}

	return out, nil
}

func ReadJdURL(cribConfigsDir string) (*jd.Output, error) {
	fileName := filepath.Join(cribConfigsDir, "jd-urls.json")

	jdURLs := types.JdURLs{}
	err := readAndUnmarshalJSON(fileName, &jdURLs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read and unmarshal JD URLs JSON")
	}

	out := &jd.Output{}
	out.UseCache = true
	out.ExternalGRPCUrl = jdURLs.GRPCExternalURL
	out.ExternalWSRPCUrl = jdURLs.WSExternalURL
	out.InternalGRPCUrl = jdURLs.GRCPInternalURL
	out.InternalWSRPCUrl = jdURLs.WSInternalURL

	return out, nil
}

func ReadNodeSetURL(cribConfigsDir string, donMetadata *cretypes.DonMetadata) (*ns.Output, error) {
	// read DON URLs
	donFileName := filepath.Join(cribConfigsDir, fmt.Sprintf("don-%s-urls.json", donMetadata.Name))
	donURLs := types.DonURLs{}
	err := readAndUnmarshalJSON(donFileName, &donURLs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read and unmarshal don URLs JSON")
	}

	// read API credentials
	credsFileName := filepath.Join(".", "crib-configs", "don-api-credentials.json")

	apiCredentials := types.DonAPICredentials{}
	err = readAndUnmarshalJSON(credsFileName, &apiCredentials)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read and unmarshal don API credentials JSON")
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
				ExternalURL:     donURLs.BootstrapNodes[i].ExternalURL,
				InternalURL:     donURLs.BootstrapNodes[i].InternalURL,
				InternalP2PUrl:  donURLs.BootstrapNodes[i].P2PInternalURL,
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
				ExternalURL:     donURLs.WorkerNodes[i].ExternalURL,
				InternalURL:     donURLs.WorkerNodes[i].InternalURL,
				InternalP2PUrl:  donURLs.WorkerNodes[i].P2PInternalURL,
				InternalIP:      donURLs.WorkerNodes[i].InternalIP,
			},
		})
	}

	return out, nil
}

func readAndUnmarshalJSON[Type any](fileName string, target *Type) error {
	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", fileName)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", fileName)
	}

	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal JSON from file %s", fileName)
	}

	return nil
}
