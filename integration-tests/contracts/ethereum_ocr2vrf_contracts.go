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
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
)

// EthereumDKG represents DKG contract
type EthereumDKG struct {
	address *common.Address
	client  blockchain.EVMClient
	dkg     *dkg.DKG
}

// EthereumVRFCoordinatorV3 represents VRFCoordinatorV3 contract
type EthereumVRFCoordinatorV3 struct {
	address          *common.Address
	client           blockchain.EVMClient
	vrfCoordinatorV3 *vrf_coordinator.VRFCoordinator
}

// EthereumVRFBeacon represents VRFBeacon contract
type EthereumVRFBeacon struct {
	address   *common.Address
	client    blockchain.EVMClient
	vrfBeacon *vrf_beacon.VRFBeacon
}

// EthereumVRFBeaconConsumer represents VRFBeaconConsumer contract
type EthereumVRFBeaconConsumer struct {
	address           *common.Address
	client            blockchain.EVMClient
	vrfBeaconConsumer *vrf_beacon_consumer.BeaconVRFConsumer
}

// EthereumVRFCoordinator represents VRF coordinator contract
type EthereumVRFCoordinator struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *solidity_vrf_coordinator_interface.VRFCoordinator
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
func (e *EthereumContractDeployer) DeployOCR2VRFCoordinator(beaconPeriodBlocksCount *big.Int, linkAddress string) (VRFCoordinatorV3, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV3", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator.DeployVRFCoordinator(auth, backend, beaconPeriodBlocksCount, common.HexToAddress(linkAddress))
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
	keyIDBytes, err := DecodeHexTo32ByteArray(keyId)
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
func DecodeHexTo32ByteArray(val string) ([32]byte, error) {
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
func (e *EthereumContractDeployer) DeployVRFBeaconConsumer(vrfCoordinatorAddress string, beaconPeriodBlockCount *big.Int) (VRFBeaconConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFBeaconConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_beacon_consumer.DeployBeaconVRFConsumer(auth, backend, common.HexToAddress(vrfCoordinatorAddress), false, beaconPeriodBlockCount)
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

func (dkgContract *EthereumDKG) Address() string {
	return dkgContract.address.Hex()
}

func (dkgContract *EthereumDKG) AddClient(keyID string, clientAddress string) error {
	keyIDBytes, err := DecodeHexTo32ByteArray(keyID)
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

func (dkgContract *EthereumDKG) WaitForTransmittedEvent(timeout time.Duration) (*dkg.DKGTransmitted, error) {
	transmittedEventsChannel := make(chan *dkg.DKGTransmitted)
	subscription, err := dkgContract.dkg.WatchTransmitted(nil, transmittedEventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err = <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for DKGTransmitted event")
		case transmittedEvent := <-transmittedEventsChannel:
			return transmittedEvent, nil
		}
	}
}

func (dkgContract *EthereumDKG) WaitForConfigSetEvent(timeout time.Duration) (*dkg.DKGConfigSet, error) {
	configSetEventsChannel := make(chan *dkg.DKGConfigSet)
	subscription, err := dkgContract.dkg.WatchConfigSet(nil, configSetEventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err = <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for DKGConfigSet event")
		case configSetEvent := <-configSetEventsChannel:
			return configSetEvent, nil
		}
	}
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

func (coordinator *EthereumVRFCoordinatorV3) FindSubscriptionID() (*big.Int, error) {
	fopts := &bind.FilterOpts{}
	owner := coordinator.client.GetDefaultWallet().Address()

	subscriptionIterator, err := coordinator.vrfCoordinatorV3.FilterSubscriptionCreated(
		fopts,
		nil,
		[]common.Address{common.HexToAddress(owner)},
	)
	if err != nil {
		return nil, err
	}

	if !subscriptionIterator.Next() {
		return nil, fmt.Errorf("expected at least 1 subID for the given owner %s", owner)
	}

	return subscriptionIterator.Event.SubId, nil
}

func (coordinator *EthereumVRFCoordinatorV3) AddConsumer(subId *big.Int, consumerAddress string) error {
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

func (coordinator *EthereumVRFCoordinatorV3) SetConfig(maxCallbackGasLimit uint32, maxCallbackArgumentsLength uint32) error {
	opts, err := coordinator.client.TransactionOpts(coordinator.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := coordinator.vrfCoordinatorV3.SetCallbackConfig(
		opts,
		vrf_coordinator.VRFCoordinatorCallbackConfig{
			MaxCallbackGasLimit:        maxCallbackGasLimit,
			MaxCallbackArgumentsLength: maxCallbackArgumentsLength, // 5 EVM words
		},
	)
	if err != nil {
		return err
	}
	return coordinator.client.ProcessTransaction(tx)
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

func (beacon *EthereumVRFBeacon) WaitForConfigSetEvent(timeout time.Duration) (*vrf_beacon.VRFBeaconConfigSet, error) {
	configSetEventsChannel := make(chan *vrf_beacon.VRFBeaconConfigSet)
	subscription, err := beacon.vrfBeacon.WatchConfigSet(nil, configSetEventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for config set event")
		case configSetEvent := <-configSetEventsChannel:
			return configSetEvent, nil
		}
	}
}

func (beacon *EthereumVRFBeacon) WaitForNewTransmissionEvent(timeout time.Duration) (*vrf_beacon.VRFBeaconNewTransmission, error) {
	newTransmissionEventsChannel := make(chan *vrf_beacon.VRFBeaconNewTransmission)
	subscription, err := beacon.vrfBeacon.WatchNewTransmission(nil, newTransmissionEventsChannel, nil)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for new transmission event")
		case newTransmissionEvent := <-newTransmissionEventsChannel:
			return newTransmissionEvent, nil
		}
	}
}

func (beacon *EthereumVRFBeacon) LatestConfigDigestAndEpoch(ctx context.Context) (vrf_beacon.LatestConfigDigestAndEpoch,
	error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(beacon.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return beacon.vrfBeacon.LatestConfigDigestAndEpoch(opts)
}

func (consumer *EthereumVRFBeaconConsumer) Address() string {
	return consumer.address.Hex()
}

func (consumer *EthereumVRFBeaconConsumer) RequestRandomness(
	numWords uint16,
	subID, confirmationDelayArg *big.Int,
) (*types.Receipt, error) {
	opts, err := consumer.client.TransactionOpts(consumer.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("TransactionOpts failed, err: %w", err)
	}
	tx, err := consumer.vrfBeaconConsumer.TestRequestRandomness(
		opts,
		numWords,
		subID,
		confirmationDelayArg,
	)
	if err != nil {
		return nil, fmt.Errorf("TestRequestRandomness failed, err: %w", err)
	}
	err = consumer.client.ProcessTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("ProcessTransaction failed, err: %w", err)
	}
	err = consumer.client.WaitForEvents()

	if err != nil {
		return nil, fmt.Errorf("WaitForEvents failed, err: %w", err)
	}
	receipt, err := consumer.client.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("GetTxReceipt failed, err: %w", err)
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confirmationDelayArg).
		Msg("RequestRandomness called")
	return receipt, nil
}

func (consumer *EthereumVRFBeaconConsumer) RedeemRandomness(
	subID, requestID *big.Int,
) error {
	opts, err := consumer.client.TransactionOpts(consumer.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := consumer.vrfBeaconConsumer.TestRedeemRandomness(
		opts,
		subID,
		requestID,
	)
	if err != nil {
		return err
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Request ID", requestID).
		Msg("RedeemRandomness called")
	return consumer.client.ProcessTransaction(tx)
}

func (consumer *EthereumVRFBeaconConsumer) RequestRandomnessFulfillment(
	numWords uint16,
	subID, confirmationDelayArg *big.Int,
	requestGasLimit uint32,
	callbackGasLimit uint32,
	arguments []byte,
) (*types.Receipt, error) {
	opts, err := consumer.client.TransactionOpts(consumer.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	// overriding gas limit because gas estimated by TestRequestRandomnessFulfillment
	// is incorrect
	opts.GasLimit = uint64(requestGasLimit)
	tx, err := consumer.vrfBeaconConsumer.TestRequestRandomnessFulfillment(
		opts,
		subID,
		numWords,
		confirmationDelayArg,
		callbackGasLimit,
		arguments,
	)
	if err != nil {
		return nil, fmt.Errorf("TestRequestRandomnessFulfillment failed, err: %w", err)
	}
	err = consumer.client.ProcessTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("ProcessTransaction failed, err: %w", err)
	}
	err = consumer.client.WaitForEvents()

	if err != nil {
		return nil, fmt.Errorf("WaitForEvents failed, err: %w", err)
	}
	receipt, err := consumer.client.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("GetTxReceipt failed, err: %w", err)
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confirmationDelayArg).
		Interface("Callback Gas Limit", callbackGasLimit).
		Msg("RequestRandomnessFulfillment called")
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
