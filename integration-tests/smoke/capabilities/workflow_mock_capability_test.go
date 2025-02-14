package capabilities_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand/v2"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	job2 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	capabilities2 "github.com/smartcontractkit/chainlink/integration-tests/smoke/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"

	pb2 "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

func TestKeystoneWithOCR3WorkflowAndMockCapabilities(t *testing.T) {
	testLogger := framework.L

	// Load test configuration√
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateInputsAndEnvVars(t, in)

	pkey := os.Getenv("PRIVATE_KEY")

	// Create a new blockchain network and Seth client to interact with it
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err, "failed to create seth client")

	chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
	require.NoError(t, err, "failed to get chain selector for chain id %d", sc.Cfg.Network.ChainID)

	// Start job distributor
	jdOutput := startJobDistributor(t, in)

	// Deploy the DONs
	nodeOutputs := make([]*WrappedNodeOutput, len(in.NodeSets))
	for i, don := range in.NodeSets {
		nodeOutputs[i] = startSingleNodeSet(t, don, bc)
	}

	// Prepare the chainlink/deployment environment, which also configures chains for nodes and job distributor
	ctfEnv, dons := buildChainlinkDeploymentEnv(t, jdOutput, bc, sc, nodeOutputs...)

	// Fund the nodes
	fundNodes(t, dons, sc)

	donTopology := buildDONTopology(t, in, dons, nodeOutputs)
	workflowDONID := mustOneDONTopologyWithFlag(t, donTopology, WorkflowDON).ID

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability)
	keystoneContractSet := deployKeystoneContracts(t, testLogger, ctfEnv, chainSelector)

	// Deploy and pre-configure workflow registry contract (using only workflow DON id)
	workflowRegistryAddr := prepareWorkflowRegistry(t, testLogger, ctfEnv, chainSelector, sc, workflowDONID)

	// Deploy and configure Keystone Feeds Consumer contract
	feedsConsumerAddress := prepareFeedsConsumer(t, testLogger, ctfEnv, chainSelector, sc, keystoneContractSet.Forwarder.Address(), in.WorkflowConfig.WorkflowName)

	donTopology, ctfEnv = configureNodes(t, testLogger, in, donTopology, ctfEnv, jdOutput, bc, keystoneContractSet, workflowRegistryAddr)

	workflowNodes := DONTopologyWithFlag(donTopology, WorkflowDON)
	createWorkflowJob(t, ctfEnv, workflowNodes[0].DON.Nodes[1:], sc)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.PriceProvider.FeedID, in.WorkflowConfig.WorkflowName, feedsConsumerAddress.Hex(), keystoneContractSet.Forwarder.Address().Hex())

			logDir := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name())

			err := os.RemoveAll(logDir)
			if err != nil {
				testLogger.Error().Err(err).Msg("failed to remove log directory")
				return
			}

			_, err = framework.SaveContainerLogs(logDir)
			if err != nil {
				testLogger.Error().Err(err).Msg("failed to save container logs")
				return
			}

			printTestDebug(t, testLogger, donTopology, bc.Nodes[0].HostWSUrl)
		}
	})

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO make it fluent!
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the workflow DON and contracts
	configureContracts(t, ctfEnv, donTopology, chainSelector)

	// It can take a while before the first report is produced, particularly on CI.
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	testLogger.Info().Msg("Connecting to mock capabilities...")
	mocksClient := newCapProxyClient()
	require.NoError(t, mocksClient.connectAll([]int{13401, 13402, 13403, 13404})) //Capability don ports
	testLogger.Info().Msg("Creating streams-trigger@1.1.0 capabilities...")
	require.NoError(t, mocksClient.CreateCapability(capabilities2.CapabilityInfo{
		ID:             "streams-trigger@1.1.0",
		CapabilityType: capabilities2.CapabilityType_Trigger,
		Description:    "mock streams-trigger capability",
		DON:            nil,
		IsLocal:        true,
	}))
	testLogger.Info().Msg("Creating write_ethereum-testnet-sepolia-arbitrum-1@1.0.0 capabilities...")
	require.NoError(t, mocksClient.CreateCapability(capabilities2.CapabilityInfo{
		ID:             "write_ethereum-testnet-sepolia-arbitrum-1@1.0.0",
		CapabilityType: capabilities2.CapabilityType_Target,
		Description:    "mock write_ethereum-testnet-sepolia-arbitrum-1",
		DON:            nil,
		IsLocal:        true,
	}))

	testLogger.Info().Msg("Hooking into mock executable capabilities")
	require.NoError(t, mocksClient.HookExecutables(testLogger))

	testLogger.Info().Msg("Waiting for feed to update...")
	startTime := time.Now()

	go sendReports(ctx, t, mocksClient, testLogger, 4) //TODO @george-dorin: Fix me!!!

	for {
		select {
		case <-ctx.Done():
			testLogger.Error().Msgf("feed did not update, timeout after %s", timeout)
			t.FailNow()
		case <-time.After(10 * time.Second):
			elapsed := time.Since(startTime).Round(time.Second)
			require.NoError(t, err, "failed to get price from Keystone Consumer contract")

			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)

		}
	}
}

func sendReports(ctx context.Context, t *testing.T, capProxy *capProxy, lggr zerolog.Logger, nrOfNodes int) {
	keyBundles := make([]ocr2key.KeyBundle, 0)
	for range nrOfNodes {
		b, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err, "cannot create key bundle")
		keyBundles = append(keyBundles, b)
	}

	for {
		select {
		case <-ctx.Done():
			lggr.Error().Msg("context canceled")
		case <-time.After(10 * time.Second):
			price := big.NewInt(int64(rand.IntN(100)))
			r := createFeedReport(t, price, time.Now().UnixMilli(), "0x000351de403f638036014add21a5abd5f464bf21d11aa356dfc6dbe4e2384e4e", keyBundles)

			out := datastreams.StreamsTriggerEvent{
				Payload:   []datastreams.FeedReport{*r},
				Metadata:  datastreams.Metadata{},
				Timestamp: time.Now().Unix(),
			}
			outputsv, err := values.WrapMap(out)
			require.NoError(t, err)
			outputBytes, err := mapToBytes(outputsv)
			require.NoError(t, err, "cannot convert payload to bytes")

			require.NoError(t, capProxy.SendTrigger("streams-trigger@1.1.0", uuid.New().String(), outputBytes))
		}
	}
}

func createWorkflowJob(t *testing.T, ctfEnv *deployment.Environment, nodes []devenv.Node, sc *seth.Client) {

	workflowJobSpec := fmt.Sprintf(`
type = "workflow"
schemaVersion = 1
name = "test-workflow-with-mock"
externalJobID = "65a75326-65f3-4a40-b8f5-407d2376ce7d"
forwardingAllowed = false
workflowID = "2ab944834c5b3e58aa71e4a4a29d6382b613f833302550277cbc3f5cc2be5226"
workflow = """
name: cciparb003
owner: '0x13e99569ce0ff981ac4e496e362a4bff8eb65265'
triggers:
 - id: streams-trigger@1.1.0
   config:
     maxFrequencyMs: 5000
     feedIds:
       - '0x000351de403f638036014add21a5abd5f464bf21d11aa356dfc6dbe4e2384e4e' # BTC/USD
       - '0x0003f2f4cae1891f647db8d73c87a7a03888bd176afdb7206853da9abfc92874' # ETH/USD
       - '0x00034db6355441c80b613f666757c63777dae7743885a9c594ca25d9f9b896ca' # LINK/USD

consensus:
 - id: offchain_reporting@1.0.0
   ref: ccip_feeds
   inputs:
     observations:
       - $(trigger.outputs)
   config:
     report_id: '0001'
     key_id: 'evm'
     aggregation_method: data_feeds
     aggregation_config:
       allowedPartialStaleness: '0.5'
       feeds:
         '0x000351de403f638036014add21a5abd5f464bf21d11aa356dfc6dbe4e2384e4e':  # BTC/USD
           deviation: '0.01'
           heartbeat: 600
           remappedID: '0x666666666666'
         '0x0003f2f4cae1891f647db8d73c87a7a03888bd176afdb7206853da9abfc92874': # ETH/USD
           deviation: '0.01'
           heartbeat: 600
           remappedID: '0x777777777777'
         '0x00034db6355441c80b613f666757c63777dae7743885a9c594ca25d9f9b896ca': # LINK/USD
           deviation: '0.01'
           heartbeat: 600
     encoder: EVM
     encoder_config:
       abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports

targets:
 - id: write_ethereum-testnet-sepolia-arbitrum-1@1.0.0
   inputs:
     signed_report: $(ccip_feeds.outputs)
   config:
     address: '0x24309990d635A6C5FF711503BfCb942dd25F96A0'
     deltaStage: 10s
     schedule: oneAtATime

"""
workflowOwner = "%s"
`, sc.MustGetRootKeyAddress().Hex())
	errCh := make(chan error, len(nodes))

	var wg sync.WaitGroup

	for _, n := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//This will also accept the job
			//TODO @george-dorin: Change the ProposeJob method so it does not approve workflow jobs
			_, err := ctfEnv.Offchain.ProposeJob(context.Background(), &job2.ProposeJobRequest{
				NodeId: n.NodeID,
				Spec:   workflowJobSpec,
				Labels: nil,
			})

			// Workflows get auto approved
			require.ErrorContains(t, err, "cannot approve an approved spec")
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "job creation/acception failed")
	}

	return

}

// TODO @george-dorin: Refactor!
type capProxy struct {
	Clients []capabilities2.ProxyClient
}

func newCapProxyClient() *capProxy {
	return &capProxy{Clients: make([]capabilities2.ProxyClient, 0)}
}

func (c *capProxy) connectAll(ports []int) error {
	for _, p := range ports {
		client, err := proxyConnectToOne(p)
		if err != nil {
			return err
		}
		c.Clients = append(c.Clients, client)
	}
	return nil
}

func (c *capProxy) CreateCapability(info capabilities2.CapabilityInfo) error {
	for _, client := range c.Clients {
		_, err := client.CreateCapability(context.TODO(), &info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *capProxy) SendTrigger(id string, eventID string, payload []byte) error {
	for _, client := range c.Clients {
		data := capabilities2.SendTriggerEventRequest{
			ID:      id,
			EventID: eventID,
			Payload: payload,
		}
		framework.L.Info().Msg(fmt.Sprintf("Sending trigger response %s:%s", id, eventID))

		_, err := client.SendTriggerEvent(context.TODO(), &data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *capProxy) HookExecutables(lggr zerolog.Logger) error {
	for _, client := range c.Clients {
		hook, errC := client.HookExecutables(context.TODO())
		if errC != nil {
			return errC
		}

		go func() {
			for {
				lggr.Info().Msg("Waiting for hook event")
				resp, err := hook.Recv()
				if err == io.EOF {
					lggr.Error().Msgf("Recieved EOF from hook %s", err)
					return
				}
				if err != nil {
					log.Fatalf("can not receive %v", err)
				}
				lggr.Info().Msgf("Got hook event %v+", resp)

				//Process request
				r := capabilities2.ExecutableResponse{
					ID:             resp.ID,
					CapabilityType: resp.CapabilityType,
					Value:          resp.Inputs,
				}
				err = hook.Send(&r)
				if err != nil {
					panic(err.Error())
				}
				lggr.Info().Msgf("Sent hook response %v+", r)

			}
		}()
	}
	return nil
}

func proxyConnectToOne(port int) (capabilities2.ProxyClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := capabilities2.NewProxyClient(conn) //TODO @george-dorin: Move the proxy pb file
	return client, nil

}

func createFeedReport(t *testing.T, price *big.Int, observationTimestamp int64,
	feedIDString string,
	keyBundles []ocr2key.KeyBundle) *datastreams.FeedReport {
	reportCtx := ocrTypes.ReportContext{}
	rawCtx := RawReportContext(reportCtx)

	bytes, err := hex.DecodeString(feedIDString)
	require.NoError(t, err)
	var feedIDBytes [32]byte
	copy(feedIDBytes[:], bytes)

	report := &datastreams.FeedReport{
		FeedID:               feedIDString,
		FullReport:           newReport(t, feedIDBytes, price, observationTimestamp),
		BenchmarkPrice:       price.Bytes(),
		ObservationTimestamp: observationTimestamp,
		Signatures:           [][]byte{},
		ReportContext:        rawCtx,
	}

	for _, key := range keyBundles {
		sig, err := key.Sign(reportCtx, report.FullReport)
		require.NoError(t, err)
		report.Signatures = append(report.Signatures, sig)
	}

	return report
}
func RawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

func newReport(t *testing.T, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	ctx := tests.Context(t)
	v3Codec := reportcodec.NewReportCodec(feedID, logger.TestLogger(t))
	raw, err := v3Codec.BuildReport(ctx, v3.ReportFields{
		BenchmarkPrice: price,

		Timestamp: uint32(timestamp),
		Bid:       big.NewInt(0),
		Ask:       big.NewInt(0),
		LinkFee:   big.NewInt(0),
		NativeFee: big.NewInt(0),
	})
	require.NoError(t, err)
	return raw
}

func mapToBytes(m *values.Map) ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	pm := make(map[string]*pb2.Value)
	for k, v := range m.Underlying {
		pm[k] = values.Proto(v)
	}
	bytes, err := proto.Marshal(pb2.NewMapValue(pm))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func bytesToMap(b []byte) (*values.Map, error) {
	var o pb2.Value
	if err := proto.Unmarshal(b, &o); err != nil {
		return nil, err
	}

	vm := values.Map{Underlying: make(map[string]values.Value)}
	for k, v := range o.GetMapValue().Fields {
		val, err := values.FromProto(v)
		if err != nil {
			return nil, err
		}
		vm.Underlying[k] = val
	}

	return &vm, nil
}
