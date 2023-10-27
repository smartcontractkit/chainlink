package vrftesthelpers

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2_plus_upgradeable_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2_upgradeable_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_malicious_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_malicious_consumer_v2_plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_reverting_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_reverting_example"
)

var (
	_ VRFConsumerContract = (*vrfConsumerContract)(nil)
)

// VRFConsumerContract is the common interface implemented by
// the example contracts used for the integration tests.
type VRFConsumerContract interface {
	CreateSubscriptionAndFund(opts *bind.TransactOpts, fundingJuels *big.Int) (*gethtypes.Transaction, error)
	CreateSubscriptionAndFundNative(opts *bind.TransactOpts, fundingAmount *big.Int) (*gethtypes.Transaction, error)
	SSubId(opts *bind.CallOpts) (*big.Int, error)
	SRequestId(opts *bind.CallOpts) (*big.Int, error)
	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subID *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*gethtypes.Transaction, error)
	SRandomWords(opts *bind.CallOpts, randomwordIdx *big.Int) (*big.Int, error)
	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*gethtypes.Transaction, error)
	TopUpSubscriptionNative(opts *bind.TransactOpts, amount *big.Int) (*gethtypes.Transaction, error)
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)
	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*gethtypes.Transaction, error)
	SetSubID(opts *bind.TransactOpts, subID *big.Int) (*gethtypes.Transaction, error)
}

type ConsumerType string

const (
	VRFConsumerV2           ConsumerType = "VRFConsumerV2"
	VRFV2PlusConsumer       ConsumerType = "VRFV2PlusConsumer"
	MaliciousConsumer       ConsumerType = "MaliciousConsumer"
	MaliciousConsumerPlus   ConsumerType = "MaliciousConsumerPlus"
	RevertingConsumer       ConsumerType = "RevertingConsumer"
	RevertingConsumerPlus   ConsumerType = "RevertingConsumerPlus"
	UpgradeableConsumer     ConsumerType = "UpgradeableConsumer"
	UpgradeableConsumerPlus ConsumerType = "UpgradeableConsumerPlus"
)

type vrfConsumerContract struct {
	consumerType            ConsumerType
	vrfConsumerV2           *vrf_consumer_v2.VRFConsumerV2
	vrfV2PlusConsumer       *vrfv2plus_consumer_example.VRFV2PlusConsumerExample
	maliciousConsumer       *vrf_malicious_consumer_v2.VRFMaliciousConsumerV2
	maliciousConsumerPlus   *vrf_malicious_consumer_v2_plus.VRFMaliciousConsumerV2Plus
	revertingConsumer       *vrfv2_reverting_example.VRFV2RevertingExample
	revertingConsumerPlus   *vrfv2plus_reverting_example.VRFV2PlusRevertingExample
	upgradeableConsumer     *vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExample
	upgradeableConsumerPlus *vrf_consumer_v2_plus_upgradeable_example.VRFConsumerV2PlusUpgradeableExample
}

func NewVRFConsumerV2(consumer *vrf_consumer_v2.VRFConsumerV2) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:  VRFConsumerV2,
		vrfConsumerV2: consumer,
	}
}

func NewVRFV2PlusConsumer(consumer *vrfv2plus_consumer_example.VRFV2PlusConsumerExample) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:      VRFV2PlusConsumer,
		vrfV2PlusConsumer: consumer,
	}
}

func NewMaliciousConsumer(consumer *vrf_malicious_consumer_v2.VRFMaliciousConsumerV2) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:      MaliciousConsumer,
		maliciousConsumer: consumer,
	}
}

func NewMaliciousConsumerPlus(consumer *vrf_malicious_consumer_v2_plus.VRFMaliciousConsumerV2Plus) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:          MaliciousConsumerPlus,
		maliciousConsumerPlus: consumer,
	}
}

func NewRevertingConsumer(consumer *vrfv2_reverting_example.VRFV2RevertingExample) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:      RevertingConsumer,
		revertingConsumer: consumer,
	}
}

func NewRevertingConsumerPlus(consumer *vrfv2plus_reverting_example.VRFV2PlusRevertingExample) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:          RevertingConsumerPlus,
		revertingConsumerPlus: consumer,
	}
}

func NewUpgradeableConsumer(consumer *vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExample) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:        UpgradeableConsumer,
		upgradeableConsumer: consumer,
	}
}

func NewUpgradeableConsumerPlus(consumer *vrf_consumer_v2_plus_upgradeable_example.VRFConsumerV2PlusUpgradeableExample) *vrfConsumerContract {
	return &vrfConsumerContract{
		consumerType:            UpgradeableConsumerPlus,
		upgradeableConsumerPlus: consumer,
	}
}

func (c *vrfConsumerContract) CreateSubscriptionAndFund(opts *bind.TransactOpts, fundingJuels *big.Int) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == UpgradeableConsumer {
		return c.upgradeableConsumer.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == UpgradeableConsumerPlus {
		return c.upgradeableConsumerPlus.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == MaliciousConsumer {
		return c.maliciousConsumer.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == MaliciousConsumerPlus {
		return c.maliciousConsumerPlus.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == RevertingConsumer {
		return c.revertingConsumer.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	if c.consumerType == RevertingConsumerPlus {
		return c.revertingConsumerPlus.CreateSubscriptionAndFund(opts, fundingJuels)
	}
	return nil, errors.New("CreateSubscriptionAndFund is not supported")
}

func (c *vrfConsumerContract) SSubId(opts *bind.CallOpts) (*big.Int, error) {
	if c.consumerType == VRFConsumerV2 {
		subID, err := c.vrfConsumerV2.SSubId(opts)
		if err != nil {
			return nil, err
		}
		return new(big.Int).SetUint64(subID), nil
	}
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.SSubId(opts)
	}
	if c.consumerType == UpgradeableConsumer {
		subID, err := c.upgradeableConsumer.SSubId(opts)
		if err != nil {
			return nil, err
		}
		return new(big.Int).SetUint64(subID), nil
	}
	if c.consumerType == UpgradeableConsumerPlus {
		return c.upgradeableConsumerPlus.SSubId(opts)
	}
	if c.consumerType == RevertingConsumer {
		subID, err := c.revertingConsumer.SSubId(opts)
		if err != nil {
			return nil, err
		}
		return new(big.Int).SetUint64(subID), nil
	}
	if c.consumerType == RevertingConsumerPlus {
		return c.revertingConsumerPlus.SSubId(opts)
	}
	return nil, errors.New("SSubId is not supported")
}

func (c *vrfConsumerContract) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.SRequestId(opts)
	}
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.SRecentRequestId(opts)
	}
	if c.consumerType == UpgradeableConsumer {
		return c.upgradeableConsumer.SRequestId(opts)
	}
	if c.consumerType == UpgradeableConsumerPlus {
		return c.upgradeableConsumerPlus.SRequestId(opts)
	}
	if c.consumerType == MaliciousConsumer {
		return c.maliciousConsumer.SRequestId(opts)
	}
	if c.consumerType == MaliciousConsumerPlus {
		return c.maliciousConsumerPlus.SRequestId(opts)
	}
	if c.consumerType == RevertingConsumer {
		return c.revertingConsumer.SRequestId(opts)
	}
	if c.consumerType == RevertingConsumerPlus {
		return c.revertingConsumerPlus.SRequestId(opts)
	}
	return nil, errors.New("SRequestId is not supported")
}

func (c *vrfConsumerContract) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subID *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.RequestRandomWords(opts, callbackGasLimit, minReqConfs, numWords, keyHash, payInEth)
	}
	if payInEth {
		return nil, errors.New("eth payment not supported")
	}
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.RequestRandomness(opts, keyHash, subID.Uint64(), minReqConfs, callbackGasLimit, numWords)
	}
	if c.consumerType == UpgradeableConsumer {
		return c.upgradeableConsumer.RequestRandomness(opts, keyHash, subID.Uint64(), minReqConfs, callbackGasLimit, numWords)
	}
	if c.consumerType == UpgradeableConsumerPlus {
		return c.upgradeableConsumerPlus.RequestRandomness(opts, keyHash, subID, minReqConfs, callbackGasLimit, numWords)
	}
	if c.consumerType == MaliciousConsumer {
		return c.maliciousConsumer.RequestRandomness(opts, keyHash)
	}
	if c.consumerType == MaliciousConsumerPlus {
		return c.maliciousConsumerPlus.RequestRandomness(opts, keyHash)
	}
	if c.consumerType == RevertingConsumer {
		return c.revertingConsumer.RequestRandomness(opts, keyHash, subID.Uint64(), minReqConfs, callbackGasLimit, numWords)
	}
	if c.consumerType == RevertingConsumerPlus {
		return c.revertingConsumerPlus.RequestRandomness(opts, keyHash, subID, minReqConfs, callbackGasLimit, numWords)
	}
	return nil, errors.New("RequestRandomness is not supported")
}

func (c *vrfConsumerContract) SRandomWords(opts *bind.CallOpts, randomwordIdx *big.Int) (*big.Int, error) {
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.SRandomWords(opts, randomwordIdx)
	}
	if c.consumerType == VRFV2PlusConsumer {
		requestID, err := c.vrfV2PlusConsumer.SRecentRequestId(opts)
		if err != nil {
			return nil, err
		}
		randomWord, err := c.vrfV2PlusConsumer.GetRandomness(opts, requestID, randomwordIdx)
		if err != nil {
			return nil, err
		}
		return randomWord, nil
	}
	if c.consumerType == UpgradeableConsumer {
		return c.upgradeableConsumer.SRandomWords(opts, randomwordIdx)
	}
	if c.consumerType == UpgradeableConsumerPlus {
		return c.upgradeableConsumerPlus.SRandomWords(opts, randomwordIdx)
	}
	return nil, errors.New("SRandomWords is not supported")
}

func (c *vrfConsumerContract) TopUpSubscription(opts *bind.TransactOpts, fundingJuels *big.Int) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.TopUpSubscription(opts, fundingJuels)
	}
	if c.consumerType == RevertingConsumer {
		return c.revertingConsumer.TopUpSubscription(opts, fundingJuels)
	}
	if c.consumerType == RevertingConsumerPlus {
		return c.revertingConsumerPlus.TopUpSubscription(opts, fundingJuels)
	}
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.TopUpSubscription(opts, fundingJuels)
	}
	return nil, errors.New("TopUpSubscription is not supported")
}

func (c *vrfConsumerContract) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.SGasAvailable(opts)
	}
	if c.consumerType == UpgradeableConsumer {
		return c.upgradeableConsumer.SGasAvailable(opts)
	}
	if c.consumerType == UpgradeableConsumerPlus {
		return c.upgradeableConsumerPlus.SGasAvailable(opts)
	}
	return nil, errors.New("SGasAvailable is not supported")
}

func (c *vrfConsumerContract) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFConsumerV2 {
		return c.vrfConsumerV2.UpdateSubscription(opts, consumers)
	}
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.UpdateSubscription(opts, consumers)
	}
	return nil, errors.New("UpdateSubscription is not supported")
}

func (c *vrfConsumerContract) SetSubID(opts *bind.TransactOpts, subID *big.Int) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFV2PlusConsumer {
		return c.vrfV2PlusConsumer.SetSubId(opts, subID)
	}
	return nil, errors.New("SetSubID is not supported")
}

func (c *vrfConsumerContract) CreateSubscriptionAndFundNative(opts *bind.TransactOpts, fundingAmount *big.Int) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFV2PlusConsumer {
		// copy object to not mutate original opts
		o := *opts
		o.Value = fundingAmount
		return c.vrfV2PlusConsumer.CreateSubscriptionAndFundNative(&o)
	}
	return nil, errors.New("CreateSubscriptionAndFundNative is not supported")
}

func (c *vrfConsumerContract) TopUpSubscriptionNative(opts *bind.TransactOpts, amount *big.Int) (*gethtypes.Transaction, error) {
	if c.consumerType == VRFV2PlusConsumer {
		// copy object to not mutate original opts
		o := *opts
		o.Value = amount
		return c.vrfV2PlusConsumer.TopUpSubscriptionNative(&o)
	}
	return nil, errors.New("TopUpSubscriptionNative is not supported")
}
