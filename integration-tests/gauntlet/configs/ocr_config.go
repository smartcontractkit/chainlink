package configs

import (
	"encoding/json"
)

type OCRConfig struct {
	Contract *Contract
	Config   *Config
}

// Contract Required fields to deploy OCR contract via gauntlet
type Contract struct {
	MaximumGasPrice           int    `json:"maximumGasPrice"`
	ReasonableGasPrice        int    `json:"reasonableGasPrice"`
	MicroLinkPerEth           int    `json:"microLinkPerEth"`
	LinkGweiPerObservation    int    `json:"linkGweiPerObservation"`
	LinkGweiPerTransmission   int    `json:"linkGweiPerTransmission"`
	MinAnswer                 int    `json:"minAnswer"`
	MaxAnswer                 int    `json:"maxAnswer"`
	Decimals                  int    `json:"decimals"`
	Description               string `json:"description"`
	Link                      string `json:"link"`
	BillingAccessController   string `json:"billingAccessController"`
	RequesterAccessController string `json:"requesterAccessController"`
}

// Config Required fields to configure OCR contract via Gauntlet
type Config struct {
	Signers                       []string `json:"signers"`
	Transmitters                  []string `json:"transmitters"`
	OcrConfigPublicKeys           []string `json:"operatorsPublicKeys"`
	OperatorsPeerIds              string   `json:"operatorsPeerIds"`
	Threshold                     int      `json:"threshold"`
	BadEpochTimeout               string   `json:"badEpochTimeout"`
	ResendInterval                string   `json:"resendInterval"`
	RoundInterval                 string   `json:"roundInterval"`
	ObservationGracePeriod        string   `json:"observationGracePeriod"`
	MaxContractValueAge           string   `json:"maxContractValueAge"`
	RelativeDeviationThresholdPPB string   `json:"relativeDeviationThresholdPPB"`
	TransmissionStageTimeout      string   `json:"transmissionStageTimeout"`
	MaxRoundCount                 int      `json:"maxRoundCount"`
	TransmissionStages            []int    `json:"transmissionStages"`
	Secret                        string   `json:"secret"`
}

func NewOCR() *OCRConfig {
	return &OCRConfig{
		Contract: &Contract{},
		Config:   &Config{},
	}
}

func (o *OCRConfig) DefaultOcrContract() {
	o.Contract = &Contract{
		MaximumGasPrice:         2000,
		ReasonableGasPrice:      10,
		MicroLinkPerEth:         102829,
		LinkGweiPerObservation:  600,
		LinkGweiPerTransmission: 3000,
		MinAnswer:               1,
		MaxAnswer:               100000,
		Decimals:                8,
		Description:             "ETH/USD",
	}
}

func (o *OCRConfig) DefaultOcrConfig() {
	o.Config = &Config{
		Threshold:                     1,
		BadEpochTimeout:               "35s",
		ResendInterval:                "17s",
		RoundInterval:                 "30s",
		ObservationGracePeriod:        "12s",
		MaxContractValueAge:           "1h",
		RelativeDeviationThresholdPPB: "10000000",
		TransmissionStageTimeout:      "60s",
		MaxRoundCount:                 6,
		TransmissionStages:            []int{1, 2, 2, 2},
		Secret:                        "awe accuse polygon tonic depart acuity onyx inform bound gilbert expire",
	}
}

// MarshalContract Returns JSON string representation of the OCR Contract that is provided to Gauntlet as --input
func (o *OCRConfig) MarshalContract() (string, error) {
	parsedConfig, err := json.Marshal(o.Contract)
	if err != nil {
		return "", err
	}
	return string(parsedConfig), nil
}

// MarshalConfig Returns JSON string representation of the OCR Config that is provided to Gauntlet as --input
func (o *OCRConfig) MarshalConfig() (string, error) {
	parsedConfig, err := json.Marshal(o.Config)
	if err != nil {
		return "", err
	}
	return string(parsedConfig), nil
}
