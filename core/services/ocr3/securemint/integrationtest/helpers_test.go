package integrationtest

import (
	"crypto/ed25519"
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
	sm_ea "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/ea"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"
)

type node struct {
	app          chainlink.Application
	clientPubKey credentials.StaticSizedPublicKey
	keyBundle    ocr2key.KeyBundle
	observedLogs *observer.ObservedLogs
}

func (node *node) addBootstrapJob(t *testing.T, spec string) *job.Job {
	job, err := ocrbootstrap.ValidatedBootstrapSpecToml(spec)
	require.NoError(t, err)
	err = node.app.AddJobV2(testutils.Context(t), &job)
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

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, _ *chainlink.Secrets) {
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
		c.OCR2.ContractPollInterval = commonconfig.MustNewDuration(500 * time.Millisecond)

		// [P2P]
		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.TraceLogging = ptr(true)

		// [P2P.V2]
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

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

	app = cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend, p2pKey, ocr2kb, csaKey, lggr.Named(dbName))
	err := app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb, observedLogs
}

func ptr[T any](t T) *T { return &t }

func createSecureMintBootstrapJob(t *testing.T, bootstrapNode node, configuratorAddress common.Address, chainID, fromBlock string) *job.Job {
	return bootstrapNode.addBootstrapJob(t, fmt.Sprintf(`
			type                              = "bootstrap"
			relay                             = "evm"
			schemaVersion                     = 1
			name                              = "bootstrap-secure-mint"
			contractID                        = "%s"
			contractConfigTrackerPollInterval = "1s"
			contractConfigConfirmations       = 1

			[relayConfig]
			chainID                           = %s
			fromBlock                         = %s
			providerType                      = "securemint"
			lloDonID                          = 1
			lloConfigMode 					  = "bluegreen"`, // Using lloConfigMode 'bluegreen' since otherwise LLO config poller won't work
		configuratorAddress.Hex(),
		chainID,
		fromBlock),
	)
}

func addSecureMintOCRJobs(
	t *testing.T,
	nodes []node,
	configuratorAddress common.Address,
) (jobIDs map[int]int32) {
	// node idx => job id
	jobIDs = make(map[int]int32)

	// Create one bridge and one SM Feed OCR job on each node
	for i, node := range nodes {
		name := "securemint-ea"

		bmBridge := createSecureMintBridge(t, name, i, node.app.BridgeORM())
		t.Logf("Created secure mint bridge %s on node %d", bmBridge, i)

		jobID := addSecureMintJob(
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

func addSecureMintJob(
	t *testing.T,
	node node,
	configuratorAddress common.Address,
	bridgeName string,
) (id int32) {

	spec := getSecureMintJobSpec(t, configuratorAddress.Hex(), node.keyBundle.ID(), node.clientPubKey[:], bridgeName)

	c := node.app.GetConfig()
	job, err := validate.ValidatedOracleSpecToml(testutils.Context(t), c.OCR2(), c.Insecure(), spec, nil)
	require.NoError(t, err)

	err = node.app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	t.Logf("Added secure mint job spec %s", job.ExternalJobID)

	return job.ID
}

func getSecureMintJobSpec(t *testing.T, ocrContractAddress, keyBundleID string, publicKey ed25519.PublicKey, bridgeName string) string {

	t.Logf("Using transmitter address %x for job", publicKey)

	return fmt.Sprintf(`
			type                              = "offchainreporting2"
			relay                             = "evm"
			schemaVersion                     = 1
			pluginType                        = "securemint"
			name                              = "secure mint spec"
			contractID                        = "%s"
			ocrKeyBundleID                    = "%s"
			transmitterID                     = "%x"
			contractConfigConfirmations       = 1
			contractConfigTrackerPollInterval = "1s"
			observationSource  = """
				// data source 1
				ds1          [type=bridge name="%s" requestData=<{ "data": $(ea_request) }>];
				ds1_parse    [type=jsonparse path="data"];

				ds1 -> ds1_parse -> answer1;

				answer1 [type=any index=0];
			"""

			allowNoBootstrappers              = false

			[relayConfig]
			chainID                           = 1337
			fromBlock                         = 1
			providerType                      = "securemint"
			lloDonID                          = 1
			lloConfigMode                     = "bluegreen"

			[pluginConfig]
			maxChains                         = 5
			token                             = "btc"
			reserves                          = "custom"
			chainSelectors                    = ["8953668971247136127", "729797994450396300"]
		`, // Using lloConfigMode 'bluegreen' since otherwise LLO config poller won't work
		ocrContractAddress, // contract address
		keyBundleID,        // ocr key bundle id
		publicKey,          // transmitter id
		bridgeName)         // bridge name
}

// Based on https://chainlink-core.slack.com/archives/C090PQH50M6/p1749483857095389?thread_ts=1749482941.061609&cid=C090PQH50M6
func createSecureMintBridge(t *testing.T, name string, i int, borm bridges.ORM) (bridgeName string) {
	ctx := testutils.Context(t)

	initialResponse := sm_ea.Response{
		Mintables: map[string]sm_ea.MintableInfo{},
		LatestBlocks: map[string]uint64{
			"8953668971247136127": 40, // "bitcoin-testnet-rootstock"
			"729797994450396300":  5,  // "telos-evm-testnet"
		},
		ReserveInfo: sm_ea.ReserveInfo{
			ReserveAmount: "1000",
			Timestamp:     time.Now().UnixMilli(),
		},
	}
	jsonInitialResp, err := json.Marshal(initialResponse)
	require.NoError(t, err)

	fullResponse := sm_ea.Response{
		Mintables: map[string]sm_ea.MintableInfo{
			"8953668971247136127": { // "bitcoin-testnet-rootstock"
				Block:    uint64(40),
				Mintable: "10",
			},
			"729797994450396300": { // "telos-evm-testnet"
				Block:    uint64(5),
				Mintable: "25",
			},
		},
		LatestBlocks: map[string]uint64{
			"8953668971247136127": 42, // "bitcoin-testnet-rootstock"
			"729797994450396300":  7,  // "telos-evm-testnet"
		},
		ReserveInfo: sm_ea.ReserveInfo{
			ReserveAmount: "500",
			Timestamp:     time.Now().UnixMilli(),
		},
	}
	jsonFullResponse, err := json.Marshal(fullResponse)
	require.NoError(t, err)

	//nolint:testifylint // allow require.NoError in the http server
	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		require.NoError(t, err)
		t.Logf("Received request for secure mint bridge %s on node %d: path %s, request body %s", name, i, req.URL.String(), string(body))

		// Parse the request body into a map to extract the 'data' field
		var requestMap map[string]any
		err = json.Unmarshal(body, &requestMap)
		require.NoError(t, err, "Failed to parse request body as map for bridge %s on node %d", name, i)

		dataField, exists := requestMap["data"]
		require.True(t, exists, "Request body should contain 'data' field for bridge %s on node %d", name, i)

		// Marshal the data field back to JSON and parse as ea.Request
		dataBytes, err := json.Marshal(dataField)
		require.NoError(t, err, "Failed to marshal data field for bridge %s on node %d", name, i)
		var eaRequest sm_ea.Request
		err = json.Unmarshal(dataBytes, &eaRequest)
		require.NoError(t, err, "Failed to parse request body as ea.Request for bridge %s on node %d", name, i)

		// Validate the parsed ea.Request
		assert.Equal(t, "btc", eaRequest.Token, "Token should be 'eth'")
		assert.Equal(t, "custom", eaRequest.Reserves, "Reserves should be 'platform'")

		// Return initial EA response if empty request (first round)
		if len(eaRequest.SupplyChains) == 0 && len(eaRequest.SupplyChainBlocks) == 0 {
			t.Logf("Received empty supply chains for secure mint bridge %s on node %d, returning initial response", name, i)
			res.WriteHeader(http.StatusOK)
			_, err = res.Write(fmt.Appendf(nil, `{"data": %s}`, string(jsonInitialResp)))
			require.NoError(t, err)
			return
		}

		// Validate non-empty request
		assert.Contains(t, eaRequest.SupplyChains, "8953668971247136127", "Supply chains should contain bitcoin-testnet-rootstock")
		assert.Contains(t, eaRequest.SupplyChains, "729797994450396300", "Supply chains should contain telos-evm-testnet")
		assert.Len(t, eaRequest.SupplyChains, 2, "Should have exactly 2 supply chains")

		assert.Len(t, eaRequest.SupplyChainBlocks, 2, "Should have exactly 2 supply chain blocks")
		assert.GreaterOrEqual(t, eaRequest.SupplyChainBlocks[0], uint64(0), "Supply chain block should be at least 0")
		assert.GreaterOrEqual(t, eaRequest.SupplyChainBlocks[1], uint64(0), "Supply chain block should be at least 0")

		// Return full EA response with mintable amounts
		res.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"data": %s}`, string(jsonFullResponse))
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

// secureMintReport mimics por.PorReport in the securemint plugin.
type secureMintReport struct {
	ConfigDigest ocr2types.ConfigDigest `json:"configDigest"`
	SeqNr        uint64                 `json:"seqNr"`
	Block        uint64                 `json:"block"`
	Mintable     *big.Int               `json:"mintable"`
}

// secureMintOffchainConfig mimics por.PorOffchainConfig in the securemint plugin.
type secureMintOffchainConfig struct {
	MaxChains uint32 // The maximum number of chains that can be tracked by the external adapter.
}

func (c *secureMintOffchainConfig) Serialize() ([]byte, error) {
	return json.Marshal(c)
}
