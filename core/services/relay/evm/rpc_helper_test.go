package evm_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"log"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	keytypes "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	_ "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting" //nolint common practice to import test mods with .
)

const sepolia_network_name = "sepolia"
const sepolia_chain_id = 11155111
const sepolia_ws_endpoint = "wss://ethereum-sepolia-rpc.publicnode.com"

// const sepolia_ws_endpoint = "wss://rpc.sepolia.org/"
// const sepolia_http_endpoint = "https://rpc.sepolia.org/"
const sepolia_http_endpoint = "https://ethereum-sepolia-rpc.publicnode.com"

// TODO move to make it reusable
type RpcHelper struct {
	//sim         *backends.SimulatedBackend
	accounts     []*bind.TransactOpts
	deployerKey  *ecdsa.PrivateKey
	senderKey    *ecdsa.PrivateKey
	txm          evmtxmgr.TxManager
	client       client.Client
	db           *sqlx.DB
	gasEstimator gas.EvmFeeEstimator

	//TODO improve, to be set during creation
	DeployPrivateKey string
	SenderPrivateKey string
}

func (h *RpcHelper) Commit() {
}

func (h RpcHelper) GasEstimator() *gas.EvmFeeEstimator {
	return &h.gasEstimator
}

func (h *RpcHelper) Backend() bind.ContractBackend {
	return h.client
}

func (h *RpcHelper) ChainReaderEVMClient(ctx context.Context, t *testing.T, ht logpoller.HeadTracker, conf evmtypes.ChainReaderConfig) client.Client {
	return h.client
}

func (h *RpcHelper) WrappedChainWriter(cw types.ChainWriter, client client.Client) types.ChainWriter {
	return cw
}

func (h *RpcHelper) Init(t *testing.T) {
	var err error
	h.deployerKey, err = createPrivateKey(h.DeployPrivateKey)
	require.NoError(t, err)
	h.senderKey, err = createPrivateKey(h.SenderPrivateKey)
	require.NoError(t, err)

	h.accounts = h.Accounts(t)

	h.db = pgtest.NewSqlxDB(t)

	h.client = h.Client(t)

	h.gasEstimator = *getEvmFeeEstimator(t, err, h)
	h.txm = h.TXM(t, h.client)
}

func createPrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {

	// Decode the hex string into a byte slice
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// Create a new big.Int from the private key bytes
	privKeyInt := new(big.Int).SetBytes(privKeyBytes)

	// Reconstruct the ecdsa.PrivateKey using the secp256k1 curve
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256() // Use the secp256k1 curve
	privKey.D = privKeyInt
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKeyBytes)

	// Print out the private and public keys
	fmt.Printf("Private Key: %x\n", privKey.D)
	fmt.Printf("Public Key X: %x\n", privKey.PublicKey.X)
	fmt.Printf("Public Key Y: %x\n", privKey.PublicKey.Y)

	// Alternatively, you can use the go-ethereum crypto library to generate the key
	privKeyECDSA, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, err
	}
	return privKeyECDSA, nil
}

func createChainClient(t *testing.T) client.Client {
	noNewHeadsThreshold := 3 * time.Minute
	selectionMode := ptr("HighestHead")
	leaseDuration := 0 * time.Second
	pollFailureThreshold := ptr(uint32(5))
	pollInterval := 10 * time.Second
	syncThreshold := ptr(uint32(5))
	nodeIsSyncingEnabled := ptr(false)
	chainTypeStr := ""
	finalizedBlockOffset := ptr[uint32](16)
	enforceRepeatableRead := ptr(true)
	deathDeclarationDelay := time.Second * 3
	noNewFinalizedBlocksThreshold := time.Second * 5
	finalizedBlockPollInterval := time.Second * 4
	nodeConfigs := []client.NodeConfig{
		{
			Name:    ptr(sepolia_network_name),
			WSURL:   ptr(sepolia_ws_endpoint),
			HTTPURL: ptr(sepolia_http_endpoint),
		},
	}
	finalityDepth := ptr(uint32(10))
	finalityTagEnabled := ptr(true)
	chainCfg, nodePool, nodes, err := client.NewClientConfigs(selectionMode, leaseDuration, chainTypeStr, nodeConfigs,
		pollFailureThreshold, pollInterval, syncThreshold, nodeIsSyncingEnabled, noNewHeadsThreshold, finalityDepth,
		finalityTagEnabled, finalizedBlockOffset, enforceRepeatableRead, deathDeclarationDelay, noNewFinalizedBlocksThreshold, finalizedBlockPollInterval)
	require.NoError(t, err)

	evmClient := client.NewEvmClient(nodePool, chainCfg, nil, logger.Test(t), big.NewInt(sepolia_chain_id), nodes, chaintype.ChainType(chainTypeStr))
	//evmClient.Dial(testutils.Context(t))
	return evmClient
}

func (h *RpcHelper) Accounts(t *testing.T) []*bind.TransactOpts {
	if h.accounts != nil {
		return h.accounts
	}
	deployer, err := bind.NewKeyedTransactorWithChainID(h.deployerKey, big.NewInt(sepolia_chain_id))
	require.NoError(t, err)

	sender, err := bind.NewKeyedTransactorWithChainID(h.senderKey, big.NewInt(sepolia_chain_id))
	require.NoError(t, err)

	return []*bind.TransactOpts{deployer, sender}
}

func (h *RpcHelper) MustGenerateRandomKey(t *testing.T) ethkey.KeyV2 {
	return cltest.MustGenerateRandomKey(t)
}

func (h *RpcHelper) GasPriceBufferPercent() int64 {
	return 0
}

//func (h *RpcHelper) Backend() bind.ContractBackend {
//	if h.sim == nil {
//		h.sim = backends.NewSimulatedBackend(
//			core.GenesisAlloc{h.accounts[0].From: {Balance: big.NewInt(math.MaxInt64)}, h.accounts[1].From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
//		cltest.Mine(h.sim, 1*time.Second)
//	}
//
//	return h.sim
//}

//	func (h *RpcHelper) Commit() {
//		h.sim.Commit()
//	}
func (h *RpcHelper) Client(t *testing.T) client.Client {
	if h.client != nil {
		return h.client
	}
	//return client.NewSimulatedBackendClient(t, h.sim, big.NewInt(1337))
	h.client = createChainClient(t)
	return h.client
}

func (h *RpcHelper) ChainID() *big.Int {
	return big.NewInt(sepolia_chain_id)
}

func (h *RpcHelper) NewSqlxDB(t *testing.T) *sqlx.DB {
	return pgtest.NewSqlxDB(t)
}

func (h *RpcHelper) Context(t *testing.T) context.Context {
	return testutils.Context(t)
}

func (h *RpcHelper) MaxWaitTimeForEvents() time.Duration {
	// From trial and error, when running on CI, sometimes the boxes get slow
	maxWaitTime := time.Second * 30
	maxWaitTimeStr, ok := os.LookupEnv("MAX_WAIT_TIME_FOR_EVENTS_S")
	if ok {
		waitS, err := strconv.ParseInt(maxWaitTimeStr, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing MAX_WAIT_TIME_FOR_EVENTS_S: %v, defaulting to %v\n", err, maxWaitTime)
		}
		maxWaitTime = time.Second * time.Duration(waitS)
	}
	return maxWaitTime
}

func (h *RpcHelper) TXM(t *testing.T, client client.Client) evmtxmgr.TxManager {
	if h.txm != nil {
		return h.txm
	}

	clconfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, secrets *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
		c.EVM[0].ChainID = (*ubig.Big)(h.ChainID())
	})

	clconfig.EVMConfigs()[0].GasEstimator.PriceMax = assets.GWei(100)
	wsUrl, err := commonconfig.ParseURL(sepolia_ws_endpoint)
	require.NoError(t, err)
	httpUrl, err := commonconfig.ParseURL(sepolia_http_endpoint)
	require.NoError(t, err)
	clconfig.EVMConfigs()[0].Nodes = append(clconfig.EVMConfigs()[0].Nodes, &evmcfg.Node{
		Name:     ptr(sepolia_network_name),
		WSURL:    wsUrl,
		HTTPURL:  httpUrl,
		SendOnly: ptr(false),
		Order:    ptr(int32(0)),
	})

	//flagsAndDeps := []interface{}{}
	//flagsAndDeps = append(flagsAndDeps, client)
	//flagsAndDeps = append(flagsAndDeps, ubig.New(h.ChainID()))
	app := cltest.NewApplicationWithConfigAndKey(t, clconfig, client, ubig.New(h.ChainID()))
	err = app.Start(h.Context(t))
	require.NoError(t, err)

	keyStore := app.KeyStore.Eth()

	keyStore.XXXTestingOnlyAdd(h.Context(t), keytypes.FromPrivateKey(h.deployerKey))
	require.NoError(t, keyStore.Add(h.Context(t), h.accounts[0].From, h.ChainID()))
	require.NoError(t, keyStore.Enable(h.Context(t), h.accounts[0].From, h.ChainID()))

	//keyStore.XXXTestingOnlyAdd(h.Context(t), keytypes.FromPrivateKey(h.senderKey))
	//require.NoError(t, keyStore.Add(h.Context(t), h.accounts[1].From, h.ChainID()))
	//require.NoError(t, keyStore.Enable(h.Context(t), h.accounts[1].From, h.ChainID()))

	chain, err := app.GetRelayers().LegacyEVMChains().Get((h.ChainID()).String())
	require.NoError(t, err)

	h.txm = chain.TxManager()
	return h.txm
}

func getEvmFeeEstimator(t *testing.T, err error, h *RpcHelper) *gas.EvmFeeEstimator {
	_, _, evmConfig := txmgr.MakeTestConfigs(t)

	estimator, err := gas.NewEstimator(logger.Test(t), h.client, evmConfig, evmConfig.GasEstimator())
	require.NoError(t, err)
	return &estimator
}
