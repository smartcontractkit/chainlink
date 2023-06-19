package v2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

// CoordinatorV2_X is an interface that allows us to use the same code for
// both the V2 and V2.5 coordinators.
type CoordinatorV2_X interface {
	Address() common.Address
	ParseRandomWordsRequested(log types.Log) (*RandomWordsRequested, error)
	GetSubscription(opts *bind.CallOpts, subID uint64) (*GetSubscription, error)
	GetConfig(opts *bind.CallOpts) (*GetConfig, error)
	ParseLog(log types.Log) (generated.AbigenLog, error)
}

type RandomWordsRequested struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	V25        *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested
}

func (r *RandomWordsRequested) Raw() types.Log {
	if r.V2 != nil {
		return r.V2.Raw
	}
	return r.V25.Raw
}

func (r *RandomWordsRequested) NumWords() uint32 {
	if r.V2 != nil {
		return r.V2.NumWords
	}
	return r.V25.NumWords
}

func (r *RandomWordsRequested) SubID() uint64 {
	if r.V2 != nil {
		return r.V2.SubId
	}
	return r.V25.SubId
}

func (r *RandomWordsRequested) MinimumRequestConfirmations() uint16 {
	if r.V2 != nil {
		return r.V2.MinimumRequestConfirmations
	}
	return r.V25.MinimumRequestConfirmations
}

func (r *RandomWordsRequested) KeyHash() [32]byte {
	if r.V2 != nil {
		return r.V2.KeyHash
	}
	return r.V25.KeyHash
}

func (r *RandomWordsRequested) RequestID() *big.Int {
	if r.V2 != nil {
		return r.V2.RequestId
	}
	return r.V25.RequestId
}

func (r *RandomWordsRequested) PreSeed() *big.Int {
	if r.V2 != nil {
		return r.V2.PreSeed
	}
	return r.V25.PreSeed
}

func (r *RandomWordsRequested) Sender() common.Address {
	if r.V2 != nil {
		return r.V2.Sender
	}
	return r.V25.Sender
}

func (r *RandomWordsRequested) CallbackGasLimit() uint32 {
	if r.V2 != nil {
		return r.V2.CallbackGasLimit
	}
	return r.V25.CallbackGasLimit
}

func (r *RandomWordsRequested) NativePayment() bool {
	if r.V2 != nil {
		return false
	}
	return r.V25.NativePayment
}

type GetSubscription struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.GetSubscription
	V25        *vrf_coordinator_v2_5.GetSubscription
}

func (s *GetSubscription) Balance() *big.Int {
	if s.V2 != nil {
		return s.V2.Balance
	}
	return s.V25.Balance
}

func (s *GetSubscription) EthBalance() *big.Int {
	if s.V2 != nil {
		panic("EthBalance not supported on V2")
	}
	return s.V25.EthBalance
}

func (s *GetSubscription) Owner() common.Address {
	if s.V2 != nil {
		return s.V2.Owner
	}
	return s.V25.Owner
}

func (s *GetSubscription) Consumers() []common.Address {
	if s.V2 != nil {
		return s.V2.Consumers
	}
	return s.V25.Consumers
}

type GetConfig struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.GetConfig
	V25        *vrf_coordinator_v2_5.SConfig
}

func (c *GetConfig) MinimumRequestConfirmations() uint16 {
	if c.V2 != nil {
		return c.V2.MinimumRequestConfirmations
	}
	return c.V25.MinimumRequestConfirmations
}

func (c *GetConfig) MaxGasLimit() uint32 {
	if c.V2 != nil {
		return c.V2.MaxGasLimit
	}
	return c.V25.MaxGasLimit
}

func (c *GetConfig) GasAfterPaymentCalculation() uint32 {
	if c.V2 != nil {
		return c.V2.GasAfterPaymentCalculation
	}
	return c.V25.GasAfterPaymentCalculation
}

func (c *GetConfig) StalenessSeconds() uint32 {
	if c.V2 != nil {
		return c.V2.StalenessSeconds
	}
	return c.V25.StalenessSeconds
}

type coordinatorV2_X struct {
	v2   *vrf_coordinator_v2.VRFCoordinatorV2
	v2_5 *vrf_coordinator_v2_5.VRFCoordinatorV25
}

// NewCoordinatorV2 returns a CoordinatorV2_X that wraps the given V2 coordinator
// contract.
func NewCoordinatorV2(coordV2 *vrf_coordinator_v2.VRFCoordinatorV2) CoordinatorV2_X {
	return &coordinatorV2_X{v2: coordV2}
}

// NewCoordinatorV2_5 returns a CoordinatorV2_X that wraps the given V2.5 coordinator
// contract.
func NewCoordinatorV2_5(coordV2_5 *vrf_coordinator_v2_5.VRFCoordinatorV25) CoordinatorV2_X {
	return &coordinatorV2_X{v2_5: coordV2_5}
}

func (c *coordinatorV2_X) Address() common.Address {
	if c.v2 != nil {
		return c.v2.Address()
	}
	return c.v2_5.Address()
}

func (c *coordinatorV2_X) ParseRandomWordsRequested(log types.Log) (*RandomWordsRequested, error) {
	if c.v2 != nil {
		parsed, err := c.v2.ParseRandomWordsRequested(log)
		return &RandomWordsRequested{
			VRFVersion: vrfcommon.V2,
			V2:         parsed,
		}, err
	}
	parsed, err := c.v2_5.ParseRandomWordsRequested(log)
	return &RandomWordsRequested{
		VRFVersion: vrfcommon.V2_5,
		V25:        parsed,
	}, err
}

func (c *coordinatorV2_X) GetSubscription(opts *bind.CallOpts, subID uint64) (*GetSubscription, error) {
	if c.v2 != nil {
		parsed, err := c.v2.GetSubscription(opts, subID)
		return &GetSubscription{
			VRFVersion: vrfcommon.V2,
			V2:         &parsed,
		}, err
	}
	parsed, err := c.v2_5.GetSubscription(opts, subID)
	return &GetSubscription{
		VRFVersion: vrfcommon.V2_5,
		V25:        &parsed,
	}, err
}

func (c *coordinatorV2_X) GetConfig(opts *bind.CallOpts) (*GetConfig, error) {
	if c.v2 != nil {
		parsed, err := c.v2.GetConfig(opts)
		return &GetConfig{
			VRFVersion: vrfcommon.V2,
			V2:         &parsed,
		}, err
	}
	parsed, err := c.v2_5.SConfig(opts)
	return &GetConfig{
		VRFVersion: vrfcommon.V2_5,
		V25:        &parsed,
	}, err
}

func (c *coordinatorV2_X) ParseLog(log types.Log) (generated.AbigenLog, error) {
	if c.v2 != nil {
		return c.v2.ParseLog(log)
	}
	return c.v2_5.ParseLog(log)
}

var (
	_ CoordinatorV2_X = (*coordinatorV2_X)(nil)
)
