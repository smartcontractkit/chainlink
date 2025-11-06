package operations

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type mockOffchainClient struct {
	sync.Mutex

	jobv1.JobServiceClient
	nodev1.NodeServiceClient
	csav1.CSAServiceClient

	nodes        []*nodev1.Node
	chainConfigs []*nodev1.ChainConfig

	listNodesErr            error
	listNodeChainConfigsErr error
}

func newOffchainClient(nodes []*nodev1.Node, chainConfigs []*nodev1.ChainConfig, listNodesErr error, listNodeChainConfigsErr error) *mockOffchainClient {
	return &mockOffchainClient{
		nodes:        nodes,
		chainConfigs: chainConfigs,

		listNodesErr:            listNodesErr,
		listNodeChainConfigsErr: listNodeChainConfigsErr,
	}
}

func (oc *mockOffchainClient) ListNodes(ctx context.Context, in *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	if oc.listNodesErr != nil {
		return &nodev1.ListNodesResponse{}, oc.listNodesErr
	}

	foundNodes := []*nodev1.Node{}
	for _, node := range oc.nodes {
		match := true
		labelLookup := make(map[string]*ptypes.Label)
		for _, label := range node.Labels {
			labelLookup[label.Key] = label
		}
		for _, selector := range in.Filter.Selectors {
			label, ok := labelLookup[selector.Key]
			if !ok || (selector.Value != nil && *selector.Value != *label.Value) {
				match = false
				break
			}
		}
		if match {
			foundNodes = append(foundNodes, node)
		}
	}

	return &nodev1.ListNodesResponse{Nodes: foundNodes}, nil
}

func (oc *mockOffchainClient) ListNodeChainConfigs(ctx context.Context, in *nodev1.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error) {
	if oc.listNodeChainConfigsErr != nil {
		return &nodev1.ListNodeChainConfigsResponse{}, oc.listNodeChainConfigsErr
	}

	foundChainConfigs := []*nodev1.ChainConfig{}
	for _, chainConfig := range oc.chainConfigs {
		for _, nodeID := range in.Filter.NodeIds {
			if chainConfig.NodeId != nodeID {
				continue
			}
			foundChainConfigs = append(foundChainConfigs, chainConfig)
		}
	}

	return &nodev1.ListNodeChainConfigsResponse{ChainConfigs: foundChainConfigs}, nil
}

func (oc *mockOffchainClient) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	return &jobv1.ProposeJobResponse{}, nil
}

var commonInput = ProposeGatewayJobInput{
	Domain: "cre",
	DONFilters: []offchain.TargetDONFilter{
		{
			Key:   "don_name",
			Value: "gateway_1_zone-b",
		},
		{
			Key:   "environment",
			Value: "staging",
		},
		{
			Key:   "product",
			Value: "cre",
		},
		{
			Key:   "zone",
			Value: "zone-b",
		},
	},
	DONs: []DON{
		{
			Name: "workflow_1_zone-b",
			Handlers: []string{
				"http-capabilities",
				"web-api-capabilities",
			},
			F: 1,
		},
	},
	GatewayKeyChainSelector:  10344971235874465080,
	GatewayRequestTimeoutSec: 5,
	JobLabels:                map[string]string{},
}

func TestProposeGatewayJob(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		input                 ProposeGatewayJobInput
		offchainClientFactory func() *mockOffchainClient
		expectError           bool
		errorMsg              string
		output                ProposeGatewayJobOutput
	}{
		{
			name:  "success",
			input: commonInput,
			offchainClientFactory: func() *mockOffchainClient {
				return newOffchainClient(nodes, chainConfigs, nil, nil)
			},
			expectError: false,
			errorMsg:    "",
			output: ProposeGatewayJobOutput{
				Specs: map[string][]string{
					"node_5": {
						"type = 'gateway'\nschemaVersion = 1\nname = 'CRE Gateway'\nexternalJobID = 'cf8aa339-6349-5e5b-9289-5c2907711200'\nforwardingAllowed = false\n\n[gatewayConfig]\n[gatewayConfig.ConnectionManagerConfig]\nAuthChallengeLen = 10\nAuthGatewayId = 'gateway-node-0'\nAuthTimestampToleranceSec = 5\nHeartbeatIntervalSec = 20\n\n[[gatewayConfig.Dons]]\nDonId = 'workflow_1_zone-b'\nF = 1\n\n[[gatewayConfig.Dons.Handlers]]\nName = 'http-capabilities'\nServiceName = 'workflows'\n\n[gatewayConfig.Dons.Handlers.Config]\nCleanUpPeriodMs = 600000\n\n[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]\nglobalBurst = 100\nglobalRPS = 500\nperSenderBurst = 100\nperSenderRPS = 100\n\n[[gatewayConfig.Dons.Handlers]]\nName = 'web-api-capabilities'\n\n[gatewayConfig.Dons.Handlers.Config]\nmaxAllowedMessageAgeSec = 1000\n\n[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]\nglobalBurst = 10\nglobalRPS = 50\nperSenderBurst = 10\nperSenderRPS = 10\n\n[[gatewayConfig.Dons.Members]]\nAddress = '0x04'\nName = 'cl-cre-one-zone-b-0 (DON workflow_1_zone-b)'\n\n[gatewayConfig.HTTPClientConfig]\nMaxResponseBytes = 50000000\nAllowedPorts = [443]\nAllowedSchemes = ['https']\nAllowedIPsCIDR = []\n\n[gatewayConfig.NodeServerConfig]\nHandshakeTimeoutMillis = 1000\nMaxRequestBytes = 100000\nPath = '/'\nPort = 5003\nReadTimeoutMillis = 1000\nRequestTimeoutMillis = 5000\nWriteTimeoutMillis = 1000\n\n[gatewayConfig.UserServerConfig]\nContentTypeHeader = 'application/jsonrpc'\nMaxRequestBytes = 100000\nPath = '/'\nPort = 5002\nReadTimeoutMillis = 5000\nRequestTimeoutMillis = 5000\nWriteTimeoutMillis = 6000\n"},
				},
			},
		},
		{
			name:  "fail - no nodes found",
			input: commonInput,
			offchainClientFactory: func() *mockOffchainClient {
				return newOffchainClient([]*nodev1.Node{}, chainConfigs, nil, nil)
			},
			expectError: true,
			errorMsg:    "no nodes with filters",
			output:      ProposeGatewayJobOutput{},
		},
		{
			name:  "fail - no chain configs found",
			input: commonInput,
			offchainClientFactory: func() *mockOffchainClient {
				return newOffchainClient(nodes, []*nodev1.ChainConfig{}, nil, nil)
			},
			expectError: true,
			errorMsg:    "no chain configs with filters",
			output:      ProposeGatewayJobOutput{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			selector := uint64(909606746561742123) // chain_selectors.TEST_90000001.Selector

			env, err := environment.New(t.Context(),
				environment.WithEVMSimulated(t, []uint64{selector}),
				environment.WithOffchainClient(tc.offchainClientFactory()),
			)
			require.NoError(t, err)

			output, err := proposeGatewayJob(operations.Bundle{GetContext: func() context.Context { return t.Context() }}, ProposeGatewayJobDeps{Env: *env}, tc.input)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, output)
			}

			require.Equal(t, tc.output, output)
		})
	}
}

func createString(x string) *string {
	return &x
}

var nodes = []*nodev1.Node{
	// Product: CRE
	{
		Id:          "node_1",
		Name:        "cl-cre-one-zone-b-bt-0",
		PublicKey:   "pk1",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don",
				Value: createString("bootstrappers_virtual_don"),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-1"),
			},
			{
				Key:   "product",
				Value: createString("cre"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk1"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-1",
				PublicKey: "p2ppk1",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_2",
		Name:        "cl-cre-one-zone-b-0",
		PublicKey:   "pk2",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don-workflow_1_zone-b",
				Value: createString(""),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-2"),
			},
			{
				Key:   "product",
				Value: createString("cre"),
			},
			{
				Key:   "type",
				Value: createString("plugin"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk2"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-2",
				PublicKey: "p2ppk2",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_3",
		Name:        "cl-cre-chain-capabilities-zone-b-bt-0",
		PublicKey:   "pk3",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don",
				Value: createString("bootstrappers_virtual_don"),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-3"),
			},
			{
				Key:   "product",
				Value: createString("cre"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk3"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-3",
				PublicKey: "p2ppk3",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_4",
		Name:        "cl-cre-chain-capabilities-zone-b-0",
		PublicKey:   "pk4",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don-chain_capabilities_zone-b",
				Value: createString(""),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-4"),
			},
			{
				Key:   "product",
				Value: createString("cre"),
			},
			{
				Key:   "type",
				Value: createString("plugin"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk4"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-4",
				PublicKey: "p2ppk4",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_5",
		Name:        "cl-cre-gateway-one-zone-b-0",
		PublicKey:   "pk5",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don-gateway_1_zone-b",
				Value: createString(""),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-5"),
			},
			{
				Key:   "product",
				Value: createString("cre"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk5"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-5",
				PublicKey: "p2ppk5",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	// Product: Other to add some noise
	{
		Id:          "node_6",
		Name:        "cl-cre-one-zone-b-bt-0",
		PublicKey:   "pk6",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don",
				Value: createString("bootstrappers_virtual_don"),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-6"),
			},
			{
				Key:   "product",
				Value: createString("other"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk6"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-6",
				PublicKey: "p2ppk6",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_7",
		Name:        "cl-cre-one-zone-b-0",
		PublicKey:   "pk7",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don-workflow_1_zone-b",
				Value: createString(""),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-7"),
			},
			{
				Key:   "product",
				Value: createString("other"),
			},
			{
				Key:   "type",
				Value: createString("plugin"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk7"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-7",
				PublicKey: "p2ppk7",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_8",
		Name:        "cl-cre-chain-capabilities-zone-b-bt-0",
		PublicKey:   "pk8",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don",
				Value: createString("bootstrappers_virtual_don"),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-8"),
			},
			{
				Key:   "product",
				Value: createString("other"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk8"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-8",
				PublicKey: "p2ppk8",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_9",
		Name:        "cl-cre-chain-capabilities-zone-b-0",
		PublicKey:   "pk9",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don-chain_capabilities_zone-b",
				Value: createString(""),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-9"),
			},
			{
				Key:   "product",
				Value: createString("other"),
			},
			{
				Key:   "type",
				Value: createString("plugin"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk9"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-9",
				PublicKey: "p2ppk9",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
	{
		Id:          "node_10",
		Name:        "cl-cre-gateway-one-zone-b-0",
		PublicKey:   "pk10",
		IsEnabled:   true,
		IsConnected: true,
		Labels: []*ptypes.Label{
			{
				Key:   "don-gateway_1_zone-b",
				Value: createString(""),
			},
			{
				Key:   "environment",
				Value: createString("staging"),
			},
			{
				Key:   "p2p_id",
				Value: createString("12D3Koo-10"),
			},
			{
				Key:   "product",
				Value: createString("other"),
			},
			{
				Key:   "zone",
				Value: createString("zone-b"),
			},
		},
		WorkflowKey: createString("wfk10"),
		P2PKeyBundles: []*nodev1.P2PKeyBundle{
			{
				PeerId:    "p2p_12D3Koo-10",
				PublicKey: "p2ppk10",
			},
		},
		Version:   "v1",
		CreatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
		UpdatedAt: &timestamppb.Timestamp{Seconds: 1767318549, Nanos: 191229000},
	},
}

var chainConfigs = []*nodev1.ChainConfig{
	{
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x01",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: true,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-1",
				PublicKey: "p2ppk1",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-1@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_1",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x02",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: true,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-1",
				PublicKey: "p2ppk1",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-1@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_1",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x03",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-2",
				PublicKey: "p2ppk2",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-1@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_2",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x04",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-2",
				PublicKey: "p2ppk2",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-1@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_2",
	}, {
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x05",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: true,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-3",
				PublicKey: "p2ppk3",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-3@cl-cre-chain-capabilities-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_3",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x06",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: true,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-3",
				PublicKey: "p2ppk3",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-3@cl-cre-chain-capabilities-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_3",
	}, {
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x07",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-4",
				PublicKey: "p2ppk4",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-3@cl-cre-chain-capabilities-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_4",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x08",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-4",
				PublicKey: "p2ppk4",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-3@cl-cre-chain-capabilities-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_4",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x10",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config:              &nodev1.OCR2Config{},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_5",
	}, {
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x11",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: true,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-6",
				PublicKey: "p2ppk6",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-6@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_6",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x12",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: true,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-6",
				PublicKey: "p2ppk6",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-6@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_6",
	}, {
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x13",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-7",
				PublicKey: "p2ppk7",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-6@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_7",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x14",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-7",
				PublicKey: "p2ppk7",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-6@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_7",
	}, {
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x15",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-8",
				PublicKey: "p2ppk8",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-8@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_8",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x16",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-8",
				PublicKey: "p2ppk8",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-8@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_8",
	}, {
		Chain: &nodev1.Chain{
			Id:   "11155111",
			Type: 1,
		},
		AccountAddress: "0x17",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-9",
				PublicKey: "p2ppk9",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-8@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_9",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x18",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config: &nodev1.OCR2Config{
			Enabled:     true,
			IsBootstrap: false,
			P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
				PeerId:    "p2p_12D3Koo-9",
				PublicKey: "p2ppk9",
			},
			OcrKeyBundle:     &nodev1.OCR2Config_OCRKeyBundle{},
			Multiaddr:        "12D3Koo-8@cl-cre-one-zone-b-bt-0:5001",
			Plugins:          &nodev1.OCR2Config_Plugins{},
			ForwarderAddress: createString(""),
		},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_9",
	},
	{
		Chain: &nodev1.Chain{
			Id:   "84532",
			Type: 1,
		},
		AccountAddress: "0x20",
		Ocr1Config: &nodev1.OCR1Config{
			P2PKeyBundle: &nodev1.OCR1Config_P2PKeyBundle{},
			OcrKeyBundle: &nodev1.OCR1Config_OCRKeyBundle{},
		},
		Ocr2Config:              &nodev1.OCR2Config{},
		AccountAddressPublicKey: createString(""),
		NodeId:                  "node_10",
	},
}
