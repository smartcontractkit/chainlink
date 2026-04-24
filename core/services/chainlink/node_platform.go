package chainlink

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

const (
	nodePlatformDomain          = "node-platform"
	nodePlatformBuildInfoEntity = "common.v1.NodeBuildInfo"
	nodePlatformJobInfoEntity   = "common.v1.NodeJobInfo"
	nodePlatformDataSchema      = "/node-platform/common/v1"
	nodePlatformBeat            = 3 * time.Minute

	nodeTransmitterTypeOCR                  = "ocr"
	nodeTransmitterTypeOCR2                 = "ocr2"
	nodeTransmitterTypeOCR2DualTransmission = "ocr2_dual_transmission"
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
		beholder.AttrKeyEntity, nodePlatformBuildInfoEntity,
		beholder.AttrKeyDataSchema, nodePlatformDataSchema,
	)
	if err != nil {
		s.eng.Errorw("failed to emit node-platform build info", "err", err)
	}
}

func (s *NodePlatformBuildInfoService) GetBeat() time.Duration {
	return s.beat
}

type NodePlatformJobInfoService struct {
	commonservices.Service
	eng *commonservices.Engine

	opts    NodePlatformJobInfoConfig
	beat    time.Duration
	emitter beholder.Emitter
}

type NodePlatformJobInfoConfig struct {
	Beat         time.Duration
	Lggr         logger.Logger
	CSAKeyStore  keystore.CSA
	JobReader    NodePlatformJobReader
	CSAPublicKey string
}

type NodePlatformJobReader interface {
	FindJobs(ctx context.Context, offset, limit int) ([]job.Job, int, error)
}

func NewNodePlatformJobInfoConfig(opts ApplicationOpts, jobReader NodePlatformJobReader) NodePlatformJobInfoConfig {
	return NodePlatformJobInfoConfig{
		Beat:        nodePlatformBeat,
		Lggr:        opts.Logger,
		CSAKeyStore: opts.KeyStore.CSA(),
		JobReader:   jobReader,
	}
}

func NewNodePlatformJobInfoService(cfg NodePlatformJobInfoConfig) NodePlatformJobInfoService {
	s := NodePlatformJobInfoService{
		opts:    cfg,
		beat:    cfg.Beat,
		emitter: beholder.GetEmitter(),
	}

	s.Service, s.eng = commonservices.Config{
		Name:  "NodePlatformJobInfo",
		Start: s.start,
	}.NewServiceEngine(cfg.Lggr)

	return s
}

func (s *NodePlatformJobInfoService) start(ctx context.Context) error {
	s.resolveCSAPublicKey(ctx)
	s.eng.GoTick(timeutil.NewTicker(s.GetBeat), s.emit)
	return nil
}

func (s *NodePlatformJobInfoService) resolveCSAPublicKey(ctx context.Context) {
	if s.opts.CSAKeyStore == nil {
		return
	}

	csaKey, err := keystore.GetDefault(ctx, s.opts.CSAKeyStore)
	if err != nil {
		s.eng.Errorw("failed to resolve CSA key for node-platform job info", "err", err)
		return
	}

	s.opts.CSAPublicKey = csaKey.PublicKeyString()
}

func (s *NodePlatformJobInfoService) emit(ctx context.Context) {
	payloadBytes, err := proto.Marshal(&commonv1.NodeJobInfo{
		CsaPublicKey: s.opts.CSAPublicKey,
		Transmitters: s.transmitters(ctx),
	})
	if err != nil {
		s.eng.Errorw("failed to marshal node-platform job info", "err", err)
		return
	}

	emitter := s.emitter
	if emitter == nil {
		emitter = beholder.GetEmitter()
	}

	err = emitter.Emit(ctx, payloadBytes,
		beholder.AttrKeyDomain, nodePlatformDomain,
		beholder.AttrKeyEntity, nodePlatformJobInfoEntity,
		beholder.AttrKeyDataSchema, nodePlatformDataSchema,
	)
	if err != nil {
		s.eng.Errorw("failed to emit node-platform job info", "err", err)
	}
}

func (s *NodePlatformJobInfoService) GetBeat() time.Duration {
	return s.beat
}

func (s *NodePlatformJobInfoService) transmitters(ctx context.Context) []*commonv1.NodeTransmitter {
	if s.opts.JobReader == nil {
		return nil
	}

	jobs, _, err := s.opts.JobReader.FindJobs(ctx, 0, math.MaxInt)
	if err != nil {
		s.eng.Warnw("failed to resolve node-platform transmitters", "err", err)
		return nil
	}

	return nodeTransmittersFromJobs(jobs)
}

func nodeTransmittersFromJobs(jobs []job.Job) []*commonv1.NodeTransmitter {
	byChain := make(map[string]map[string]map[string]struct{})
	for _, jb := range jobs {
		if jb.OCROracleSpec != nil {
			addOCRTransmitter(byChain, jb.OCROracleSpec)
		}
		if jb.OCR2OracleSpec != nil {
			addOCR2Transmitters(byChain, jb.OCR2OracleSpec)
		}
	}
	return sortedNodeTransmitters(byChain)
}

func addOCRTransmitter(byChain map[string]map[string]map[string]struct{}, spec *job.OCROracleSpec) {
	if spec == nil || spec.TransmitterAddress == nil || spec.EVMChainID == nil {
		return
	}
	addNodeTransmitter(byChain, spec.EVMChainID.String(), nodeTransmitterTypeOCR, spec.TransmitterAddress.String())
}

func addOCR2Transmitters(byChain map[string]map[string]map[string]struct{}, spec *job.OCR2OracleSpec) {
	if spec == nil {
		return
	}
	chainID := ocr2ChainID(spec)
	if chainID == "" {
		return
	}

	addNullableTransmitterID(byChain, chainID, nodeTransmitterTypeOCR2, spec.TransmitterID)
	if sendingKeys, err := job.SendingKeysForJob(spec); err == nil {
		for _, sendingKey := range sendingKeys {
			addNodeTransmitter(byChain, chainID, nodeTransmitterTypeOCR2, sendingKey)
		}
	}
	addNodeTransmitter(byChain, chainID, nodeTransmitterTypeOCR2DualTransmission, dualTransmissionTransmitterAddress(spec.RelayConfig))
}

func ocr2ChainID(spec *job.OCR2OracleSpec) string {
	if relayID, err := spec.RelayID(); err == nil {
		return strings.TrimSpace(relayID.ChainID)
	}
	if chainID := strings.TrimSpace(spec.ChainID); chainID != "" {
		return chainID
	}
	return jsonConfigString(spec.RelayConfig, "chainID")
}

func addNullableTransmitterID(byChain map[string]map[string]map[string]struct{}, chainID, addressType string, transmitterID null.String) {
	if !transmitterID.Valid {
		return
	}
	addNodeTransmitter(byChain, chainID, addressType, transmitterID.String)
}

func dualTransmissionTransmitterAddress(config job.JSONConfig) string {
	if !jsonConfigBool(config, "enableDualTransmission") {
		return ""
	}
	raw, ok := config["dualTransmission"]
	if !ok {
		return ""
	}
	value, ok := raw.(map[string]any)
	if !ok {
		return ""
	}
	return jsonConfigString(value, "transmitterAddress")
}

func jsonConfigBool(config map[string]any, key string) bool {
	raw, ok := config[key]
	if !ok {
		return false
	}
	switch value := raw.(type) {
	case bool:
		return value
	case string:
		return strings.EqualFold(strings.TrimSpace(value), "true")
	default:
		return false
	}
}

func jsonConfigString(config map[string]any, key string) string {
	raw, ok := config[key]
	if !ok {
		return ""
	}
	switch value := raw.(type) {
	case string:
		return strings.TrimSpace(value)
	case fmt.Stringer:
		return strings.TrimSpace(value.String())
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return strings.TrimSpace(fmt.Sprint(value))
	default:
		return ""
	}
}

func addNodeTransmitter(byChain map[string]map[string]map[string]struct{}, chainID, addressType, address string) {
	chainID = strings.TrimSpace(chainID)
	addressType = strings.TrimSpace(addressType)
	address = strings.TrimSpace(address)
	if chainID == "" || addressType == "" || address == "" {
		return
	}
	if byChain[chainID] == nil {
		byChain[chainID] = make(map[string]map[string]struct{})
	}
	if byChain[chainID][addressType] == nil {
		byChain[chainID][addressType] = make(map[string]struct{})
	}
	byChain[chainID][addressType][address] = struct{}{}
}

func sortedNodeTransmitters(byChain map[string]map[string]map[string]struct{}) []*commonv1.NodeTransmitter {
	if len(byChain) == 0 {
		return nil
	}
	chainIDs := make([]string, 0, len(byChain))
	for chainID := range byChain {
		chainIDs = append(chainIDs, chainID)
	}
	sort.Strings(chainIDs)

	transmitters := make([]*commonv1.NodeTransmitter, 0, len(chainIDs))
	for _, chainID := range chainIDs {
		addressesByType := make(map[string]*commonv1.NodeTransmitterAddresses, len(byChain[chainID]))
		for addressType, addressesSet := range byChain[chainID] {
			addresses := make([]string, 0, len(addressesSet))
			for address := range addressesSet {
				addresses = append(addresses, address)
			}
			sort.Strings(addresses)
			if len(addresses) == 0 {
				continue
			}
			addressesByType[addressType] = &commonv1.NodeTransmitterAddresses{Values: addresses}
		}
		if len(addressesByType) == 0 {
			continue
		}
		transmitters = append(transmitters, &commonv1.NodeTransmitter{
			ChainId:   chainID,
			Addresses: addressesByType,
		})
	}
	return transmitters
}
