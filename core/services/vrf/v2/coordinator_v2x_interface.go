package v2

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/log"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/extraargs"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

var (
	_ CoordinatorV2_X = (*coordinatorV2_5)(nil)
)

// CoordinatorV2_X is an interface that allows us to use the same code for
// both the V2 and V2Plus coordinators.
type CoordinatorV2_X interface { //nolint:revive // V2_X naming matches coordinator version convention
	Address() common.Address
	ParseRandomWordsRequested(log types.Log) (RandomWordsRequested, error)
	ParseRandomWordsFulfilled(log types.Log) (RandomWordsFulfilled, error)
	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subID *big.Int, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error)
	AddConsumer(opts *bind.TransactOpts, subID *big.Int, consumer common.Address) (*types.Transaction, error)
	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)
	GetSubscription(opts *bind.CallOpts, subID *big.Int) (Subscription, error)
	GetConfig(opts *bind.CallOpts) (Config, error)
	ParseLog(log types.Log) (generated.AbigenLog, error)
	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)
	Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)
	WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)
	LogsWithTopics(keyHash common.Hash) map[common.Hash][][]log.Topic
	Version() vrfcommon.Version
	RegisterProvingKey(opts *bind.TransactOpts, oracle *common.Address, publicProvingKey [2]*big.Int, maxGasPrice *uint64) (*types.Transaction, error)
	FilterSubscriptionCreated(opts *bind.FilterOpts, subID []*big.Int) (SubscriptionCreatedIterator, error)
	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subID []*big.Int, sender []common.Address) (RandomWordsRequestedIterator, error)
	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestID []*big.Int, subID []*big.Int) (RandomWordsFulfilledIterator, error)
	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
	RemoveConsumer(opts *bind.TransactOpts, subID *big.Int, consumer common.Address) (*types.Transaction, error)
	CancelSubscription(opts *bind.TransactOpts, subID *big.Int, to common.Address) (*types.Transaction, error)
	GetCommitment(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error)
	Migrate(opts *bind.TransactOpts, subID *big.Int, newCoordinator common.Address) (*types.Transaction, error)
	FundSubscriptionWithNative(opts *bind.TransactOpts, subID *big.Int, amount *big.Int) (*types.Transaction, error)
	// RandomWordsRequestedTopic returns the log topic of the RandomWordsRequested log
	RandomWordsRequestedTopic() common.Hash
	// RandomWordsFulfilledTopic returns the log topic of the RandomWordsFulfilled log
	RandomWordsFulfilledTopic() common.Hash
}

type coordinatorV2_5 struct {
	vrfVersion  vrfcommon.Version
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25Interface
}

func NewCoordinatorV2_5(c vrf_coordinator_v2_5.VRFCoordinatorV25Interface) CoordinatorV2_X {
	return &coordinatorV2_5{
		vrfVersion:  vrfcommon.V2Plus,
		coordinator: c,
	}
}

func (c *coordinatorV2_5) RandomWordsRequestedTopic() common.Hash {
	return vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRandomWordsRequested{}.Topic()
}

func (c *coordinatorV2_5) RandomWordsFulfilledTopic() common.Hash {
	return vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRandomWordsFulfilled{}.Topic()
}

func (c *coordinatorV2_5) Address() common.Address {
	return c.coordinator.Address()
}

func (c *coordinatorV2_5) ParseRandomWordsRequested(log types.Log) (RandomWordsRequested, error) {
	parsed, err := c.coordinator.ParseRandomWordsRequested(log)
	if err != nil {
		return nil, err
	}
	return NewV2_5RandomWordsRequested(parsed), nil
}

func (c *coordinatorV2_5) ParseRandomWordsFulfilled(log types.Log) (RandomWordsFulfilled, error) {
	parsed, err := c.coordinator.ParseRandomWordsFulfilled(log)
	if err != nil {
		return nil, err
	}
	return NewV2_5RandomWordsFulfilled(parsed), nil
}

func (c *coordinatorV2_5) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subID *big.Int, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error) {
	extraArgs, err := extraargs.EncodeV1(payInEth)
	if err != nil {
		return nil, err
	}
	req := vrf_coordinator_v2_5.VRFV2PlusClientRandomWordsRequest{
		KeyHash:              keyHash,
		SubId:                subID,
		RequestConfirmations: requestConfirmations,
		CallbackGasLimit:     callbackGasLimit,
		NumWords:             numWords,
		ExtraArgs:            extraArgs,
	}
	return c.coordinator.RequestRandomWords(opts, req)
}

func (c *coordinatorV2_5) AddConsumer(opts *bind.TransactOpts, subID *big.Int, consumer common.Address) (*types.Transaction, error) {
	return c.coordinator.AddConsumer(opts, subID, consumer)
}

func (c *coordinatorV2_5) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return c.coordinator.CreateSubscription(opts)
}

func (c *coordinatorV2_5) GetSubscription(opts *bind.CallOpts, subID *big.Int) (Subscription, error) {
	sub, err := c.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return nil, err
	}
	return NewV2_5Subscription(sub), nil
}

func (c *coordinatorV2_5) GetConfig(opts *bind.CallOpts) (Config, error) {
	config, err := c.coordinator.SConfig(opts)
	if err != nil {
		return nil, err
	}
	return NewV2_5Config(config), nil
}

func (c *coordinatorV2_5) ParseLog(log types.Log) (generated.AbigenLog, error) {
	return c.coordinator.ParseLog(log)
}

func (c *coordinatorV2_5) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return nil, errors.New("oracle withdraw not implemented for v2.5")
}

func (c *coordinatorV2_5) Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return c.coordinator.Withdraw(opts, recipient)
}

func (c *coordinatorV2_5) WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return c.coordinator.WithdrawNative(opts, recipient)
}

func (c *coordinatorV2_5) LogsWithTopics(keyHash common.Hash) map[common.Hash][][]log.Topic {
	return map[common.Hash][][]log.Topic{
		vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested{}.Topic(): {
			{
				log.Topic(keyHash),
			},
		},
	}
}

func (c *coordinatorV2_5) Version() vrfcommon.Version {
	return c.vrfVersion
}

func (c *coordinatorV2_5) RegisterProvingKey(opts *bind.TransactOpts, oracle *common.Address, publicProvingKey [2]*big.Int, maxGasPrice *uint64) (*types.Transaction, error) {
	if oracle != nil {
		return nil, errors.New("oracle address not supported for registering proving key in v2.5")
	}
	if maxGasPrice == nil {
		return nil, errors.New("max gas price is required for registering proving key in v2.5")
	}
	return c.coordinator.RegisterProvingKey(opts, publicProvingKey, *maxGasPrice)
}

func (c *coordinatorV2_5) FilterSubscriptionCreated(opts *bind.FilterOpts, subID []*big.Int) (SubscriptionCreatedIterator, error) {
	it, err := c.coordinator.FilterSubscriptionCreated(opts, subID)
	if err != nil {
		return nil, err
	}
	return NewV2_5SubscriptionCreatedIterator(it), nil
}

func (c *coordinatorV2_5) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subID []*big.Int, sender []common.Address) (RandomWordsRequestedIterator, error) {
	it, err := c.coordinator.FilterRandomWordsRequested(opts, keyHash, subID, sender)
	if err != nil {
		return nil, err
	}
	return NewV2_5RandomWordsRequestedIterator(it), nil
}

func (c *coordinatorV2_5) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestID []*big.Int, subID []*big.Int) (RandomWordsFulfilledIterator, error) {
	it, err := c.coordinator.FilterRandomWordsFulfilled(opts, requestID, subID)
	if err != nil {
		return nil, err
	}
	return NewV2_5RandomWordsFulfilledIterator(it), nil
}

func (c *coordinatorV2_5) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return c.coordinator.TransferOwnership(opts, to)
}

func (c *coordinatorV2_5) RemoveConsumer(opts *bind.TransactOpts, subID *big.Int, consumer common.Address) (*types.Transaction, error) {
	return c.coordinator.RemoveConsumer(opts, subID, consumer)
}

func (c *coordinatorV2_5) CancelSubscription(opts *bind.TransactOpts, subID *big.Int, to common.Address) (*types.Transaction, error) {
	return c.coordinator.CancelSubscription(opts, subID, to)
}

func (c *coordinatorV2_5) GetCommitment(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
	return c.coordinator.SRequestCommitments(opts, requestID)
}

func (c *coordinatorV2_5) Migrate(opts *bind.TransactOpts, subID *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return c.coordinator.Migrate(opts, subID, newCoordinator)
}

func (c *coordinatorV2_5) FundSubscriptionWithNative(opts *bind.TransactOpts, subID *big.Int, amount *big.Int) (*types.Transaction, error) {
	if opts == nil {
		return nil, errors.New("*bind.TransactOpts cannot be nil")
	}
	o := *opts
	o.Value = amount
	return c.coordinator.FundSubscriptionWithNative(&o, subID)
}

var (
	_ RandomWordsRequestedIterator = (*v2_5RandomWordsRequestedIterator)(nil)
)

type RandomWordsRequestedIterator interface {
	Next() bool
	Error() error
	Close() error
	Event() RandomWordsRequested
}

type v2_5RandomWordsRequestedIterator struct {
	vrfVersion vrfcommon.Version
	iterator   *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequestedIterator
}

func NewV2_5RandomWordsRequestedIterator(it *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequestedIterator) RandomWordsRequestedIterator {
	return &v2_5RandomWordsRequestedIterator{
		vrfVersion: vrfcommon.V2Plus,
		iterator:   it,
	}
}

func (it *v2_5RandomWordsRequestedIterator) Next() bool {
	return it.iterator.Next()
}

func (it *v2_5RandomWordsRequestedIterator) Error() error {
	return it.iterator.Error()
}

func (it *v2_5RandomWordsRequestedIterator) Close() error {
	return it.iterator.Close()
}

func (it *v2_5RandomWordsRequestedIterator) Event() RandomWordsRequested {
	return NewV2_5RandomWordsRequested(it.iterator.Event)
}

var (
	_ RandomWordsRequested = (*v2_5RandomWordsRequested)(nil)
)

type RandomWordsRequested interface {
	Raw() types.Log
	NumWords() uint32
	SubID() *big.Int
	MinimumRequestConfirmations() uint16
	KeyHash() [32]byte
	RequestID() *big.Int
	PreSeed() *big.Int
	Sender() common.Address
	CallbackGasLimit() uint32
	NativePayment() bool
}

type v2_5RandomWordsRequested struct {
	vrfVersion vrfcommon.Version
	event      *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested
}

func NewV2_5RandomWordsRequested(event *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested) RandomWordsRequested {
	return &v2_5RandomWordsRequested{
		vrfVersion: vrfcommon.V2Plus,
		event:      event,
	}
}

func (r *v2_5RandomWordsRequested) Raw() types.Log {
	return r.event.Raw
}

func (r *v2_5RandomWordsRequested) NumWords() uint32 {
	return r.event.NumWords
}

func (r *v2_5RandomWordsRequested) SubID() *big.Int {
	return r.event.SubId
}

func (r *v2_5RandomWordsRequested) MinimumRequestConfirmations() uint16 {
	return r.event.MinimumRequestConfirmations
}

func (r *v2_5RandomWordsRequested) KeyHash() [32]byte {
	return r.event.KeyHash
}

func (r *v2_5RandomWordsRequested) RequestID() *big.Int {
	return r.event.RequestId
}

func (r *v2_5RandomWordsRequested) PreSeed() *big.Int {
	return r.event.PreSeed
}

func (r *v2_5RandomWordsRequested) Sender() common.Address {
	return r.event.Sender
}

func (r *v2_5RandomWordsRequested) CallbackGasLimit() uint32 {
	return r.event.CallbackGasLimit
}

func (r *v2_5RandomWordsRequested) NativePayment() bool {
	nativePayment, err := extraargs.DecodeV1(r.event.ExtraArgs)
	if err != nil {
		panic(err)
	}
	return nativePayment
}

var (
	_ RandomWordsFulfilledIterator = (*v2_5RandomWordsFulfilledIterator)(nil)
)

type RandomWordsFulfilledIterator interface {
	Next() bool
	Error() error
	Close() error
	Event() RandomWordsFulfilled
}

type v2_5RandomWordsFulfilledIterator struct {
	vrfVersion vrfcommon.Version
	iterator   *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilledIterator
}

func NewV2_5RandomWordsFulfilledIterator(it *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilledIterator) RandomWordsFulfilledIterator {
	return &v2_5RandomWordsFulfilledIterator{
		vrfVersion: vrfcommon.V2Plus,
		iterator:   it,
	}
}

func (it *v2_5RandomWordsFulfilledIterator) Next() bool {
	return it.iterator.Next()
}

func (it *v2_5RandomWordsFulfilledIterator) Error() error {
	return it.iterator.Error()
}

func (it *v2_5RandomWordsFulfilledIterator) Close() error {
	return it.iterator.Close()
}

func (it *v2_5RandomWordsFulfilledIterator) Event() RandomWordsFulfilled {
	return NewV2_5RandomWordsFulfilled(it.iterator.Event)
}

var (
	_ RandomWordsFulfilled = (*v2_5RandomWordsFulfilled)(nil)
)

type RandomWordsFulfilled interface {
	RequestID() *big.Int
	Success() bool
	SubID() *big.Int
	Payment() *big.Int
	Raw() types.Log
	NativePayment() bool
}

type v2_5RandomWordsFulfilled struct {
	vrfVersion vrfcommon.Version
	event      *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled
}

func NewV2_5RandomWordsFulfilled(event *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled) RandomWordsFulfilled {
	return &v2_5RandomWordsFulfilled{
		vrfVersion: vrfcommon.V2Plus,
		event:      event,
	}
}

func (rwf *v2_5RandomWordsFulfilled) RequestID() *big.Int {
	return rwf.event.RequestId
}

func (rwf *v2_5RandomWordsFulfilled) Success() bool {
	return rwf.event.Success
}

func (rwf *v2_5RandomWordsFulfilled) SubID() *big.Int {
	return rwf.event.SubId
}

func (rwf *v2_5RandomWordsFulfilled) Payment() *big.Int {
	return rwf.event.Payment
}

func (rwf *v2_5RandomWordsFulfilled) Raw() types.Log {
	return rwf.event.Raw
}

func (rwf *v2_5RandomWordsFulfilled) NativePayment() bool {
	return rwf.event.NativePayment
}

var (
	_ SubscriptionCreatedIterator = (*v2_5SubscriptionCreatedIterator)(nil)
)

type SubscriptionCreatedIterator interface {
	Next() bool
	Error() error
	Close() error
	Event() SubscriptionCreated
}

type v2_5SubscriptionCreatedIterator struct {
	vrfVersion vrfcommon.Version
	iterator   *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreatedIterator
}

func NewV2_5SubscriptionCreatedIterator(it *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreatedIterator) SubscriptionCreatedIterator {
	return &v2_5SubscriptionCreatedIterator{
		vrfVersion: vrfcommon.V2Plus,
		iterator:   it,
	}
}

func (it *v2_5SubscriptionCreatedIterator) Next() bool {
	return it.iterator.Next()
}

func (it *v2_5SubscriptionCreatedIterator) Error() error {
	return it.iterator.Error()
}

func (it *v2_5SubscriptionCreatedIterator) Close() error {
	return it.iterator.Close()
}

func (it *v2_5SubscriptionCreatedIterator) Event() SubscriptionCreated {
	return NewV2_5SubscriptionCreated(it.iterator.Event)
}

var (
	_ SubscriptionCreated = (*v2_5SubscriptionCreated)(nil)
)

type SubscriptionCreated interface {
	Owner() common.Address
	SubID() *big.Int
}

type v2_5SubscriptionCreated struct {
	vrfVersion vrfcommon.Version
	event      *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreated
}

func NewV2_5SubscriptionCreated(event *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreated) SubscriptionCreated {
	return &v2_5SubscriptionCreated{
		vrfVersion: vrfcommon.V2Plus,
		event:      event,
	}
}

func (sc *v2_5SubscriptionCreated) Owner() common.Address {
	return sc.event.Owner
}

func (sc *v2_5SubscriptionCreated) SubID() *big.Int {
	return sc.event.SubId
}

var (
	_ Subscription = (*v2_5Subscription)(nil)
)

type Subscription interface {
	Balance() *big.Int
	NativeBalance() *big.Int
	Owner() common.Address
	Consumers() []common.Address
	Version() vrfcommon.Version
}

type v2_5Subscription struct {
	vrfVersion vrfcommon.Version
	event      vrf_coordinator_v2_5.GetSubscription
}

func NewV2_5Subscription(event vrf_coordinator_v2_5.GetSubscription) Subscription {
	return &v2_5Subscription{
		vrfVersion: vrfcommon.V2Plus,
		event:      event,
	}
}

func (s *v2_5Subscription) Balance() *big.Int {
	return s.event.Balance
}

func (s *v2_5Subscription) NativeBalance() *big.Int {
	return s.event.NativeBalance
}

func (s *v2_5Subscription) Owner() common.Address {
	return s.event.SubOwner
}

func (s *v2_5Subscription) Consumers() []common.Address {
	return s.event.Consumers
}

func (s *v2_5Subscription) Version() vrfcommon.Version {
	return s.vrfVersion
}

var (
	_ Config = (*v2_5Config)(nil)
)

type Config interface {
	MinimumRequestConfirmations() uint16
	MaxGasLimit() uint32
	GasAfterPaymentCalculation() uint32
	StalenessSeconds() uint32
}

type v2_5Config struct {
	vrfVersion vrfcommon.Version
	config     vrf_coordinator_v2_5.SConfig
}

func NewV2_5Config(config vrf_coordinator_v2_5.SConfig) Config {
	return &v2_5Config{
		vrfVersion: vrfcommon.V2Plus,
		config:     config,
	}
}

func (c *v2_5Config) MinimumRequestConfirmations() uint16 {
	return c.config.MinimumRequestConfirmations
}

func (c *v2_5Config) MaxGasLimit() uint32 {
	return c.config.MaxGasLimit
}

func (c *v2_5Config) GasAfterPaymentCalculation() uint32 {
	return c.config.GasAfterPaymentCalculation
}

func (c *v2_5Config) StalenessSeconds() uint32 {
	return c.config.StalenessSeconds
}

type VRFProof struct {
	VRFVersion vrfcommon.Version
	V2Plus     vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof
}

func FromV2PlusProof(proof vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof) VRFProof {
	return VRFProof{
		VRFVersion: vrfcommon.V2Plus,
		V2Plus:     proof,
	}
}

func ToV2PlusProofs(proofs []VRFProof) []vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof {
	v2Proofs := make([]vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof, len(proofs))
	for i, proof := range proofs {
		v2Proofs[i] = proof.V2Plus
	}
	return v2Proofs
}

type RequestCommitment struct {
	VRFVersion vrfcommon.Version
	V2Plus     vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment
}

func ToV2PlusCommitments(commitments []RequestCommitment) []vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment {
	v2PlusCommitments := make([]vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment, len(commitments))
	for i, commitment := range commitments {
		v2PlusCommitments[i] = commitment.V2Plus
	}
	return v2PlusCommitments
}

func NewRequestCommitment(val any) RequestCommitment {
	switch val := val.(type) {
	case vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment:
		return RequestCommitment{VRFVersion: vrfcommon.V2Plus, V2Plus: val}
	default:
		panic(fmt.Sprintf("NewRequestCommitment: unknown type %T", val))
	}
}

func (r *RequestCommitment) Get() any {
	return r.V2Plus
}

func (r *RequestCommitment) NativePayment() bool {
	nativePayment, err := extraargs.DecodeV1(r.V2Plus.ExtraArgs)
	if err != nil {
		panic(err)
	}
	return nativePayment
}

func (r *RequestCommitment) NumWords() uint32 {
	return r.V2Plus.NumWords
}

func (r *RequestCommitment) Sender() common.Address {
	return r.V2Plus.Sender
}

func (r *RequestCommitment) BlockNum() uint64 {
	return r.V2Plus.BlockNum
}

func (r *RequestCommitment) SubID() *big.Int {
	return r.V2Plus.SubId
}

func (r *RequestCommitment) CallbackGasLimit() uint32 {
	return r.V2Plus.CallbackGasLimit
}
