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

	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
)

// EthereumDKG represents DKG contract
type EthereumDKG struct {
	address common.Address
	client  *seth.Client
	dkg     *dkg.DKG
}

// EthereumVRFCoordinatorV3 represents VRFCoordinatorV3 contract
type EthereumVRFCoordinatorV3 struct {
	address          common.Address
	client           *seth.Client
	vrfCoordinatorV3 *vrf_coordinator.VRFCoordinator
}

// EthereumVRFBeacon represents VRFBeacon contract
type EthereumVRFBeacon struct {
	address   common.Address
	client    *seth.Client
	vrfBeacon *vrf_beacon.VRFBeacon
}

// EthereumVRFBeaconConsumer represents VRFBeaconConsumer contract
type EthereumVRFBeaconConsumer struct {
	address           common.Address
	client            *seth.Client
	vrfBeaconConsumer *vrf_beacon_consumer.BeaconVRFConsumer
}

// DeployDKG deploys DKG contract
func DeployDKG(seth *seth.Client) (DKG, error) {
	abi, err := dkg.DKGMetaData.GetAbi()
	if err != nil {
		return &EthereumDKG{}, fmt.Errorf("failed to get DKG ABI: %w", err)
	}

	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"DKG",
		*abi,
		common.FromHex(dkg.DKGMetaData.Bin))
	if err != nil {
		return &EthereumDKG{}, fmt.Errorf("DKG instance deployment have failed: %w", err)
	}

	contract, err := dkg.NewDKG(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumDKG{}, fmt.Errorf("failed to instantiate DKG instance: %w", err)
	}

	return &EthereumDKG{
		client:  seth,
		dkg:     contract,
		address: data.Address,
	}, err
}

// DeployOCR2VRFCoordinator deploys CR2VRFCoordinator contract
func DeployOCR2VRFCoordinator(seth *seth.Client, beaconPeriodBlocksCount *big.Int, linkAddress string) (VRFCoordinatorV3, error) {
	abi, err := vrf_coordinator.VRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinatorV3{}, fmt.Errorf("failed to get VRFCoordinatorV3 ABI: %w", err)
	}

	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorV3",
		*abi,
		common.FromHex(vrf_coordinator.VRFCoordinatorMetaData.Bin),
		beaconPeriodBlocksCount, common.HexToAddress(linkAddress))
	if err != nil {
		return &EthereumVRFCoordinatorV3{}, fmt.Errorf("VRFCoordinatorV3 instance deployment have failed: %w", err)
	}

	contract, err := vrf_coordinator.NewVRFCoordinator(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFCoordinatorV3{}, fmt.Errorf("failed to instantiate VRFCoordinatorV3 instance: %w", err)
	}

	return &EthereumVRFCoordinatorV3{
		client:           seth,
		vrfCoordinatorV3: contract,
		address:          data.Address,
	}, err
}

// DeployVRFBeacon deploys DeployVRFBeacon contract
func DeployVRFBeacon(seth *seth.Client, vrfCoordinatorAddress string, linkAddress string, dkgAddress string, keyId string) (VRFBeacon, error) {
	abi, err := vrf_beacon.VRFBeaconMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFBeacon{}, fmt.Errorf("failed to get VRFBeacon ABI: %w", err)
	}
	keyIDBytes, err := DecodeHexTo32ByteArray(keyId)
	if err != nil {
		return nil, err
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFBeacon",
		*abi,
		common.FromHex(vrf_beacon.VRFBeaconMetaData.Bin),
		common.HexToAddress(linkAddress), common.HexToAddress(vrfCoordinatorAddress), common.HexToAddress(dkgAddress), keyIDBytes)
	if err != nil {
		return &EthereumVRFBeacon{}, fmt.Errorf("VRFBeacon instance deployment have failed: %w", err)
	}

	contract, err := vrf_beacon.NewVRFBeacon(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFBeacon{}, fmt.Errorf("failed to instantiate VRFBeacon instance: %w", err)
	}

	return &EthereumVRFBeacon{
		client:    seth,
		vrfBeacon: contract,
		address:   data.Address,
	}, err
}

// DeployBatchBlockhashStore deploys DeployBatchBlockhashStore contract
func DeployBatchBlockhashStore(seth *seth.Client, blockhashStoreAddr string) (BatchBlockhashStore, error) {
	abi, err := batch_blockhash_store.BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return &EthereumBatchBlockhashStore{}, fmt.Errorf("failed to get BatchBlockhashStore ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"BatchBlockhashStore",
		*abi,
		common.FromHex(batch_blockhash_store.BatchBlockhashStoreMetaData.Bin),
		common.HexToAddress(blockhashStoreAddr))
	if err != nil {
		return &EthereumBatchBlockhashStore{}, fmt.Errorf("BatchBlockhashStore instance deployment have failed: %w", err)
	}

	contract, err := batch_blockhash_store.NewBatchBlockhashStore(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumBatchBlockhashStore{}, fmt.Errorf("failed to instantiate BatchBlockhashStore instance: %w", err)
	}

	return &EthereumBatchBlockhashStore{
		client:              seth,
		batchBlockhashStore: contract,
		address:             data.Address,
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
func DeployVRFBeaconConsumer(seth *seth.Client, vrfCoordinatorAddress string, beaconPeriodBlockCount *big.Int) (VRFBeaconConsumer, error) {
	abi, err := vrf_beacon_consumer.BeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFBeaconConsumer{}, fmt.Errorf("failed to get VRFBeaconConsumer ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFBeaconConsumer",
		*abi,
		common.FromHex(vrf_beacon_consumer.BeaconVRFConsumerMetaData.Bin),
		common.HexToAddress(vrfCoordinatorAddress), false, beaconPeriodBlockCount)
	if err != nil {
		return &EthereumVRFBeaconConsumer{}, fmt.Errorf("VRFBeaconConsumer instance deployment have failed: %w", err)
	}

	contract, err := vrf_beacon_consumer.NewBeaconVRFConsumer(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFBeaconConsumer{}, fmt.Errorf("failed to instantiate VRFBeaconConsumer instance: %w", err)
	}

	return &EthereumVRFBeaconConsumer{
		client:            seth,
		vrfBeaconConsumer: contract,
		address:           data.Address,
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
	_, err = dkgContract.client.Decode(dkgContract.dkg.AddClient(
		dkgContract.client.NewTXOpts(),
		keyIDBytes,
		common.HexToAddress(clientAddress),
	))
	return err
}

func (dkgContract *EthereumDKG) SetConfig(
	signerAddresses []common.Address,
	transmitterAddresses []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	_, err := dkgContract.client.Decode(dkgContract.dkg.SetConfig(
		dkgContract.client.NewTXOpts(),
		signerAddresses,
		transmitterAddresses,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	))
	return err
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
	_, err := coordinator.client.Decode(coordinator.vrfCoordinatorV3.SetProducer(
		coordinator.client.NewTXOpts(),
		common.HexToAddress(producerAddress),
	))
	return err
}

func (coordinator *EthereumVRFCoordinatorV3) CreateSubscription() error {
	_, err := coordinator.client.Decode(coordinator.vrfCoordinatorV3.CreateSubscription(
		coordinator.client.NewTXOpts(),
	))
	return err
}

func (coordinator *EthereumVRFCoordinatorV3) FindSubscriptionID() (*big.Int, error) {
	fopts := &bind.FilterOpts{}
	owner := coordinator.client.MustGetRootKeyAddress()

	subscriptionIterator, err := coordinator.vrfCoordinatorV3.FilterSubscriptionCreated(
		fopts,
		nil,
		[]common.Address{owner},
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
	_, err := coordinator.client.Decode(coordinator.vrfCoordinatorV3.AddConsumer(
		coordinator.client.NewTXOpts(),
		subId,
		common.HexToAddress(consumerAddress),
	))
	return err
}

func (coordinator *EthereumVRFCoordinatorV3) SetConfig(maxCallbackGasLimit uint32, maxCallbackArgumentsLength uint32) error {
	_, err := coordinator.client.Decode(coordinator.vrfCoordinatorV3.SetCallbackConfig(
		coordinator.client.NewTXOpts(),
		vrf_coordinator.VRFCoordinatorCallbackConfig{
			MaxCallbackGasLimit:        maxCallbackGasLimit,
			MaxCallbackArgumentsLength: maxCallbackArgumentsLength, // 5 EVM words
		},
	))
	return err
}

func (beacon *EthereumVRFBeacon) Address() string {
	return beacon.address.Hex()
}

func (beacon *EthereumVRFBeacon) SetPayees(transmitterAddresses []common.Address, payeesAddresses []common.Address) error {
	_, err := beacon.client.Decode(beacon.vrfBeacon.SetPayees(
		beacon.client.NewTXOpts(),
		transmitterAddresses,
		payeesAddresses,
	))
	return err
}

func (beacon *EthereumVRFBeacon) SetConfig(
	signerAddresses []common.Address,
	transmitterAddresses []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	_, err := beacon.client.Decode(beacon.vrfBeacon.SetConfig(
		beacon.client.NewTXOpts(),
		signerAddresses,
		transmitterAddresses,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	))
	return err
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
		From:    beacon.client.MustGetRootKeyAddress(),
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
	tx, err := consumer.client.Decode(consumer.vrfBeaconConsumer.TestRequestRandomness(
		consumer.client.NewTXOpts(),
		numWords,
		subID,
		confirmationDelayArg,
	))
	if err != nil {
		return nil, fmt.Errorf("TestRequestRandomness failed, err: %w", err)
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confirmationDelayArg).
		Msg("RequestRandomness called")
	return tx.Receipt, nil
}

func (consumer *EthereumVRFBeaconConsumer) RedeemRandomness(
	subID, requestID *big.Int,
) error {
	_, err := consumer.client.Decode(consumer.vrfBeaconConsumer.TestRedeemRandomness(
		consumer.client.NewTXOpts(),
		subID,
		requestID,
	))
	if err != nil {
		return err
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Request ID", requestID).
		Msg("RedeemRandomness called")
	return nil
}

func (consumer *EthereumVRFBeaconConsumer) RequestRandomnessFulfillment(
	numWords uint16,
	subID, confirmationDelayArg *big.Int,
	requestGasLimit uint32,
	callbackGasLimit uint32,
	arguments []byte,
) (*types.Receipt, error) {
	opts := consumer.client.NewTXOpts()
	// overriding gas limit because gas estimated by TestRequestRandomnessFulfillment
	// is incorrect
	opts.GasLimit = uint64(requestGasLimit)
	tx, err := consumer.client.Decode(consumer.vrfBeaconConsumer.TestRequestRandomnessFulfillment(
		opts,
		subID,
		numWords,
		confirmationDelayArg,
		callbackGasLimit,
		arguments,
	))
	if err != nil {
		return nil, fmt.Errorf("TestRequestRandomnessFulfillment failed, err: %w", err)
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confirmationDelayArg).
		Interface("Callback Gas Limit", callbackGasLimit).
		Msg("RequestRandomnessFulfillment called")
	return tx.Receipt, nil
}

func (consumer *EthereumVRFBeaconConsumer) IBeaconPeriodBlocks(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    consumer.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	return consumer.vrfBeaconConsumer.IBeaconPeriodBlocks(opts)
}

func (consumer *EthereumVRFBeaconConsumer) GetRequestIdsBy(ctx context.Context, nextBeaconOutputHeight *big.Int, confDelay *big.Int) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    consumer.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	return consumer.vrfBeaconConsumer.SRequestsIDs(opts, nextBeaconOutputHeight, confDelay)
}

func (consumer *EthereumVRFBeaconConsumer) GetRandomnessByRequestId(ctx context.Context, requestID *big.Int, numWordIndex *big.Int) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    consumer.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	return consumer.vrfBeaconConsumer.SReceivedRandomnessByRequestID(opts, requestID, numWordIndex)
}
