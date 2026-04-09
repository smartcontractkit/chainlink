package contracts

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

type OffChainAggregatorV2Config struct {
	DeltaProgress                           *config.Duration                   `toml:",omitempty"`
	DeltaResend                             *config.Duration                   `toml:",omitempty"`
	DeltaRound                              *config.Duration                   `toml:",omitempty"`
	DeltaGrace                              *config.Duration                   `toml:",omitempty"`
	DeltaStage                              *config.Duration                   `toml:",omitempty"`
	RMax                                    uint8                              `toml:"-"`
	S                                       []int                              `toml:"-"`
	Oracles                                 []confighelper.OracleIdentityExtra `toml:"-"`
	ReportingPluginConfig                   []byte                             `toml:"-"`
	MaxDurationQuery                        *config.Duration                   `toml:",omitempty"`
	MaxDurationObservation                  *config.Duration                   `toml:",omitempty"`
	MaxDurationReport                       *config.Duration                   `toml:",omitempty"`
	MaxDurationShouldAcceptFinalizedReport  *config.Duration                   `toml:",omitempty"`
	MaxDurationShouldTransmitAcceptedReport *config.Duration                   `toml:",omitempty"`
	F                                       int                                `toml:"-"`
	OnchainConfig                           []byte                             `toml:"-"`
}

type ChainlinkKeyExporter interface {
	ExportEVMKeysForChain(string) ([]*nodeclient.ExportedEVMKey, error)
}

type ChainlinkNodeWithKeysAndAddress interface {
	MustReadOCRKeys() (*nodeclient.OCRKeys, error)
	MustReadP2PKeys() (*nodeclient.P2PKeys, error)
	PrimaryEthAddress() (string, error)
	EthAddresses() ([]string, error)
	ChainlinkKeyExporter
}

func ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*nodeclient.ChainlinkK8sClient) []ChainlinkNodeWithKeysAndAddress {
	nodes := make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
	for i, node := range k8sNodes {
		nodes[i] = node
	}
	return nodes
}
