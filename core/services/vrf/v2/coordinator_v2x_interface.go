package v2

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var extraArgsV1Tag = crypto.Keccak256([]byte("VRF ExtraArgsV1"))[:4]

// CoordinatorV2_X is an interface that allows us to use the same code for
// both the V2 and V2Plus coordinators.
type CoordinatorV2_X interface {
	Address() common.Address
	ParseRandomWordsRequested(log types.Log) (*RandomWordsRequested, error)
	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error)
	AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)
	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)
	GetSubscription(opts *bind.CallOpts, subID uint64) (*GetSubscription, error)
	GetConfig(opts *bind.CallOpts) (*GetConfig, error)
	ParseLog(log types.Log) (generated.AbigenLog, error)
	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)
	LogsWithTopics(keyHash common.Hash) map[common.Hash][][]log.Topic
	Version() vrfcommon.Version
	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)
	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*SubscriptionCreatedIterator, error)
	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*RandomWordsRequestedIterator, error)
	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestID []*big.Int) (*RandomWordsFulfilledIterator, error)
	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
	RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)
	CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error)
	GetCommitment(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error)
}

type RandomWordsRequestedIterator struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequestedIterator
	V2Plus     *vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsRequestedIterator
}

func (it *RandomWordsRequestedIterator) Next() bool {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Next()
	}
	return it.V2Plus.Next()
}

func (it *RandomWordsRequestedIterator) Error() error {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Error()
	}
	return it.V2Plus.Error()
}

func (it *RandomWordsRequestedIterator) Close() error {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Close()
	}
	return it.V2Plus.Close()
}

func (it *RandomWordsRequestedIterator) Event() *RandomWordsRequested {
	if it.VRFVersion == vrfcommon.V2 {
		return &RandomWordsRequested{
			VRFVersion: it.VRFVersion,
			V2:         it.V2.Event,
		}
	}
	return &RandomWordsRequested{
		VRFVersion: it.VRFVersion,
		V2Plus:     it.V2Plus.Event,
	}
}

type RandomWordsFulfilledIterator struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilledIterator
	V2Plus     *vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsFulfilledIterator
}

func (it *RandomWordsFulfilledIterator) Next() bool {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Next()
	}
	return it.V2Plus.Next()
}

func (it *RandomWordsFulfilledIterator) Error() error {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Error()
	}
	return it.V2Plus.Error()
}

func (it *RandomWordsFulfilledIterator) Close() error {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Close()
	}
	return it.V2Plus.Close()
}

func (it *RandomWordsFulfilledIterator) Event() *RandomWordsFulfilled {
	if it.VRFVersion == vrfcommon.V2 {
		return &RandomWordsFulfilled{
			VRFVersion: it.VRFVersion,
			V2:         it.V2.Event,
		}
	}
	return &RandomWordsFulfilled{
		VRFVersion: it.VRFVersion,
		V2Plus:     it.V2Plus.Event,
	}
}

type RandomWordsFulfilled struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	V2Plus     *vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsFulfilled
}

func (rwf *RandomWordsFulfilled) RequestID() *big.Int {
	if rwf.VRFVersion == vrfcommon.V2 {
		return rwf.V2.RequestId
	}
	return rwf.V2Plus.RequestId
}

func (rwf *RandomWordsFulfilled) Success() bool {
	if rwf.VRFVersion == vrfcommon.V2 {
		return rwf.V2.Success
	}
	return rwf.V2Plus.Success
}

func (rwf *RandomWordsFulfilled) NativePayment() bool {
	if rwf.VRFVersion == vrfcommon.V2 {
		return false
	}
	return rwf.V2Plus.NativePayment
}

func (rwf *RandomWordsFulfilled) Payment() *big.Int {
	if rwf.VRFVersion == vrfcommon.V2 {
		return rwf.V2.Payment
	}
	return rwf.V2Plus.Payment
}

func (rwf *RandomWordsFulfilled) Raw() types.Log {
	if rwf.VRFVersion == vrfcommon.V2 {
		return rwf.V2.Raw
	}
	return rwf.V2Plus.Raw
}

type SubscriptionCreatedIterator struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreatedIterator
	V2Plus     *vrf_coordinator_v2plus.VRFCoordinatorV2PlusSubscriptionCreatedIterator
}

func (it *SubscriptionCreatedIterator) Next() bool {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Next()
	}
	return it.V2Plus.Next()
}

func (it *SubscriptionCreatedIterator) Error() error {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Error()
	}
	return it.V2Plus.Error()
}

func (it *SubscriptionCreatedIterator) Close() error {
	if it.VRFVersion == vrfcommon.V2 {
		return it.V2.Close()
	}
	return it.V2Plus.Close()
}

func (it *SubscriptionCreatedIterator) Event() *SubscriptionCreated {
	if it.VRFVersion == vrfcommon.V2 {
		return &SubscriptionCreated{
			VRFVersion: it.VRFVersion,
			V2:         it.V2.Event,
		}
	}
	return &SubscriptionCreated{
		VRFVersion: it.VRFVersion,
		V2Plus:     it.V2Plus.Event,
	}
}

type SubscriptionCreated struct {
	VRFVersion vrfcommon.Version
	V2         *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreated
	V2Plus     *vrf_coordinator_v2plus.VRFCoordinatorV2PlusSubscriptionCreated
}

func (sc *SubscriptionCreated) Owner() common.Address {
	if sc.VRFVersion == vrfcommon.V2 {
		return sc.V2.Owner
	}
	return sc.V2Plus.Owner
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

type VRFProof struct {
	VRFVersion vrfcommon.Version
	V2         vrf_coordinator_v2.VRFProof
	V2Plus     vrf_coordinator_v2plus.VRFProof
}

func FromV2Proof(proof vrf_coordinator_v2.VRFProof) VRFProof {
	return VRFProof{
		VRFVersion: vrfcommon.V2,
		V2:         proof,
	}
}

func FromV2PlusProof(proof vrf_coordinator_v2plus.VRFProof) VRFProof {
	return VRFProof{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     proof,
	}
}

func ToV2Proofs(proofs []VRFProof) []vrf_coordinator_v2.VRFProof {
	v2Proofs := make([]vrf_coordinator_v2.VRFProof, len(proofs))
	for i, proof := range proofs {
		v2Proofs[i] = proof.V2
	}
	return v2Proofs
}

func ToV2PlusProofs(proofs []VRFProof) []vrf_coordinator_v2plus.VRFProof {
	v2Proofs := make([]vrf_coordinator_v2plus.VRFProof, len(proofs))
	for i, proof := range proofs {
		v2Proofs[i] = proof.V2Plus
	}
	return v2Proofs
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

func (c *coordinatorV2_X) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.RequestRandomWords(opts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
	}
	extraArgs, err := GetExtraArgsV1(payInEth)
	if err != nil {
		return nil, err
	}
	req := vrf_coordinator_v2plus.VRFV2PlusClientRandomWordsRequest{
		KeyHash:              keyHash,
		SubId:                subId,
		RequestConfirmations: requestConfirmations,
		CallbackGasLimit:     callbackGasLimit,
		NumWords:             numWords,
		ExtraArgs:            extraArgs,
	}
	return c.v2plus.RequestRandomWords(opts, req)
}

func GetExtraArgsV1(nativePayment bool) ([]byte, error) {
	encodedArgs, err := utils.ABIEncode(`[{"type":"bool"}]`, nativePayment)
	if err != nil {
		return nil, err
	}

	return append(extraArgsV1Tag, encodedArgs...), nil
}

func (c *coordinatorV2_X) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.CreateSubscription(opts)
	}
	return c.v2plus.CreateSubscription(opts)
}

func (c *coordinatorV2_X) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.AddConsumer(opts, subId, consumer)
	}
	return c.v2plus.AddConsumer(opts, subId, consumer)
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

func (c *coordinatorV2_X) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.RegisterProvingKey(opts, oracle, publicProvingKey)
	}
	return c.v2plus.RegisterProvingKey(opts, oracle, publicProvingKey)
}

func (c *coordinatorV2_X) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*SubscriptionCreatedIterator, error) {
	if c.v2 != nil {
		it, err := c.v2.FilterSubscriptionCreated(opts, subId)
		if err != nil {
			return nil, err
		}
		return &SubscriptionCreatedIterator{
			VRFVersion: vrfcommon.V2,
			V2:         it,
		}, nil
	}
	it, err := c.v2plus.FilterSubscriptionCreated(opts, subId)
	if err != nil {
		return nil, err
	}
	return &SubscriptionCreatedIterator{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     it,
	}, nil
}

func (c *coordinatorV2_X) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*RandomWordsRequestedIterator, error) {
	if c.v2 != nil {
		it, err := c.v2.FilterRandomWordsRequested(opts, keyHash, subId, sender)
		if err != nil {
			return nil, err
		}
		return &RandomWordsRequestedIterator{
			VRFVersion: vrfcommon.V2,
			V2:         it,
		}, nil
	}
	it, err := c.v2plus.FilterRandomWordsRequested(opts, keyHash, subId, sender)
	if err != nil {
		return nil, err
	}
	return &RandomWordsRequestedIterator{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     it,
	}, nil
}

func (c *coordinatorV2_X) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestID []*big.Int) (*RandomWordsFulfilledIterator, error) {
	if c.v2 != nil {
		it, err := c.v2.FilterRandomWordsFulfilled(opts, requestID)
		if err != nil {
			return nil, err
		}
		return &RandomWordsFulfilledIterator{
			VRFVersion: vrfcommon.V2,
			V2:         it,
		}, nil
	}
	it, err := c.v2plus.FilterRandomWordsFulfilled(opts, requestID)
	if err != nil {
		return nil, err
	}
	return &RandomWordsFulfilledIterator{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     it,
	}, nil
}

func (c *coordinatorV2_X) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.OracleWithdraw(opts, recipient, amount)
	}
	return c.v2plus.OracleWithdraw(opts, recipient, amount)
}

func (c *coordinatorV2_X) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.TransferOwnership(opts, to)
	}
	return c.v2plus.TransferOwnership(opts, to)
}

func (c *coordinatorV2_X) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.RemoveConsumer(opts, subId, consumer)
	}
	return c.v2plus.RemoveConsumer(opts, subId, consumer)
}

func (c *coordinatorV2_X) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	if c.v2 != nil {
		return c.v2.CancelSubscription(opts, subId, to)
	}
	return c.v2plus.CancelSubscription(opts, subId, to)
}

func (c *coordinatorV2_X) GetCommitment(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
	if c.v2 != nil {
		return c.v2.GetCommitment(opts, requestID)
	}
	return c.v2plus.SRequestCommitments(opts, requestID)
}

var (
	_ CoordinatorV2_X = (*coordinatorV2_X)(nil)
)
