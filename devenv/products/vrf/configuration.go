package vrf

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "vrf"}).Logger()

type Configurator struct {
	Config []*VRF `toml:"vrf"`
}

type VRF struct {
	GasSettings *products.GasSettings `toml:"gas_settings"`
	Out         *Out                  `toml:"out"`
}

type Out struct {
	ConsumerAddress     string `toml:"consumer_address"`
	CoordinatorAddress  string `toml:"coordinator_address"`
	KeyHash             string `toml:"key_hash"`
	JobID               string `toml:"job_id"`
	PublicKeyCompressed string `toml:"public_key_compressed"`
	ExternalJobID       string `toml:"external_job_id"`
	ChainID             string `toml:"chain_id"`
}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.Config = cfg.Config
	return nil
}

func (m *Configurator) Store(path string, instanceIdx int) error {
	if err := products.Store(".", m); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateNodesConfig(
	ctx context.Context,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) (string, error) {
	return products.DefaultLegacyCLNodeConfig(bc)
}

func (m *Configurator) GenerateNodesSecrets(
	_ context.Context,
	_ *fake.Input,
	_ []*blockchain.Input,
	_ []*nodeset.Input,
) (string, error) {
	return "", nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	instanceIdx int,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) error {
	L.Info().Msg("Connecting to CL nodes")
	cls, err := clclient.New(ns[0].Out.CLNodes)
	if err != nil {
		return fmt.Errorf("failed to connect to CL nodes: %w", err)
	}

	c, auth, rootAddr, err := products.ETHClient(ctx, bc[0].Out.Nodes[0].ExternalWSUrl, m.Config[0].GasSettings.FeeCapMultiplier, m.Config[0].GasSettings.TipCapMultiplier)
	if err != nil {
		return fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	// Deploy Link Token
	linkAddr, linkTx, lt, err := link_token.DeployLinkToken(auth, c)
	if err != nil {
		return fmt.Errorf("could not deploy link token contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, linkTx)
	if err != nil {
		return err
	}
	L.Info().Str("Address", linkAddr.Hex()).Msg("Deployed link token contract")

	tx, err := lt.GrantMintRole(auth, common.HexToAddress(rootAddr))
	if err != nil {
		return fmt.Errorf("could not grant mint role: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return err
	}

	// Deploy BlockHashStore
	bhsAddr, bhsTx, _, err := blockhash_store.DeployBlockhashStore(auth, c)
	if err != nil {
		return fmt.Errorf("could not deploy blockhash store: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, bhsTx)
	if err != nil {
		return err
	}
	L.Info().Str("Address", bhsAddr.Hex()).Msg("Deployed blockhash store")

	// Deploy VRF Coordinator
	coordAddr, coordTx, coordinator, err := solidity_vrf_coordinator_interface.DeployVRFCoordinator(auth, c, linkAddr, bhsAddr)
	if err != nil {
		return fmt.Errorf("could not deploy VRF coordinator: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, coordTx)
	if err != nil {
		return err
	}
	L.Info().Str("Address", coordAddr.Hex()).Msg("Deployed VRF coordinator")

	// Deploy VRF Consumer
	consumerAddr, consumerTx, _, err := solidity_vrf_consumer_interface.DeployVRFConsumer(auth, c, coordAddr, linkAddr)
	if err != nil {
		return fmt.Errorf("could not deploy VRF consumer: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, consumerTx)
	if err != nil {
		return err
	}
	L.Info().Str("Address", consumerAddr.Hex()).Msg("Deployed VRF consumer")

	// Deploy VRF v1 library
	_, vrfLibTx, _, err := solidity_vrf_wrapper.DeployVRF(auth, c)
	if err != nil {
		return fmt.Errorf("could not deploy VRF library: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, vrfLibTx)
	if err != nil {
		return err
	}
	L.Info().Msg("Deployed VRF v1 library")

	// Mint LINK to consumer
	L.Info().Msgf("Minting LINK for consumer address: %s", consumerAddr)
	tx, err = lt.Mint(auth, consumerAddr, big.NewInt(2e18))
	if err != nil {
		return fmt.Errorf("could not mint link to consumer: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return err
	}

	// Fund CL node transmitter
	transmitters := make([]common.Address, 0)
	for _, nc := range cls {
		addr, cErr := nc.ReadPrimaryETHKey(bc[0].Out.ChainID)
		if cErr != nil {
			return cErr
		}
		transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
	}
	pkey := products.NetworkPrivateKey()
	if pkey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}
	for _, addr := range transmitters {
		if cErr := products.FundAddressEIP1559(ctx, c, pkey, addr.String(), 10); cErr != nil {
			return cErr
		}
	}

	// Create VRF key on node
	vrfKey, err := cls[0].MustCreateVRFKey()
	if err != nil {
		return fmt.Errorf("could not create VRF key: %w", err)
	}
	L.Info().Interface("Key", vrfKey).Msg("Created VRF proving key")
	pubKeyCompressed := vrfKey.Data.ID

	// Create VRF job
	jobUUID := uuid.New()
	pipelineSpec := &clclient.VRFTxPipelineSpec{
		Address: coordAddr.String(),
	}
	observationSource, err := pipelineSpec.String()
	if err != nil {
		return fmt.Errorf("could not build VRF pipeline spec: %w", err)
	}

	job, err := cls[0].MustCreateJob(&clclient.VRFJobSpec{
		Name:                     fmt.Sprintf("vrf-%s", jobUUID),
		CoordinatorAddress:       coordAddr.String(),
		MinIncomingConfirmations: 1,
		PublicKey:                pubKeyCompressed,
		ExternalJobID:            jobUUID.String(),
		EVMChainID:               bc[0].ChainID,
		ObservationSource:        observationSource,
	})
	if err != nil {
		return fmt.Errorf("could not create VRF job: %w", err)
	}
	L.Info().Str("JobID", job.Data.ID).Msg("Created VRF job")

	// Register proving key on coordinator
	oracleAddr := transmitters[0]
	provingKey, err := encodeOnChainVRFProvingKey(vrfKey)
	if err != nil {
		return fmt.Errorf("could not encode VRF proving key: %w", err)
	}
	jobIDBytes := encodeOnChainExternalJobID(jobUUID)

	tx, err = coordinator.RegisterProvingKey(auth, big.NewInt(1), oracleAddr, provingKey, jobIDBytes)
	if err != nil {
		return fmt.Errorf("could not register proving key: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return err
	}
	L.Info().Msg("Registered VRF proving key on coordinator")

	// Compute key hash
	keyHash, err := coordinator.HashOfKey(&bind.CallOpts{Context: ctx}, provingKey)
	if err != nil {
		return fmt.Errorf("could not compute key hash: %w", err)
	}
	L.Info().Str("KeyHash", hex.EncodeToString(keyHash[:])).Msg("Computed key hash")

	m.Config[0].Out = &Out{
		ConsumerAddress:     consumerAddr.String(),
		CoordinatorAddress:  coordAddr.String(),
		KeyHash:             hex.EncodeToString(keyHash[:]),
		JobID:               job.Data.ID,
		PublicKeyCompressed: pubKeyCompressed,
		ExternalJobID:       jobUUID.String(),
		ChainID:             bc[0].ChainID,
	}
	return nil
}

func encodeOnChainVRFProvingKey(vrfKey *clclient.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Data.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var ok bool
	provingKey[0], ok = new(big.Int).SetString(uncompressed[2:66], 16)
	if !ok {
		return [2]*big.Int{}, errors.New("cannot convert first half of VRF key to *big.Int")
	}
	provingKey[1], ok = new(big.Int).SetString(uncompressed[66:], 16)
	if !ok {
		return [2]*big.Int{}, errors.New("cannot convert second half of VRF key to *big.Int")
	}
	return provingKey, nil
}

func encodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}
