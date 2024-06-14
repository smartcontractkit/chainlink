package vrfv2_actions

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type VRFV2EncodedProvingKey [2]*big.Int

// VRFV2JobInfo defines a jobs into and proving key info
type VRFV2JobInfo struct {
	Job               *client.Job
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFV2EncodedProvingKey
	KeyHash           [32]byte
}

type VRFV2Contracts struct {
	Coordinator      contracts.VRFCoordinatorV2
	BHS              contracts.BlockHashStore
	LoadTestConsumer contracts.VRFv2LoadTestConsumer
}
