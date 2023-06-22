package vrf_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmlogger "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2_upgradeable_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_malicious_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_single_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_proxy_admin"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_reverting_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_transparent_upgradeable_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// vrfConsumerContract is the common interface implemented by
// the example contracts used for the integration tests.
type vrfConsumerContract interface {
	CreateSubscriptionAndFund(opts *bind.TransactOpts, fundingJuels *big.Int) (*gethtypes.Transaction, error)
	SSubId(opts *bind.CallOpts) (uint64, error)
	SRequestId(opts *bind.CallOpts) (*big.Int, error)
	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*gethtypes.Transaction, error)
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)
}

type coordinatorV2Universe struct {
	// Golang wrappers of solidity contracts
	consumerContracts         []*vrf_consumer_v2.VRFConsumerV2
	consumerContractAddresses []common.Address

	vrfOwner        *vrf_owner.VRFOwner
	vrfOwnerAddress common.Address

	oldRootContract                    *vrf_coordinator_v2.VRFCoordinatorV2
	oldRootContractAddress             common.Address
	oldBatchCoordinatorContract        *batch_vrf_coordinator_v2.BatchVRFCoordinatorV2
	oldBatchCoordinatorContractAddress common.Address

	rootContract                     *vrf_coordinator_v2.VRFCoordinatorV2
	rootContractAddress              common.Address
	batchCoordinatorContract         *batch_vrf_coordinator_v2.BatchVRFCoordinatorV2
	batchCoordinatorContractAddress  common.Address
	linkContract                     *link_token_interface.LinkToken
	linkContractAddress              common.Address
	linkEthFeedAddress               common.Address
	bhsContract                      *blockhash_store.BlockhashStore
	bhsContractAddress               common.Address
	batchBHSContract                 *batch_blockhash_store.BatchBlockhashStore
	batchBHSContractAddress          common.Address
	maliciousConsumerContract        *vrf_malicious_consumer_v2.VRFMaliciousConsumerV2
	maliciousConsumerContractAddress common.Address
	revertingConsumerContract        *vrfv2_reverting_example.VRFV2RevertingExample
	revertingConsumerContractAddress common.Address
	// This is a VRFConsumerV2Upgradeable wrapper that points to the proxy address.
	consumerProxyContract        *vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExample
	consumerProxyContractAddress common.Address
	proxyAdminAddress            common.Address

	// Abstract representation of the ethereum blockchain
	backend        *backends.SimulatedBackend
	coordinatorABI *abi.ABI
	consumerABI    *abi.ABI

	// Cast of participants
	vrfConsumers []*bind.TransactOpts // Authors of consuming contracts that request randomness
	sergey       *bind.TransactOpts   // Owns all the LINK initially
	neil         *bind.TransactOpts   // Node operator running VRF service
	ned          *bind.TransactOpts   // Secondary node operator
	nallory      *bind.TransactOpts   // Oracle transactor
	evil         *bind.TransactOpts   // Author of a malicious consumer contract
	reverter     *bind.TransactOpts   // Author of always reverting contract
}

var (
	weiPerUnitLink = decimal.RequireFromString("10000000000000000")
)

func newVRFCoordinatorV2Universe(t *testing.T, key ethkey.KeyV2, numConsumers int) coordinatorV2Universe {
	testutils.SkipShort(t, "VRFCoordinatorV2Universe")
	oracleTransactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	require.NoError(t, err)
	var (
		sergey       = testutils.MustNewSimTransactor(t)
		neil         = testutils.MustNewSimTransactor(t)
		ned          = testutils.MustNewSimTransactor(t)
		evil         = testutils.MustNewSimTransactor(t)
		reverter     = testutils.MustNewSimTransactor(t)
		nallory      = oracleTransactor
		vrfConsumers []*bind.TransactOpts
	)

	// Create consumer contract deployer identities
	for i := 0; i < numConsumers; i++ {
		vrfConsumers = append(vrfConsumers, testutils.MustNewSimTransactor(t))
	}

	genesisData := core.GenesisAlloc{
		sergey.From:   {Balance: assets.Ether(1000).ToInt()},
		neil.From:     {Balance: assets.Ether(1000).ToInt()},
		ned.From:      {Balance: assets.Ether(1000).ToInt()},
		nallory.From:  {Balance: assets.Ether(1000).ToInt()},
		evil.From:     {Balance: assets.Ether(1000).ToInt()},
		reverter.From: {Balance: assets.Ether(1000).ToInt()},
	}
	for _, consumer := range vrfConsumers {
		genesisData[consumer.From] = core.GenesisAccount{
			Balance: assets.Ether(1000).ToInt(),
		}
	}

	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	consumerABI, err := abi.JSON(strings.NewReader(
		vrf_consumer_v2.VRFConsumerV2ABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		vrf_coordinator_v2.VRFCoordinatorV2ABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")

	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			evil, backend, 18, weiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)

	// Deploy blockhash store
	bhsAddress, _, bhsContract, err := blockhash_store.DeployBlockhashStore(neil, backend)
	require.NoError(t, err, "failed to deploy BlockhashStore contract to simulated ethereum blockchain")

	// Deploy batch blockhash store
	batchBHSAddress, _, batchBHSContract, err := batch_blockhash_store.DeployBatchBlockhashStore(neil, backend, bhsAddress)
	require.NoError(t, err, "failed to deploy BatchBlockhashStore contract to simulated ethereum blockchain")

	// Deploy VRF V2 coordinator
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			neil, backend, linkAddress, bhsAddress, linkEthFeed /* linkEth*/)
	require.NoError(t, err, "failed to deploy VRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy batch VRF V2 coordinator
	batchCoordinatorAddress, _, batchCoordinatorContract, err :=
		batch_vrf_coordinator_v2.DeployBatchVRFCoordinatorV2(
			neil, backend, coordinatorAddress,
		)
	require.NoError(t, err, "failed to deploy BatchVRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy old VRF v2 coordinator from bytecode
	err, oldRootContractAddress, oldRootContract := deployOldCoordinator(
		t, linkAddress, bhsAddress, linkEthFeed, backend, neil)

	// Deploy the VRFOwner contract, which will own the VRF coordinator
	// in some tests.
	// Don't transfer ownership now because it'll unnecessarily complicate
	// tests that don't really use this code path (which will be 99.9% of all
	// real-world use cases).
	vrfOwnerAddress, _, vrfOwner, err := vrf_owner.DeployVRFOwner(
		neil, backend, oldRootContractAddress,
	)
	require.NoError(t, err, "failed to deploy VRFOwner contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy batch VRF V2 coordinator
	oldBatchCoordinatorAddress, _, oldBatchCoordinatorContract, err :=
		batch_vrf_coordinator_v2.DeployBatchVRFCoordinatorV2(
			neil, backend, coordinatorAddress,
		)
	require.NoError(t, err, "failed to deploy BatchVRFCoordinatorV2 contract wrapping old vrf coordinator v2 to simulated ethereum blockchain")
	backend.Commit()

	// Create the VRF consumers.
	var (
		consumerContracts         []*vrf_consumer_v2.VRFConsumerV2
		consumerContractAddresses []common.Address
	)
	for _, author := range vrfConsumers {
		// Deploy a VRF consumer. It has a starting balance of 500 LINK.
		consumerContractAddress, _, consumerContract, err :=
			vrf_consumer_v2.DeployVRFConsumerV2(
				author, backend, coordinatorAddress, linkAddress)
		require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
		_, err = linkContract.Transfer(sergey, consumerContractAddress, assets.Ether(500).ToInt()) // Actually, LINK
		require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")

		consumerContracts = append(consumerContracts, consumerContract)
		consumerContractAddresses = append(consumerContractAddresses, consumerContractAddress)

		backend.Commit()
	}

	// Deploy malicious consumer with 1 link
	maliciousConsumerContractAddress, _, maliciousConsumerContract, err :=
		vrf_malicious_consumer_v2.DeployVRFMaliciousConsumerV2(
			evil, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFMaliciousConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, assets.Ether(1).ToInt()) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy upgradeable consumer, proxy, and proxy admin
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_upgradeable_example.DeployVRFConsumerV2UpgradeableExample(neil, backend)
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	backend.Commit()

	proxyAdminAddress, _, proxyAdmin, err := vrfv2_proxy_admin.DeployVRFV2ProxyAdmin(neil, backend)
	require.NoError(t, err)
	backend.Commit()

	// provide abi-encoded initialize function call on the implementation contract
	// so that it's called upon the proxy construction, to initialize it.
	upgradeableAbi, err := vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.GetAbi()
	require.NoError(t, err)
	initializeCalldata, err := upgradeableAbi.Pack("initialize", coordinatorAddress, linkAddress)
	hexified := hexutil.Encode(initializeCalldata)
	t.Log("initialize calldata:", hexified, "coordinator:", coordinatorAddress.String(), "link:", linkAddress)
	require.NoError(t, err)
	proxyAddress, _, _, err := vrfv2_transparent_upgradeable_proxy.DeployVRFV2TransparentUpgradeableProxy(
		neil, backend, upgradeableConsumerAddress, proxyAdminAddress, initializeCalldata)
	require.NoError(t, err)

	_, err = linkContract.Transfer(sergey, proxyAddress, assets.Ether(500).ToInt()) // Actually, LINK
	require.NoError(t, err)
	backend.Commit()

	implAddress, err := proxyAdmin.GetProxyImplementation(nil, proxyAddress)
	require.NoError(t, err)
	t.Log("impl address:", implAddress.String())
	require.Equal(t, upgradeableConsumerAddress, implAddress)

	proxiedConsumer, err := vrf_consumer_v2_upgradeable_example.NewVRFConsumerV2UpgradeableExample(
		proxyAddress, backend)
	require.NoError(t, err)

	cAddress, err := proxiedConsumer.COORDINATOR(nil)
	require.NoError(t, err)
	t.Log("coordinator address in proxy to upgradeable consumer:", cAddress.String())
	require.Equal(t, coordinatorAddress, cAddress)

	lAddress, err := proxiedConsumer.LINKTOKEN(nil)
	require.NoError(t, err)
	t.Log("link address in proxy to upgradeable consumer:", lAddress.String())
	require.Equal(t, linkAddress, lAddress)

	// Deploy always reverting consumer
	revertingConsumerContractAddress, _, revertingConsumerContract, err := vrfv2_reverting_example.DeployVRFV2RevertingExample(
		reverter, backend, coordinatorAddress, linkAddress,
	)
	require.NoError(t, err, "failed to deploy VRFRevertingExample contract to simulated eth blockchain")
	_, err = linkContract.Transfer(sergey, revertingConsumerContractAddress, assets.Ether(500).ToInt()) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFRevertingExample contract on simulated eth blockchain")
	backend.Commit()

	// Set the configuration on the coordinator.
	_, err = coordinatorContract.SetConfig(neil,
		uint16(1),                              // minRequestConfirmations
		uint32(2.5e6),                          // gas limit
		uint32(60*60*24),                       // stalenessSeconds
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(1e16),                       // 0.01 eth per link fallbackLinkPrice
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier2: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier3: uint32(100),
			FulfillmentFlatFeeLinkPPMTier4: uint32(10),
			FulfillmentFlatFeeLinkPPMTier5: uint32(1),
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40),
		},
	)
	require.NoError(t, err, "failed to set coordinator configuration")
	backend.Commit()

	// Set the configuration on the old coordinator.
	_, err = oldRootContract.SetConfig(neil,
		uint16(1),                              // minRequestConfirmations
		uint32(2.5e6),                          // gas limit
		uint32(60*60*24),                       // stalenessSeconds
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(1e16),                       // 0.01 eth per link fallbackLinkPrice
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier2: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier3: uint32(100),
			FulfillmentFlatFeeLinkPPMTier4: uint32(10),
			FulfillmentFlatFeeLinkPPMTier5: uint32(1),
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40),
		},
	)
	require.NoError(t, err, "failed to set old coordinator configuration")
	backend.Commit()

	return coordinatorV2Universe{
		vrfConsumers:              vrfConsumers,
		consumerContracts:         consumerContracts,
		consumerContractAddresses: consumerContractAddresses,

		batchCoordinatorContract:        batchCoordinatorContract,
		batchCoordinatorContractAddress: batchCoordinatorAddress,

		vrfOwner:        vrfOwner,
		vrfOwnerAddress: vrfOwnerAddress,

		oldRootContractAddress:             oldRootContractAddress,
		oldRootContract:                    oldRootContract,
		oldBatchCoordinatorContract:        oldBatchCoordinatorContract,
		oldBatchCoordinatorContractAddress: oldBatchCoordinatorAddress,

		revertingConsumerContract:        revertingConsumerContract,
		revertingConsumerContractAddress: revertingConsumerContractAddress,

		consumerProxyContract:        proxiedConsumer,
		consumerProxyContractAddress: proxiedConsumer.Address(),
		proxyAdminAddress:            proxyAdminAddress,

		rootContract:                     coordinatorContract,
		rootContractAddress:              coordinatorAddress,
		linkContract:                     linkContract,
		linkContractAddress:              linkAddress,
		linkEthFeedAddress:               linkEthFeed,
		bhsContract:                      bhsContract,
		bhsContractAddress:               bhsAddress,
		batchBHSContract:                 batchBHSContract,
		batchBHSContractAddress:          batchBHSAddress,
		maliciousConsumerContract:        maliciousConsumerContract,
		maliciousConsumerContractAddress: maliciousConsumerContractAddress,
		backend:                          backend,
		coordinatorABI:                   &coordinatorABI,
		consumerABI:                      &consumerABI,
		sergey:                           sergey,
		neil:                             neil,
		ned:                              ned,
		nallory:                          nallory,
		evil:                             evil,
		reverter:                         reverter,
	}
}

func deployOldCoordinator(
	t *testing.T,
	linkAddress common.Address,
	bhsAddress common.Address,
	linkEthFeed common.Address,
	backend *backends.SimulatedBackend,
	neil *bind.TransactOpts,
) (
	error,
	common.Address,
	*vrf_coordinator_v2.VRFCoordinatorV2,
) {
	bytecode := hexutil.MustDecode("0x60e06040523480156200001157600080fd5b506040516200608c3803806200608c8339810160408190526200003491620001b1565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e8565b5050506001600160601b0319606093841b811660805290831b811660a052911b1660c052620001fb565b6001600160a01b038116331415620001435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ac57600080fd5b919050565b600080600060608486031215620001c757600080fd5b620001d28462000194565b9250620001e26020850162000194565b9150620001f26040850162000194565b90509250925092565b60805160601c60a05160601c60c05160601c615e2762000265600039600081816105260152613bd901526000818161061d015261402401526000818161036d01528181611599015281816125960152818161302c0152818161318201526138360152615e276000f3fe608060405234801561001057600080fd5b506004361061025b5760003560e01c80636f64f03f11610145578063ad178361116100bd578063d2f9f9a71161008c578063e72f6e3011610071578063e72f6e30146106fa578063e82ad7d41461070d578063f2fde38b1461073057600080fd5b8063d2f9f9a7146106d4578063d7ae1d30146106e757600080fd5b8063ad17836114610618578063af198b971461063f578063c3f909d41461066f578063caf70c4a146106c157600080fd5b80638da5cb5b11610114578063a21a23e4116100f9578063a21a23e4146105da578063a47c7696146105e2578063a4c0ed361461060557600080fd5b80638da5cb5b146105a95780639f87fad7146105c757600080fd5b80636f64f03f146105685780637341c10c1461057b57806379ba50971461058e578063823597401461059657600080fd5b8063356dac71116101d85780635fbbc0d2116101a757806366316d8d1161018c57806366316d8d1461050e578063689c45171461052157806369bcdb7d1461054857600080fd5b80635fbbc0d21461040057806364d51a2a1461050657600080fd5b8063356dac71146103b457806340d6bb82146103bc5780634cb48a54146103da5780635d3b1d30146103ed57600080fd5b806308821d581161022f57806315c48b841161021457806315c48b841461030e578063181f5a77146103295780631b6b6d231461036857600080fd5b806308821d58146102cf57806312b58349146102e257600080fd5b80620122911461026057806302bcc5b61461028057806304c357cb1461029557806306bfa637146102a8575b600080fd5b610268610743565b60405161027793929190615964565b60405180910390f35b61029361028e366004615792565b6107bf565b005b6102936102a33660046157ad565b61086b565b60055467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610277565b6102936102dd3660046154a3565b610a60565b6005546801000000000000000090046bffffffffffffffffffffffff165b604051908152602001610277565b61031660c881565b60405161ffff9091168152602001610277565b604080518082018252601681527f565246436f6f7264696e61746f72563220312e302e30000000000000000000006020820152905161027791906158f1565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610277565b600a54610300565b6103c56101f481565b60405163ffffffff9091168152602001610277565b6102936103e836600461563c565b610c3f565b6103006103fb366004615516565b611036565b600c546040805163ffffffff80841682526401000000008404811660208301526801000000000000000084048116928201929092526c010000000000000000000000008304821660608201527001000000000000000000000000000000008304909116608082015262ffffff740100000000000000000000000000000000000000008304811660a0830152770100000000000000000000000000000000000000000000008304811660c08301527a0100000000000000000000000000000000000000000000000000008304811660e08301527d01000000000000000000000000000000000000000000000000000000000090920490911661010082015261012001610277565b610316606481565b61029361051c36600461545b565b611444565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b610300610556366004615779565b60009081526009602052604090205490565b6102936105763660046153a0565b6116ad565b6102936105893660046157ad565b6117f7565b610293611a85565b6102936105a4366004615792565b611b82565b60005473ffffffffffffffffffffffffffffffffffffffff1661038f565b6102936105d53660046157ad565b611d7c565b6102b66121fd565b6105f56105f0366004615792565b6123ed565b6040516102779493929190615b02565b6102936106133660046153d4565b612537565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b61065261064d366004615574565b6127a8565b6040516bffffffffffffffffffffffff9091168152602001610277565b600b546040805161ffff8316815263ffffffff6201000084048116602083015267010000000000000084048116928201929092526b010000000000000000000000909204166060820152608001610277565b6103006106cf3660046154bf565b612c6d565b6103c56106e2366004615792565b612c9d565b6102936106f53660046157ad565b612e92565b610293610708366004615385565b612ff3565b61072061071b366004615792565b613257565b6040519015158152602001610277565b61029361073e366004615385565b6134ae565b600b546007805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff169391928391908301828280156107ad57602002820191906000526020600020905b815481526020019060010190808311610799575b50505050509050925092509250909192565b6107c76134bf565b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1661082d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205461086890829073ffffffffffffffffffffffffffffffffffffffff16613542565b50565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806108d4576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610940576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600b546601000000000000900460ff1615610987576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff848116911614610a5a5767ffffffffffffffff841660008181526003602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b610a686134bf565b604080518082018252600091610a97919084906002908390839080828437600092019190915250612c6d915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1680610af9576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101839052602401610937565b600082815260066020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b600754811015610be9578260078281548110610b4c57610b4c615dbc565b90600052602060002001541415610bd7576007805460009190610b7190600190615c76565b81548110610b8157610b81615dbc565b906000526020600020015490508060078381548110610ba257610ba2615dbc565b6000918252602090912001556007805480610bbf57610bbf615d8d565b60019003818190600052602060002001600090559055505b80610be181615cba565b915050610b2e565b508073ffffffffffffffffffffffffffffffffffffffff167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610c3291815260200190565b60405180910390a2505050565b610c476134bf565b60c861ffff87161115610c9a576040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff871660048201819052602482015260c86044820152606401610937565b60008213610cd7576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610937565b6040805160a0808201835261ffff891680835263ffffffff89811660208086018290526000868801528a831660608088018290528b85166080988901819052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000001690971762010000909502949094177fffffffffffffffffffffffffffffffffff000000000000000000ffffffffffff166701000000000000009092027fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff16919091176b010000000000000000000000909302929092179093558651600c80549489015189890151938a0151978a0151968a015160c08b015160e08c01516101008d01519588167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009099169890981764010000000093881693909302929092177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000958716959095027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16949094176c0100000000000000000000000098861698909802979097177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000096909416959095027fffffffffffffffffff000000ffffffffffffffffffffffffffffffffffffffff16929092177401000000000000000000000000000000000000000062ffffff92831602177fffffff000000000000ffffffffffffffffffffffffffffffffffffffffffffff1677010000000000000000000000000000000000000000000000958216959095027fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff16949094177a01000000000000000000000000000000000000000000000000000092851692909202919091177cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167d0100000000000000000000000000000000000000000000000000000000009390911692909202919091178155600a84905590517fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2916110269189918991899189918991906159c3565b60405180910390a1505050505050565b600b546000906601000000000000900460ff1615611080576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff851660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff166110e6576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260026020908152604080832067ffffffffffffffff808a1685529252909120541680611156576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152336024820152604401610937565b600b5461ffff9081169086161080611172575060c861ffff8616115b156111c257600b546040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff8088166004830152909116602482015260c86044820152606401610937565b600b5463ffffffff620100009091048116908516111561122957600b546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff8087166004830152620100009092049091166024820152604401610937565b6101f463ffffffff8416111561127b576040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff841660048201526101f46024820152604401610937565b6000611288826001615bd2565b6040805160208082018c9052338284015267ffffffffffffffff808c16606084015284166080808401919091528351808403909101815260a08301845280519082012060c083018d905260e080840182905284518085039091018152610100909301909352815191012091925060009182916040805160208101849052439181019190915267ffffffffffffffff8c16606082015263ffffffff808b166080830152891660a08201523360c0820152919350915060e001604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012060008681526009835283902055848352820183905261ffff8a169082015263ffffffff808916606083015287166080820152339067ffffffffffffffff8b16908c907f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a97729060a00160405180910390a45033600090815260026020908152604080832067ffffffffffffffff808d16855292529091208054919093167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009091161790915591505095945050505050565b600b546601000000000000900460ff161561148b576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600860205260409020546bffffffffffffffffffffffff808316911610156114e5576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260086020526040812080548392906115129084906bffffffffffffffffffffffff16615c8d565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600560088282829054906101000a90046bffffffffffffffffffffffff166115699190615c8d565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b815260040161162192919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b602060405180830381600087803b15801561163b57600080fd5b505af115801561164f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061167391906154db565b6116a9576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b6116b56134bf565b6040805180820182526000916116e4919084906002908390839080828437600092019190915250612c6d915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1615611746576040517f4a0b8fa700000000000000000000000000000000000000000000000000000000815260048101829052602401610937565b600081815260066020908152604080832080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091556007805460018101825594527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610c32565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611860576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216146118c7576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610937565b600b546601000000000000900460ff161561190e576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206002015460641415611965576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416156119ac57610a5a565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260026020818152604080842067ffffffffffffffff8a1680865290835281852080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600384528286209094018054948501815585529382902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001685179055905192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610a51565b60015473ffffffffffffffffffffffffffffffffffffffff163314611b06576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610937565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600b546601000000000000900460ff1615611bc9576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16611c2f576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611cd15767ffffffffffffffff8116600090815260036020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610937565b67ffffffffffffffff81166000818152600360209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611de5576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611e4c576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610937565b600b546601000000000000900460ff1615611e93576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416611f2e576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff84166024820152604401610937565b67ffffffffffffffff8416600090815260036020908152604080832060020180548251818502810185019093528083529192909190830182828015611fa957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611f7e575b50505050509050600060018251611fc09190615c76565b905060005b825181101561215f578573ffffffffffffffffffffffffffffffffffffffff16838281518110611ff757611ff7615dbc565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16141561214d57600083838151811061202f5761202f615dbc565b6020026020010151905080600360008a67ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600201838154811061207557612075615dbc565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a1681526003909152604090206002018054806120ef576120ef615d8d565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061215f565b8061215781615cba565b915050611fc5565b5073ffffffffffffffffffffffffffffffffffffffff8516600081815260026020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b600b546000906601000000000000900460ff1615612247576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805467ffffffffffffffff1690600061226183615cf3565b82546101009290920a67ffffffffffffffff8181021990931691831602179091556005541690506000806040519080825280602002602001820160405280156122b4578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff888116808552600484528685209551865493516bffffffffffffffffffffffff9091167fffffffffffffffffffffffff0000000000000000000000000000000000000000948516176c010000000000000000000000009190931602919091179094558451606081018652338152808301848152818701888152958552600384529590932083518154831673ffffffffffffffffffffffffffffffffffffffff918216178255955160018201805490931696169590951790559151805194955090936123a592600285019201906150c5565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff81166000908152600360205260408120548190819060609073ffffffffffffffffffffffffffffffffffffffff1661245a576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff80861660009081526004602090815260408083205460038352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff8616966c010000000000000000000000009096049095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561252157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116124f6575b5050505050905093509350935093509193509193565b600b546601000000000000900460ff161561257e576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146125ed576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612627576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061263582840184615792565b67ffffffffffffffff811660009081526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1661269e576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260046020526040812080546bffffffffffffffffffffffff16918691906126d58385615bfe565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084600560088282829054906101000a90046bffffffffffffffffffffffff1661272c9190615bfe565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f88287846127939190615bba565b604080519283526020830191909152016121ed565b600b546000906601000000000000900460ff16156127f2576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a9050600080600061280687876139b5565b9250925092506000866060015163ffffffff1667ffffffffffffffff81111561283157612831615deb565b60405190808252806020026020018201604052801561285a578160200160208202803683370190505b50905060005b876060015163ffffffff168110156128ce5760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c8282815181106128b1576128b1615dbc565b6020908102919091010152806128c681615cba565b915050612860565b506000838152600960205260408082208290555181907f1fe543e300000000000000000000000000000000000000000000000000000000906129169087908690602401615ab4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff166601000000000000179055908a015160808b01519192506000916129e49163ffffffff169084613d04565b600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff1690556020808c01805167ffffffffffffffff9081166000908152600490935260408084205492518216845290922080549394506c01000000000000000000000000918290048316936001939192600c92612a68928692900416615bd2565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506000612abf8a600b600001600b9054906101000a900463ffffffff1663ffffffff16612ab985612c9d565b3a613d52565b6020808e015167ffffffffffffffff166000908152600490915260409020549091506bffffffffffffffffffffffff80831691161015612b2b576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020808d015167ffffffffffffffff1660009081526004909152604081208054839290612b679084906bffffffffffffffffffffffff16615c8d565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008b81526006602090815260408083205473ffffffffffffffffffffffffffffffffffffffff1683526008909152812080548594509092612bd091859116615bfe565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550877f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4888386604051612c53939291909283526bffffffffffffffffffffffff9190911660208301521515604082015260600190565b60405180910390a299505050505050505050505b92915050565b600081604051602001612c8091906158e3565b604051602081830303815290604052805190602001209050919050565b6040805161012081018252600c5463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301527001000000000000000000000000000000008104909216608082015262ffffff740100000000000000000000000000000000000000008304811660a08301819052770100000000000000000000000000000000000000000000008404821660c08401527a0100000000000000000000000000000000000000000000000000008404821660e08401527d0100000000000000000000000000000000000000000000000000000000009093041661010082015260009167ffffffffffffffff841611612dbb575192915050565b8267ffffffffffffffff168160a0015162ffffff16108015612df057508060c0015162ffffff168367ffffffffffffffff1611155b15612dff576020015192915050565b8267ffffffffffffffff168160c0015162ffffff16108015612e3457508060e0015162ffffff168367ffffffffffffffff1611155b15612e43576040015192915050565b8267ffffffffffffffff168160e0015162ffffff16108015612e79575080610100015162ffffff168367ffffffffffffffff1611155b15612e88576060015192915050565b6080015192915050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680612efb576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614612f62576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610937565b600b546601000000000000900460ff1615612fa9576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612fb284613257565b15612fe9576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5a8484613542565b612ffb6134bf565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561308357600080fd5b505afa158015613097573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130bb91906154fd565b6005549091506801000000000000000090046bffffffffffffffffffffffff168181111561311f576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610937565b818110156132525760006131338284615c76565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b1580156131c857600080fd5b505af11580156131dc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061320091906154db565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a1505b505050565b67ffffffffffffffff811660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff9081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561330657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116132db575b505050505081525050905060005b8160400151518110156134a45760005b60075481101561349157600061345a6007838154811061334657613346615dbc565b90600052602060002001548560400151858151811061336757613367615dbc565b602002602001015188600260008960400151898151811061338a5761338a615dbc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808f168352935220541660408051602080820187905273ffffffffffffffffffffffffffffffffffffffff959095168183015267ffffffffffffffff9384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b506000818152600960205260409020549091501561347e5750600195945050505050565b508061348981615cba565b915050613324565b508061349c81615cba565b915050613314565b5060009392505050565b6134b66134bf565b61086881613e5a565b60005473ffffffffffffffffffffffffffffffffffffffff163314613540576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610937565b565b600b546601000000000000900460ff1615613589576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff821660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff90811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561363457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613609575b5050509190925250505067ffffffffffffffff80851660009081526004602090815260408083208151808301909252546bffffffffffffffffffffffff81168083526c01000000000000000000000000909104909416918101919091529293505b83604001515181101561373b5760026000856040015183815181106136bc576136bc615dbc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff8a168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690558061373381615cba565b915050613695565b5067ffffffffffffffff8516600090815260036020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000009081168255600182018054909116905590613796600283018261514f565b505067ffffffffffffffff8516600090815260046020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055600580548291906008906138069084906801000000000000000090046bffffffffffffffffffffffff16615c8d565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85836bffffffffffffffffffffffff166040518363ffffffff1660e01b81526004016138be92919073ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b602060405180830381600087803b1580156138d857600080fd5b505af11580156138ec573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061391091906154db565b613946576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff861681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8716917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a25050505050565b60008060006139c78560000151612c6d565b60008181526006602052604090205490935073ffffffffffffffffffffffffffffffffffffffff1680613a29576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101859052602401610937565b6080860151604051613a48918691602001918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815281516020928301206000818152600990935291205490935080613ac5576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c01519251613b3e968b96909594910195865267ffffffffffffffff948516602087015292909316604085015263ffffffff908116606085015291909116608083015273ffffffffffffffffffffffffffffffffffffffff1660a082015260c00190565b604051602081830303815290604052805190602001208114613b8c576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b855167ffffffffffffffff164080613cb05786516040517fe9413d3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063e9413d389060240160206040518083038186803b158015613c3057600080fd5b505afa158015613c44573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c6891906154fd565b905080613cb05786516040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610937565b6000886080015182604051602001613cd2929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050613cf78982613f50565b9450505050509250925092565b60005a611388811015613d1657600080fd5b611388810390508460408204820311613d2e57600080fd5b50823b613d3a57600080fd5b60008083516020850160008789f190505b9392505050565b600080613d5d613fd9565b905060008113613d9c576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610937565b6000815a613daa8989615bba565b613db49190615c76565b613dc686670de0b6b3a7640000615c39565b613dd09190615c39565b613dda9190615c25565b90506000613df363ffffffff871664e8d4a51000615c39565b9050613e0b816b033b2e3c9fd0803ce8000000615c76565b821115613e44576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613e4e8183615bba565b98975050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415613eda576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610937565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613f848360000151846020015185604001518660600151868860a001518960c001518a60e001518b61010001516140ed565b60038360200151604051602001613f9c929190615aa0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b600b54604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600092670100000000000000900463ffffffff169182151591849182917f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b15801561407f57600080fd5b505afa158015614093573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140b791906157d7565b5094509092508491505080156140db57506140d28242615c76565b8463ffffffff16105b156140e55750600a545b949350505050565b6140f6896143c4565b61415c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610937565b614165886143c4565b6141cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610937565b6141d4836143c4565b61423a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610937565b614243826143c4565b6142a9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610937565b6142b5878a888761451f565b61431b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610937565b60006143278a876146c2565b9050600061433a898b878b868989614726565b9050600061434b838d8d8a866148ae565b9050808a146143b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610937565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11614451576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610937565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f116144de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610937565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096145188360005b602002015161490c565b1492915050565b600073ffffffffffffffffffffffffffffffffffffffff821661459e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610937565b6020840151600090600116156145b557601c6145b8565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa15801561466f573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b6146ca61516d565b6146f7600184846040516020016146e3939291906158c2565b604051602081830303815290604052614964565b90505b614703816143c4565b612c6757805160408051602081019290925261471f91016146e3565b90506146fa565b61472e61516d565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f90819006910614156147c1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610937565b6147cc8789886149cd565b614832576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610937565b61483d8486856149cd565b6148a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610937565b613e4e868484614b5a565b6000600286868685876040516020016148cc96959493929190615850565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b61496c61516d565b61497582614c89565b815261498a61498582600061450e565b614cde565b6020820181905260029006600114156149c8576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0390525b919050565b600082614a36576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610937565b83516020850151600090614a4c90600290615d1b565b15614a5857601c614a5b565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614adb573d6000803e3d6000fd5b505050602060405103519050600086604051602001614afa919061583e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b614b6261516d565b835160208086015185519186015160009384938493614b8393909190614d18565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114614c17576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610937565b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80614c5057614c50615d5e565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81106149c857604080516020808201939093528151808203840181529082019091528051910120614c91565b6000612c67826002614d117ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f6001615bba565b901c614eae565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a0890506000614dc083838585614fa2565b9098509050614dd188828e88614ffa565b9098509050614de288828c87614ffa565b90985090506000614df58d878b85614ffa565b9098509050614e0688828686614fa2565b9098509050614e1788828e89614ffa565b9098509050818114614e9a577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183099650614e9e565b8196505b5050505050509450945094915050565b600080614eb961518b565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a0820152614f056151a9565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa925082614f98576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610937565b5195945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b82805482825590600052602060002090810192821561513f579160200282015b8281111561513f57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906150e5565b5061514b9291506151c7565b5090565b508054600082559060005260206000209081019061086891906151c7565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b8082111561514b57600081556001016151c8565b803573ffffffffffffffffffffffffffffffffffffffff811681146149c857600080fd5b8060408101831015612c6757600080fd5b600082601f83011261522257600080fd5b6040516040810181811067ffffffffffffffff8211171561524557615245615deb565b806040525080838560408601111561525c57600080fd5b60005b600281101561527e57813583526020928301929091019060010161525f565b509195945050505050565b600060a0828403121561529b57600080fd5b60405160a0810181811067ffffffffffffffff821117156152be576152be615deb565b6040529050806152cd83615353565b81526152db60208401615353565b60208201526152ec6040840161533f565b60408201526152fd6060840161533f565b606082015261530e608084016151dc565b60808201525092915050565b803561ffff811681146149c857600080fd5b803562ffffff811681146149c857600080fd5b803563ffffffff811681146149c857600080fd5b803567ffffffffffffffff811681146149c857600080fd5b805169ffffffffffffffffffff811681146149c857600080fd5b60006020828403121561539757600080fd5b613d4b826151dc565b600080606083850312156153b357600080fd5b6153bc836151dc565b91506153cb8460208501615200565b90509250929050565b600080600080606085870312156153ea57600080fd5b6153f3856151dc565b935060208501359250604085013567ffffffffffffffff8082111561541757600080fd5b818701915087601f83011261542b57600080fd5b81358181111561543a57600080fd5b88602082850101111561544c57600080fd5b95989497505060200194505050565b6000806040838503121561546e57600080fd5b615477836151dc565b915060208301356bffffffffffffffffffffffff8116811461549857600080fd5b809150509250929050565b6000604082840312156154b557600080fd5b613d4b8383615200565b6000604082840312156154d157600080fd5b613d4b8383615211565b6000602082840312156154ed57600080fd5b81518015158114613d4b57600080fd5b60006020828403121561550f57600080fd5b5051919050565b600080600080600060a0868803121561552e57600080fd5b8535945061553e60208701615353565b935061554c6040870161531a565b925061555a6060870161533f565b91506155686080870161533f565b90509295509295909350565b60008082840361024081121561558957600080fd5b6101a08082121561559957600080fd5b6155a1615b90565b91506155ad8686615211565b82526155bc8660408701615211565b60208301526080850135604083015260a0850135606083015260c085013560808301526155eb60e086016151dc565b60a08301526101006155ff87828801615211565b60c0840152615612876101408801615211565b60e0840152610180860135818401525081935061563186828701615289565b925050509250929050565b6000806000806000808688036101c081121561565757600080fd5b6156608861531a565b965061566e6020890161533f565b955061567c6040890161533f565b945061568a6060890161533f565b935060808801359250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60830112156156c557600080fd5b6156cd615b90565b91506156db60a08a0161533f565b82526156e960c08a0161533f565b60208301526156fa60e08a0161533f565b604083015261010061570d818b0161533f565b606084015261571d828b0161533f565b608084015261572f6101408b0161532c565b60a08401526157416101608b0161532c565b60c08401526157536101808b0161532c565b60e08401526157656101a08b0161532c565b818401525050809150509295509295509295565b60006020828403121561578b57600080fd5b5035919050565b6000602082840312156157a457600080fd5b613d4b82615353565b600080604083850312156157c057600080fd5b6157c983615353565b91506153cb602084016151dc565b600080600080600060a086880312156157ef57600080fd5b6157f88661536b565b94506020860151935060408601519250606086015191506155686080870161536b565b8060005b6002811015610a5a57815184526020938401939091019060010161581f565b615848818361581b565b604001919050565b868152615860602082018761581b565b61586d606082018661581b565b61587a60a082018561581b565b61588760e082018461581b565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b8381526158d2602082018461581b565b606081019190915260800192915050565b60408101612c67828461581b565b600060208083528351808285015260005b8181101561591e57858101830151858201604001528201615902565b81811115615930576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b818110156159b557845183529383019391830191600101615999565b509098975050505050505050565b60006101c08201905061ffff8816825263ffffffff808816602084015280871660408401528086166060840152846080840152835481811660a0850152615a1760c08501838360201c1663ffffffff169052565b615a2e60e08501838360401c1663ffffffff169052565b615a466101008501838360601c1663ffffffff169052565b615a5e6101208501838360801c1663ffffffff169052565b62ffffff60a082901c811661014086015260b882901c811661016086015260d082901c1661018085015260e81c6101a090930192909252979650505050505050565b82815260608101613d4b602083018461581b565b6000604082018483526020604081850152818551808452606086019150828701935060005b81811015615af557845183529383019391830191600101615ad9565b5090979650505050505050565b6000608082016bffffffffffffffffffffffff87168352602067ffffffffffffffff87168185015273ffffffffffffffffffffffffffffffffffffffff80871660408601526080606086015282865180855260a087019150838801945060005b81811015615b80578551841683529484019491840191600101615b62565b50909a9950505050505050505050565b604051610120810167ffffffffffffffff81118282101715615bb457615bb4615deb565b60405290565b60008219821115615bcd57615bcd615d2f565b500190565b600067ffffffffffffffff808316818516808303821115615bf557615bf5615d2f565b01949350505050565b60006bffffffffffffffffffffffff808316818516808303821115615bf557615bf5615d2f565b600082615c3457615c34615d5e565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615c7157615c71615d2f565b500290565b600082821015615c8857615c88615d2f565b500390565b60006bffffffffffffffffffffffff83811690831681811015615cb257615cb2615d2f565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615cec57615cec615d2f565b5060010190565b600067ffffffffffffffff80831681811415615d1157615d11615d2f565b6001019392505050565b600082615d2a57615d2a615d5e565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a")
	ctorArgs, err := utils.ABIEncode(`[{"type":"address"}, {"type":"address"}, {"type":"address"}]`, linkAddress, bhsAddress, linkEthFeed)
	require.NoError(t, err)
	bytecode = append(bytecode, ctorArgs...)
	nonce, err := backend.PendingNonceAt(context.Background(), neil.From)
	require.NoError(t, err)
	gasPrice, err := backend.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	unsignedTx := gethtypes.NewContractCreation(nonce, big.NewInt(0), 15e6, gasPrice, bytecode)
	signedTx, err := neil.Signer(neil.From, unsignedTx)
	require.NoError(t, err)
	err = backend.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err, "could not deploy old vrf coordinator to simulated blockchain")
	backend.Commit()
	receipt, err := backend.TransactionReceipt(context.Background(), signedTx.Hash())
	require.NoError(t, err)
	oldRootContractAddress := receipt.ContractAddress
	require.NotEqual(t, common.HexToAddress("0x0"), oldRootContractAddress, "old vrf coordinator address equal to zero address, deployment failed")
	oldRootContract, err := vrf_coordinator_v2.NewVRFCoordinatorV2(oldRootContractAddress, backend)
	require.NoError(t, err, "could not create wrapper object for old vrf coordinator v2")
	return err, oldRootContractAddress, oldRootContract
}

// Send eth from prefunded account.
// Amount is number of ETH not wei.
func sendEth(t *testing.T, key ethkey.KeyV2, ec *backends.SimulatedBackend, to common.Address, eth int) {
	nonce, err := ec.PendingNonceAt(testutils.Context(t), key.Address)
	require.NoError(t, err)
	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: assets.GWei(10).ToInt(), // block base fee in sim
		Gas:       uint64(21_000),
		To:        &to,
		Value:     big.NewInt(0).Mul(big.NewInt(int64(eth)), big.NewInt(1e18)),
		Data:      nil,
	})
	signedTx, err := gethtypes.SignTx(tx, gethtypes.NewLondonSigner(big.NewInt(1337)), key.ToEcdsaPrivKey())
	require.NoError(t, err)
	err = ec.SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	ec.Commit()
}

func subscribeVRF(
	t *testing.T,
	author *bind.TransactOpts,
	consumerContract vrfConsumerContract,
	coordinatorContract vrf_coordinator_v2.VRFCoordinatorV2Interface,
	backend *backends.SimulatedBackend,
	fundingJuels *big.Int,
) (vrf_coordinator_v2.GetSubscription, uint64) {
	_, err := consumerContract.CreateSubscriptionAndFund(author, fundingJuels)
	require.NoError(t, err)
	backend.Commit()

	subID, err := consumerContract.SSubId(nil)
	require.NoError(t, err)

	sub, err := coordinatorContract.GetSubscription(nil, subID)
	require.NoError(t, err)
	return sub, subID
}

func createVRFJobs(
	t *testing.T,
	fromKeys [][]ethkey.KeyV2,
	app *cltest.TestApplication,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	uni coordinatorV2Universe,
	batchEnabled bool,
	gasLanePrices ...*assets.Wei,
) (jobs []job.Job) {
	if len(gasLanePrices) != len(fromKeys) {
		t.Fatalf("must provide one gas lane price for each set of from addresses. len(gasLanePrices) != len(fromKeys) [%d != %d]",
			len(gasLanePrices), len(fromKeys))
	}
	// Create separate jobs for each gas lane and register their keys
	for i, keys := range fromKeys {
		var keyStrs []string
		for _, k := range keys {
			keyStrs = append(keyStrs, k.Address.String())
		}

		vrfkey, err := app.GetKeyStore().VRF().Create()
		require.NoError(t, err)

		jid := uuid.New()
		incomingConfs := 2
		s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			JobID:                    jid.String(),
			Name:                     fmt.Sprintf("vrf-primary-%d", i),
			CoordinatorAddress:       coordinatorAddress.Hex(),
			BatchCoordinatorAddress:  batchCoordinatorAddress.Hex(),
			BatchFulfillmentEnabled:  batchEnabled,
			MinIncomingConfirmations: incomingConfs,
			PublicKey:                vrfkey.PublicKey.String(),
			FromAddresses:            keyStrs,
			BackoffInitialDelay:      10 * time.Millisecond,
			BackoffMaxDelay:          time.Second,
			V2:                       true,
			GasLanePrice:             gasLanePrices[i],
			VRFOwnerAddress:          uni.vrfOwnerAddress.Hex(),
		}).Toml()
		jb, err := vrf.ValidatedVRFSpec(s)
		t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfkey.PublicKey.MustHash())
		require.NoError(t, err)
		err = app.JobSpawner().CreateJob(&jb)
		require.NoError(t, err)
		registerProvingKeyHelper(t, uni, coordinator, vrfkey)
		jobs = append(jobs, jb)
	}
	// Wait until all jobs are active and listening for logs
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		var count int
		for _, jb := range jbs {
			if jb.Type == job.VRF {
				count++
			}
		}
		return count == len(fromKeys)
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	// Unfortunately the lb needs heads to be able to backfill logs to new subscribers.
	// To avoid confirming
	// TODO: it could just backfill immediately upon receiving a new subscriber? (though would
	// only be useful for tests, probably a more robust way is to have the job spawner accept a signal that a
	// job is fully up and running and not add it to the active jobs list before then)
	time.Sleep(2 * time.Second)

	return
}

func requestRandomnessForWrapper(
	t *testing.T,
	vrfWrapperConsumer vrfv2_wrapper_consumer_example.VRFV2WrapperConsumerExample,
	consumerOwner *bind.TransactOpts,
	keyHash common.Hash,
	subID uint64,
	numWords uint32,
	cbGasLimit uint32,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	uni coordinatorV2Universe,
	wrapperOverhead uint32,
) (*big.Int, uint64) {
	minRequestConfirmations := uint16(3)
	_, err := vrfWrapperConsumer.MakeRequest(
		consumerOwner,
		cbGasLimit,
		minRequestConfirmations,
		numWords,
	)
	require.NoError(t, err)
	uni.backend.Commit()

	iter, err := coordinator.FilterRandomWordsRequested(nil, nil, []uint64{subID}, nil)
	require.NoError(t, err, "could not filter RandomWordsRequested events")

	var events []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for iter.Next() {
		events = append(events, iter.Event)
	}

	wrapperIter, err := vrfWrapperConsumer.FilterWrapperRequestMade(nil, nil)
	require.NoError(t, err, "could not filter WrapperRequestMade events")

	wrapperConsumerEvents := []*vrfv2_wrapper_consumer_example.VRFV2WrapperConsumerExampleWrapperRequestMade{}
	for wrapperIter.Next() {
		wrapperConsumerEvents = append(wrapperConsumerEvents, wrapperIter.Event)
	}

	event := events[len(events)-1]
	wrapperConsumerEvent := wrapperConsumerEvents[len(wrapperConsumerEvents)-1]
	require.Equal(t, event.RequestId, wrapperConsumerEvent.RequestId, "request ID in consumer log does not match request ID in coordinator log")
	require.Equal(t, keyHash.Bytes(), event.KeyHash[:], "key hash of event (%s) and of request not equal (%s)", hex.EncodeToString(event.KeyHash[:]), keyHash.String())
	require.Equal(t, cbGasLimit+(cbGasLimit/63+1)+wrapperOverhead, event.CallbackGasLimit, "callback gas limit of event and of request not equal")
	require.Equal(t, minRequestConfirmations, event.MinimumRequestConfirmations, "min request confirmations of event and of request not equal")
	require.Equal(t, numWords, event.NumWords, "num words of event and of request not equal")

	return event.RequestId, event.Raw.BlockNumber
}

// requestRandomness requests randomness from the given vrf consumer contract
// and asserts that the request ID logged by the RandomWordsRequested event
// matches the request ID that is returned and set by the consumer contract.
// The request ID and request block number are then returned to the caller.
func requestRandomnessAndAssertRandomWordsRequestedEvent(
	t *testing.T,
	vrfConsumerHandle vrfConsumerContract,
	consumerOwner *bind.TransactOpts,
	keyHash common.Hash,
	subID uint64,
	numWords uint32,
	cbGasLimit uint32,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	uni coordinatorV2Universe,
) (requestID *big.Int, requestBlockNumber uint64) {
	minRequestConfirmations := uint16(2)
	_, err := vrfConsumerHandle.RequestRandomness(
		consumerOwner,
		keyHash,
		subID,
		minRequestConfirmations,
		cbGasLimit,
		numWords,
	)
	require.NoError(t, err)

	uni.backend.Commit()

	iter, err := coordinator.FilterRandomWordsRequested(nil, nil, []uint64{subID}, nil)
	require.NoError(t, err, "could not filter RandomWordsRequested events")

	var events []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for iter.Next() {
		events = append(events, iter.Event)
	}

	requestID, err = vrfConsumerHandle.SRequestId(nil)
	require.NoError(t, err)

	event := events[len(events)-1]
	require.Equal(t, event.RequestId, requestID, "request ID in contract does not match request ID in log")
	require.Equal(t, keyHash.Bytes(), event.KeyHash[:], "key hash of event (%s) and of request not equal (%s)", hex.EncodeToString(event.KeyHash[:]), keyHash.String())
	require.Equal(t, cbGasLimit, event.CallbackGasLimit, "callback gas limit of event and of request not equal")
	require.Equal(t, minRequestConfirmations, event.MinimumRequestConfirmations, "min request confirmations of event and of request not equal")
	require.Equal(t, numWords, event.NumWords, "num words of event and of request not equal")

	return requestID, event.Raw.BlockNumber
}

// subscribeAndAssertSubscriptionCreatedEvent subscribes the given consumer contract
// to VRF and funds the subscription with the given fundingJuels amount. It returns the
// subscription ID of the resulting subscription.
func subscribeAndAssertSubscriptionCreatedEvent(
	t *testing.T,
	vrfConsumerHandle vrfConsumerContract,
	consumerOwner *bind.TransactOpts,
	consumerContractAddress common.Address,
	fundingJuels *big.Int,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	uni coordinatorV2Universe,
) uint64 {
	// Create a subscription and fund with LINK.
	sub, subID := subscribeVRF(t, consumerOwner, vrfConsumerHandle, coordinator, uni.backend, fundingJuels)
	require.Equal(t, uint64(1), subID)
	require.Equal(t, fundingJuels.String(), sub.Balance.String())

	// Assert the subscription event in the coordinator contract.
	iter, err := coordinator.FilterSubscriptionCreated(nil, []uint64{subID})
	require.NoError(t, err)
	found := false
	for iter.Next() {
		if iter.Event.Owner != consumerContractAddress {
			require.FailNowf(t, "SubscriptionCreated event contains wrong owner address", "expected: %+v, actual: %+v", consumerContractAddress, iter.Event.Owner)
		} else {
			found = true
		}
	}
	require.True(t, found, "could not find SubscriptionCreated event for subID %d", subID)

	return subID
}

func assertRandomWordsFulfilled(
	t *testing.T,
	requestID *big.Int,
	expectedSuccess bool,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
) (rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
	// Check many times in case there are delays processing the event
	// this could happen occasionally and cause flaky tests.
	numChecks := 3
	found := false
	for i := 0; i < numChecks; i++ {
		filter, err := coordinator.FilterRandomWordsFulfilled(nil, []*big.Int{requestID})
		require.NoError(t, err)

		for filter.Next() {
			require.Equal(t, expectedSuccess, filter.Event.Success, "fulfillment event success not correct, expected: %+v, actual: %+v", expectedSuccess, filter.Event.Success)
			require.Equal(t, requestID, filter.Event.RequestId)
			found = true
			rwfe = filter.Event
		}

		if found {
			break
		}

		// Wait a bit and try again.
		time.Sleep(time.Second)
	}
	require.True(t, found, "RandomWordsFulfilled event not found")
	return
}

func assertNumRandomWords(
	t *testing.T,
	contract vrfConsumerContract,
	numWords uint32,
) {
	var err error
	for i := uint32(0); i < numWords; i++ {
		_, err = contract.SRandomWords(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
	}
}

func mine(t *testing.T, requestID *big.Int, subID uint64, uni coordinatorV2Universe, db *sqlx.DB) bool {
	return gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		var txs []txmgr.DbEthTx
		err := db.Select(&txs, `
		SELECT * FROM eth_txes
		WHERE eth_txes.state = 'confirmed'
			AND eth_txes.meta->>'RequestID' = $1
			AND CAST(eth_txes.meta->>'SubId' AS NUMERIC) = $2 LIMIT 1
		`, common.BytesToHash(requestID.Bytes()).String(), subID)
		require.NoError(t, err)
		t.Log("num txs", len(txs))
		return len(txs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func mineBatch(t *testing.T, requestIDs []*big.Int, subID uint64, uni coordinatorV2Universe, db *sqlx.DB) bool {
	requestIDMap := map[string]bool{}
	for _, requestID := range requestIDs {
		requestIDMap[common.BytesToHash(requestID.Bytes()).String()] = false
	}
	return gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		var txs []txmgr.DbEthTx
		err := db.Select(&txs, `
		SELECT * FROM eth_txes
		WHERE eth_txes.state = 'confirmed'
			AND CAST(eth_txes.meta->>'SubId' AS NUMERIC) = $1
		`, subID)
		require.NoError(t, err)
		for _, tx := range txs {
			var evmTx txmgr.Tx
			txmgr.DbEthTxToEthTx(tx, &evmTx)
			meta, err := evmTx.GetMeta()
			require.NoError(t, err)
			t.Log("meta:", meta)
			for _, requestID := range meta.RequestIDs {
				if _, ok := requestIDMap[requestID.String()]; ok {
					requestIDMap[requestID.String()] = true
				}
			}
		}
		foundAll := true
		for _, found := range requestIDMap {
			foundAll = foundAll && found
		}
		t.Log("requestIDMap:", requestIDMap)
		return foundAll
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func TestVRFV2Integration_SingleConsumer_ForceFulfillment(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerForcedFulfillment(
		t,
		ownerKey,
		uni,
		uni.oldRootContract,
		uni.oldRootContractAddress,
		uni.oldBatchCoordinatorContractAddress,
		false, // batchEnabled
	)
}

func TestVRFV2Integration_SingleConsumer_ForceFulfillment_BatchEnabled(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerForcedFulfillment(
		t,
		ownerKey,
		uni,
		uni.oldRootContract,
		uni.oldRootContractAddress,
		uni.oldBatchCoordinatorContractAddress,
		true, // batchEnabled
	)
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		5,     // number of requests to send
		false, // don't send big callback
	)
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment_BigGasCallback(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		5,    // number of requests to send
		true, // send big callback
	)
}

func TestVRFV2Integration_SingleConsumer_HappyPath(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPath(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress)
}

func TestVRFV2Integration_SingleConsumer_EOA_Request(t *testing.T) {
	t.Parallel()
	testEoa(t, false)
}

func TestVRFV2Integration_SingleConsumer_EOA_Request_Batching_Enabled(t *testing.T) {
	t.Parallel()
	testEoa(t, true)
}

func testEoa(t *testing.T, batchingEnabled bool) {
	gasLimit := int64(2_500_000)

	finalityDepth := uint32(50)

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_eoa_request", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.EVM[0].FinalityDepth = ptr(finalityDepth)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni)

	// Createa a new subscription.
	_, err := uni.rootContract.CreateSubscription(consumer)
	require.NoError(t, err)
	uni.backend.Commit()

	// Add the EOA as a consumer.
	_, err = uni.rootContract.AddConsumer(consumer, subID+1, consumer.From)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund the subscription with 1 LINK.
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID+1)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(uni.sergey, uni.rootContractAddress, big.NewInt(1e18), b)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		batchingEnabled,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make a randomness request with the EOA. This request is impossible to fulfill.
	numWords := uint32(1)
	minRequestConfirmations := uint16(2)
	_, err = uni.rootContract.RequestRandomWords(consumer, keyHash, subID+1, minRequestConfirmations, uint32(200_000), numWords)
	require.NoError(t, err)
	uni.backend.Commit()

	// Ensure request is not fulfilled.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Create query to fetch the application's log broadcasts.
	var broadcastsBeforeFinality []evmlogger.LogBroadcast
	var broadcastsAfterFinality []evmlogger.LogBroadcast
	query := `SELECT block_hash, consumed, log_index, job_id FROM log_broadcasts`
	q := pg.NewQ(app.GetSqlxDB(), app.Logger, app.Config.Database())

	// Execute the query.
	err = q.Select(&broadcastsBeforeFinality, query)
	require.NoError(t, err)

	// Ensure there is only one log broadcast (our EOA request), and that
	// it hasn't been marked as consumed yet.
	require.Equal(t, 1, len(broadcastsBeforeFinality))
	require.Equal(t, false, broadcastsBeforeFinality[0].Consumed)

	// Create new blocks until the finality depth has elapsed.
	for i := 0; i < int(finalityDepth); i++ {
		uni.backend.Commit()
	}

	// Ensure the request is still not fulfilled.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Execute the query for log broadcasts again after finality depth has elapsed.
	err = q.Select(&broadcastsAfterFinality, query)
	require.NoError(t, err)

	// Ensure that there is still only one log broadcast (our EOA request), but that
	// it has been marked as "consumed," such that it won't be retried.
	require.Equal(t, 1, len(broadcastsAfterFinality))
	require.Equal(t, true, broadcastsAfterFinality[0].Consumed)

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_EIP150_HappyPath(t *testing.T) {
	t.Parallel()
	callBackGasLimit := int64(2_500_000)            // base callback gas.
	eip150Fee := callBackGasLimit / 64              // premium needed for callWithExactGas
	coordinatorFulfillmentOverhead := int64(90_000) // fixed gas used in coordinator fulfillment
	gasLimit := callBackGasLimit + eip150Fee + coordinatorFulfillmentOverhead

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_eip150_happypath", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_EIP150_Revert(t *testing.T) {
	t.Parallel()
	callBackGasLimit := int64(2_500_000)            // base callback gas.
	eip150Fee := int64(0)                           // no premium given for callWithExactGas
	coordinatorFulfillmentOverhead := int64(90_000) // fixed gas used in coordinator fulfillment
	gasLimit := callBackGasLimit + eip150Fee + coordinatorFulfillmentOverhead

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_eip150_revert", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni)

	// Simulation should not pass.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func deployWrapper(t *testing.T, uni coordinatorV2Universe, wrapperOverhead uint32, coordinatorOverhead uint32, keyHash common.Hash) (
	wrapper *vrfv2_wrapper.VRFV2Wrapper,
	wrapperAddress common.Address,
	wrapperConsumer *vrfv2_wrapper_consumer_example.VRFV2WrapperConsumerExample,
	wrapperConsumerAddress common.Address,
) {
	wrapperAddress, _, wrapper, err := vrfv2_wrapper.DeployVRFV2Wrapper(uni.neil, uni.backend, uni.linkContractAddress, uni.linkEthFeedAddress, uni.rootContractAddress)
	require.NoError(t, err)
	uni.backend.Commit()

	_, err = wrapper.SetConfig(uni.neil, wrapperOverhead, coordinatorOverhead, 0, keyHash, 10)
	require.NoError(t, err)
	uni.backend.Commit()

	wrapperConsumerAddress, _, wrapperConsumer, err = vrfv2_wrapper_consumer_example.DeployVRFV2WrapperConsumerExample(uni.neil, uni.backend, uni.linkContractAddress, wrapperAddress)
	require.NoError(t, err)
	uni.backend.Commit()

	return
}

func TestVRFV2Integration_SingleConsumer_Wrapper(t *testing.T) {
	t.Parallel()
	wrapperOverhead := uint32(30_000)
	coordinatorOverhead := uint32(90_000)

	callBackGasLimit := int64(100_000) // base callback gas.
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_wrapper", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](3_500_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	wrapper, _, consumer, consumerAddress := deployWrapper(t, uni, wrapperOverhead, coordinatorOverhead, keyHash)

	// Fetch Subscription ID for Wrapper.
	wrapperSubID, err := wrapper.SUBSCRIPTIONID(nil)
	require.NoError(t, err)

	// Fund Subscription.
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, wrapperSubID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(uni.sergey, uni.rootContractAddress, assets.Ether(100).ToInt(), b)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund Consumer Contract.
	_, err = uni.linkContract.Transfer(uni.sergey, consumerAddress, assets.Ether(100).ToInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Make the first randomness request.
	numWords := uint32(1)
	requestID, _ := requestRandomnessForWrapper(t, *consumer, uni.neil, keyHash, wrapperSubID, numWords, uint32(callBackGasLimit), uni.rootContract, uni, wrapperOverhead)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, wrapperSubID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract)

	t.Log("Done!")
}

func TestVRFV2Integration_Wrapper_High_Gas(t *testing.T) {
	t.Parallel()
	wrapperOverhead := uint32(30_000)
	coordinatorOverhead := uint32(90_000)

	key1 := cltest.MustGenerateRandomKey(t)
	callBackGasLimit := int64(2_000_000) // base callback gas.
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_wrapper_high_gas_revert", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](3_500_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	wrapper, _, consumer, consumerAddress := deployWrapper(t, uni, wrapperOverhead, coordinatorOverhead, keyHash)

	// Fetch Subscription ID for Wrapper.
	wrapperSubID, err := wrapper.SUBSCRIPTIONID(nil)
	require.NoError(t, err)

	// Fund Subscription.
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, wrapperSubID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(uni.sergey, uni.rootContractAddress, assets.Ether(100).ToInt(), b)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund Consumer Contract.
	_, err = uni.linkContract.Transfer(uni.sergey, consumerAddress, assets.Ether(100).ToInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Make the first randomness request.
	numWords := uint32(1)
	requestID, _ := requestRandomnessForWrapper(t, *consumer, uni.neil, keyHash, wrapperSubID, numWords, uint32(callBackGasLimit), uni.rootContract, uni, wrapperOverhead)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, wrapperSubID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract)

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_NeedsBlockhashStore(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 2)
	testMultipleConsumersNeedBHS(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress)
}

func TestVRFV2Integration_SingleConsumer_BlockHeaderFeeder(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testBlockHeaderFeeder(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress)
}

func TestVRFV2Integration_SingleConsumer_NeedsTopUp(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerNeedsTopUp(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		assets.Ether(1).ToInt(),   // initial funding of 1 LINK
		assets.Ether(100).ToInt(), // top up of 100 LINK
	)
}

func TestVRFV2Integration_SingleConsumer_BigGasCallback_Sandwich(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(100)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_bigcallback_sandwich", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(100), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](5_000_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, assets.Ether(2).ToInt(), uni.rootContract, uni)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests, each one block apart, which contain a single low-gas request sandwiched between two high-gas requests.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	callbackGasLimits := []uint32{2_500_000, 50_000, 1_500_000}
	for _, limit := range callbackGasLimits {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni)
		reqIDs = append(reqIDs, requestID)
		uni.backend.Commit()
	}

	// Assert that we've completed 0 runs before adding 3 new requests.
	runs, err := app.PipelineORM().GetAllRuns()
	require.NoError(t, err)
	assert.Equal(t, 0, len(runs))
	assert.Equal(t, 3, len(reqIDs))

	// Wait for the 50_000 gas randomness request to be enqueued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// After the first successful request, no more will be enqueued.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, 3*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, reqIDs[1], subID, uni, db)

	// Assert the random word was fulfilled
	assertRandomWordsFulfilled(t, reqIDs[1], false, uni.rootContract)

	// Assert that we've still only completed 1 run before adding new requests.
	runs, err = app.PipelineORM().GetAllRuns()
	require.NoError(t, err)
	assert.Equal(t, 1, len(runs))

	// Make some randomness requests, each one block apart, this time without a low-gas request present in the callbackGasLimit slice.
	callbackGasLimits = []uint32{2_500_000, 2_500_000, 2_500_000}
	for _, limit := range callbackGasLimits {
		_, _ = requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni)
		uni.backend.Commit()
	}

	// Fulfillment will not be enqueued because subscriber doesn't have enough LINK for any of the requests.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_MultipleGasLanes(t *testing.T) {
	cheapKey := cltest.MustGenerateRandomKey(t)
	expensiveKey := cltest.MustGenerateRandomKey(t)
	cheapGasLane := assets.GWei(10)
	expensiveGasLane := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_multiplegaslanes", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Cheap gas lane.
			Key:          ptr(cheapKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: cheapGasLane},
		}, v2.KeySpecific{
			// Expensive gas lane.
			Key:          ptr(expensiveKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: expensiveGasLane},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, cheapKey, expensiveKey)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni)

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, cheapKey.Address, 10)
	sendEth(t, ownerKey, uni.backend, expensiveKey.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF jobs.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{cheapKey}, {expensiveKey}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		cheapGasLane, expensiveGasLane)
	cheapHash := jbs[0].VRFSpec.PublicKey.MustHash()
	expensiveHash := jbs[1].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	cheapRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, cheapHash, subID, numWords, 500_000, uni.rootContract, uni)

	// Wait for fulfillment to be queued for cheap key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, cheapRequestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, cheapRequestID, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	expensiveRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, expensiveHash, subID, numWords, 500_000, uni.rootContract, uni)

	// We should not have any new fulfillments until a top up.
	gomega.NewWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through. 100 LINK should do the trick.
	_, err := consumerContract.TopUpSubscription(consumer, decimal.RequireFromString("100e18").BigInt())
	require.NoError(t, err)

	// Wait for fulfillment to be queued for expensive key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, expensiveRequestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, expensiveRequestID, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_AlwaysRevertingCallback_StillFulfilled(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_alwaysrevertingcallback", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key)
	consumer := uni.reverter
	consumerContract := uni.revertingConsumerContract
	consumerContractAddress := uni.revertingConsumerContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the randomness request.
	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni.rootContract, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, false, uni.rootContract)
	t.Log("Done!")
}

func TestVRFV2Integration_ConsumerProxy_HappyPath(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_consumerproxy_happypath", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			Key:          ptr(key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)
	consumerOwner := uni.neil
	consumerContract := uni.consumerProxyContract
	consumerContractAddress := uni.consumerProxyContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(
		t, consumerContract, consumerOwner, consumerContractAddress,
		assets.Ether(5).ToInt(), uni.rootContract, uni)

	// Create gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1, key2}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID1, true, uni.rootContract)

	// Gas available will be around 724,385, which means that 750,000 - 724,385 = 25,615 gas was used.
	// This is ~20k more than what the non-proxied consumer uses.
	// So to be safe, users should probably over-estimate their fulfillment gas by ~25k.
	gasAvailable, err := consumerContract.SGasAvailable(nil)
	require.NoError(t, err)
	t.Log("gas available after proxied callback:", gasAvailable)

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni, db)
	assertRandomWordsFulfilled(t, requestID2, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.PendingNonceAt(testutils.Context(t), key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.PendingNonceAt(testutils.Context(t), key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}

func TestVRFV2Integration_ConsumerProxy_CoordinatorZeroAddress(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)

	// Deploy another upgradeable consumer, proxy, and proxy admin
	// to test vrfCoordinator != 0x0 condition.
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_upgradeable_example.DeployVRFConsumerV2UpgradeableExample(uni.neil, uni.backend)
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	uni.backend.Commit()

	// Deployment should revert if we give the 0x0 address for the coordinator.
	upgradeableAbi, err := vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.GetAbi()
	require.NoError(t, err)
	initializeCalldata, err := upgradeableAbi.Pack("initialize",
		common.BytesToAddress(common.LeftPadBytes([]byte{}, 20)), // zero address for the coordinator
		uni.linkContractAddress)
	require.NoError(t, err)
	_, _, _, err = vrfv2_transparent_upgradeable_proxy.DeployVRFV2TransparentUpgradeableProxy(
		uni.neil, uni.backend, upgradeableConsumerAddress, uni.proxyAdminAddress, initializeCalldata)
	require.Error(t, err)
}

func simulatedOverrides(t *testing.T, defaultGasPrice *assets.Wei, ks ...v2.KeySpecific) func(*chainlink.Config, *chainlink.Secrets) {
	return func(c *chainlink.Config, s *chainlink.Secrets) {
		require.Zero(t, testutils.SimulatedChainID.Cmp(c.EVM[0].ChainID.ToInt()))
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
		if defaultGasPrice != nil {
			c.EVM[0].GasEstimator.PriceDefault = defaultGasPrice
		}
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](2_000_000)

		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = models.MustNewDuration(0) // Head sampling disabled

		c.EVM[0].Transactions.ResendAfterThreshold = models.MustNewDuration(0)
		c.EVM[0].Transactions.ReaperThreshold = models.MustNewDuration(100 * time.Millisecond)

		c.EVM[0].FinalityDepth = ptr[uint32](15)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
		c.EVM[0].MinContractPayment = assets.NewLinkFromJuels(100)
		c.EVM[0].KeySpecific = ks
	}
}

func registerProvingKeyHelper(t *testing.T, uni coordinatorV2Universe, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, vrfkey vrfkey.KeyV2) {
	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = coordinator.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
}

func TestExternalOwnerConsumerExample(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	random := testutils.MustNewSimTransactor(t)
	genesisData := core.GenesisAlloc{
		owner.From:  {Balance: assets.Ether(10).ToInt()},
		random.From: {Balance: assets.Ether(10).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		owner, backend)
	require.NoError(t, err)
	backend.Commit()
	coordinatorAddress, _, coordinator, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner, backend, linkAddress, common.Address{}, common.Address{})
	require.NoError(t, err)
	_, err = coordinator.SetConfig(owner, uint16(1), uint32(10000), 1, 1, big.NewInt(10), vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: 0,
		FulfillmentFlatFeeLinkPPMTier2: 0,
		FulfillmentFlatFeeLinkPPMTier3: 0,
		FulfillmentFlatFeeLinkPPMTier4: 0,
		FulfillmentFlatFeeLinkPPMTier5: 0,
		ReqsForTier2:                   big.NewInt(0),
		ReqsForTier3:                   big.NewInt(0),
		ReqsForTier4:                   big.NewInt(0),
		ReqsForTier5:                   big.NewInt(0),
	})
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(owner, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err)
	backend.Commit()
	_, err = linkContract.Transfer(owner, consumerAddress, assets.Ether(2).ToInt())
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress}, []*big.Int{assets.Ether(999_999_998).ToInt(), assets.Ether(2).ToInt()})

	// Create sub, fund it and assign consumer
	_, err = coordinator.CreateSubscription(owner)
	require.NoError(t, err)
	backend.Commit()
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, uint64(1))
	require.NoError(t, err)
	_, err = linkContract.TransferAndCall(owner, coordinatorAddress, big.NewInt(0), b)
	require.NoError(t, err)
	_, err = coordinator.AddConsumer(owner, 1, consumerAddress)
	require.NoError(t, err)
	_, err = consumer.RequestRandomWords(random, 1, 1, 1, 1, [32]byte{})
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(owner, 1, 1, 1, 1, [32]byte{})
	require.NoError(t, err)

	// Reassign ownership, check that only new owner can request
	_, err = consumer.TransferOwnership(owner, random.From)
	require.NoError(t, err)
	_, err = consumer.RequestRandomWords(owner, 1, 1, 1, 1, [32]byte{})
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(random, 1, 1, 1, 1, [32]byte{})
	require.NoError(t, err)
}

func TestSimpleConsumerExample(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	random := testutils.MustNewSimTransactor(t)
	genesisData := core.GenesisAlloc{
		owner.From: {Balance: assets.Ether(10).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		owner, backend)
	require.NoError(t, err)
	backend.Commit()
	coordinatorAddress, _, _, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner, backend, linkAddress, common.Address{}, common.Address{})
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_single_consumer_example.DeployVRFSingleConsumerExample(owner, backend, coordinatorAddress, linkAddress, 1, 1, 1, [32]byte{})
	require.NoError(t, err)
	backend.Commit()
	_, err = linkContract.Transfer(owner, consumerAddress, assets.Ether(2).ToInt())
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress}, []*big.Int{assets.Ether(999_999_998).ToInt(), assets.Ether(2).ToInt()})
	_, err = consumer.TopUpSubscription(owner, assets.Ether(1).ToInt())
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(999_999_998).ToInt(), assets.Ether(1).ToInt(), assets.Ether(1).ToInt()})
	// Non-owner cannot withdraw
	_, err = consumer.Withdraw(random, assets.Ether(1).ToInt(), owner.From)
	require.Error(t, err)
	_, err = consumer.Withdraw(owner, assets.Ether(1).ToInt(), owner.From)
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(999_999_999).ToInt(), assets.Ether(0).ToInt(), assets.Ether(1).ToInt()})
	_, err = consumer.Unsubscribe(owner, owner.From)
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(1_000_000_000).ToInt(), assets.Ether(0).ToInt(), assets.Ether(0).ToInt()})
}

func TestIntegrationVRFV2(t *testing.T) {
	t.Parallel()
	// Reconfigure the sim chain with a default gas price of 1 gwei,
	// max gas limit of 2M and a key specific max 10 gwei price.
	// Keep the prices low so we can operate with small link balance subscriptions.
	gasPrice := assets.GWei(1)
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrf_v2_integration", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasPrice, v2.KeySpecific{
			Key:          &key.EIP55Address,
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	keys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Zero(t, key.Cmp(keys[0]))

	require.NoError(t, app.Start(testutils.Context(t)))

	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Create and fund a subscription.
	// We should see that our subscription has 1 link.
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
	}, []*big.Int{
		assets.Ether(500).ToInt(), // 500 link
		big.NewInt(0),             // 0 link
	})
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = carolContract.CreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		assets.Ether(499).ToInt(),
		assets.Ether(1).ToInt(),
		big.NewInt(0),
	})
	subId, err := carolContract.SSubId(nil)
	require.NoError(t, err)
	subStart, err := uni.rootContract.GetSubscription(nil, subId)
	require.NoError(t, err)

	// Make a request for random words.
	// By requesting 500k callback with a configured eth gas limit default of 500k,
	// we ensure that the job is indeed adjusting the gaslimit to suit the users request.
	gasRequested := 500_000
	nw := 10
	requestedIncomingConfs := 3
	_, err = carolContract.RequestRandomness(carol, keyHash, subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
	require.NoError(t, err)

	// Oracle tries to withdraw before its fulfilled should fail
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, big.NewInt(1000))
	require.Error(t, err)

	for i := 0; i < requestedIncomingConfs; i++ {
		uni.backend.Commit()
	}

	// We expect the request to be serviced
	// by the node.
	var runs []pipeline.Run
	gomega.NewWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		// It is possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However, we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		uni.backend.Commit()
		return len(runs) == 1 && runs[0].State == pipeline.RunStatusCompleted
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Wait for the request to be fulfilled on-chain.
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	gomega.NewWithT(t).Eventually(func() bool {
		rfIterator, err2 := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
		require.NoError(t, err2, "failed to logs")
		uni.backend.Commit()
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, testutils.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())
	assert.True(t, rf[0].Success, "expected callback to succeed")
	fulfillReceipt, err := uni.backend.TransactionReceipt(testutils.Context(t), rf[0].Raw.TxHash)
	require.NoError(t, err)

	// Assert all the random words received by the consumer are different and non-zero.
	seen := make(map[string]struct{})
	var rw *big.Int
	for i := 0; i < nw; i++ {
		rw, err = carolContract.SRandomWords(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
		_, ok := seen[rw.String()]
		assert.False(t, ok)
		seen[rw.String()] = struct{}{}
	}

	// We should have exactly as much gas as we requested
	// after accounting for function look up code, argument decoding etc.
	// which should be fixed in this test.
	ga, err := carolContract.SGasAvailable(nil)
	require.NoError(t, err)
	gaDecoding := big.NewInt(0).Add(ga, big.NewInt(3701))
	assert.Equal(t, 0, gaDecoding.Cmp(big.NewInt(int64(gasRequested))), "expected gas available %v to exceed gas requested %v", gaDecoding, gasRequested)
	t.Log("gas available", ga.String())

	// Assert that we were only charged for how much gas we actually used.
	// We should be charged for the verification + our callbacks execution in link.
	subEnd, err := uni.rootContract.GetSubscription(nil, subId)
	require.NoError(t, err)
	var (
		end   = decimal.RequireFromString(subEnd.Balance.String())
		start = decimal.RequireFromString(subStart.Balance.String())
		wei   = decimal.RequireFromString("1000000000000000000")
		gwei  = decimal.RequireFromString("1000000000")
	)
	t.Log("end balance", end)
	linkWeiCharged := start.Sub(end)
	// Remove flat fee of 0.001 to get fee for just gas.
	linkCharged := linkWeiCharged.Sub(decimal.RequireFromString("1000000000000000")).Div(wei)
	gasPriceD := decimal.NewFromBigInt(gasPrice.ToInt(), 0)
	t.Logf("subscription charged %s with gas prices of %s gwei and %s ETH per LINK\n", linkCharged, gasPriceD.Div(gwei), weiPerUnitLink.Div(wei))
	expected := decimal.RequireFromString(strconv.Itoa(int(fulfillReceipt.GasUsed))).Mul(gasPriceD).Div(weiPerUnitLink)
	t.Logf("expected sub charge gas use %v %v off by %v", fulfillReceipt.GasUsed, expected, expected.Sub(linkCharged))
	// The expected sub charge should be within 200 gas of the actual gas usage.
	// wei/link * link / wei/gas = wei / (wei/gas) = gas
	gasDiff := linkCharged.Sub(expected).Mul(weiPerUnitLink).Div(gasPriceD).Abs().IntPart()
	t.Log("gasDiff", gasDiff)
	assert.Less(t, gasDiff, int64(200))

	// If the oracle tries to withdraw more than it was paid it should fail.
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, linkWeiCharged.Add(decimal.NewFromInt(1)).BigInt())
	require.Error(t, err)

	// Assert the oracle can withdraw its payment.
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, linkWeiCharged.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		assets.Ether(499).ToInt(),
		subFunding.Sub(linkWeiCharged).BigInt(),
		linkWeiCharged.BigInt(),
	})

	// We should see the response count present
	chain, err := app.Chains.EVM.Get(big.NewInt(1337))
	require.NoError(t, err)

	q := pg.NewQ(app.GetSqlxDB(), app.Logger, app.Config.Database())
	counts := vrf.GetStartingResponseCountsV2(q, app.Logger, chain.Client().ConfiguredChainID().Uint64(), chain.Config().EVM().FinalityDepth())
	t.Log(counts, rf[0].RequestId.String())
	assert.Equal(t, uint64(1), counts[rf[0].RequestId.String()])
}

func TestMaliciousConsumer(t *testing.T) {
	t.Parallel()
	config, _ := heavyweight.FullTestDBV2(t, "vrf_v2_integration_malicious", func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](2_000_000)
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(1)
		c.EVM[0].GasEstimator.PriceDefault = assets.GWei(1)
		c.EVM[0].GasEstimator.FeeCapDefault = assets.GWei(1)
	})
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	err := app.GetKeyStore().Unlock(cltest.Password)
	require.NoError(t, err)
	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)

	jid := uuid.New()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		FromAddresses:            []string{key.Address.String()},
		CoordinatorAddress:       uni.rootContractAddress.String(),
		BatchCoordinatorAddress:  uni.batchCoordinatorContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		GasLanePrice:             assets.GWei(1),
		PublicKey:                vrfkey.PublicKey.String(),
		V2:                       true,
	}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	err = app.JobSpawner().CreateJob(&jb)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)

	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)

	_, err = uni.maliciousConsumerContract.SetKeyHash(carol,
		vrfkey.PublicKey.MustHash())
	require.NoError(t, err)
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.maliciousConsumerContract.CreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Send a re-entrant request
	_, err = uni.maliciousConsumerContract.RequestRandomness(carol)
	require.NoError(t, err)

	// We expect the request to be serviced
	// by the node.
	var attempts []txmgr.TxAttempt
	gomega.NewWithT(t).Eventually(func() bool {
		//runs, err = app.PipelineORM().GetAllRuns()
		attempts, _, err = app.TxmStorageService().TxAttempts(0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		t.Log("attempts", attempts)
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].Tx.State == txmgrcommon.TxConfirmed
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// The fulfillment tx should succeed
	ch, err := app.GetChains().EVM.Default()
	require.NoError(t, err)
	r, err := ch.Client().TransactionReceipt(testutils.Context(t), attempts[0].Hash)
	require.NoError(t, err)
	require.Equal(t, uint64(1), r.Status)

	// The user callback should have errored
	it, err := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
	require.NoError(t, err)
	var fulfillments []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	for it.Next() {
		fulfillments = append(fulfillments, it.Event)
	}
	require.Equal(t, 1, len(fulfillments))
	require.Equal(t, false, fulfillments[0].Success)

	// It should not have succeeded in placing another request.
	it2, err2 := uni.rootContract.FilterRandomWordsRequested(nil, nil, nil, nil)
	require.NoError(t, err2)
	var requests []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for it2.Next() {
		requests = append(requests, it2.Event)
	}
	require.Equal(t, 1, len(requests))
}

func TestRequestCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()

	t.Run("non-proxied consumer", func(tt *testing.T) {
		carol := uni.vrfConsumers[0]
		carolContract := uni.consumerContracts[0]
		carolContractAddress := uni.consumerContractAddresses[0]

		_, err = carolContract.CreateSubscriptionAndFund(carol,
			big.NewInt(1000000000000000000)) // 0.1 LINK
		require.NoError(tt, err)
		uni.backend.Commit()
		subId, err := carolContract.SSubId(nil)
		require.NoError(tt, err)
		// Ensure even with large number of consumers its still cheap
		var addrs []common.Address
		for i := 0; i < 99; i++ {
			addrs = append(addrs, testutils.NewAddress())
		}
		_, err = carolContract.UpdateSubscription(carol, addrs)
		require.NoError(tt, err)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			carolContractAddress, uni.consumerABI,
			"requestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
		tt.Log("gas estimate of non-proxied testRequestRandomness:", estimate)
		// V2 should be at least (87000-134000)/134000 = 35% cheaper
		// Note that a second call drops further to 68998 gas, but would also drop in V1.
		assert.Less(tt, estimate, uint64(90_000),
			"requestRandomness tx gas cost more than expected")
	})

	t.Run("proxied consumer", func(tt *testing.T) {
		consumerOwner := uni.neil
		consumerContract := uni.consumerProxyContract
		consumerContractAddress := uni.consumerProxyContractAddress

		// Create a subscription and fund with 5 LINK.
		tx, err := consumerContract.CreateSubscriptionAndFund(consumerOwner, assets.Ether(5).ToInt())
		require.NoError(tt, err)
		uni.backend.Commit()
		r, err := uni.backend.TransactionReceipt(testutils.Context(t), tx.Hash())
		require.NoError(tt, err)
		t.Log("gas used by proxied CreateSubscriptionAndFund:", r.GasUsed)

		subId, err := consumerContract.SSubId(nil)
		require.NoError(tt, err)
		_, err = uni.rootContract.GetSubscription(nil, subId)
		require.NoError(tt, err)

		// Ensure even with large number of consumers it's still cheap
		var addrs []common.Address
		for i := 0; i < 99; i++ {
			addrs = append(addrs, testutils.NewAddress())
		}
		_, err = consumerContract.UpdateSubscription(consumerOwner, addrs)

		theAbi := evmtypes.MustGetABI(vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.ABI)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			consumerContractAddress, &theAbi,
			"requestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
		tt.Log("gas estimate of proxied requestRandomness:", estimate)
		// There is some gas overhead of the delegatecall that is made by the proxy
		// to the logic contract. See https://www.evm.codes/#f4?fork=grayGlacier for a detailed
		// breakdown of the gas costs of a delegatecall.
		assert.Less(tt, estimate, uint64(96_000),
			"proxied testRequestRandomness tx gas cost more than expected")
	})
}

func TestMaxConsumersCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))
	_, err := carolContract.CreateSubscriptionAndFund(carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := carolContract.SSubId(nil)
	require.NoError(t, err)
	var addrs []common.Address
	for i := 0; i < 98; i++ {
		addrs = append(addrs, testutils.NewAddress())
	}
	_, err = carolContract.UpdateSubscription(carol, addrs)
	// Ensure even with max number of consumers its still reasonable gas costs.
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"removeConsumer", subId, carolContractAddress)
	t.Log(estimate)
	assert.Less(t, estimate, uint64(310000))
	estimate = estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"addConsumer", subId, testutils.NewAddress())
	t.Log(estimate)
	assert.Less(t, estimate, uint64(100000))
}

func TestFulfillmentCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()

	var (
		nonProxiedConsumerGasEstimate uint64
		proxiedConsumerGasEstimate    uint64
	)
	t.Run("non-proxied consumer", func(tt *testing.T) {
		carol := uni.vrfConsumers[0]
		carolContract := uni.consumerContracts[0]
		carolContractAddress := uni.consumerContractAddresses[0]

		_, err = carolContract.CreateSubscriptionAndFund(carol,
			big.NewInt(1000000000000000000)) // 0.1 LINK
		require.NoError(tt, err)
		uni.backend.Commit()
		subId, err := carolContract.SSubId(nil)
		require.NoError(tt, err)

		gasRequested := 50_000
		nw := 1
		requestedIncomingConfs := 3
		_, err = carolContract.RequestRandomness(carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
		require.NoError(t, err)
		for i := 0; i < requestedIncomingConfs; i++ {
			uni.backend.Commit()
		}

		requestLog := FindLatestRandomnessRequestedLog(tt, uni.rootContract, vrfkey.PublicKey.MustHash())
		s, err := proof.BigToSeed(requestLog.PreSeed)
		require.NoError(t, err)
		proof, rc, err := proof.GenerateProofResponseV2(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
			PreSeed:          s,
			BlockHash:        requestLog.Raw.BlockHash,
			BlockNum:         requestLog.Raw.BlockNumber,
			SubId:            subId,
			CallbackGasLimit: uint32(gasRequested),
			NumWords:         uint32(nw),
			Sender:           carolContractAddress,
		})
		require.NoError(tt, err)
		nonProxiedConsumerGasEstimate = estimateGas(tt, uni.backend, common.Address{},
			uni.rootContractAddress, uni.coordinatorABI,
			"fulfillRandomWords", proof, rc)
		t.Log("non-proxied consumer fulfillment gas estimate:", nonProxiedConsumerGasEstimate)
		// Establish very rough bounds on fulfillment cost
		assert.Greater(tt, nonProxiedConsumerGasEstimate, uint64(120_000))
		assert.Less(tt, nonProxiedConsumerGasEstimate, uint64(500_000))
	})

	t.Run("proxied consumer", func(tt *testing.T) {
		consumerOwner := uni.neil
		consumerContract := uni.consumerProxyContract
		consumerContractAddress := uni.consumerProxyContractAddress

		_, err = consumerContract.CreateSubscriptionAndFund(consumerOwner, assets.Ether(5).ToInt())
		require.NoError(t, err)
		uni.backend.Commit()
		subId, err := consumerContract.SSubId(nil)
		require.NoError(t, err)
		gasRequested := 50_000
		nw := 1
		requestedIncomingConfs := 3
		_, err = consumerContract.RequestRandomness(consumerOwner, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
		require.NoError(t, err)
		for i := 0; i < requestedIncomingConfs; i++ {
			uni.backend.Commit()
		}

		requestLog := FindLatestRandomnessRequestedLog(t, uni.rootContract, vrfkey.PublicKey.MustHash())
		require.Equal(tt, subId, requestLog.SubId)
		s, err := proof.BigToSeed(requestLog.PreSeed)
		require.NoError(t, err)
		proof, rc, err := proof.GenerateProofResponseV2(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
			PreSeed:          s,
			BlockHash:        requestLog.Raw.BlockHash,
			BlockNum:         requestLog.Raw.BlockNumber,
			SubId:            subId,
			CallbackGasLimit: uint32(gasRequested),
			NumWords:         uint32(nw),
			Sender:           consumerContractAddress,
		})
		require.NoError(t, err)
		proxiedConsumerGasEstimate = estimateGas(t, uni.backend, common.Address{},
			uni.rootContractAddress, uni.coordinatorABI,
			"fulfillRandomWords", proof, rc)
		t.Log("proxied consumer fulfillment gas estimate", proxiedConsumerGasEstimate)
		// Establish very rough bounds on fulfillment cost
		assert.Greater(t, proxiedConsumerGasEstimate, uint64(120_000))
		assert.Less(t, proxiedConsumerGasEstimate, uint64(500_000))
	})
}

func TestStartingCountsV1(t *testing.T) {
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, "vrf_test_starting_counts", nil)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES (1337, NOW(), NOW())`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO evm_heads (hash, number, parent_hash, created_at, timestamp, evm_chain_id)
	VALUES ($1, 4, $2, NOW(), NOW(), 1337)`, utils.NewHash(), utils.NewHash())
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	q := pg.NewQ(db, lggr, cfg.Database())
	finalityDepth := 3
	counts := vrf.GetStartingResponseCountsV1(q, lggr, 1337, uint32(finalityDepth))
	assert.Equal(t, 0, len(counts))
	ks := keystore.New(db, utils.FastScryptParams, lggr, cfg.Database())
	err = ks.Unlock(testutils.Password)
	require.NoError(t, err)
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := evmtypes.Nonce(0), evmtypes.Nonce(1), evmtypes.Nonce(2), evmtypes.Nonce(3)
	reqID := utils.PadByteToHash(0x10)
	m1 := txmgr.TxMeta{
		RequestID: &reqID,
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	md1_ := datatypes.JSON(md1)
	reqID2 := utils.PadByteToHash(0x11)
	m2 := txmgr.TxMeta{
		RequestID: &reqID2,
	}
	md2, err := json.Marshal(&m2)
	md2_ := datatypes.JSON(md2)
	require.NoError(t, err)
	chainID := utils.NewBig(big.NewInt(1337))
	confirmedTxes := []txmgr.Tx{
		{
			Sequence:           &n1,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &datatypes.JSON{},
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		{
			Sequence:           &n2,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md1_,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		{
			Sequence:           &n3,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md2_,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		{
			Sequence:           &n4,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md2_,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
	}
	// add unconfirmed txes
	unconfirmedTxes := []txmgr.Tx{}
	for i := int64(4); i < 6; i++ {
		reqID3 := utils.PadByteToHash(0x12)
		md, err := json.Marshal(&txmgr.TxMeta{
			RequestID: &reqID3,
		})
		require.NoError(t, err)
		md1 := datatypes.JSON(md)
		newNonce := evmtypes.Nonce(i + 1)
		unconfirmedTxes = append(unconfirmedTxes, txmgr.Tx{
			Sequence:           &newNonce,
			FromAddress:        k.Address,
			Error:              null.String{},
			CreatedAt:          b,
			State:              txmgrcommon.TxUnconfirmed,
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			Meta:               &md1,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		})
	}
	txes := append(confirmedTxes, unconfirmedTxes...)
	sql := `INSERT INTO eth_txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, broadcast_at, initial_broadcast_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id)
VALUES (:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :state, :created_at, :broadcast_at, :initial_broadcast_at, :meta, :subject, :evm_chain_id, :min_confirmations, :pipeline_task_run_id);`
	for _, tx := range txes {
		dbEtx := txmgr.DbEthTxFromEthTx(&tx)
		_, err = db.NamedExec(sql, &dbEtx)
		txmgr.DbEthTxToEthTx(dbEtx, &tx)
		require.NoError(t, err)
	}

	// add eth_tx_attempts for confirmed
	broadcastBlock := int64(1)
	txAttempts := []txmgr.TxAttempt{}
	for i := range confirmedTxes {
		txAttempts = append(txAttempts, txmgr.TxAttempt{
			TxID:                    int64(i + 1),
			TxFee:                   gas.EvmFee{Legacy: assets.NewWeiI(100)},
			SignedRawTx:             []byte(`blah`),
			Hash:                    utils.NewHash(),
			BroadcastBeforeBlockNum: &broadcastBlock,
			State:                   txmgrtypes.TxAttemptBroadcast,
			CreatedAt:               time.Now(),
			ChainSpecificFeeLimit:   uint32(100),
		})
	}
	// add eth_tx_attempts for unconfirmed
	for i := range unconfirmedTxes {
		txAttempts = append(txAttempts, txmgr.TxAttempt{
			TxID:                  int64(i + 1 + len(confirmedTxes)),
			TxFee:                 gas.EvmFee{Legacy: assets.NewWeiI(100)},
			SignedRawTx:           []byte(`blah`),
			Hash:                  utils.NewHash(),
			State:                 txmgrtypes.TxAttemptInProgress,
			CreatedAt:             time.Now(),
			ChainSpecificFeeLimit: uint32(100),
		})
	}
	for _, txAttempt := range txAttempts {
		t.Log("tx attempt eth tx id: ", txAttempt.TxID)
	}
	sql = `INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, state, created_at, chain_specific_gas_limit)
		VALUES (:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :state, :created_at, :chain_specific_gas_limit)`
	for _, attempt := range txAttempts {
		dbAttempt := txmgr.DbEthTxAttemptFromEthTxAttempt(&attempt)
		_, err = db.NamedExec(sql, &dbAttempt)
		txmgr.DbEthTxAttemptToEthTxAttempt(dbAttempt, &attempt)
		require.NoError(t, err)
	}

	// add eth_receipts
	receipts := []txmgr.Receipt{}
	for i := 0; i < 4; i++ {
		receipts = append(receipts, txmgr.Receipt{
			BlockHash:        utils.NewHash(),
			TxHash:           txAttempts[i].Hash,
			BlockNumber:      broadcastBlock,
			TransactionIndex: 1,
			Receipt:          evmtypes.Receipt{},
			CreatedAt:        time.Now(),
		})
	}
	sql = `INSERT INTO eth_receipts (block_hash, tx_hash, block_number, transaction_index, receipt, created_at)
		VALUES (:block_hash, :tx_hash, :block_number, :transaction_index, :receipt, :created_at)`
	for _, r := range receipts {
		_, err := db.NamedExec(sql, &r)
		require.NoError(t, err)
	}

	counts = vrf.GetStartingResponseCountsV1(q, lggr, 1337, uint32(finalityDepth))
	assert.Equal(t, 3, len(counts))
	assert.Equal(t, uint64(1), counts[utils.PadByteToHash(0x10)])
	assert.Equal(t, uint64(2), counts[utils.PadByteToHash(0x11)])
	assert.Equal(t, uint64(2), counts[utils.PadByteToHash(0x12)])

	countsV2 := vrf.GetStartingResponseCountsV2(q, lggr, 1337, uint32(finalityDepth))
	t.Log(countsV2)
	assert.Equal(t, 3, len(countsV2))
	assert.Equal(t, uint64(1), countsV2[big.NewInt(0x10).String()])
	assert.Equal(t, uint64(2), countsV2[big.NewInt(0x11).String()])
	assert.Equal(t, uint64(2), countsV2[big.NewInt(0x12).String()])
}

func FindLatestRandomnessRequestedLog(t *testing.T,
	coordContract *vrf_coordinator_v2.VRFCoordinatorV2,
	keyHash [32]byte) *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested {
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	gomega.NewWithT(t).Eventually(func() bool {
		rfIterator, err2 := coordContract.FilterRandomWordsRequested(nil, [][32]byte{keyHash}, nil, []common.Address{})
		require.NoError(t, err2, "failed to logs")
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) >= 1
	}, testutils.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())
	latest := len(rf) - 1
	return rf[latest]
}

func AssertLinkBalances(t *testing.T, linkContract *link_token_interface.LinkToken, addresses []common.Address, balances []*big.Int) {
	require.Equal(t, len(addresses), len(balances))
	for i, a := range addresses {
		b, err := linkContract.BalanceOf(nil, a)
		require.NoError(t, err)
		assert.Equal(t, balances[i].String(), b.String(), "invalid balance for %v", a)
	}
}

func ptr[T any](t T) *T { return &t }
