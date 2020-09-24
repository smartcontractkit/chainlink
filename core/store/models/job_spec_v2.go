package models

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type (
	JobSpecV2 struct {
		ID                            int32 `gorm: "primary_key"`
		OffchainreportingOracleSpecID int32
		OffchainreportingOracleSpec   *OffchainReportingOracleSpec `gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
		PipelineSpecID                int32
	}

	OffchainReportingOracleSpec struct {
		ID                                int32                              `toml:"-"                 gorm:"primary_key"`
		ContractAddress                   common.Address                     `toml:"contractAddress"`
		P2PPeerID                         string                             `toml:"p2pPeerID"         gorm:"column:p2p_peer_id"`
		P2PBootstrapPeers                 OffchainReportingP2PBootstrapPeers `toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:jsonb"`
		OffchainreportingKeyBundleID      Sha256Hash                         `toml:"-"                 gorm:"type:bytea"`
		OffchainreportingKeyBundle        *OffchainReportingKeyBundle        `toml:"keyBundle"         gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
		MonitoringEndpoint                string                             `toml:"monitoringEndpoint"`
		TransmitterAddress                common.Address                     `toml:"transmitterAddress"`
		ObservationTimeout                Interval                           `toml:"observationTimeout" gorm:"type:bigint"`
		BlockchainTimeout                 Interval                           `toml:"blockchainTimeout" gorm:"type:bigint"`
		ContractConfigTrackerPollInterval Interval                           `toml:"contractConfigTrackerPollInterval" gorm:"type:bigint"`
		ContractConfigConfirmations       uint16                             `toml:"contractConfigConfirmations"`
		CreatedAt                         time.Time                          `toml:"-"`
		UpdatedAt                         time.Time                          `toml:"-"`
	}

	OffchainReportingP2PBootstrapPeer struct {
		PeerID    string `toml:"peerID"`
		Multiaddr string `toml:"multiAddr"`
	}

	OffchainReportingKeyBundle struct {
		ID                     Sha256Hash `toml:"-"                      gorm:"primary_key;type:bytea"`
		EncryptedPrivKeyBundle JSON       `toml:"encryptedPrivKeyBundle" gorm:"type:jsonb"`
		CreatedAt              time.Time  `toml:"-"`
	}
)

func (s *OffchainReportingOracleSpec) BeforeCreate() error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	s.OffchainreportingKeyBundleID = s.OffchainreportingKeyBundle.ID
	return nil
}

func (s *OffchainReportingOracleSpec) BeforeSave() error {
	s.UpdatedAt = time.Now()
	s.OffchainreportingKeyBundleID = s.OffchainreportingKeyBundle.ID
	return nil
}

func (JobSpecV2) TableName() string                   { return "jobs" }
func (OffchainReportingOracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
func (OffchainReportingKeyBundle) TableName() string  { return "offchainreporting_key_bundles" }

type OffchainReportingP2PBootstrapPeers []OffchainReportingP2PBootstrapPeer

func (p OffchainReportingP2PBootstrapPeers) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *OffchainReportingP2PBootstrapPeers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, p)
}

func (b *OffchainReportingKeyBundle) BeforeSave(scope *gorm.Scope) error {
	hash := sha256.New()
	copy(b.ID[:], hash.Sum(b.EncryptedPrivKeyBundle.Bytes()))
	scope.SetColumn("id", b.ID)
	scope.Set("gorm:insert_option", "ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id")
	return nil
}

func (b *OffchainReportingKeyBundle) BeforeCreate(scope *gorm.Scope) error {
	hash := sha256.New()
	copy(b.ID[:], hash.Sum(b.EncryptedPrivKeyBundle.Bytes()))
	b.CreatedAt = time.Now()
	scope.SetColumn("id", b.ID)
	scope.Set("gorm:insert_option", "ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id")
	return nil
}
