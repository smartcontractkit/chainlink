package capabilities_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/capabilities"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
)

// Copy-paste of the OCR3 Workflow test, used for development of the proxy capability
func TestKeystoneWithOCR3WorkflowWithProxy(t *testing.T) {
	testLogger := framework.L

	// we need to use double-pointers, so that what's captured in the cleanup function is a pointer, not the actual object,
	// which is only set later in the test, after the cleanup function is defined
	var nodes **ns.Output
	var wsRPCURL *string

	// clean up is LIFO, so we need to make sure we execute the debug report transmission after logs are written down
	// by function added to clean up by framework.Load() method
	t.Cleanup(func() {
		if t.Failed() {
			if nodes == nil {
				testLogger.Warn().Msg("nodeset output is nil, skipping debug report transmission")
				return
			}
			printTestDebug(t, testLogger, *nodes, *wsRPCURL)
		}
	})

	// Load test configuration
	in, err := framework.Load[WorkflowTestConfig](t)
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

	// Start job distributor
	jdOutput := startJobDistributor(t, in)

	// Deploy the DON
	nodeOutput := startNodes(t, in, bc)

	// Prepare the chainlink/deployment environment
	ctfEnv, don, chainSelector := buildChainlinkDeploymentEnv(t, jdOutput, nodeOutput, bc, sc)

	// Fund the nodes¡™¡
	fundNodes(t, don, sc)

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability)
	keystoneContractSet := deployKeystoneContracts(t, testLogger, ctfEnv, chainSelector)

	// Deploy and pre-configure workflow registry contract
	workflowRegistryAddr := prepareWorkflowRegistry(t, testLogger, ctfEnv, chainSelector, sc, in.WorkflowConfig.DonID)

	// Deploy and configure Keystone Feeds Consumer contract
	feedsConsumerAddress := prepareFeedsConsumer(t, testLogger, ctfEnv, chainSelector, sc, keystoneContractSet.Forwarder.Address(), in.WorkflowConfig.WorkflowName)

	// Register the workflow (either via CRE CLI or by calling the workflow registry directly)
	registerWorkflow(t, in, sc, keystoneContractSet.CapabilitiesRegistry.Address(), workflowRegistryAddr, feedsConsumerAddress, in.WorkflowConfig.DonID, chainSelector, in.WorkflowConfig.WorkflowName, pkey, bc.Nodes[0].HostHTTPUrl)

	// Create OCR3 and capability jobs for each node JD
	ns, nodeClient := configureNodes(t, don, in, bc, keystoneContractSet.CapabilitiesRegistry.Address(), workflowRegistryAddr, keystoneContractSet.Forwarder.Address())
	// JD client needs to be reinitialised after restarting nodes
	ctfEnv = ptr.Ptr(reinitialiseJDClient(t, ctfEnv, jdOutput, nodeOutput))
	createNodeJobsWithJd(t, ctfEnv, don, bc, keystoneContractSet) //TODO @george-dorin: Split this so we don't create the cron cap job

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, feedsConsumerAddress.Hex(), keystoneContractSet.Forwarder.Address().Hex())
		}
	})

	//Start proxy cap job
	//TODO @george-dorin: Switch to jd and don't create on bootstrap node
	for i := range nodeClient {
		_, _, err2 := nodeClient[i].CreateJobRaw(`
					type = "standardcapabilities"
					schemaVersion = 1
					externalJobID = "44fc1902-82ef-4cb6-b94f-03d7dc1617f2"
					name = "proxy-capability"
					forwardingAllowed = false
					command = "/home/capabilities/amd64_capproxy"
					config = ""
				`)
		require.NoError(t, err2)
	}

	//TODO @george-dorin: FixME!
	time.Sleep(time.Second * 10)
	//We have hardcode the grpc server ports as 13300-13304
	//1. Connect to each grpc server
	proxyClients := newCapProxyClient()
	require.NoError(t, proxyClients.connectAll([]int{13301, 13302, 13303, 13304}))
	//2. List all capabilities
	for _, c := range proxyClients.Clients {
		r, err2 := c.List(context.TODO(), &capabilities.ListRequest{})
		require.NoError(t, err2)
		framework.L.Info().Msg(fmt.Sprintf("Cap-Proxy got List() response %+v", r.CapInfos))
	}
	//3. Register as cron-trigger@1.0.0
	require.NoError(t, proxyClients.CreateTriggerCap("cron-trigger@1.0.0"))

	//4. Send data to all fake crons
	payload := cron.Payload{
		ScheduledExecutionTime: time.Now().UTC().Format(time.RFC3339Nano),
		ActualExecutionTime:    time.Now().UTC().Format(time.RFC3339Nano),
	}
	wrappedPayload, err := values.WrapMap(payload)

	bytes, err := proto.Marshal(values.Proto(wrappedPayload))
	require.NoError(t, err)
	require.NoError(t, proxyClients.SendTrigger("cron-trigger@1.0.0", uuid.New().String(), bytes))

	// set variables that are needed for the cleanup function, which debugs report transmissions
	nodes = &ns
	wsRPCURL = &bc.Nodes[0].HostWSUrl

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO make it fluent!
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the workflow DON and contracts
	configureWorkflowDON(t, ctfEnv, don, chainSelector)

	// It can take a while before the first report is produced, particularly on CI.
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	testLogger.Info().Msg("Waiting for feed to update...")
	startTime := time.Now()
	feedBytes := common.HexToHash(in.WorkflowConfig.FeedID)

	for {
		select {
		case <-ctx.Done():
			testLogger.Error().Msgf("feed did not update, timeout after %s", timeout)
			t.FailNow()
		case <-time.After(10 * time.Second):
			//Send trigger message
			payload = cron.Payload{
				ScheduledExecutionTime: time.Now().UTC().Format(time.RFC3339Nano),
				ActualExecutionTime:    time.Now().UTC().Format(time.RFC3339Nano),
			}
			wrappedPayload, err = values.WrapMap(payload)

			bytes, err = proto.Marshal(values.Proto(wrappedPayload))
			require.NoError(t, err)
			require.NoError(t, proxyClients.SendTrigger("cron-trigger@1.0.0", uuid.New().String(), bytes))

			elapsed := time.Since(startTime).Round(time.Second)
			price, _, err := feedsConsumerInstance.GetPrice(
				sc.NewCallOpts(),
				feedBytes,
			)
			require.NoError(t, err, "failed to get price from Keystone Consumer contract")

			if price.String() != "0" {
				testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
				return
			}
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}
	}
}

// TODO @george-dorin: Refactor!
type capProxy struct {
	Clients   []capabilities.ProxyClient
	pID2CapID map[string][]string
}

func newCapProxyClient() *capProxy {
	return &capProxy{Clients: make([]capabilities.ProxyClient, 0), pID2CapID: make(map[string][]string, 0)}
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

func (c *capProxy) CreateTriggerCap(id string) error {
	for _, client := range c.Clients {
		r, err := client.CreateTrigger(context.TODO(), &capabilities.CapabilityInfo{
			ProxyID:        "",
			ID:             id,
			CapabilityType: capabilities.CapabilityType(1),
			Description:    "",
			DON:            nil,
			IsLocal:        true,
		})
		if err != nil {
			return err
		}
		c.pID2CapID[id] = append(c.pID2CapID[id], r.ProxyID)
	}
	return nil
}

func (c *capProxy) SendTrigger(id string, eventID string, payload []byte) error {
	for i, client := range c.Clients {
		data := capabilities.SendTriggerRequest{
			ProxyID: c.pID2CapID[id][i],
			EventID: eventID,
			Payload: payload,
		}
		framework.L.Info().Msg(fmt.Sprintf("Sending trigger response to %s: %v+", c.pID2CapID[id][i], data))
		_, err := client.SendTrigger(context.TODO(), &data)
		if err != nil {
			return err
		}
	}
	return nil
}

func proxyConnectToOne(port int) (capabilities.ProxyClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := capabilities.NewProxyClient(conn) //TODO @george-dorin: Move the proxy pb file
	return client, nil
}
