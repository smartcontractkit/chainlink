package config

import (
	"sort"
	"strings"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

type OCR2Config struct {
	OnChainConfig        *OCR2OnChainConfig
	OffChainConfig       *OCROffChainConfig
	PayeeConfig          *PayeeConfig
	ProposalAcceptConfig *ProposalAcceptConfig
	NodeKeys             []client.NodeKeysBundle
	VaultAddress         string
	Secret               string
	ProposalID           string
}

type OCR2OnChainConfig struct {
	Oracles    []Operator `json:"oracles"`
	F          int        `json:"f"`
	ProposalID string     `json:"proposalId"`
}

type OffchainConfig struct {
	DeltaProgressNanoseconds                           int64                 `json:"deltaProgressNanoseconds"`
	DeltaResendNanoseconds                             int64                 `json:"deltaResendNanoseconds"`
	DeltaRoundNanoseconds                              int64                 `json:"deltaRoundNanoseconds"`
	DeltaGraceNanoseconds                              int64                 `json:"deltaGraceNanoseconds"`
	DeltaStageNanoseconds                              int64                 `json:"deltaStageNanoseconds"`
	RMax                                               int                   `json:"rMax"`
	S                                                  []int                 `json:"s"`
	OffchainPublicKeys                                 []string              `json:"offchainPublicKeys"`
	PeerIds                                            []string              `json:"peerIds"`
	ReportingPluginConfig                              ReportingPluginConfig `json:"reportingPluginConfig"`
	MaxDurationQueryNanoseconds                        int64                 `json:"maxDurationQueryNanoseconds"`
	MaxDurationObservationNanoseconds                  int64                 `json:"maxDurationObservationNanoseconds"`
	MaxDurationReportNanoseconds                       int64                 `json:"maxDurationReportNanoseconds"`
	MaxDurationShouldAcceptFinalizedReportNanoseconds  int64                 `json:"maxDurationShouldAcceptFinalizedReportNanoseconds"`
	MaxDurationShouldTransmitAcceptedReportNanoseconds int64                 `json:"maxDurationShouldTransmitAcceptedReportNanoseconds"`
	ConfigPublicKeys                                   []string              `json:"configPublicKeys"`
}

type ReportingPluginConfig struct {
	AlphaReportInfinite bool `json:"alphaReportInfinite"`
	AlphaReportPpb      int  `json:"alphaReportPpb"`
	AlphaAcceptInfinite bool `json:"alphaAcceptInfinite"`
	AlphaAcceptPpb      int  `json:"alphaAcceptPpb"`
	DeltaCNanoseconds   int  `json:"deltaCNanoseconds"`
}

// TODO - Decouple all OCR2 config structs to be reusable between chains
type OCROffChainConfig struct {
	ProposalID     string         `json:"proposalId"`
	OffchainConfig OffchainConfig `json:"offchainConfig"`
	UserSecret     string         `json:"userSecret"`
}

type Operator struct {
	Signer      string `json:"signer"`
	Transmitter string `json:"transmitter"`
	Payee       string `json:"payee"`
}

type PayeeConfig struct {
	Operators  []Operator `json:"operators"`
	ProposalID string     `json:"proposalId"`
}

type ProposalAcceptConfig struct {
	ProposalID     string         `json:"proposalId"`
	Version        int            `json:"version"`
	F              int            `json:"f"`
	Oracles        []Operator     `json:"oracles"`
	OffchainConfig OffchainConfig `json:"offchainConfig"`
	RandomSecret   string         `json:"randomSecret"`
}

type OCR2TransmitConfig struct {
	MinAnswer     string `json:"minAnswer"`
	MaxAnswer     string `json:"maxAnswer"`
	Transmissions string `json:"transmissions"`
}

type OCR2BillingConfig struct {
	ObservationPaymentGjuels  int `json:"ObservationPaymentGjuels"`
	TransmissionPaymentGjuels int `json:"TransmissionPaymentGjuels"`
}

type StoreFeedConfig struct {
	Store       string `json:"store"`
	Granularity int    `json:"granularity"`
	LiveLength  int    `json:"liveLength"`
	Decimals    int    `json:"decimals"`
	Description string `json:"description"`
}

type StoreWriterConfig struct {
	Transmissions string `json:"transmissions"`
}

func NewOCR2Config(nodeKeys []client.NodeKeysBundle, proposalID string, vaultAddress string, secret string) *OCR2Config {
	var oracles []Operator

	nodeKeysSorted := make([]client.NodeKeysBundle, len(nodeKeys))
	copy(nodeKeysSorted, nodeKeys)

	// We have to sort by on_chain_pub_key for the config digest
	sort.Slice(nodeKeysSorted, func(i, j int) bool {
		return nodeKeysSorted[i].OCR2Key.Data.Attributes.OnChainPublicKey < nodeKeysSorted[j].OCR2Key.Data.Attributes.OnChainPublicKey
	})

	for _, nodeKey := range nodeKeysSorted {
		oracles = append(oracles, Operator{
			Signer:      strings.Replace(nodeKey.OCR2Key.Data.Attributes.OnChainPublicKey, "ocr2on_solana_", "", 1),
			Transmitter: nodeKey.TXKey.Data.Attributes.PublicKey,
			Payee:       vaultAddress,
		})
	}

	return &OCR2Config{
		OnChainConfig: &OCR2OnChainConfig{
			Oracles:    oracles,
			F:          1,
			ProposalID: proposalID,
		},
		OffChainConfig:       &OCROffChainConfig{},
		PayeeConfig:          &PayeeConfig{},
		ProposalAcceptConfig: &ProposalAcceptConfig{},
		NodeKeys:             nodeKeysSorted,
		VaultAddress:         vaultAddress,
		Secret:               secret,
		ProposalID:           proposalID,
	}
}

func (o *OCR2Config) Default() {
	o.OffChainConfig.OffchainConfig.ReportingPluginConfig = ReportingPluginConfig{
		AlphaReportInfinite: false,
		AlphaReportPpb:      0,
		AlphaAcceptInfinite: false,
		AlphaAcceptPpb:      0,
		DeltaCNanoseconds:   0,
	}
	offchainPublicKeys := make([]string, len(o.NodeKeys))
	peerIds := make([]string, len(o.NodeKeys))
	configPublicKeys := make([]string, len(o.NodeKeys))
	s := make([]int, len(o.NodeKeys))

	for i := range s {
		s[i] = 1
	}

	for i, key := range o.NodeKeys {
		offchainPublicKeys[i] = strings.Replace(key.OCR2Key.Data.Attributes.OffChainPublicKey, "ocr2off_solana_", "", 1)
		peerIds[i] = key.PeerID
		configPublicKeys[i] = strings.Replace(key.OCR2Key.Data.Attributes.ConfigPublicKey, "ocr2cfg_solana_", "", 1)
	}
	o.OffChainConfig = &OCROffChainConfig{
		UserSecret: o.Secret,
		ProposalID: o.ProposalID,
		OffchainConfig: OffchainConfig{
			DeltaProgressNanoseconds:          int64(8000000000),  // 8s
			DeltaResendNanoseconds:            int64(5000000000),  // 5s
			DeltaRoundNanoseconds:             int64(3000000000),  // 3s
			DeltaGraceNanoseconds:             int64(400000000),   // 400ms
			DeltaStageNanoseconds:             int64(10000000000), // 10s
			RMax:                              3,
			S:                                 s,
			OffchainPublicKeys:                offchainPublicKeys,
			PeerIds:                           peerIds,
			ConfigPublicKeys:                  configPublicKeys,
			ReportingPluginConfig:             o.OffChainConfig.OffchainConfig.ReportingPluginConfig,
			MaxDurationQueryNanoseconds:       int64(0),
			MaxDurationObservationNanoseconds: int64(1000000000), // 1s
			MaxDurationReportNanoseconds:      int64(1000000000), // 1s
			MaxDurationShouldAcceptFinalizedReportNanoseconds:  int64(1000000000), // 1s
			MaxDurationShouldTransmitAcceptedReportNanoseconds: int64(1000000000), // 1s
		},
	}
	o.PayeeConfig = &PayeeConfig{
		Operators:  o.OnChainConfig.Oracles,
		ProposalID: o.ProposalID,
	}
	o.ProposalAcceptConfig = &ProposalAcceptConfig{
		ProposalID:     o.ProposalID,
		Version:        2,
		F:              1,
		Oracles:        o.OnChainConfig.Oracles,
		OffchainConfig: o.OffChainConfig.OffchainConfig,
		RandomSecret:   o.Secret,
	}
}
