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
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

const (
	nodePlatformDomain     = "node-platform"
	nodePlatformEntity     = "common.v1.NodeBuildInfo"
	nodePlatformDataSchema = "/node-platform/common/v1"
)

type NodePlatformBuildInfo struct {
	commonservices.Service
	eng *commonservices.Engine

	opts    NodePlatformBuildInfoConfig
	beat    time.Duration
	emitter beholder.Emitter
}

type NodePlatformBuildInfoConfig struct {
	Beat         time.Duration
	Lggr         logger.Logger
	CSAPublicKey string
	CommitSHA    string
	DockerTag    string
	VersionTag   string
	Version      string
}

func NewNodePlatformBuildInfoConfig(opts ApplicationOpts) NodePlatformBuildInfoConfig {
	csaKey := ""
	csaKeys, err := opts.KeyStore.CSA().GetAll()
	if err != nil {
		opts.Logger.Errorw("failed to get CSA keys for node-platform build info", "err", err)
	}

	if len(csaKeys) > 0 {
		csaKey = csaKeys[0].PublicKeyString()
	} else {
		opts.Logger.Warn("no CSA key found for node-platform build info")
	}

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
		dockerTag = static.DockerTag
	}

	return NodePlatformBuildInfoConfig{
		Beat:         opts.Config.Telemetry().HeartbeatInterval(),
		Lggr:         opts.Logger,
		CSAPublicKey: csaKey,
		CommitSHA:    static.Sha,
		DockerTag:    dockerTag,
		VersionTag:   versionTag,
		Version:      version,
	}
}

type NodePlatformBuildInfoOpt func(*NodePlatformBuildInfo)

func WithNodePlatformBuildInfoEmitter(emitter beholder.Emitter) NodePlatformBuildInfoOpt {
	return func(h *NodePlatformBuildInfo) {
		h.emitter = emitter
	}
}

func NewNodePlatformBuildInfo(cfg NodePlatformBuildInfoConfig, opts ...NodePlatformBuildInfoOpt) NodePlatformBuildInfo {
	h := NodePlatformBuildInfo{
		opts:    cfg,
		beat:    cfg.Beat,
		emitter: beholder.GetEmitter(),
	}

	for _, opt := range opts {
		opt(&h)
	}

	h.Service, h.eng = commonservices.Config{
		Name:  "NodePlatformBuildInfo",
		Start: h.start,
	}.NewServiceEngine(cfg.Lggr)

	return h
}

func (h *NodePlatformBuildInfo) start(_ context.Context) error {
	h.eng.GoTick(timeutil.NewTicker(h.GetBeat), h.emit)
	return nil
}

func (h *NodePlatformBuildInfo) emit(ctx context.Context) {
	payloadBytes, err := proto.Marshal(&commonv1.NodeBuildInfo{
		CsaPublicKey: h.opts.CSAPublicKey,
		CommitSha:    h.opts.CommitSHA,
		DockerTag:    h.opts.DockerTag,
		VersionTag:   h.opts.VersionTag,
		Version:      h.opts.Version,
	})
	if err != nil {
		h.eng.Errorw("failed to marshal node-platform build info", "err", err)
		return
	}

	emitter := h.emitter
	if emitter == nil {
		emitter = beholder.GetEmitter()
	}

	err = emitter.Emit(ctx, payloadBytes,
		beholder.AttrKeyDomain, nodePlatformDomain,
		beholder.AttrKeyEntity, nodePlatformEntity,
		beholder.AttrKeyDataSchema, nodePlatformDataSchema,
	)
	if err != nil {
		h.eng.Errorw("failed to emit node-platform build info", "err", err)
	}
}

func (h *NodePlatformBuildInfo) GetBeat() time.Duration {
	return h.beat
}
