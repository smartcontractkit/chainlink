package mockllo

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"text/template"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/cre"
)

type FeedWithStreamID struct {
	Feed     string `json:"feed"`
	StreamID int32  `json:"streamID"`
}

func (f FeedWithStreamID) MustFeedConfig() FeedConfig {
	feedID, err2 := datastreams.NewFeedID(f.Feed)
	if err2 != nil {
		panic("invalid feedID: " + err2.Error())
	}
	feedBytes := feedID.Bytes()
	return FeedConfig{
		FeedIDsIndex: f.StreamID,
		Deviation:    "0.001",
		Heartbeat:    3600,
		RemappedID:   "0x" + hex.EncodeToString(feedBytes[:]),
	}
}

func CreateLLOFeedReport(lggr logger.Logger, price decimal.Decimal, timestamp uint64,
	feeds []FeedWithStreamID, keyBundles []ocr2key.KeyBundle) (*capabilities.OCRTriggerEvent, string, error) {
	values := make([]datastreamsllo.StreamValue, 0)

	priceBytes, err := price.MarshalBinary()
	if err != nil {
		return nil, "", err
	}
	streams := make([]llotypes.Stream, 0)

	for _, f := range feeds {
		dec := &datastreamsllo.Decimal{}
		err2 := dec.UnmarshalBinary(priceBytes)
		if err2 != nil {
			return nil, "", err2
		}
		values = append(values, dec)
		streams = append(streams, llotypes.Stream{
			StreamID: llotypes.StreamID(f.StreamID), //nolint:gosec // G115 don't care in test code
		})
	}

	reportCodec := cre.NewReportCodecCapabilityTrigger(lggr, 1)

	report := datastreamsllo.Report{
		ObservationTimestampNanoseconds: timestamp,
		Values:                          values,
	}

	reportBytes, err := reportCodec.Encode(report, llotypes.ChannelDefinition{
		Streams: streams,
	})
	if err != nil {
		return nil, "", err
	}
	eventID := reportCodec.EventID(report)

	event := &capabilities.OCRTriggerEvent{
		ConfigDigest: []byte{0: 1, 31: 2},
		SeqNr:        0,
		Report:       reportBytes,
		Sigs:         make([]capabilities.OCRAttributedOnchainSignature, 0, len(keyBundles)),
	}

	for i, key := range keyBundles {
		sig, err2 := key.Sign3(ocrTypes.ConfigDigest(event.ConfigDigest), event.SeqNr, reportBytes)
		if err2 != nil {
			return nil, "", err
		}
		event.Sigs = append(event.Sigs, capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(i), //nolint:gosec // G115 don't care in test code
			Signature: sig,
		})
	}

	return event, eventID, nil
}

func GenerateFeedAddresses(streamIDs ...int32) []FeedWithStreamID {
	feedsAddresses := make([]FeedWithStreamID, 0)
	for _, streamID := range streamIDs {
		_, id := NewFeedIDDF2()
		feedsAddresses = append(feedsAddresses, FeedWithStreamID{
			Feed:     id,
			StreamID: streamID,
		})
	}
	return feedsAddresses
}

// NewFeedIDDF2 creates a random Data Feeds 2.0 format https://docs.google.com/document/d/13ciwTx8lSUfyz1IdETwpxlIVSn1lwYzGtzOBBTpl5Vg/edit?tab=t.0#heading=h.dxx2wwn1dmoz
func NewFeedIDDF2() ([32]byte, string) {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	if err != nil {
		panic("cannot create feedID: " + err.Error())
	}

	buf[0] = 0x01 // format byte
	buf[5] = 0x00 // attribute
	buf[6] = 0x03 // attribute
	buf[7] = 0x00 // data type byte

	for i := 8; i < 16; i++ {
		buf[i] = 0x00
	}

	return buf, "0x" + hex.EncodeToString(buf[:])
}

type FeedConfig struct {
	FeedIDsIndex int32  `json:"feedIDsIndex"`
	Deviation    string `json:"deviation"`
	Heartbeat    int32  `json:"heartbeat"`
	RemappedID   string `json:"remappedID"`
}

// TODO shouldn't consumer address be configurable?
func WorkflowsJob(nodeID string, workflowName string, feeds []FeedConfig) *jobv1.ProposeJobRequest {
	const workflowTemplateLoad = `
 type = "workflow"
 schemaVersion = 1
 name = "{{ .WorkflowName }}"
 externalJobID = "{{ .JobID }}"
 workflow = """
 name: "{{ .WorkflowName }}"
 owner: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
 triggers:
  - id: streams-trigger@2.0.0
    config:
      feedIds:
 {{- range .Feeds }}
        - '{{ .FeedIDsIndex }}'
 {{- end }}
 consensus:
   - id: "offchain_reporting@1.0.0"
     ref: "evm_median"
     inputs:
       observations:
         - "$(trigger.outputs)"
     config:
       report_id: "0001"
       key_id: "evm"
       aggregation_method: "llo_streams"
       aggregation_config:
         streams:
{{- range .Feeds }}
           "{{ .FeedIDsIndex }}":
             deviation: "{{ .Deviation }}"
             heartbeat: {{ .Heartbeat }}
             remappedID: {{ .RemappedID }}
{{- end }}
       encoder: "EVM"
       encoder_config:
         abi: "(bytes32 RemappedID, uint224 Price, uint32 Timestamp)[] Reports"
 targets:
   - id: write_ethereum_mock@1.0.0
     inputs:
       signed_report: "$(evm_median.outputs)"
     config:
       address: "0xEB739A9641938934D21A325A0C6b26126D48926A"
       params: ["$(report)"]
       abi: "receive(report bytes)"
       deltaStage: 2s
       schedule: allAtOnce
 """
 `

	tmpl, err := template.New("workflow").Parse(workflowTemplateLoad)

	if err != nil {
		panic(err)
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, map[string]interface{}{
		"WorkflowName": workflowName,
		"Feeds":        feeds,
		"JobID":        uuid.NewString(),
	})
	if err != nil {
		panic(err)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String()}
}
