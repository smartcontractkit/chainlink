package view

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

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
	Signers               [][]byte         `json:"signers"`
	Transmitters          []common.Address `json:"transmitters"`
	F                     uint8            `json:"f"`
	OnchainConfig         []byte           `json:"onchainConfig"`
	OffchainConfigVersion uint64           `json:"offchainConfigVersion"`
	OffchainConfig        []byte           `json:"offchainConfig"`
}

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

	return OCR3ConfigView{
		Signers:               config.Event.Signers,
		Transmitters:          config.Event.Transmitters,
		F:                     config.Event.F,
		OnchainConfig:         config.Event.OnchainConfig,
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
