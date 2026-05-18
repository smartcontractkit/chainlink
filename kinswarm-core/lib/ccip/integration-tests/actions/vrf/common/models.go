package common

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

type VRFEncodedProvingKey [2]*big.Int

// VRFV2PlusKeyData defines a jobs into and proving key info
type VRFKeyData struct {
	VRFKey            *client.VRFKey
	EncodedProvingKey VRFEncodedProvingKey
	KeyHash           [32]byte
	PubKeyCompressed  string
}

type VRFNodeType int

const (
	VRF VRFNodeType = iota + 1
	BHS
	BHF
)

func (n VRFNodeType) String() string {
	return [...]string{"VRF", "BHS", "BHF"}[n-1]
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
	CoordinatorV2          contracts.VRFCoordinatorV2
	BatchCoordinatorV2     contracts.BatchVRFCoordinatorV2
	CoordinatorV2Plus      contracts.VRFCoordinatorV2_5
	BatchCoordinatorV2Plus contracts.BatchVRFCoordinatorV2Plus
	VRFOwner               contracts.VRFOwner
	BHS                    contracts.BlockHashStore
	BatchBHS               contracts.BatchBlockhashStore
	VRFV2Consumers         []contracts.VRFv2LoadTestConsumer
	VRFV2PlusConsumer      []contracts.VRFv2PlusLoadTestConsumer
	LinkToken              contracts.LinkToken
	MockETHLINKFeed        contracts.VRFMockETHLINKFeed
	LinkNativeFeedAddress  string
}

type VRFOwnerConfig struct {
	OwnerAddress string
	UseVRFOwner  bool
}

type VRFJobSpecConfig struct {
	ForwardingAllowed             bool
	CoordinatorAddress            string
	BatchCoordinatorAddress       string
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

type NewEnvConfig struct {
	NodesToCreate                   []VRFNodeType
	NumberOfTxKeysToCreate          int
	UseVRFOwner                     bool
	UseTestCoordinator              bool
	ChainlinkNodeLogScannerSettings test_env.ChainlinkNodeLogScannerSettings
}

type VRFEnvConfig struct {
	TestConfig tc.TestConfig
	ChainID    int64
	CleanupFn  func()
}
