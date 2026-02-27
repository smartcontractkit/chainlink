package pkg

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
)

const (
	GatewayHandlerTypeWebAPICapabilities = "web-api-capabilities"
	GatewayHandlerTypeHTTPCapabilities   = "http-capabilities"
	GatewayHandlerTypeVault              = "vault"

	ServiceNameWorkflows = "workflows"
	ServiceNameVault     = "vault"

	minimumRequestTimeoutSec = 5
)

// HandlerServiceName returns the service name for a given handler type.
func HandlerServiceName(handlerType string) string {
	switch handlerType {
	case GatewayHandlerTypeVault:
		return ServiceNameVault
	case GatewayHandlerTypeHTTPCapabilities, GatewayHandlerTypeWebAPICapabilities:
		return ServiceNameWorkflows
	default:
		return handlerType
	}
}

type TargetDONMember struct {
	Address string
	Name    string
}

type TargetDON struct {
	ID      string
	F       int
	Members []TargetDONMember
}

type GatewayServiceConfig struct {
	ServiceName string
	Handlers    []string
	DONs        []string
}

type GatewayJob struct {
	DONs              []TargetDON
	Services          []GatewayServiceConfig
	JobName           string
	RequestTimeoutSec int
	AllowedPorts      []int
	AllowedSchemes    []string
	AllowedIPsCIDR    []string
	AuthGatewayID     string
	ExternalJobID     string
}

func (g GatewayJob) Validate() error {
	if g.JobName == "" {
		return errors.New("must provide job name")
	}

	if len(g.DONs) == 0 {
		return errors.New("must provide at least one DON")
	}

	if len(g.Services) == 0 {
		return errors.New("must provide at least one service")
	}

	// We impose a lower bound to account for other timeouts which are hardcoded,
	// including Read/WriteTimeoutMillis, and handler-specific timeouts like the vault handler timeout.
	if g.RequestTimeoutSec < minimumRequestTimeoutSec {
		return errors.New("request timeout must be at least" + strconv.Itoa(minimumRequestTimeoutSec) + " seconds")
	}

	for _, port := range g.AllowedPorts {
		if port < 1 || port > 65535 {
			return errors.New("allowed port out of range: " + strconv.Itoa(port))
		}
	}

	for _, scheme := range g.AllowedSchemes {
		if scheme != "http" && scheme != "https" {
			return errors.New("allowed scheme must be either http or https: " + scheme)
		}
	}

	if len(g.AllowedIPsCIDR) > 0 {
		for _, cidr := range g.AllowedIPsCIDR {
			if _, _, err := net.ParseCIDR(cidr); err != nil {
				return errors.New("invalid CIDR format: " + cidr)
			}
		}
	}

	return nil
}

func (g GatewayJob) Resolve(gatewayNodeIdx int) (string, error) {
	externalJobID := g.ExternalJobID
	if externalJobID == "" {
		externalJobID = uuid.NewSHA1(uuid.Nil, []byte(g.JobName)).String()
	}

	shardedDONs, services, err := g.buildServicesAndShardedDONs()
	if err != nil {
		return "", err
	}

	requestTimeout := time.Duration(g.RequestTimeoutSec) * time.Second
	config := gatewayConfig{
		ConnectionManagerConfig: connectionManagerConfig{
			AuthChallengeLen:          10,
			AuthGatewayID:             "gateway-node-" + strconv.Itoa(gatewayNodeIdx),
			AuthTimestampToleranceSec: 5,
			HeartbeatIntervalSec:      20,
		},
		NodeServerConfig: nodeServerConfig{
			HandshakeTimeoutMillis: 1_000,
			MaxRequestBytes:        100_000,
			Path:                   "/",
			Port:                   5_003,
			ReadTimeoutMillis:      1_000,
			RequestTimeoutMillis:   int(requestTimeout.Milliseconds()),
			WriteTimeoutMillis:     1_000,
		},
		UserServerConfig: userServerConfig{
			ContentTypeHeader:    "application/jsonrpc",
			MaxRequestBytes:      100_000,
			Path:                 "/",
			Port:                 5_002,
			ReadTimeoutMillis:    int(requestTimeout.Milliseconds()),
			RequestTimeoutMillis: int(requestTimeout.Milliseconds()),
			WriteTimeoutMillis:   int(requestTimeout.Milliseconds() + 1000),
		},
		HTTPClientConfig: httpClientConfig{
			MaxResponseBytes: 50_000_000,
			AllowedPorts:     []int{443},
			AllowedSchemes:   []string{"https"},
		},
		ShardedDONs: shardedDONs,
		Services:    services,
	}

	if len(g.AllowedPorts) > 0 {
		config.HTTPClientConfig.AllowedPorts = g.AllowedPorts
	}

	if len(g.AllowedSchemes) > 0 {
		config.HTTPClientConfig.AllowedSchemes = g.AllowedSchemes
	}

	if len(g.AllowedIPsCIDR) > 0 {
		config.HTTPClientConfig.AllowedIPsCIDR = g.AllowedIPsCIDR
	}

	if g.AuthGatewayID != "" {
		config.ConnectionManagerConfig.AuthGatewayID = g.AuthGatewayID
	}

	spec := &gatewaySpec{
		Type:              "gateway",
		SchemaVersion:     1,
		Name:              g.JobName,
		ExternalJobID:     externalJobID,
		ForwardingAllowed: false,
		GatewayConfig:     config,
	}
	b, marshalErr := toml.Marshal(spec)
	if marshalErr != nil {
		return "", marshalErr
	}

	return string(b), nil
}

func (g GatewayJob) buildServicesAndShardedDONs() ([]shardedDON, []service, error) {
	shardedDONs := make([]shardedDON, len(g.DONs))
	for i, don := range g.DONs {
		nodes := make([]member, len(don.Members))
		for j, mem := range don.Members {
			nodes[j] = member(mem)
		}
		shardedDONs[i] = shardedDON{
			DonName: don.ID,
			F:       don.F,
			Shards:  []shard{{Nodes: nodes}},
		}
	}

	services := make([]service, 0, len(g.Services))
	for _, svcCfg := range g.Services {
		var handlers []handler
		for _, ht := range svcCfg.Handlers {
			switch ht {
			case GatewayHandlerTypeWebAPICapabilities:
				handlers = append(handlers, newDefaultWebAPICapabilitiesHandler())
			case GatewayHandlerTypeVault:
				handlers = append(handlers, newDefaultVaultHandler(g.RequestTimeoutSec))
			case GatewayHandlerTypeHTTPCapabilities:
				handlers = append(handlers, newDefaultHTTPCapabilitiesHandler())
			default:
				return nil, nil, errors.New("unknown handler type: " + ht)
			}
		}
		services = append(services, service{
			ServiceName: svcCfg.ServiceName,
			Handlers:    handlers,
			DONs:        svcCfg.DONs,
		})
	}

	return shardedDONs, services, nil
}

type webAPICapabilitiesHandlerConfig struct {
	MaxAllowedMessageAgeSec int                   `toml:"maxAllowedMessageAgeSec"`
	NodeRateLimiter         nodeRateLimiterConfig `toml:"NodeRateLimiter"`
}

func newDefaultWebAPICapabilitiesHandler() handler {
	return handler{
		Name: GatewayHandlerTypeWebAPICapabilities,
		Config: webAPICapabilitiesHandlerConfig{
			MaxAllowedMessageAgeSec: 1_000,
			NodeRateLimiter: nodeRateLimiterConfig{
				GlobalBurst:    10,
				GlobalRPS:      50,
				PerSenderBurst: 10,
				PerSenderRPS:   10,
			},
		},
	}
}

type vaultHandlerConfig struct {
	RequestTimeoutSec int                   `toml:"requestTimeoutSec"`
	NodeRateLimiter   nodeRateLimiterConfig `toml:"NodeRateLimiter"`
}

func newDefaultVaultHandler(requestTimeoutSec int) handler {
	return handler{
		Name:        "vault",
		ServiceName: "vault",
		Config: vaultHandlerConfig{
			// must be lower than the overall gateway request timeout.
			// so we allow for the response to be sent back.
			RequestTimeoutSec: requestTimeoutSec - 1,
			NodeRateLimiter: nodeRateLimiterConfig{
				GlobalBurst:    10,
				GlobalRPS:      50,
				PerSenderBurst: 10,
				PerSenderRPS:   10,
			},
		},
	}
}

type gatewaySpec struct {
	Type              string        `toml:"type"`
	SchemaVersion     int           `toml:"schemaVersion"`
	Name              string        `toml:"name"`
	ExternalJobID     string        `toml:"externalJobID"`
	ForwardingAllowed bool          `toml:"forwardingAllowed"`
	GatewayConfig     gatewayConfig `toml:"gatewayConfig"`
}

type gatewayConfig struct {
	ConnectionManagerConfig connectionManagerConfig `toml:"ConnectionManagerConfig"`
	ShardedDONs             []shardedDON            `toml:"ShardedDONs"`
	Services                []service               `toml:"Services"`
	HTTPClientConfig        httpClientConfig        `toml:"HTTPClientConfig"`
	NodeServerConfig        nodeServerConfig        `toml:"NodeServerConfig"`
	UserServerConfig        userServerConfig        `toml:"UserServerConfig"`
}

type service struct {
	ServiceName string    `toml:"ServiceName"`
	Handlers    []handler `toml:"Handlers"`
	DONs        []string  `toml:"DONs"`
}

type shardedDON struct {
	DonName string  `toml:"DonName"`
	F       int     `toml:"F"`
	Shards  []shard `toml:"Shards"`
}

type shard struct {
	Nodes []member `toml:"Nodes"`
}

type connectionManagerConfig struct {
	AuthChallengeLen          int    `toml:"AuthChallengeLen"`
	AuthGatewayID             string `toml:"AuthGatewayId"`
	AuthTimestampToleranceSec int    `toml:"AuthTimestampToleranceSec"`
	HeartbeatIntervalSec      int    `toml:"HeartbeatIntervalSec"`
}

type handler struct {
	Name        string `toml:"Name"`
	ServiceName string `toml:"ServiceName,omitempty"`
	Config      any    `toml:"Config"`
}

type member struct {
	Address string `toml:"Address"`
	Name    string `toml:"Name"`
}

type httpClientConfig struct {
	MaxResponseBytes int      `toml:"MaxResponseBytes"`
	AllowedPorts     []int    `toml:"AllowedPorts"`
	AllowedSchemes   []string `toml:"AllowedSchemes"`
	AllowedIPsCIDR   []string `toml:"AllowedIPsCIDR"`
}

type nodeServerConfig struct {
	HandshakeTimeoutMillis int    `toml:"HandshakeTimeoutMillis"`
	MaxRequestBytes        int    `toml:"MaxRequestBytes"`
	Path                   string `toml:"Path"`
	Port                   int    `toml:"Port"`
	ReadTimeoutMillis      int    `toml:"ReadTimeoutMillis"`
	RequestTimeoutMillis   int    `toml:"RequestTimeoutMillis"`
	WriteTimeoutMillis     int    `toml:"WriteTimeoutMillis"`
}

type userServerConfig struct {
	ContentTypeHeader    string `toml:"ContentTypeHeader"`
	MaxRequestBytes      int    `toml:"MaxRequestBytes"`
	Path                 string `toml:"Path"`
	Port                 int    `toml:"Port"`
	ReadTimeoutMillis    int    `toml:"ReadTimeoutMillis"`
	RequestTimeoutMillis int    `toml:"RequestTimeoutMillis"`
	WriteTimeoutMillis   int    `toml:"WriteTimeoutMillis"`
}

type nodeRateLimiterConfig struct {
	GlobalBurst    int `toml:"globalBurst"`
	GlobalRPS      int `toml:"globalRPS"`
	PerSenderBurst int `toml:"perSenderBurst"`
	PerSenderRPS   int `toml:"perSenderRPS"`
}

type httpCapabilitiesHandlerConfig struct {
	NodeRateLimiter nodeRateLimiterConfig `toml:"NodeRateLimiter"`
	CleanUpPeriodMs int                   `toml:"CleanUpPeriodMs"`
}

func newDefaultHTTPCapabilitiesHandler() handler {
	return handler{
		Name:        GatewayHandlerTypeHTTPCapabilities,
		ServiceName: "workflows",
		Config: httpCapabilitiesHandlerConfig{
			NodeRateLimiter: nodeRateLimiterConfig{
				GlobalBurst:    100,
				GlobalRPS:      500,
				PerSenderBurst: 100,
				PerSenderRPS:   100,
			},
			CleanUpPeriodMs: 10 * 60 * 1000, // 10 minutes
		},
	}
}
