package ocr

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/seth"

	client2 "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func SetupCluster(
	l zerolog.Logger,
	seth *seth.Client,
	workerNodes []*client.ChainlinkK8sClient,
) (common.Address, error) {
	err := actions_seth.FundChainlinkNodesFromRootAddress(l, seth, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(workerNodes), big.NewFloat(3))
	if err != nil {
		return common.Address{}, err
	}
	linkContract, err := contracts.DeployLinkTokenContract(l, seth)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(linkContract.Address()), nil
}

func SetupFeed(
	l zerolog.Logger,
	seth *seth.Client,
	lta common.Address,
	msClient *client2.MockserverClient,
	bootstrapNode *client.ChainlinkK8sClient,
	workerNodes []*client.ChainlinkK8sClient,
) ([]contracts.OffchainAggregator, error) {
	ocrInstances, err := actions_seth.DeployOCRv1Contracts(l, seth, 1, lta, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(workerNodes))
	if err != nil {
		return nil, err
	}
	err = actions.CreateOCRJobs(ocrInstances, bootstrapNode, workerNodes, 5, msClient, fmt.Sprint(seth.ChainID))
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
