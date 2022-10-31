package actions

//revive:disable:dot-imports
import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	types2 "github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func BuildOCRConfigVars(chainlinkNodes []*client.Chainlink, registryConfig contracts.KeeperRegistrySettings, registrar string) contracts.OCRConfig {
	S, oracleIdentities := getOracleIdentities(chainlinkNodes)

	signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		10*time.Second,       // deltaProgress time.Duration,
		10*time.Second,       // deltaResend time.Duration,
		5*time.Second,        // deltaRound time.Duration,
		500*time.Millisecond, // deltaGrace time.Duration,
		2*time.Second,        // deltaStage time.Duration,
		3,                    // rMax uint8,
		S,                    // s []int,
		oracleIdentities,     // oracles []OracleIdentityExtra,
		types2.OffchainConfig{
			PerformLockoutWindow: 100 * 12 * 1000, // ~100 block lockout (on goerli)
			UniqueReports:        false,           // set quorum requirements
		}.Encode(), // reportingPluginConfig []byte,
		50*time.Millisecond, // maxDurationQuery time.Duration,
		time.Second,         // maxDurationObservation time.Duration,
		time.Second,         // maxDurationReport time.Duration,
		time.Second,         // maxDurationShouldAcceptFinalizedReport time.Duration,
		time.Second,         // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                   // f int,
		nil,                 // onchainConfig []byte,
	)
	Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail ContractSetConfigArgsForTests")

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		Expect(len(signer) != 20).ShouldNot(HaveOccurred(), fmt.Errorf("OnChainPublicKey has wrong length for address"))
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		Expect(!common.IsHexAddress(string(transmitter))).ShouldNot(HaveOccurred(), fmt.Errorf("TransmitAccount is not a valid Ethereum address"))
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	onchainConfig, err := encodeOnChainConfig(registryConfig, registrar)
	Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail encoding config")

	return contracts.OCRConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func encodeOnChainConfig(registryConfig contracts.KeeperRegistrySettings, registrar string) ([]byte, error) {
	configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address registrar)")
	onchainConfig, err := abi.Encode(map[string]interface{}{
		"paymentPremiumPPB":    registryConfig.PaymentPremiumPPB,
		"flatFeeMicroLink":     registryConfig.FlatFeeMicroLINK,
		"checkGasLimit":        registryConfig.CheckGasLimit,
		"stalenessSeconds":     registryConfig.StalenessSeconds,
		"gasCeilingMultiplier": registryConfig.GasCeilingMultiplier,
		"minUpkeepSpend":       registryConfig.MinUpkeepSpend,
		"maxPerformGas":        registryConfig.MaxPerformGas,
		"maxCheckDataSize":     registryConfig.MaxCheckDataSize,
		"maxPerformDataSize":   registryConfig.MaxPerformDataSize,
		"fallbackGasPrice":     registryConfig.FallbackGasPrice,
		"fallbackLinkPrice":    registryConfig.FallbackLinkPrice,
		"transcoder":           ZeroAddress.Hex(),
		"registrar":            registrar,
	}, configType)
	return onchainConfig, err
}

func getOracleIdentities(chainlinkNodes []*client.Chainlink) ([]int, []confighelper.OracleIdentityExtra) {
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(chainlinkNodes))
	var wg sync.WaitGroup
	for i, cl := range chainlinkNodes {
		wg.Add(1)
		go func(i int, cl *client.Chainlink) {
			defer wg.Done()

			address, err := cl.PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting primary ETH address from OCR node: index %d", i)
			ocr2Keys, err := cl.MustReadOCR2Keys()
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading OCR2 keys from node")
			fmt.Printf("%+v", *ocr2Keys)
			ocr2Config := ocr2Keys.Data[0].Attributes

			keys, err := cl.MustReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from node")
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			Expect(err).ShouldNot(HaveOccurred(), fmt.Errorf("failed to decode %s: %v", ocr2Config.OffChainPublicKey, err))

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			Expect(n != ed25519.PublicKeySize).ShouldNot(HaveOccurred(), fmt.Errorf("wrong num elements copied"))

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			Expect(err).ShouldNot(HaveOccurred(), fmt.Errorf("failed to decode %s: %v", ocr2Config.ConfigPublicKey, err))

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			Expect(n != ed25519.PublicKeySize).ShouldNot(HaveOccurred(), fmt.Errorf("wrong num elements copied"))

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			Expect(err).ShouldNot(HaveOccurred(), fmt.Errorf("failed to decode %s: %v", ocr2Config.OnChainPublicKey, err))

			sharedSecretEncryptionPublicKeys[i] = configPkBytesFixed
			oracleIdentities[i] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(address),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[i] = 1
		}(i, cl)
	}
	wg.Wait()
	return S, oracleIdentities
}

// CreateOCRKeeperJobs bootstraps the first node and to the other nodes sends ocr jobs
func CreateOCRKeeperJobs(chainlinkNodes []*client.Chainlink, mockserver *ctfClient.MockserverClient, registryAddr string, chainID int64) {
	bootstrapNode := chainlinkNodes[0]
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
	_, err = bootstrapNode.MustCreateJob(&client.AutoBootstrapOCR2JobSpec{
		ContractID: registryAddr,
		ChainID:    int(chainID),
	})
	Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bootstrap job on bootstrap node")
	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, "0.0.0.0", 6690)

	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCRKeys()
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
		nodeOCRKeyId := nodeOCRKeys.Data[0].ID

		_, err = chainlinkNodes[nodeIndex].MustCreateJob(&client.AutoOCR2JobSpec{
			ContractID:         registryAddr,           // registryAddr
			OCRKeyBundleID:     nodeOCRKeyId,           // get node ocr2config.ID
			TransmitterID:      nodeTransmitterAddress, // node addr
			P2Pv2Bootstrappers: P2Pv2Bootstrapper,      // bootstrap node key and address <p2p-key>@bootstrap:8000
			ChainID:            int(chainID),
		})
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
	}
}

// DeployAutoOCRRegistryAndRegistrar registry and registrar
func DeployAutoOCRRegistryAndRegistrar(
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	numberOfUpkeeps int,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar) {
	registry := deployRegistry(registryVersion, registrySettings, contractDeployer, client, linkToken)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err := linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfUpkeeps))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrar := deployRegistrar(registryVersion, registry, linkToken, contractDeployer, client)

	return registry, registrar
}

func DeployConsumers(registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, linkToken contracts.LinkToken, contractDeployer contracts.ContractDeployer, client blockchain.EVMClient, numberOfUpkeeps int, linkFundsForEachUpkeep *big.Int, upkeepGasLimit uint32) ([]contracts.KeeperConsumer, []*big.Int) {
	upkeeps := DeployKeeperConsumers(contractDeployer, client, numberOfUpkeeps)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses)
	return upkeeps, upkeepIds
}

func deployRegistrar(registryVersion ethereum.KeeperRegistryVersion, registry contracts.KeeperRegistry, linkToken contracts.LinkToken, contractDeployer contracts.ContractDeployer, client blockchain.EVMClient) contracts.KeeperRegistrar {
	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar, err := contractDeployer.DeployKeeperRegistrar(registryVersion, linkToken.Address(), registrarSettings)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperRegistrar contract shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar to deploy")
	return registrar
}

func deployRegistry(registryVersion ethereum.KeeperRegistryVersion, registrySettings contracts.KeeperRegistrySettings, contractDeployer contracts.ContractDeployer, client blockchain.EVMClient, linkToken contracts.LinkToken) contracts.KeeperRegistry {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	// Deploy the transcoder here, and then set it to the registry
	transcoder := DeployUpkeepTranscoder(contractDeployer, client)
	registry := DeployKeeperRegistry(contractDeployer, client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  transcoder.Address(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        registrySettings,
		},
	)
	return registry
}
