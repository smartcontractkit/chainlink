package ocr2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

type Node struct {
	App                  *cltest.TestApplication
	PeerID               string
	Transmitter          common.Address
	EffectiveTransmitter common.Address
	KeyBundle            ocr2key.KeyBundle
}

func SetupOCR2Contracts(t *testing.T) (*bind.TransactOpts, *simulated.Backend, common.Address, *ocr2aggregator.OCR2Aggregator) {
	owner := testutils.MustNewSimTransactor(t)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := types.GenesisAlloc{owner.From: {Balance: sb}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := simulated.NewBackend(genesisData, simulated.WithBlockGasLimit(gasLimit))
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b.Client())
	require.NoError(t, err)
	b.Commit()
	accessAddress, _, _, err := testoffchainaggregator2.DeploySimpleWriteAccessController(owner, b.Client())
	require.NoError(t, err, "failed to deploy test access controller contract")
	b.Commit()

	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))
	ocrContractAddress, _, ocrContract, err := ocr2aggregator.DeployOCR2Aggregator(
		owner,
		b.Client(),
		linkTokenAddress, // _link common.Address,
		minAnswer,        // -2**191
		maxAnswer,        // 2**191 - 1
		accessAddress,
		accessAddress,
		9,
		"TEST",
	)
	// Ensure we have finality depth worth of blocks to start.
	for i := 0; i < 20; i++ {
		b.Commit()
	}
	require.NoError(t, err)
	_, err = linkContract.Transfer(owner, ocrContractAddress, big.NewInt(1000))
	require.NoError(t, err)
	b.Commit()
	return owner, b, ocrContractAddress, ocrContract
}

func SetupNodeOCR2(
	t *testing.T,
	owner *bind.TransactOpts,
	port int,
	useForwarder bool,
	b *simulated.Backend,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
) *Node {
	ctx := testutils.Context(t)
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(port)))
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(5 * time.Second)
		c.EVM[0].Transactions.ForwardersEnabled = &useForwarder
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address
	effectiveTransmitter := sendingKeys[0].Address

	// Fund the transmitter address with some ETH
	n, err := b.Client().NonceAt(testutils.Context(t), owner.From, nil)
	require.NoError(t, err)

	tx := cltest.NewLegacyTransaction(
		n, transmitter,
		assets.Ether(1).ToInt(),
		21000,
		assets.GWei(1).ToInt(),
		nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.Client().SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	kb, err := app.GetKeyStore().OCR2().Create(ctx, "evm")
	require.NoError(t, err)

	if useForwarder {
		// deploy a forwarder
		faddr, _, authorizedForwarder, err2 := authorized_forwarder.DeployAuthorizedForwarder(owner, b.Client(), common.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB"), owner.From, common.Address{}, []byte{})
		require.NoError(t, err2)
		b.Commit()

		// set EOA as an authorized sender for the forwarder
		_, err2 = authorizedForwarder.SetAuthorizedSenders(owner, []common.Address{transmitter})
		require.NoError(t, err2)
		b.Commit()

		// add forwarder address to be tracked in db
		forwarderORM := forwarders.NewORM(app.GetDB())
		chainID, err := b.Client().ChainID(testutils.Context(t))
		require.NoError(t, err)
		_, err2 = forwarderORM.CreateForwarder(testutils.Context(t), faddr, ubig.Big(*chainID))
		require.NoError(t, err2)

		effectiveTransmitter = faddr
	}
	return &Node{
		App:                  app,
		PeerID:               p2pKey.PeerID().Raw(),
		Transmitter:          transmitter,
		EffectiveTransmitter: effectiveTransmitter,
		KeyBundle:            kb,
	}
}

func RunTestIntegrationOCR2(t *testing.T) {
	for _, test := range []struct {
		name                string
		chainReaderAndCodec bool
	}{
		{"legacy", false},
		{"chain-reader", true},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			owner, b, ocrContractAddress, ocrContract := SetupOCR2Contracts(t)

			lggr := logger.TestLogger(t)
			bootstrapNodePort := freeport.GetOne(t)
			bootstrapNode := SetupNodeOCR2(t, owner, bootstrapNodePort, false /* useForwarders */, b, nil)

			var (
				oracles      []confighelper2.OracleIdentityExtra
				transmitters []common.Address
				kbs          []ocr2key.KeyBundle
				apps         []*cltest.TestApplication
			)
			ports := freeport.GetN(t, 4)
			for i := 0; i < 4; i++ {
				node := SetupNodeOCR2(t, owner, ports[i], false /* useForwarders */, b, []commontypes.BootstrapperLocator{
					// Supply the bootstrap IP and port as a V2 peer address
					{PeerID: bootstrapNode.PeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
				})

				kbs = append(kbs, node.KeyBundle)
				apps = append(apps, node.App)
				transmitters = append(transmitters, node.Transmitter)

				oracles = append(oracles, confighelper2.OracleIdentityExtra{
					OracleIdentity: confighelper2.OracleIdentity{
						OnchainPublicKey:  node.KeyBundle.PublicKey(),
						TransmitAccount:   ocrtypes2.Account(node.Transmitter.String()),
						OffchainPublicKey: node.KeyBundle.OffchainPublicKey(),
						PeerID:            node.PeerID,
					},
					ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
				})
			}

			blockBeforeConfig := InitOCR2(t, lggr, b, ocrContract, owner, bootstrapNode, oracles, transmitters, transmitters, func(blockNum int64) string {
				return fmt.Sprintf(`
type				= "bootstrap"
name				= "bootstrap"
relay				= "evm"
schemaVersion		= 1
contractID			= "%s"
[relayConfig]
chainID 			= 1337
fromBlock = %d
`, ocrContractAddress, blockNum)
			})

			tick := time.NewTicker(1 * time.Second)
			defer tick.Stop()
			go func() {
				for range tick.C {
					b.Commit()
				}
			}()

			var jids []int32
			var servers, slowServers = make([]*httptest.Server, 4), make([]*httptest.Server, 4)
			// We expect metadata of:
			//  latestAnswer:nil // First call
			//  latestAnswer:0
			//  latestAnswer:10
			//  latestAnswer:20
			//  latestAnswer:30
			var metaLock sync.Mutex
			expectedMeta := map[string]struct{}{
				"0": {}, "10": {}, "20": {}, "30": {},
			}
			returnData := int(10)
			for i := 0; i < 4; i++ {
				s := i
				require.NoError(t, apps[i].Start(testutils.Context(t)))

				// API speed is > observation timeout set in ContractSetConfigArgsForIntegrationTest
				slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					time.Sleep(5 * time.Second)
					var result string
					metaLock.Lock()
					result = fmt.Sprintf(`{"data":%d}`, returnData)
					metaLock.Unlock()
					res.WriteHeader(http.StatusOK)
					t.Logf("Slow Bridge %d returning data:10", s)
					_, err := res.Write([]byte(result))
					assert.NoError(t, err)
				}))
				t.Cleanup(slowServers[s].Close)
				servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					b, err := io.ReadAll(req.Body)
					if !assert.NoError(t, err) {
						res.WriteHeader(http.StatusInternalServerError)
						return
					}
					var m bridges.BridgeMetaDataJSON
					if !assert.NoError(t, json.Unmarshal(b, &m)) {
						res.WriteHeader(http.StatusInternalServerError)
						return
					}
					var result string
					metaLock.Lock()
					result = fmt.Sprintf(`{"data":%d}`, returnData)
					metaLock.Unlock()
					if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
						t.Logf("Bridge %d deleting %s, from request body: %s", s, m.Meta.LatestAnswer, b)
						metaLock.Lock()
						delete(expectedMeta, m.Meta.LatestAnswer.String())
						metaLock.Unlock()
					}
					res.WriteHeader(http.StatusOK)
					_, err = res.Write([]byte(result))
					assert.NoError(t, err)
				}))
				t.Cleanup(servers[s].Close)
				u, _ := url.Parse(servers[i].URL)
				bridgeName := fmt.Sprintf("bridge%d", i)
				require.NoError(t, apps[i].BridgeORM().CreateBridgeType(testutils.Context(t), &bridges.BridgeType{
					Name: bridges.BridgeName(bridgeName),
					URL:  models.WebURL(*u),
				}))

				var chainReaderSpec string
				if test.chainReaderAndCodec {
					chainReaderSpec = `
[relayConfig.chainReader.contracts.median]
contractPollingFilter.genericEventNames = ["LatestRoundRequested"]

contractABI = '''
[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "requester",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "configDigest",
        "type": "bytes32"
      },
      {
        "indexed": false,
        "internalType": "uint32",
        "name": "epoch",
        "type": "uint32"
      },
      {
        "indexed": false,
        "internalType": "uint8",
        "name": "round",
        "type": "uint8"
      }
    ],
    "name": "RoundRequested",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "latestTransmissionDetails",
    "outputs": [
      {
        "internalType": "bytes32",
        "name": "configDigest",
        "type": "bytes32"
      },
      {
        "internalType": "uint32",
        "name": "epoch",
        "type": "uint32"
      },
      {
        "internalType": "uint8",
        "name": "round",
        "type": "uint8"
      },
      {
        "internalType": "int192",
        "name": "latestAnswer_",
        "type": "int192"
      },
      {
        "internalType": "uint64",
        "name": "latestTimestamp_",
        "type": "uint64"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]
'''

[relayConfig.chainReader.contracts.median.configs]
LatestRoundRequested = '''
{
  "chainSpecificName": "RoundRequested",
  "readType": "event"
}
'''
LatestTransmissionDetails = '''
{
  "chainSpecificName": "latestTransmissionDetails",
  "outputModifications": [
    {
      "Fields": [
        "LatestTimestamp_"
      ],
      "type": "epoch to time"
    },
    {
      "Fields": {
        "LatestAnswer_": "LatestAnswer",
        "LatestTimestamp_": "LatestTimestamp"
      },
      "type": "rename"
    }
  ]
}
'''

[relayConfig.codec.configs.MedianReport]
typeABI = '''
[
  {
    "Name": "Timestamp",
    "Type": "uint32"
  },
  {
    "Name": "Observers",
    "Type": "bytes32"
  },
  {
    "Name": "Observations",
    "Type": "int192[]"
  },
  {
    "Name": "JuelsPerFeeCoin",
    "Type": "int192"
  }
]
'''
`
				}
				ocrJob, err := validate.ValidatedOracleSpecToml(testutils.Context(t), apps[i].Config.OCR2(), apps[i].Config.Insecure(), fmt.Sprintf(`
type               = "offchainreporting2"
relay              = "evm"
schemaVersion      = 1
pluginType         = "median"
name               = "web oracle spec"
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
observationSource  = """
    // data source 1
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1_multiply [type=multiply times=%d];

    // data source 2
    ds2          [type=http method=GET url="%s"];
    ds2_parse    [type=jsonparse path="data"];
    ds2_multiply [type=multiply times=%d];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""

[relayConfig]
chainID = 1337
fromBlock = %d
%s

[pluginConfig]
juelsPerFeeCoinSource = """
		// data source 1
		ds1          [type=bridge name="%s"];
		ds1_parse    [type=jsonparse path="data"];
		ds1_multiply [type=multiply times=%d];

		// data source 2
		ds2          [type=http method=GET url="%s"];
		ds2_parse    [type=jsonparse path="data"];
		ds2_multiply [type=multiply times=%d];

		ds1 -> ds1_parse -> ds1_multiply -> answer1;
		ds2 -> ds2_parse -> ds2_multiply -> answer1;

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
`, ocrContractAddress, kbs[i].ID(), transmitters[i], bridgeName, i, slowServers[i].URL, i, blockBeforeConfig.Number().Int64(), chainReaderSpec, bridgeName, i, slowServers[i].URL, i, bridgeName), nil)
				require.NoError(t, err)
				err = apps[i].AddJobV2(testutils.Context(t), &ocrJob)
				require.NoError(t, err)
				jids = append(jids, ocrJob.ID)
			}

			// Watch for OCR2AggregatorTransmitted events
			start := uint64(0)
			txEvents := make(chan *ocr2aggregator.OCR2AggregatorTransmitted)
			_, err := ocrContract.WatchTransmitted(&bind.WatchOpts{Start: &start, Context: testutils.Context(t)}, txEvents)
			require.NoError(t, err)
			newTxEvents := make(chan *ocr2aggregator.OCR2AggregatorNewTransmission)
			_, err = ocrContract.WatchNewTransmission(&bind.WatchOpts{Start: &start, Context: testutils.Context(t)}, newTxEvents, []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			require.NoError(t, err)

			go func() {
				var newTxEvent *ocr2aggregator.OCR2AggregatorNewTransmission
				select {
				case txEvent := <-txEvents:
					t.Logf("txEvent: %v", txEvent)
					if newTxEvent != nil {
						assert.Equal(t, uint64(txEvent.Epoch), newTxEvent.EpochAndRound.Uint64())
					}
				case newTxEvent = <-newTxEvents:
					t.Logf("newTxEvent: %v", newTxEvent)
				}
			}()

			ctx := testutils.Context(t)
			for trial := 0; trial < 2; trial++ {
				var retVal int

				metaLock.Lock()
				returnData = 10 * (trial + 1)
				retVal = returnData
				for i := 0; i < 4; i++ {
					expectedMeta[strconv.Itoa(returnData*i)] = struct{}{}
				}
				metaLock.Unlock()

				// Assert that all the OCR jobs get a run with valid values eventually.
				var wg sync.WaitGroup
				for i := 0; i < 4; i++ {
					ic := i
					wg.Add(1)
					go func() {
						defer wg.Done()
						completedRuns, err2 := apps[ic].JobORM().FindPipelineRunIDsByJobID(ctx, jids[ic], 0, 1000)
						if !assert.NoError(t, err2) {
							return
						}
						// Want at least 2 runs so we see all the metadata.
						pr := cltest.WaitForPipelineComplete(t, ic, jids[ic], len(completedRuns)+2, 7, apps[ic].JobORM(), 2*time.Minute, 5*time.Second)
						jb, err2 := pr[0].Outputs.MarshalJSON()
						if !assert.NoError(t, err2) {
							return
						}
						assert.Equal(t, []byte(fmt.Sprintf("[\"%d\"]", retVal*ic)), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1])
					}()
				}
				wg.Wait()

				// Trail #1: 4 oracles reporting 0, 10, 20, 30. Answer should be 20 (results[4/2]).
				// Trial #2: 4 oracles reporting 0, 20, 40, 60. Answer should be 40 (results[4/2]).
				gomega.NewGomegaWithT(t).Eventually(func() string {
					answer, err2 := ocrContract.LatestAnswer(nil)
					require.NoError(t, err2)
					return answer.String()
				}, 1*time.Minute, 200*time.Millisecond).Should(gomega.Equal(strconv.Itoa(2 * retVal)))

				for _, app := range apps {
					jobs, _, err2 := app.JobORM().FindJobs(ctx, 0, 1000)
					require.NoError(t, err2)
					// No spec errors
					for _, j := range jobs {
						ignore := 0
						for i := range j.JobSpecErrors {
							// Non-fatal timing related error, ignore for testing.
							if strings.Contains(j.JobSpecErrors[i].Description, "leader's phase conflicts tGrace timeout") {
								ignore++
							}
						}
						require.Len(t, j.JobSpecErrors, ignore)
					}
				}
				em := map[string]struct{}{}
				metaLock.Lock()
				maps.Copy(em, expectedMeta)
				metaLock.Unlock()
				assert.Empty(t, em, "expected metadata %v", em)

				t.Logf("======= Summary =======")
				roundID, err2 := ocrContract.LatestRound(nil)
				require.NoError(t, err2)
				for i := 0; i <= int(roundID.Int64()); i++ {
					roundData, err3 := ocrContract.GetRoundData(nil, big.NewInt(int64(i)))
					require.NoError(t, err3)
					t.Logf("RoundId: %d, AnsweredInRound: %d, Answer: %d, StartedAt: %v, UpdatedAt: %v", roundData.RoundId, roundData.AnsweredInRound, roundData.Answer, roundData.StartedAt, roundData.UpdatedAt)
				}

				expectedAnswer := big.NewInt(2 * int64(retVal))

				// Assert we can read the latest config digest and epoch after a report has been submitted.
				contractABI, err2 := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
				require.NoError(t, err2)
				apps[0].GetRelayers().LegacyEVMChains().Slice()
				ct, err2 := evm.NewOCRContractTransmitter(testutils.Context(t), ocrContractAddress, apps[0].GetRelayers().LegacyEVMChains().Slice()[0].Client(), contractABI, nil, apps[0].GetRelayers().LegacyEVMChains().Slice()[0].LogPoller(), lggr, apps[0].KeyStore.Eth())
				require.NoError(t, err2)
				configDigest, epoch, err2 := ct.LatestConfigDigestAndEpoch(testutils.Context(t))
				require.NoError(t, err2)
				details, err2 := ocrContract.LatestConfigDetails(nil)
				require.NoError(t, err2)
				assert.True(t, bytes.Equal(configDigest[:], details.ConfigDigest[:]))
				digestAndEpoch, err2 := ocrContract.LatestConfigDigestAndEpoch(nil)
				require.NoError(t, err2)
				assert.Equal(t, digestAndEpoch.Epoch, epoch)
				latestTransmissionDetails, err2 := ocrContract.LatestTransmissionDetails(nil)
				require.NoError(t, err2)
				assert.Equal(t, expectedAnswer, latestTransmissionDetails.LatestAnswer)
				require.NoError(t, err2)
				newTransmissionEvents, err2 := ocrContract.FilterTransmitted(&bind.FilterOpts{Start: 0, End: nil})
				require.NoError(t, err2)
				for newTransmissionEvents.Next() {
					assert.Equal(t, 3, newTransmissionEvents.Event.Epoch)
				}
			}
		})
	}
}

func InitOCR2(t *testing.T, lggr logger.Logger, b *simulated.Backend,
	ocrContract *ocr2aggregator.OCR2Aggregator,
	owner *bind.TransactOpts,
	bootstrapNode *Node,
	oracles []confighelper2.OracleIdentityExtra,
	transmitters []common.Address,
	payees []common.Address,
	specFn func(int64) string,
) (
	blockBeforeConfig *types.Block,
) {
	lggr.Debugw("Setting Payees on OraclePlugin Contract", "transmitters", payees)
	_, err := ocrContract.SetPayees(
		owner,
		transmitters,
		payees,
	)
	require.NoError(t, err)
	b.Commit()
	blockBeforeConfig, err = b.Client().BlockByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	signers, effectiveTransmitters, threshold, _, encodedConfigVersion, encodedConfig, err := confighelper2.ContractSetConfigArgsForEthereumIntegrationTest(
		oracles,
		1,
		1000000000/100, // threshold PPB
	)
	require.NoError(t, err)

	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer)
	require.NoError(t, err)

	lggr.Debugw("Setting Config on Oracle Contract",
		"signers", signers,
		"transmitters", transmitters,
		"effectiveTransmitters", effectiveTransmitters,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", encodedConfigVersion,
	)
	_, err = ocrContract.SetConfig(
		owner,
		signers,
		effectiveTransmitters,
		threshold,
		onchainConfig,
		encodedConfigVersion,
		encodedConfig,
	)
	require.NoError(t, err)
	b.Commit()

	err = bootstrapNode.App.Start(testutils.Context(t))
	require.NoError(t, err)

	chainSet := bootstrapNode.App.GetRelayers().LegacyEVMChains()
	require.NotNil(t, chainSet)
	ocrJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(specFn(blockBeforeConfig.Number().Int64()))
	require.NoError(t, err)
	err = bootstrapNode.App.AddJobV2(testutils.Context(t), &ocrJob)
	require.NoError(t, err)
	return
}

func ptr[T any](v T) *T { return &v }
