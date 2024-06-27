package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_wrapper"
)

// EthereumBatchBlockhashStore represents BatchBlockhashStore contract
type EthereumBatchBlockhashStore struct {
	address             *common.Address
	client              blockchain.EVMClient
	batchBlockhashStore *batch_blockhash_store.BatchBlockhashStore
}

// EthereumBlockhashStore represents a blockhash store for VRF contract
type EthereumBlockhashStore struct {
	address        *common.Address
	client         blockchain.EVMClient
	blockHashStore *blockhash_store.BlockhashStore
}

// EthereumVRFConsumer represents VRF consumer contract
type EthereumVRFConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *solidity_vrf_consumer_interface.VRFConsumer
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

// EthereumVRF represents a VRF contract
type EthereumVRF struct {
	client  blockchain.EVMClient
	vrf     *solidity_vrf_wrapper.VRF
	address *common.Address
}

// DeployVRFContract deploy VRF contract
func (e *EthereumContractDeployer) DeployVRFContract() (VRF, error) {
	address, _, instance, err := e.client.DeployContract("VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return solidity_vrf_wrapper.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:  e.client,
		vrf:     instance.(*solidity_vrf_wrapper.VRF),
		address: address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func (e *EthereumContractDeployer) DeployBlockhashStore() (BlockHashStore, error) {
	address, _, instance, err := e.client.DeployContract("BlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return blockhash_store.DeployBlockhashStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBlockhashStore{
		client:         e.client,
		blockHashStore: instance.(*blockhash_store.BlockhashStore),
		address:        address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return solidity_vrf_coordinator_interface.DeployVRFCoordinator(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinator{
		client:      e.client,
		coordinator: instance.(*solidity_vrf_coordinator_interface.VRFCoordinator),
		address:     address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return solidity_vrf_consumer_interface.DeployVRFConsumer(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumer{
		client:   e.client,
		consumer: instance.(*solidity_vrf_consumer_interface.VRFConsumer),
		address:  address,
	}, err
}

func (v *EthereumBlockhashStore) Address() string {
	return v.address.Hex()
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

func (v *EthereumVRFConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFConsumer) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(ethereum.CallMsg{
		To: v.address,
	})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
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

// Fund sends specified currencies to the contract
func (v *EthereumVRF) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(ethereum.CallMsg{
		To: v.address,
	})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

// ProofLength returns the PROOFLENGTH call from the VRF contract
func (v *EthereumVRF) ProofLength(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return v.vrf.PROOFLENGTH(opts)
}

func (v *EthereumBatchBlockhashStore) Address() string {
	return v.address.Hex()
}
