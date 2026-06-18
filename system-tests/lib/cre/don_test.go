package cre

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	webclient "github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

type fakeGQLClient struct {
	webclient.Client

	fetchOCR2KeyBundleIDFn            func(context.Context, string) (string, error)
	createJobDistributorChainConfigFn func(context.Context, webclient.JobDistributorChainConfigInput) (string, error)

	fetchOCR2KeyBundleIDCalls            []string
	createJobDistributorChainConfigCalls []webclient.JobDistributorChainConfigInput
}

func (f *fakeGQLClient) FetchOCR2KeyBundleID(ctx context.Context, chainType string) (string, error) {
	f.fetchOCR2KeyBundleIDCalls = append(f.fetchOCR2KeyBundleIDCalls, chainType)
	if f.fetchOCR2KeyBundleIDFn != nil {
		return f.fetchOCR2KeyBundleIDFn(ctx, chainType)
	}
	return "bundle-id", nil
}

func (f *fakeGQLClient) CreateJobDistributorChainConfig(ctx context.Context, in webclient.JobDistributorChainConfigInput) (string, error) {
	f.createJobDistributorChainConfigCalls = append(f.createJobDistributorChainConfigCalls, in)
	if f.createJobDistributorChainConfigFn != nil {
		return f.createJobDistributorChainConfigFn(ctx, in)
	}
	return "config-id", nil
}

type fakeJDChainConfigLister struct {
	chainIDsByNodeID map[string]map[string]struct{}
}

func (f *fakeJDChainConfigLister) ListNodeChainConfigs(_ context.Context, in *nodev1.ListNodeChainConfigsRequest, _ ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error) {
	resp := &nodev1.ListNodeChainConfigsResponse{}
	if in.GetFilter() == nil || len(in.GetFilter().GetNodeIds()) == 0 {
		return resp, nil
	}

	for chainID := range f.chainIDsByNodeID[in.GetFilter().GetNodeIds()[0]] {
		resp.ChainConfigs = append(resp.ChainConfigs, &nodev1.ChainConfig{
			Chain: &nodev1.Chain{Id: chainID},
		})
	}

	return resp, nil
}

type fakeBlockchain struct {
	chainID     uint64
	chainFamily string
}

func (f fakeBlockchain) ChainSelector() uint64 { return f.chainID }
func (f fakeBlockchain) ChainID() uint64       { return f.chainID }
func (f fakeBlockchain) ChainFamily() string   { return f.chainFamily }
func (f fakeBlockchain) IsFamily(chainFamily string) bool {
	return f.chainFamily == chainFamily
}
func (f fakeBlockchain) Fund(context.Context, string, uint64) error { return nil }
func (f fakeBlockchain) CtfOutput() *blockchain.Output              { return nil }
func (f fakeBlockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	return nil, nil
}

var _ blockchains.Blockchain = fakeBlockchain{}

func TestCreateJDChainConfigsSkipsExistingConfigs(t *testing.T) {
	t.Parallel()

	node := mustNewTestNode(t)
	gql := &fakeGQLClient{}
	node.Clients.GQLClient = gql
	jd := &fakeJDChainConfigLister{
		chainIDsByNodeID: map[string]map[string]struct{}{
			node.JobDistributorDetails.NodeID: {"111": {}},
		},
	}

	err := createJDChainConfigs(context.Background(), node, []blockchains.Blockchain{
		fakeBlockchain{chainID: 111, chainFamily: blockchain.FamilyEVM},
	}, jd)

	require.NoError(t, err)
	require.Empty(t, gql.createJobDistributorChainConfigCalls)
}

func TestCreateJDChainConfigsCreatesMissingConfigsAndReusesBundleIDs(t *testing.T) {
	t.Parallel()

	node := mustNewTestNode(t)
	jd := &fakeJDChainConfigLister{
		chainIDsByNodeID: map[string]map[string]struct{}{
			node.JobDistributorDetails.NodeID: {},
		},
	}
	gql := &fakeGQLClient{
		createJobDistributorChainConfigFn: func(_ context.Context, in webclient.JobDistributorChainConfigInput) (string, error) {
			jd.chainIDsByNodeID[node.JobDistributorDetails.NodeID][in.ChainID] = struct{}{}
			return "created-" + in.ChainID, nil
		},
	}
	node.Clients.GQLClient = gql

	err := createJDChainConfigs(context.Background(), node, []blockchains.Blockchain{
		fakeBlockchain{chainID: 111, chainFamily: blockchain.FamilyEVM},
		fakeBlockchain{chainID: 222, chainFamily: blockchain.FamilyEVM},
	}, jd)

	require.NoError(t, err)
	require.Len(t, gql.createJobDistributorChainConfigCalls, 2)
	require.Equal(t, []string{"EVM"}, gql.fetchOCR2KeyBundleIDCalls)
	require.Equal(t, "bundle-id", node.Keys.OCR2BundleIDs["evm"])
}

func TestCreateJDChainConfigsFailsVerificationOnTimeout(t *testing.T) {
	node := mustNewTestNode(t)
	jd := &fakeJDChainConfigLister{
		chainIDsByNodeID: map[string]map[string]struct{}{
			node.JobDistributorDetails.NodeID: {},
		},
	}
	node.Clients.GQLClient = &fakeGQLClient{}

	originalTimeout := jdChainConfigPollTimeout
	jdChainConfigPollTimeout = 5 * time.Millisecond
	defer func() {
		jdChainConfigPollTimeout = originalTimeout
	}()

	err := createJDChainConfigs(context.Background(), node, []blockchains.Blockchain{
		fakeBlockchain{chainID: 111, chainFamily: blockchain.FamilyEVM},
	}, jd)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create JD chain configuration")
}

func TestAptosAccountForNode_UsesMetadataKeyWithoutCallingNodeAPI(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		http.NotFound(w, r)
	}))
	t.Cleanup(server.Close)

	expected, err := crecrypto.NormalizeAptosAccount("0x1")
	require.NoError(t, err)

	node := &Node{
		Name: "node-1",
		Keys: &secrets.NodeKeys{
			Aptos: &crecrypto.AptosKey{Account: expected},
		},
		Clients: NodeClients{
			RestClient: &clclient.ChainlinkClient{
				APIClient: resty.New().SetBaseURL(server.URL),
				Config:    &clclient.Config{URL: server.URL},
			},
		},
	}

	account, err := aptosAccountForNode(node)
	require.NoError(t, err)
	require.Equal(t, expected, account)
	require.Zero(t, hits.Load(), "node API must not be called when metadata already has the Aptos key")
}

func TestAptosAccountForNode_FallsBackToNodeAPIAndCachesKey(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/keys/aptos" {
			t.Errorf("unexpected path: got %q want %q", r.URL.Path, "/v2/keys/aptos")
			http.Error(w, fmt.Sprintf("unexpected path %q", r.URL.Path), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"data":[{"attributes":{"account":"0x1","publicKey":"0xabc123"}}]}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	node := &Node{
		Name: "node-1",
		Keys: &secrets.NodeKeys{},
		Clients: NodeClients{
			RestClient: &clclient.ChainlinkClient{
				APIClient: resty.New().SetBaseURL(server.URL),
				Config:    &clclient.Config{URL: server.URL},
			},
		},
	}

	account, err := aptosAccountForNode(node)
	require.NoError(t, err)

	expected, err := crecrypto.NormalizeAptosAccount("0x1")
	require.NoError(t, err)
	require.Equal(t, expected, account)
	require.NotNil(t, node.Keys.Aptos)
	require.Equal(t, expected, node.Keys.Aptos.Account)
}

func mustNewTestNode(t *testing.T) *Node {
	t.Helper()

	p2pKey, err := crecrypto.NewP2PKey("password")
	require.NoError(t, err)
	evmKey, err := crecrypto.NewEVMKey("password", 111)
	require.NoError(t, err)

	return &Node{
		Name: "node-1",
		Keys: &secrets.NodeKeys{
			P2PKey: p2pKey,
			EVM: map[uint64]*crecrypto.EVMKey{
				111: evmKey,
				222: {PublicAddress: evmKey.PublicAddress},
			},
		},
		Addresses: Addresses{
			AdminAddress: "0xadmin",
		},
		JobDistributorDetails: &JobDistributorDetails{
			NodeID: "node-id-1",
			JDID:   "jd-id-1",
		},
		Clients: NodeClients{},
		Roles:   Roles{RoleWorker},
	}
}
