package offchainreporting

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	KeyBundle struct {
		ID                     models.Sha256Hash `gorm:"primary_key"`
		EncryptedPrivKeyBundle models.JSON
		CreatedAt              time.Time
	}

	OracleSpec struct {
		ID                 int32              `toml:"-" gorm:"primary_key"`
		ContractAddress    common.Address     `toml:"contractAddress"`
		P2PPeerID          string             `toml:"p2pPeerID"`
		P2PBootstrapPeers  []P2PBootstrapPeer `toml:"p2pBootstrapPeers"`
		KeyBundle          KeyBundle          `toml:"keyBundle"`
		MonitoringEndpoint string             `toml:"monitoringEndpoint"`
		TransmitterAddress common.Address     `toml:"transmitterAddress"`
		ObservationTimeout time.Duration      `toml:"observationTimeout"`
		LogLevel           ocrtypes.LogLevel  `toml:"logLevel,omitempty"`

		ObservationSource pipeline.TaskDAG `toml:"observationSource" gorm:"-"`
	}

	P2PBootstrapPeer struct {
		PeerID    string `toml:"peerID"`
		Multiaddr string `toml:"multiAddr"`
	}
)

func (KeyBundle) TableName() string  { return "offchainreporting_key_bundles" }
func (OracleSpec) TableName() string { return "offchainreporting_oracle_specs" }

const JobType pipeline.JobType = "offchainreporting"

// OracleSpec conforms to the pipeline.Spec interface
var _ pipeline.JobSpec = OracleSpec{}

func (spec OracleSpec) JobID() *models.ID {
	return spec.UUID
}

func (spec OracleSpec) JobType() pipeline.JobType {
	return JobType
}

func (spec OracleSpec) TaskDAG() TaskDAG {
	return spec.ObservationSource
}
