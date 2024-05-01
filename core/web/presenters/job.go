package presenters

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// JobSpecType defines the spec type of the job
type JobSpecType string

func (t JobSpecType) String() string {
	return string(t)
}

const (
	DirectRequestJobSpec     JobSpecType = "directrequest"
	FluxMonitorJobSpec       JobSpecType = "fluxmonitor"
	OffChainReportingJobSpec JobSpecType = "offchainreporting"
	KeeperJobSpec            JobSpecType = "keeper"
	CronJobSpec              JobSpecType = "cron"
	VRFJobSpec               JobSpecType = "vrf"
	WebhookJobSpec           JobSpecType = "webhook"
	BlockhashStoreJobSpec    JobSpecType = "blockhashstore"
	BlockHeaderFeederJobSpec JobSpecType = "blockheaderfeeder"
	BootstrapJobSpec         JobSpecType = "bootstrap"
	GatewayJobSpec           JobSpecType = "gateway"
	WorkflowJobSpec          JobSpecType = "workflow"
)

// DirectRequestSpec defines the spec details of a DirectRequest Job
type DirectRequestSpec struct {
	ContractAddress          types.EIP55Address       `json:"contractAddress"`
	MinIncomingConfirmations clnull.Uint32            `json:"minIncomingConfirmations"`
	MinContractPayment       *commonassets.Link       `json:"minContractPaymentLinkJuels"`
	Requesters               models.AddressCollection `json:"requesters"`
	Initiator                string                   `json:"initiator"`
	CreatedAt                time.Time                `json:"createdAt"`
	UpdatedAt                time.Time                `json:"updatedAt"`
	EVMChainID               *big.Big                 `json:"evmChainID"`
}

// NewDirectRequestSpec initializes a new DirectRequestSpec from a
// job.DirectRequestSpec
func NewDirectRequestSpec(spec *job.DirectRequestSpec) *DirectRequestSpec {
	return &DirectRequestSpec{
		ContractAddress:          spec.ContractAddress,
		MinIncomingConfirmations: spec.MinIncomingConfirmations,
		MinContractPayment:       spec.MinContractPayment,
		Requesters:               spec.Requesters,
		// This is hardcoded to runlog. When we support other initiators, we need
		// to change this
		Initiator:  "runlog",
		CreatedAt:  spec.CreatedAt,
		UpdatedAt:  spec.UpdatedAt,
		EVMChainID: spec.EVMChainID,
	}
}

// FluxMonitorSpec defines the spec details of a FluxMonitor Job
type FluxMonitorSpec struct {
	ContractAddress     types.EIP55Address `json:"contractAddress"`
	Threshold           float32            `json:"threshold"`
	AbsoluteThreshold   float32            `json:"absoluteThreshold"`
	PollTimerPeriod     string             `json:"pollTimerPeriod"`
	PollTimerDisabled   bool               `json:"pollTimerDisabled"`
	IdleTimerPeriod     string             `json:"idleTimerPeriod"`
	IdleTimerDisabled   bool               `json:"idleTimerDisabled"`
	DrumbeatEnabled     bool               `json:"drumbeatEnabled"`
	DrumbeatSchedule    *string            `json:"drumbeatSchedule"`
	DrumbeatRandomDelay *string            `json:"drumbeatRandomDelay"`
	MinPayment          *commonassets.Link `json:"minPayment"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	EVMChainID          *big.Big           `json:"evmChainID"`
}

// NewFluxMonitorSpec initializes a new DirectFluxMonitorSpec from a
// job.FluxMonitorSpec
func NewFluxMonitorSpec(spec *job.FluxMonitorSpec) *FluxMonitorSpec {
	var drumbeatSchedulePtr *string
	if spec.DrumbeatEnabled {
		drumbeatSchedulePtr = &spec.DrumbeatSchedule
	}
	var drumbeatRandomDelayPtr *string
	if spec.DrumbeatRandomDelay > 0 {
		drumbeatRandomDelay := spec.DrumbeatRandomDelay.String()
		drumbeatRandomDelayPtr = &drumbeatRandomDelay
	}
	return &FluxMonitorSpec{
		ContractAddress:     spec.ContractAddress,
		Threshold:           float32(spec.Threshold),
		AbsoluteThreshold:   float32(spec.AbsoluteThreshold),
		PollTimerPeriod:     spec.PollTimerPeriod.String(),
		PollTimerDisabled:   spec.PollTimerDisabled,
		IdleTimerPeriod:     spec.IdleTimerPeriod.String(),
		IdleTimerDisabled:   spec.IdleTimerDisabled,
		DrumbeatEnabled:     spec.DrumbeatEnabled,
		DrumbeatSchedule:    drumbeatSchedulePtr,
		DrumbeatRandomDelay: drumbeatRandomDelayPtr,
		MinPayment:          spec.MinPayment,
		CreatedAt:           spec.CreatedAt,
		UpdatedAt:           spec.UpdatedAt,
		EVMChainID:          spec.EVMChainID,
	}
}

// OffChainReportingSpec defines the spec details of a OffChainReporting Job
type OffChainReportingSpec struct {
	ContractAddress                        types.EIP55Address  `json:"contractAddress"`
	P2PV2Bootstrappers                     pq.StringArray      `json:"p2pv2Bootstrappers"`
	IsBootstrapPeer                        bool                `json:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                *models.Sha256Hash  `json:"keyBundleID"`
	TransmitterAddress                     *types.EIP55Address `json:"transmitterAddress"`
	ObservationTimeout                     models.Interval     `json:"observationTimeout"`
	BlockchainTimeout                      models.Interval     `json:"blockchainTimeout"`
	ContractConfigTrackerSubscribeInterval models.Interval     `json:"contractConfigTrackerSubscribeInterval"`
	ContractConfigTrackerPollInterval      models.Interval     `json:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations            uint16              `json:"contractConfigConfirmations"`
	CreatedAt                              time.Time           `json:"createdAt"`
	UpdatedAt                              time.Time           `json:"updatedAt"`
	EVMChainID                             *big.Big            `json:"evmChainID"`
	DatabaseTimeout                        *models.Interval    `json:"databaseTimeout"`
	ObservationGracePeriod                 *models.Interval    `json:"observationGracePeriod"`
	ContractTransmitterTransmitTimeout     *models.Interval    `json:"contractTransmitterTransmitTimeout"`
	CollectTelemetry                       bool                `json:"collectTelemetry,omitempty"`
}

// NewOffChainReportingSpec initializes a new OffChainReportingSpec from a
// job.OCROracleSpec
func NewOffChainReportingSpec(spec *job.OCROracleSpec) *OffChainReportingSpec {
	return &OffChainReportingSpec{
		ContractAddress:                        spec.ContractAddress,
		P2PV2Bootstrappers:                     spec.P2PV2Bootstrappers,
		IsBootstrapPeer:                        spec.IsBootstrapPeer,
		EncryptedOCRKeyBundleID:                spec.EncryptedOCRKeyBundleID,
		TransmitterAddress:                     spec.TransmitterAddress,
		ObservationTimeout:                     spec.ObservationTimeout,
		BlockchainTimeout:                      spec.BlockchainTimeout,
		ContractConfigTrackerSubscribeInterval: spec.ContractConfigTrackerSubscribeInterval,
		ContractConfigTrackerPollInterval:      spec.ContractConfigTrackerPollInterval,
		ContractConfigConfirmations:            spec.ContractConfigConfirmations,
		CreatedAt:                              spec.CreatedAt,
		UpdatedAt:                              spec.UpdatedAt,
		EVMChainID:                             spec.EVMChainID,
		DatabaseTimeout:                        spec.DatabaseTimeout,
		ObservationGracePeriod:                 spec.ObservationGracePeriod,
		ContractTransmitterTransmitTimeout:     spec.ContractTransmitterTransmitTimeout,
		CollectTelemetry:                       spec.CaptureEATelemetry,
	}
}

// OffChainReporting2Spec defines the spec details of a OffChainReporting2 Job
type OffChainReporting2Spec struct {
	ContractID                        string                 `json:"contractID"`
	Relay                             string                 `json:"relay"` // RelayID.Network
	RelayConfig                       map[string]interface{} `json:"relayConfig"`
	P2PV2Bootstrappers                pq.StringArray         `json:"p2pv2Bootstrappers"`
	OCRKeyBundleID                    null.String            `json:"ocrKeyBundleID"`
	TransmitterID                     null.String            `json:"transmitterID"`
	ObservationTimeout                models.Interval        `json:"observationTimeout"`
	BlockchainTimeout                 models.Interval        `json:"blockchainTimeout"`
	ContractConfigTrackerPollInterval models.Interval        `json:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations       uint16                 `json:"contractConfigConfirmations"`
	CreatedAt                         time.Time              `json:"createdAt"`
	UpdatedAt                         time.Time              `json:"updatedAt"`
	CollectTelemetry                  bool                   `json:"collectTelemetry"`
}

// NewOffChainReporting2Spec initializes a new OffChainReportingSpec from a
// job.OCR2OracleSpec
func NewOffChainReporting2Spec(spec *job.OCR2OracleSpec) *OffChainReporting2Spec {
	return &OffChainReporting2Spec{
		ContractID:                        spec.ContractID,
		Relay:                             spec.Relay,
		RelayConfig:                       spec.RelayConfig,
		P2PV2Bootstrappers:                spec.P2PV2Bootstrappers,
		OCRKeyBundleID:                    spec.OCRKeyBundleID,
		TransmitterID:                     spec.TransmitterID,
		BlockchainTimeout:                 spec.BlockchainTimeout,
		ContractConfigTrackerPollInterval: spec.ContractConfigTrackerPollInterval,
		ContractConfigConfirmations:       spec.ContractConfigConfirmations,
		CreatedAt:                         spec.CreatedAt,
		UpdatedAt:                         spec.UpdatedAt,
		CollectTelemetry:                  spec.CaptureEATelemetry,
	}
}

// PipelineSpec defines the spec details of the pipeline
type PipelineSpec struct {
	ID           int32  `json:"id"`
	JobID        int32  `json:"jobID"`
	DotDAGSource string `json:"dotDagSource"`
}

// NewPipelineSpec generates a new PipelineSpec from a pipeline.Spec
func NewPipelineSpec(spec *pipeline.Spec) PipelineSpec {
	return PipelineSpec{
		ID:           spec.ID,
		JobID:        spec.JobID,
		DotDAGSource: spec.DotDagSource,
	}
}

// KeeperSpec defines the spec details of a Keeper Job
type KeeperSpec struct {
	ContractAddress types.EIP55Address `json:"contractAddress"`
	FromAddress     types.EIP55Address `json:"fromAddress"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
	EVMChainID      *big.Big           `json:"evmChainID"`
}

// NewKeeperSpec generates a new KeeperSpec from a job.KeeperSpec
func NewKeeperSpec(spec *job.KeeperSpec) *KeeperSpec {
	return &KeeperSpec{
		ContractAddress: spec.ContractAddress,
		FromAddress:     spec.FromAddress,
		CreatedAt:       spec.CreatedAt,
		UpdatedAt:       spec.UpdatedAt,
		EVMChainID:      spec.EVMChainID,
	}
}

// WebhookSpec defines the spec details of a Webhook Job
type WebhookSpec struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewWebhookSpec generates a new WebhookSpec from a job.WebhookSpec
func NewWebhookSpec(spec *job.WebhookSpec) *WebhookSpec {
	return &WebhookSpec{
		CreatedAt: spec.CreatedAt,
		UpdatedAt: spec.UpdatedAt,
	}
}

// CronSpec defines the spec details of a Cron Job
type CronSpec struct {
	CronSchedule string    `json:"schedule" tom:"schedule"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NewCronSpec generates a new CronSpec from a job.CronSpec
func NewCronSpec(spec *job.CronSpec) *CronSpec {
	return &CronSpec{
		CronSchedule: spec.CronSchedule,
		CreatedAt:    spec.CreatedAt,
		UpdatedAt:    spec.UpdatedAt,
	}
}

type VRFSpec struct {
	BatchCoordinatorAddress       *types.EIP55Address   `json:"batchCoordinatorAddress"`
	BatchFulfillmentEnabled       bool                  `json:"batchFulfillmentEnabled"`
	CustomRevertsPipelineEnabled  *bool                 `json:"customRevertsPipelineEnabled,omitempty"`
	BatchFulfillmentGasMultiplier float64               `json:"batchFulfillmentGasMultiplier"`
	CoordinatorAddress            types.EIP55Address    `json:"coordinatorAddress"`
	PublicKey                     secp256k1.PublicKey   `json:"publicKey"`
	FromAddresses                 []types.EIP55Address  `json:"fromAddresses"`
	PollPeriod                    commonconfig.Duration `json:"pollPeriod"`
	MinIncomingConfirmations      uint32                `json:"confirmations"`
	CreatedAt                     time.Time             `json:"createdAt"`
	UpdatedAt                     time.Time             `json:"updatedAt"`
	EVMChainID                    *big.Big              `json:"evmChainID"`
	ChunkSize                     uint32                `json:"chunkSize"`
	RequestTimeout                commonconfig.Duration `json:"requestTimeout"`
	BackoffInitialDelay           commonconfig.Duration `json:"backoffInitialDelay"`
	BackoffMaxDelay               commonconfig.Duration `json:"backoffMaxDelay"`
	GasLanePrice                  *assets.Wei           `json:"gasLanePrice"`
	RequestedConfsDelay           int64                 `json:"requestedConfsDelay"`
	VRFOwnerAddress               *types.EIP55Address   `json:"vrfOwnerAddress,omitempty"`
}

func NewVRFSpec(spec *job.VRFSpec) *VRFSpec {
	return &VRFSpec{
		BatchCoordinatorAddress:       spec.BatchCoordinatorAddress,
		BatchFulfillmentEnabled:       spec.BatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: float64(spec.BatchFulfillmentGasMultiplier),
		CustomRevertsPipelineEnabled:  &spec.CustomRevertsPipelineEnabled,
		CoordinatorAddress:            spec.CoordinatorAddress,
		PublicKey:                     spec.PublicKey,
		FromAddresses:                 spec.FromAddresses,
		PollPeriod:                    *commonconfig.MustNewDuration(spec.PollPeriod),
		MinIncomingConfirmations:      spec.MinIncomingConfirmations,
		CreatedAt:                     spec.CreatedAt,
		UpdatedAt:                     spec.UpdatedAt,
		EVMChainID:                    spec.EVMChainID,
		ChunkSize:                     spec.ChunkSize,
		RequestTimeout:                *commonconfig.MustNewDuration(spec.RequestTimeout),
		BackoffInitialDelay:           *commonconfig.MustNewDuration(spec.BackoffInitialDelay),
		BackoffMaxDelay:               *commonconfig.MustNewDuration(spec.BackoffMaxDelay),
		GasLanePrice:                  spec.GasLanePrice,
		RequestedConfsDelay:           spec.RequestedConfsDelay,
		VRFOwnerAddress:               spec.VRFOwnerAddress,
	}
}

// BlockhashStoreSpec defines the job parameters for a blockhash store feeder job.
type BlockhashStoreSpec struct {
	CoordinatorV1Address           *types.EIP55Address  `json:"coordinatorV1Address"`
	CoordinatorV2Address           *types.EIP55Address  `json:"coordinatorV2Address"`
	CoordinatorV2PlusAddress       *types.EIP55Address  `json:"coordinatorV2PlusAddress"`
	WaitBlocks                     int32                `json:"waitBlocks"`
	LookbackBlocks                 int32                `json:"lookbackBlocks"`
	HeartbeatPeriod                time.Duration        `json:"heartbeatPeriod"`
	BlockhashStoreAddress          types.EIP55Address   `json:"blockhashStoreAddress"`
	TrustedBlockhashStoreAddress   *types.EIP55Address  `json:"trustedBlockhashStoreAddress"`
	TrustedBlockhashStoreBatchSize int32                `json:"trustedBlockhashStoreBatchSize"`
	PollPeriod                     time.Duration        `json:"pollPeriod"`
	RunTimeout                     time.Duration        `json:"runTimeout"`
	EVMChainID                     *big.Big             `json:"evmChainID"`
	FromAddresses                  []types.EIP55Address `json:"fromAddresses"`
	CreatedAt                      time.Time            `json:"createdAt"`
	UpdatedAt                      time.Time            `json:"updatedAt"`
}

// NewBlockhashStoreSpec creates a new BlockhashStoreSpec for the given parameters.
func NewBlockhashStoreSpec(spec *job.BlockhashStoreSpec) *BlockhashStoreSpec {
	return &BlockhashStoreSpec{
		CoordinatorV1Address:           spec.CoordinatorV1Address,
		CoordinatorV2Address:           spec.CoordinatorV2Address,
		CoordinatorV2PlusAddress:       spec.CoordinatorV2PlusAddress,
		WaitBlocks:                     spec.WaitBlocks,
		LookbackBlocks:                 spec.LookbackBlocks,
		HeartbeatPeriod:                spec.HeartbeatPeriod,
		BlockhashStoreAddress:          spec.BlockhashStoreAddress,
		TrustedBlockhashStoreAddress:   spec.TrustedBlockhashStoreAddress,
		TrustedBlockhashStoreBatchSize: spec.TrustedBlockhashStoreBatchSize,
		PollPeriod:                     spec.PollPeriod,
		RunTimeout:                     spec.RunTimeout,
		EVMChainID:                     spec.EVMChainID,
		FromAddresses:                  spec.FromAddresses,
	}
}

// BlockHeaderFeederSpec defines the job parameters for a blcok header feeder job.
type BlockHeaderFeederSpec struct {
	CoordinatorV1Address       *types.EIP55Address  `json:"coordinatorV1Address"`
	CoordinatorV2Address       *types.EIP55Address  `json:"coordinatorV2Address"`
	CoordinatorV2PlusAddress   *types.EIP55Address  `json:"coordinatorV2PlusAddress"`
	WaitBlocks                 int32                `json:"waitBlocks"`
	LookbackBlocks             int32                `json:"lookbackBlocks"`
	BlockhashStoreAddress      types.EIP55Address   `json:"blockhashStoreAddress"`
	BatchBlockhashStoreAddress types.EIP55Address   `json:"batchBlockhashStoreAddress"`
	PollPeriod                 time.Duration        `json:"pollPeriod"`
	RunTimeout                 time.Duration        `json:"runTimeout"`
	EVMChainID                 *big.Big             `json:"evmChainID"`
	FromAddresses              []types.EIP55Address `json:"fromAddresses"`
	GetBlockhashesBatchSize    uint16               `json:"getBlockhashesBatchSize"`
	StoreBlockhashesBatchSize  uint16               `json:"storeBlockhashesBatchSize"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewBlockHeaderFeederSpec creates a new BlockHeaderFeederSpec for the given parameters.
func NewBlockHeaderFeederSpec(spec *job.BlockHeaderFeederSpec) *BlockHeaderFeederSpec {
	return &BlockHeaderFeederSpec{
		CoordinatorV1Address:       spec.CoordinatorV1Address,
		CoordinatorV2Address:       spec.CoordinatorV2Address,
		CoordinatorV2PlusAddress:   spec.CoordinatorV2PlusAddress,
		WaitBlocks:                 spec.WaitBlocks,
		LookbackBlocks:             spec.LookbackBlocks,
		BlockhashStoreAddress:      spec.BlockhashStoreAddress,
		BatchBlockhashStoreAddress: spec.BatchBlockhashStoreAddress,
		PollPeriod:                 spec.PollPeriod,
		RunTimeout:                 spec.RunTimeout,
		EVMChainID:                 spec.EVMChainID,
		FromAddresses:              spec.FromAddresses,
		GetBlockhashesBatchSize:    spec.GetBlockhashesBatchSize,
		StoreBlockhashesBatchSize:  spec.StoreBlockhashesBatchSize,
	}
}

// BootstrapSpec defines the spec details of a BootstrapSpec Job
type BootstrapSpec struct {
	ContractID                             string                 `json:"contractID"`
	Relay                                  string                 `json:"relay"` // RelayID.Network
	RelayConfig                            map[string]interface{} `json:"relayConfig"`
	BlockchainTimeout                      models.Interval        `json:"blockchainTimeout"`
	ContractConfigTrackerSubscribeInterval models.Interval        `json:"contractConfigTrackerSubscribeInterval"`
	ContractConfigTrackerPollInterval      models.Interval        `json:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations            uint16                 `json:"contractConfigConfirmations"`
	CreatedAt                              time.Time              `json:"createdAt"`
	UpdatedAt                              time.Time              `json:"updatedAt"`
}

// NewBootstrapSpec initializes a new BootstrapSpec from a job.BootstrapSpec
func NewBootstrapSpec(spec *job.BootstrapSpec) *BootstrapSpec {
	return &BootstrapSpec{
		ContractID:                        spec.ContractID,
		Relay:                             spec.Relay,
		RelayConfig:                       spec.RelayConfig,
		BlockchainTimeout:                 spec.BlockchainTimeout,
		ContractConfigTrackerPollInterval: spec.ContractConfigTrackerPollInterval,
		ContractConfigConfirmations:       spec.ContractConfigConfirmations,
		CreatedAt:                         spec.CreatedAt,
		UpdatedAt:                         spec.UpdatedAt,
	}
}

type GatewaySpec struct {
	GatewayConfig map[string]interface{} `json:"gatewayConfig"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

func NewGatewaySpec(spec *job.GatewaySpec) *GatewaySpec {
	return &GatewaySpec{
		GatewayConfig: spec.GatewayConfig,
		CreatedAt:     spec.CreatedAt,
		UpdatedAt:     spec.UpdatedAt,
	}
}

type WorkflowSpec struct {
	Workflow      string    `json:"workflow"`
	WorkflowID    string    `json:"workflowId"`
	WorkflowOwner string    `json:"workflowOwner"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewWorkflowSpec(spec *job.WorkflowSpec) *WorkflowSpec {
	return &WorkflowSpec{
		Workflow:      spec.Workflow,
		WorkflowID:    spec.WorkflowID,
		WorkflowOwner: spec.WorkflowOwner,
		CreatedAt:     spec.CreatedAt,
		UpdatedAt:     spec.UpdatedAt,
	}
}

// JobError represents errors on the job
type JobError struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Occurrences uint      `json:"occurrences"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewJobError(e job.SpecError) JobError {
	return JobError{
		ID:          e.ID,
		Description: e.Description,
		Occurrences: e.Occurrences,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// JobResource represents a JobResource
type JobResource struct {
	JAID
	Name                   string                  `json:"name"`
	StreamID               *uint32                 `json:"streamID,omitempty"`
	Type                   JobSpecType             `json:"type"`
	SchemaVersion          uint32                  `json:"schemaVersion"`
	GasLimit               clnull.Uint32           `json:"gasLimit"`
	ForwardingAllowed      bool                    `json:"forwardingAllowed"`
	MaxTaskDuration        models.Interval         `json:"maxTaskDuration"`
	ExternalJobID          uuid.UUID               `json:"externalJobID"`
	DirectRequestSpec      *DirectRequestSpec      `json:"directRequestSpec"`
	FluxMonitorSpec        *FluxMonitorSpec        `json:"fluxMonitorSpec"`
	CronSpec               *CronSpec               `json:"cronSpec"`
	OffChainReportingSpec  *OffChainReportingSpec  `json:"offChainReportingOracleSpec"`
	OffChainReporting2Spec *OffChainReporting2Spec `json:"offChainReporting2OracleSpec"`
	KeeperSpec             *KeeperSpec             `json:"keeperSpec"`
	VRFSpec                *VRFSpec                `json:"vrfSpec"`
	WebhookSpec            *WebhookSpec            `json:"webhookSpec"`
	BlockhashStoreSpec     *BlockhashStoreSpec     `json:"blockhashStoreSpec"`
	BlockHeaderFeederSpec  *BlockHeaderFeederSpec  `json:"blockHeaderFeederSpec"`
	BootstrapSpec          *BootstrapSpec          `json:"bootstrapSpec"`
	GatewaySpec            *GatewaySpec            `json:"gatewaySpec"`
	WorkflowSpec           *WorkflowSpec           `json:"workflowSpec"`
	PipelineSpec           PipelineSpec            `json:"pipelineSpec"`
	Errors                 []JobError              `json:"errors"`
}

// NewJobResource initializes a new JSONAPI job resource
func NewJobResource(j job.Job) *JobResource {
	resource := &JobResource{
		JAID:              NewJAIDInt32(j.ID),
		Name:              j.Name.ValueOrZero(),
		StreamID:          j.StreamID,
		Type:              JobSpecType(j.Type),
		SchemaVersion:     j.SchemaVersion,
		GasLimit:          j.GasLimit,
		ForwardingAllowed: j.ForwardingAllowed,
		MaxTaskDuration:   j.MaxTaskDuration,
		PipelineSpec:      NewPipelineSpec(j.PipelineSpec),
		ExternalJobID:     j.ExternalJobID,
	}

	switch j.Type {
	case job.DirectRequest:
		resource.DirectRequestSpec = NewDirectRequestSpec(j.DirectRequestSpec)
	case job.FluxMonitor:
		resource.FluxMonitorSpec = NewFluxMonitorSpec(j.FluxMonitorSpec)
	case job.Cron:
		resource.CronSpec = NewCronSpec(j.CronSpec)
	case job.OffchainReporting:
		resource.OffChainReportingSpec = NewOffChainReportingSpec(j.OCROracleSpec)
	case job.OffchainReporting2:
		resource.OffChainReporting2Spec = NewOffChainReporting2Spec(j.OCR2OracleSpec)
	case job.Keeper:
		resource.KeeperSpec = NewKeeperSpec(j.KeeperSpec)
	case job.VRF:
		resource.VRFSpec = NewVRFSpec(j.VRFSpec)
	case job.Webhook:
		resource.WebhookSpec = NewWebhookSpec(j.WebhookSpec)
	case job.BlockhashStore:
		resource.BlockhashStoreSpec = NewBlockhashStoreSpec(j.BlockhashStoreSpec)
	case job.BlockHeaderFeeder:
		resource.BlockHeaderFeederSpec = NewBlockHeaderFeederSpec(j.BlockHeaderFeederSpec)
	case job.Bootstrap:
		resource.BootstrapSpec = NewBootstrapSpec(j.BootstrapSpec)
	case job.Gateway:
		resource.GatewaySpec = NewGatewaySpec(j.GatewaySpec)
	case job.Stream:
		// no spec; nothing to do
	case job.Workflow:
		resource.WorkflowSpec = NewWorkflowSpec(j.WorkflowSpec)
	case job.LegacyGasStationServer, job.LegacyGasStationSidecar:
		// unsupported
	}

	jes := []JobError{}
	for _, e := range j.JobSpecErrors {
		jes = append(jes, NewJobError((e)))
	}
	resource.Errors = jes

	return resource
}

// GetName implements the api2go EntityNamer interface
func (r JobResource) GetName() string {
	return "jobs"
}
