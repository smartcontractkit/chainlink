package test_env

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// VRFV2JobInfo defines a jobs into and proving key info
type VRFV2JobInfo struct {
	Job               *client.Job
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFV2EncodedProvingKey
	KeyHash           [32]byte
}

func (m *ClNode) CreateVRFv2Job(coordinatorV2 contracts.VRFCoordinatorV2, c *blockchain.EthereumClient) (*VRFV2JobInfo, error) {
	vrfKey, err := m.API.MustCreateVRFKey()
	if err != nil {
		return nil, err
	}
	pubKeyCompressed := vrfKey.Data.ID
	jobUUID := uuid.New()
	spec := &client.VRFV2TxPipelineSpec{
		Address: coordinatorV2.Address(),
	}
	ost, err := spec.String()
	if err != nil {
		return nil, err
	}
	nativeTokenPrimaryKeyAddress, err := m.API.PrimaryEthAddress()
	if err != nil {
		return nil, err
	}
	job, err := m.API.MustCreateJob(&client.VRFV2JobSpec{
		Name:                     fmt.Sprintf("vrf-%s", jobUUID),
		CoordinatorAddress:       coordinatorV2.Address(),
		FromAddresses:            []string{nativeTokenPrimaryKeyAddress},
		EVMChainID:               c.GetChainID().String(),
		MinIncomingConfirmations: 1,
		PublicKey:                pubKeyCompressed,
		ExternalJobID:            jobUUID.String(),
		ObservationSource:        ost,
		BatchFulfillmentEnabled:  false,
	})
	if err != nil {
		return nil, err
	}
	provingKey, err := VRFV2RegisterProvingKey(vrfKey, nativeTokenPrimaryKeyAddress, coordinatorV2)
	keyHash, err := coordinatorV2.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreatingProvingKey)
	}
	return &VRFV2JobInfo{
		Job:               job,
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
	}, nil
}
