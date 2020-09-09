package offchainreporting

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	OffchainreportingKeyBundle struct {
		ID                     models.Sha256Hash `gorm:"primary_key"`
		EncryptedPrivKeyBundle models.JSON
		CreatedAt              time.Time
	}

	OffchainreportingOracleSpec struct {
		ID                    int32              `gorm:"primary_key"`
		ContractAddress       common.Address     `toml:"contractAddress"`
		P2PPeerID             string             `toml:"p2pPeerID"`
		P2PBootstrapPeers     []P2PBootstrapPeer `toml:"p2pBootstrapPeers"`
		KeyBundle             KeyBundle          `toml:"keyBundle"`
		MonitoringEndpoint    string             `toml:"monitoringEndpoint"`
		TransmitterAddress    common.Address     `toml:"transmitterAddress"`
		ObservationTimeout    time.Duration      `toml:"observationTimeout"`
		DataFetchPipelineSpec pipeline.PipelineSpec
		LogLevel              ocrtypes.LogLevel `toml:"logLevel,omitempty"`
	}

	P2PBootstrapPeer struct {
		PeerID    string `toml:"peerID"`
		Multiaddr string `toml:"multiAddr"`
	}
)
