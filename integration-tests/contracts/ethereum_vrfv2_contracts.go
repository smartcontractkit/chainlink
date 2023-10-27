package contracts

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2_consumer_wrapper"
)

// EthereumVRFCoordinatorV2 represents VRFV2 coordinator contract
type EthereumVRFCoordinatorV2 struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *vrf_coordinator_v2.VRFCoordinatorV2
}

// EthereumVRFConsumerV2 represents VRFv2 consumer contract
type EthereumVRFConsumerV2 struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_consumer_v2.VRFConsumerV2
}

// EthereumVRFv2Consumer represents VRFv2 consumer contract
type EthereumVRFv2Consumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_v2_consumer_wrapper.VRFv2Consumer
}

// EthereumVRFv2LoadTestConsumer represents VRFv2 consumer contract for performing Load Tests
type EthereumVRFv2LoadTestConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics
}

// DeployVRFCoordinatorV2 deploys VRFV2 coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator_v2.DeployVRFCoordinatorV2(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr), common.HexToAddress(linkEthFeedAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2{
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_v2.VRFCoordinatorV2),
		address:     address,
	}, err
}

// DeployVRFConsumerV2 deploys VRFv@ consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumerV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_consumer_v2.DeployVRFConsumerV2(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumerV2{
		client:   e.client,
		consumer: instance.(*vrf_consumer_v2.VRFConsumerV2),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFv2Consumer(coordinatorAddr string) (VRFv2Consumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFv2Consumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_v2_consumer_wrapper.DeployVRFv2Consumer(auth, backend, common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2Consumer{
		client:   e.client,
		consumer: instance.(*vrf_v2_consumer_wrapper.VRFv2Consumer),
		address:  address,
	}, err
}

// DeployVRFv2LoadTestConsumer(coordinatorAddr string) (VRFv2Consumer, error)
func (e *EthereumContractDeployer) DeployVRFv2LoadTestConsumer(coordinatorAddr string) (VRFv2LoadTestConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2LoadTestWithMetrics", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_load_test_with_metrics.DeployVRFV2LoadTestWithMetrics(auth, backend, common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2LoadTestConsumer{
		client:   e.client,
		consumer: instance.(*vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics),
		address:  address,
	}, err
}

func (v *EthereumVRFCoordinatorV2) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	hash, err := v.coordinator.HashOfKey(opts, pubKey)
	if err != nil {
		return [32]byte{}, err
	}
	return hash, nil
}

func (v *EthereumVRFCoordinatorV2) GetSubscription(ctx context.Context, subID uint64) (vrf_coordinator_v2.GetSubscription, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	subscription, err := v.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return vrf_coordinator_v2.GetSubscription{}, err
	}
	return subscription, nil
}

func (v *EthereumVRFCoordinatorV2) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetConfig(
		opts,
		minimumRequestConfirmations,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		fallbackWeiPerUnitLink,
		feeConfig,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2) RegisterProvingKey(
	oracleAddr string,
	publicProvingKey [2]*big.Int,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, common.HexToAddress(oracleAddr), publicProvingKey)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2) CreateSubscription() error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.CreateSubscription(opts)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2) AddConsumer(subId uint64, consumerAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.AddConsumer(
		opts,
		subId,
		common.HexToAddress(consumerAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// GetAllRandomWords get all VRFv2 randomness output words
func (v *EthereumVRFConsumerV2) GetAllRandomWords(ctx context.Context, num int) ([]*big.Int, error) {
	words := make([]*big.Int, 0)
	for i := 0; i < num; i++ {
		word, err := v.consumer.SRandomWords(&bind.CallOpts{
			From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
			Context: ctx,
		}, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

// LoadExistingConsumer loads an EthereumVRFConsumerV2 with a specified address
func (v *EthereumVRFConsumerV2) LoadExistingConsumer(address string, client blockchain.EVMClient) error {
	a := common.HexToAddress(address)
	consumer, err := vrf_consumer_v2.NewVRFConsumerV2(a, client.(*blockchain.EthereumClient).Client)
	if err != nil {
		return err
	}
	v.client = client
	v.consumer = consumer
	v.address = &a
	return nil
}

// CreateFundedSubscription create funded subscription for VRFv2 randomness
func (v *EthereumVRFConsumerV2) CreateFundedSubscription(funds *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.CreateSubscriptionAndFund(opts, funds)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// TopUpSubscriptionFunds add funds to a VRFv2 subscription
func (v *EthereumVRFConsumerV2) TopUpSubscriptionFunds(funds *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.TopUpSubscription(opts, funds)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFConsumerV2) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFv2Consumer) Address() string {
	return v.address.Hex()
}

// CurrentSubscription get current VRFv2 subscription
func (v *EthereumVRFConsumerV2) CurrentSubscription() (uint64, error) {
	return v.consumer.SSubId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
}

// GasAvailable get available gas after randomness fulfilled
func (v *EthereumVRFConsumerV2) GasAvailable() (*big.Int, error) {
	return v.consumer.SGasAvailable(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
}

func (v *EthereumVRFConsumerV2) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(ethereum.CallMsg{
		To: v.address,
	})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFConsumerV2) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.RequestRandomness(opts, hash, subID, confs, gasLimit, numWords)
	if err != nil {
		return err
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confs).
		Interface("Callback Gas Limit", gasLimit).
		Interface("KeyHash", hex.EncodeToString(hash[:])).
		Interface("Consumer Contract", v.address).
		Msg("RequestRandomness called")
	return v.client.ProcessTransaction(tx)
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFv2Consumer) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.RequestRandomWords(opts, subID, gasLimit, confs, numWords, hash)
	if err != nil {
		return err
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confs).
		Interface("Callback Gas Limit", gasLimit).
		Interface("KeyHash", hex.EncodeToString(hash[:])).
		Interface("Consumer Contract", v.address).
		Msg("RequestRandomness called")
	return v.client.ProcessTransaction(tx)
}

// RandomnessOutput get VRFv2 randomness output (word)
func (v *EthereumVRFConsumerV2) RandomnessOutput(ctx context.Context, arg0 *big.Int) (*big.Int, error) {
	return v.consumer.SRandomWords(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, arg0)
}

func (v *EthereumVRFv2LoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomness(keyHash [32]byte, subID uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, requestCount uint16) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := v.consumer.RequestRandomWords(opts, subID, requestConfirmations, keyHash, callbackGasLimit, numWords, requestCount)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFv2Consumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2_consumer_wrapper.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2Consumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.LastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFv2LoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2LoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFv2LoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	requestCount, err := v.consumer.SRequestCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	fulfilmentCount, err := v.consumer.SResponseCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	averageFulfillmentInMillions, err := v.consumer.SAverageFulfillmentInMillions(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	slowestFulfillment, err := v.consumer.SSlowestFulfillment(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	fastestFulfillment, err := v.consumer.SFastestFulfillment(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}

	return &VRFLoadTestMetrics{
		requestCount,
		fulfilmentCount,
		averageFulfillmentInMillions,
		slowestFulfillment,
		fastestFulfillment,
	}, nil
}
