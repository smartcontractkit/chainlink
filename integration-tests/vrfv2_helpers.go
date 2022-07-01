package integration_tests

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"

	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
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
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	c blockchain.EVMClient,
	linkEthFeedAddress string,
) (contracts.VRFCoordinatorV2, contracts.VRFConsumerV2, contracts.BlockHashStore) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	Expect(err).ShouldNot(HaveOccurred())
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedAddress)
	Expect(err).ShouldNot(HaveOccurred())
	consumer, err := contractDeployer.DeployVRFConsumerV2(linkTokenContract.Address(), coordinator.Address())
	Expect(err).ShouldNot(HaveOccurred())
	err = c.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())

	return coordinator, consumer, bhs
}

func CreateVRFV2Jobs(
	chainlinkNodes []client.Chainlink,
	coordinator contracts.VRFCoordinatorV2,
	c blockchain.EVMClient,
	minIncomingConfirmations int,
) []VRFV2JobInfo {
	jobInfo := make([]VRFV2JobInfo, 0)
	for _, n := range chainlinkNodes {
		vrfKey, err := n.CreateVRFKey()
		Expect(err).ShouldNot(HaveOccurred())
		log.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.NewV4()
		os := &VRFV2TxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())
		oracleAddr, err := n.PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred())
		job, err := n.CreateJob(&VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddress:              oracleAddr,
			EVMChainID:               c.GetChainID().String(),
			MinIncomingConfirmations: minIncomingConfirmations,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		provingKey := VRFV2RegisterProvingKey(vrfKey, oracleAddr, coordinator)
		keyHash, err := coordinator.HashOfKey(context.Background(), provingKey)
		Expect(err).ShouldNot(HaveOccurred(), "Should be able to create a keyHash from the proving key")
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
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) VRFV2EncodedProvingKey {
	provingKey, err := EncodeOnChainVRFProvingKey(*vrfKey)
	Expect(err).ShouldNot(HaveOccurred())
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	Expect(err).ShouldNot(HaveOccurred())
	return provingKey
}
