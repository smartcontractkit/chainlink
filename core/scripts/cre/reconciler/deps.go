package reconciler

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

// K8sAPI is the Kubernetes surface the reconciler needs. Implemented by *infra.K8sClient;
// fakes implement it in tests.
type K8sAPI interface {
	GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*infra.NodeAPIInfo, error)
	GetNodeSecretsToml(ctx context.Context, nodeName, namespace string) (string, error)
	RestartNodePods(ctx context.Context, nodeName, namespace string) error
	CopyFilesToPod(ctx context.Context, namespace, podName, container, destDir string, localPaths []string) error
}

// NodeClient reads key material from a single node's admin API. Implemented by a thin
// wrapper over chainlink-testing-framework clclient (see nodeapi.go); fakes implement it in tests.
type NodeClient interface {
	ReadCSAKey() (string, error)                   // returns key WITHOUT "csa_" prefix, "" if none
	ReadPeerID() (string, error)                   // P2P PeerID, "" if none
	ReadEVMAddresses() (map[string]string, error)  // chainID(string) -> address
	ReadOCR2BundleIDs() (map[string]string, error) // chain family (lowercase) -> bundle ID
}

// NodeDialer opens a NodeClient for a node's API endpoint.
type NodeDialer interface {
	Dial(apiURL, email, password string) (NodeClient, error)
}
