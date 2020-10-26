package models

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

type (
	JobSpecV2 struct {
		ID                            int32 `gorm:"primary_key"`
		OffchainreportingOracleSpecID int32
		OffchainreportingOracleSpec   *OffchainReportingOracleSpec `gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
		PipelineSpecID                int32
	}

	OffchainReportingOracleSpec struct {
		ID                                     int32          `toml:"-"                 gorm:"primary_key"`
		ContractAddress                        EIP55Address   `toml:"contractAddress"`
		P2PPeerID                              PeerID         `toml:"p2pPeerID"         gorm:"column:p2p_peer_id"`
		P2PBootstrapPeers                      pq.StringArray `toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:text[]"`
		Name                                   string         `toml:"name" gorm:"column:name;type:text"`
		IsBootstrapPeer                        bool           `toml:"isBootstrapPeer"`
		EncryptedOCRKeyBundleID                Sha256Hash     `toml:"keyBundleID"                 gorm:"type:bytea"`
		MonitoringEndpoint                     string         `toml:"monitoringEndpoint"`
		TransmitterAddress                     EIP55Address   `toml:"transmitterAddress"`
		ObservationTimeout                     Interval       `toml:"observationTimeout" gorm:"type:bigint"`
		BlockchainTimeout                      Interval       `toml:"blockchainTimeout" gorm:"type:bigint"`
		ContractConfigTrackerSubscribeInterval Interval       `toml:"contractConfigTrackerSubscribeInterval"`
		ContractConfigTrackerPollInterval      Interval       `toml:"contractConfigTrackerPollInterval" gorm:"type:bigint"`
		ContractConfigConfirmations            uint16         `toml:"contractConfigConfirmations"`
		CreatedAt                              time.Time      `toml:"-"`
		UpdatedAt                              time.Time      `toml:"-"`
	}

	PeerID peer.ID
)

func (p PeerID) String() string {
	return peer.ID(p).String()
}

func (p *PeerID) UnmarshalText(bs []byte) error {
	peerID, err := peer.Decode(string(bs))
	if err != nil {
		return errors.Wrapf(err, `PeerID#UnmarshalText("%v")`, string(bs))
	}
	*p = PeerID(peerID)
	return nil
}

func (p *PeerID) Scan(value interface{}) error {
	s, is := value.(string)
	if !is {
		return errors.Errorf("PeerID#Scan got %T, expected string", value)
	}
	*p = PeerID("")
	return p.UnmarshalText([]byte(s))
}

func (p PeerID) Value() (driver.Value, error) {
	return peer.Encode(peer.ID(p)), nil
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
func (OffchainReportingOracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
