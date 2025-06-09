package llo_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	"github.com/smartcontractkit/wsrpc/credentials"
)

type Node struct {
	App          chainlink.Application
	ClientPubKey credentials.StaticSizedPublicKey
	KeyBundle    ocr2key.KeyBundle
	ObservedLogs *observer.ObservedLogs
}

func (node *Node) addBootstrapJob(t *testing.T, spec string) *job.Job {
	job, err := ocrbootstrap.ValidatedBootstrapSpecToml(spec)
	require.NoError(t, err)
	err = node.App.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	return &job
}

func setupNode(
	t *testing.T,
	port int,
	dbName string,
	backend evmtypes.Backend,
	csaKey csakey.KeyV2,
	f func(*chainlink.Config),
) (app chainlink.Application, peerID string, clientPubKey credentials.StaticSizedPublicKey, ocr2kb ocr2key.KeyBundle, observedLogs *observer.ObservedLogs) {
	k := big.NewInt(int64(port)) // keys unique to port
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(k)
	rdr := keystest.NewRandReaderFromSeed(int64(port))
	ocr2kb = ocr2key.MustNewInsecure(rdr, chaintype.EVM)

	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		// TODO(gg): potentially update node config here

		// set finality depth to 1 so we don't have to wait for multiple blocks
		c.EVM[0].FinalityDepth = ptr[uint32](1)

		// [JobPipeline]
		c.JobPipeline.MaxSuccessfulRuns = ptr(uint64(1000))
		c.JobPipeline.VerboseLogging = ptr(true)

		// [Feature]
		c.Feature.UICSAKeys = ptr(true)
		c.Feature.LogPoller = ptr(true)
		c.Feature.FeedsManager = ptr(false)

		// [OCR]
		c.OCR.Enabled = ptr(false)

		// [OCR2]
		c.OCR2.Enabled = ptr(true)
		c.OCR2.ContractPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)

		// [P2P]
		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.TraceLogging = ptr(true)

		// [P2P.V2]
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		// [Mercury]
		c.Mercury.VerboseLogging = ptr(true)

		// [Log]
		c.Log.Level = ptr(toml.LogLevel(zapcore.DebugLevel)) // generally speaking we want debug level for logs unless overridden

		// [EVM.Transactions]
		for _, evmCfg := range c.EVM {
			evmCfg.Transactions.Enabled = ptr(false) // don't need txmgr
		}

		// Optional overrides
		if f != nil {
			f(c)
		}
	})

	lggr, observedLogs := logger.TestLoggerObserved(t, config.Log().Level())
	if backend != nil {
		app = cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend, p2pKey, ocr2kb, csaKey, lggr.Named(dbName))
	} else {
		app = cltest.NewApplicationWithConfig(t, config, p2pKey, ocr2kb, csaKey, lggr.Named(dbName))
	}
	err := app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb, observedLogs
}

func ptr[T any](t T) *T { return &t }

func createSecureMintBootstrapJob(t *testing.T, bootstrapNode Node, configuratorAddress common.Address, chainID, fromBlock string) *job.Job {
	return bootstrapNode.addBootstrapJob(t, fmt.Sprintf(`
		type                              = "bootstrap"
		relay                             = "evm"
		schemaVersion                     = 1
		name                              = "bootstrap-secure-mint"
		contractID                        = "%s"
		contractConfigTrackerPollInterval = "1s"
		contractConfigConfirmations = 1

		[relayConfig]
		chainID = %s
		fromBlock = %s

		providerType = "securemint"`,
		configuratorAddress.Hex(),
		chainID,
		fromBlock),
	)
}

func addSecureMintOCRJobs(
	t *testing.T,
	nodes []Node,
	configuratorAddress common.Address,
) (jobIDs map[int]int32) {

	// node idx => job id
	jobIDs = make(map[int]int32)

	// Create one bridge and one SM Feed OCR job on each node
	for i, node := range nodes {
		name := "securemint-ea"
		bridgeResp := por.ExternalAdapterPayload{
			Mintables: por.Mintables{
				por.ChainSelector(uint64(1)): por.BlockMintablePair{
					Block:    por.BlockNumber(1),
					Mintable: big.NewInt(1000000000),
				},
			},
			LatestRelevantBlocks: por.Blocks{
				por.ChainSelector(uint64(1)): por.BlockNumber(1),
			},
			ReserveInfo: por.ReserveInfo{
				ReserveAmount: big.NewInt(1000000000),
				Timestamp:     time.Now(),
			},
		}
		bmBridge := createSecureMintBridge(t, name, i, bridgeResp, node.App.BridgeORM())
		t.Logf("Created secure mint bridge %s on node %d", bmBridge, i)

		addresses, err := node.App.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		t.Logf("Using transmitter address %s for node %d", addresses[0].String(), i)

		jobID := addSecureMintJob(i,
			t,
			node,
			configuratorAddress,
			bmBridge,
		)
		jobIDs[i] = jobID
		t.Logf("Added secure mint job with id %d on node %d", jobID, i)
	}
	return jobIDs
}

func addSecureMintJob(i int,
	t *testing.T,
	node Node,
	configuratorAddress common.Address,
	bridgeName string,
) (id int32) {

	// TODO(gg): validate SM spec
	// job, err := streams.ValidatedStreamSpec(spec)
	// require.NoError(t, err)

	addresses, err := node.App.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	spec := getSecureMintJobSpec(configuratorAddress.Hex(), node.KeyBundle.ID(), addresses[0].String(), bridgeName)

	c := node.App.GetConfig()

	job, err := validate.ValidatedOracleSpecToml(testutils.Context(t), c.OCR2(), c.Insecure(), spec, nil)
	require.NoError(t, err)

	err = node.App.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	t.Logf("Added secure mint job spec %s", job.ExternalJobID)

	return job.ID
}

func getSecureMintJobSpec(ocrContractAddress, keyBundleID, transmitterAddress, bridgeName string) string {

	// TODO(gg): check/update EA request/response format
	// TODO(gg): update pluginConfig
	// TODO(gg): is `answer1 [type=any index=0];` correct? Does it actually enable the plugin to come to consensus?

	return fmt.Sprintf(`
type               = "offchainreporting2"
relay              = "evm"
schemaVersion      = 1
pluginType         = "securemint"
name               = "secure mint spec"
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
observationSource  = """
    // data source 1
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];

    ds1 -> ds1_parse -> answer1;

	answer1 [type=any index=0];
"""

allowNoBootstrappers = false

[relayConfig]
chainID = 1337
fromBlock = 1

[pluginConfig]
juelsPerFeeCoinSource = """
		// data source 1
		ds1          [type=bridge name="%s"];
		ds1_parse    [type=jsonparse path="data"];
		ds1_multiply [type=multiply times=1];

		ds1 -> ds1_parse -> ds1_multiply -> answer1;

	answer1 [type=median index=0];
"""
gasPriceSubunitsSource = """
		// data source
		dsp          [type=bridge name="%s"];
		dsp_parse    [type=jsonparse path="data"];
		dsp -> dsp_parse;
"""
[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "1m"
`,
		ocrContractAddress, // contract address
		keyBundleID,        // ocr key bundle id
		transmitterAddress, // transmitter id
		bridgeName,         // bridge name
		bridgeName,         // bridge name
		bridgeName)         // bridge name
}

func createSecureMintBridge(t *testing.T, name string, i int, response por.ExternalAdapterPayload, borm bridges.ORM) (bridgeName string) {
	ctx := testutils.Context(t)
	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// TODO(gg): assert on the EA request format here
		// require.JSONEq(t, `{"meta":{"latestAnswer":"", "updatedAt": ""}}`, string(b))

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		require.NoError(t, err)

		t.Logf("Received request for secure mint bridge %s on node %d: path %s, request body %s", name, i, req.URL.String(), string(body))

		jsonResp, err := json.Marshal(response)
		require.NoError(t, err)

		res.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"data": %s}`, string(jsonResp))
		t.Logf("Responding from secure mint bridge %s on node %d with: %s", name, i, resp)
		_, err = res.Write([]byte(resp))
		require.NoError(t, err)
	}))
	t.Cleanup(func() {
		t.Logf("Closing secure mint bridge %s on node %d with url %s", name, i, bridge.URL)
		bridge.Close()
	})
	t.Logf("Created secure mint bridge %s on node %d with URL %s", name, i, bridge.URL)
	u, _ := url.Parse(bridge.URL)
	bridgeName = fmt.Sprintf("bridge-%s-%d", name, i)
	require.NoError(t, borm.CreateBridgeType(ctx, &bridges.BridgeType{
		Name: bridges.BridgeName(bridgeName),
		URL:  models.WebURL(*u),
	}))

	return bridgeName
}
