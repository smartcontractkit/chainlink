package resolver

import (
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/gqlscalar"
)

type SpecResolver struct {
	j job.Job
}

func NewSpec(j job.Job) *SpecResolver {
	return &SpecResolver{j: j}
}

func (r *SpecResolver) ToCronSpec() (*CronSpecResolver, bool) {
	if r.j.Type != job.Cron {
		return nil, false
	}

	return &CronSpecResolver{spec: *r.j.CronSpec}, true
}

func (r *SpecResolver) ToDirectRequestSpec() (*DirectRequestSpecResolver, bool) {
	if r.j.Type != job.DirectRequest {
		return nil, false
	}

	return &DirectRequestSpecResolver{spec: *r.j.DirectRequestSpec}, true
}

func (r *SpecResolver) ToFluxMonitorSpec() (*FluxMonitorSpecResolver, bool) {
	if r.j.Type != job.FluxMonitor {
		return nil, false
	}

	return &FluxMonitorSpecResolver{spec: *r.j.FluxMonitorSpec}, true
}

func (r *SpecResolver) ToKeeperSpec() (*KeeperSpecResolver, bool) {
	if r.j.Type != job.Keeper {
		return nil, false
	}

	return &KeeperSpecResolver{spec: *r.j.KeeperSpec}, true
}

func (r *SpecResolver) ToOCRSpec() (*OCRSpecResolver, bool) {
	if r.j.Type != job.OffchainReporting {
		return nil, false
	}

	return &OCRSpecResolver{spec: *r.j.OCROracleSpec}, true
}

func (r *SpecResolver) ToOCR2Spec() (*OCR2SpecResolver, bool) {
	if r.j.Type != job.OffchainReporting2 {
		return nil, false
	}

	return &OCR2SpecResolver{spec: *r.j.OCR2OracleSpec}, true
}

func (r *SpecResolver) ToVRFSpec() (*VRFSpecResolver, bool) {
	if r.j.Type != job.VRF {
		return nil, false
	}

	return &VRFSpecResolver{spec: *r.j.VRFSpec}, true
}

func (r *SpecResolver) ToWebhookSpec() (*WebhookSpecResolver, bool) {
	if r.j.Type != job.Webhook {
		return nil, false
	}

	return &WebhookSpecResolver{spec: *r.j.WebhookSpec}, true
}

// ToBlockhashStoreSpec returns the BlockhashStoreSpec from the SpecResolver if the job is a
// BlockhashStore job.
func (r *SpecResolver) ToBlockhashStoreSpec() (*BlockhashStoreSpecResolver, bool) {
	if r.j.Type != job.BlockhashStore {
		return nil, false
	}

	return &BlockhashStoreSpecResolver{spec: *r.j.BlockhashStoreSpec}, true
}

// ToBlockHeaderFeederSpec returns the BlockHeaderFeederSpec from the SpecResolver if the job is a
// BlockHeaderFeeder job.
func (r *SpecResolver) ToBlockHeaderFeederSpec() (*BlockHeaderFeederSpecResolver, bool) {
	if r.j.Type != job.BlockHeaderFeeder {
		return nil, false
	}

	return &BlockHeaderFeederSpecResolver{spec: *r.j.BlockHeaderFeederSpec}, true
}

// ToBootstrapSpec resolves to the Booststrap Spec Resolver
func (r *SpecResolver) ToBootstrapSpec() (*BootstrapSpecResolver, bool) {
	if r.j.Type != job.Bootstrap {
		return nil, false
	}

	return &BootstrapSpecResolver{spec: *r.j.BootstrapSpec}, true
}

func (r *SpecResolver) ToGatewaySpec() (*GatewaySpecResolver, bool) {
	if r.j.Type != job.Gateway {
		return nil, false
	}

	return &GatewaySpecResolver{spec: *r.j.GatewaySpec}, true
}

func (r *SpecResolver) ToWorkflowSpec() (*WorkflowSpecResolver, bool) {
	if r.j.Type != job.Workflow {
		return nil, false
	}

	return &WorkflowSpecResolver{spec: *r.j.WorkflowSpec}, true
}

func (r *SpecResolver) ToStandardCapabilitiesSpec() (*StandardCapabilitiesSpecResolver, bool) {
	if r.j.Type != job.StandardCapabilities {
		return nil, false
	}

	return &StandardCapabilitiesSpecResolver{spec: *r.j.StandardCapabilitiesSpec}, true
}

func (r *SpecResolver) ToStreamSpec() (*StreamSpecResolver, bool) {
	if r.j.Type != job.Stream {
		return nil, false
	}
	res := &StreamSpecResolver{}
	if r.j.StreamID != nil {
		sid := strconv.FormatUint(uint64(*r.j.StreamID), 10)
		res.streamID = &sid
	}

	return res, true
}

func (r *SpecResolver) ToCCIPSpec() (*CCIPSpecResolver, bool) {
	if r.j.Type != job.CCIP {
		return nil, false
	}

	return &CCIPSpecResolver{spec: *r.j.CCIPSpec}, true
}

type CronSpecResolver struct {
	spec job.CronSpec
}

func (r *CronSpecResolver) Schedule() string {
	return r.spec.CronSchedule
}

// EVMChainID resolves the spec's evm chain id.
func (r *CronSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// CreatedAt resolves the spec's created at timestamp.
func (r *CronSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

type DirectRequestSpecResolver struct {
	spec job.DirectRequestSpec
}

// ContractAddress resolves the spec's contract address.
func (r *DirectRequestSpecResolver) ContractAddress() string {
	return r.spec.ContractAddress.String()
}

// CreatedAt resolves the spec's created at timestamp.
func (r *DirectRequestSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// EVMChainID resolves the spec's evm chain id.
func (r *DirectRequestSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// MinIncomingConfirmations resolves the spec's min incoming confirmations.
func (r *DirectRequestSpecResolver) MinIncomingConfirmations() int32 {
	if r.spec.MinIncomingConfirmations.Valid {
		return int32(r.spec.MinIncomingConfirmations.Uint32)
	}

	return 0
}

// MinContractPaymentLinkJuels resolves the spec's min contract payment link.
func (r *DirectRequestSpecResolver) MinContractPaymentLinkJuels() string {
	return r.spec.MinContractPayment.String()
}

// Requesters resolves the spec's evm chain id.
func (r *DirectRequestSpecResolver) Requesters() *[]string {
	if r.spec.Requesters == nil {
		return nil
	}

	requesters := r.spec.Requesters.ToStrings()

	return &requesters
}

type FluxMonitorSpecResolver struct {
	spec job.FluxMonitorSpec
}

// AbsoluteThreshold resolves the spec's absolute deviation threshold.
func (r *FluxMonitorSpecResolver) AbsoluteThreshold() float64 {
	return float64(r.spec.AbsoluteThreshold)
}

// ContractAddress resolves the spec's contract address.
func (r *FluxMonitorSpecResolver) ContractAddress() string {
	return r.spec.ContractAddress.String()
}

// CreatedAt resolves the spec's created at timestamp.
func (r *FluxMonitorSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// AbsoluteThreshold resolves the spec's absolute threshold.
func (r *FluxMonitorSpecResolver) DrumbeatEnabled() bool {
	return r.spec.DrumbeatEnabled
}

// DrumbeatRandomDelay resolves the spec's drumbeat random delay.
func (r *FluxMonitorSpecResolver) DrumbeatRandomDelay() *string {
	var delay *string
	if r.spec.DrumbeatRandomDelay > 0 {
		drumbeatRandomDelay := r.spec.DrumbeatRandomDelay.String()
		delay = &drumbeatRandomDelay
	}

	return delay
}

// DrumbeatSchedule resolves the spec's drumbeat schedule.
func (r *FluxMonitorSpecResolver) DrumbeatSchedule() *string {
	if r.spec.DrumbeatEnabled {
		return &r.spec.DrumbeatSchedule
	}

	return nil
}

// EVMChainID resolves the spec's evm chain id.
func (r *FluxMonitorSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// IdleTimerDisabled resolves the spec's idle timer disabled flag.
func (r *FluxMonitorSpecResolver) IdleTimerDisabled() bool {
	return r.spec.IdleTimerDisabled
}

// IdleTimerPeriod resolves the spec's idle timer period.
func (r *FluxMonitorSpecResolver) IdleTimerPeriod() string {
	return r.spec.IdleTimerPeriod.String()
}

// MinPayment resolves the spec's min payment.
func (r *FluxMonitorSpecResolver) MinPayment() *string {
	if r.spec.MinPayment != nil {
		min := r.spec.MinPayment.String()

		return &min
	}
	return nil
}

// PollTimerDisabled resolves the spec's poll timer disabled flag.
func (r *FluxMonitorSpecResolver) PollTimerDisabled() bool {
	return r.spec.PollTimerDisabled
}

// PollTimerPeriod resolves the spec's poll timer period.
func (r *FluxMonitorSpecResolver) PollTimerPeriod() string {
	return r.spec.PollTimerPeriod.String()
}

// Threshold resolves the spec's deviation threshold.
func (r *FluxMonitorSpecResolver) Threshold() float64 {
	return float64(r.spec.Threshold)
}

type KeeperSpecResolver struct {
	spec job.KeeperSpec
}

// ContractAddress resolves the spec's contract address.
func (r *KeeperSpecResolver) ContractAddress() string {
	return r.spec.ContractAddress.String()
}

// CreatedAt resolves the spec's created at timestamp.
func (r *KeeperSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// EVMChainID resolves the spec's evm chain id.
func (r *KeeperSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// FromAddress resolves the spec's from contract address.
//
// Because VRF has an non required field of the same name, we have to be
// consistent in our return value of using a *string instead of a string even
// though this is a required field for the KeeperSpec.
//
// http://spec.graphql.org/draft/#sec-Field-Selection-Merging
func (r *KeeperSpecResolver) FromAddress() *string {
	addr := r.spec.FromAddress.String()

	return &addr
}

type OCRSpecResolver struct {
	spec job.OCROracleSpec
}

// BlockchainTimeout resolves the spec's blockchain timeout.
func (r *OCRSpecResolver) BlockchainTimeout() *string {
	if r.spec.BlockchainTimeout.Duration() == 0 {
		return nil
	}

	timeout := r.spec.BlockchainTimeout.Duration().String()

	return &timeout
}

// ContractAddress resolves the spec's contract address.
func (r *OCRSpecResolver) ContractAddress() string {
	return r.spec.ContractAddress.String()
}

// ContractConfigConfirmations resolves the spec's confirmations config.
func (r *OCRSpecResolver) ContractConfigConfirmations() *int32 {
	if r.spec.ContractConfigConfirmations == 0 {
		return nil
	}

	confirmations := int32(r.spec.ContractConfigConfirmations)

	return &confirmations
}

// ContractConfigTrackerPollInterval resolves the spec's contract tracker poll
// interval config.
func (r *OCRSpecResolver) ContractConfigTrackerPollInterval() *string {
	if r.spec.ContractConfigTrackerPollInterval.Duration() == 0 {
		return nil
	}

	interval := r.spec.ContractConfigTrackerPollInterval.Duration().String()

	return &interval
}

// ContractConfigTrackerSubscribeInterval resolves the spec's tracker subscribe
// interval config.
func (r *OCRSpecResolver) ContractConfigTrackerSubscribeInterval() *string {
	if r.spec.ContractConfigTrackerSubscribeInterval.Duration() == 0 {
		return nil
	}

	interval := r.spec.ContractConfigTrackerSubscribeInterval.Duration().String()

	return &interval
}

// CreatedAt resolves the spec's created at timestamp.
func (r *OCRSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// EVMChainID resolves the spec's evm chain id.
func (r *OCRSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// DatabaseTimeout resolves the spec's database timeout.
func (r *OCRSpecResolver) DatabaseTimeout() string {
	return r.spec.DatabaseTimeout.Duration().String()
}

// ObservationGracePeriod resolves the spec's observation grace period.
func (r *OCRSpecResolver) ObservationGracePeriod() string {
	return r.spec.ObservationGracePeriod.Duration().String()
}

// ContractTransmitterTransmitTimeout resolves the spec's contract transmitter transmit timeout.
func (r *OCRSpecResolver) ContractTransmitterTransmitTimeout() string {
	return r.spec.ContractTransmitterTransmitTimeout.Duration().String()
}

// IsBootstrapPeer resolves whether spec is a bootstrap peer.
func (r *OCRSpecResolver) IsBootstrapPeer() bool {
	return r.spec.IsBootstrapPeer
}

// KeyBundleID resolves the spec's key bundle id.
func (r *OCRSpecResolver) KeyBundleID() *string {
	if r.spec.EncryptedOCRKeyBundleID == nil {
		return nil
	}

	bundleID := r.spec.EncryptedOCRKeyBundleID.String()

	return &bundleID
}

// ObservationTimeout resolves the spec's observation timeout
func (r *OCRSpecResolver) ObservationTimeout() *string {
	if r.spec.ObservationTimeout.Duration() == 0 {
		return nil
	}

	timeout := r.spec.ObservationTimeout.Duration().String()

	return &timeout
}

// P2PV2Bootstrappers resolves the OCR1 spec's p2pv2 bootstrappers
func (r *OCRSpecResolver) P2PV2Bootstrappers() *[]string {
	if len(r.spec.P2PV2Bootstrappers) == 0 {
		return nil
	}

	peers := []string(r.spec.P2PV2Bootstrappers)

	return &peers
}

// TransmitterAddress resolves the spec's transmitter address
func (r *OCRSpecResolver) TransmitterAddress() *string {
	if r.spec.TransmitterAddress == nil {
		return nil
	}

	addr := r.spec.TransmitterAddress.String()
	return &addr
}

type OCR2SpecResolver struct {
	spec job.OCR2OracleSpec
}

// BlockchainTimeout resolves the spec's blockchain timeout.
func (r *OCR2SpecResolver) BlockchainTimeout() *string {
	if r.spec.BlockchainTimeout.Duration() == 0 {
		return nil
	}

	timeout := r.spec.BlockchainTimeout.Duration().String()

	return &timeout
}

// ContractID resolves the spec's contract address.
func (r *OCR2SpecResolver) ContractID() string {
	return r.spec.ContractID
}

// ContractConfigConfirmations resolves the spec's confirmations config.
func (r *OCR2SpecResolver) ContractConfigConfirmations() *int32 {
	if r.spec.ContractConfigConfirmations == 0 {
		return nil
	}

	confirmations := int32(r.spec.ContractConfigConfirmations)

	return &confirmations
}

// ContractConfigTrackerPollInterval resolves the spec's contract tracker poll
// interval config.
func (r *OCR2SpecResolver) ContractConfigTrackerPollInterval() *string {
	if r.spec.ContractConfigTrackerPollInterval.Duration() == 0 {
		return nil
	}

	interval := r.spec.ContractConfigTrackerPollInterval.Duration().String()

	return &interval
}

// CreatedAt resolves the spec's created at timestamp.
func (r *OCR2SpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// OcrKeyBundleID resolves the spec's key bundle id.
func (r *OCR2SpecResolver) OcrKeyBundleID() *string {
	if !r.spec.OCRKeyBundleID.Valid {
		return nil
	}

	return &r.spec.OCRKeyBundleID.String
}

// MonitoringEndpoint resolves the spec's monitoring endpoint
func (r *OCR2SpecResolver) MonitoringEndpoint() *string {
	if !r.spec.MonitoringEndpoint.Valid {
		return nil
	}

	return &r.spec.MonitoringEndpoint.String
}

// P2PV2Bootstrappers resolves the OCR2 spec's p2pv2 bootstrappers
func (r *OCR2SpecResolver) P2PV2Bootstrappers() *[]string {
	if len(r.spec.P2PV2Bootstrappers) == 0 {
		return nil
	}

	peers := []string(r.spec.P2PV2Bootstrappers)

	return &peers
}

// Relay resolves the spec's relay
func (r *OCR2SpecResolver) Relay() string {
	return r.spec.Relay
}

// RelayConfig resolves the spec's relay config
func (r *OCR2SpecResolver) RelayConfig() gqlscalar.Map {
	return gqlscalar.Map(r.spec.RelayConfig)
}

// PluginType resolves the spec's plugin type
func (r *OCR2SpecResolver) PluginType() string {
	return string(r.spec.PluginType)
}

// PluginConfig resolves the spec's plugin config
func (r *OCR2SpecResolver) PluginConfig() gqlscalar.Map {
	return gqlscalar.Map(r.spec.PluginConfig)
}

// TransmitterID resolves the spec's transmitter id
func (r *OCR2SpecResolver) TransmitterID() *string {
	if !r.spec.TransmitterID.Valid {
		return nil
	}

	addr := r.spec.TransmitterID.String
	return &addr
}

// FeedID resolves the spec's feed ID
func (r *OCR2SpecResolver) FeedID() *string {
	if r.spec.FeedID == nil {
		return nil
	}
	feedID := r.spec.FeedID.String()
	return &feedID
}

func (r *OCR2SpecResolver) AllowNoBootstrappers() bool {
	return r.spec.AllowNoBootstrappers
}

type VRFSpecResolver struct {
	spec job.VRFSpec
}

// MinIncomingConfirmations resolves the spec's min incoming confirmations.
func (r *VRFSpecResolver) MinIncomingConfirmations() int32 {
	return int32(r.spec.MinIncomingConfirmations)
}

// CoordinatorAddress resolves the spec's coordinator address.
func (r *VRFSpecResolver) CoordinatorAddress() string {
	return r.spec.CoordinatorAddress.String()
}

// CreatedAt resolves the spec's created at timestamp.
func (r *VRFSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// EVMChainID resolves the spec's evm chain id.
func (r *VRFSpecResolver) EVMChainID() *string {
	if r.spec.EVMChainID == nil {
		return nil
	}

	chainID := r.spec.EVMChainID.String()

	return &chainID
}

// FromAddresses resolves the spec's from addresses.
func (r *VRFSpecResolver) FromAddresses() *[]string {
	if len(r.spec.FromAddresses) == 0 {
		return nil
	}

	var addresses []string
	for _, a := range r.spec.FromAddresses {
		addresses = append(addresses, a.Address().String())
	}
	return &addresses
}

// PollPeriod resolves the spec's poll period.
func (r *VRFSpecResolver) PollPeriod() string {
	return r.spec.PollPeriod.String()
}

// PublicKey resolves the spec's public key.
func (r *VRFSpecResolver) PublicKey() string {
	return r.spec.PublicKey.String()
}

// RequestedConfsDelay resolves the spec's requested conf delay.
func (r *VRFSpecResolver) RequestedConfsDelay() int32 {
	// GraphQL doesn't support 64 bit integers, so we have to cast.
	return int32(r.spec.RequestedConfsDelay)
}

// RequestTimeout resolves the spec's request timeout.
func (r *VRFSpecResolver) RequestTimeout() string {
	return r.spec.RequestTimeout.String()
}

// BatchCoordinatorAddress resolves the spec's batch coordinator address.
func (r *VRFSpecResolver) BatchCoordinatorAddress() *string {
	if r.spec.BatchCoordinatorAddress == nil {
		return nil
	}
	addr := r.spec.BatchCoordinatorAddress.String()
	return &addr
}

// BatchFulfillmentEnabled resolves the spec's batch fulfillment enabled flag.
func (r *VRFSpecResolver) BatchFulfillmentEnabled() bool {
	return r.spec.BatchFulfillmentEnabled
}

// BatchFulfillmentGasMultiplier resolves the spec's batch fulfillment gas multiplier.
func (r *VRFSpecResolver) BatchFulfillmentGasMultiplier() float64 {
	return float64(r.spec.BatchFulfillmentGasMultiplier)
}

// CustomRevertsPipelineEnabled resolves the spec's custom reverts pipeline enabled flag.
func (r *VRFSpecResolver) CustomRevertsPipelineEnabled() *bool {
	return &r.spec.CustomRevertsPipelineEnabled
}

// ChunkSize resolves the spec's chunk size.
func (r *VRFSpecResolver) ChunkSize() int32 {
	return int32(r.spec.ChunkSize)
}

// BackoffInitialDelay resolves the spec's backoff initial delay.
func (r *VRFSpecResolver) BackoffInitialDelay() string {
	return r.spec.BackoffInitialDelay.String()
}

// BackoffMaxDelay resolves the spec's backoff max delay.
func (r *VRFSpecResolver) BackoffMaxDelay() string {
	return r.spec.BackoffMaxDelay.String()
}

// GasLanePrice resolves the spec's gas lane price.
func (r *VRFSpecResolver) GasLanePrice() *string {
	if r.spec.GasLanePrice == nil {
		return nil
	}
	gasLanePriceGWei := r.spec.GasLanePrice.String()
	return &gasLanePriceGWei
}

// VRFOwnerAddress resolves the spec's vrf owner address.
func (r *VRFSpecResolver) VRFOwnerAddress() *string {
	if r.spec.VRFOwnerAddress == nil {
		return nil
	}
	vrfOwnerAddress := r.spec.VRFOwnerAddress.String()
	return &vrfOwnerAddress
}

type WebhookSpecResolver struct {
	spec job.WebhookSpec
}

// CreatedAt resolves the spec's created at timestamp.
func (r *WebhookSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// BlockhashStoreSpecResolver exposes the job parameters for a BlockhashStoreSpec.
type BlockhashStoreSpecResolver struct {
	spec job.BlockhashStoreSpec
}

// CoordinatorV1Address returns the address of the V1 Coordinator, if any.
func (b *BlockhashStoreSpecResolver) CoordinatorV1Address() *string {
	if b.spec.CoordinatorV1Address == nil {
		return nil
	}
	addr := b.spec.CoordinatorV1Address.String()
	return &addr
}

// CoordinatorV2Address returns the address of the V2 Coordinator, if any.
func (b *BlockhashStoreSpecResolver) CoordinatorV2Address() *string {
	if b.spec.CoordinatorV2Address == nil {
		return nil
	}
	addr := b.spec.CoordinatorV2Address.String()
	return &addr
}

// CoordinatorV2PlusAddress returns the address of the V2Plus Coordinator, if any.
func (b *BlockhashStoreSpecResolver) CoordinatorV2PlusAddress() *string {
	if b.spec.CoordinatorV2PlusAddress == nil {
		return nil
	}
	addr := b.spec.CoordinatorV2PlusAddress.String()
	return &addr
}

// WaitBlocks returns the job's WaitBlocks param.
func (b *BlockhashStoreSpecResolver) WaitBlocks() int32 {
	return b.spec.WaitBlocks
}

// LookbackBlocks returns the job's LookbackBlocks param.
func (b *BlockhashStoreSpecResolver) LookbackBlocks() int32 {
	return b.spec.LookbackBlocks
}

// HeartbeatPeriod returns the job's HeartbeatPeriod param.
func (b *BlockhashStoreSpecResolver) HeartbeatPeriod() string {
	return b.spec.HeartbeatPeriod.String()
}

// BlockhashStoreAddress returns the job's BlockhashStoreAddress param.
func (b *BlockhashStoreSpecResolver) BlockhashStoreAddress() string {
	return b.spec.BlockhashStoreAddress.String()
}

// TrustedBlockhashStoreAddress returns the address of the job's TrustedBlockhashStoreAddress, if any.
func (b *BlockhashStoreSpecResolver) TrustedBlockhashStoreAddress() *string {
	if b.spec.TrustedBlockhashStoreAddress == nil {
		return nil
	}
	addr := b.spec.TrustedBlockhashStoreAddress.String()
	return &addr
}

// BatchBlockhashStoreAddress returns the job's BatchBlockhashStoreAddress param.
func (b *BlockhashStoreSpecResolver) TrustedBlockhashStoreBatchSize() int32 {
	return b.spec.TrustedBlockhashStoreBatchSize
}

// PollPeriod return's the job's PollPeriod param.
func (b *BlockhashStoreSpecResolver) PollPeriod() string {
	return b.spec.PollPeriod.String()
}

// RunTimeout return's the job's RunTimeout param.
func (b *BlockhashStoreSpecResolver) RunTimeout() string {
	return b.spec.RunTimeout.String()
}

// EVMChainID returns the job's EVMChainID param.
func (b *BlockhashStoreSpecResolver) EVMChainID() *string {
	chainID := b.spec.EVMChainID.String()
	return &chainID
}

// FromAddress returns the job's FromAddress param, if any.
func (b *BlockhashStoreSpecResolver) FromAddresses() *[]string {
	if b.spec.FromAddresses == nil {
		return nil
	}
	var addresses []string
	for _, a := range b.spec.FromAddresses {
		addresses = append(addresses, a.Address().String())
	}
	return &addresses
}

// CreatedAt resolves the spec's created at timestamp.
func (b *BlockhashStoreSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: b.spec.CreatedAt}
}

// BlockHeaderFeederSpecResolver exposes the job parameters for a BlockHeaderFeederSpec.
type BlockHeaderFeederSpecResolver struct {
	spec job.BlockHeaderFeederSpec
}

// CoordinatorV1Address returns the address of the V1 Coordinator, if any.
func (b *BlockHeaderFeederSpecResolver) CoordinatorV1Address() *string {
	if b.spec.CoordinatorV1Address == nil {
		return nil
	}
	addr := b.spec.CoordinatorV1Address.String()
	return &addr
}

// CoordinatorV2Address returns the address of the V2 Coordinator, if any.
func (b *BlockHeaderFeederSpecResolver) CoordinatorV2Address() *string {
	if b.spec.CoordinatorV2Address == nil {
		return nil
	}
	addr := b.spec.CoordinatorV2Address.String()
	return &addr
}

// CoordinatorV2PlusAddress returns the address of the V2 Coordinator Plus, if any.
func (b *BlockHeaderFeederSpecResolver) CoordinatorV2PlusAddress() *string {
	if b.spec.CoordinatorV2PlusAddress == nil {
		return nil
	}
	addr := b.spec.CoordinatorV2PlusAddress.String()
	return &addr
}

// WaitBlocks returns the job's WaitBlocks param.
func (b *BlockHeaderFeederSpecResolver) WaitBlocks() int32 {
	return b.spec.WaitBlocks
}

// LookbackBlocks returns the job's LookbackBlocks param.
func (b *BlockHeaderFeederSpecResolver) LookbackBlocks() int32 {
	return b.spec.LookbackBlocks
}

// BlockhashStoreAddress returns the job's BlockhashStoreAddress param.
func (b *BlockHeaderFeederSpecResolver) BlockhashStoreAddress() string {
	return b.spec.BlockhashStoreAddress.String()
}

// BatchBlockhashStoreAddress returns the job's BatchBlockhashStoreAddress param.
func (b *BlockHeaderFeederSpecResolver) BatchBlockhashStoreAddress() string {
	return b.spec.BatchBlockhashStoreAddress.String()
}

// PollPeriod return's the job's PollPeriod param.
func (b *BlockHeaderFeederSpecResolver) PollPeriod() string {
	return b.spec.PollPeriod.String()
}

// RunTimeout return's the job's RunTimeout param.
func (b *BlockHeaderFeederSpecResolver) RunTimeout() string {
	return b.spec.RunTimeout.String()
}

// EVMChainID returns the job's EVMChainID param.
func (b *BlockHeaderFeederSpecResolver) EVMChainID() *string {
	chainID := b.spec.EVMChainID.String()
	return &chainID
}

// FromAddress returns the job's FromAddress param, if any.
func (b *BlockHeaderFeederSpecResolver) FromAddresses() *[]string {
	if b.spec.FromAddresses == nil {
		return nil
	}
	var addresses []string
	for _, a := range b.spec.FromAddresses {
		addresses = append(addresses, a.Address().String())
	}
	return &addresses
}

// GetBlockhashesBatchSize returns the job's GetBlockhashesBatchSize param.
func (b *BlockHeaderFeederSpecResolver) GetBlockhashesBatchSize() int32 {
	return int32(b.spec.GetBlockhashesBatchSize)
}

// StoreBlockhashesBatchSize returns the job's StoreBlockhashesBatchSize param.
func (b *BlockHeaderFeederSpecResolver) StoreBlockhashesBatchSize() int32 {
	return int32(b.spec.StoreBlockhashesBatchSize)
}

// CreatedAt resolves the spec's created at timestamp.
func (b *BlockHeaderFeederSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: b.spec.CreatedAt}
}

// BootstrapSpecResolver defines the Bootstrap Spec Resolver
type BootstrapSpecResolver struct {
	spec job.BootstrapSpec
}

// ID resolves the Bootstrap spec ID
func (r *BootstrapSpecResolver) ID() graphql.ID {
	return graphql.ID(stringutils.FromInt32(r.spec.ID))
}

// ContractID resolves the spec's contract address
func (r *BootstrapSpecResolver) ContractID() string {
	return r.spec.ContractID
}

// Relay resolves the spec's relay
func (r *BootstrapSpecResolver) Relay() string {
	return r.spec.Relay
}

// RelayConfig resolves the spec's relay config
func (r *BootstrapSpecResolver) RelayConfig() gqlscalar.Map {
	return gqlscalar.Map(r.spec.RelayConfig)
}

// MonitoringEndpoint resolves the spec's monitoring endpoint
func (r *BootstrapSpecResolver) MonitoringEndpoint() *string {
	if !r.spec.MonitoringEndpoint.Valid {
		return nil
	}

	return &r.spec.MonitoringEndpoint.String
}

// BlockchainTimeout resolves the spec's blockchain timeout
func (r *BootstrapSpecResolver) BlockchainTimeout() *string {
	if r.spec.BlockchainTimeout.Duration() == 0 {
		return nil
	}

	interval := r.spec.BlockchainTimeout.Duration().String()

	return &interval
}

// ContractConfigTrackerPollInterval resolves the spec's contract tracker poll
// interval config.
func (r *BootstrapSpecResolver) ContractConfigTrackerPollInterval() *string {
	if r.spec.ContractConfigTrackerPollInterval.Duration() == 0 {
		return nil
	}

	interval := r.spec.ContractConfigTrackerPollInterval.Duration().String()

	return &interval
}

// ContractConfigConfirmations resolves the spec's confirmations config.
func (r *BootstrapSpecResolver) ContractConfigConfirmations() *int32 {
	if r.spec.ContractConfigConfirmations == 0 {
		return nil
	}

	confirmations := int32(r.spec.ContractConfigConfirmations)

	return &confirmations
}

// RelayConfig resolves the spec's onchain signing strategy config
func (r *OCR2SpecResolver) OnchainSigningStrategy() gqlscalar.Map {
	return gqlscalar.Map(r.spec.OnchainSigningStrategy)
}

// CreatedAt resolves the spec's created at timestamp.
func (r *BootstrapSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

type GatewaySpecResolver struct {
	spec job.GatewaySpec
}

func (r *GatewaySpecResolver) ID() graphql.ID {
	return graphql.ID(stringutils.FromInt32(r.spec.ID))
}

func (r *GatewaySpecResolver) GatewayConfig() gqlscalar.Map {
	return gqlscalar.Map(r.spec.GatewayConfig)
}

func (r *GatewaySpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

type WorkflowSpecResolver struct {
	spec job.WorkflowSpec
}

func (r *WorkflowSpecResolver) ID() graphql.ID {
	return graphql.ID(stringutils.FromInt32(r.spec.ID))
}

func (r *WorkflowSpecResolver) WorkflowID() string {
	return r.spec.WorkflowID
}

func (r *WorkflowSpecResolver) Workflow() string {
	return r.spec.Workflow
}

func (r *WorkflowSpecResolver) WorkflowOwner() string {
	return r.spec.WorkflowOwner
}

func (r *WorkflowSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

func (r *WorkflowSpecResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.UpdatedAt}
}

type StandardCapabilitiesSpecResolver struct {
	spec job.StandardCapabilitiesSpec
}

func (r *StandardCapabilitiesSpecResolver) ID() graphql.ID {
	return graphql.ID(stringutils.FromInt32(r.spec.ID))
}

func (r *StandardCapabilitiesSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

func (r *StandardCapabilitiesSpecResolver) Command() string {
	return r.spec.Command
}

func (r *StandardCapabilitiesSpecResolver) Config() *string {
	return &r.spec.Config
}

type StreamSpecResolver struct {
	streamID *string
}

func (r *StreamSpecResolver) StreamID() *string {
	return r.streamID
}

type CCIPSpecResolver struct {
	spec job.CCIPSpec
}

func (r *CCIPSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

func (r *CCIPSpecResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.UpdatedAt}
}

func (r *CCIPSpecResolver) ID() graphql.ID {
	return graphql.ID(stringutils.FromInt32(r.spec.ID))
}

func (r *CCIPSpecResolver) CapabilityLabelledName() string {
	return r.spec.CapabilityLabelledName
}

func (r *CCIPSpecResolver) CapabilityVersion() string {
	return r.spec.CapabilityVersion
}

func (r *CCIPSpecResolver) P2PV2Bootstrappers() *[]string {
	if len(r.spec.P2PV2Bootstrappers) == 0 {
		return nil
	}

	peers := []string(r.spec.P2PV2Bootstrappers)

	return &peers
}

func (r *CCIPSpecResolver) OCRKeyBundleIDs() gqlscalar.Map {
	return gqlscalar.Map(r.spec.OCRKeyBundleIDs)
}

func (r *CCIPSpecResolver) RelayConfigs() gqlscalar.Map {
	return gqlscalar.Map(r.spec.RelayConfigs)
}

func (r *CCIPSpecResolver) PluginConfig() gqlscalar.Map {
	return gqlscalar.Map(r.spec.PluginConfig)
}

func (r *CCIPSpecResolver) P2PKeyID() string {
	return r.spec.P2PKeyID
}
