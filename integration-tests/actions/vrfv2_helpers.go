package actions

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type VRFV2EncodedProvingKey [2]*big.Int

// VRFV2JobInfo defines a jobs into and proving key info
type VRFV2JobInfo struct {
	Job            *client.Job
	VRFKey         *client.VRFKey
	ProvingKey     VRFV2EncodedProvingKey
	ProvingKeyHash [32]byte
}

func DeployVRFV2Contracts(
	t *testing.T,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	c blockchain.EVMClient,
	linkEthFeedAddress string,
) (contracts.VRFCoordinatorV2, contracts.VRFConsumerV2, contracts.BlockHashStore) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	require.NoError(t, err, "Error deploying blockhash store")
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedAddress)
	require.NoError(t, err, "Error deploying VRFv2 Coordinator")
	consumer, err := contractDeployer.DeployVRFConsumerV2(linkTokenContract.Address(), coordinator.Address())
	require.NoError(t, err, "Error deploying VRFv2 Consumer")
	err = c.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	return coordinator, consumer, bhs
}

func CreateVRFV2Jobs(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	coordinator contracts.VRFCoordinatorV2,
	c blockchain.EVMClient,
	minIncomingConfirmations int,
) []VRFV2JobInfo {
	jobInfo := make([]VRFV2JobInfo, 0)
	for _, n := range chainlinkNodes {
		vrfKey, err := n.MustCreateVRFKey()
		require.NoError(t, err, "Error creating VRF key")
		log.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.NewV4()
		os := &client.VRFV2TxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		require.NoError(t, err, "Error getting job string")
		oracleAddr, err := n.PrimaryEthAddress()
		require.NoError(t, err, "Error getting node's primary ETH key")
		job, err := n.MustCreateJob(&client.VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddresses:            []string{oracleAddr},
			EVMChainID:               c.GetChainID().String(),
			MinIncomingConfirmations: minIncomingConfirmations,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		require.NoError(t, err, "Error creating VRFv2 job")
		provingKey := VRFV2RegisterProvingKey(t, vrfKey, oracleAddr, coordinator)
		keyHash, err := coordinator.HashOfKey(context.Background(), provingKey)
		require.NoError(t, err, "Error creating a keyHash from the proving key")
		ji := VRFV2JobInfo{
			Job:            job,
			VRFKey:         vrfKey,
			ProvingKey:     provingKey,
			ProvingKeyHash: keyHash,
		}
		jobInfo = append(jobInfo, ji)
	}
	return jobInfo
}

func VRFV2RegisterProvingKey(
	t *testing.T,
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) VRFV2EncodedProvingKey {
	provingKey, err := EncodeOnChainVRFProvingKey(*vrfKey)
	require.NoError(t, err, "Error encoding proving key")
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	require.NoError(t, err, "Error registering proving keys")
	return provingKey
}
