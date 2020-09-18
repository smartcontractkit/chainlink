package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type (
	OffchainReportingOracleSpec struct {
		ID                                int32                               `toml:"-" gorm:"primary_key"`
		JID                               int32                               `toml:"-" gorm:"column:job_id"`
		ContractAddress                   common.Address                      `toml:"contractAddress"`
		P2PPeerID                         string                              `toml:"p2pPeerID"`
		P2PBootstrapPeers                 []OffchainReportingP2PBootstrapPeer `toml:"p2pBootstrapPeers"`
		KeyBundle                         OffchainReportingKeyBundle          `toml:"keyBundle"`
		MonitoringEndpoint                string                              `toml:"monitoringEndpoint"`
		TransmitterAddress                common.Address                      `toml:"transmitterAddress"`
		ObservationTimeout                time.Duration                       `toml:"observationTimeout"`
		BlockchainTimeout                 time.Duration
		ContractConfigTrackerPollInterval time.Duration
		ContractConfigConfirmations       uint16
		DataFetchPipelineSpecID           int32
	}

	OffchainReportingP2PBootstrapPeer struct {
		PeerID    string `toml:"peerID"`
		Multiaddr string `toml:"multiAddr"`
	}

	OffchainReportingKeyBundle struct {
		ID                     Sha256Hash `gorm:"primary_key"`
		EncryptedPrivKeyBundle JSON
		CreatedAt              time.Time
	}
)

func (OffchainReportingOracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
