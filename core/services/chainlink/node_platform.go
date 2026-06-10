package chainlink

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	chainsel "github.com/smartcontractkit/chain-selectors"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
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
	nodeSubmitterFieldTransmitterKeys                    = "transmitterKeys"
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
	Beat               time.Duration
	Lggr               logger.Logger
	CSAKeyStore        keystore.CSA
	JobReader          NodePlatformJobReader
	SubmitterKeyReader NodePlatformSubmitterKeyReader
	CSAPublicKey       string
}

type NodePlatformJobReader interface {
	FindJobs(ctx context.Context, offset, limit int) ([]job.Job, int, error)
}

type NodePlatformSubmitterKeyReader interface {
	SubmitterKeys(ctx context.Context) (map[commontypes.RelayID][]string, error)
}

type nodePlatformRelayerIDReader interface {
	RelayIDs() []commontypes.RelayID
}

type relayerIDsReaderFunc func() []commontypes.RelayID

func (f relayerIDsReaderFunc) RelayIDs() []commontypes.RelayID {
	return f()
}

type nodePlatformSubmitterKeyReader struct {
	keyStore keystore.Master
	relayers nodePlatformRelayerIDReader
}

func (r nodePlatformSubmitterKeyReader) SubmitterKeys(ctx context.Context) (map[commontypes.RelayID][]string, error) {
	if r.keyStore == nil || r.relayers == nil {
		return nil, nil
	}

	out := make(map[commontypes.RelayID][]string)
	for _, relayID := range r.relayers.RelayIDs() {
		keys, err := r.submitterKeysForRelay(ctx, relayID)
		if err != nil {
			return nil, err
		}
		if len(keys) == 0 {
			continue
		}
		out[relayID] = keys
	}
	return out, nil
}

func (r nodePlatformSubmitterKeyReader) submitterKeysForRelay(ctx context.Context, relayID commontypes.RelayID) ([]string, error) {
	switch relayID.Network {
	case relay.NetworkEVM:
		chainID, ok := new(big.Int).SetString(relayID.ChainID, 10)
		if !ok {
			return nil, fmt.Errorf("error parsing chain ID, expected big int: %s", relayID.ChainID)
		}
		ethKeys, err := r.keyStore.Eth().EnabledAddressesForChain(ctx, chainID)
		if err != nil {
			return nil, fmt.Errorf("error getting enabled EVM addresses for chain %s: %w", chainID.String(), err)
		}
		keys := make([]string, 0, len(ethKeys))
		for _, key := range ethKeys {
			keys = append(keys, key.Hex())
		}
		return keys, nil
	case relay.NetworkSolana:
		return nodePlatformKeyIDs(r.keyStore.Solana())
	case relay.NetworkAptos:
		return nodePlatformKeyIDs(r.keyStore.Aptos())
	case relay.NetworkCosmos:
		return nodePlatformKeyIDs(r.keyStore.Cosmos())
	case relay.NetworkStarkNet:
		return nodePlatformKeyIDs(r.keyStore.StarkNet())
	case relay.NetworkTON:
		return nodePlatformKeyIDs(r.keyStore.TON())
	case relay.NetworkSui:
		return nodePlatformKeyIDs(r.keyStore.Sui())
	default:
		return nil, nil
	}
}

func nodePlatformKeyIDs[K keystore.Key](ks interface{ GetAll() ([]K, error) }) ([]string, error) {
	keys, err := ks.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all keys: %w", err)
	}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, key.ID())
	}
	return out, nil
}

func NewNodePlatformJobInfoConfig(opts ApplicationOpts, jobReader NodePlatformJobReader, relayerReader ...RelayerChainInteroperators) NodePlatformJobInfoConfig {
	cfg := NodePlatformJobInfoConfig{
		Beat:        nodePlatformBeat,
		Lggr:        opts.Logger,
		CSAKeyStore: opts.KeyStore.CSA(),
		JobReader:   jobReader,
	}
	if opts.KeyStore != nil && len(relayerReader) > 0 && relayerReader[0] != nil {
		cfg.SubmitterKeyReader = nodePlatformSubmitterKeyReader{
			keyStore: opts.KeyStore,
			relayers: relayerIDsReaderFunc(func() []commontypes.RelayID {
				return slices.Collect(maps.Keys(relayerReader[0].GetIDToRelayerMap()))
			}),
		}
	}
	return cfg
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

	builder := newNodeSubmitterAddressBuilder()
	var (
		submitterKeys       map[commontypes.RelayID][]string
		submitterKeysLoaded bool
	)
	for offset := 0; ; {
		jobs, count, err := s.opts.JobReader.FindJobs(ctx, offset, nodePlatformJobInfoPageSize)
		if err != nil {
			s.eng.Warnw("failed to resolve node-platform submitter addresses", "offset", offset, "limit", nodePlatformJobInfoPageSize, "err", err)
			return nil
		}

		if !submitterKeysLoaded && needsSubmitterKeys(jobs) && s.opts.SubmitterKeyReader != nil {
			submitterKeys, err = s.opts.SubmitterKeyReader.SubmitterKeys(ctx)
			if err != nil {
				s.eng.Warnw("failed to resolve node-platform submitter keys", "err", err)
				submitterKeys = nil
			}
			submitterKeysLoaded = true
		}

		builder.addJobs(jobs, submitterKeys)

		offset += len(jobs)
		if len(jobs) == 0 || offset >= count || len(jobs) < nodePlatformJobInfoPageSize {
			break
		}
	}

	return builder.build()
}

func needsSubmitterKeys(jobs []job.Job) bool {
	return slices.ContainsFunc(jobs, func(jb job.Job) bool {
		return jb.CCIPSpec != nil || jb.CCVExecutorSpec != nil
	})
}

type nodeSubmitterAddressKey struct {
	chainID    string
	jobType    string
	pluginType string
	fieldPath  string
}

type nodeSubmitterAddressBuilder struct {
	bySource map[nodeSubmitterAddressKey]map[string]struct{}
}

func newNodeSubmitterAddressBuilder() *nodeSubmitterAddressBuilder {
	return &nodeSubmitterAddressBuilder{bySource: make(map[nodeSubmitterAddressKey]map[string]struct{})}
}

func (b *nodeSubmitterAddressBuilder) addJobs(jobs []job.Job, submitterKeys map[commontypes.RelayID][]string) {
	for _, jb := range jobs {
		b.addOCRSubmitterAddress(jb)
		b.addOCR2SubmitterAddresses(jb)
		b.addVRFSubmitterAddresses(jb)
		b.addBlockhashStoreSubmitterAddresses(jb)
		b.addBlockHeaderFeederSubmitterAddresses(jb)
		b.addStandardCapabilitiesSubmitterAddress(jb)
		b.addCCIPSubmitterAddresses(jb, submitterKeys)
		b.addCCVExecutorSubmitterAddresses(jb, submitterKeys)
		b.addPipelineETHTxSubmitterAddresses(jb)
	}
}

func (b *nodeSubmitterAddressBuilder) addOCRSubmitterAddress(jb job.Job) {
	spec := jb.OCROracleSpec
	if spec == nil || spec.TransmitterAddress == nil || spec.EVMChainID == nil {
		return
	}
	b.add(spec.EVMChainID.String(), jobType(jb, job.OffchainReporting), "", nodeSubmitterFieldTransmitterAddress, spec.TransmitterAddress.String())
}

func (b *nodeSubmitterAddressBuilder) addOCR2SubmitterAddresses(jb job.Job) {
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
		b.add(chainID, jobType(jb, job.OffchainReporting2), pluginType, nodeSubmitterFieldTransmitterID, spec.TransmitterID.String)
	}
	if sendingKeys, err := job.SendingKeysForJob(spec); err == nil {
		b.add(chainID, jobType(jb, job.OffchainReporting2), pluginType, nodeSubmitterFieldRelayConfigSendingKeys, sendingKeys...)
	}
	b.add(chainID, jobType(jb, job.OffchainReporting2), pluginType, nodeSubmitterFieldDualTransmissionTransmitterAddress, dualTransmissionTransmitterAddress(spec.RelayConfig))
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

func (b *nodeSubmitterAddressBuilder) addVRFSubmitterAddresses(jb job.Job) {
	spec := jb.VRFSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	b.add(spec.EVMChainID.String(), jobType(jb, job.VRF), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func (b *nodeSubmitterAddressBuilder) addBlockhashStoreSubmitterAddresses(jb job.Job) {
	spec := jb.BlockhashStoreSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	b.add(spec.EVMChainID.String(), jobType(jb, job.BlockhashStore), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func (b *nodeSubmitterAddressBuilder) addBlockHeaderFeederSubmitterAddresses(jb job.Job) {
	spec := jb.BlockHeaderFeederSpec
	if spec == nil || spec.EVMChainID == nil {
		return
	}
	b.add(spec.EVMChainID.String(), jobType(jb, job.BlockHeaderFeeder), "", nodeSubmitterFieldFromAddresses, eip55AddressStrings(spec.FromAddresses)...)
}

func (b *nodeSubmitterAddressBuilder) addStandardCapabilitiesSubmitterAddress(jb job.Job) {
	spec := jb.StandardCapabilitiesSpec
	if spec == nil || !spec.OracleFactory.Enabled {
		return
	}
	b.add(spec.OracleFactory.ChainID, jobType(jb, job.StandardCapabilities), "", nodeSubmitterFieldOracleFactoryTransmitterID, spec.OracleFactory.TransmitterID)
}

func (b *nodeSubmitterAddressBuilder) addCCIPSubmitterAddresses(jb job.Job, submitterKeys map[commontypes.RelayID][]string) {
	if jb.CCIPSpec == nil || len(jb.CCIPSpec.P2PV2Bootstrappers) == 0 || len(submitterKeys) == 0 {
		return
	}
	for relayID, addresses := range submitterKeys {
		if len(addresses) == 0 {
			continue
		}
		b.add(relayID.ChainID, jobType(jb, job.CCIP), "", nodeSubmitterFieldTransmitterKeys, addresses[0])
	}
}

func (b *nodeSubmitterAddressBuilder) addCCVExecutorSubmitterAddresses(jb job.Job, submitterKeys map[commontypes.RelayID][]string) {
	spec := jb.CCVExecutorSpec
	if spec == nil || len(submitterKeys) == 0 {
		return
	}
	for _, chainID := range ccvExecutorEVMChainIDs(spec.ExecutorConfig) {
		relayID := commontypes.NewRelayID(relay.NetworkEVM, chainID)
		b.add(chainID, jobType(jb, job.CCVExecutor), "", nodeSubmitterFieldFromAddresses, submitterKeys[relayID]...)
	}
}

func (b *nodeSubmitterAddressBuilder) addPipelineETHTxSubmitterAddresses(jb job.Job) {
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
		b.add(chainID, jobType(jb, ""), "", nodeSubmitterFieldObservationSourceETHTxFrom, addresses...)
	}
}

func ccvExecutorEVMChainIDs(raw string) []string {
	tree, err := toml.Load(strings.TrimSpace(raw))
	if err != nil {
		return nil
	}
	chainConfigRaw := tree.Get("chain_configuration")
	chainConfig, ok := chainConfigRaw.(*toml.Tree)
	if !ok || chainConfig == nil {
		return nil
	}
	chainIDs := make([]string, 0, len(chainConfig.Keys()))
	for _, selector := range chainConfig.Keys() {
		selectorID, err := strconv.ParseUint(selector, 10, 64)
		if err != nil {
			continue
		}
		chainID, err := chainsel.GetChainIDFromSelector(selectorID)
		if err != nil {
			continue
		}
		chainIDs = append(chainIDs, chainID)
	}
	return chainIDs
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

func (b *nodeSubmitterAddressBuilder) add(chainID, jobType, pluginType, fieldPath string, addresses ...string) {
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
	if b.bySource == nil {
		b.bySource = make(map[nodeSubmitterAddressKey]map[string]struct{})
	}
	for _, address := range addresses {
		address = strings.TrimSpace(address)
		if address == "" {
			continue
		}
		if b.bySource[key] == nil {
			b.bySource[key] = make(map[string]struct{})
		}
		b.bySource[key][address] = struct{}{}
	}
}

func (b *nodeSubmitterAddressBuilder) build() []*commonv1.NodeSubmitterAddress {
	if len(b.bySource) == 0 {
		return nil
	}
	keys := make([]nodeSubmitterAddressKey, 0, len(b.bySource))
	for key := range b.bySource {
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
		if len(b.bySource[key]) == 0 {
			continue
		}
		addresses := slices.Sorted(maps.Keys(b.bySource[key]))
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
