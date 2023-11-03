package zksync

import (
	"context"
	"fmt"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	helmEth "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet"
)

type ZKSyncClient struct {
	GRunner              *gauntlet.GauntletRunner
	LinkAddr             string
	OCRAddr              string
	ContractLoader       contracts.ContractLoader
	LinkContract         contracts.LinkToken
	OCRContract          contracts.OffchainAggregator
	AccessControllerAddr string
	L2RPC                string
	NKeys                []client.NodeKeysBundle
	Transmitters         []string
	Payees               []string
	Signers              []string
	PeerIds              []string
	OcrConfigPubKeys     []string
	Client               blockchain.EVMClient
	OcrInstance          []contracts.OffchainAggregator
	ChainClient          blockchain.EVMClient
	ChainlinkNodes       []*client.ChainlinkK8sClient
	Mockserver           *ctfClient.MockserverClient
}

func Setup(L2RPC string, privateKey string, client blockchain.EVMClient) (*ZKSyncClient, error) {
	g, err := gauntlet.Setup(L2RPC, privateKey)
	if err != nil {
		return nil, err
	}

	return &ZKSyncClient{
		GRunner:  g,
		LinkAddr: "",
		OCRAddr:  "",
		L2RPC:    L2RPC,
		Client:   client,
	}, nil

}

func (z *ZKSyncClient) DeployLinkToken() error {
	output, err := z.GRunner.ExecuteCommand([]string{"token:deploy"})
	if err != nil {
		return err
	}
	z.LinkAddr = output.Responses[0].Tx.Address
	return nil
}

func (z *ZKSyncClient) DeployAccessController() error {
	output, err := z.GRunner.ExecuteCommand([]string{"access_controller:deploy"})
	if err != nil {
		return err
	}
	z.AccessControllerAddr = output.Responses[0].Tx.Address
	return nil
}

func (z *ZKSyncClient) DeployOCR(ocrContractValues string) error {
	output, err := z.GRunner.ExecuteCommand([]string{"ocr:deploy", fmt.Sprintf("--input=%s", ocrContractValues)})
	if err != nil {
		return err
	}
	z.OCRAddr = output.Responses[0].Tx.Address
	return nil
}

func (z *ZKSyncClient) AddAccess(ocrAddress string) error {
	_, err := z.GRunner.ExecuteCommand([]string{"access_controller:add_access", fmt.Sprintf("--address=%s", ocrAddress), z.AccessControllerAddr})
	if err != nil {
		return err
	}
	return nil
}

func (z *ZKSyncClient) SetPayees(ocrAddress string, payees []string, transmitters []string) error {
	_, err := z.GRunner.ExecuteCommand([]string{"ocr:set_payees", ocrAddress, fmt.Sprintf("--transmitters=%s", strings.Join(transmitters, ",")), fmt.Sprintf("--payees=%s", strings.Join(payees, ","))})
	if err != nil {
		return err
	}
	return nil
}
func (z *ZKSyncClient) CreateKeys(chainlinkNodes []*client.ChainlinkClient) error {
	var err error

	z.NKeys, _, err = client.CreateNodeKeysBundle(chainlinkNodes, "evm", "280")
	if err != nil {
		return err
	}
	for _, key := range z.NKeys {
		z.PeerIds = append(z.PeerIds, key.PeerID)
		z.OcrConfigPubKeys = append(z.OcrConfigPubKeys, strings.Replace(key.OCRKeys.Data[0].Attributes.OffChainPublicKey, "ocroff_", "", 1))
		z.Transmitters = append(z.Transmitters, strings.Replace(key.EthAddress, "0x", "", 1))
		z.Signers = append(z.Signers, strings.Replace(key.OCRKeys.Data[0].Attributes.OnChainSigningAddress, "ocrsad_", "", 1))
		z.Payees = append(z.Payees, strings.Replace(z.Client.GetDefaultWallet().Address(), "0x", "", 1))
	}

	return nil
}

func (z *ZKSyncClient) SetConfig(ocrAddress, ocrConfigValues string) error {
	_, err := z.GRunner.ExecuteCommand([]string{"ocr:set_config", ocrAddress, fmt.Sprintf("--input=%s", ocrConfigValues)})
	if err != nil {
		return err
	}
	return nil
}

func (z *ZKSyncClient) FundNodes(chainlinkClient blockchain.EVMClient) error {
	for _, key := range z.NKeys {
		toAddress := common.HexToAddress(key.TXKey.Data.ID)
		log.Info().Stringer("toAddress", toAddress).Msg("Funding node")
		amount := big.NewInt(1e17)
		callMsg := ethereum.CallMsg{
			From:  common.HexToAddress(chainlinkClient.GetDefaultWallet().Address()),
			To:    &toAddress,
			Value: amount,
		}
		log.Debug().Interface("CallMsg", callMsg).Msg("Estimating gas")
		gasEstimates, err := chainlinkClient.EstimateGas(callMsg)
		if err != nil {
			return fmt.Errorf("estimating gas: %w", err)
		}
		log.Debug().Stringer("toAddress", toAddress).Stringer("amount", amount).Interface("gasEstimates", gasEstimates).Msg("Transferring funds")
		err = chainlinkClient.Fund(toAddress.String(), big.NewFloat(0).SetInt(amount), gasEstimates)
		if err != nil {
			return fmt.Errorf("funding %q: %w", toAddress, err)
		}

		log.Info().Stringer("toAddress", toAddress).Stringer("amount", amount).Msg("Transferred funds")

		// TO-DO Link funding seems to hang but tx is present on chain
		// err = z.LinkContract.Transfer(key.TXKey.Data.ID, big.NewInt(100000000))
		// if err != nil {
		//	return err
		// }
	}
	return nil
}

func (z *ZKSyncClient) DeployContracts(chainlinkClient blockchain.EVMClient, ocrContractValues *gauntlet.OCRContract, ocrConfigValues *gauntlet.OCRConfig, l zerolog.Logger) error {
	err := z.DeployLinkToken()
	if err != nil {
		return err
	}
	z.ContractLoader, err = contracts.NewContractLoader(chainlinkClient, l)
	if err != nil {
		return err
	}
	z.LinkContract, err = z.ContractLoader.LoadLINKToken(common.HexToAddress(z.LinkAddr).String())

	err = z.FundNodes(chainlinkClient)
	if err != nil {
		return err
	}

	err = z.DeployAccessController()
	if err != nil {
		return err
	}

	ocrContractValues.Link = z.LinkAddr
	ocrContractValues.BillingAccessController = z.AccessControllerAddr
	ocrContractValues.RequesterAccessController = z.AccessControllerAddr

	ocrConfigValues.Signers = z.Signers
	ocrConfigValues.Transmitters = z.Transmitters
	ocrConfigValues.OcrConfigPublicKeys = z.OcrConfigPubKeys
	ocrConfigValues.OperatorsPeerIds = strings.Join(z.PeerIds, ",")

	ocrJsonContract, err := ocrContractValues.MarshalOCR()
	if err != nil {
		return err
	}
	err = z.DeployOCR(ocrJsonContract)
	if err != nil {
		return err
	}
	z.OCRContract, err = z.ContractLoader.LoadOcrContract(common.HexToAddress(z.OCRAddr))
	if err != nil {
		return err
	}

	err = z.AddAccess(z.OCRAddr)
	if err != nil {
		return err
	}

	err = z.SetPayees(z.OCRAddr, z.Payees, z.Transmitters)
	if err != nil {
		return err
	}

	ocrJsonConfig, err := ocrConfigValues.MarshalOCRConfig()
	if err != nil {
		return err
	}
	err = z.SetConfig(z.OCRAddr, ocrJsonConfig)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZKSyncClient) DeployOCRFeed(testEnvironment *environment.Environment, chainClient blockchain.EVMClient, chainlinkNodes []*client.ChainlinkK8sClient, testNetwork blockchain.EVMNetwork, l zerolog.Logger) error {
	z.ChainClient = chainClient
	z.ChainlinkNodes = chainlinkNodes

	var chainlinkClients []*client.ChainlinkClient
	for _, k8sClient := range z.ChainlinkNodes {
		chainlinkClients = append(chainlinkClients, k8sClient.ChainlinkClient)
	}

	err := z.CreateKeys(chainlinkClients)
	if err != nil {
		return err
	}

	err = z.DeployContracts(chainClient, gauntlet.DefaultOcrContract(), gauntlet.DefaultOcrConfig(), l)
	if err != nil {
		return err
	}

	z.Mockserver, err = ctfClient.ConnectMockServer(testEnvironment)
	if err != nil {
		return err
	}
	chainClient.ParallelTransactions(true)

	err = chainClient.WaitForEvents()
	if err != nil {
		return err
	}
	z.OcrInstance = []contracts.OffchainAggregator{
		z.OCRContract,
	}

	// Set Config
	transmitterAddresses, err := actions.ChainlinkNodeAddresses(z.ChainlinkNodes[1:])
	if err != nil {
		return err
	}

	// Exclude the first node, which will be used as a bootstrapper
	err = z.OcrInstance[0].SetConfig(
		z.ChainlinkNodes[1:],
		contracts.DefaultOffChainAggregatorConfig(len(z.ChainlinkNodes[1:])),
		transmitterAddresses,
	)
	if err != nil {
		return err
	}

	bootstrapNode, workerNodes := z.ChainlinkNodes[0], z.ChainlinkNodes[1:]

	err = actions.CreateOCRJobs(z.OcrInstance, bootstrapNode, workerNodes, 5, z.Mockserver, "280")
	if err != nil {
		return err
	}
	return nil
}

func (z *ZKSyncClient) RequestOCRRound(roundNumber int64, value int, l zerolog.Logger) (*big.Int, error) {
	err := actions.SetAllAdapterResponsesToTheSameValue(value, z.OcrInstance, z.ChainlinkNodes, z.Mockserver)
	if err != nil {
		return nil, err
	}
	err = actions.StartNewRound(roundNumber, z.OcrInstance, z.ChainClient, l)
	if err != nil {
		return nil, err
	}

	answer, err := z.OcrInstance[0].GetLatestAnswer(context.Background())
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func SetupOCRTest(t *testing.T) (
	testEnvironment *environment.Environment,
	testNetwork blockchain.EVMNetwork,
	err error,
) {
	var ocrEnvVars = map[string]any{}
	testNetwork = networks.MustGetSelectedNetworksFromEnv()[0]
	evmConfig := helmEth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = helmEth.New(&helmEth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	chainlinkChart := chainlink.New(0, map[string]interface{}{
		"toml":     client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, config.DefaultOCRNetworkDetailTomlConfig, testNetwork),
		"replicas": 6,
	})

	useEnvVars := strings.ToLower(os.Getenv("TEST_USE_ENV_VAR_CONFIG"))
	if useEnvVars == "true" {
		chainlinkChart = chainlink.NewVersioned(0, "0.0.11", map[string]any{
			"replicas": 6,
			"env":      ocrEnvVars,
		})
	}

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix:    fmt.Sprintf("performance-ocr-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:               t,
		PreventPodEviction: true,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlinkChart)
	err = testEnvironment.Run()
	if err != nil {
		return nil, testNetwork, err
	}
	return testEnvironment, testNetwork, nil
}
