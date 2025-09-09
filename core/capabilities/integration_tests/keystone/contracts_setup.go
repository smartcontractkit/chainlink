package keystone

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	data_feeds_cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
)

func SetupForwarderContract(t *testing.T, reportCreator *framework.DON,
	backend *framework.EthBlockchain) (common.Address, *forwarder.KeystoneForwarder) {
	addr, _, fwd, err := forwarder.DeployKeystoneForwarder(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	signers := make([]common.Address, 0, len(reportCreator.GetPeerIDsAndOCRSigners()))
	for _, p := range reportCreator.GetPeerIDsAndOCRSigners() {
		signers = append(signers, common.HexToAddress(p.Signer))
	}

	_, err = fwd.SetConfig(backend.TransactionOpts(), reportCreator.GetID(), reportCreator.GetConfigVersion(), reportCreator.GetF(), signers)
	require.NoError(t, err)
	backend.Commit()

	return addr, fwd
}

func SetupConsumerContract(t *testing.T, backend *framework.EthBlockchain,
	forwarderAddress common.Address, workflowOwner string, workflowName string) (common.Address, *feeds_consumer.KeystoneFeedsConsumer) {
	addr, _, consumer, err := feeds_consumer.DeployKeystoneFeedsConsumer(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = consumer.SetConfig(backend.TransactionOpts(), []common.Address{forwarderAddress}, []common.Address{ownerAddr}, [][10]byte{nameBytes})
	require.NoError(t, err)

	backend.Commit()

	return addr, consumer
}

// secureMintFeedDataID is the Data ID generated for a secure mint feed using the feed_id.go generator in CLD
const secureMintFeedDataID = "0x01c508f42b0201320000000000000000"

func SetupDataFeedsCacheContract(t *testing.T, backend *framework.EthBlockchain,
	forwarderAddress common.Address, workflowOwner string, workflowName string) (common.Address, *data_feeds_cache.DataFeedsCache) {
	addr, _, dataFeedsCache, err := data_feeds_cache.DeployDataFeedsCache(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = dataFeedsCache.SetFeedAdmin(backend.TransactionOpts(), backend.TransactionOpts().From, true)
	require.NoError(t, err)
	backend.Commit()

	feedIDBytes := [16]byte{}
	copy(feedIDBytes[:], common.FromHex(secureMintFeedDataID))

	_, err = dataFeedsCache.SetDecimalFeedConfigs(backend.TransactionOpts(), [][16]byte{feedIDBytes}, []string{"securemint"},
		[]data_feeds_cache.DataFeedsCacheWorkflowMetadata{
			{
				AllowedSender:        forwarderAddress,
				AllowedWorkflowOwner: ownerAddr,
				AllowedWorkflowName:  nameBytes,
			},
		})
	require.NoError(t, err)
	backend.Commit()

	return addr, dataFeedsCache
}
