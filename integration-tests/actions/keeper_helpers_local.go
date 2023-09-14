package actions

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func CreateKeeperJobsLocal(
	l zerolog.Logger,
	chainlinkNodes []*client.ChainlinkClient,
	keeperRegistry contracts.KeeperRegistry,
	ocrConfig contracts.OCRv2Config,
	evmChainID string,
) ([]*client.Job, error) {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	if err != nil {
		l.Error().Err(err).Msg("Reading ETH Keys from Chainlink Client shouldn't fail")
		return nil, err
	}
	nodeAddresses, err := ChainlinkNodeAddressesLocal(chainlinkNodes)
	if err != nil {
		l.Error().Err(err).Msg("Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
		return nil, err
	}
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees, ocrConfig)
	if err != nil {
		l.Error().Err(err).Msg("Setting keepers in the registry shouldn't fail")
		return nil, err
	}
	jobs := []*client.Job{}
	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		if err != nil {
			l.Error().Err(err).Msg("Error retrieving chainlink node address")
			return nil, err
		}
		job, err := chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress,
			EVMChainID:               evmChainID,
			MinIncomingConfirmations: 1,
		})
		if err != nil {
			l.Error().Err(err).Msg("Creating KeeperV2 Job shouldn't fail")
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
