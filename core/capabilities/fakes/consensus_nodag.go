package fakes

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/cre-sdk-go/sdk"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	valuespb "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type fakeConsensusNoDAG struct {
	services.Service
	eng *services.Engine

	signers      []ocr2key.KeyBundle
	configDigest ocr2types.ConfigDigest
	seqNr        uint32
}

var _ services.Service = (*fakeConsensus)(nil)
var _ consensusserver.ConsensusCapability = (*fakeConsensusNoDAG)(nil)

func NewFakeConsensusNoDAG(signers []ocr2key.KeyBundle, lggr logger.Logger) *fakeConsensusNoDAG {
	configDigest := ocr2types.ConfigDigest{}
	for i := range len(configDigest) {
		configDigest[i] = byte(i)
	}
	fc := &fakeConsensusNoDAG{
		signers:      signers,
		configDigest: configDigest,
		seqNr:        1,
	}
	fc.Service, fc.eng = services.Config{
		Name:  "fakeConsensusNoDAG",
		Start: fc.start,
		Close: fc.close,
	}.NewServiceEngine(lggr)
	return fc
}

func (fc *fakeConsensusNoDAG) start(ctx context.Context) error {
	return nil
}

func (fc *fakeConsensusNoDAG) close() error {
	return nil
}

// NOTE: This fake capability currently bounces back the request payload, ignoring everything else.
// When the real NoDAG consensus OCR plugin is ready, it should be used here, similarly to how the V1 fake works.
func (fc *fakeConsensusNoDAG) Simple(ctx context.Context, metadata capabilities.RequestMetadata, input *sdkpb.SimpleConsensusInputs) (*sdkpb.ConsensusOutputs, error) {
	fc.eng.Infow("Executing Fake Consensus NoDAG", "input", input)
	observation := input.GetValue()
	if observation == nil {
		return nil, errors.New("input value cannot be nil")
	}

	switch input.Descriptors.EncoderName {
	case "proto", "": // mode-switch (default)
		mapProto := &valuespb.Map{
			Fields: map[string]*valuespb.Value{
				sdk.ConsensusResponseMapKeyMetadata: {Value: &valuespb.Value_StringValue{StringValue: "fake_metadata"}},
				sdk.ConsensusResponseMapKeyPayload:  observation,
			},
		}
		rawMap, err := proto.Marshal(mapProto)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input value: %w", err)
		}
		return &sdkpb.ConsensusOutputs{
			RawReport: rawMap,
			// other fields are unused by mode-switch calls, use fake ones
			ConfigDigest:  []byte("fake_config_digest"),
			SeqNr:         42,
			ReportContext: []byte("fake_report_context"),
			Sigs: []*sdkpb.AttributedSignature{
				{
					SignerId:  3,
					Signature: []byte("fake_signature_value"),
				},
			},
		}, nil
	case "evm": // report-gen for EVM
		encodedBytes := observation.GetBytesValue()
		if len(encodedBytes) == 0 {
			return nil, errors.New("input value for EVM encoder needs to be a byte array and cannot be empty or nil")
		}

		// prepend metadata
		meta := consensustypes.Metadata{
			Version:          1,
			ExecutionID:      metadata.WorkflowExecutionID,
			Timestamp:        100,
			DONID:            metadata.WorkflowDonID,
			DONConfigVersion: metadata.WorkflowDonConfigVersion,
			WorkflowID:       metadata.WorkflowID,
			WorkflowName:     metadata.WorkflowName,
			WorkflowOwner:    metadata.WorkflowOwner,
			ReportID:         "0001",
		}
		rawOutput, err := evm.PrependMetadataFields(meta, encodedBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to prepend metadata fields: %w", err)
		}

		// sign the report
		sigs := []*sdkpb.AttributedSignature{}
		var idx uint32
		for _, signer := range fc.signers {
			sig, err := signer.Sign3(fc.configDigest, uint64(fc.seqNr), rawOutput)
			if err != nil {
				return nil, fmt.Errorf("failed to sign with signer %s: %w", signer.ID(), err)
			}
			sigs = append(sigs, &sdkpb.AttributedSignature{
				SignerId:  idx,
				Signature: sig,
			})
			idx++
		}

		return &sdkpb.ConsensusOutputs{
			RawReport:     rawOutput,
			ConfigDigest:  fc.configDigest[:],
			SeqNr:         uint64(fc.seqNr),
			ReportContext: reportContext(fc.configDigest[:], fc.seqNr),
			Sigs:          sigs,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported encoder name: %s", input.Descriptors.EncoderName)
	}
}

func reportContext(configDigest []byte, seqNr uint32) []byte {
	// report context is the config digest + the sequence number padded with zeros
	seqToEpoch := make([]byte, 32)
	binary.BigEndian.PutUint32(seqToEpoch[32-5:32-1], seqNr)
	zeros := make([]byte, 32)
	return append(append(configDigest, seqToEpoch...), zeros...)
}

func (fc *fakeConsensusNoDAG) Description() string {
	return "Fake OCR Consensus NoDAG"
}

func (fc *fakeConsensusNoDAG) Initialise(
	_ context.Context,
	_ string,
	_ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory,
	_ core.GatewayConnector,
	_ core.Keystore) error {
	return nil
}

func SeedForKeys() io.Reader {
	byteArray := make([]byte, 10000)
	for i := range 10000 {
		byteArray[i] = byte((420666 + i) % 256)
	}
	return bytes.NewReader(byteArray)
}
