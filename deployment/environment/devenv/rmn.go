package devenv

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	tcLog "github.com/testcontainers/testcontainers-go/log"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
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

type RMNKeys struct {
	OffchainPublicKey   ed25519.PublicKey
	EVMOnchainPublicKey common.Address
}

func GenerateRMNKeyStore(lggr zerolog.Logger, image string, version string, platform string) (keys RMNKeys, fileString string, passphrase string, err error) {
	container, err := docker.StartContainerWithRetry(lggr, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			AutoRemove: false,
			Image:      fmt.Sprintf("%s:%s", image, version),
			Env: map[string]string{
				"AFN_PASSPHRASE": DefaultAFNPassphrase,
			},
			Cmd:           []string{"afn2proxy", "--generate", "--keystore", RMNKeyStore},
			WaitingFor:    tcwait.ForExit(),
			ImagePlatform: platform,
		},
		Started: true,
		Logger:  &lggr,
	})
	defer func() {
		if terminateErr := container.Terminate(context.Background()); terminateErr != nil {
			log.Printf("Failed to stop container: %v", terminateErr)
		}
	}()

	if err != nil {
		return RMNKeys{}, "", "", err
	}

	// Copy the file from container
	reader, err := container.CopyFileFromContainer(context.Background(), "/app/"+RMNKeyStore)
	if err != nil {
		log.Printf("Failed to copy file: %v", err)
		return RMNKeys{}, "", "", err
	}
	defer reader.Close()

	fileContents, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Failed to read file contents: %v", err)
		return RMNKeys{}, "", "", err
	}

	fileString = string(fileContents)

	address, publicKey, err := extractKeys(fileContents)
	if err != nil {
		return RMNKeys{}, "", "", err
	}

	keys = RMNKeys{
		OffchainPublicKey:   publicKey,
		EVMOnchainPublicKey: address,
	}
	passphrase = DefaultAFNPassphrase

	return keys, fileString, passphrase, nil
}

func GeneratePeerID(lggr zerolog.Logger, image string, version string, imagePlatform string) (peerID p2ptypes.PeerID, content string, passphrase string, err error) {
	container, err := docker.StartContainerWithRetry(lggr, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			AutoRemove: false,
			Image:      fmt.Sprintf("%s:%s", image, version),
			Env: map[string]string{
				"RAGEPROXY_PASSPHRASE": DefaultAFNPassphrase,
			},
			Cmd:           []string{"rageproxy", "--generate", "--keystore", ProxyKeyStore},
			WaitingFor:    tcwait.ForExit(),
			ImagePlatform: imagePlatform,
		},
		Started: true,
		Logger:  &lggr,
	})
	defer (func() {
		err := container.Terminate(context.Background())
		if err != nil {
			log.Printf("Failed to stop container: %v", err)
		}
	})()

	if err != nil {
		return p2ptypes.PeerID{}, "", "", err
	}

	// Copy the file from container
	reader, err := container.CopyFileFromContainer(context.Background(), "/app/"+ProxyKeyStore)
	if err != nil {
		return p2ptypes.PeerID{}, "", "", err
	}
	defer reader.Close()

	fileContents, err := io.ReadAll(reader)
	if err != nil {
		return p2ptypes.PeerID{}, "", "", err
	}

	fileString := string(fileContents)

	peerID, err = extractPeerID(fileContents)
	if err != nil {
		return p2ptypes.PeerID{}, "", "", err
	}
	return peerID, fileString, DefaultAFNPassphrase, nil
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

	l := tcLog.Default()
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

	l := tcLog.Default()
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

func (rmn *RMNCluster) Restart(ctx context.Context) error {
	for _, node := range rmn.Nodes {
		_, _, err := node.RMN.Container.Exec(ctx, []string{"rm", "-f", "/app/cache/v4/*"})
		if err != nil {
			return err
		}

		err = node.RMN.Container.Stop(ctx, nil)
		if err != nil {
			return err
		}

		err = node.RMN.Container.Start(ctx)
		if err != nil {
			return err
		}
	}
	return nil
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
