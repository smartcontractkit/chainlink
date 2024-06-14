package vrfv2plus

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type VRFV2PlusEncodedProvingKey [2]*big.Int

// VRFV2PlusKeyData defines a jobs into and proving key info
type VRFV2PlusKeyData struct {
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFV2PlusEncodedProvingKey
	KeyHash           [32]byte
}

type VRFV2PlusData struct {
	VRFV2PlusKeyData
	VRFJob            *client.Job
	PrimaryEthAddress string
	ChainID           *big.Int
}

type VRFV2_5Contracts struct {
	Coordinator       contracts.VRFCoordinatorV2_5
	BHS               contracts.BlockHashStore
	LoadTestConsumers []contracts.VRFv2PlusLoadTestConsumer
}

type VRFV2PlusWrapperContracts struct {
	VRFV2PlusWrapper  contracts.VRFV2PlusWrapper
	LoadTestConsumers []contracts.VRFv2PlusWrapperLoadTestConsumer
}
