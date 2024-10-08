package devenv

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	DefaultAFNPasphrase          = "my-not-so-secret-passphrase"
	DefaultRageProxyPort         = "8081"
	DefaultRageProxyListenerPort = "8080"
	DefaultDiscovererDbPath      = "/app/rageproxy-discoverer-db.json"
)

var (
	DefaultRageProxySharedConfig = ProxySharedConfig{
		Host: HostConfig{
			DurationBetweenDials: 1000000000,
		},
		Discoverer: DiscovererConfig{
			DeltaReconcile: 1000000000,
		},
	}
)

type Networking struct {
	RageProxy     string   `toml:"rageproxy"`
	Bootstrappers []string `toml:"bootstrappers"`
}

type HomeChain struct {
	Name                 string `toml:"name"`
	CapabilitiesRegistry string `toml:"capabilities_registry"`
	CCIPConfig           string `toml:"ccip_config"`
	RMNHome              string `toml:"rmn_home"`
}

type Stability struct {
	Type string `toml:"type"`
}

type RemoteChain struct {
	Name             string    `toml:"name"`
	Stability        Stability `toml:"stability"`
	StartBlockNumber int       `toml:"start_block_number"`
	OffRamp          string    `toml:"off_ramp"`
	RMNRemote        string    `toml:"rmn_remote"`
}

type SharedConfig struct {
	Networking   Networking    `toml:"networking"`
	HomeChain    HomeChain     `toml:"home_chain"`
	RemoteChains []RemoteChain `toml:"remote_chains"`
}

type LocalConfig struct {
	Chains []Chain `toml:"chains"`
}

type Chain struct {
	Name string `toml:"name"`
	RPC  string `toml:"rpc"`
}

type ProxyLocalConfig struct {
	ListenAddresses   []string `json:"ListenAddresses"`
	AnnounceAddresses []string `json:"AnnounceAddresses"`
	ProxyAddress      string   `json:"ProxyAddress"`
	DiscovererDbPath  string   `json:"DiscovererDbPath"`
}

type ProxySharedConfig struct {
	Host       HostConfig       `json:"Host"`
	Discoverer DiscovererConfig `json:"Discoverer"`
}

type HostConfig struct {
	DurationBetweenDials int64 `json:"DurationBetweenDials"`
}

type DiscovererConfig struct {
	DeltaReconcile int64 `json:"DeltaReconcile"`
}

type RMNConfig struct {
	Shared      SharedConfig
	Local       LocalConfig
	ProxyLocal  ProxyLocalConfig
	ProxyShared ProxySharedConfig
}

func (rmn RMNConfig) afn2ProxyLocalConfigFile() (string, error) {
	data, err := toml.Marshal(rmn.Local)
	if err != nil {
		return "", fmt.Errorf("failed to marshal afn2Proxy local config: %v", err)
	}
	return CreateTempFile(data, "afn2proxy_local")
}

func (rmn RMNConfig) afn2ProxySharedConfigFile() (string, error) {
	data, err := toml.Marshal(rmn.Shared)
	if err != nil {
		return "", fmt.Errorf("failed to marshal afn2Proxy shared config: %v", err)
	}
	return CreateTempFile(data, "afn2proxy_shared")
}

func (rmn RMNConfig) rageProxyShared() (string, error) {
	data, err := json.Marshal(rmn.ProxyShared)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rageProxy shared config: %v", err)
	}
	return CreateTempFile(data, "rageproxy_shared")
}

func (rmn RMNConfig) rageProxyLocal() (string, error) {
	data, err := json.Marshal(rmn.ProxyLocal)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rageProxy local config: %v", err)
	}
	return CreateTempFile(data, "rageproxy_local")
}

func CreateTempFile(data []byte, pattern string) (string, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for %s: %v", pattern, err)
	}
	_, err = file.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to write  %s: %v", pattern, err)
	}
	return file.Name(), nil
}

type RMNNode struct {
	PeerID p2ptypes.PeerID
}

type RMNClusterInput struct {
	Config       map[string]RMNConfig
	ProxyImage   string
	ProxyVersion string
	AFNImage     string
	AFNVersion   string
}

func NewRMNClusterInput(config ccip.RMNConfig) RMNClusterInput {
	input := RMNClusterInput{
		Config: make(map[string]RMNConfig),
	}
	for i := 1; i <= *config.NoOfNodes; i++ {
		nodeName := fmt.Sprintf("rmn-node%d", i)
		input.Config[nodeName] = RMNConfig{}
		input.ProxyImage = *config.ProxyImage
		input.ProxyVersion = *config.ProxyVersion
		input.AFNImage = *config.AFNImage
		input.AFNVersion = *config.AFNVersion
	}
	return input
}

type RageProxy struct {
	test_env.EnvComponent
	AlwaysPullImage   bool
	proxyListenerPort string
	proxyPort         string
	Local             ProxyLocalConfig
	Shared            ProxySharedConfig
}

func (proxy *RageProxy) Start(t *testing.T, lggr zerolog.Logger, config RMNConfig) (tc.Container, error) {
	cReq, err := proxy.getContainerRequest(config)
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
		ContainerRequest: *cReq,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return nil, err
	}
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

func (afn *AFN2Proxy) Start(t *testing.T, lggr zerolog.Logger, config RMNConfig, reuse bool) (tc.Container, error) {
	cReq, err := afn.getContainerRequest(config)
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
		ContainerRequest: *cReq,
		Started:          true,
		Reuse:            reuse,
		Logger:           l,
	})
	if err != nil {
		return nil, err
	}
	afn.Container = container
	return container, nil
}

type RMNOptions = func(rmn *RMNOffChainCluster)

type RMNNodeComponents struct {
	Config *RMNConfig
	AFN    *AFN2Proxy
	Proxy  *RageProxy
}

type RMNOffChainCluster struct {
	components map[string]RMNNodeComponents
	t          *testing.T
	l          zerolog.Logger
}

func (rmn *RMNOffChainCluster) Start(t *testing.T, lggr zerolog.Logger) error {
	for _, component := range rmn.components {
		_, err := component.Proxy.Start(t, lggr, *component.Config)
		if err != nil {
			return err
		}
		_, err = component.AFN.Start(t, lggr, *component.Config, false)
		if err != nil {
			return err
		}
	}
	// TODO Return RMNNode with p2p peerID
	return nil
}

func WithAddedChainRPC(nodeName string, chains ...Chain) RMNOptions {
	return func(rmn *RMNOffChainCluster) {
		comp := rmn.components[nodeName]
		comp.Config.Local.Chains = append(comp.Config.Local.Chains, chains...)
		rmn.components[nodeName] = comp
	}
}

func WithCCIPState(t *testing.T, state ccipdeployment.CCIPOnChainState, homeChainSel uint64) RMNOptions {
	return func(rmn *RMNOffChainCluster) {
		for k, comp := range rmn.components {
			homeChainState, ok := state.Chains[homeChainSel]
			require.Truef(t, ok, "home chain %d not found in state", homeChainSel)
			chain, found := chainsel.ChainBySelector(homeChainSel)
			require.Truef(t, found, "chain name not found for chain %d", homeChainSel)
			capReg := homeChainState.CapabilityRegistry
			require.NotNil(t, capReg, "capability registry not found for home chain %d", homeChainSel)
			ccipConfig := homeChainState.CCIPConfig
			require.NotNil(t, ccipConfig, "ccip config not found for home chain %d", homeChainSel)
			// TODO: Add RMNHome to CCIPOnChainState
			/*rmnHome := homeChainState.RMNHome
			require.NotNil(t, rmnHome, "rmn home not found for home chain %d", homeChainSel)*/
			comp.Config.Shared.HomeChain = HomeChain{
				Name:                 chain.Name,
				CapabilitiesRegistry: capReg.Address().String(),
				CCIPConfig:           ccipConfig.Address().String(),
				RMNHome:              "", // TODO: Add RMNHome to CCIPOnChainState
			}
			for sel, remoteChain := range state.Chains {
				rChain, found := chainsel.ChainBySelector(sel)
				require.Truef(t, found, "chain name not found for chain %d", sel)
				offRamp := remoteChain.OffRamp
				require.NotNil(t, offRamp, "off ramp not found for remote chain %d", sel)
				rmnRemote := remoteChain.RMNRemote
				require.NotNil(t, rmnRemote, "rmn remote not found for remote chain %d", sel)
				comp.Config.Shared.RemoteChains = append(comp.Config.Shared.RemoteChains, RemoteChain{
					Name:      rChain.Name,
					OffRamp:   offRamp.Address().String(),
					RMNRemote: rmnRemote.Address().String(),
					// TODO Add Stability to CCIPOnChainState
					Stability: Stability{
						Type: "FinalityTag",
					},
					// TODO: Add StartBlockNumber to CCIPOnChainState
					StartBlockNumber: 0,
				})
			}
			rmn.components[k] = comp
		}
	}
}

func WithAddedBootstrapper(bootstrappers ...string) RMNOptions {
	return func(rmn *RMNOffChainCluster) {
		for nodeName := range rmn.components {
			comp := rmn.components[nodeName]
			comp.Config.Shared.Networking.Bootstrappers = append(comp.Config.Shared.Networking.Bootstrappers, bootstrappers...)
			rmn.components[nodeName] = comp
		}
	}
}

func WithRageProxyListenerPort(port string) RMNOptions {
	return func(rmn *RMNOffChainCluster) {
		for k, comp := range rmn.components {
			comp.Proxy.proxyListenerPort = port
			comp.Proxy.Local.ListenAddresses = []string{fmt.Sprintf("127.0.0.1:%s", port)}
			rmn.components[k] = comp
		}
	}
}

func WithRageProxyPort(port string) RMNOptions {
	return func(rmn *RMNOffChainCluster) {
		for k, comp := range rmn.components {
			comp.Proxy.proxyPort = port
			comp.Proxy.Local.ProxyAddress = fmt.Sprintf("127.0.0.1:%s", port)
			comp.Config.Shared.Networking.RageProxy = comp.Proxy.Local.ProxyAddress
			rmn.components[k] = comp
		}
	}
}

func WithDiscovererDbPath(path string) RMNOptions {
	return func(rmn *RMNOffChainCluster) {
		for k, comp := range rmn.components {
			comp.Proxy.Local.DiscovererDbPath = path
			rmn.components[k] = comp
		}
	}
}

func NewRMNCluster(
	t *testing.T,
	l zerolog.Logger,
	networks []string,
	input RMNClusterInput,
	logStream *logstream.LogStream,
	opts ...RMNOptions,
) (*RMNOffChainCluster, error) {
	rmn := &RMNOffChainCluster{
		t:          t,
		l:          l,
		components: make(map[string]RMNNodeComponents),
	}
	for name, config := range input.Config {
		afn, err := NewAFN2ProxyComponent(networks, name, input.AFNImage, input.AFNVersion, config.Shared, config.Local, logStream)
		if err != nil {
			return nil, err
		}
		proxy, err := NewRage2ProxyComponent(networks, name, input.ProxyImage, input.ProxyVersion, config.ProxyLocal, config.ProxyShared, logStream)
		if err != nil {
			return nil, err
		}
		rmn.components[name] = RMNNodeComponents{
			Config: &config,
			AFN:    afn,
			Proxy:  proxy,
		}
	}

	for _, opt := range opts {
		opt(rmn)
	}
	return rmn, nil
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
	rmn := &RageProxy{
		EnvComponent: test_env.EnvComponent{
			ContainerName:    rageName,
			ContainerImage:   imageName,
			ContainerVersion: imageVersion,
			Networks:         networks,
			LogStream:        logStream,
		},
		AlwaysPullImage:   true,
		proxyListenerPort: DefaultRageProxyListenerPort,
		proxyPort:         DefaultRageProxyPort,
		Local:             local,
		Shared:            shared,
	}
	if rmn.Local.ListenAddresses == nil {
		rmn.Local.ListenAddresses = []string{fmt.Sprintf("127.0.0.1:%s", DefaultRageProxyListenerPort)}
	}
	if rmn.Local.ProxyAddress == "" {
		rmn.Local.ProxyAddress = fmt.Sprintf("127.0.0.1:%s", DefaultRageProxyPort)
	}
	if rmn.Local.DiscovererDbPath == "" {
		rmn.Local.DiscovererDbPath = DefaultDiscovererDbPath
	}
	if rmn.Shared == (ProxySharedConfig{}) {
		rmn.Shared = DefaultRageProxySharedConfig
	}
	return rmn, nil
}

func (rmn *AFN2Proxy) getContainerRequest(config RMNConfig) (*tc.ContainerRequest, error) {
	localAFN2Proxy, err := config.afn2ProxyLocalConfigFile()
	if err != nil {
		return nil, err
	}
	sharedAFN2Proxy, err := config.afn2ProxySharedConfigFile()
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
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
	}, nil
}

func (proxy *RageProxy) getContainerRequest(config RMNConfig) (*tc.ContainerRequest, error) {
	sharedRageProxy, err := config.rageProxyShared()
	if err != nil {
		return nil, err
	}
	localRageProxy, err := config.rageProxyLocal()
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
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
	}, nil
}
