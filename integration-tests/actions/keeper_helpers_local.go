package actions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func CreateKeeperJobsLocal(
	t *testing.T,
	chainlinkNodes []*client.ChainlinkClient,
	keeperRegistry contracts.KeeperRegistry,
	ocrConfig contracts.OCRv2Config,
) []*client.Job {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	require.NoError(t, err, "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddressesLocal(chainlinkNodes)
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees, ocrConfig)
	require.NoError(t, err, "Setting keepers in the registry shouldn't fail")
	jobs := []*client.Job{}
	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		require.NoError(t, err, "Error retrieving chainlink node address")
		job, err := chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress,
			MinIncomingConfirmations: 1,
		})
		require.NoError(t, err, "Creating KeeperV2 Job shouldn't fail")
		jobs = append(jobs, job)
	}
	return jobs
}
