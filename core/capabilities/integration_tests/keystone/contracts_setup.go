package keystone

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
)

func SetupForwarderContract(t *testing.T, reportCreator *framework.DON,
	backend *framework.EthBlockchain) (common.Address, *forwarder.KeystoneForwarder) {
	addr, _, fwd, err := forwarder.DeployKeystoneForwarder(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	signers := make([]common.Address, 0, len(reportCreator.GetPeerIDs()))
	for _, p := range reportCreator.GetPeerIDs() {
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
