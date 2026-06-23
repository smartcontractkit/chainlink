package ocr2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestIntegration_OCR2(t *testing.T) {
	t.Parallel()
	RunTestIntegrationOCR2(t)
}

func TestIntegration_OCR2_ForwarderFlow(t *testing.T) {
	t.Parallel()
	owner, b, ocrContractAddress, ocrContract, nodeConfig := SetupOCR2Contracts(t)

	lggr := logger.TestLogger(t)
	bootstrapNodePort := freeport.GetOne(t)
	bootstrapNode := SetupNodeOCR2(t, owner, bootstrapNodePort, true /* useForwarders */, b, nil, nodeConfig)

	oracles := make([]confighelper2.OracleIdentityExtra, 0, 4)
	transmitters := make([]common.Address, 0, 4)
	forwarderContracts := make([]common.Address, 0, 4)
	kbs := make([]ocr2key.KeyBundle, 0, 4)
	apps := make([]*cltest.TestApplication, 0, 4)
	ports := freeport.GetN(t, 4)
	for i := range uint16(4) {
		node := SetupNodeOCR2(t, owner, ports[i], true /* useForwarders */, b, []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapNode.PeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, nodeConfig)

		// Effective transmitter should be a forwarder not an EOA.
		require.NotEqual(t, node.EffectiveTransmitter, node.Transmitter)

		kbs = append(kbs, node.KeyBundle)
		apps = append(apps, node.App)
		forwarderContracts = append(forwarderContracts, node.EffectiveTransmitter)
		transmitters = append(transmitters, node.Transmitter)

		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  node.KeyBundle.PublicKey(),
				TransmitAccount:   ocrtypes2.Account(node.EffectiveTransmitter.String()),
				OffchainPublicKey: node.KeyBundle.OffchainPublicKey(),
				PeerID:            node.PeerID,
			},
			ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
		})
	}

	blockBeforeConfig := InitOCR2(t, lggr, b, ocrContract, owner, bootstrapNode, oracles, forwarderContracts, transmitters, func(int64) string {
		return fmt.Sprintf(`
type				= "bootstrap"
name				= "bootstrap"
relay				= "evm"
schemaVersion		= 1
forwardingAllowed   = true
contractID			= "%s"
[relayConfig]
chainID 			= 1337
`, ocrContractAddress)
	})

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

	jids := make([]int32, 0, 4)
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
	ctx := t.Context()
	for i := range 4 {
		s := i
		require.NoError(t, apps[i].Start(t.Context()))

		// API speed is > observation timeout set in ContractSetConfigArgsForIntegrationTest
		slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			<-req.Context().Done()
			res.WriteHeader(http.StatusOK)
			_, err := res.Write([]byte(`{"data":10}`))
			assert.NoError(t, err)
		}))
		t.Cleanup(func() {
			slowServers[s].Close()
		})
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
			if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
				metaLock.Lock()
				delete(expectedMeta, m.Meta.LatestAnswer.String())
				metaLock.Unlock()
			}
			res.WriteHeader(http.StatusOK)
			_, err = res.Write([]byte(`{"data":10}`))
			assert.NoError(t, err)
		}))
		t.Cleanup(func() {
			servers[s].Close()
		})
		u, _ := url.Parse(servers[i].URL)
		require.NoError(t, apps[i].BridgeORM().CreateBridgeType(ctx, &bridges.BridgeType{
			Name: bridges.BridgeName(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		}))

		ocrJob, err := validate.ValidatedOracleSpecToml(t.Context(), apps[i].Config.OCR2(), apps[i].Config.Insecure(), fmt.Sprintf(`
type               = "offchainreporting2"
relay              = "evm"
schemaVersion      = 1
pluginType         = "median"
name               = "web oracle spec"
forwardingAllowed  = true
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
    ds2          [type=http method=GET url="%s" timeout="2s"];
    ds2_parse    [type=jsonparse path="data"];
    ds2_multiply [type=multiply times=%d];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
		// data source 1
		ds1          [type=bridge name="%s"];
		ds1_parse    [type=jsonparse path="data"];
		ds1_multiply [type=multiply times=%d];

		// data source 2
		ds2          [type=http method=GET url="%s" timeout="2s"];
		ds2_parse    [type=jsonparse path="data"];
		ds2_multiply [type=multiply times=%d];

		ds1 -> ds1_parse -> ds1_multiply -> answer1;
		ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""
[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "1m"
`, ocrContractAddress, kbs[i].ID(), transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i, fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i), nil)
		require.NoError(t, err)
		err = apps[i].AddJobV2(t.Context(), &ocrJob)
		require.NoError(t, err)
		jids = append(jids, ocrJob.ID)
	}

	// Once all the jobs are added, replay to ensure we have the configSet logs.
	for _, app := range apps {
		require.NoError(t, mustEVMChain(t, app).LogPoller().Replay(t.Context(), blockBeforeConfig.Number().Int64()))
	}

	require.NoError(t, mustEVMChain(t, bootstrapNode.App).LogPoller().Replay(t.Context(), blockBeforeConfig.Number().Int64()))

	// Assert that all the OCR jobs get a run with valid values eventually.
	var wg sync.WaitGroup
	for i := range 4 {
		ic := i
		wg.Go(func() {
			// Want at least 2 runs so we see all the metadata.
			pr := cltest.WaitForPipelineComplete(t, ic, jids[ic], 2, 7, apps[ic].JobORM(), 2*time.Minute, 100*time.Millisecond)
			jb, err := pr[0].Outputs.MarshalJSON()
			if !assert.NoError(t, err) { //nolint:testifylint // require.NoError inside a goroutine is unsafe
				return
			}
			assert.Equal(t, fmt.Appendf(nil, "[\"%d\"]", 10*ic), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1])
		})
	}
	wg.Wait()

	require.Eventually(t, func() bool {
		answer, err := ocrContract.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.String() == "20"
	}, 1*time.Minute, 200*time.Millisecond)

	for _, app := range apps {
		jobs, _, err := app.JobORM().FindJobs(ctx, 0, 1000)
		require.NoError(t, err)
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
	require.Eventually(t, func() bool {
		metaLock.Lock()
		defer metaLock.Unlock()
		return len(expectedMeta) == 0
	}, 1*time.Minute, 200*time.Millisecond)

	// Assert we can read the latest config digest and epoch after a report has been submitted.
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	require.NoError(t, err)
	store := keys.NewStore(keystore.NewEthSigner(apps[0].KeyStore.Eth(), testutils.SimulatedChainID))
	chain := mustEVMChain(t, apps[0])
	ct, err := transmitter.NewOCRContractTransmitter(t.Context(), ocrContractAddress, chain.Client(), contractABI, nil, chain.LogPoller(), lggr, store)
	require.NoError(t, err)
	configDigest, epoch, err := ct.LatestConfigDigestAndEpoch(t.Context())
	require.NoError(t, err)
	details, err := ocrContract.LatestConfigDetails(nil)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(configDigest[:], details.ConfigDigest[:]))
	digestAndEpoch, err := ocrContract.LatestConfigDigestAndEpoch(nil)
	require.NoError(t, err)
	assert.Equal(t, digestAndEpoch.Epoch, epoch)
}
