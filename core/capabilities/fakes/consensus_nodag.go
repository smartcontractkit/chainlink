package fakes

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	valuespb "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
)

type fakeConsensusNoDAG struct {
	services.Service
	eng *services.Engine
}

var _ services.Service = (*fakeConsensus)(nil)
var _ consensusserver.ConsensusCapability = (*fakeConsensusNoDAG)(nil)

func NewFakeConsensusNoDAG(lggr logger.Logger) *fakeConsensusNoDAG {
	fc := &fakeConsensusNoDAG{}
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
	fakeOutput := &sdkpb.ConsensusOutputs{
		ConfigDigest:  []byte("fake_config_digest"),
		SeqNr:         42,
		ReportContext: []byte("fake_report_context"),
		Sigs: []*sdkpb.AttributedSignature{
			{
				SignerId:  3,
				Signature: []byte("fake_signature_value"),
			},
		},
	}

	switch input.Descriptors.EncoderName {
	case "proto", "": // mode-switch (default)
		mapProto := &valuespb.Map{
			Fields: map[string]*valuespb.Value{
				sdk.ConsensusResponseMapKeyMetadata: {Value: &valuespb.Value_StringValue{StringValue: "fake_metadata"}},
				sdk.ConsensusResponseMapKeyPayload:  input.GetValue(),
			},
		}
		rawMap, err := proto.Marshal(mapProto)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input value: %w", err)
		}
		fakeOutput.RawReport = rawMap
		// other fields are unused by mode-switch calls, keeping the fake ones
		return fakeOutput, nil
	case "evm": // report-gen for EVM
		encodedBytes := input.GetValue().GetBytesValue()
		if len(encodedBytes) == 0 {
			return nil, errors.New("input value for EVM encoder needs to be a byte array and cannot be empty or nil")
		}
		// prepend metadata
		// TODO: do real EVM metadata encoding here
		rawOutput := append([]byte("fake_evm_metadata"), encodedBytes...)
		fakeOutput.RawReport = rawOutput
		// TODO: add real EVM signatures here
		return fakeOutput, nil
	default:
		return nil, fmt.Errorf("unsupported encoder name: %s", input.Descriptors.EncoderName)
	}
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
