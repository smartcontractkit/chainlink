package v2

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
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
	LogsWithTopics(keyHash common.Hash) map[common.Hash][][]log.Topic
	Version() vrfcommon.Version
}

type RandomWordsRequested struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	V2Plus     *vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsRequested
}

func (r *RandomWordsRequested) Raw() types.Log {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.Raw
	}
	return r.V2Plus.Raw
}

func (r *RandomWordsRequested) NumWords() uint32 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.NumWords
	}
	return r.V2Plus.NumWords
}

func (r *RandomWordsRequested) SubID() uint64 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.SubId
	}
	return r.V2Plus.SubId
}

func (r *RandomWordsRequested) MinimumRequestConfirmations() uint16 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.MinimumRequestConfirmations
	}
	return r.V2Plus.MinimumRequestConfirmations
}

func (r *RandomWordsRequested) KeyHash() [32]byte {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.KeyHash
	}
	return r.V2Plus.KeyHash
}

func (r *RandomWordsRequested) RequestID() *big.Int {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.RequestId
	}
	return r.V2Plus.RequestId
}

func (r *RandomWordsRequested) PreSeed() *big.Int {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.PreSeed
	}
	return r.V2Plus.PreSeed
}

func (r *RandomWordsRequested) Sender() common.Address {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.Sender
	}
	return r.V2Plus.Sender
}

func (r *RandomWordsRequested) CallbackGasLimit() uint32 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.CallbackGasLimit
	}
	return r.V2Plus.CallbackGasLimit
}

func (r *RandomWordsRequested) NativePayment() bool {
	if r.VRFVersion == vrfcommon.V2 {
		return false
	}
	return r.V2Plus.NativePayment
}

type GetSubscription struct {
	VRFVersion vrfcommon.Version
	V2         vrf_coordinator_v2.GetSubscription
	V2Plus     vrf_coordinator_v2plus.GetSubscription
}

func (s *GetSubscription) Balance() *big.Int {
	if s.VRFVersion == vrfcommon.V2 {
		return s.V2.Balance
	}
	return s.V2Plus.Balance
}

func (s *GetSubscription) EthBalance() *big.Int {
	if s.VRFVersion == vrfcommon.V2 {
		panic("EthBalance not supported on V2")
	}
	return s.V2Plus.EthBalance
}

func (s *GetSubscription) Owner() common.Address {
	if s.VRFVersion == vrfcommon.V2 {
		return s.V2.Owner
	}
	return s.V2Plus.Owner
}

func (s *GetSubscription) Consumers() []common.Address {
	if s.VRFVersion == vrfcommon.V2 {
		return s.V2.Consumers
	}
	return s.V2Plus.Consumers
}

type GetConfig struct {
	VRFVersion vrfcommon.Version
	V2         vrf_coordinator_v2.GetConfig
	V2Plus     vrf_coordinator_v2plus.SConfig
}

func (c *GetConfig) MinimumRequestConfirmations() uint16 {
	if c.VRFVersion == vrfcommon.V2 {
		return c.V2.MinimumRequestConfirmations
	}
	return c.V2Plus.MinimumRequestConfirmations
}

func (c *GetConfig) MaxGasLimit() uint32 {
	if c.VRFVersion == vrfcommon.V2 {
		return c.V2.MaxGasLimit
	}
	return c.V2Plus.MaxGasLimit
}

func (c *GetConfig) GasAfterPaymentCalculation() uint32 {
	if c.VRFVersion == vrfcommon.V2 {
		return c.V2.GasAfterPaymentCalculation
	}
	return c.V2Plus.GasAfterPaymentCalculation
}

func (c *GetConfig) StalenessSeconds() uint32 {
	if c.VRFVersion == vrfcommon.V2 {
		return c.V2.StalenessSeconds
	}
	return c.V2Plus.StalenessSeconds
}

type RequestCommitment struct {
	VRFVersion vrfcommon.Version
	V2         vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment
	V2Plus     vrf_coordinator_v2plus.VRFCoordinatorV2PlusRequestCommitment
}

func ToV2Commitments(commitments []RequestCommitment) []vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment {
	v2Commitments := make([]vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment, len(commitments))
	for i, commitment := range commitments {
		v2Commitments[i] = commitment.V2
	}
	return v2Commitments
}

func ToV2PlusCommitments(commitments []RequestCommitment) []vrf_coordinator_v2plus.VRFCoordinatorV2PlusRequestCommitment {
	v2PlusCommitments := make([]vrf_coordinator_v2plus.VRFCoordinatorV2PlusRequestCommitment, len(commitments))
	for i, commitment := range commitments {
		v2PlusCommitments[i] = commitment.V2Plus
	}
	return v2PlusCommitments
}

func NewRequestCommitment(val any) RequestCommitment {
	switch val := val.(type) {
	case vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment:
		return RequestCommitment{VRFVersion: vrfcommon.V2, V2: val}
	case vrf_coordinator_v2plus.VRFCoordinatorV2PlusRequestCommitment:
		return RequestCommitment{VRFVersion: vrfcommon.V2Plus, V2Plus: val}
	default:
		panic(fmt.Sprintf("NewRequestCommitment: unknown type %T", val))
	}
}

func (r *RequestCommitment) Get() any {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2
	}
	return r.V2Plus
}

func (r *RequestCommitment) NativePayment() bool {
	if r.VRFVersion == vrfcommon.V2 {
		return false
	}
	return r.V2Plus.NativePayment
}

func (r *RequestCommitment) NumWords() uint32 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.NumWords
	}
	return r.V2Plus.NumWords
}

func (r *RequestCommitment) Sender() common.Address {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.Sender
	}
	return r.V2Plus.Sender
}

func (r *RequestCommitment) BlockNum() uint64 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.BlockNum
	}
	return r.V2Plus.BlockNum
}

func (r *RequestCommitment) SubID() uint64 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.SubId
	}
	return r.V2Plus.SubId
}

func (r *RequestCommitment) CallbackGasLimit() uint32 {
	if r.VRFVersion == vrfcommon.V2 {
		return r.V2.CallbackGasLimit
	}
	return r.V2Plus.CallbackGasLimit
}

type coordinatorV2_X struct {
	v2     *vrf_coordinator_v2.VRFCoordinatorV2
	v2plus *vrf_coordinator_v2plus.VRFCoordinatorV2Plus
}

// NewCoordinatorV2 returns a CoordinatorV2_X that wraps the given V2 coordinator
// contract.
func NewCoordinatorV2(coordV2 *vrf_coordinator_v2.VRFCoordinatorV2) CoordinatorV2_X {
	return &coordinatorV2_X{v2: coordV2}
}

// NewCoordinatorV2Plus returns a CoordinatorV2_X that wraps the given V2.5 coordinator
// contract.
func NewCoordinatorV2Plus(coordV2Plus *vrf_coordinator_v2plus.VRFCoordinatorV2Plus) CoordinatorV2_X {
	return &coordinatorV2_X{v2plus: coordV2Plus}
}

func (c *coordinatorV2_X) Version() vrfcommon.Version {
	if c.v2 != nil {
		return vrfcommon.V2
	}
	return vrfcommon.V2Plus
}

func (c *coordinatorV2_X) LogsWithTopics(keyHash common.Hash) map[common.Hash][][]log.Topic {
	if c.v2 != nil {
		return map[common.Hash][][]log.Topic{
			vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(): {
				{
					log.Topic(keyHash),
				},
			},
		}
	}
	return map[common.Hash][][]log.Topic{
		vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsRequested{}.Topic(): {
			{
				log.Topic(keyHash),
			},
		},
	}
}

func (c *coordinatorV2_X) Address() common.Address {
	if c.v2 != nil {
		return c.v2.Address()
	}
	return c.v2plus.Address()
}

func (c *coordinatorV2_X) ParseRandomWordsRequested(log types.Log) (*RandomWordsRequested, error) {
	if c.v2 != nil {
		parsed, err := c.v2.ParseRandomWordsRequested(log)
		return &RandomWordsRequested{
			VRFVersion: vrfcommon.V2,
			V2:         parsed,
		}, err
	}
	parsed, err := c.v2plus.ParseRandomWordsRequested(log)
	return &RandomWordsRequested{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     parsed,
	}, err
}

func (c *coordinatorV2_X) GetSubscription(opts *bind.CallOpts, subID uint64) (*GetSubscription, error) {
	if c.v2 != nil {
		sub, err := c.v2.GetSubscription(opts, subID)
		return &GetSubscription{
			VRFVersion: vrfcommon.V2,
			V2:         sub,
		}, err
	}
	sub, err := c.v2plus.GetSubscription(opts, subID)
	return &GetSubscription{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     sub,
	}, err
}

func (c *coordinatorV2_X) GetConfig(opts *bind.CallOpts) (*GetConfig, error) {
	if c.v2 != nil {
		cfg, err := c.v2.GetConfig(opts)
		return &GetConfig{
			VRFVersion: vrfcommon.V2,
			V2:         cfg,
		}, err
	}
	cfg, err := c.v2plus.SConfig(opts)
	return &GetConfig{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     cfg,
	}, err
}

func (c *coordinatorV2_X) ParseLog(log types.Log) (generated.AbigenLog, error) {
	if c.v2 != nil {
		return c.v2.ParseLog(log)
	}
	return c.v2plus.ParseLog(log)
}

var (
	_ CoordinatorV2_X = (*coordinatorV2_X)(nil)
)
