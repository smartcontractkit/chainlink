package pkg

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
)

const (
	GatewayHandlerTypeWebAPICapabilities = "web-api-capabilities"
	GatewayHandlerTypeHTTPCapabilities   = "http-capabilities"
	GatewayHandlerTypeVault              = "vault"
)

type TargetDONMember struct {
	Address string
	Name    string
}

type TargetDON struct {
	ID       string
	Members  []TargetDONMember
	Handlers []string
}

type GatewayJob struct {
	TargetDONs    []TargetDON
	JobName       string
	ExternalJobID string
}

func (g GatewayJob) Validate() error {
	if g.JobName == "" {
		return errors.New("must provide job name")
	}

	if len(g.TargetDONs) == 0 {
		return errors.New("must provide at least one target DON")
	}

	return nil
}

func (g GatewayJob) Resolve(gatewayNodeIdx int) (string, error) {
	externalJobID := g.ExternalJobID
	if externalJobID == "" {
		externalJobID = uuid.NewSHA1(uuid.Nil, []byte(g.JobName)).String()
	}

	dons := []don{}
	for _, targetDON := range g.TargetDONs {
		ms := []member{}
		for _, mem := range targetDON.Members {
			ms = append(ms, member(mem))
		}

		hs := []handler{}
		for _, ht := range targetDON.Handlers {
			switch ht {
			case GatewayHandlerTypeWebAPICapabilities:
				hs = append(hs, newDefaultWebAPICapabilitiesHandler())
			case GatewayHandlerTypeVault:
				hs = append(hs, newDefaultVaultHandler())
			case GatewayHandlerTypeHTTPCapabilities:
				hs = append(hs, newDefaultHTTPCapabilitiesHandler())
			default:
				return "", errors.New("unknown handler type: " + ht)
			}
		}

		d := don{
			DonID:    targetDON.ID,
			Members:  ms,
			Handlers: hs,
		}
		dons = append(dons, d)
	}

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
			RequestTimeoutMillis:   10_000,
			WriteTimeoutMillis:     1_000,
		},
		UserServerConfig: userServerConfig{
			ContentTypeHeader:    "application/jsonrpc",
			MaxRequestBytes:      100_000,
			Path:                 "/",
			Port:                 5_002,
			ReadTimeoutMillis:    80_000,
			RequestTimeoutMillis: 80_000,
			WriteTimeoutMillis:   80_000,
		},
		HTTPClientConfig: httpClientConfig{
			MaxResponseBytes: 50_000_000,
		},
		Dons: dons,
	}

	spec := &gatewaySpec{
		Type:              "gateway",
		SchemaVersion:     1,
		Name:              g.JobName,
		ExternalJobID:     externalJobID,
		ForwardingAllowed: false,
		GatewayConfig:     config,
	}
	b, err := toml.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(b), nil
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

func newDefaultVaultHandler() handler {
	return handler{
		Name:        "vault",
		ServiceName: "vault",
		Config: vaultHandlerConfig{
			RequestTimeoutSec: 70,
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
	Dons                    []don                   `toml:"Dons"`
	HTTPClientConfig        httpClientConfig        `toml:"HTTPClientConfig"`
	NodeServerConfig        nodeServerConfig        `toml:"NodeServerConfig"`
	UserServerConfig        userServerConfig        `toml:"UserServerConfig"`
}

type connectionManagerConfig struct {
	AuthChallengeLen          int    `toml:"AuthChallengeLen"`
	AuthGatewayID             string `toml:"AuthGatewayId"`
	AuthTimestampToleranceSec int    `toml:"AuthTimestampToleranceSec"`
	HeartbeatIntervalSec      int    `toml:"HeartbeatIntervalSec"`
}

type don struct {
	DonID    string    `toml:"DonId"`
	Handlers []handler `toml:"Handlers"`
	Members  []member  `toml:"Members"`
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
	MaxResponseBytes int `toml:"MaxResponseBytes"`
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
				GlobalBurst:    10,
				GlobalRPS:      50,
				PerSenderBurst: 10,
				PerSenderRPS:   10,
			},
			CleanUpPeriodMs: 86400000, // 24 hours
		},
	}
}
