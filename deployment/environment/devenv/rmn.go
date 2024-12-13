package devenv

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/exec"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	RMNKeyStore   = "keystore/afn2proxy-keystore.json"
	ProxyKeyStore = "keystore/rageproxy-keystore.json"
)

type RageProxy struct {
	test_env.EnvComponent
	proxyListenerPort string
	proxyPort         string
	Passphrase        string
	Local             ProxyLocalConfig
	Shared            ProxySharedConfig

	// Generated on first time boot.
	// Needed for RMHHome.
	PeerID p2ptypes.PeerID
}

func NewRage2ProxyComponent(
	networks []string,
	name,
	imageName,
	imageVersion string,
	local ProxyLocalConfig,
	shared ProxySharedConfig,
) (*RageProxy, error) {
	rageName := fmt.Sprintf("%s-proxy-%s", name, uuid.NewString()[0:8])

	// TODO support multiple listeners
	_, listenPort, err := net.SplitHostPort(local.ListenAddresses[0])
	if err != nil {
		return nil, err
	}
	_, proxyPort, err := net.SplitHostPort(local.ProxyAddress)
	if err != nil {
		return nil, err
	}

	rmn := &RageProxy{
		EnvComponent: test_env.EnvComponent{
			ContainerName:    rageName,
			ContainerImage:   imageName,
			ContainerVersion: imageVersion,
			Networks:         networks,
		},
		Passphrase:        DefaultAFNPassphrase,
		proxyListenerPort: listenPort,
		proxyPort:         proxyPort,
		Local:             local,
		Shared:            shared,
	}
	return rmn, nil
}

func extractPeerID(b []byte) (p2ptypes.PeerID, error) {
	var keystore struct {
		AdditionalData string `json:"additionalData"`
	}
	if err := json.Unmarshal(b, &keystore); err != nil {
		return p2ptypes.PeerID{}, err
	}
	var additionalData struct {
		PeerID string `json:"PeerID"`
	}
	if err := json.Unmarshal([]byte(keystore.AdditionalData), &additionalData); err != nil {
		return p2ptypes.PeerID{}, err
	}
	var peerID p2ptypes.PeerID
	if err := peerID.UnmarshalText([]byte(additionalData.PeerID)); err != nil {
		return p2ptypes.PeerID{}, err
	}
	return peerID, nil
}

func (proxy *RageProxy) Start(t *testing.T, lggr zerolog.Logger, networks []string) (tc.Container, error) {
	sharedRageProxy, err := proxy.Shared.rageProxyShared()
	if err != nil {
		return nil, err
	}
	localRageProxy, err := proxy.Local.rageProxyLocal()
	if err != nil {
		return nil, err
	}

	l := tc.Logger
	if t != nil {
		l = logging.CustomT{
			T: t,
			L: lggr,
		}
	}
	container, err := docker.StartContainerWithRetry(lggr, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Name:     proxy.ContainerName,
			Networks: networks,
			Image:    fmt.Sprintf("%s:%s", proxy.ContainerImage, proxy.ContainerVersion),
			Env: map[string]string{
				"RAGEPROXY_PASSPHRASE": proxy.Passphrase,
			},
			ExposedPorts: []string{
				test_env.NatPortFormat(proxy.proxyPort),
				test_env.NatPortFormat(proxy.proxyListenerPort),
			},
			Files: []tc.ContainerFile{
				{
					HostFilePath:      sharedRageProxy,
					ContainerFilePath: "/app/cfg/rageproxy-shared.json",
					FileMode:          0644,
				},
				{
					HostFilePath:      localRageProxy,
					ContainerFilePath: "/app/cfg/rageproxy-local.json",
					FileMode:          0644,
				},
			},
			WaitingFor: tcwait.ForExec([]string{"cat", ProxyKeyStore}),
			LifecycleHooks: []tc.ContainerLifecycleHooks{
				{
					PostStarts:    proxy.PostStartsHooks,
					PostStops:     proxy.PostStopsHooks,
					PreTerminates: proxy.PreTerminatesHooks,
				},
			},
		},
		Started: true,
		Logger:  l,
	})
	if err != nil {
		return nil, err
	}
	_, reader, err := container.Exec(context.Background(), []string{
		"cat", ProxyKeyStore}, exec.Multiplexed())
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to cat keystore")
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	peerID, err := extractPeerID(b)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to extract peerID %s", string(b))
	}
	proxy.PeerID = peerID
	proxy.Container = container
	return container, nil
}

type AFN2Proxy struct {
	test_env.EnvComponent
	AFNPassphrase string
	Shared        SharedConfig
	Local         LocalConfig

	// Generated on boot
	OffchainPublicKey   ed25519.PublicKey // RMNHome
	EVMOnchainPublicKey common.Address    // RMNRemote
}

func NewAFN2ProxyComponent(
	networks []string,
	name,
	imageName,
	imageVersion string,
	shared SharedConfig,
	local LocalConfig) (*AFN2Proxy, error) {
	afnName := fmt.Sprintf("%s-%s", name, uuid.NewString()[0:8])
	rmn := &AFN2Proxy{
		EnvComponent: test_env.EnvComponent{
			ContainerName:    afnName,
			ContainerImage:   imageName,
			ContainerVersion: imageVersion,
			Networks:         networks,
		},
		AFNPassphrase: DefaultAFNPassphrase,
		Shared:        shared,
		Local:         local,
	}

	return rmn, nil
}

func extractKeys(b []byte) (common.Address, ed25519.PublicKey, error) {
	var keystore struct {
		AssociatedData string `json:"associated_data"`
	}
	if err := json.Unmarshal(b, &keystore); err != nil {
		return common.Address{}, ed25519.PublicKey{}, err
	}
	var associatedData struct {
		OffchainPublicKey   string `json:"offchain_public_key"`
		EVMOnchainPublicKey string `json:"evm_onchain_public_key"`
	}
	if err := json.Unmarshal([]byte(keystore.AssociatedData), &associatedData); err != nil {
		return common.Address{}, ed25519.PublicKey{}, err
	}
	offchainKey, err := hexutil.Decode(associatedData.OffchainPublicKey)
	if err != nil {
		return common.Address{}, ed25519.PublicKey{}, err
	}
	if len(offchainKey) != ed25519.PublicKeySize {
		return common.Address{}, ed25519.PublicKey{}, fmt.Errorf("invalid offchain public key: %x", offchainKey)
	}
	return common.HexToAddress(associatedData.EVMOnchainPublicKey), offchainKey, nil
}

func (rmn *AFN2Proxy) Start(t *testing.T, lggr zerolog.Logger, reuse bool, networks []string) (tc.Container, error) {
	localAFN2Proxy, err := rmn.Local.afn2ProxyLocalConfigFile()
	if err != nil {
		return nil, err
	}
	sharedAFN2Proxy, err := rmn.Shared.afn2ProxySharedConfigFile()
	if err != nil {
		return nil, err
	}

	l := tc.Logger
	if t != nil {
		l = logging.CustomT{
			T: t,
			L: lggr,
		}
	}
	container, err := docker.StartContainerWithRetry(lggr, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Name:     rmn.ContainerName,
			Networks: networks,
			Image:    fmt.Sprintf("%s:%s", rmn.ContainerImage, rmn.ContainerVersion),
			Env: map[string]string{
				"AFN_PASSPHRASE": rmn.AFNPassphrase,
			},
			Files: []tc.ContainerFile{
				{
					HostFilePath:      sharedAFN2Proxy,
					ContainerFilePath: "/app/cfg/afn2proxy-shared.toml",
					FileMode:          0644,
				},
				{
					HostFilePath:      localAFN2Proxy,
					ContainerFilePath: "/app/cfg/afn2proxy-local.toml",
					FileMode:          0644,
				},
			},
			WaitingFor: tcwait.ForExec([]string{"cat", RMNKeyStore}),
			LifecycleHooks: []tc.ContainerLifecycleHooks{
				{
					PostStarts:    rmn.PostStartsHooks,
					PostStops:     rmn.PostStopsHooks,
					PreTerminates: rmn.PreTerminatesHooks,
				},
			},
		},
		Started: true,
		Reuse:   reuse,
		Logger:  l,
	})
	if err != nil {
		return nil, err
	}
	_, reader, err := container.Exec(context.Background(), []string{
		"cat", RMNKeyStore}, exec.Multiplexed())
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to cat keystore")
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	onchainPubKey, offchainPubKey, err := extractKeys(b)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to extract peerID %s", string(b))
	}
	rmn.OffchainPublicKey = offchainPubKey
	rmn.EVMOnchainPublicKey = onchainPubKey
	rmn.Container = container
	return container, nil
}

type RMNNode struct {
	RMN   AFN2Proxy
	Proxy RageProxy
}

func (n *RMNNode) Start(t *testing.T, lggr zerolog.Logger, networks []string) error {
	_, err := n.Proxy.Start(t, lggr, networks)
	if err != nil {
		return err
	}

	_, err = n.RMN.Start(t, lggr, false, networks)
	if err != nil {
		return err
	}
	return nil
}

type RMNCluster struct {
	Nodes map[string]RMNNode
	t     *testing.T
	l     zerolog.Logger
}

// NewRMNCluster creates a new RMNCluster with the given configuration
// and starts it.
func NewRMNCluster(
	t *testing.T,
	l zerolog.Logger,
	networks []string,
	config map[string]RMNConfig,
	proxyImage string,
	proxyVersion string,
	rmnImage string,
	rmnVersion string,
) (*RMNCluster, error) {
	rmn := &RMNCluster{
		t:     t,
		l:     l,
		Nodes: make(map[string]RMNNode),
	}
	for name, rmnConfig := range config {
		proxy, err := NewRage2ProxyComponent(networks, name, proxyImage, proxyVersion, rmnConfig.ProxyLocal, rmnConfig.ProxyShared)
		if err != nil {
			return nil, err
		}
		_, err = proxy.Start(t, l, networks)
		if err != nil {
			return nil, err
		}

		// TODO: Hack here is we overwrite the host with the container name
		// since the RMN node needs to be able to reach its own proxy container.
		proxyName, err := proxy.Container.Name(context.Background())
		if err != nil {
			return nil, err
		}
		_, port, err := net.SplitHostPort(rmnConfig.Local.Networking.RageProxy)
		if err != nil {
			return nil, err
		}
		rmnConfig.Local.Networking.RageProxy = strings.TrimPrefix(fmt.Sprintf("%s:%s", proxyName, port), "/")
		afn, err := NewAFN2ProxyComponent(networks, name, rmnImage, rmnVersion, rmnConfig.Shared, rmnConfig.Local)
		if err != nil {
			return nil, err
		}
		_, err = afn.Start(t, l, false, networks)
		if err != nil {
			return nil, err
		}
		rmn.Nodes[name] = RMNNode{
			RMN:   *afn,
			Proxy: *proxy,
		}
	}
	return rmn, nil
}
