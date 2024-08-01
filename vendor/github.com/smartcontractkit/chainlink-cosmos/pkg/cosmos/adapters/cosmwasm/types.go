package cosmwasm

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type TransmitMsg struct {
	Transmit TransmitPayload `json:"transmit"`
}

type TransmitPayload struct {
	ReportContext []byte   `json:"report_context"`
	Report        []byte   `json:"report"`
	Signatures    [][]byte `json:"signatures"`
}

type ConfigDetails struct {
	BlockNumber  uint64             `json:"block_number"`
	ConfigDigest types.ConfigDigest `json:"config_digest"`
}

type LatestTransmissionDetails struct {
	LatestConfigDigest types.ConfigDigest `json:"latest_config_digest"`
	Epoch              uint32             `json:"epoch"`
	Round              uint8              `json:"round"`
	LatestAnswer       string             `json:"latest_answer"`
	LatestTimestamp    int64              `json:"latest_timestamp"`
}

type LatestConfigDigestAndEpoch struct {
	ConfigDigest types.ConfigDigest `json:"config_digest"`
	Epoch        uint32             `json:"epoch"`
}
