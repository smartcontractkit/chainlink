package contracts

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

type OffChainAggregatorV2Config struct {
	DeltaProgress                           *commonconfig.Duration             `toml:",omitempty"`
	DeltaResend                             *commonconfig.Duration             `toml:",omitempty"`
	DeltaRound                              *commonconfig.Duration             `toml:",omitempty"`
	DeltaGrace                              *commonconfig.Duration             `toml:",omitempty"`
	DeltaStage                              *commonconfig.Duration             `toml:",omitempty"`
	RMax                                    uint8                              `toml:"-"`
	S                                       []int                              `toml:"-"`
	Oracles                                 []confighelper.OracleIdentityExtra `toml:"-"`
	ReportingPluginConfig                   []byte                             `toml:"-"`
	MaxDurationQuery                        *commonconfig.Duration             `toml:",omitempty"`
	MaxDurationObservation                  *commonconfig.Duration             `toml:",omitempty"`
	MaxDurationReport                       *commonconfig.Duration             `toml:",omitempty"`
	MaxDurationShouldAcceptFinalizedReport  *commonconfig.Duration             `toml:",omitempty"`
	MaxDurationShouldTransmitAcceptedReport *commonconfig.Duration             `toml:",omitempty"`
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
