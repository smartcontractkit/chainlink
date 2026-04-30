package chainlink

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

const (
	nodePlatformDomain          = "node-platform"
	nodePlatformBuildInfoEntity = "common.v1.NodeBuildInfo"
	nodePlatformJobInfoEntity   = "common.v1.NodeJobInfo"
	nodePlatformDataSchema      = "/node-platform/common/v1"
	nodePlatformBeat            = 3 * time.Minute
	nodePlatformJobInfoPageSize = 1000

	nodeSubmitterFieldTransmitterAddress                 = "transmitterAddress"
	nodeSubmitterFieldTransmitterID                      = "transmitterID"
	nodeSubmitterFieldRelayConfigSendingKeys             = "relayConfig.sendingKeys"
	nodeSubmitterFieldDualTransmissionTransmitterAddress = "relayConfig.dualTransmission.transmitterAddress"
	nodeSubmitterFieldFromAddresses                      = "fromAddresses"
	nodeSubmitterFieldOracleFactoryTransmitterID         = "oracle_factory.transmitter_id"
	nodeSubmitterFieldObservationSourceETHTxFrom         = "observationSource.ethtx.from"
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
		CsaPublicKey:       s.opts.CSAPublicKey,
		SubmitterAddresses: s.submitterAddresses(ctx),
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

func (s *NodePlatformJobInfoService) submitterAddresses(ctx context.Context) []*commonv1.NodeSubmitterAddress {
	if s.opts.JobReader == nil {
		return nil
	}

	bySource := make(map[nodeSubmitterAddressKey]map[string]struct{})
	for offset := 0; ; {
		jobs, count, err := s.opts.JobReader.FindJobs(ctx, offset, nodePlatformJobInfoPageSize)
		if err != nil {
			s.eng.Warnw("failed to resolve node-platform submitter addresses", "offset", offset, "limit", nodePlatformJobInfoPageSize, "err", err)
			return nil
		}

		addNodeSubmitterAddressesFromJobs(bySource, jobs)

		offset += len(jobs)
		if len(jobs) == 0 || offset >= count || len(jobs) < nodePlatformJobInfoPageSize {
			break
		}
	}

	return sortedNodeSubmitterAddresses(bySource)
}

type nodeSubmitterAddressKey struct {
	chainID    string
	jobType    string
	pluginType string
	fieldPath  string
}

func addNodeSubmitterAddressesFromJobs(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jobs []job.Job) {
	for _, jb := range jobs {
		addOCRSubmitterAddress(bySource, jb)
		addOCR2SubmitterAddresses(bySource, jb)
		addVRFSubmitterAddresses(bySource, jb)
		addBlockhashStoreSubmitterAddresses(bySource, jb)
		addBlockHeaderFeederSubmitterAddresses(bySource, jb)
		addLegacyGasStationServerSubmitterAddresses(bySource, jb)
		addStandardCapabilitiesSubmitterAddress(bySource, jb)
		addPipelineETHTxSubmitterAddresses(bySource, jb)
	}
}

func addOCRSubmitterAddress(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.OCROracleSpec
	if spec == nil || spec.TransmitterAddress == nil || spec.EVMChainID == nil {
		return
	}
	addNodeSubmitterAddress(bySource, spec.EVMChainID.String(), jobType(jb, job.OffchainReporting), "", nodeSubmitterFieldTransmitterAddress, spec.TransmitterAddress.String())
}

func addOCR2SubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.OCR2OracleSpec
	if spec == nil || !isOnChainOCR2Plugin(spec.PluginType) {
		return
	}
	chainID := ocr2ChainID(spec)
	if chainID == "" {
		return
	}

	pluginType := string(spec.PluginType)
	if spec.TransmitterID.Valid {
		addNodeSubmitterAddress(bySource, chainID, jobType(jb, job.OffchainReporting2), pluginType, nodeSubmitterFieldTransmitterID, spec.TransmitterID.String)
	}
	if sendingKeys, err := job.SendingKeysForJob(spec); err == nil {
		addNodeSubmitterAddress(bySource, chainID, jobType(jb, job.OffchainReporting2), pluginType, nodeSubmitterFieldRelayConfigSendingKeys, sendingKeys...)
	}
	addNodeSubmitterAddress(bySource, chainID, jobType(jb, job.OffchainReporting2), pluginType, nodeSubmitterFieldDualTransmissionTransmitterAddress, dualTransmissionTransmitterAddress(spec.RelayConfig))
}

func isOnChainOCR2Plugin(pluginType commontypes.OCR2PluginType) bool {
	switch pluginType {
	case commontypes.Mercury, commontypes.LLO:
		return false
	default:
		return true
	}
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

func addVRFSubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.VRFSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	addNodeSubmitterAddress(bySource, spec.EVMChainID.String(), jobType(jb, job.VRF), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func addBlockhashStoreSubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.BlockhashStoreSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	addNodeSubmitterAddress(bySource, spec.EVMChainID.String(), jobType(jb, job.BlockhashStore), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func addBlockHeaderFeederSubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.BlockHeaderFeederSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	addNodeSubmitterAddress(bySource, spec.EVMChainID.String(), jobType(jb, job.BlockHeaderFeeder), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func addLegacyGasStationServerSubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.LegacyGasStationServerSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	addNodeSubmitterAddress(bySource, spec.EVMChainID.String(), jobType(jb, job.LegacyGasStationServer), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func addStandardCapabilitiesSubmitterAddress(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	spec := jb.StandardCapabilitiesSpec
	if spec == nil || !spec.OracleFactory.Enabled {
		return
	}
	addNodeSubmitterAddress(bySource, spec.OracleFactory.ChainID, jobType(jb, job.StandardCapabilities), "", nodeSubmitterFieldOracleFactoryTransmitterID, spec.OracleFactory.TransmitterID)
}

func addPipelineETHTxSubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}, jb job.Job) {
	p, ok := jobPipeline(jb)
	if !ok {
		return
	}
	for _, task := range p.Tasks {
		ethTxTask, ok := task.(*pipeline.ETHTxTask)
		if !ok {
			continue
		}
		addresses := staticPipelineAddresses(ethTxTask.From)
		if len(addresses) == 0 {
			continue
		}
		chainID := staticPipelineString(ethTxTask.EVMChainID)
		if chainID == "" {
			chainID = jobEVMChainID(jb)
		}
		addNodeSubmitterAddress(bySource, chainID, jobType(jb, ""), "", nodeSubmitterFieldObservationSourceETHTxFrom, addresses...)
	}
}

func jobPipeline(jb job.Job) (*pipeline.Pipeline, bool) {
	if len(jb.Pipeline.Tasks) > 0 {
		return &jb.Pipeline, true
	}
	if jb.PipelineSpec == nil {
		return nil, false
	}
	p, err := jb.PipelineSpec.GetOrParsePipeline()
	if err != nil {
		return nil, false
	}
	return p, true
}

func staticPipelineAddresses(raw string) []string {
	raw = staticPipelineString(raw)
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "[") {
		var values []string
		if err := json.Unmarshal([]byte(raw), &values); err == nil {
			return values
		}
		var anyValues []any
		if err := json.Unmarshal([]byte(raw), &anyValues); err != nil {
			return nil
		}
		addresses := make([]string, 0, len(anyValues))
		for _, value := range anyValues {
			if address := strings.TrimSpace(fmt.Sprint(value)); address != "" {
				addresses = append(addresses, address)
			}
		}
		return addresses
	}
	return []string{raw}
}

func staticPipelineString(raw string) string {
	value := strings.Trim(strings.TrimSpace(raw), `'"`)
	if value == "" || strings.Contains(value, "$(") {
		return ""
	}
	return value
}

func jobEVMChainID(jb job.Job) string {
	switch {
	case jb.DirectRequestSpec != nil && jb.DirectRequestSpec.EVMChainID != nil:
		return jb.DirectRequestSpec.EVMChainID.String()
	case jb.FluxMonitorSpec != nil && jb.FluxMonitorSpec.EVMChainID != nil:
		return jb.FluxMonitorSpec.EVMChainID.String()
	case jb.OCROracleSpec != nil && jb.OCROracleSpec.EVMChainID != nil:
		return jb.OCROracleSpec.EVMChainID.String()
	case jb.VRFSpec != nil && jb.VRFSpec.EVMChainID != nil:
		return jb.VRFSpec.EVMChainID.String()
	case jb.BlockhashStoreSpec != nil && jb.BlockhashStoreSpec.EVMChainID != nil:
		return jb.BlockhashStoreSpec.EVMChainID.String()
	case jb.BlockHeaderFeederSpec != nil && jb.BlockHeaderFeederSpec.EVMChainID != nil:
		return jb.BlockHeaderFeederSpec.EVMChainID.String()
	case jb.LegacyGasStationServerSpec != nil && jb.LegacyGasStationServerSpec.EVMChainID != nil:
		return jb.LegacyGasStationServerSpec.EVMChainID.String()
	case jb.LegacyGasStationSidecarSpec != nil && jb.LegacyGasStationSidecarSpec.EVMChainID != nil:
		return jb.LegacyGasStationSidecarSpec.EVMChainID.String()
	default:
		return ""
	}
}

func eip55AddressStrings[T fmt.Stringer](addresses []T) []string {
	if len(addresses) == 0 {
		return nil
	}
	out := make([]string, 0, len(addresses))
	for _, address := range addresses {
		out = append(out, address.String())
	}
	return out
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

func jobType(jb job.Job, fallback job.Type) string {
	if jb.Type != "" {
		return jb.Type.String()
	}
	return fallback.String()
}

func addNodeSubmitterAddress(bySource map[nodeSubmitterAddressKey]map[string]struct{}, chainID, jobType, pluginType, fieldPath string, addresses ...string) {
	chainID = strings.TrimSpace(chainID)
	jobType = strings.TrimSpace(jobType)
	pluginType = strings.TrimSpace(pluginType)
	fieldPath = strings.TrimSpace(fieldPath)
	if chainID == "" || jobType == "" || fieldPath == "" {
		return
	}
	key := nodeSubmitterAddressKey{
		chainID:    chainID,
		jobType:    jobType,
		pluginType: pluginType,
		fieldPath:  fieldPath,
	}
	for _, address := range addresses {
		address = strings.TrimSpace(address)
		if address == "" {
			continue
		}
		if bySource[key] == nil {
			bySource[key] = make(map[string]struct{})
		}
		bySource[key][address] = struct{}{}
	}
}

func sortedNodeSubmitterAddresses(bySource map[nodeSubmitterAddressKey]map[string]struct{}) []*commonv1.NodeSubmitterAddress {
	if len(bySource) == 0 {
		return nil
	}
	keys := make([]nodeSubmitterAddressKey, 0, len(bySource))
	for key := range bySource {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].chainID != keys[j].chainID {
			return keys[i].chainID < keys[j].chainID
		}
		if keys[i].jobType != keys[j].jobType {
			return keys[i].jobType < keys[j].jobType
		}
		if keys[i].pluginType != keys[j].pluginType {
			return keys[i].pluginType < keys[j].pluginType
		}
		return keys[i].fieldPath < keys[j].fieldPath
	})

	out := make([]*commonv1.NodeSubmitterAddress, 0, len(keys))
	for _, key := range keys {
		addresses := make([]string, 0, len(bySource[key]))
		for address := range bySource[key] {
			addresses = append(addresses, address)
		}
		sort.Strings(addresses)
		if len(addresses) == 0 {
			continue
		}
		out = append(out, &commonv1.NodeSubmitterAddress{
			ChainId:    key.chainID,
			JobType:    key.jobType,
			PluginType: key.pluginType,
			FieldPath:  key.fieldPath,
			Addresses:  addresses,
		})
	}
	return out
}
