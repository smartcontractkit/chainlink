package chainlink

import (
	"cmp"
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

const (
	nodePlatformDomain          = "node-platform"
	nodePlatformBuildInfoEntity = "common.v1.NodeInfo"

	nodePlatformDataSchema = "/node-platform/common/v1"
	nodePlatformBeat       = 3 * time.Minute
)

// HealthChecker is an interface for github.com/smartcontractkit/chainlink-common@/pkg/services/health.go
// it registers all services in chainlink/core/services/chainlink/application.go
// and then we call IsHealthy to receive system-wide health-check and errors for each service
type HealthChecker interface {
	IsHealthy() (bool, map[string]error)
}

// NodePlatformInfoService is a basic info service for node
// provides node build info and a health-check/errors
type NodePlatformInfoService struct {
	commonservices.Service
	eng     *commonservices.Engine
	beat    time.Duration
	emitter beholder.Emitter
	hc      HealthChecker

	opts NodePlatformInfoConfig
}

type NodePlatformInfoConfig struct {
	Beat time.Duration
	Lggr logger.Logger
	// Key info
	CSAKeyStore  keystore.CSA
	CSAPublicKey string
	// Build info
	CommitSHA  string
	DockerTag  string
	VersionTag string
	Version    string
}

func NewNodePlatformInfoConfig(opts ApplicationOpts) NodePlatformInfoConfig {
	return NodePlatformInfoConfig{
		Beat:        nodePlatformBeat,
		Lggr:        opts.Logger,
		CSAKeyStore: opts.KeyStore.CSA(),
		CommitSHA:   static.Sha,
		DockerTag:   cmp.Or(opts.DockerTag, string(static.Unset)),
		VersionTag:  cmp.Or(opts.VersionTag, static.VersionTag),
		Version:     cmp.Or(opts.Version, static.Version),
	}
}

func NewNodePlatformInfoService(hc HealthChecker, cfg NodePlatformInfoConfig) NodePlatformInfoService {
	s := NodePlatformInfoService{
		opts:    cfg,
		beat:    cfg.Beat,
		emitter: beholder.GetEmitter(),
		hc:      hc,
	}

	s.Service, s.eng = commonservices.Config{
		Name:  "NodePlatformInfo",
		Start: s.start,
	}.NewServiceEngine(cfg.Lggr)

	return s
}

func (s *NodePlatformInfoService) start(ctx context.Context) error {
	s.resolveCSAPublicKey(ctx)
	s.eng.GoTick(timeutil.NewTicker(s.GetBeat), s.emit)
	return nil
}

func (s *NodePlatformInfoService) resolveCSAPublicKey(ctx context.Context) {
	if s.opts.CSAKeyStore == nil {
		return
	}

	csaKey, err := keystore.GetDefault(ctx, s.opts.CSAKeyStore)
	if err != nil {
		s.eng.Errorw("failed to resolve CSA key for node-platform build info", "err", err)
		return
	}

	s.opts.CSAPublicKey = csaKey.PublicKeyString()
}

func (s *NodePlatformInfoService) emit(ctx context.Context) {
	// report current health check state and errors
	var (
		allServicesHealthy    bool
		healthCheckErrs       map[string]error
		healthCheckErrStrings = make(map[string]string)
	)
	if s.hc != nil {
		allServicesHealthy, healthCheckErrs = s.hc.IsHealthy()
		if len(healthCheckErrs) > 0 {
			for k, v := range healthCheckErrs {
				healthCheckErrStrings[k] = v.Error()
			}
		}
	}

	payloadBytes, err := proto.Marshal(&commonv1.NodeInfo{
		CsaPublicKey: s.opts.CSAPublicKey,
		CommitSha:    s.opts.CommitSHA,
		DockerTag:    s.opts.DockerTag,
		VersionTag:   s.opts.VersionTag,
		Version:      s.opts.Version,
		Healthy:      allServicesHealthy,
		HealthErrors: healthCheckErrStrings,
	})
	if err != nil {
		s.eng.Errorw("failed to marshal node-platform build info", "err", err)
		return
	}

	emitter := s.emitter
	if emitter == nil {
		emitter = beholder.GetEmitter()
	}

	err = emitter.Emit(ctx, payloadBytes,
		beholder.AttrKeyDomain, nodePlatformDomain,
		beholder.AttrKeyEntity, nodePlatformBuildInfoEntity,
		beholder.AttrKeyDataSchema, nodePlatformDataSchema,
	)
	if err != nil {
		s.eng.Errorw("failed to emit node-platform build info", "err", err)
	}
}

func (s *NodePlatformInfoService) GetBeat() time.Duration {
	return s.beat
}
