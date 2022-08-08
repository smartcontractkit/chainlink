// Package actions enables common chainlink interactions
package actions

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// GinkgoSuite provides the default setup for running a Ginkgo test suite
func GinkgoSuite() {
	logging.Init()
	gomega.RegisterFailHandler(ginkgo.Fail)
}

// ContractDeploymentInterval After how many contract actions to wait before starting any more
// Example: When deploying 1000 contracts, stop every ContractDeploymentInterval have been deployed to wait before continuing
var ContractDeploymentInterval = 200

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodes(
	nodes []*client.Chainlink,
	client blockchain.EVMClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = client.Fund(toAddress, amount)
		if err != nil {
			return err
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodesLink(
	nodes []*client.Chainlink,
	blockchain blockchain.EVMClient,
	linkToken contracts.LinkToken,
	linkAmount *big.Int,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = linkToken.Transfer(toAddress, linkAmount)
		if err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddresses(nodes []*client.Chainlink) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// SetChainlinkAPIPageSize specifies the page size from the Chainlink API, useful for high volume testing
func SetChainlinkAPIPageSize(nodes []*client.Chainlink, pageSize int) {
	for _, n := range nodes {
		n.SetPageSize(pageSize)
	}
}

// ExtractRequestIDFromJobRun extracts RequestID from job runs response
func ExtractRequestIDFromJobRun(jobDecodeData client.RunsResponseData) ([]byte, error) {
	var taskRun client.TaskRun
	for _, tr := range jobDecodeData.Attributes.TaskRuns {
		if tr.Type == "ethabidecodelog" {
			taskRun = tr
		}
	}
	var decodeLogTaskRun *client.DecodeLogTaskRun
	if err := json.Unmarshal([]byte(taskRun.Output), &decodeLogTaskRun); err != nil {
		return nil, err
	}
	rqInts := decodeLogTaskRun.RequestID
	return rqInts, nil
}

// EncodeOnChainVRFProvingKey encodes uncompressed public VRF key to on-chain representation
func EncodeOnChainVRFProvingKey(vrfKey client.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Data.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var set1 bool
	var set2 bool
	// strip 0x to convert to int
	provingKey[0], set1 = new(big.Int).SetString(uncompressed[2:66], 16)
	if !set1 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(uncompressed[66:], 16)
	if !set2 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	return provingKey, nil
}

// GetMockserverInitializerDataForOTPE creates mocked weiwatchers data needed for otpe
func GetMockserverInitializerDataForOTPE(
	OCRInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
) (interface{}, error) {
	var contractsInfo []ctfClient.ContractInfoJSON

	for index, OCRInstance := range OCRInstances {
		contractInfo := ctfClient.ContractInfoJSON{
			ContractVersion: 4,
			Path:            fmt.Sprintf("contract_%d", index),
			Status:          "live",
			ContractAddress: OCRInstance.Address(),
		}

		contractsInfo = append(contractsInfo, contractInfo)
	}

	contractsInitializer := ctfClient.HttpInitializer{
		Request:  ctfClient.HttpRequest{Path: "/contracts.json"},
		Response: ctfClient.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []ctfClient.NodeInfoJSON

	for _, chainlink := range chainlinkNodes {
		ocrKeys, err := chainlink.MustReadOCRKeys()
		if err != nil {
			return nil, err
		}
		nodeInfo := ctfClient.NodeInfoJSON{
			NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
			ID:          ocrKeys.Data[0].ID,
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}

	nodesInitializer := ctfClient.HttpInitializer{
		Request:  ctfClient.HttpRequest{Path: "/nodes.json"},
		Response: ctfClient.HttpResponse{Body: nodesInfo},
	}
	initializers := []ctfClient.HttpInitializer{contractsInitializer, nodesInitializer}
	return initializers, nil
}

// TeardownSuite tears down networks/clients and environment and creates a logs folder for failed tests in the
// specified path. Can also accept a testreporter (if one was used) to log further results
func TeardownSuite(
	env *environment.Environment,
	logsFolderPath string,
	chainlinkNodes []*client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	c blockchain.EVMClient,
) error {
	if err := testreporters.WriteTeardownLogs(env, optionalTestReporter); err != nil {
		return errors.Wrap(err, "Error dumping environment logs, leaving environment running for manual retrieval")
	}
	if c != nil && chainlinkNodes != nil && len(chainlinkNodes) > 0 && !c.NetworkSimulated() {
		if err := LogChainlinkKeys(chainlinkNodes); err != nil {
			log.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
				Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
					"Environment is left running so you can try manually!")
		}
	} else {
		log.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
	}
	// nolint
	if c != nil {
		c.Close()
	}

	keepEnvs := os.Getenv("KEEP_ENVIRONMENTS")
	if keepEnvs == "" {
		keepEnvs = "NEVER"
	}

	switch strings.ToUpper(keepEnvs) {
	case "ALWAYS":
	case "ONFAIL":
		if ginkgo.CurrentSpecReport().Failed() {
			return env.Shutdown()
		}
	case "NEVER":
		return env.Shutdown()
	default:
		log.Warn().Str("Invalid Keep Value", keepEnvs).
			Msg("Invalid 'keep_environments' value, see the 'framework.yaml' file")
	}
	return nil
}

// TeardownRemoteSuite is used when running a test within a remote-test-runner, like for long-running performance and
// soak tests
func TeardownRemoteSuite(
	env *environment.Environment,
	chainlinkNodes []*client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	client blockchain.EVMClient,
) error {
	var err error
	if err = testreporters.SendReport(env, "./", optionalTestReporter); err != nil {
		log.Warn().Err(err).Msg("Error writing test report")
	}
	if err = LogChainlinkKeys(chainlinkNodes); err != nil {
		log.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}
	return err
}

// logChainlinkKeys retrieves and decrypts funded keys on the Chainlink nodes, and logs them.
// This is used for tests on real networks, and WILL LOG PRIVATE KEY INFO OF THE NODES. Use only for tests where the
// keys aren't used for anything else, and the nodes are ephemeral. This will also use a significant amount of RAM.
// TODO: Modify method to directly transfer funds instead of logging keys.
func LogChainlinkKeys(chainlinkNodes []*client.Chainlink) error {
	var (
		keysMutex     sync.Mutex
		keysToDecrypt = [][]byte{}
	)

	fundsErrGroup := new(errgroup.Group)
	for _, n := range chainlinkNodes {
		node := n
		fundsErrGroup.Go(func() error {
			keys, err := node.ExportEVMKeys()
			if err != nil {
				return err
			}
			for _, key := range keys {
				log.Debug().Str("Password", client.ChainlinkKeyPassword).Interface("Key", key).Msg("Decrypting Key")
				keyJson, err := json.Marshal(key)
				if err != nil {
					return err
				}
				keysMutex.Lock()
				keysToDecrypt = append(keysToDecrypt, keyJson)
				keysMutex.Unlock()
			}
			return nil
		})
	}
	if err := fundsErrGroup.Wait(); err != nil {
		return err
	}

	for _, key := range keysToDecrypt {
		log.Debug().Msg("Decrypting Key. This can take some time (and a good bit of RAM)")
		decrypted, err := keystore.DecryptKey(key, client.ChainlinkKeyPassword)
		if err != nil {
			return err
		}
		log.Info().Str("Key", fmt.Sprintf("%x", crypto.FromECDSA(decrypted.PrivateKey))).Msg("Decrypted Chainlink Node Key")
	}
	return nil
}

// FundAddresses will fund a list of addresses with an amount of native currency
func FundAddresses(blockchain blockchain.EVMClient, amount *big.Float, addresses ...string) error {
	for _, address := range addresses {
		if err := blockchain.Fund(address, amount); err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// EncodeOnChainExternalJobID encodes external job uuid to on-chain representation
func EncodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}
