package common

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

type VRFEncodedProvingKey [2]*big.Int

// VRFV2PlusKeyData defines a jobs into and proving key info
type VRFKeyData struct {
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFEncodedProvingKey
	KeyHash           [32]byte
}

type VRFNodeType int

const (
	VRF VRFNodeType = iota + 1
	BHS
)

func (n VRFNodeType) String() string {
	return [...]string{"VRF", "BHS"}[n-1]
}

func (n VRFNodeType) Index() int {
	return int(n)
}

type VRFNode struct {
	CLNode              *test_env.ClNode
	Job                 *client.Job
	TXKeyAddressStrings []string
}

type VRFContracts struct {
	CoordinatorV2     contracts.VRFCoordinatorV2
	CoordinatorV2Plus contracts.VRFCoordinatorV2_5
	VRFOwner          contracts.VRFOwner
	BHS               contracts.BlockHashStore
	VRFV2Consumer     []contracts.VRFv2LoadTestConsumer
	VRFV2PlusConsumer []contracts.VRFv2PlusLoadTestConsumer
}

type VRFOwnerConfig struct {
	OwnerAddress string
	UseVRFOwner  bool
}

type VRFJobSpecConfig struct {
	ForwardingAllowed             bool
	CoordinatorAddress            string
	FromAddresses                 []string
	EVMChainID                    string
	MinIncomingConfirmations      int
	PublicKey                     string
	BatchFulfillmentEnabled       bool
	BatchFulfillmentGasMultiplier float64
	EstimateGasMultiplier         float64
	PollPeriod                    time.Duration
	RequestTimeout                time.Duration
	VRFOwnerConfig                *VRFOwnerConfig
	SimulationBlock               *string
}

type VRFLoadTestConsumer interface {
	GetLoadTestMetrics(ctx context.Context) (*contracts.VRFLoadTestMetrics, error)
}
