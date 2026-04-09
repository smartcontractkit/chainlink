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
	GatewayHandlerTypeConfidentialRelay  = "confidential-compute-relay"

	ServiceNameWorkflows    = "workflows"
	ServiceNameVault        = "vault"
	ServiceNameConfidential = "confidential"

	minimumRequestTimeoutSec = 5
)

// HandlerServiceName returns the service name for a given handler type.
func HandlerServiceName(handlerType string) string {
	switch handlerType {
	case GatewayHandlerTypeVault:
		return ServiceNameVault
	case GatewayHandlerTypeHTTPCapabilities, GatewayHandlerTypeWebAPICapabilities:
		return ServiceNameWorkflows
	case GatewayHandlerTypeConfidentialRelay:
		return ServiceNameConfidential
	default:
		return handlerType
	}
}

type TargetDONMember struct {
	Address string
	Name    string
}

type TargetDON struct {
	ID       string
	F        int
	Members  []TargetDONMember
	Handlers []string // used only in legacy (don-centric) format
}

type GatewayServiceConfig struct {
	ServiceName string
	Handlers    []string
	DONs        []string
}

type GatewayJob struct {
	ServiceCentricFormatEnabled bool

	// Don-centric format (ServiceCentricFormatEnabled == false): handlers are per-DON.
	TargetDONs []TargetDON

	// Service-centric format (ServiceCentricFormatEnabled == true): handlers are per-service, DONs referenced by name.
	DONs     []TargetDON
	Services []GatewayServiceConfig

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

	if g.ServiceCentricFormatEnabled {
		if len(g.DONs) == 0 {
			return errors.New("must provide at least one DON")
		}
		if len(g.Services) == 0 {
			return errors.New("must provide at least one service")
		}
	} else if len(g.TargetDONs) == 0 {
		return errors.New("must provide at least one target DON")
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

	requestTimeout := time.Duration(g.RequestTimeoutSec) * time.Second
	connCfg := connectionManagerConfig{
		AuthChallengeLen:          10,
		AuthGatewayID:             "gateway-node-" + strconv.Itoa(gatewayNodeIdx),
		AuthTimestampToleranceSec: 5,
		HeartbeatIntervalSec:      20,
	}
	nodeCfg := nodeServerConfig{
		HandshakeTimeoutMillis: 1_000,
		MaxRequestBytes:        100_000,
		Path:                   "/",
		Port:                   5_003,
		ReadTimeoutMillis:      1_000,
		RequestTimeoutMillis:   int(requestTimeout.Milliseconds()),
		WriteTimeoutMillis:     1_000,
	}
	userCfg := userServerConfig{
		ContentTypeHeader:    "application/jsonrpc",
		MaxRequestBytes:      100_000,
		Path:                 "/",
		Port:                 5_002,
		ReadTimeoutMillis:    int(requestTimeout.Milliseconds()),
		RequestTimeoutMillis: int(requestTimeout.Milliseconds()),
		WriteTimeoutMillis:   int(requestTimeout.Milliseconds() + 1000),
	}
	httpCfg := httpClientConfig{
		MaxResponseBytes: 50_000_000,
		AllowedPorts:     []int{443},
		AllowedSchemes:   []string{"https"},
	}

	if len(g.AllowedPorts) > 0 {
		httpCfg.AllowedPorts = g.AllowedPorts
	}
	if len(g.AllowedSchemes) > 0 {
		httpCfg.AllowedSchemes = g.AllowedSchemes
	}
	if len(g.AllowedIPsCIDR) > 0 {
		httpCfg.AllowedIPsCIDR = g.AllowedIPsCIDR
	}
	if g.AuthGatewayID != "" {
		connCfg.AuthGatewayID = g.AuthGatewayID
	}

	var gwConfig any
	if g.ServiceCentricFormatEnabled {
		shardedDONs, services, err := g.buildServicesAndShardedDONs()
		if err != nil {
			return "", err
		}
		gwConfig = gatewayConfigServiceCentric{
			ConnectionManagerConfig: connCfg,
			ShardedDONs:             shardedDONs,
			Services:                services,
			HTTPClientConfig:        httpCfg,
			NodeServerConfig:        nodeCfg,
			UserServerConfig:        userCfg,
		}
	} else {
		dons, err := g.buildLegacyDons()
		if err != nil {
			return "", err
		}
		gwConfig = gatewayConfigLegacy{
			ConnectionManagerConfig: connCfg,
			Dons:                    dons,
			HTTPClientConfig:        httpCfg,
			NodeServerConfig:        nodeCfg,
			UserServerConfig:        userCfg,
		}
	}

	spec := &gatewaySpec{
		Type:              "gateway",
		SchemaVersion:     1,
		Name:              g.JobName,
		ExternalJobID:     externalJobID,
		ForwardingAllowed: false,
		GatewayConfig:     gwConfig,
	}
	b, marshalErr := toml.Marshal(spec)
	if marshalErr != nil {
		return "", marshalErr
	}

	return string(b), nil
}

func (g GatewayJob) buildLegacyDons() ([]legacyDON, error) {
	dons := make([]legacyDON, 0, len(g.TargetDONs))
	for _, targetDON := range g.TargetDONs {
		ms := make([]member, 0, len(targetDON.Members))
		for _, mem := range targetDON.Members {
			ms = append(ms, member(mem))
		}

		hs := make([]handler, 0, len(targetDON.Handlers))
		for _, ht := range targetDON.Handlers {
			switch ht {
			case GatewayHandlerTypeWebAPICapabilities:
				hs = append(hs, newDefaultWebAPICapabilitiesHandler())
			case GatewayHandlerTypeVault:
				hs = append(hs, newDefaultVaultHandler(g.RequestTimeoutSec))
			case GatewayHandlerTypeHTTPCapabilities:
				hs = append(hs, newDefaultHTTPCapabilitiesHandler())
			case GatewayHandlerTypeConfidentialRelay:
				// -1 so the handler times out before the gateway, allowing a clean error response.
				// TODO: the vault handler does the same -1 internally; unify both to use this pattern.
				hs = append(hs, newDefaultConfidentialRelayHandler(g.RequestTimeoutSec-1))
			default:
				return nil, errors.New("unknown handler type: " + ht)
			}
		}

		dons = append(dons, legacyDON{
			DonID:    targetDON.ID,
			F:        targetDON.F,
			Handlers: hs,
			Members:  ms,
		})
	}
	return dons, nil
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
			case GatewayHandlerTypeConfidentialRelay:
				// -1 so the handler times out before the gateway, allowing a clean error response.
				handlers = append(handlers, newDefaultConfidentialRelayHandler(g.RequestTimeoutSec-1))
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
		Name:        GatewayHandlerTypeVault,
		ServiceName: ServiceNameVault,
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
	Type              string `toml:"type"`
	SchemaVersion     int    `toml:"schemaVersion"`
	Name              string `toml:"name"`
	ExternalJobID     string `toml:"externalJobID"`
	ForwardingAllowed bool   `toml:"forwardingAllowed"`
	GatewayConfig     any    `toml:"gatewayConfig"`
}

type gatewayConfigLegacy struct {
	ConnectionManagerConfig connectionManagerConfig `toml:"ConnectionManagerConfig"`
	Dons                    []legacyDON             `toml:"Dons"`
	HTTPClientConfig        httpClientConfig        `toml:"HTTPClientConfig"`
	NodeServerConfig        nodeServerConfig        `toml:"NodeServerConfig"`
	UserServerConfig        userServerConfig        `toml:"UserServerConfig"`
}

type gatewayConfigServiceCentric struct {
	ConnectionManagerConfig connectionManagerConfig `toml:"ConnectionManagerConfig"`
	ShardedDONs             []shardedDON            `toml:"ShardedDONs"`
	Services                []service               `toml:"Services"`
	HTTPClientConfig        httpClientConfig        `toml:"HTTPClientConfig"`
	NodeServerConfig        nodeServerConfig        `toml:"NodeServerConfig"`
	UserServerConfig        userServerConfig        `toml:"UserServerConfig"`
}

type legacyDON struct {
	DonID    string    `toml:"DonId"`
	F        int       `toml:"F"`
	Handlers []handler `toml:"Handlers"`
	Members  []member  `toml:"Members"`
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
	CleanUpPeriodMs int                   `toml:"CleanUpPeriodMs"`
	NodeRateLimiter nodeRateLimiterConfig `toml:"NodeRateLimiter"`
}

func newDefaultHTTPCapabilitiesHandler() handler {
	return handler{
		Name:        GatewayHandlerTypeHTTPCapabilities,
		ServiceName: ServiceNameWorkflows,
		Config: httpCapabilitiesHandlerConfig{
			CleanUpPeriodMs: 10 * 60 * 1000, // 10 minutes
			NodeRateLimiter: nodeRateLimiterConfig{
				GlobalBurst:    100,
				GlobalRPS:      500,
				PerSenderBurst: 100,
				PerSenderRPS:   100,
			},
		},
	}
}

type confidentialRelayHandlerConfig struct {
	RequestTimeoutSec int `toml:"requestTimeoutSec"`
}

func newDefaultConfidentialRelayHandler(requestTimeoutSec int) handler {
	return handler{
		Name:        GatewayHandlerTypeConfidentialRelay,
		ServiceName: ServiceNameConfidential,
		Config: confidentialRelayHandlerConfig{
			RequestTimeoutSec: requestTimeoutSec,
		},
	}
}
