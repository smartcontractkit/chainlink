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
	sm_ea "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/external_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
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

		bmBridge := createSecureMintBridge(t, name, i, node.App.BridgeORM())
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

	addresses, err := node.App.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	c := node.App.GetConfig()

	spec := getSecureMintJobSpec(configuratorAddress.Hex(), node.KeyBundle.ID(), addresses[0].String(), bridgeName)

	job, err := validate.ValidatedOracleSpecToml(testutils.Context(t), c.OCR2(), c.Insecure(), spec, nil)
	require.NoError(t, err)

	err = node.App.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	t.Logf("Added secure mint job spec %s", job.ExternalJobID)

	return job.ID
}

func getSecureMintJobSpec(ocrContractAddress, keyBundleID, transmitterAddress, bridgeName string) string {

	// TODO(gg): update pluginConfig

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
    ds1          [type=bridge name="%s" requestData=<{ "data": $(ea_request) }>];
    ds1_parse    [type=jsonparse path="data"];

    ds1 -> ds1_parse -> answer1;

	answer1 [type=any index=0];
"""

allowNoBootstrappers = false

[relayConfig]
chainID = 1337
fromBlock = 1

[pluginConfig]
maxChains = 5
`,
		ocrContractAddress, // contract address
		keyBundleID,        // ocr key bundle id
		transmitterAddress, // transmitter id
		bridgeName)         // bridge name
}

//https://chainlink-core.slack.com/archives/C090PQH50M6/p1749483857095389?thread_ts=1749482941.061609&cid=C090PQH50M6:
/**
Input
{
    "data": {
        "token": "eth",
        "reserves": "platform",
        "supplyChains": [
            "5009297550715157269"
        ],
        "supplyChainBlocks": [
            0
        ]
    }
}
Output
{
    "data": {
        "mintables": {
            "5009297550715157269": {
                "mintable": "0",
                "block": 0
            }
        },
        "reserveInfo": {
            "reserveAmount": "10332550000000000000000",
            "timestamp": 1749483841486
        },
        "latestRelevantBlocks": {
            "5009297550715157269": 22667990
        },
        "supplyDetails": {
            "supply": "47550052000000000000000000",
            "premint": "0",
            "chains": {
                "5009297550715157269": {
                    "latest_block": 22667990,
                    "response_block": 0,
                    "request_block": 22667990,
                    "mintable": "0",
                    "token_supply": "44153737311060787567559446",
                    "token_native_mint": "0",
                    "token_ccip_mint": "1637953921482741588493980",
                    "token_ccip_burn": "5034268610421954020934534",
                    "token_pre_mint": "0",
                    "aggregate_pre_mint": false
                }
            }
        }
    },
    "statusCode": 200,
    "result": 0,
    "timestamps": {
        "providerDataRequestedUnixMs": 1749483841817,
        "providerDataReceivedUnixMs": 1749483841984
    },
    "meta": {
        "adapterName": "SECURE_MINT",
        "metrics": {
            "feedId": "{\"token\":\"eth\",\"reserves\":\"platform\",\"supplyChains\":[\"5009297550715157269\"],\"supplyChainBlocks\":[0]}"
        }
    }
}
*/

func createSecureMintBridge(t *testing.T, name string, i int, borm bridges.ORM) (bridgeName string) {
	ctx := testutils.Context(t)

	initialResponse := sm_ea.EAResponse{
		Mintables: map[string]sm_ea.MintableInfo{},
		LatestRelevantBlocks: map[string]uint64{
			"8953668971247136127": 5, // "bitcoin-testnet-rootstock"
			"729797994450396300":  5, // "telos-evm-testnet"
		},
		ReserveInfo: sm_ea.ReserveInfo{
			ReserveAmount: "1000",
			Timestamp:     time.Now().UnixMilli(),
		},
	}
	jsonInitialResp, err := json.Marshal(initialResponse)
	require.NoError(t, err)

	laterResponse := sm_ea.EAResponse{
		Mintables: map[string]sm_ea.MintableInfo{
			"8953668971247136127": sm_ea.MintableInfo{
				Block:    uint64(5),
				Mintable: "10",
			},
			"729797994450396300": sm_ea.MintableInfo{
				Block:    uint64(5),
				Mintable: "25",
			},
		},
		LatestRelevantBlocks: map[string]uint64{
			"8953668971247136127": 8, // "bitcoin-testnet-rootstock"
			"729797994450396300":  7, // "telos-evm-testnet"
		},
		ReserveInfo: sm_ea.ReserveInfo{
			ReserveAmount: "500",
			Timestamp:     time.Now().UnixMilli(),
		},
	}
	jsonLaterResp, err := json.Marshal(laterResponse)
	require.NoError(t, err)

	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// TODO(gg): assert on the EA request format here
		// require.JSONEq(t, `{"meta":{"latestAnswer":"", "updatedAt": ""}}`, string(b))

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		require.NoError(t, err)

		// 		ds1 [type=bridge name="%s" timeout=0 requestData=<{"data": {"address": "0x1234"}}>]

		// ds1 [type=bridge name=\"bridge-api0\" requestData="{\\\"data\\": {\\\"from\\\":\\\"LINK\\\",\\\"to\\\":\\\"ETH\\\"}}"];

		//     submit [type=bridge name="substrate-adapter1" requestData=<{ "value": $(parse) }>]

		// 		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// 			b, err := io.ReadAll(req.Body)
		// 			require.NoError(t, err)
		// 			var m bridges.BridgeMetaDataJSON
		// 			require.NoError(t, json.Unmarshal(b, &m))
		// 			if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
		// 				metaLock.Lock()
		// 				delete(expectedMeta, m.Meta.LatestAnswer.String())
		// 				metaLock.Unlock()
		// 			}
		// 			res.WriteHeader(http.StatusOK)
		// 			_, err = res.Write([]byte(`{"data":10}`))
		// 			require.NoError(t, err)

		t.Logf("Received request for secure mint bridge %s on node %d: path %s, request body %s", name, i, req.URL.String(), string(body))

		// First parse the request body into a map to extract the data field
		var requestMap map[string]any
		err = json.Unmarshal(body, &requestMap)
		require.NoError(t, err, "Failed to parse request body as map for bridge %s on node %d", name, i)

		// Extract the data field
		dataField, exists := requestMap["data"]
		require.True(t, exists, "Request body should contain 'data' field for bridge %s on node %d", name, i)

		// Marshal the data field back to JSON and parse as EARequest
		dataBytes, err := json.Marshal(dataField)
		require.NoError(t, err, "Failed to marshal data field for bridge %s on node %d", name, i)
		var eaRequest sm_ea.EARequest
		err = json.Unmarshal(dataBytes, &eaRequest)
		require.NoError(t, err, "Failed to parse request body as EARequest for bridge %s on node %d", name, i)

		// Assert on the parsed EARequest
		assert.Equal(t, "eth", eaRequest.Token, "Token should be 'eth'")
		assert.Equal(t, "platform", eaRequest.Reserves, "Reserves should be 'platform'")

		if len(eaRequest.SupplyChains) == 0 && len(eaRequest.SupplyChains) == 0 {
			t.Logf("Received empty supply chains for secure mint bridge %s on node %d, returning initial response", name, i)
			res.WriteHeader(http.StatusOK)
			_, err = res.Write([]byte(fmt.Sprintf(`{"data": %s}`, string(jsonInitialResp))))
			require.NoError(t, err)
			return
		}

		assert.Contains(t, eaRequest.SupplyChains, "8953668971247136127", "Supply chains should contain bitcoin-testnet-rootstock")
		assert.Contains(t, eaRequest.SupplyChains, "729797994450396300", "Supply chains should contain telos-evm-testnet")
		assert.Len(t, eaRequest.SupplyChains, 2, "Should have exactly 2 supply chains")

		assert.Len(t, eaRequest.SupplyChainBlocks, 2, "Should have exactly 2 supply chain blocks")
		assert.GreaterOrEqual(t, eaRequest.SupplyChainBlocks[0], uint64(5), "Supply chain block should be at least 5 (based on initial EA response)")
		assert.GreaterOrEqual(t, eaRequest.SupplyChainBlocks[1], uint64(5), "Supply chain block should be at least 5 (based on initial EA response)")

		// {
		//     "data": {
		//         "token": "eth",
		//         "reserves": "platform",
		//         "supplyChains": [
		//             "5009297550715157269"
		//         ],
		//         "supplyChainBlocks": [
		//             0
		//         ]
		//     }
		// }

		// if body == nil || string(body) == `{"data":{"token":"eth","reserves":"platform"}}` {
		// 	t.Logf("Received empty request body for secure mint bridge %s on node %d, returning initial response", name, i)
		// 	res.WriteHeader(http.StatusOK)
		// 	_, err = res.Write([]byte(fmt.Sprintf(`{"data": %s}`, string(jsonInitialResp))))
		// 	require.NoError(t, err)
		// 	return
		// }

		// Check if the request body contains the expected data

		// assert.JSONEqf(t, `{"data":{"token":"eth","reserves":"platform","supplyChains":["8953668971247136127", "729797994450396300"],"supplyChainBlocks":[5, 5]}}`, string(body),
		// 	"Request body does not match empty body or expected format for secure mint bridge %s on node %d", name, i)

		res.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"data": %s}`, string(jsonLaterResp))
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
