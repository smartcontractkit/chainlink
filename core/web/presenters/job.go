package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// JobSpecType defines the the the spec type of the job
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
)

// DirectRequestSpec defines the spec details of a DirectRequest Job
type DirectRequestSpec struct {
	ContractAddress          models.EIP55Address `json:"contractAddress"`
	OnChainJobSpecID         string              `json:"onChainJobSpecID"`
	MinIncomingConfirmations clnull.Uint32       `json:"minIncomingConfirmations"`
	Initiator                string              `json:"initiator"`
	CreatedAt                time.Time           `json:"createdAt"`
	UpdatedAt                time.Time           `json:"updatedAt"`
}

// NewDirectRequestSpec initializes a new DirectRequestSpec from a
// job.DirectRequestSpec
func NewDirectRequestSpec(spec *job.DirectRequestSpec) *DirectRequestSpec {
	return &DirectRequestSpec{
		ContractAddress:          spec.ContractAddress,
		OnChainJobSpecID:         spec.OnChainJobSpecID.String(),
		MinIncomingConfirmations: spec.MinIncomingConfirmations,
		// This is hardcoded to runlog. When we support other intiators, we need
		// to change this
		Initiator: "runlog",
		CreatedAt: spec.CreatedAt,
		UpdatedAt: spec.UpdatedAt,
	}
}

// FluxMonitorSpec defines the spec details of a FluxMonitor Job
type FluxMonitorSpec struct {
	ContractAddress   models.EIP55Address `json:"contractAddress"`
	Precision         int32               `json:"precision"`
	Threshold         float32             `json:"threshold"`
	AbsoluteThreshold float32             `json:"absoluteThreshold"`
	PollTimerPeriod   string              `json:"pollTimerPeriod"`
	PollTimerDisabled bool                `json:"pollTimerDisabled"`
	IdleTimerPeriod   string              `json:"idleTimerPeriod"`
	IdleTimerDisabled bool                `json:"idleTimerDisabled"`
	MinPayment        *assets.Link        `json:"minPayment"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
}

// NewFluxMonitorSpec initializes a new DirectFluxMonitorSpec from a
// job.FluxMonitorSpec
func NewFluxMonitorSpec(spec *job.FluxMonitorSpec) *FluxMonitorSpec {
	return &FluxMonitorSpec{
		ContractAddress:   spec.ContractAddress,
		Precision:         spec.Precision,
		Threshold:         spec.Threshold,
		AbsoluteThreshold: spec.AbsoluteThreshold,
		PollTimerPeriod:   spec.PollTimerPeriod.String(),
		PollTimerDisabled: spec.PollTimerDisabled,
		IdleTimerPeriod:   spec.IdleTimerPeriod.String(),
		IdleTimerDisabled: spec.IdleTimerDisabled,
		MinPayment:        spec.MinPayment,
		CreatedAt:         spec.CreatedAt,
		UpdatedAt:         spec.UpdatedAt,
	}
}

// OffChainReportingSpec defines the spec details of a OffChainReporting Job
type OffChainReportingSpec struct {
	ContractAddress                        models.EIP55Address  `json:"contractAddress"`
	P2PPeerID                              *models.PeerID       `json:"p2pPeerID"`
	P2PBootstrapPeers                      pq.StringArray       `json:"p2pBootstrapPeers"`
	IsBootstrapPeer                        bool                 `json:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                *models.Sha256Hash   `json:"keyBundleID"`
	TransmitterAddress                     *models.EIP55Address `json:"transmitterAddress"`
	ObservationTimeout                     models.Interval      `json:"observationTimeout"`
	BlockchainTimeout                      models.Interval      `json:"blockchainTimeout"`
	ContractConfigTrackerSubscribeInterval models.Interval      `json:"contractConfigTrackerSubscribeInterval"`
	ContractConfigTrackerPollInterval      models.Interval      `json:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations            uint16               `json:"contractConfigConfirmations"`
	CreatedAt                              time.Time            `json:"createdAt"`
	UpdatedAt                              time.Time            `json:"updatedAt"`
}

// NewOffChainReportingSpec initializes a new OffChainReportingSpec from a
// job.OffchainReportingOracleSpec
func NewOffChainReportingSpec(spec *job.OffchainReportingOracleSpec) *OffChainReportingSpec {
	return &OffChainReportingSpec{
		ContractAddress:                        spec.ContractAddress,
		P2PPeerID:                              spec.P2PPeerID,
		P2PBootstrapPeers:                      spec.P2PBootstrapPeers,
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
	}
}

// PipelineSpec defines the spec details of the pipeline
type PipelineSpec struct {
	ID           int32  `json:"id"`
	DotDAGSource string `json:"dotDagSource"`
}

// NewPipelineSpec generates a new PipelineSpec from a pipeline.Spec
func NewPipelineSpec(spec *pipeline.Spec) PipelineSpec {
	return PipelineSpec{
		ID:           spec.ID,
		DotDAGSource: spec.DotDagSource,
	}
}

// KeeperSpec defines the spec details of a Keeper Job
type KeeperSpec struct {
	ContractAddress models.EIP55Address `json:"contractAddress"`
	FromAddress     models.EIP55Address `json:"fromAddress"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

// NewKeeperSpec generates a new KeeperSpec from a job.KeeperSpec
func NewKeeperSpec(spec *job.KeeperSpec) *KeeperSpec {
	return &KeeperSpec{
		ContractAddress: spec.ContractAddress,
		FromAddress:     spec.FromAddress,
		CreatedAt:       spec.CreatedAt,
		UpdatedAt:       spec.UpdatedAt,
	}
}

// WebhookSpec defines the spec details of a Webhook Job
type WebhookSpec struct {
	OnChainJobSpecID string    `json:"onChainJobSpecID"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// NewWebhookSpec generates a new WebhookSpec from a job.WebhookSpec
func NewWebhookSpec(spec *job.WebhookSpec) *WebhookSpec {
	return &WebhookSpec{
		OnChainJobSpecID: spec.OnChainJobSpecID.String(),
		CreatedAt:        spec.CreatedAt,
		UpdatedAt:        spec.UpdatedAt,
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
	CoordinatorAddress models.EIP55Address `toml:"coordinatorAddress"`
	PublicKey          secp256k1.PublicKey `toml:"publicKey"`
	Confirmations      uint32              `toml:"confirmations"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
}

func NewVRFSpec(spec *job.VRFSpec) *VRFSpec {
	return &VRFSpec{
		CoordinatorAddress: spec.CoordinatorAddress,
		PublicKey:          spec.PublicKey,
		Confirmations:      spec.Confirmations,
		CreatedAt:          spec.CreatedAt,
		UpdatedAt:          spec.UpdatedAt,
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
	Name                  string                 `json:"name"`
	Type                  JobSpecType            `json:"type"`
	SchemaVersion         uint32                 `json:"schemaVersion"`
	MaxTaskDuration       models.Interval        `json:"maxTaskDuration"`
	DirectRequestSpec     *DirectRequestSpec     `json:"directRequestSpec"`
	FluxMonitorSpec       *FluxMonitorSpec       `json:"fluxMonitorSpec"`
	CronSpec              *CronSpec              `json:"cronSpec"`
	OffChainReportingSpec *OffChainReportingSpec `json:"offChainReportingOracleSpec"`
	KeeperSpec            *KeeperSpec            `json:"keeperSpec"`
	VRFSpec               *VRFSpec               `json:"vrfSpec"`
	WebhookSpec           *WebhookSpec           `json:"webhookSpec"`
	PipelineSpec          PipelineSpec           `json:"pipelineSpec"`
	Errors                []JobError             `json:"errors"`
}

// NewJobResource initializes a new JSONAPI job resource
func NewJobResource(j job.Job) *JobResource {
	resource := &JobResource{
		JAID:            NewJAIDInt32(j.ID),
		Name:            j.Name.ValueOrZero(),
		Type:            JobSpecType(j.Type),
		SchemaVersion:   j.SchemaVersion,
		MaxTaskDuration: j.MaxTaskDuration,
		PipelineSpec:    NewPipelineSpec(j.PipelineSpec),
	}

	switch j.Type {
	case job.DirectRequest:
		resource.DirectRequestSpec = NewDirectRequestSpec(j.DirectRequestSpec)
	case job.FluxMonitor:
		resource.FluxMonitorSpec = NewFluxMonitorSpec(j.FluxMonitorSpec)
	case job.Cron:
		resource.CronSpec = NewCronSpec(j.CronSpec)
	case job.OffchainReporting:
		resource.OffChainReportingSpec = NewOffChainReportingSpec(j.OffchainreportingOracleSpec)
	case job.Keeper:
		resource.KeeperSpec = NewKeeperSpec(j.KeeperSpec)
	case job.VRF:
		resource.VRFSpec = NewVRFSpec(j.VRFSpec)
	case job.Webhook:
		resource.WebhookSpec = NewWebhookSpec(j.WebhookSpec)
	}

	jes := []JobError{}
	for _, e := range j.JobSpecErrors {
		jes = append(jes, NewJobError((e)))
	}
	resource.Errors = jes

	return resource
}

// NewJobResources initializes a slice of JSONAPI job resources
func NewJobResources(js []job.Job) []JobResource {
	rs := []JobResource{}

	for _, j := range js {
		rs = append(rs, *NewJobResource(j))
	}

	return rs
}

// GetName implements the api2go EntityNamer interface
func (r JobResource) GetName() string {
	return "jobs"
}
