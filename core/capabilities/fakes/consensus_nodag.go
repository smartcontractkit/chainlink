package fakes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/report"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

type fakeConsensusNoDAG struct {
	services.Service
	eng *services.Engine

	signers      []ocr2key.KeyBundle
	configDigest ocr2types.ConfigDigest
	seqNr        uint64
}

var (
	_ services.Service                    = (*fakeConsensus)(nil)
	_ consensusserver.ConsensusCapability = (*fakeConsensusNoDAG)(nil)
)

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
func (fc *fakeConsensusNoDAG) Simple(ctx context.Context, metadata capabilities.RequestMetadata, input *sdkpb.SimpleConsensusInputs) (*capabilities.ResponseAndMetadata[*valuespb.Value], caperrors.Error) {
	fc.eng.Infow("Executing Fake Consensus NoDAG: Simple()", "input", input, "metadata", metadata)

	switch obs := input.Observation.(type) {
	case *sdkpb.SimpleConsensusInputs_Value:
		if obs.Value == nil {
			return nil, caperrors.NewPublicUserError(errors.New("input value cannot be nil"), caperrors.InvalidArgument)
		}
		responseAndMetadata := capabilities.ResponseAndMetadata[*valuespb.Value]{
			Response:         obs.Value,
			ResponseMetadata: capabilities.ResponseMetadata{},
		}
		return &responseAndMetadata, nil
	case *sdkpb.SimpleConsensusInputs_Error:
		return nil, caperrors.NewPublicSystemError(errors.New(obs.Error), caperrors.Unknown)
	case nil:
		return nil, caperrors.NewPublicUserError(errors.New("input observation cannot be nil"), caperrors.InvalidArgument)
	default:
		return nil, caperrors.NewPublicUserError(errors.New("unknown observation type"), caperrors.InvalidArgument)
	}
}

func (fc *fakeConsensusNoDAG) Report(ctx context.Context, metadata capabilities.RequestMetadata, input *sdkpb.ReportRequest) (*capabilities.ResponseAndMetadata[*sdkpb.ReportResponse], caperrors.Error) {
	fc.eng.Infow("Executing Fake Consensus NoDAG: Report()", "input", input, "metadata", metadata)
	// Prepare EVM metadata that will be prepended to all reports
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

	switch input.EncoderName {
	case "evm", "EVM": // report-gen for EVM
		if len(input.EncodedPayload) == 0 {
			return nil, caperrors.NewPublicUserError(errors.New("input value for EVM encoder needs to be a byte array and cannot be empty or nil"), caperrors.InvalidArgument)
		}

		// Prepend EVM metadata
		rawOutput, err := meta.Encode()
		if err != nil {
			return nil, caperrors.NewPublicSystemError(fmt.Errorf("failed to prepend metadata fields: %w", err), caperrors.Internal)
		}
		rawOutput = append(rawOutput, input.EncodedPayload...)

		// sign the report
		sigs := []*sdkpb.AttributedSignature{}
		var idx uint32
		for _, signer := range fc.signers {
			sig, err := signer.Sign3(fc.configDigest, fc.seqNr, rawOutput)
			if err != nil {
				return nil, caperrors.NewPublicSystemError(fmt.Errorf("failed to sign with signer %s: %w", signer.ID(), err), caperrors.Internal)
			}
			sigs = append(sigs, &sdkpb.AttributedSignature{
				SignerId:  idx,
				Signature: sig,
			})
			idx++
		}

		reportResponse := &sdkpb.ReportResponse{
			RawReport:     rawOutput,
			ConfigDigest:  fc.configDigest[:],
			SeqNr:         fc.seqNr,
			ReportContext: report.GenerateReportContext(fc.seqNr, fc.configDigest),
			Sigs:          sigs,
		}
		responseAndMetadata := capabilities.ResponseAndMetadata[*sdkpb.ReportResponse]{
			Response:         reportResponse,
			ResponseMetadata: capabilities.ResponseMetadata{},
		}
		return &responseAndMetadata, nil

	default:
		return nil, caperrors.NewPublicUserError(fmt.Errorf("unsupported encoder name: %s", input.EncoderName), caperrors.InvalidArgument)
	}
}

func (fc *fakeConsensusNoDAG) Description() string {
	return "Fake OCR Consensus NoDAG"
}

func (fc *fakeConsensusNoDAG) Initialise(
	_ context.Context,
	_ core.StandardCapabilitiesDependencies,
) error {
	return nil
}

func SeedForKeys() io.Reader {
	byteArray := make([]byte, 10000)
	for i := range 10000 {
		byteArray[i] = byte((420666 + i) % 256)
	}
	return bytes.NewReader(byteArray)
}
