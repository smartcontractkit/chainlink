package vrfv2plus

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type VRFV2PlusEncodedProvingKey [2]*big.Int

// VRFV2PlusJobInfo defines a jobs into and proving key info
type VRFV2PlusJobInfo struct {
	Job               *client.Job
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFV2PlusEncodedProvingKey
	KeyHash           [32]byte
}

type VRFV2PlusContracts struct {
	Coordinator      contracts.VRFCoordinatorV2Plus
	BHS              contracts.BlockHashStore
	LoadTestConsumer contracts.VRFv2PlusLoadTestConsumer
}
