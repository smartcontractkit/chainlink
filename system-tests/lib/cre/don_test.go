package cre

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	webclient "github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
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

func mustNewTestNode(t *testing.T) *Node {
	t.Helper()

	p2pKey, err := crypto.NewP2PKey("password")
	require.NoError(t, err)
	evmKey, err := crypto.NewEVMKey("password", 111)
	require.NoError(t, err)

	return &Node{
		Name: "node-1",
		Keys: &secrets.NodeKeys{
			P2PKey: p2pKey,
			EVM: map[uint64]*crypto.EVMKey{
				111: evmKey,
				222: &crypto.EVMKey{PublicAddress: evmKey.PublicAddress},
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
