package chainlink

import (
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
	nodePlatformDomain     = "node-platform"
	nodePlatformEntity     = "common.v1.NodeBuildInfo"
	nodePlatformDataSchema = "/node-platform/common/v1"
	nodePlatformBeat       = 3 * time.Minute
)

type NodePlatformBuildInfoService struct {
	commonservices.Service
	eng *commonservices.Engine

	opts    NodePlatformBuildInfoConfig
	beat    time.Duration
	emitter beholder.Emitter
}

type NodePlatformBuildInfoConfig struct {
	Beat         time.Duration
	Lggr         logger.Logger
	CSAKeyStore  keystore.CSA
	CSAPublicKey string
	CommitSHA    string
	DockerTag    string
	VersionTag   string
	Version      string
}

func NewNodePlatformBuildInfoConfig(opts ApplicationOpts) NodePlatformBuildInfoConfig {
	version := opts.Version
	if version == "" {
		version = static.Version
	}

	versionTag := opts.VersionTag
	if versionTag == "" {
		versionTag = static.VersionTag
	}

	dockerTag := opts.DockerTag
	if dockerTag == "" {
		dockerTag = static.Unset
	}

	return NodePlatformBuildInfoConfig{
		Beat:        nodePlatformBeat,
		Lggr:        opts.Logger,
		CSAKeyStore: opts.KeyStore.CSA(),
		CommitSHA:   static.Sha,
		DockerTag:   dockerTag,
		VersionTag:  versionTag,
		Version:     version,
	}
}

func NewNodePlatformBuildInfoService(cfg NodePlatformBuildInfoConfig) NodePlatformBuildInfoService {
	s := NodePlatformBuildInfoService{
		opts:    cfg,
		beat:    cfg.Beat,
		emitter: beholder.GetEmitter(),
	}

	s.Service, s.eng = commonservices.Config{
		Name:  "NodePlatformBuildInfo",
		Start: s.start,
	}.NewServiceEngine(cfg.Lggr)

	return s
}

func (s *NodePlatformBuildInfoService) start(ctx context.Context) error {
	s.resolveCSAPublicKey(ctx)
	s.eng.GoTick(timeutil.NewTicker(s.GetBeat), s.emit)
	return nil
}

func (s *NodePlatformBuildInfoService) resolveCSAPublicKey(ctx context.Context) {
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

func (s *NodePlatformBuildInfoService) emit(ctx context.Context) {
	payloadBytes, err := proto.Marshal(&commonv1.NodeBuildInfo{
		CsaPublicKey: s.opts.CSAPublicKey,
		CommitSha:    s.opts.CommitSHA,
		DockerTag:    s.opts.DockerTag,
		VersionTag:   s.opts.VersionTag,
		Version:      s.opts.Version,
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
		beholder.AttrKeyEntity, nodePlatformEntity,
		beholder.AttrKeyDataSchema, nodePlatformDataSchema,
	)
	if err != nil {
		s.eng.Errorw("failed to emit node-platform build info", "err", err)
	}
}

func (s *NodePlatformBuildInfoService) GetBeat() time.Duration {
	return s.beat
}
