package discovery

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

// fakeK8s implements K8sAPI for tests.
type fakeK8s struct {
	apiInfo map[string]*infra.NodeAPIInfo
	errors  map[string]error
}

func (f *fakeK8s) GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*infra.NodeAPIInfo, error) {
	if err, ok := f.errors[nodeName]; ok {
		return nil, err
	}
	return f.apiInfo[nodeName], nil
}

// fakeClient implements Client for tests.
type fakeClient struct {
	csa       string
	peerID    string
	evmAddrs  map[string]string
	ocr2IDs   map[string]string
	aptosAddr string
	solAddr   string
	csaErr    error
	peerErr   error
	evmErr    error
	ocr2Err   error
	aptosErr  error
	solErr    error
}

func (f *fakeClient) ReadCSAKey() (string, error) {
	return f.csa, f.csaErr
}

func (f *fakeClient) ReadPeerID() (string, error) {
	return f.peerID, f.peerErr
}

func (f *fakeClient) ReadEVMAddresses() (map[string]string, error) {
	return f.evmAddrs, f.evmErr
}

func (f *fakeClient) ReadOCR2BundleIDs() (map[string]string, error) {
	return f.ocr2IDs, f.ocr2Err
}

func (f *fakeClient) ReadAptosKeys() (string, error) {
	return f.aptosAddr, f.aptosErr
}

func (f *fakeClient) ReadSolanaKeys() (string, error) {
	return f.solAddr, f.solErr
}

// fakeDialer implements Dialer for tests.
type fakeDialer struct {
	clients map[string]*fakeClient
	errors  map[string]error
}

func (f *fakeDialer) Dial(apiURL, email, password string) (Client, error) {
	// Simplistic: match by apiURL
	for url, client := range f.clients {
		if url == apiURL {
			return client, nil
		}
	}
	// Return error if URL not found
	if err, ok := f.errors[apiURL]; ok {
		return nil, err
	}
	return nil, nil
}

func TestDiscovery_SuccessfulDiscovery(t *testing.T) {
	t.Parallel()

	k8s := &fakeK8s{
		apiInfo: map[string]*infra.NodeAPIInfo{
			"node-0": {URL: "http://node-0:6688", Email: "test", Password: "test"},
		},
	}
	dialer := &fakeDialer{
		clients: map[string]*fakeClient{
			"http://node-0:6688": {
				csa:      "csa_key_abc",
				peerID:   "12D3KooWABC",
				evmAddrs: map[string]string{"1": "0xabc"},
				ocr2IDs:  map[string]string{"evm": "bundle_id_123"},
			},
		},
	}

	nodes := []domain.ChartNodeInfo{{Name: "node-0", NodeType: domain.RoleStandard}}
	cv := &domain.ChartValues{Nodes: nodes, Namespace: "default"}

	result, err := Run(context.Background(), zerolog.Nop(), nodes, cv, k8s, dialer, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	info, ok := result["node-0"]
	require.True(t, ok)
	require.Equal(t, "csa_key_abc", info.CSAKey)
	require.Equal(t, "12D3KooWABC", info.PeerID)
	require.Equal(t, "0xabc", info.EVMAddress["1"])
	require.Equal(t, "bundle_id_123", info.OCR2BundleIDs["evm"])
}

func TestDiscovery_SkipsOCR2ReadForBootstrapAndGatewayNodes(t *testing.T) {
	t.Parallel()

	k8s := &fakeK8s{
		apiInfo: map[string]*infra.NodeAPIInfo{
			"node-boot": {URL: "http://node-boot:6688", Email: "test", Password: "test"},
			"node-gw":   {URL: "http://node-gw:6688", Email: "test", Password: "test"},
		},
	}
	dialer := &fakeDialer{
		clients: map[string]*fakeClient{
			// ocr2IDs/ocr2Err set to prove the field would be populated/erroring
			// if ReadOCR2BundleIDs were called; it must not be for these roles.
			"http://node-boot:6688": {csa: "csa_boot", ocr2IDs: map[string]string{"evm": "should-not-be-read"}},
			"http://node-gw:6688":   {csa: "csa_gw", ocr2Err: context.Canceled},
		},
	}

	nodes := []domain.ChartNodeInfo{
		{Name: "node-boot", NodeType: domain.RoleBootstrap},
		{Name: "node-gw", NodeType: domain.RoleGateway},
	}
	cv := &domain.ChartValues{Nodes: nodes, Namespace: "default"}

	result, err := Run(context.Background(), zerolog.Nop(), nodes, cv, k8s, dialer, nil)

	require.NoError(t, err)
	require.Len(t, result, 2)

	bootInfo, ok := result["node-boot"]
	require.True(t, ok)
	require.Equal(t, "csa_boot", bootInfo.CSAKey)
	require.Empty(t, bootInfo.OCR2BundleIDs)

	gwInfo, ok := result["node-gw"]
	require.True(t, ok)
	require.Equal(t, "csa_gw", gwInfo.CSAKey)
	require.Empty(t, gwInfo.OCR2BundleIDs)
}

func TestDiscovery_ReadsNonEVMKeysOnlyWhenNeeded(t *testing.T) {
	t.Parallel()

	k8s := &fakeK8s{
		apiInfo: map[string]*infra.NodeAPIInfo{
			"node-solana": {URL: "http://node-solana:6688", Email: "test", Password: "test"},
			"node-plain":  {URL: "http://node-plain:6688", Email: "test", Password: "test"},
		},
	}
	dialer := &fakeDialer{
		clients: map[string]*fakeClient{
			"http://node-solana:6688": {csa: "csa-solana", solAddr: "sol-addr-1", aptosErr: context.Canceled},
			// aptosAddr/aptosErr set to prove ReadAptosKeys/ReadSolanaKeys are never
			// called for a node whose DON doesn't need that family.
			"http://node-plain:6688": {csa: "csa-plain", solAddr: "should-not-be-read", aptosAddr: "should-not-be-read"},
		},
	}

	nodes := []domain.ChartNodeInfo{
		{Name: "node-solana", NodeType: domain.RoleStandard},
		{Name: "node-plain", NodeType: domain.RoleStandard},
	}
	cv := &domain.ChartValues{Nodes: nodes, Namespace: "default"}
	nonEVMFamilies := map[string][]string{"node-solana": {"solana"}}

	result, err := Run(context.Background(), zerolog.Nop(), nodes, cv, k8s, dialer, nonEVMFamilies)

	require.NoError(t, err)
	require.Equal(t, "sol-addr-1", result["node-solana"].SolanaAddress)
	require.Empty(t, result["node-solana"].AptosAddress)
	require.Empty(t, result["node-plain"].SolanaAddress)
	require.Empty(t, result["node-plain"].AptosAddress)
}

func TestDiscovery_SkipsFailedNodes(t *testing.T) {
	t.Parallel()

	k8s := &fakeK8s{
		apiInfo: map[string]*infra.NodeAPIInfo{
			"node-0": {URL: "http://node-0:6688", Email: "test", Password: "test"},
			"node-1": {URL: "http://node-1:6688", Email: "test", Password: "test"},
		},
	}
	dialer := &fakeDialer{
		clients: map[string]*fakeClient{
			"http://node-0:6688": {
				csa:      "csa_key_abc",
				peerID:   "12D3KooWABC",
				evmAddrs: map[string]string{"1": "0xabc"},
			},
		},
		errors: map[string]error{
			"http://node-1:6688": context.Canceled,
		},
	}

	nodes := []domain.ChartNodeInfo{
		{Name: "node-0", NodeType: domain.RoleStandard},
		{Name: "node-1", NodeType: domain.RoleStandard},
	}
	cv := &domain.ChartValues{Nodes: nodes, Namespace: "default"}

	result, err := Run(context.Background(), zerolog.Nop(), nodes, cv, k8s, dialer, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	_, ok := result["node-0"]
	require.True(t, ok)
	_, ok = result["node-1"]
	require.False(t, ok)
}
