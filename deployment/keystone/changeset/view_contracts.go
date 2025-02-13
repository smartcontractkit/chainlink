package changeset

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

type KeystoneChainView struct {
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	// OCR3ConfigView is a map of OCR3 contract addresses to their configuration view
	OCR3ConfigView map[string]OCR3ConfigView `json:"ocr3ConfigViews,omitempty"`
}

type OCR3ConfigView struct {
	Signers               []string            `json:"signers"`
	Transmitters          []ocr2types.Account `json:"transmitters"`
	F                     uint8               `json:"f"`
	OnchainConfig         []byte              `json:"onchainConfig"`
	OffchainConfigVersion uint64              `json:"offchainConfigVersion"`
	OffchainConfig        interface{}         `json:"offchainConfig"` // TODO: we need a struct here to hold the values
}

var ErrOCR3NotConfigured = errors.New("OCR3 not configured")

func GenerateOCR3ConfigView(ocr3Cap ocr3_capability.OCR3Capability) (OCR3ConfigView, error) {
	details, err := ocr3Cap.LatestConfigDetails(nil)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	blockNumber := uint64(details.BlockNumber)
	config, err := ocr3Cap.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNumber,
		End:     &blockNumber,
		Context: nil,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}

	if config.Event == nil {
		return OCR3ConfigView{}, ErrOCR3NotConfigured
	}
	var signers []ocr2types.OnchainPublicKey
	var readableSigners []string
	for _, s := range config.Event.Signers {
		signers = append(signers, s)
		readableSigners = append(readableSigners, string(s))
	}
	var transmitters []ocr2types.Account
	for _, t := range config.Event.Transmitters {
		transmitters = append(transmitters, ocr2types.Account(t.String()))
	}
	// `PublicConfigFromContractConfig` returns the `ocr2types.PublicConfig` that contains all the `OracleConfig` fields we need, including the
	// report plugin config.
	_, err = ocr3confighelper.PublicConfigFromContractConfig(true, ocr2types.ContractConfig{
		ConfigDigest:          config.Event.ConfigDigest,
		ConfigCount:           config.Event.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     config.Event.F,
		OnchainConfig:         nil, // empty onChain config
		OffchainConfigVersion: config.Event.OffchainConfigVersion,
		OffchainConfig:        config.Event.OffchainConfig,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}

	// TODO: make human readable
	return OCR3ConfigView{
		Signers:               readableSigners,
		Transmitters:          transmitters,
		F:                     config.Event.F,
		OnchainConfig:         nil, // empty onChain config
		OffchainConfigVersion: config.Event.OffchainConfigVersion,
		OffchainConfig:        config.Event.OffchainConfig,
	}, nil
}

func NewKeystoneChainView() KeystoneChainView {
	return KeystoneChainView{
		CapabilityRegistry: make(map[string]common_v1_0.CapabilityRegistryView),
		OCR3ConfigView:     make(map[string]OCR3ConfigView),
	}
}

type KeystoneView struct {
	Chains map[string]KeystoneChainView `json:"chains,omitempty"`
	Nops   map[string]view.NopView      `json:"nops,omitempty"`
}

func (v KeystoneView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias KeystoneView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
