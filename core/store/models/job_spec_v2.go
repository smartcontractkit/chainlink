package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

type (
	JobSpecV2 struct {
		ID                            int32                        `json:"-" gorm:"primary_key"`
		OffchainreportingOracleSpecID int32                        `json:"-"`
		OffchainreportingOracleSpec   *OffchainReportingOracleSpec `json:"offChainReportingOracleSpec" gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
		PipelineSpecID                int32                        `json:"-"`
	}

	OffchainReportingOracleSpec struct {
		ID                                     int32          `json:"-" toml:"-"                 gorm:"primary_key"`
		ContractAddress                        EIP55Address   `json:"contractAddress" toml:"contractAddress"`
		P2PPeerID                              PeerID         `json:"p2pPeerID" toml:"p2pPeerID"         gorm:"column:p2p_peer_id"`
		P2PBootstrapPeers                      pq.StringArray `json:"p2pBootstrapPeers" toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:text[]"`
		IsBootstrapPeer                        bool           `json:"isBootstrapPeer" toml:"isBootstrapPeer"`
		EncryptedOCRKeyBundleID                Sha256Hash     `json:"keyBundleID" toml:"keyBundleID"                 gorm:"type:bytea"`
		MonitoringEndpoint                     string         `json:"monitoringEndpoint" toml:"monitoringEndpoint"`
		TransmitterAddress                     EIP55Address   `json:"transmitterAddress" toml:"transmitterAddress"`
		ObservationTimeout                     Interval       `json:"observationTimeout" toml:"observationTimeout" gorm:"type:bigint"`
		BlockchainTimeout                      Interval       `json:"blockchainTimeout" toml:"blockchainTimeout" gorm:"type:bigint"`
		ContractConfigTrackerSubscribeInterval Interval       `json:"contractConfigTrackerSubscribeInterval" toml:"contractConfigTrackerSubscribeInterval"`
		ContractConfigTrackerPollInterval      Interval       `json:"contractConfigTrackerPollInterval" toml:"contractConfigTrackerPollInterval" gorm:"type:bigint"`
		ContractConfigConfirmations            uint16         `json:"contractConfigConfirmations" toml:"contractConfigConfirmations"`
		CreatedAt                              time.Time      `json:"createdAt" toml:"-"`
		UpdatedAt                              time.Time      `json:"updatedAt" toml:"-"`
	}

	PeerID peer.ID
)

func (js JobSpecV2) GetID() string {
	return fmt.Sprintf("%v", js.ID)
}

func (js *JobSpecV2) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	js.ID = int32(ID)
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

func (p PeerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(peer.Encode(peer.ID(p)))
}

func (p *PeerID) UnmarshalJSON(input []byte) error {
	var result string
	if err := json.Unmarshal(input, &result); err != nil {
		return err
	}

	peerId, err := peer.Decode(result)
	if err != nil {
		return err
	}

	*p = PeerID(peerId)
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
func (OffchainReportingOracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
