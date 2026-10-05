package fakes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
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
)

// Report encoders. Every encoder emits the shared 109-byte metadata header
// followed by the raw payload; they differ only in the key type signing it.
const (
	encoderEVM     = "evm"
	encoderSolana  = "solana"
	encoderStellar = "stellar"
)

type fakeConsensusNoDAG struct {
	services.Service
	eng *services.Engine

	// signers is the default signer set. signersByEncoder overrides it per
	// encoder (keyed by lowercase encoder name) when a chain needs a different
	// signature scheme.
	signers          []ocr2key.KeyBundle
	signersByEncoder map[string][]ocr2key.KeyBundle
	configDigest     ocr2types.ConfigDigest
	seqNr            uint64
}

var _ consensusserver.ConsensusCapability = (*fakeConsensusNoDAG)(nil)

// FakeConsensusOption configures NewFakeConsensusNoDAG.
type FakeConsensusOption func(*fakeConsensusNoDAG)

// WithEncoderSigners signs reports for the given encoder (case-insensitive)
// with signers instead of the default set, e.g. to match the signer set
// configured on a real forwarder.
func WithEncoderSigners(encoderName string, signers []ocr2key.KeyBundle) FakeConsensusOption {
	return func(fc *fakeConsensusNoDAG) {
		fc.signersByEncoder[strings.ToLower(encoderName)] = signers
	}
}

// NewFakeConsensusNoDAG returns a single-node consensus fake. signers sign EVM
// and Solana reports. Stellar reports need ed25519 (keccak) signatures, so
// unless overridden via WithEncoderSigners they are signed by deterministic
// Stellar keys generated here, one per default signer.
func NewFakeConsensusNoDAG(signers []ocr2key.KeyBundle, lggr logger.Logger, opts ...FakeConsensusOption) *fakeConsensusNoDAG {
	configDigest := ocr2types.ConfigDigest{}
	for i := range configDigest {
		configDigest[i] = byte(i)
	}
	fc := &fakeConsensusNoDAG{
		signers:          signers,
		signersByEncoder: map[string][]ocr2key.KeyBundle{},
		configDigest:     configDigest,
		seqNr:            1,
	}
	for _, opt := range opts {
		opt(fc)
	}
	if _, ok := fc.signersByEncoder[encoderStellar]; !ok {
		fc.signersByEncoder[encoderStellar] = deterministicSigners(corekeys.Stellar, max(len(signers), 1))
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

// Simple bounces back the observation value, reshaping frequency_list fields
// to their correct output type. The simulator runs a single node, so there is
// exactly one observation: frequency_list turns a single value T into
// [{value: T, count: 1}]. All other aggregation types have matching input and
// output types, so the raw value is returned unchanged.
func (fc *fakeConsensusNoDAG) Simple(ctx context.Context, metadata capabilities.RequestMetadata, input *sdkpb.SimpleConsensusInputs) (*capabilities.ResponseAndMetadata[*valuespb.Value], caperrors.Error) {
	fc.eng.Infow("Executing Fake Consensus NoDAG: Simple()", "input", input, "metadata", metadata)

	switch obs := input.Observation.(type) {
	case *sdkpb.SimpleConsensusInputs_Value:
		if obs.Value == nil {
			return nil, caperrors.NewPublicUserError(errors.New("input value cannot be nil"), caperrors.InvalidArgument)
		}
		response := applyFrequencyListShape(obs.Value, input.Descriptors)
		responseAndMetadata := capabilities.ResponseAndMetadata[*valuespb.Value]{
			Response:         response,
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

	encoder := strings.ToLower(input.EncoderName)
	switch encoder {
	case encoderEVM, encoderSolana, encoderStellar:
		if len(input.EncodedPayload) == 0 {
			return nil, caperrors.NewPublicUserError(fmt.Errorf("input value for %s encoder needs to be a byte array and cannot be empty or nil", input.EncoderName), caperrors.InvalidArgument)
		}

		// Prepend the shared 109-byte metadata header
		rawOutput, err := meta.Encode()
		if err != nil {
			return nil, caperrors.NewPublicSystemError(fmt.Errorf("failed to prepend metadata fields: %w", err), caperrors.Internal)
		}
		rawOutput = append(rawOutput, input.EncodedPayload...)

		reportContext := report.GenerateReportContext(fc.seqNr, fc.configDigest)

		// sign the report
		sigs := []*sdkpb.AttributedSignature{}
		var idx uint32
		for _, signer := range fc.signersFor(encoder) {
			var sig []byte
			if encoder == encoderStellar {
				// The Stellar forwarder verifies ed25519 over
				// keccak256(keccak256(raw_report) ‖ report_context); the
				// keyring returns public_key ‖ signature, the pair the
				// forwarder's Ed25519Signature carries.
				sig, err = signer.SignBlob(stellarReportDigest(rawOutput, reportContext))
			} else {
				sig, err = signer.Sign3(fc.configDigest, fc.seqNr, rawOutput)
			}
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
			ReportContext: reportContext,
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

// stellarReportDigest is the digest the Stellar CRE forwarder verifies.
func stellarReportDigest(rawReport, reportContext []byte) []byte {
	inner := crypto.Keccak256(rawReport)
	return crypto.Keccak256(append(inner, reportContext...))
}

// signersFor returns the signer set for a lowercase encoder name.
func (fc *fakeConsensusNoDAG) signersFor(encoder string) []ocr2key.KeyBundle {
	if s, ok := fc.signersByEncoder[encoder]; ok {
		return s
	}
	return fc.signers
}

// deterministicSigners returns n distinct, reproducible key bundles of the
// given chain type, drawn sequentially from SeedForKeys.
func deterministicSigners(chainType corekeys.ChainType, n int) []ocr2key.KeyBundle {
	seed := SeedForKeys()
	out := make([]ocr2key.KeyBundle, n)
	for i := range out {
		out[i] = ocr2key.MustNewInsecure(seed, chainType)
	}
	return out
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
