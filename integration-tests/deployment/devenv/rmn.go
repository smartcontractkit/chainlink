package devenv

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
)

type RageProxy struct {
	test_env.EnvComponent
	AlwaysPullImage   bool
	proxyListenerPort string
	proxyPort         string
	Local             ProxyLocalConfig
	Shared            ProxySharedConfig
}

func NewRage2ProxyComponent(
	networks []string,
	name,
	imageName,
	imageVersion string,
	local ProxyLocalConfig,
	shared ProxySharedConfig,
	logStream *logstream.LogStream,
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
			LogStream:        logStream,
		},
		AlwaysPullImage:   true,
		proxyListenerPort: listenPort,
		proxyPort:         proxyPort,
		Local:             local,
		Shared:            shared,
	}
	return rmn, nil
}

type Keystore struct {
	AdditionalData map[string]interface{} `json:"additionalData"`
}

func (proxy *RageProxy) Start(t *testing.T, lggr zerolog.Logger) (tc.Container, error) {
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
			Name:            proxy.ContainerName,
			AlwaysPullImage: proxy.AlwaysPullImage,
			Image:           fmt.Sprintf("%s:%s", proxy.ContainerImage, proxy.ContainerVersion),
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
	// Cat inside container /app/keystore/rageproxy-keystore.json
	//_, reader, err := container.Exec(context.Background(), []string{"cat", "/app/keystore/rageproxy-keystore.json"})
	//if err != nil {
	//	return nil, err
	//}
	//b, err := ioutil.ReadAll(reader)
	//if err != nil {
	//	return nil, err
	//}
	//json.Unmarshal(b)
	proxy.Container = container
	return container, nil
}

type AFN2Proxy struct {
	test_env.EnvComponent
	AlwaysPullImage bool
	AFNPassphrase   string
	Shared          SharedConfig
	Local           LocalConfig
}

func NewAFN2ProxyComponent(
	networks []string,
	name,
	imageName,
	imageVersion string,
	shared SharedConfig,
	local LocalConfig,
	logStream *logstream.LogStream) (*AFN2Proxy, error) {
	afnName := fmt.Sprintf("%s-%s", name, uuid.NewString()[0:8])
	rmn := &AFN2Proxy{
		EnvComponent: test_env.EnvComponent{
			ContainerName:    afnName,
			ContainerImage:   imageName,
			ContainerVersion: imageVersion,
			Networks:         networks,
			LogStream:        logStream,
		},
		AlwaysPullImage: true,
		AFNPassphrase:   DefaultAFNPasphrase,
		Shared:          shared,
		Local:           local,
	}

	return rmn, nil
}

func (rmn *AFN2Proxy) Start(t *testing.T, lggr zerolog.Logger, reuse bool) (tc.Container, error) {
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
			Name:            rmn.ContainerName,
			AlwaysPullImage: rmn.AlwaysPullImage,
			Image:           fmt.Sprintf("%s:%s", rmn.ContainerImage, rmn.ContainerVersion),
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
	rmn.Container = container
	return container, nil
}

type RMNNode struct {
	AFN   AFN2Proxy
	Proxy RageProxy
}

func (n *RMNNode) Start(t *testing.T, lggr zerolog.Logger) error {
	_, err := n.Proxy.Start(t, lggr)
	if err != nil {
		return err
	}
	_, err = n.AFN.Start(t, lggr, false)
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
	logStream *logstream.LogStream,
) (*RMNCluster, error) {
	rmn := &RMNCluster{
		t:     t,
		l:     l,
		Nodes: make(map[string]RMNNode),
	}
	for name, rmnConfig := range config {
		proxy, err := NewRage2ProxyComponent(networks, name, proxyImage, proxyVersion, rmnConfig.ProxyLocal, rmnConfig.ProxyShared, logStream)
		if err != nil {
			return nil, err
		}
		_, err = proxy.Start(t, l)
		if err != nil {
			return nil, err
		}
		proxyName, err := proxy.Container.Name(context.Background())
		if err != nil {
			return nil, err
		}
		_, port, err := net.SplitHostPort(rmnConfig.Shared.Networking.RageProxy)
		if err != nil {
			return nil, err
		}
		rmnConfig.Shared.Networking.RageProxy = fmt.Sprintf("%s:%s", proxyName, port)
		afn, err := NewAFN2ProxyComponent(networks, name, rmnImage, rmnVersion, rmnConfig.Shared, rmnConfig.Local, logStream)
		if err != nil {
			return nil, err
		}
		_, err = afn.Start(t, l, false)
		if err != nil {
			return nil, err
		}
		rmn.Nodes[name] = RMNNode{
			AFN:   *afn,
			Proxy: *proxy,
		}
	}
	return rmn, nil
}
