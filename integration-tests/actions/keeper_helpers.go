package actions

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var ZeroAddress = common.Address{}

func CreateKeeperJobsWithKeyIndex(
	t *testing.T,
	chainlinkNodes []*client.ChainlinkK8sClient,
	keeperRegistry contracts.KeeperRegistry,
	keyIndex int,
	ocrConfig contracts.OCRv2Config,
	evmChainID string,
) {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddresses, err := primaryNode.EthAddresses()
	require.NoError(t, err, "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddressesAtIndex(chainlinkNodes, keyIndex)
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddresses[keyIndex])
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees, ocrConfig)
	require.NoError(t, err, "Setting keepers in the registry shouldn't fail")

	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.EthAddresses()
		require.NoError(t, err, "Error retrieving chainlink node address")
		_, err = chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress[keyIndex],
			EVMChainID:               evmChainID,
			MinIncomingConfirmations: 1,
		})
		require.NoError(t, err, "Creating KeeperV2 Job shouldn't fail")
	}
}

func DeleteKeeperJobsWithId(t *testing.T, chainlinkNodes []*client.ChainlinkK8sClient, id int) {
	for _, chainlinkNode := range chainlinkNodes {
		err := chainlinkNode.MustDeleteJob(strconv.Itoa(id))
		require.NoError(t, err, "Deleting KeeperV2 Job shouldn't fail")
	}
}
