package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/batch_blockhash_store"
)

// DeployVRFContract deploy VRF contract
func (e *EthereumContractDeployer) DeployVRFContract() (VRF, error) {
	address, _, instance, err := e.client.DeployContract("VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:  e.client,
		vrf:     instance.(*ethereum.VRF),
		address: address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func (e *EthereumContractDeployer) DeployBlockhashStore() (BlockHashStore, error) {
	address, _, instance, err := e.client.DeployContract("BlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployBlockhashStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBlockhashStore{
		client:         e.client,
		blockHashStore: instance.(*ethereum.BlockhashStore),
		address:        address,
	}, err
}

// DeployVRFCoordinatorV2 deploys VRFV2 coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFCoordinatorV2(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr), common.HexToAddress(linkEthFeedAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2{
		client:      e.client,
		coordinator: instance.(*ethereum.VRFCoordinatorV2),
		address:     address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFCoordinator(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinator{
		client:      e.client,
		coordinator: instance.(*ethereum.VRFCoordinator),
		address:     address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFConsumer(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumer{
		client:   e.client,
		consumer: instance.(*ethereum.VRFConsumer),
		address:  address,
	}, err
}

// DeployVRFConsumerV2 deploys VRFv@ consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumerV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFConsumerV2(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumerV2{
		client:   e.client,
		consumer: instance.(*ethereum.VRFConsumerV2),
		address:  address,
	}, err
}

// DeployDKG deploys DKG contract
func (e *EthereumContractDeployer) DeployDKG() (DKG, error) {
	address, _, instance, err := e.client.DeployContract("DKG", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return dkg.DeployDKG(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumDKG{
		client:  e.client,
		dkg:     instance.(*dkg.DKG),
		address: address,
	}, err
}

// DeployOCR2VRFCoordinator deploys CR2VRFCoordinator contract
func (e *EthereumContractDeployer) DeployOCR2VRFCoordinator(beaconPeriodBlocksCount *big.Int, linkAddress string, linkEthFeedAddress string) (VRFCoordinatorV3, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV3", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator.DeployVRFCoordinator(auth, backend, beaconPeriodBlocksCount, common.HexToAddress(linkAddress), common.HexToAddress(linkEthFeedAddress))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV3{
		client:           e.client,
		vrfCoordinatorV3: instance.(*vrf_coordinator.VRFCoordinator),
		address:          address,
	}, err
}

// DeployVRFBeacon deploys DeployVRFBeacon contract
func (e *EthereumContractDeployer) DeployVRFBeacon(vrfCoordinatorAddress string, linkAddress string, dkgAddress string, keyId string) (VRFBeacon, error) {
	keyIDBytes, err := decodeHexTo32ByteArray(keyId)
	if err != nil {
		return nil, err
	}
	address, _, instance, err := e.client.DeployContract("VRFBeacon", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_beacon.DeployVRFBeacon(auth, backend, common.HexToAddress(linkAddress), common.HexToAddress(vrfCoordinatorAddress), common.HexToAddress(dkgAddress), keyIDBytes)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFBeacon{
		client:    e.client,
		vrfBeacon: instance.(*vrf_beacon.VRFBeacon),
		address:   address,
	}, err
}

// DeployBatchBlockhashStore deploys DeployBatchBlockhashStore contract
func (e *EthereumContractDeployer) DeployBatchBlockhashStore(blockhashStoreAddr string) (BatchBlockhashStore, error) {
	address, _, instance, err := e.client.DeployContract("BatchBlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return batch_blockhash_store.DeployBatchBlockhashStore(auth, backend, common.HexToAddress(blockhashStoreAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBatchBlockhashStore{
		client:              e.client,
		batchBlockhashStore: instance.(*batch_blockhash_store.BatchBlockhashStore),
		address:             address,
	}, err
}

// todo - solve import cycle
func decodeHexTo32ByteArray(val string) ([32]byte, error) {
	var byteArray [32]byte
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return [32]byte{}, err
	}
	if len(decoded) != 32 {
		return [32]byte{}, fmt.Errorf("expected value to be 32 bytes but received %d bytes", len(decoded))
	}
	copy(byteArray[:], decoded)
	return byteArray, err
}

// DeployVRFBeaconConsumer deploys VRFv@ consumer contract
func (e *EthereumContractDeployer) DeployVRFBeaconConsumer(vrfCoordinatorV3Address string, beaconPeriodBlockCount *big.Int) (VRFBeaconConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFBeaconConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_beacon_consumer.DeployBeaconVRFConsumer(auth, backend, common.HexToAddress(vrfCoordinatorV3Address), false, beaconPeriodBlockCount)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFBeaconConsumer{
		client:            e.client,
		vrfBeaconConsumer: instance.(*vrf_beacon_consumer.BeaconVRFConsumer),
		address:           address,
	}, err
}

// EthereumBlockhashStore represents a blockhash store for VRF contract
type EthereumBlockhashStore struct {
	address        *common.Address
	client         blockchain.EVMClient
	blockHashStore *ethereum.BlockhashStore
}

func (v *EthereumBlockhashStore) Address() string {
	return v.address.Hex()
}

// EthereumVRFCoordinatorV2 represents VRFV2 coordinator contract
type EthereumVRFCoordinatorV2 struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *ethereum.VRFCoordinatorV2
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

func (v *EthereumVRFCoordinatorV2) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig ethereum.VRFCoordinatorV2FeeConfig) error {
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

// EthereumVRFCoordinator represents VRF coordinator contract
type EthereumVRFCoordinator struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *ethereum.VRFCoordinator
}

func (v *EthereumVRFCoordinator) Address() string {
	return v.address.Hex()
}

// HashOfKey get a hash of proving key to use it as a request ID part for VRF
func (v *EthereumVRFCoordinator) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
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

// RegisterProvingKey register VRF proving key
func (v *EthereumVRFCoordinator) RegisterProvingKey(
	fee *big.Int,
	oracleAddr string,
	publicProvingKey [2]*big.Int,
	jobID [32]byte,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, fee, common.HexToAddress(oracleAddr), publicProvingKey, jobID)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// EthereumVRFConsumerV2 represents VRFv2 consumer contract
type EthereumVRFConsumerV2 struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *ethereum.VRFConsumerV2
}

// CurrentSubscription get current VRFv2 subscription
func (v *EthereumVRFConsumerV2) CurrentSubscription() (uint64, error) {
	return v.consumer.SSubId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
}

// CreateFundedSubscription create funded subscription for VRFv2 randomness
func (v *EthereumVRFConsumerV2) CreateFundedSubscription(funds *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.TestCreateSubscriptionAndFund(opts, funds)
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

// GasAvailable get available gas after randomness fulfilled
func (v *EthereumVRFConsumerV2) GasAvailable() (*big.Int, error) {
	return v.consumer.SGasAvailable(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
}

func (v *EthereumVRFConsumerV2) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFConsumerV2) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.TestRequestRandomness(opts, hash, subID, confs, gasLimit, numWords)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// RandomnessOutput get VRFv2 randomness output (word)
func (v *EthereumVRFConsumerV2) RandomnessOutput(ctx context.Context, arg0 *big.Int) (*big.Int, error) {
	return v.consumer.SRandomWords(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, arg0)
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
	consumer, err := ethereum.NewVRFConsumerV2(a, client.(*blockchain.EthereumClient).Client)
	if err != nil {
		return err
	}
	v.client = client
	v.consumer = consumer
	v.address = &a
	return nil
}

// EthereumVRFConsumer represents VRF consumer contract
type EthereumVRFConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *ethereum.VRFConsumer
}

func (v *EthereumVRFConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFConsumer) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

// RequestRandomness requests VRF randomness
func (v *EthereumVRFConsumer) RequestRandomness(hash [32]byte, fee *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.TestRequestRandomness(opts, hash, fee)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// CurrentRoundID helper roundID counter in consumer to check when all randomness requests are finished
func (v *EthereumVRFConsumer) CurrentRoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return v.consumer.CurrentRoundID(opts)
}

func (v *EthereumVRFConsumer) WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error {
	ethEventChan := make(chan *ethereum.VRFConsumerPerfMetricsEvent)
	sub, err := v.consumer.WatchPerfMetricsEvent(&bind.WatchOpts{}, ethEventChan)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &PerfEvent{
				Contract:       v,
				RequestID:      event.RequestId,
				Round:          event.RoundID,
				BlockTimestamp: event.Timestamp,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// RandomnessOutput get VRF randomness output
func (v *EthereumVRFConsumer) RandomnessOutput(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	out, err := v.consumer.RandomnessOutput(opts)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VRFConsumerRoundConfirmer is a header subscription that awaits for a certain VRF round to be completed
type VRFConsumerRoundConfirmer struct {
	consumer VRFConsumer
	roundID  *big.Int
	doneChan chan struct{}
	context  context.Context
	cancel   context.CancelFunc
	done     bool
}

// NewVRFConsumerRoundConfirmer provides a new instance of a NewVRFConsumerRoundConfirmer
func NewVRFConsumerRoundConfirmer(
	contract VRFConsumer,
	roundID *big.Int,
	timeout time.Duration,
) *VRFConsumerRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &VRFConsumerRoundConfirmer{
		consumer: contract,
		roundID:  roundID,
		doneChan: make(chan struct{}),
		context:  ctx,
		cancel:   ctxCancel,
	}
}

// ReceiveHeader will query the latest VRFConsumer round and check to see whether the round has confirmed
func (f *VRFConsumerRoundConfirmer) ReceiveHeader(header blockchain.NodeHeader) error {
	if f.done {
		return nil
	}
	roundID, err := f.consumer.CurrentRoundID(context.Background())
	if err != nil {
		return err
	}
	l := log.Debug().
		Str("Contract Address", f.consumer.Address()).
		Int64("Waiting for Round", f.roundID.Int64()).
		Int64("Current round ID", roundID.Int64()).
		Uint64("Header Number", header.Number.Uint64())
	if roundID.Int64() == f.roundID.Int64() {
		randomness, err := f.consumer.RandomnessOutput(context.Background())
		if err != nil {
			return err
		}
		l.Uint64("Randomness", randomness.Uint64()).
			Msg("VRFConsumer round completed")
		f.done = true
		f.doneChan <- struct{}{}
	} else {
		l.Msg("Waiting for VRFConsumer round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (f *VRFConsumerRoundConfirmer) Wait() error {
	for {
		select {
		case <-f.doneChan:
			f.cancel()
			return nil
		case <-f.context.Done():
			return fmt.Errorf("timeout waiting for VRFConsumer round to confirm: %d", f.roundID)
		}
	}
}

// EthereumVRF represents a VRF contract
type EthereumVRF struct {
	client  blockchain.EVMClient
	vrf     *ethereum.VRF
	address *common.Address
}

// Fund sends specified currencies to the contract
func (v *EthereumVRF) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

// ProofLength returns the PROOFLENGTH call from the VRF contract
func (v *EthereumVRF) ProofLength(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return v.vrf.PROOFLENGTH(opts)
}

// EthereumDKG represents DKG contract
type EthereumDKG struct {
	address *common.Address
	client  blockchain.EVMClient
	dkg     *dkg.DKG
}

func (dkgContract *EthereumDKG) Address() string {
	return dkgContract.address.Hex()
}

func (dkgContract *EthereumDKG) AddClient(keyID string, clientAddress string) error {
	keyIDBytes, err := decodeHexTo32ByteArray(keyID)
	if err != nil {
		return err
	}
	opts, err := dkgContract.client.TransactionOpts(dkgContract.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := dkgContract.dkg.AddClient(
		opts,
		keyIDBytes,
		common.HexToAddress(clientAddress),
	)
	if err != nil {
		return err
	}
	return dkgContract.client.ProcessTransaction(tx)

}

func (dkgContract *EthereumDKG) SetConfig(
	signerAddresses []common.Address,
	transmitterAddresses []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	opts, err := dkgContract.client.TransactionOpts(dkgContract.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := dkgContract.dkg.SetConfig(
		opts,
		signerAddresses,
		transmitterAddresses,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	if err != nil {
		return err
	}
	return dkgContract.client.ProcessTransaction(tx)
}

func (dkgContract *EthereumDKG) WaitForTransmittedEvent() (*dkg.DKGTransmitted, error) {
	transmittedEventsChannel := make(chan *dkg.DKGTransmitted)
	subscription, err := dkgContract.dkg.WatchTransmitted(nil, transmittedEventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()
	transmittedEvent := <-transmittedEventsChannel
	if err != nil {
		return nil, err
	}
	return transmittedEvent, nil
}

func (dkgContract *EthereumDKG) WaitForConfigSetEvent() (*dkg.DKGConfigSet, error) {

	configSetEventsChannel := make(chan *dkg.DKGConfigSet)
	subscription, err := dkgContract.dkg.WatchConfigSet(nil, configSetEventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()
	configSetEvent := <-configSetEventsChannel
	if err != nil {
		return nil, err
	}
	return configSetEvent, nil
}

// EthereumVRFCoordinatorV3 represents VRFCoordinatorV3 contract
type EthereumVRFCoordinatorV3 struct {
	address          *common.Address
	client           blockchain.EVMClient
	vrfCoordinatorV3 *vrf_coordinator.VRFCoordinator
}

func (coordinator *EthereumVRFCoordinatorV3) Address() string {
	return coordinator.address.Hex()
}

func (coordinator *EthereumVRFCoordinatorV3) SetProducer(producerAddress string) error {
	opts, err := coordinator.client.TransactionOpts(coordinator.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := coordinator.vrfCoordinatorV3.SetProducer(
		opts,
		common.HexToAddress(producerAddress),
	)
	if err != nil {
		return err
	}
	return coordinator.client.ProcessTransaction(tx)
}

func (coordinator *EthereumVRFCoordinatorV3) CreateSubscription() error {
	opts, err := coordinator.client.TransactionOpts(coordinator.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := coordinator.vrfCoordinatorV3.CreateSubscription(
		opts,
	)
	if err != nil {
		return err
	}
	return coordinator.client.ProcessTransaction(tx)
}

func (coordinator *EthereumVRFCoordinatorV3) AddConsumer(subId uint64, consumerAddress string) error {
	opts, err := coordinator.client.TransactionOpts(coordinator.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := coordinator.vrfCoordinatorV3.AddConsumer(
		opts,
		subId,
		common.HexToAddress(consumerAddress),
	)
	if err != nil {
		return err
	}
	return coordinator.client.ProcessTransaction(tx)
}

// EthereumVRFBeacon represents VRFBeacon contract
type EthereumVRFBeacon struct {
	address   *common.Address
	client    blockchain.EVMClient
	vrfBeacon *vrf_beacon.VRFBeacon
}

func (beacon *EthereumVRFBeacon) Address() string {
	return beacon.address.Hex()
}

func (beacon *EthereumVRFBeacon) SetPayees(transmitterAddresses []common.Address, payeesAddresses []common.Address) error {
	opts, err := beacon.client.TransactionOpts(beacon.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := beacon.vrfBeacon.SetPayees(
		opts,
		transmitterAddresses,
		payeesAddresses,
	)
	if err != nil {
		return err
	}
	return beacon.client.ProcessTransaction(tx)
}

func (beacon *EthereumVRFBeacon) SetConfig(
	signerAddresses []common.Address,
	transmitterAddresses []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	opts, err := beacon.client.TransactionOpts(beacon.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := beacon.vrfBeacon.SetConfig(
		opts,
		signerAddresses,
		transmitterAddresses,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	if err != nil {
		return err
	}
	return beacon.client.ProcessTransaction(tx)
}

func (beacon *EthereumVRFBeacon) WaitForConfigSetEvent() (*vrf_beacon.VRFBeaconConfigSet, error) {
	configSetEventsChannel := make(chan *vrf_beacon.VRFBeaconConfigSet)
	subscription, err := beacon.vrfBeacon.WatchConfigSet(nil, configSetEventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()
	configSetEvent := <-configSetEventsChannel
	if err != nil {
		return nil, err
	}
	return configSetEvent, nil
}

func (beacon *EthereumVRFBeacon) WaitForNewTransmissionEvent() (*vrf_beacon.VRFBeaconNewTransmission, error) {
	newTransmissionEventsChannel := make(chan *vrf_beacon.VRFBeaconNewTransmission)
	subscription, err := beacon.vrfBeacon.WatchNewTransmission(nil, newTransmissionEventsChannel, nil, nil)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()
	newTransmissionEvent := <-newTransmissionEventsChannel
	if err != nil {
		return nil, err
	}
	return newTransmissionEvent, nil
}

func (beacon *EthereumVRFBeacon) LatestConfigDigestAndEpoch(ctx context.Context) (vrf_beacon.LatestConfigDigestAndEpoch,
	error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(beacon.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return beacon.vrfBeacon.LatestConfigDigestAndEpoch(opts)
}

// EthereumVRFBeaconConsumer represents VRFBeaconConsumer contract
type EthereumVRFBeaconConsumer struct {
	address           *common.Address
	client            blockchain.EVMClient
	vrfBeaconConsumer *vrf_beacon_consumer.BeaconVRFConsumer
	abi               string
}

func (consumer *EthereumVRFBeaconConsumer) Address() string {
	return consumer.address.Hex()
}

func (consumer *EthereumVRFBeaconConsumer) RequestRandomness(
	numWords uint16,
	subID uint64,
	confirmationDelayArg *big.Int,
) (*types.Receipt, error) {
	opts, err := consumer.client.TransactionOpts(consumer.client.GetDefaultWallet())
	if err != nil {
		return nil, errors.Wrap(err, "TransactionOpts failed")
	}
	tx, err := consumer.vrfBeaconConsumer.TestRequestRandomness(
		opts,
		numWords,
		subID,
		confirmationDelayArg,
	)
	if err != nil {
		return nil, errors.Wrap(err, "TestRequestRandomness failed")
	}
	err = consumer.client.ProcessTransaction(tx)
	if err != nil {
		return nil, errors.Wrap(err, "ProcessTransaction failed")
	}
	err = consumer.client.WaitForEvents()

	if err != nil {
		return nil, errors.Wrap(err, "WaitForEvents failed")
	}
	receipt, err := consumer.client.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "GetTxReceipt failed")
	}
	return receipt, nil
}

func (consumer *EthereumVRFBeaconConsumer) RedeemRandomness(
	requestID *big.Int,
) error {
	opts, err := consumer.client.TransactionOpts(consumer.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := consumer.vrfBeaconConsumer.TestRedeemRandomness(
		opts,
		requestID,
	)
	if err != nil {
		return err
	}
	return consumer.client.ProcessTransaction(tx)
}

func (consumer *EthereumVRFBeaconConsumer) RequestRandomnessFulfillment(
	numWords uint16,
	subID uint64,
	confirmationDelayArg *big.Int,
	callbackGasLimit uint32,
	arguments []byte,
) (*types.Receipt, error) {
	opts, err := consumer.client.TransactionOpts(consumer.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := consumer.vrfBeaconConsumer.TestRequestRandomnessFulfillment(
		opts,
		subID,
		numWords,
		confirmationDelayArg,
		callbackGasLimit,
		arguments,
	)
	if err != nil {
		return nil, errors.Wrap(err, "TestRequestRandomnessFulfillment failed")
	}
	err = consumer.client.ProcessTransaction(tx)
	if err != nil {
		return nil, errors.Wrap(err, "ProcessTransaction failed")
	}
	err = consumer.client.WaitForEvents()

	if err != nil {
		return nil, errors.Wrap(err, "WaitForEvents failed")
	}
	receipt, err := consumer.client.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "GetTxReceipt failed")
	}
	return receipt, nil
}

func (consumer *EthereumVRFBeaconConsumer) IBeaconPeriodBlocks(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(consumer.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return consumer.vrfBeaconConsumer.IBeaconPeriodBlocks(opts)
}

func (consumer *EthereumVRFBeaconConsumer) GetRequestIdsBy(ctx context.Context, nextBeaconOutputHeight *big.Int, confDelay *big.Int) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(consumer.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return consumer.vrfBeaconConsumer.SRequestsIDs(opts, nextBeaconOutputHeight, confDelay)
}

func (consumer *EthereumVRFBeaconConsumer) GetRandomnessByRequestId(ctx context.Context, requestID *big.Int, numWordIndex *big.Int) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(consumer.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return consumer.vrfBeaconConsumer.SReceivedRandomnessByRequestID(opts, requestID, numWordIndex)
}

// EthereumBatchBlockhashStore represents BatchBlockhashStore contract
type EthereumBatchBlockhashStore struct {
	address             *common.Address
	client              blockchain.EVMClient
	batchBlockhashStore *batch_blockhash_store.BatchBlockhashStore
}

func (v *EthereumBatchBlockhashStore) Address() string {
	return v.address.Hex()
}
