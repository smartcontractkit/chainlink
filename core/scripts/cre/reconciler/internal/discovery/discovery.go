package discovery

import (
	"context"
	"runtime"
	"slices"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

// K8sClient interface for discovery (minimal subset needed).
type K8sClient interface {
	GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*infra.NodeAPIInfo, error)
}

// Dialer interface for discovering nodes.
type Dialer interface {
	Dial(apiURL, email, password string) (Client, error)
}

// Client reads key material from a node.
type Client interface {
	ReadCSAKey() (string, error)
	ReadPeerID() (string, error)
	ReadEVMAddresses() (map[string]string, error)
	ReadOCR2BundleIDs() (map[string]string, error)
	ReadAptosKeys() (string, error)
	ReadSolanaKeys() (string, error)
}

// Run discovers runtime info for all nodes concurrently (bounded), returning a map keyed by node name.
// Individual node failures are logged and skipped (matching current behavior: warn + continue), so a
// partial map is returned; the caller's completeness checks gate progress.
//
// nonEVMFamilies maps node name -> the non-EVM chain families ("aptos"/"solana") that
// node's DON actually declares a capability for, so discoverOne only calls the matching
// Read*Keys when the node is expected to have that key at all (calling it otherwise
// would just fail — most nodes have no Aptos/Solana key configured).
func Run(ctx context.Context, log zerolog.Logger, nodes []domain.ChartNodeInfo, cv *domain.ChartValues,
	k8s K8sClient, dialer Dialer, nonEVMFamilies map[string][]string) (map[string]domain.NodeRuntimeInfo, error) {
	type result struct {
		name string
		info domain.NodeRuntimeInfo
		ok   bool
	}
	results := make([]result, len(nodes))

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(min(len(nodes), 2*runtime.GOMAXPROCS(0)))
	for i := range nodes {
		node := nodes[i]
		idx := i
		g.Go(func() error {
			info, ok := discoverOne(gctx, log, node, cv, k8s, dialer, nonEVMFamilies[node.Name])
			results[idx] = result{name: node.Name, info: info, ok: ok}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	out := make(map[string]domain.NodeRuntimeInfo)
	for _, r := range results {
		if r.ok {
			out[r.name] = r.info
		}
	}
	return out, nil
}

// discoverOne contains the per-node discovery logic: K8s API lookup + node API dial + key reads.
func discoverOne(ctx context.Context, log zerolog.Logger, node domain.ChartNodeInfo, cv *domain.ChartValues,
	k8s K8sClient, dialer Dialer, neededNonEVMFamilies []string) (domain.NodeRuntimeInfo, bool) {
	var info domain.NodeRuntimeInfo
	info.NodeType = string(node.NodeType)

	namespace := cv.GetNodeNamespace(node.Name)

	apiInfo, err := k8s.GetNodeAPIInfo(ctx, node.Name, namespace)
	if err != nil {
		log.Warn().Err(err).Str("node", node.Name).Msg("Failed to get node API info")
		return info, false
	}

	info.APIURL = apiInfo.URL

	client, err := dialer.Dial(apiInfo.URL, apiInfo.Email, apiInfo.Password)
	if err != nil {
		log.Warn().Err(err).Str("node", node.Name).Msg("Failed to connect to node API")
		return info, false
	}

	// CSA key
	csaKey, err := client.ReadCSAKey()
	if err != nil {
		log.Warn().Err(err).Str("node", node.Name).Msg("Failed to read CSA key")
	} else if csaKey != "" {
		info.CSAKey = csaKey
	}

	// Peer ID
	peerID, err := client.ReadPeerID()
	if err != nil {
		log.Warn().Err(err).Str("node", node.Name).Msg("Failed to read P2P peer ID")
	} else if peerID != "" {
		info.PeerID = peerID
	}

	// EVM addresses
	evmAddrs, err := client.ReadEVMAddresses()
	if err != nil {
		log.Warn().Err(err).Str("node", node.Name).Msg("Failed to read EVM addresses")
	} else if len(evmAddrs) > 0 {
		info.EVMAddress = evmAddrs
	}

	// OCR2 bundle IDs — bootstrap and gateway nodes are not OCR signers and don't expose these.
	if node.NodeType != domain.RoleBootstrap && node.NodeType != domain.RoleGateway {
		ocr2Bundles, err := client.ReadOCR2BundleIDs()
		if err != nil {
			log.Warn().Err(err).Str("node", node.Name).Msg("Failed to read OCR2 bundle IDs")
		} else if len(ocr2Bundles) > 0 {
			info.OCR2BundleIDs = ocr2Bundles
		}
	}

	// Aptos/Solana addresses — only for nodes whose DON actually declares the matching
	// capability; most nodes have no such key configured and the read would just fail.
	if slices.Contains(neededNonEVMFamilies, "aptos") {
		aptosAddr, err := client.ReadAptosKeys()
		if err != nil {
			log.Warn().Err(err).Str("node", node.Name).Msg("Failed to read Aptos keys")
		} else if aptosAddr != "" {
			info.AptosAddress = aptosAddr
		}
	}
	if slices.Contains(neededNonEVMFamilies, "solana") {
		solanaAddr, err := client.ReadSolanaKeys()
		if err != nil {
			log.Warn().Err(err).Str("node", node.Name).Msg("Failed to read Solana keys")
		} else if solanaAddr != "" {
			info.SolanaAddress = solanaAddr
		}
	}

	return info, true
}
