package vrftesthelpers

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_consumer_interface_v08"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_request_id"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_request_id_v08"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

var (
	WeiPerUnitLink = decimal.RequireFromString("10000000000000000")
)

func GenerateProofResponseFromProof(p vrfkey.Proof, s proof.PreSeedData) (
	proof.MarshaledOnChainResponse, error) {
	return proof.GenerateProofResponseFromProof(p, s)
}

func CreateAndStartBHSJob(
	t *testing.T,
	fromAddresses []string,
	app *cltest.TestApplication,
	bhsAddress, coordinatorV1Address, coordinatorV2Address, coordinatorV2PlusAddress string,
	trustedBlockhashStoreAddress string, trustedBlockhashStoreBatchSize int32, lookback int,
	heartbeatPeriod time.Duration, waitBlocks int,
) job.Job {
	jid := uuid.New()
	s := testspecs.GenerateBlockhashStoreSpec(testspecs.BlockhashStoreSpecParams{
		JobID:                          jid.String(),
		Name:                           "blockhash-store",
		CoordinatorV1Address:           coordinatorV1Address,
		CoordinatorV2Address:           coordinatorV2Address,
		CoordinatorV2PlusAddress:       coordinatorV2PlusAddress,
		WaitBlocks:                     waitBlocks,
		LookbackBlocks:                 lookback,
		HeartbeatPeriod:                heartbeatPeriod,
		BlockhashStoreAddress:          bhsAddress,
		TrustedBlockhashStoreAddress:   trustedBlockhashStoreAddress,
		TrustedBlockhashStoreBatchSize: trustedBlockhashStoreBatchSize,
		PollPeriod:                     time.Second,
		RunTimeout:                     10 * time.Second,
		EVMChainID:                     1337,
		FromAddresses:                  fromAddresses,
	})
	jb, err := blockhashstore.ValidatedSpec(s.Toml())
	require.NoError(t, err)

	ctx := testutils.Context(t)
	require.NoError(t, app.JobSpawner().CreateJob(ctx, nil, &jb))
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			if jb.Type == job.BlockhashStore {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	return jb
}

func CreateAndStartBlockHeaderFeederJob(
	t *testing.T,
	fromAddresses []string,
	app *cltest.TestApplication,
	bhsAddress, batchBHSAddress, coordinatorV1Address, coordinatorV2Address, coordinatorV2PlusAddress string,
) job.Job {
	jid := uuid.New()
	s := testspecs.GenerateBlockHeaderFeederSpec(testspecs.BlockHeaderFeederSpecParams{
		JobID:                      jid.String(),
		Name:                       "block-header-feeder",
		CoordinatorV1Address:       coordinatorV1Address,
		CoordinatorV2Address:       coordinatorV2Address,
		CoordinatorV2PlusAddress:   coordinatorV2PlusAddress,
		WaitBlocks:                 256,
		LookbackBlocks:             1000,
		BlockhashStoreAddress:      bhsAddress,
		BatchBlockhashStoreAddress: batchBHSAddress,
		PollPeriod:                 15 * time.Second,
		RunTimeout:                 15 * time.Second,
		EVMChainID:                 1337,
		FromAddresses:              fromAddresses,
		GetBlockhashesBatchSize:    20,
		StoreBlockhashesBatchSize:  20,
	})
	jb, err := blockheaderfeeder.ValidatedSpec(s.Toml())
	require.NoError(t, err)

	ctx := testutils.Context(t)
	require.NoError(t, app.JobSpawner().CreateJob(ctx, nil, &jb))
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			if jb.Type == job.BlockHeaderFeeder {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	return jb
}

// CoordinatorUniverse represents the universe in which a randomness request occurs and
// is fulfilled.
type CoordinatorUniverse struct {
	// Golang wrappers ofr solidity contracts
	RootContract               *solidity_vrf_coordinator_interface.VRFCoordinator
	LinkContract               *link_token_interface.LinkToken
	BHSContract                *blockhash_store.BlockhashStore
	ConsumerContract           *solidity_vrf_consumer_interface.VRFConsumer
	RequestIDBase              *solidity_vrf_request_id.VRFRequestIDBaseTestHelper
	ConsumerContractV08        *solidity_vrf_consumer_interface_v08.VRFConsumer
	RequestIDBaseV08           *solidity_vrf_request_id_v08.VRFRequestIDBaseTestHelper
	RootContractAddress        common.Address
	ConsumerContractAddress    common.Address
	ConsumerContractAddressV08 common.Address
	LinkContractAddress        common.Address
	BHSContractAddress         common.Address

	// Abstraction representation of the ethereum blockchain
	Backend        *backends.SimulatedBackend
	CoordinatorABI *abi.ABI
	ConsumerABI    *abi.ABI
	// Cast of participants
	Sergey *bind.TransactOpts // Owns all the LINK initially
	Neil   *bind.TransactOpts // Node operator running VRF service
	Ned    *bind.TransactOpts // Secondary node operator
	Carol  *bind.TransactOpts // Author of consuming contract which requests randomness
}

var oneEth = big.NewInt(1000000000000000000) // 1e18 wei

func NewVRFCoordinatorUniverseWithV08Consumer(t *testing.T, key ethkey.KeyV2) CoordinatorUniverse {
	cu := NewVRFCoordinatorUniverse(t, key)
	consumerContractAddress, _, consumerContract, err :=
		solidity_vrf_consumer_interface_v08.DeployVRFConsumer(
			cu.Carol, cu.Backend, cu.RootContractAddress, cu.LinkContractAddress)
	require.NoError(t, err, "failed to deploy v08 VRFConsumer contract to simulated ethereum blockchain")
	_, _, requestIDBase, err :=
		solidity_vrf_request_id_v08.DeployVRFRequestIDBaseTestHelper(cu.Neil, cu.Backend)
	require.NoError(t, err, "failed to deploy v08 VRFRequestIDBaseTestHelper contract to simulated ethereum blockchain")
	cu.ConsumerContractAddressV08 = consumerContractAddress
	cu.RequestIDBaseV08 = requestIDBase
	cu.ConsumerContractV08 = consumerContract
	_, err = cu.LinkContract.Transfer(cu.Sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	cu.Backend.Commit()
	return cu
}

// newVRFCoordinatorUniverse sets up all identities and contracts associated with
// testing the solidity VRF contracts involved in randomness request workflow
func NewVRFCoordinatorUniverse(t *testing.T, keys ...ethkey.KeyV2) CoordinatorUniverse {
	var oracleTransactors []*bind.TransactOpts
	for _, key := range keys {
		oracleTransactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
		require.NoError(t, err)
		oracleTransactors = append(oracleTransactors, oracleTransactor)
	}

	var (
		sergey = testutils.MustNewSimTransactor(t)
		neil   = testutils.MustNewSimTransactor(t)
		ned    = testutils.MustNewSimTransactor(t)
		carol  = testutils.MustNewSimTransactor(t)
	)
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(1000).ToInt()},
		neil.From:   {Balance: assets.Ether(1000).ToInt()},
		ned.From:    {Balance: assets.Ether(1000).ToInt()},
		carol.From:  {Balance: assets.Ether(1000).ToInt()},
	}

	for _, t := range oracleTransactors {
		genesisData[t.From] = core.GenesisAccount{Balance: assets.Ether(1000).ToInt()}
	}

	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	consumerABI, err := abi.JSON(strings.NewReader(
		solidity_vrf_consumer_interface.VRFConsumerABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		solidity_vrf_coordinator_interface.VRFCoordinatorABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")
	bhsAddress, _, bhsContract, err := blockhash_store.DeployBlockhashStore(neil, backend)
	require.NoError(t, err, "failed to deploy BlockhashStore contract to simulated ethereum blockchain")
	coordinatorAddress, _, coordinatorContract, err :=
		solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			neil, backend, linkAddress, bhsAddress)
	require.NoError(t, err, "failed to deploy VRFCoordinator contract to simulated ethereum blockchain")
	consumerContractAddress, _, consumerContract, err :=
		solidity_vrf_consumer_interface.DeployVRFConsumer(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
	_, _, requestIDBase, err :=
		solidity_vrf_request_id.DeployVRFRequestIDBaseTestHelper(neil, backend)
	require.NoError(t, err, "failed to deploy VRFRequestIDBaseTestHelper contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	backend.Commit()
	return CoordinatorUniverse{
		RootContract:            coordinatorContract,
		RootContractAddress:     coordinatorAddress,
		LinkContract:            linkContract,
		LinkContractAddress:     linkAddress,
		BHSContract:             bhsContract,
		BHSContractAddress:      bhsAddress,
		ConsumerContract:        consumerContract,
		RequestIDBase:           requestIDBase,
		ConsumerContractAddress: consumerContractAddress,
		Backend:                 backend,
		CoordinatorABI:          &coordinatorABI,
		ConsumerABI:             &consumerABI,
		Sergey:                  sergey,
		Neil:                    neil,
		Ned:                     ned,
		Carol:                   carol,
	}
}
