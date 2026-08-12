package cre

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-confidential-compute/tests/testhelpers"
	cctypes "github.com/smartcontractkit/chainlink-confidential-compute/types"
	"github.com/smartcontractkit/chainlink-confidential-compute/util"
	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/confidentialrelay"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const (
	// confidentialWorkflowsApp is the enclave application, the DON capability flag
	// and the capability binary name; all three share this value.
	confidentialWorkflowsApp = "confidential-workflows"

	// confidentialWorkflowsCapVersion is the version the capability registers under.
	confidentialWorkflowsCapVersion = "1.0.0-alpha"

	// confidentialGatewayProxyPort is the fixed port the enclaves are told to reach
	// the CRE gateway on. It must be known before the enclaves start, which is why
	// the proxy in front of it resolves its target lazily.
	confidentialGatewayProxyPort = 9999

	// confidentialStorageKeyHex is a deterministic ed25519 seed the enclave uses to
	// authenticate to the fake storage service. The fake does not verify the JWT.
	confidentialStorageKeyHex = "0000000000000000000000000000000000000000000000000000000000000001"

	confidentialWorkflowBinaryFilename = "workflow-test-confidential.br.b64"
	confidentialWorkflowConfigFilename = "workflow-test-config.json"

	// confidentialWorkflowSrcRelDir is the WASM workflow this test compiles,
	// relative to the chainlink-confidential-compute checkout. The source is not
	// vendored here on purpose: it depends on cre-sdk-go versions that predate the
	// removal of the in-TEE HTTP API, which this repository's dependency
	// validation rejects, and the compiled artifact is covered by .gitignore.
	confidentialWorkflowSrcRelDir = "tests/e2e/testdata/workflow"

	// confidentialEnclaveRegion is recorded on each enclave descriptor. Descriptor
	// hashes cover it, so it must match what the enclave itself reports.
	confidentialEnclaveRegion = "us-west-2"
)

// ---------------------------------------------------------------------------
// Deferred gateway proxy
// ---------------------------------------------------------------------------

// deferredGatewayProxy is a reverse proxy on a fixed port that returns 502 until
// SetTarget is called with the real gateway URL. This resolves a chicken-and-egg
// problem: the enclaves are told their gateway URL at startup, but the real URL
// is only known once the CRE environment is up.
type deferredGatewayProxy struct {
	mu     sync.RWMutex
	target *url.URL
	server *http.Server
}

func newDeferredGatewayProxy(t *testing.T, port int) *deferredGatewayProxy {
	t.Helper()

	p := &deferredGatewayProxy{}
	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			p.mu.RLock()
			defer p.mu.RUnlock()
			if p.target != nil {
				req.URL.Scheme = p.target.Scheme
				req.URL.Host = p.target.Host
				req.Host = p.target.Host
			}
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.mu.RLock()
		hasTarget := p.target != nil
		p.mu.RUnlock()
		if !hasTarget {
			http.Error(w, "gateway not ready", http.StatusBadGateway)
			return
		}
		rp.ServeHTTP(w, r)
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	require.NoError(t, err, "failed to listen on port %d for gateway proxy", port)

	p.server = &http.Server{Handler: handler}
	go func() { _ = p.server.Serve(listener) }()
	t.Cleanup(func() { _ = p.server.Close() })

	return p
}

func (p *deferredGatewayProxy) SetTarget(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.target = u
	return nil
}

// ---------------------------------------------------------------------------
// Fake CRE storage service
// ---------------------------------------------------------------------------

// fakeStorageService is a minimal in-process CRE storage NodeService. The enclave
// fetches the workflow binary itself: it calls DownloadArtifact over JWT-authed
// gRPC, gets a pre-signed URL, downloads it and verifies the hash. This fake
// returns the URL of the base64 WASM server the test stands up.
type fakeStorageService struct {
	storage_service.UnimplementedNodeServiceServer
	mu  sync.Mutex
	url string
}

func (f *fakeStorageService) setURL(u string) {
	f.mu.Lock()
	f.url = u
	f.mu.Unlock()
}

func (f *fakeStorageService) DownloadArtifact(_ context.Context, req *storage_service.DownloadArtifactRequest) (*storage_service.DownloadArtifactResponse, error) {
	f.mu.Lock()
	u := f.url
	f.mu.Unlock()

	// Mirror real storage-service semantics: the id must be a bare artifact id,
	// not a full URL. Rejecting the URL shape here keeps a regression failing in
	// this test rather than only in a live environment.
	if strings.Contains(req.GetId(), "://") {
		return nil, status.Errorf(codes.NotFound, "fake storage: artifact with id %q not found (expected a bare id, not a URL)", req.GetId())
	}
	if u == "" {
		return nil, fmt.Errorf("fake storage: artifact url not set yet")
	}
	return &storage_service.DownloadArtifactResponse{Url: u}, nil
}

// startFakeStorageService starts a gRPC NodeService bound to 0.0.0.0 and returns
// the address the enclave dials it at, plus the service so the test can set the
// artifact URL once the WASM server is up.
func startFakeStorageService(t *testing.T, enclaveHost string) (string, *fakeStorageService) {
	t.Helper()

	lis, err := net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err, "fake storage listener")

	svc := &fakeStorageService{}
	grpcSrv := grpc.NewServer()
	storage_service.RegisterNodeServiceServer(grpcSrv, svc)
	go func() { _ = grpcSrv.Serve(lis) }()
	t.Cleanup(grpcSrv.Stop)

	return fmt.Sprintf("%s:%d", enclaveHost, lis.Addr().(*net.TCPAddr).Port), svc
}

// ---------------------------------------------------------------------------
// Confidential relay feature wrapper
// ---------------------------------------------------------------------------

// testConfidentialRelayFeature wraps the real ConfidentialRelay feature and
// injects trusted measurements into the DON's capability config before
// PreEnvStartup runs, so the relay handler will accept attestations from the
// enclaves this test started.
type testConfidentialRelayFeature struct {
	inner    confidentialrelay.ConfidentialRelay
	pcrsJSON string
}

func (f *testConfidentialRelayFeature) Flag() crelib.CapabilityFlag {
	return f.inner.Flag()
}

func (f *testConfidentialRelayFeature) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *crelib.DonMetadata,
	topology *crelib.Topology,
	creEnv *crelib.Environment,
) (*crelib.PreEnvStartupOutput, error) {
	if don.CapabilityConfigs == nil {
		don.CapabilityConfigs = make(map[crelib.CapabilityFlag]crelib.CapabilityConfig)
	}
	cfg, ok := don.CapabilityConfigs[crelib.ConfidentialRelayCapability]
	if !ok {
		cfg = crelib.CapabilityConfig{Values: make(map[string]any)}
	}
	if cfg.Values == nil {
		cfg.Values = make(map[string]any)
	}
	cfg.Values["trustedPCRs"] = f.pcrsJSON
	don.CapabilityConfigs[crelib.ConfidentialRelayCapability] = cfg

	return f.inner.PreEnvStartup(ctx, testLogger, don, topology, creEnv)
}

func (f *testConfidentialRelayFeature) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *crelib.Don,
	dons *crelib.Dons,
	creEnv *crelib.Environment,
) error {
	return f.inner.PostEnvStartup(ctx, testLogger, don, dons, creEnv)
}

// newTestConfidentialRelayFeature builds the relay feature for a set of enclaves.
// Fake enclaves emit a sentinel attestation document instead of real PCRs, so the
// trusted value is the fake measurements placeholder and attestation validation is
// relaxed. Real Nitro enclaves keep full validation against their measurements.
func newTestConfidentialRelayFeature(t *testing.T, enclaves []cctypes.Enclave, fake bool) crelib.Feature {
	t.Helper()

	var pcrsJSON string
	if fake {
		// Marshaling the raw "fake-measurements" bytes as json.RawMessage would fail
		// since it is not valid JSON, so encode it as a JSON string array.
		b, err := json.Marshal([]string{cctypes.FakeMeasurements})
		require.NoError(t, err, "failed to marshal fake measurements")
		pcrsJSON = string(b)
	} else {
		// Each enclave bakes in per-CID WireGuard keys, so measurements differ per
		// enclave. The relay accepts a JSON array and tries each until one matches.
		var allPCRs []json.RawMessage
		for _, enc := range enclaves {
			for _, tv := range enc.TrustedValues {
				if string(tv) != "invalid" {
					allPCRs = append(allPCRs, json.RawMessage(tv))
				}
			}
		}
		b, err := json.Marshal(allPCRs)
		require.NoError(t, err, "failed to marshal PCR measurements")
		pcrsJSON = string(b)
	}

	return crelib.Feature(&testConfidentialRelayFeature{
		inner:    confidentialrelay.ConfidentialRelay{TrustEnclaves: fake, RequireBFTQuorum: true},
		pcrsJSON: pcrsJSON,
	})
}

// ---------------------------------------------------------------------------
// Workflow artifacts
// ---------------------------------------------------------------------------

// confidentialWorkflowArtifacts holds everything derived from compiling the test
// workflow: the URLs the syncer and enclave fetch it from, and the on-disk
// directory the test copies into the DON containers.
type confidentialWorkflowArtifacts struct {
	BinaryURL   string
	ConfigURL   string
	ArtifactDir string
	BinaryHash  []byte
}

// buildAndServeConfidentialWorkflow compiles the test workflow to wasip1/wasm
// from the chainlink-confidential-compute checkout, brotli-compresses and
// base64-encodes it (the format the syncer and the enclave both expect), and
// serves the binary and its config over HTTP bound to 0.0.0.0 so the host, the
// Docker containers and the enclaves can all fetch them.
//
// Compiling from the checkout rather than vendoring the source keeps the
// workflow single-sourced and keeps its cre-sdk-go pins out of this
// repository's module graph.
func buildAndServeConfidentialWorkflow(t *testing.T, ccRoot string, configJSON string, hostIP string) confidentialWorkflowArtifacts {
	t.Helper()

	srcDir := filepath.Join(ccRoot, confidentialWorkflowSrcRelDir)
	require.DirExists(t, srcDir, "confidential workflow source not found in the chainlink-confidential-compute checkout")

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "workflow-test.wasm")

	cmd := exec.Command("go", "build", "-o", outFile, ".")
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm", "CGO_ENABLED=0")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "compiling confidential workflow WASM: %s", string(output))

	raw, err := os.ReadFile(outFile)
	require.NoError(t, err, "reading compiled WASM")

	var compressed bytes.Buffer
	w := brotli.NewWriter(&compressed)
	_, err = w.Write(raw)
	require.NoError(t, err, "brotli compressing WASM")
	require.NoError(t, w.Close(), "closing brotli writer")

	binary := compressed.Bytes()
	hash := sha256.Sum256(binary)
	encoded := base64.StdEncoding.EncodeToString(binary)

	// The syncer's file fetcher reads both files from disk inside the container,
	// so they have to exist as real files the test can copy in.
	require.NoError(t,
		os.WriteFile(filepath.Join(tmpDir, confidentialWorkflowBinaryFilename), []byte(encoded), 0o600),
		"staging workflow binary artifact")
	require.NoError(t,
		os.WriteFile(filepath.Join(tmpDir, confidentialWorkflowConfigFilename), []byte(configJSON), 0o600),
		"staging workflow config artifact")

	mux := http.NewServeMux()
	mux.HandleFunc("/"+confidentialWorkflowBinaryFilename, func(rw http.ResponseWriter, _ *http.Request) {
		_, _ = rw.Write([]byte(encoded))
	})
	mux.HandleFunc("/"+confidentialWorkflowConfigFilename, func(rw http.ResponseWriter, _ *http.Request) {
		_, _ = rw.Write([]byte(configJSON))
	})

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err, "workflow artifact listener")
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(listener) }()
	t.Cleanup(func() { _ = srv.Close() })

	port := listener.Addr().(*net.TCPAddr).Port
	base := fmt.Sprintf("http://%s:%d/", hostIP, port)

	return confidentialWorkflowArtifacts{
		BinaryURL:   base + confidentialWorkflowBinaryFilename,
		ConfigURL:   base + confidentialWorkflowConfigFilename,
		ArtifactDir: tmpDir,
		BinaryHash:  hash[:],
	}
}

// ---------------------------------------------------------------------------
// Enclave configuration
// ---------------------------------------------------------------------------

// configureEnclaves POSTs the enclave config to every enclave's config endpoint.
// Until this lands, an enclave has no signer set and no master public key, so it
// rejects every incoming compute request.
//
// The signer set is the workflow DON's worker P2P IDs, and F is derived as
// 2*don.F + 1 to match the relay DON quorum the enclave expects.
func configureEnclaves(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	configURLs []string,
	vaultPublicKey string,
) {
	t.Helper()

	workers, err := testEnv.Dons.MustWorkflowDON().Workers()
	require.NoError(t, err, "failed to get worker nodes from topology")
	require.NotEmpty(t, workers, "workflow DON has no worker nodes")

	signers := make([][]byte, 0, len(workers))
	for _, node := range workers {
		signers = append(signers, node.Keys.P2PKey.PeerID[:])
	}

	masterPublicKey, err := hex.DecodeString(vaultPublicKey)
	require.NoError(t, err, "failed to hex-decode vault public key")

	// don.F for an N-node DON is N/3; the enclave's own F/T is 2*don.F + 1.
	donF := uint32(len(workers) / 3)
	quorum := 2*donF + 1

	config := cctypes.EnclaveConfig{
		Signers:         signers,
		MasterPublicKey: masterPublicKey,
		T:               quorum,
		F:               quorum,
	}
	configBytes, err := json.Marshal(config)
	require.NoError(t, err, "failed to marshal enclave config")

	enclaveType := cctypes.EnclaveTypeNitro
	if UseFakeEnclaveForTest() {
		enclaveType = cctypes.EnclaveTypeFake
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // local test enclaves use self-signed certs
		},
	}

	for i, configURL := range configURLs {
		_, err := util.SetNodeConfig(
			t.Context(),
			cctypes.Enclave{
				EnclaveURL:    configURL,
				EnclaveType:   enclaveType,
				TrustedValues: [][]byte{},
				Region:        confidentialEnclaveRegion,
			},
			cctypes.ConfigRequest{Config: configBytes},
			&client,
		)
		require.NoError(t, err, "failed to set config on enclave %d (%s)", i, configURL)
		testLogger.Info().Int("enclave", i).Str("configURL", configURL).Msg("Enclave configured")
	}
}

// UseFakeEnclaveForTest reports whether the harness selected fake enclaves.
func UseFakeEnclaveForTest() bool {
	return testhelpers.UseFakeEnclave()
}

// storeConfidentialWorkflowSecret encrypts a secret to the vault's public key and
// stores it in the vault DON through the gateway, so the workflow's GetSecret call
// resolves. Reuses the vault request helpers already in this package.
func storeConfidentialWorkflowSecret(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	testLogger zerolog.Logger,
	gatewayURL string,
	vaultPublicKey string,
	secretKey string,
	secretValue string,
) {
	t.Helper()

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	sethClient := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	owner := sethClient.MustGetRootKeyAddress().Hex()

	wfRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(), //nolint:staticcheck // mirrors system-tests usage
		keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)
	wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sethClient.Client)
	require.NoError(t, err, "failed to build workflow registry wrapper")

	// The vault DON only accepts secrets from an owner linked in the registry.
	requireVaultLinkOwner(t, sethClient, common.HexToAddress(wfRegAddr),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()])

	parsedKey := mustVaultPublicKey(t, vaultPublicKey)
	encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner(secretValue, parsedKey, sethClient.MustGetRootKeyAddress())
	require.NoError(t, err, "failed to encrypt secret for the vault DON")

	auth := newAllowlistVaultRequestAuth(owner, sethClient, wfReg)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretKey, owner, gatewayURL, []string{"main"})

	testLogger.Info().Str("key", secretKey).Str("owner", owner).Msg("Stored workflow secret in the vault DON")
}
