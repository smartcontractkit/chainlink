package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_test_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_mock_ethlink_aggregator"
)

// EthereumBlockhashStore represents a blockhash store for VRF contract
type EthereumBlockhashStore struct {
	address        *common.Address
	client         *seth.Client
	blockHashStore *blockhash_store.BlockhashStore
}

// EthereumVRFCoordinator represents VRF coordinator contract
type EthereumVRFCoordinator struct {
	address     *common.Address
	client      *seth.Client
	coordinator *solidity_vrf_coordinator_interface.VRFCoordinator
}

type EthereumVRFCoordinatorTestV2 struct {
	address     *common.Address
	client      *seth.Client
	coordinator *vrf_coordinator_test_v2.VRFCoordinatorTestV2
}

func (v *EthereumVRFCoordinatorTestV2) Address() string {
	return v.address.Hex()
}

// EthereumVRFConsumer represents VRF consumer contract
type EthereumVRFConsumer struct {
	address  *common.Address
	client   *seth.Client
	consumer *solidity_vrf_consumer_interface.VRFConsumer
}

// EthereumVRF represents a VRF contract
type EthereumVRF struct {
	client  *seth.Client
	vrf     *solidity_vrf_wrapper.VRF
	address *common.Address
}

type EthereumVRFMockETHLINKAggregator struct {
	client   *seth.Client
	address  *common.Address
	contract *vrf_mock_ethlink_aggregator.VRFMockETHLINKAggregator
}

// EthereumBatchBlockhashStore represents BatchBlockhashStore contract
type EthereumBatchBlockhashStore struct {
	address             common.Address
	client              *seth.Client
	batchBlockhashStore *batch_blockhash_store.BatchBlockhashStore
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
	logFields := map[string]any{
		"Contract Address":  f.consumer.Address(),
		"Waiting for Round": f.roundID.Int64(),
		"Current Round ID":  roundID.Int64(),
		"Header Number":     header.Number.Uint64(),
	}
	if roundID.Int64() == f.roundID.Int64() {
		randomness, err := f.consumer.RandomnessOutput(context.Background())
		if err != nil {
			return err
		}
		log.Info().Fields(logFields).Uint64("Randomness", randomness.Uint64()).Msg("VRFConsumer round completed")
		f.done = true
		f.doneChan <- struct{}{}
	} else {
		log.Debug().Fields(logFields).Msg("Waiting for VRFConsumer round")
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

func (v *EthereumBatchBlockhashStore) Address() string {
	return v.address.Hex()
}

func (a *EthereumVRFMockETHLINKAggregator) Address() string {
	return a.address.Hex()
}

func (a *EthereumVRFMockETHLINKAggregator) LatestRoundData() (*big.Int, error) {
	data, err := a.contract.LatestRoundData(a.client.NewCallOpts())
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (a *EthereumVRFMockETHLINKAggregator) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := a.contract.LatestRoundData(a.client.NewCallOpts())
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

func (a *EthereumVRFMockETHLINKAggregator) SetBlockTimestampDeduction(blockTimestampDeduction *big.Int) error {
	_, err := a.client.Decode(a.contract.SetBlockTimestampDeduction(a.client.NewTXOpts(), blockTimestampDeduction))
	return err
}

// DeployVRFContract deploy VRFv1 contract
func DeployVRFv1Contract(seth *seth.Client) (VRF, error) {
	abi, err := solidity_vrf_wrapper.VRFMetaData.GetAbi()
	if err != nil {
		return &EthereumVRF{}, fmt.Errorf("failed to get VRF ABI: %w", err)
	}

	vrfDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRF",
		*abi,
		common.FromHex(solidity_vrf_wrapper.VRFMetaData.Bin))
	if err != nil {
		return &EthereumVRF{}, fmt.Errorf("VRF instance deployment have failed: %w", err)
	}

	vrf, err := solidity_vrf_wrapper.NewVRF(vrfDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRF{}, fmt.Errorf("failed to instantiate VRF instance: %w", err)
	}

	return &EthereumVRF{
		client:  seth,
		vrf:     vrf,
		address: &vrfDeploymentData.Address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func DeployBlockhashStore(seth *seth.Client) (BlockHashStore, error) {
	abi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return &EthereumBlockhashStore{}, fmt.Errorf("failed to get BlockhashStore ABI: %w", err)
	}

	storeDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"BlockhashStore",
		*abi,
		common.FromHex(blockhash_store.BlockhashStoreMetaData.Bin))
	if err != nil {
		return &EthereumBlockhashStore{}, fmt.Errorf("BlockhashStore instance deployment have failed: %w", err)
	}

	store, err := blockhash_store.NewBlockhashStore(storeDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumBlockhashStore{}, fmt.Errorf("failed to instantiate BlockhashStore instance: %w", err)
	}

	return &EthereumBlockhashStore{
		client:         seth,
		blockHashStore: store,
		address:        &storeDeploymentData.Address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func DeployVRFCoordinator(seth *seth.Client, linkAddr, bhsAddr string) (VRFCoordinator, error) {
	abi, err := solidity_vrf_coordinator_interface.VRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinator{}, fmt.Errorf("failed to get VRFCoordinator ABI: %w", err)
	}

	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinator",
		*abi,
		common.FromHex(solidity_vrf_coordinator_interface.VRFCoordinatorMetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return &EthereumVRFCoordinator{}, fmt.Errorf("VRFCoordinator instance deployment have failed: %w", err)
	}

	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFCoordinator{}, fmt.Errorf("failed to instantiate VRFCoordinator instance: %w", err)
	}

	return &EthereumVRFCoordinator{
		client:      seth,
		coordinator: coordinator,
		address:     &coordinatorDeploymentData.Address,
	}, err
}

func DeployVRFCoordinatorTestV2(seth *seth.Client, linkAddr, bhsAddr, linkEthFeedAddr string) (*EthereumVRFCoordinatorTestV2, error) {
	abi, err := vrf_coordinator_test_v2.VRFCoordinatorTestV2MetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinatorTestV2{}, fmt.Errorf("failed to get VRFCoordinatorTestV2 ABI: %w", err)
	}

	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorTestV2",
		*abi,
		common.FromHex(vrf_coordinator_test_v2.VRFCoordinatorTestV2MetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(bhsAddr),
		common.HexToAddress(linkEthFeedAddr))
	if err != nil {
		return &EthereumVRFCoordinatorTestV2{}, fmt.Errorf("VRFCoordinatorTestV2 instance deployment have failed: %w", err)
	}

	coordinator, err := vrf_coordinator_test_v2.NewVRFCoordinatorTestV2(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFCoordinatorTestV2{}, fmt.Errorf("failed to instantiate VRFCoordinatorTestV2 instance: %w", err)
	}

	return &EthereumVRFCoordinatorTestV2{
		client:      seth,
		coordinator: coordinator,
		address:     &coordinatorDeploymentData.Address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func DeployVRFConsumer(seth *seth.Client, linkAddr, coordinatorAddr string) (VRFConsumer, error) {
	abi, err := solidity_vrf_consumer_interface.VRFConsumerMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFConsumer{}, fmt.Errorf("failed to get VRFConsumer ABI: %w", err)
	}

	consumerDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFConsumer",
		*abi,
		common.FromHex(solidity_vrf_consumer_interface.VRFConsumerMetaData.Bin),
		common.HexToAddress(coordinatorAddr),
		common.HexToAddress(linkAddr),
	)
	if err != nil {
		return &EthereumVRFConsumer{}, fmt.Errorf("VRFConsumer instance deployment have failed: %w", err)
	}

	consumer, err := solidity_vrf_consumer_interface.NewVRFConsumer(consumerDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFConsumer{}, fmt.Errorf("failed to instantiate VRFConsumer instance: %w", err)
	}

	return &EthereumVRFConsumer{
		client:   seth,
		consumer: consumer,
		address:  &consumerDeploymentData.Address,
	}, err
}

func DeployVRFMockETHLINKFeed(seth *seth.Client, answer *big.Int) (VRFMockETHLINKFeed, error) {
	abi, err := vrf_mock_ethlink_aggregator.VRFMockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFMockETHLINKAggregator{}, fmt.Errorf("failed to get VRFMockETHLINKAggregator ABI: %w", err)
	}

	deployment, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFMockETHLINKAggregator",
		*abi,
		common.FromHex(vrf_mock_ethlink_aggregator.VRFMockETHLINKAggregatorMetaData.Bin),
		answer,
	)
	if err != nil {
		return &EthereumVRFMockETHLINKAggregator{}, fmt.Errorf("VRFMockETHLINKAggregator deployment have failed: %w", err)
	}

	contract, err := vrf_mock_ethlink_aggregator.NewVRFMockETHLINKAggregator(deployment.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFMockETHLINKAggregator{}, fmt.Errorf("failed to instantiate VRFMockETHLINKAggregator instance: %w", err)
	}

	return &EthereumVRFMockETHLINKAggregator{
		client:   seth,
		contract: contract,
		address:  &deployment.Address,
	}, err
}

func (v *EthereumBlockhashStore) Address() string {
	return v.address.Hex()
}

func (v *EthereumBlockhashStore) GetBlockHash(ctx context.Context, blockNumber *big.Int) ([32]byte, error) {
	blockHash, err := v.blockHashStore.GetBlockhash(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, blockNumber)
	if err != nil {
		return [32]byte{}, err
	}
	return blockHash, nil
}

func (v *EthereumVRFCoordinator) Address() string {
	return v.address.Hex()
}

// HashOfKey get a hash of proving key to use it as a request ID part for VRF
func (v *EthereumVRFCoordinator) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	hash, err := v.coordinator.HashOfKey(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, pubKey)
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
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(v.client.NewTXOpts(), fee, common.HexToAddress(oracleAddr), publicProvingKey, jobID))
	return err
}

func (v *EthereumVRFConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFConsumer) Fund(_ *big.Float) error {
	panic("do not use this function, use actions_seth.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

// RequestRandomness requests VRF randomness
func (v *EthereumVRFConsumer) RequestRandomness(hash [32]byte, fee *big.Int) error {
	_, err := v.client.Decode(v.consumer.TestRequestRandomness(v.client.NewTXOpts(), hash, fee))
	return err
}

// CurrentRoundID helper roundID counter in consumer to check when all randomness requests are finished
func (v *EthereumVRFConsumer) CurrentRoundID(ctx context.Context) (*big.Int, error) {
	return v.consumer.CurrentRoundID(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

// RandomnessOutput get VRF randomness output
func (v *EthereumVRFConsumer) RandomnessOutput(ctx context.Context) (*big.Int, error) {
	return v.consumer.RandomnessOutput(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

// Fund sends specified currencies to the contract
func (v *EthereumVRF) Fund(_ *big.Float) error {
	panic("do not use this function, use actions_seth.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

// ProofLength returns the PROOFLENGTH call from the VRF contract
func (v *EthereumVRF) ProofLength(ctx context.Context) (*big.Int, error) {
	return v.vrf.PROOFLENGTH(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}
