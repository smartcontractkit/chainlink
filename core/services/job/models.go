package job

import (
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/lib/pq"
	null "gopkg.in/guregu/null.v4"
)

type (
	IDEmbed struct {
		ID int32 `json:"-" toml:"-"                 gorm:"primary_key"`
	}

	JobSpecV2 struct {
		IDEmbed
		OffchainreportingOracleSpecID *int32                       `json:"-"`
		OffchainreportingOracleSpec   *OffchainReportingOracleSpec `json:"offChainReportingOracleSpec" gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
		EthRequestEventSpecID         *int32                       `json:"-"`
		EthRequestEventSpec           *EthRequestEventSpec         `json:"ethRequestEventSpec" gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
		PipelineSpecID                int32                        `json:"-"`
		PipelineSpec                  *PipelineSpec                `json:"pipelineSpec"`
		JobSpecErrors                 []JobSpecErrorV2             `json:"errors" gorm:"foreignKey:JobID"`
		Type                          string                       `json:"type"`
		SchemaVersion                 uint32                       `json:"schemaVersion"`
		Name                          null.String                  `json:"name"`
		MaxTaskDuration               models.Interval              `json:"maxTaskDuration"`
	}

	JobSpecErrorV2 struct {
		ID          int64     `json:"id" gorm:"primary_key"`
		JobID       int32     `json:"-"`
		Description string    `json:"description"`
		Occurrences uint      `json:"occurrences"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	PipelineRun struct {
		ID int64 `json:"-" gorm:"primary_key"`
	}

	PipelineSpec struct {
		IDEmbed
		DotDagSource string    `json:"dotDagSource"`
		CreatedAt    time.Time `json:"-"`
	}

	// TODO: remove pointers when upgrading to gormv2
	// which has https://github.com/go-gorm/gorm/issues/2748 fixed.
	OffchainReportingOracleSpec struct {
		IDEmbed
		ContractAddress                        models.EIP55Address  `json:"contractAddress" toml:"contractAddress"`
		P2PPeerID                              *models.PeerID       `json:"p2pPeerID" toml:"p2pPeerID" gorm:"column:p2p_peer_id;default:null"`
		P2PBootstrapPeers                      pq.StringArray       `json:"p2pBootstrapPeers" toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:text[]"`
		IsBootstrapPeer                        bool                 `json:"isBootstrapPeer" toml:"isBootstrapPeer"`
		EncryptedOCRKeyBundleID                *models.Sha256Hash   `json:"keyBundleID" toml:"keyBundleID"                 gorm:"type:bytea"`
		MonitoringEndpoint                     string               `json:"monitoringEndpoint" toml:"monitoringEndpoint"`
		TransmitterAddress                     *models.EIP55Address `json:"transmitterAddress" toml:"transmitterAddress"`
		ObservationTimeout                     models.Interval      `json:"observationTimeout" toml:"observationTimeout" gorm:"type:bigint;default:null"`
		BlockchainTimeout                      models.Interval      `json:"blockchainTimeout" toml:"blockchainTimeout" gorm:"type:bigint;default:null"`
		ContractConfigTrackerSubscribeInterval models.Interval      `json:"contractConfigTrackerSubscribeInterval" toml:"contractConfigTrackerSubscribeInterval" gorm:"default:null"`
		ContractConfigTrackerPollInterval      models.Interval      `json:"contractConfigTrackerPollInterval" toml:"contractConfigTrackerPollInterval" gorm:"type:bigint;default:null"`
		ContractConfigConfirmations            uint16               `json:"contractConfigConfirmations" toml:"contractConfigConfirmations" gorm:"default:null"`
		CreatedAt                              time.Time            `json:"createdAt" toml:"-"`
		UpdatedAt                              time.Time            `json:"updatedAt" toml:"-"`
	}

	EthRequestEventSpec struct {
		IDEmbed
		ContractAddress models.EIP55Address `json:"contractAddress" toml:"contractAddress"`
		CreatedAt       time.Time           `json:"createdAt" toml:"-"`
		UpdatedAt       time.Time           `json:"updatedAt" toml:"-"`
	}
)

const (
	EthRequestEventJobType = "ethrequestevent"
)

func (id IDEmbed) GetID() string {
	return fmt.Sprintf("%v", id.ID)
}

func (id *IDEmbed) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	id.ID = int32(ID)
	return nil
}

func (s OffchainReportingOracleSpec) GetID() string {
	return fmt.Sprintf("%v", s.ID)
}

func (s *OffchainReportingOracleSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

func (pr PipelineRun) GetID() string {
	return fmt.Sprintf("%v", pr.ID)
}

func (pr *PipelineRun) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	pr.ID = int64(ID)
	return nil
}

func (s *OffchainReportingOracleSpec) BeforeCreate() error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *OffchainReportingOracleSpec) BeforeSave() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (JobSpecV2) TableName() string                   { return "jobs" }
func (JobSpecErrorV2) TableName() string              { return "job_spec_errors_v2" }
func (OffchainReportingOracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
func (EthRequestEventSpec) TableName() string         { return "eth_request_event_specs" }
