package ocr

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	client2 "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func SetupCluster(
	cc blockchain.EVMClient,
	cd contracts.ContractDeployer,
	workerNodes []*client.ChainlinkK8sClient,
) (contracts.LinkToken, error) {
	err := actions.FundChainlinkNodes(workerNodes, cc, big.NewFloat(3))
	if err != nil {
		return nil, err
	}
	lt, err := cd.DeployLinkTokenContract()
	if err != nil {
		return nil, err
	}
	return lt, nil
}

func SetupFeed(
	cc blockchain.EVMClient,
	msClient *client2.MockserverClient,
	cd contracts.ContractDeployer,
	bootstrapNode *client.ChainlinkK8sClient,
	workerNodes []*client.ChainlinkK8sClient,
	lt contracts.LinkToken,
) ([]contracts.OffchainAggregator, error) {
	ocrInstances, err := actions.DeployOCRContracts(1, lt, cd, workerNodes, cc)
	if err != nil {
		return nil, err
	}
	err = actions.CreateOCRJobs(ocrInstances, bootstrapNode, workerNodes, 5, msClient, cc.GetChainID().String())
	if err != nil {
		return nil, err
	}
	return ocrInstances, nil
}

func SimulateEAActivity(
	l zerolog.Logger,
	eaChangeInterval time.Duration,
	ocrInstances []contracts.OffchainAggregator,
	workerNodes []*client.ChainlinkK8sClient,
	msClient *client2.MockserverClient,
) {
	go func() {
		for {
			time.Sleep(eaChangeInterval)
			if err := actions.SetAllAdapterResponsesToTheSameValue(rand.Intn(1000), ocrInstances, workerNodes, msClient); err != nil {
				l.Error().Err(err).Msg("failed to update mockserver responses")
			}
		}
	}()
}
