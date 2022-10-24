package actions

//revive:disable:dot-imports
import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// DeployBenchmarkKeeperContracts deploys a set amount of keeper Benchmark contracts registered to a single registry
func DeployBenchmarkKeeperContracts(
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registrySettings *contracts.KeeperRegistrySettings,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn, // How much gas should be burned on performUpkeep() calls
	firstEligibleBuffer int64, // How many blocks to add to randomised first eligible block, set to 0 to disable randomised first eligible block
	predeployedContracts []string, // Array of addresses of predeployed consumer addresses to load
	upkeepResetterAddress string,
) (contracts.KeeperRegistry, []contracts.KeeperConsumerBenchmark, []*big.Int) {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(contractDeployer, client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  ZeroAddress.Hex(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        *registrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of KeeperConsumerBenchmark contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployKeeperConsumersBenchmark(contractDeployer, client, numberOfContracts, blockRange, blockInterval, checkGasToBurn, performGasToBurn, firstEligibleBuffer, predeployedContracts, upkeepResetterAddress)

	upkeepsAddresses := []string{}
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	linkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(blockRange/blockInterval))
	gasPrice := big.NewInt(0).Mul(registrySettings.FallbackGasPrice, big.NewInt(2))
	minLinkBalance := big.NewInt(0).
		Add(big.NewInt(0).
			Mul(big.NewInt(0).
				Div(big.NewInt(0).Mul(gasPrice, big.NewInt(int64(upkeepGasLimit+80000))), registrySettings.FallbackLinkPrice),
				big.NewInt(1e18+0)),
			big.NewInt(0))

	linkFunds = big.NewInt(0).Add(linkFunds, minLinkBalance)

	upkeepIds := RegisterUpkeepContracts(linkToken, linkFunds, client, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses)

	return registry, upkeeps, upkeepIds
}

func ResetUpkeeps(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn, // How much gas should be burned on performUpkeep() calls
	firstEligibleBuffer int64, // How many blocks to add to randomised first eligible block
	predeployedContracts []string,
	upkeepResetterAddr string,
) {
	contractLoader, err := contracts.NewContractLoader(client)
	Expect(err).ShouldNot(HaveOccurred(), "Error loading upkeep contract")
	upkeepChunkSize := 500
	upkeepChunks := make([][]string, int(math.Ceil(float64(numberOfContracts)/float64(upkeepChunkSize))))
	upkeepResetter, err := contractLoader.LoadUpkeepResetter(common.HexToAddress(upkeepResetterAddr))
	log.Info().Str("UpkeepResetter Address", upkeepResetter.Address()).Msg("Loaded UpkeepResetter")
	if err != nil {
		upkeepResetter, err = contractDeployer.DeployUpkeepResetter()
		log.Info().Str("UpkeepResetter Address", upkeepResetter.Address()).Msg("Deployed UpkeepResetter")
		if err != nil {
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Upkeep Resetter shouldn't fail")
		}
	}
	iter := 0
	upkeepChunks[iter] = make([]string, 0)
	for count := 0; count < numberOfContracts; count++ {
		if count != 0 && count%upkeepChunkSize == 0 {
			iter++
			upkeepChunks[iter] = make([]string, 0)
		}
		upkeepChunks[iter] = append(upkeepChunks[iter], predeployedContracts[count])
	}
	log.Debug().Int("UpkeepChunk length", len(upkeepChunks))
	for it, upkeepChunk := range upkeepChunks {
		err := upkeepResetter.ResetManyConsumerBenchmark(context.Background(), upkeepChunk, big.NewInt(blockRange),
			big.NewInt(blockInterval), big.NewInt(firstEligibleBuffer), big.NewInt(checkGasToBurn), big.NewInt(performGasToBurn))
		log.Info().Int("Number of Contracts", len(upkeepChunk)).Int("Batch", it).Msg("Resetting batch of Contracts")
		if err != nil {
			Expect(err).ShouldNot(HaveOccurred(), "Resetting upkeeps shouldn't fail")
		}
		err = client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for resetting upkeeps")
	}
}

func DeployKeeperConsumersBenchmark(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn, // How much gas should be burned on performUpkeep() calls
	firstEligibleBuffer int64, // How many blocks to add to randomised first eligible block
	predeployedContracts []string,
	upkeepResetterAddr string,
) []contracts.KeeperConsumerBenchmark {
	upkeeps := make([]contracts.KeeperConsumerBenchmark, 0)

	if len(predeployedContracts) >= numberOfContracts {
		contractLoader, err := contracts.NewContractLoader(client)
		if err != nil {
			log.Error().Err(err).Msg("Loading Contract Loader shouldn't fail")
		}
		for count, address := range predeployedContracts {
			if count < numberOfContracts {
				keeperConsumerInstance, err := contractLoader.LoadKeeperConsumerBenchmark(common.HexToAddress(address))
				if err != nil {
					log.Error().Err(err).Int("count", count+1).Str("UpkeepAddress", address).Msg("Loading KeeperConsumerBenchmark instance shouldn't fail")
					Expect(err).ShouldNot(HaveOccurred(), "Failed to load KeeperConsumerBenchmark")
				}
				upkeeps = append(upkeeps, keeperConsumerInstance)
			}
		}
		// Reset upkeeps so that they are not eligible when being registered
		ResetUpkeeps(contractDeployer, client, numberOfContracts, blockRange, blockInterval, checkGasToBurn,
			performGasToBurn, 10000, predeployedContracts, upkeepResetterAddr)
		return upkeeps
	}

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerBenchmark(
			big.NewInt(blockRange),
			big.NewInt(blockInterval),
			big.NewInt(checkGasToBurn),
			big.NewInt(performGasToBurn),
			big.NewInt(firstEligibleBuffer),
		)
		if err != nil {
			log.Error().Err(err).Int("count", contractCount+1).Msg("Deploying KeeperConsumerBenchmark instance %d shouldn't fail")
			keeperConsumerInstance, err = contractDeployer.DeployKeeperConsumerBenchmark(
				big.NewInt(blockRange),
				big.NewInt(blockInterval),
				big.NewInt(checkGasToBurn),
				big.NewInt(performGasToBurn),
				big.NewInt(firstEligibleBuffer),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Error deploying KeeperConsumerBenchmark")
		}
		//Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerBenchmark instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Benchmark Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumerBenchmark deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}
