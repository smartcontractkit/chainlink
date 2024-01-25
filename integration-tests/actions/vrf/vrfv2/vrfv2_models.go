package vrfv2

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
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
	Coordinator       contracts.VRFCoordinatorV2
	VRFOwner          contracts.VRFOwner
	BHS               contracts.BlockHashStore
	LoadTestConsumers []contracts.VRFv2LoadTestConsumer
}

type VRFV2WrapperContracts struct {
	VRFV2Wrapper      contracts.VRFV2Wrapper
	LoadTestConsumers []contracts.VRFv2WrapperLoadTestConsumer
}

// VRFV2PlusKeyData defines a jobs into and proving key info
type VRFV2KeyData struct {
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFV2EncodedProvingKey
	KeyHash           [32]byte
}

type VRFV2Data struct {
	VRFV2KeyData
	PrimaryEthAddress string
	ChainID           *big.Int
}

type VRFNodeType int

const (
	VRF VRFNodeType = iota + 1
	VRF_Backup
	BHS
	BHS_Backup
	BHF
)

func (n VRFNodeType) String() string {
	return [...]string{"VRF", "VRF_Backup", "BHS", "BHS_Backup", "BHF"}[n-1]
}

func (n VRFNodeType) Index() int {
	return int(n)
}

type VRFNode struct {
	CLNode *test_env.ClNode
	Job    *client.Job
}
