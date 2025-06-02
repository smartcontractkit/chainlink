//nolint:revive // helpers for specific version
package testhelpers_1_4_0

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	types3 "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/freeport"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/utils/pointer" //nolint:staticcheck // tests

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	types4 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	pb "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	evmUtils "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	commit_store_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/commit_store"
	evm_2_evm_offramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_offramp"
	evm_2_evm_onramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_onramp"
	evmcapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	configv2 "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	feeds2 "github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	feedsMocks "github.com/smartcontractkit/chainlink/v2/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	ksMocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	clutils "github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

const (
	execSpecTemplate = `
		type = "offchainreporting2"
		schemaVersion = 1
		name = "ccip-exec-1"
		externalJobID      = "67ffad71-d90f-4fe3-b4e4-494924b707fb"
		forwardingAllowed = false
		maxTaskDuration = "0s"
		contractID = "%s"
		contractConfigConfirmations = 1
		contractConfigTrackerPollInterval = "20s"
		ocrKeyBundleID = "%s"
		relay = "evm"
		pluginType = "ccip-execution"
		transmitterID = "%s"

		[relayConfig]
		chainID = 1_337

		[pluginConfig]
		destStartBlock = 50

	    [pluginConfig.USDCConfig]
	    AttestationAPI = "http://blah.com"
	    SourceMessageTransmitterAddress = "%s"
	    SourceTokenAddress = "%s"
		AttestationAPITimeoutSeconds = 10
	`
	commitSpecTemplatePipeline = `
		type = "offchainreporting2"
		schemaVersion = 1
		name = "ccip-commit-1"
		externalJobID = "13c997cf-1a14-4ab7-9068-07ee6d2afa55"
		forwardingAllowed = false
		maxTaskDuration = "0s"
		contractID = "%s"
		contractConfigConfirmations = 1
		contractConfigTrackerPollInterval = "20s"
		ocrKeyBundleID = "%s"
		relay = "evm"
		pluginType = "ccip-commit"
		transmitterID = "%s"

		[relayConfig]
		chainID = 1_337

		[pluginConfig]
		destStartBlock = 50
		offRamp = "%s"
		tokenPricesUSDPipeline = """
		%s
		"""
	`
	commitSpecTemplateDynamicPriceGetter = `
		type = "offchainreporting2"
		schemaVersion = 1
		name = "ccip-commit-1"
		externalJobID = "13c997cf-1a14-4ab7-9068-07ee6d2afa55"
		forwardingAllowed = false
		maxTaskDuration = "0s"
		contractID = "%s"
		contractConfigConfirmations = 1
		contractConfigTrackerPollInterval = "20s"
		ocrKeyBundleID = "%s"
		relay = "evm"
		pluginType = "ccip-commit"
		transmitterID = "%s"

		[relayConfig]
		chainID = 1_337

		[pluginConfig]
		destStartBlock = 50
		offRamp = "%s"
		priceGetterConfig = """
		%s
		"""
	`
)

type Node struct {
	App             chainlink.Application
	Transmitter     common.Address
	PaymentReceiver common.Address
	KeyBundle       ocr2key.KeyBundle
}

func (node *Node) FindJobIDForContract(t *testing.T, addr common.Address) int32 {
	jobs := node.App.JobSpawner().ActiveJobs()
	for _, j := range jobs {
		if j.Type == job.OffchainReporting2 && j.OCR2OracleSpec.ContractID == addr.Hex() {
			return j.ID
		}
	}
	t.Fatalf("Could not find job for contract %s", addr.Hex())
	return 0
}

func (node *Node) EventuallyNodeUsesUpdatedPriceRegistry(t *testing.T, ccipContracts CCIPIntegrationTestHarness) logpoller.Log {
	c, err := node.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(ccipContracts.Dest.ChainID, 10))
	require.NoError(t, err)
	var log logpoller.Log
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		ccipContracts.Dest.Chain.Commit()
		log, err := c.LogPoller().LatestLogByEventSigWithConfs(
			testutils.Context(t),
			v1_2_0.UsdPerUnitGasUpdated,
			ccipContracts.Dest.PriceRegistry.Address(),
			0,
		)
		// err can be transient errors such as sql row set empty
		if err != nil {
			return false
		}
		return log != nil
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue(), "node is not using updated price registry %s", ccipContracts.Dest.PriceRegistry.Address().Hex())
	return log
}

func (node *Node) EventuallyNodeUsesNewCommitConfig(t *testing.T, ccipContracts CCIPIntegrationTestHarness, commitCfg ccipdata.CommitOnchainConfig) logpoller.Log {
	c, err := node.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(ccipContracts.Dest.ChainID, 10))
	require.NoError(t, err)
	var log logpoller.Log
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		ccipContracts.Dest.Chain.Commit()
		log, err := c.LogPoller().LatestLogByEventSigWithConfs(
			testutils.Context(t),
			evmrelay.OCR2AggregatorLogDecoder.EventSig(),
			ccipContracts.Dest.CommitStore.Address(),
			0,
		)
		require.NoError(t, err)
		var latestCfg ccipdata.CommitOnchainConfig
		if log != nil {
			latestCfg, err = DecodeCommitOnChainConfig(log.Data)
			require.NoError(t, err)
			return latestCfg == commitCfg
		}
		return false
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue(), "node is using old cfg")
	return log
}

func (node *Node) EventuallyNodeUsesNewExecConfig(t *testing.T, ccipContracts CCIPIntegrationTestHarness, execCfg v1_2_0.ExecOnchainConfig) logpoller.Log {
	c, err := node.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(ccipContracts.Dest.ChainID, 10))
	require.NoError(t, err)
	var log logpoller.Log
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		ccipContracts.Dest.Chain.Commit()
		log, err := c.LogPoller().LatestLogByEventSigWithConfs(
			testutils.Context(t),
			evmrelay.OCR2AggregatorLogDecoder.EventSig(),
			ccipContracts.Dest.OffRamp.Address(),
			0,
		)
		require.NoError(t, err)
		var latestCfg v1_2_0.ExecOnchainConfig
		if log != nil {
			latestCfg, err = DecodeExecOnChainConfig(log.Data)
			require.NoError(t, err)
			return latestCfg == execCfg
		}
		return false
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue(), "node is using old cfg")
	return log
}

//nolint:gosec // safe cast in tests
func (node *Node) EventuallyHasReqSeqNum(t *testing.T, ccipContracts *CCIPIntegrationTestHarness, onRamp common.Address, seqNum int) logpoller.Log {
	c, err := node.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(ccipContracts.Source.ChainID, 10))
	require.NoError(t, err)
	var log logpoller.Log
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		ccipContracts.Dest.Chain.Commit()
		lgs, err := c.LogPoller().LogsDataWordRange(
			testutils.Context(t),
			v1_2_0.CCIPSendRequestEventSig,
			onRamp,
			v1_2_0.CCIPSendRequestSeqNumIndex,
			abihelpers.EvmWord(uint64(seqNum)),
			abihelpers.EvmWord(uint64(seqNum)),
			1,
		)
		require.NoError(t, err)
		t.Log("Send requested", len(lgs))
		if len(lgs) == 1 {
			log = lgs[0]
			return true
		}
		return false
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue(), "eventually has seq num")
	return log
}

//nolint:gosec // safe cast in tests
func (node *Node) EventuallyHasExecutedSeqNums(t *testing.T, ccipContracts *CCIPIntegrationTestHarness, offRamp common.Address, minSeqNum int, maxSeqNum int) []logpoller.Log {
	c, err := node.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(ccipContracts.Dest.ChainID, 10))
	require.NoError(t, err)
	var logs []logpoller.Log
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		ccipContracts.Dest.Chain.Commit()
		lgs, err := c.LogPoller().IndexedLogsTopicRange(
			testutils.Context(t),
			v1_2_0.ExecutionStateChangedEvent,
			offRamp,
			v1_2_0.ExecutionStateChangedSeqNrIndex,
			abihelpers.EvmWord(uint64(minSeqNum)),
			abihelpers.EvmWord(uint64(maxSeqNum)),
			1,
		)
		require.NoError(t, err)
		t.Logf("Have executed logs %d want %d", len(lgs), maxSeqNum-minSeqNum+1)
		if len(lgs) == maxSeqNum-minSeqNum+1 {
			logs = lgs
			t.Logf("Seq Num %d-%d executed", minSeqNum, maxSeqNum)
			return true
		}
		return false
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue(), "eventually has not executed seq num")
	return logs
}

//nolint:gosec // safe to casts in tests
func (node *Node) ConsistentlySeqNumHasNotBeenExecuted(t *testing.T, ccipContracts *CCIPIntegrationTestHarness, offRamp common.Address, seqNum int) logpoller.Log {
	c, err := node.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(ccipContracts.Dest.ChainID, 10))
	require.NoError(t, err)
	var log logpoller.Log
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		ccipContracts.Source.Chain.Commit()
		ccipContracts.Dest.Chain.Commit()
		lgs, err := c.LogPoller().IndexedLogsTopicRange(
			testutils.Context(t),
			v1_2_0.ExecutionStateChangedEvent,
			offRamp,
			v1_2_0.ExecutionStateChangedSeqNrIndex,
			abihelpers.EvmWord(uint64(seqNum)),
			abihelpers.EvmWord(uint64(seqNum)),
			1,
		)
		require.NoError(t, err)
		t.Log("Executed logs", lgs)
		if len(lgs) == 1 {
			log = lgs[0]
			return true
		}
		return false
	}, 10*time.Second, 1*time.Second).Should(gomega.BeFalse(), "seq number got executed")
	return log
}

func (node *Node) AddJob(t *testing.T, spec *integrationtesthelpers.OCR2TaskJobSpec) {
	specString, err := spec.String()
	require.NoError(t, err)
	ccipJob, err := validate.ValidatedOracleSpecToml(
		testutils.Context(t),
		node.App.GetConfig().OCR2(),
		node.App.GetConfig().Insecure(),
		specString,
		// FIXME Ani
		nil,
	)
	require.NoError(t, err)
	err = node.App.AddJobV2(context.Background(), &ccipJob)
	require.NoError(t, err)
}

func (node *Node) AddBootstrapJob(t *testing.T, spec *integrationtesthelpers.OCR2TaskJobSpec) {
	specString, err := spec.String()
	require.NoError(t, err)
	ccipJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(specString)
	require.NoError(t, err)
	err = node.App.AddJobV2(context.Background(), &ccipJob)
	require.NoError(t, err)
}

func (node *Node) AddJobsWithSpec(t *testing.T, jobSpec *integrationtesthelpers.OCR2TaskJobSpec) {
	// set node specific values
	jobSpec.OCR2OracleSpec.OCRKeyBundleID.SetValid(node.KeyBundle.ID())
	jobSpec.OCR2OracleSpec.TransmitterID.SetValid(node.Transmitter.Hex())
	node.AddJob(t, jobSpec)
}

func setupNodeCCIP(
	t *testing.T,
	owner *bind.TransactOpts,
	port int64,
	dbName string,
	sourceChain *backends.SimulatedBackend, destChain *backends.SimulatedBackend,
	sourceChainID *big.Int, destChainID *big.Int,
	bootstrapPeerID string,
	bootstrapPort int64,
) (chainlink.Application, string, common.Address, ocr2key.KeyBundle) {
	trueRef, falseRef := true, false

	// Do not want to load fixtures as they contain a dummy chainID.
	loglevel := configv2.LogLevel(zap.DebugLevel)
	config, db := heavyweight.FullTestDBNoFixturesV2(t, func(c *chainlink.Config, _ *chainlink.Secrets) {
		p2pAddresses := []string{
			fmt.Sprintf("127.0.0.1:%d", port),
		}
		c.Log.Level = &loglevel
		c.Feature.CCIP = &trueRef
		c.Feature.UICSAKeys = &trueRef
		c.Feature.FeedsManager = &trueRef
		c.OCR.Enabled = &falseRef
		c.OCR.DefaultTransactionQueueDepth = pointer.Uint32(200)
		c.OCR2.Enabled = &trueRef
		c.Feature.LogPoller = &trueRef
		c.P2P.V2.Enabled = &trueRef

		dur, err := config.NewDuration(500 * time.Millisecond)
		if err != nil {
			panic(err)
		}
		c.P2P.V2.DeltaDial = &dur

		dur2, err := config.NewDuration(5 * time.Second)
		if err != nil {
			panic(err)
		}

		c.P2P.V2.DeltaReconcile = &dur2
		c.P2P.V2.ListenAddresses = &p2pAddresses
		c.P2P.V2.AnnounceAddresses = &p2pAddresses

		c.EVM = []*toml.EVMConfig{createConfigV2Chain(sourceChainID), createConfigV2Chain(destChainID)}

		if bootstrapPeerID != "" {
			// Supply the bootstrap IP and port as a V2 peer address
			c.P2P.V2.DefaultBootstrappers = &[]commontypes.BootstrapperLocator{
				{
					PeerID: bootstrapPeerID, Addrs: []string{
						fmt.Sprintf("127.0.0.1:%d", bootstrapPort),
					},
				},
			}
		}
	})

	lggr := logger.TestLogger(t)

	// The in-memory geth sim does not let you create a custom ChainID, it will always be 1337.
	// In particular this means that if you sign an eip155 tx, the chainID used MUST be 1337
	// and the CHAINID op code will always emit 1337. To work around this to simulate a "multichain"
	// test, we fake different chainIDs using the wrapped sim cltest.SimulatedBackend so the RPC
	// appears to operate on different chainIDs and we use an EthKeyStoreSim wrapper which always
	// signs 1337 see https://github.com/smartcontractkit/chainlink-ccip/blob/a24dd436810250a458d27d8bb3fb78096afeb79c/core/services/ocr2/plugins/ccip/testhelpers/simulated_backend.go#L35
	sourceClient := client.NewSimulatedBackendClient(t, sourceChain.Backend, sourceChainID)
	destClient := client.NewSimulatedBackendClient(t, destChain.Backend, destChainID)
	csaKeyStore := ksMocks.NewCSA(t)

	key, err := csakey.NewV2()
	require.NoError(t, err)
	csaKeyStore.On("EnsureKey", mock.Anything).Return(nil)
	csaKeyStore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
	keyStore := NewKsa(db, lggr, csaKeyStore)
	ctx := testutils.Context(t)
	app, err := chainlink.NewApplication(ctx, chainlink.ApplicationOpts{
		CREOpts: chainlink.CREOpts{
			CapabilitiesRegistry: evmcapabilities.NewRegistry(lggr),
		},
		Config:   config,
		DS:       db,
		KeyStore: keyStore,
		EVMFactoryConfigFn: func(fc *chainlink.EVMFactoryConfig) {
			fc.GenEthClient = func(chainID *big.Int) client.Client {
				if chainID.String() == sourceChainID.String() {
					return sourceClient
				} else if chainID.String() == destChainID.String() {
					return destClient
				}
				t.Fatalf("invalid chain ID %v", chainID.String())
				return nil
			}
		},
		Logger:                   lggr,
		ExternalInitiatorManager: nil,
		CloseLogger:              lggr.Sync,
		UnrestrictedHTTPClient:   &http.Client{},
		RestrictedHTTPClient:     &http.Client{},
		AuditLogger:              audit.NoopLogger,
	})

	require.NoError(t, err)
	require.NoError(t, app.GetKeyStore().Unlock(ctx, "password"))
	_, err = app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()

	testCtx := testutils.Context(t)
	_, err = app.GetKeyStore().Eth().Create(testCtx, destChainID)
	require.NoError(t, err)
	sendingKeys, err := app.GetKeyStore().Eth().EnabledKeysForChain(testCtx, destChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address
	s, err := app.GetKeyStore().Eth().GetState(testCtx, sendingKeys[0].ID(), destChainID)
	require.NoError(t, err)
	lggr.Debug(fmt.Sprintf("Transmitter address %s chainID %s", transmitter, s.EVMChainID.String()))

	// Fund the commitTransmitter address with some ETH
	n, err := destChain.NonceAt(context.Background(), owner.From, nil)
	require.NoError(t, err)

	tx := types3.NewTransaction(n, transmitter, big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = destChain.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	destChain.Commit()

	kb, err := app.GetKeyStore().OCR2().Create(ctx, chaintype.EVM)
	require.NoError(t, err)
	return app, peerID.Raw(), transmitter, kb
}

func createConfigV2Chain(chainID *big.Int) *toml.EVMConfig {
	// NOTE: For the executor jobs, the default of 500k is insufficient for a 3 message batch
	defaultGasLimit := uint64(5000000)
	tr := true

	sourceC := toml.Defaults((*evmUtils.Big)(chainID))
	sourceC.GasEstimator.LimitDefault = &defaultGasLimit
	fixedPrice := "FixedPrice"
	sourceC.GasEstimator.Mode = &fixedPrice
	d, _ := config.NewDuration(100 * time.Millisecond)
	sourceC.LogPollInterval = &d
	fd := uint32(2)
	sourceC.FinalityDepth = &fd
	return &toml.EVMConfig{
		ChainID: (*evmUtils.Big)(chainID),
		Enabled: &tr,
		Chain:   sourceC,
		Nodes:   toml.EVMNodes{&toml.Node{}},
	}
}

type CCIPIntegrationTestHarness struct {
	CCIPContracts
	Nodes     []Node
	Bootstrap Node
}

func SetupCCIPIntegrationTH(t *testing.T, sourceChainID, sourceChainSelector, destChainId, destChainSelector uint64) CCIPIntegrationTestHarness {
	return CCIPIntegrationTestHarness{
		CCIPContracts: SetupCCIPContracts(t, sourceChainID, sourceChainSelector, destChainId, destChainSelector),
	}
}

//nolint:testifylint //require is used for assertions in handlers
func (c *CCIPIntegrationTestHarness) CreatePricesPipeline(t *testing.T) (string, *httptest.Server, *httptest.Server) {
	linkUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"UsdPerLink": "8000000000000000000"}`))
		require.NoError(t, err)
	}))
	ethUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"UsdPerETH": "1700000000000000000000"}`))
		require.NoError(t, err)
	}))
	sourceWrappedNative, err := c.Source.Router.GetWrappedNative(nil)
	require.NoError(t, err)
	destWrappedNative, err := c.Dest.Router.GetWrappedNative(nil)
	require.NoError(t, err)
	tokenPricesUSDPipeline := fmt.Sprintf(`
// Price 1
link [type=http method=GET url="%s"];
link_parse [type=jsonparse path="UsdPerLink"];
link->link_parse;
eth [type=http method=GET url="%s"];
eth_parse [type=jsonparse path="UsdPerETH"];
eth->eth_parse;
merge [type=merge left="{}" right="{\\\"%s\\\":$(link_parse), \\\"%s\\\":$(eth_parse), \\\"%s\\\":$(eth_parse)}"];`,
		linkUSD.URL, ethUSD.URL, c.Dest.LinkToken.Address(), sourceWrappedNative, destWrappedNative)

	return tokenPricesUSDPipeline, linkUSD, ethUSD
}

func (c *CCIPIntegrationTestHarness) AddAllJobs(t *testing.T, jobParams integrationtesthelpers.CCIPJobSpecParams) {
	jobParams.OffRamp = c.Dest.OffRamp.Address()

	commitSpec, err := jobParams.CommitJobSpec()
	require.NoError(t, err)
	geExecutionSpec, err := jobParams.ExecutionJobSpec()
	require.NoError(t, err)
	nodes := c.Nodes
	for _, node := range nodes {
		node.AddJobsWithSpec(t, commitSpec)
		node.AddJobsWithSpec(t, geExecutionSpec)
	}
}

func (c *CCIPIntegrationTestHarness) jobSpecProposal(t *testing.T, specTemplate string, f func() (*integrationtesthelpers.OCR2TaskJobSpec, error), feedsManagerId int64, version int32, opts ...any) feeds2.ProposeJobArgs {
	spec, err := f()
	require.NoError(t, err)

	args := []any{spec.OCR2OracleSpec.ContractID}
	args = append(args, opts...)

	return feeds2.ProposeJobArgs{
		FeedsManagerID: feedsManagerId,
		RemoteUUID:     uuid.New(),
		Multiaddrs:     nil,
		Version:        version,
		Spec:           fmt.Sprintf(specTemplate, args...),
	}
}

func (c *CCIPIntegrationTestHarness) SetupFeedsManager(t *testing.T) {
	ctx := testutils.Context(t)
	for _, node := range c.Nodes {
		f := node.App.GetFeedsService()

		managers, err := f.ListManagers(ctx)
		require.NoError(t, err)
		if len(managers) > 0 {
			// Use at most one feeds manager, don't register if one already exists
			continue
		}

		secret := utils.RandomBytes32()
		pkey, err := crypto.PublicKeyFromHex(hex.EncodeToString(secret[:]))
		require.NoError(t, err)

		m := feeds2.RegisterManagerParams{
			Name:      "CCIP",
			URI:       "http://localhost:8080",
			PublicKey: *pkey,
		}

		_, err = f.RegisterManager(testutils.Context(t), m)
		require.NoError(t, err)

		connManager := feedsMocks.NewConnectionsManager(t)
		connManager.On("GetClient", mock.Anything).Maybe().Return(NoopFeedsClient{}, nil)
		connManager.On("Close").Maybe().Return()
		connManager.On("IsConnected", mock.Anything).Maybe().Return(true)
		f.Unsafe_SetConnectionsManager(connManager)
	}
}

func (c *CCIPIntegrationTestHarness) ApproveJobSpecs(t *testing.T, jobParams integrationtesthelpers.CCIPJobSpecParams) {
	ctx := testutils.Context(t)

	for _, node := range c.Nodes {
		f := node.App.GetFeedsService()
		managers, err := f.ListManagers(ctx)
		require.NoError(t, err)
		require.Len(t, managers, 1, "expected exactly one feeds manager")

		execSpec := c.jobSpecProposal(
			t,
			execSpecTemplate,
			jobParams.ExecutionJobSpec,
			managers[0].ID,
			1,
			node.KeyBundle.ID(),
			node.Transmitter.Hex(),
			utils.RandomAddress().String(),
			utils.RandomAddress().String(),
		)
		execID, err := f.ProposeJob(ctx, &execSpec)
		require.NoError(t, err)

		err = f.ApproveSpec(ctx, execID, true)
		require.NoError(t, err)

		var commitSpec feeds2.ProposeJobArgs
		if jobParams.TokenPricesUSDPipeline != "" {
			commitSpec = c.jobSpecProposal(
				t,
				commitSpecTemplatePipeline,
				jobParams.CommitJobSpec,
				managers[0].ID,
				2,
				node.KeyBundle.ID(),
				node.Transmitter.Hex(),
				jobParams.OffRamp.String(),
				jobParams.TokenPricesUSDPipeline,
			)
		} else {
			commitSpec = c.jobSpecProposal(
				t,
				commitSpecTemplateDynamicPriceGetter,
				jobParams.CommitJobSpec,
				managers[0].ID,
				2,
				node.KeyBundle.ID(),
				node.Transmitter.Hex(),
				jobParams.OffRamp.String(),
				jobParams.PriceGetterConfig,
			)
		}

		commitID, err := f.ProposeJob(ctx, &commitSpec)
		require.NoError(t, err)

		err = f.ApproveSpec(ctx, commitID, true)
		require.NoError(t, err)
	}
}

func (c *CCIPIntegrationTestHarness) AllNodesHaveReqSeqNum(t *testing.T, seqNum int, onRampOpts ...common.Address) logpoller.Log {
	var log logpoller.Log
	nodes := c.Nodes
	var onRamp common.Address
	if len(onRampOpts) > 0 {
		onRamp = onRampOpts[0]
	} else {
		require.NotNil(t, c.Source.OnRamp, "no onramp configured")
		onRamp = c.Source.OnRamp.Address()
	}
	for _, node := range nodes {
		log = node.EventuallyHasReqSeqNum(t, c, onRamp, seqNum)
	}
	return log
}

func (c *CCIPIntegrationTestHarness) AllNodesHaveExecutedSeqNums(t *testing.T, minSeqNum int, maxSeqNum int, offRampOpts ...common.Address) []logpoller.Log {
	var logs []logpoller.Log
	nodes := c.Nodes
	var offRamp common.Address

	if len(offRampOpts) > 0 {
		offRamp = offRampOpts[0]
	} else {
		require.NotNil(t, c.Dest.OffRamp, "no offramp configured")
		offRamp = c.Dest.OffRamp.Address()
	}
	for _, node := range nodes {
		logs = node.EventuallyHasExecutedSeqNums(t, c, offRamp, minSeqNum, maxSeqNum)
	}
	return logs
}

func (c *CCIPIntegrationTestHarness) NoNodesHaveExecutedSeqNum(t *testing.T, seqNum int, offRampOpts ...common.Address) logpoller.Log {
	var log logpoller.Log
	nodes := c.Nodes
	var offRamp common.Address
	if len(offRampOpts) > 0 {
		offRamp = offRampOpts[0]
	} else {
		require.NotNil(t, c.Dest.OffRamp, "no offramp configured")
		offRamp = c.Dest.OffRamp.Address()
	}
	for _, node := range nodes {
		log = node.ConsistentlySeqNumHasNotBeenExecuted(t, c, offRamp, seqNum)
	}
	return log
}

func (c *CCIPIntegrationTestHarness) EventuallyCommitReportAccepted(t *testing.T, currentBlock uint64, commitStoreOpts ...common.Address) commit_store_1_2_0.CommitStoreCommitReport {
	var commitStore *commit_store_1_2_0.CommitStore
	var err error
	if len(commitStoreOpts) > 0 {
		commitStore, err = commit_store_1_2_0.NewCommitStore(commitStoreOpts[0], c.Dest.Chain)
		require.NoError(t, err)
	} else {
		require.NotNil(t, c.Dest.CommitStore, "no commitStore configured")
		commitStore = c.Dest.CommitStore
	}
	g := gomega.NewGomegaWithT(t)
	var report commit_store_1_2_0.CommitStoreCommitReport
	g.Eventually(func() bool {
		it, err := commitStore.FilterReportAccepted(&bind.FilterOpts{Start: currentBlock})
		g.Expect(err).NotTo(gomega.HaveOccurred(), "Error filtering ReportAccepted event")
		g.Expect(it.Next()).To(gomega.BeTrue(), "No ReportAccepted event found")
		report = it.Event.Report
		if report.MerkleRoot != [32]byte{} {
			t.Log("Report Accepted by commitStore")
			return true
		}
		return false
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue(), "report has not been committed")
	return report
}

func (c *CCIPIntegrationTestHarness) EventuallyExecutionStateChangedToSuccess(t *testing.T, seqNum []uint64, blockNum uint64, offRampOpts ...common.Address) {
	var offRamp *evm_2_evm_offramp_1_2_0.EVM2EVMOffRamp
	var err error
	if len(offRampOpts) > 0 {
		offRamp, err = evm_2_evm_offramp_1_2_0.NewEVM2EVMOffRamp(offRampOpts[0], c.Dest.Chain)
		require.NoError(t, err)
	} else {
		require.NotNil(t, c.Dest.OffRamp, "no offRamp configured")
		offRamp = c.Dest.OffRamp
	}
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		it, err := offRamp.FilterExecutionStateChanged(&bind.FilterOpts{Start: blockNum}, seqNum, [][32]byte{})
		require.NoError(t, err)
		for it.Next() {
			if cciptypes.MessageExecutionState(it.Event.State) == cciptypes.ExecutionStateSuccess {
				t.Logf("ExecutionStateChanged event found for seqNum %d", it.Event.SequenceNumber)
				return true
			}
		}
		c.Source.Chain.Commit()
		c.Dest.Chain.Commit()
		return false
	}, testutils.WaitTimeout(t), time.Second).
		Should(gomega.BeTrue(), "ExecutionStateChanged Event")
}

//nolint:gosec // safe to casts in tests
func (c *CCIPIntegrationTestHarness) EventuallyReportCommitted(t *testing.T, maxSeqNr int, commitStoreOpts ...common.Address) uint64 {
	var commitStore *commit_store_1_2_0.CommitStore
	var err error
	var committedSeqNum uint64
	if len(commitStoreOpts) > 0 {
		commitStore, err = commit_store_1_2_0.NewCommitStore(commitStoreOpts[0], c.Dest.Chain)
		require.NoError(t, err)
	} else {
		require.NotNil(t, c.Dest.CommitStore, "no commitStore configured")
		commitStore = c.Dest.CommitStore
	}
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		minSeqNum, err := commitStore.GetExpectedNextSequenceNumber(nil)
		require.NoError(t, err)
		c.Source.Chain.Commit()
		c.Dest.Chain.Commit()
		t.Log("next expected seq num reported", minSeqNum)
		committedSeqNum = minSeqNum
		return minSeqNum > uint64(maxSeqNr)
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue(), "report has not been committed")
	return committedSeqNum
}

func (c *CCIPIntegrationTestHarness) EventuallySendRequested(t *testing.T, seqNum uint64, onRampOpts ...common.Address) {
	var onRamp *evm_2_evm_onramp_1_2_0.EVM2EVMOnRamp
	var err error
	if len(onRampOpts) > 0 {
		onRamp, err = evm_2_evm_onramp_1_2_0.NewEVM2EVMOnRamp(onRampOpts[0], c.Source.Chain)
		require.NoError(t, err)
	} else {
		require.NotNil(t, c.Source.OnRamp, "no onRamp configured")
		onRamp = c.Source.OnRamp
	}
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		it, err := onRamp.FilterCCIPSendRequested(nil)
		require.NoError(t, err)
		for it.Next() {
			if it.Event.Message.SequenceNumber == seqNum {
				t.Log("sendRequested generated for", seqNum)
				return true
			}
		}
		c.Source.Chain.Commit()
		c.Dest.Chain.Commit()
		return false
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue(), "sendRequested has not been generated")
}

//nolint:gosec // safe to cast in tests
func (c *CCIPIntegrationTestHarness) ConsistentlyReportNotCommitted(t *testing.T, maxSeqNr int, commitStoreOpts ...common.Address) {
	var commitStore *commit_store_1_2_0.CommitStore
	var err error
	if len(commitStoreOpts) > 0 {
		commitStore, err = commit_store_1_2_0.NewCommitStore(commitStoreOpts[0], c.Dest.Chain)
		require.NoError(t, err)
	} else {
		require.NotNil(t, c.Dest.CommitStore, "no commitStore configured")
		commitStore = c.Dest.CommitStore
	}
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		minSeqNum, err := commitStore.GetExpectedNextSequenceNumber(nil)
		require.NoError(t, err)
		c.Source.Chain.Commit()
		c.Dest.Chain.Commit()
		t.Log("min seq num reported", minSeqNum)
		return minSeqNum > uint64(maxSeqNr)
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeFalse(), "report has been committed")
}

func (c *CCIPIntegrationTestHarness) SetupAndStartNodes(ctx context.Context, t *testing.T, bootstrapNodePort int64) (Node, []Node, int64) {
	appBootstrap, bootstrapPeerID, bootstrapTransmitter, bootstrapKb := setupNodeCCIP(t, c.Dest.User, bootstrapNodePort,
		"bootstrap_ccip", c.Source.Chain, c.Dest.Chain, big.NewInt(0).SetUint64(c.Source.ChainID),
		big.NewInt(0).SetUint64(c.Dest.ChainID), "", 0)
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	err := appBootstrap.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, appBootstrap.Stop())
	})
	bootstrapNode := Node{
		App:         appBootstrap,
		Transmitter: bootstrapTransmitter,
		KeyBundle:   bootstrapKb,
	}
	// Set up the minimum 4 oracles all funded with destination ETH
	for i := int64(0); i < 4; i++ {
		app, peerID, transmitter, kb := setupNodeCCIP(
			t,
			c.Dest.User,
			int64(freeport.GetOne(t)),
			fmt.Sprintf("oracle_ccip%d", i),
			c.Source.Chain,
			c.Dest.Chain,
			big.NewInt(0).SetUint64(c.Source.ChainID),
			big.NewInt(0).SetUint64(c.Dest.ChainID),
			bootstrapPeerID,
			bootstrapNodePort,
		)
		nodes = append(nodes, Node{
			App:         app,
			Transmitter: transmitter,
			KeyBundle:   kb,
		})
		offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   types4.Account(transmitter.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
		err = app.Start(ctx)
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, app.Stop())
		})
	}

	c.Oracles = oracles
	commitOnchainConfig := c.CreateDefaultCommitOnchainConfig(t)
	commitOffchainConfig := c.CreateDefaultCommitOffchainConfig(t)
	execOnchainConfig := c.CreateDefaultExecOnchainConfig(t)
	execOffchainConfig := c.CreateDefaultExecOffchainConfig(t)

	configBlock := c.SetupOnchainConfig(t, commitOnchainConfig, commitOffchainConfig, execOnchainConfig, execOffchainConfig)
	c.Nodes = nodes
	c.Bootstrap = bootstrapNode
	return bootstrapNode, nodes, configBlock
}

func (c *CCIPIntegrationTestHarness) SetUpNodesAndJobs(t *testing.T, pricePipeline string, priceGetterConfig string, usdcAttestationAPI string) integrationtesthelpers.CCIPJobSpecParams {
	// setup Jobs
	ctx := context.Background()
	// Starts nodes and configures them in the OCR contracts.
	bootstrapNode, _, configBlock := c.SetupAndStartNodes(ctx, t, int64(freeport.GetOne(t)))

	jobParams := c.NewCCIPJobSpecParams(pricePipeline, priceGetterConfig, configBlock, usdcAttestationAPI)

	// Add the bootstrap job
	c.Bootstrap.AddBootstrapJob(t, jobParams.BootstrapJob(c.Dest.CommitStore.Address().Hex()))
	c.AddAllJobs(t, jobParams)

	// Replay for bootstrap.
	bc, err := bootstrapNode.App.GetRelayers().LegacyEVMChains().Get(strconv.FormatUint(c.Dest.ChainID, 10))
	require.NoError(t, err)
	require.NoError(t, bc.LogPoller().Replay(context.Background(), configBlock))
	c.Dest.Chain.Commit()

	return jobParams
}

//nolint:gosec // safe to cast in tests
func (c *CCIPIntegrationTestHarness) NewCCIPJobSpecParams(tokenPricesUSDPipeline string, priceGetterConfig string, configBlock int64, usdcAttestationAPI string) integrationtesthelpers.CCIPJobSpecParams {
	return integrationtesthelpers.CCIPJobSpecParams{
		CommitStore:            c.Dest.CommitStore.Address(),
		OffRamp:                c.Dest.OffRamp.Address(),
		DestEvmChainId:         c.Dest.ChainID,
		SourceChainName:        "SimulatedSource",
		DestChainName:          "SimulatedDest",
		TokenPricesUSDPipeline: tokenPricesUSDPipeline,
		PriceGetterConfig:      priceGetterConfig,
		DestStartBlock:         uint64(configBlock),
		USDCAttestationAPI:     usdcAttestationAPI,
	}
}

func DecodeCommitOnChainConfig(encoded []byte) (ccipdata.CommitOnchainConfig, error) {
	var onchainConfig ccipdata.CommitOnchainConfig
	unpacked, err := abihelpers.DecodeOCR2Config(encoded)
	if err != nil {
		return onchainConfig, err
	}
	onChainCfg := unpacked.OnchainConfig
	onchainConfig, err = abihelpers.DecodeAbiStruct[ccipdata.CommitOnchainConfig](onChainCfg)
	if err != nil {
		return onchainConfig, err
	}
	return onchainConfig, nil
}

func DecodeExecOnChainConfig(encoded []byte) (v1_2_0.ExecOnchainConfig, error) {
	var onchainConfig v1_2_0.ExecOnchainConfig
	unpacked, err := abihelpers.DecodeOCR2Config(encoded)
	if err != nil {
		return onchainConfig, errors.Wrap(err, "failed to unpack log data")
	}
	onChainCfg := unpacked.OnchainConfig
	onchainConfig, err = abihelpers.DecodeAbiStruct[v1_2_0.ExecOnchainConfig](onChainCfg)
	if err != nil {
		return onchainConfig, err
	}
	return onchainConfig, nil
}

type ksa struct {
	keystore.Master
	csa keystore.CSA
}

func (k *ksa) CSA() keystore.CSA {
	return k.csa
}

func NewKsa(db *sqlx.DB, lggr logger.Logger, csa keystore.CSA) *ksa {
	return &ksa{
		Master: keystore.New(db, clutils.FastScryptParams, lggr),
		csa:    csa,
	}
}

type NoopFeedsClient struct{}

func (n NoopFeedsClient) ApprovedJob(context.Context, *pb.ApprovedJobRequest) (*pb.ApprovedJobResponse, error) {
	return &pb.ApprovedJobResponse{}, nil
}

func (n NoopFeedsClient) Healthcheck(context.Context, *pb.HealthcheckRequest) (*pb.HealthcheckResponse, error) {
	return &pb.HealthcheckResponse{}, nil
}

func (n NoopFeedsClient) UpdateNode(context.Context, *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	return &pb.UpdateNodeResponse{}, nil
}

func (n NoopFeedsClient) RejectedJob(context.Context, *pb.RejectedJobRequest) (*pb.RejectedJobResponse, error) {
	return &pb.RejectedJobResponse{}, nil
}

func (n NoopFeedsClient) CancelledJob(context.Context, *pb.CancelledJobRequest) (*pb.CancelledJobResponse, error) {
	return &pb.CancelledJobResponse{}, nil
}
