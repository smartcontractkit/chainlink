package evm

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

var (
	_ llo.ReportCodec = ReportCodecEVMStreamlined{}
)

func NewReportCodecStreamlined() ReportCodecEVMStreamlined {
	return ReportCodecEVMStreamlined{}
}

type ReportCodecEVMStreamlined struct{}

func (rc ReportCodecEVMStreamlined) Encode(ctx context.Context, r llo.Report, cd llotypes.ChannelDefinition) (payload []byte, err error) {
	opts := ReportFormatEVMStreamlinedOpts{}
	if err = (&opts).Decode(cd.Opts); err != nil {
		return nil, fmt.Errorf("failed to decode opts; got: '%s'; %w", cd.Opts, err)
	}

	if opts.FeedID == nil {
		payload = append(
			encodePackedUint32(uint32(cd.ReportFormat)),
			encodePackedUint32(r.ChannelID)...,
		)
	} else {
		// Must assume feed ID is universally unique so contains sufficient
		// domain separation and info about the report format.
		payload = opts.FeedID.Bytes()
	}
	payload = append(payload, encodePackedUint64(r.ValidAfterNanoseconds)...)
	// Pack-encode the rest of the values
	for i, encoder := range opts.ABI {
		b, err := encoder.EncodePacked(r.Values[i])
		if err != nil {
			return nil, fmt.Errorf("failed to encode stream value at index %d; %w", i, err)
		}
		payload = append(payload, b...)
	}

	return payload, nil
}

func encodePackedUint32(value uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		panic(err)
	}
	return buf.Bytes() // Returns 4 bytes
}

func encodePackedUint64(value uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		panic(err)
	}
	return buf.Bytes() // Returns 8 bytes
}

func (rc ReportCodecEVMStreamlined) Verify(_ context.Context, cd llotypes.ChannelDefinition) error {
	opts := ReportFormatEVMStreamlinedOpts{}
	if err := (&opts).Decode(cd.Opts); err != nil {
		return fmt.Errorf("failed to decode opts; got: '%s'; %w", cd.Opts, err)
	}
	if len(opts.ABI) != len(cd.Streams) {
		return fmt.Errorf("ABI length mismatch; expected: %d, got: %d", len(cd.Streams), len(opts.ABI))
	}
	return nil
}

func (rc ReportCodecEVMStreamlined) Pack(digest types.ConfigDigest, seqNr uint64, report ocr2types.Report, sigs []ocr2types.AttributedOnchainSignature) ([]byte, error) {
	payload, err := rc.packPayload(digest, report, sigs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode report payload; %w", err)
	}
	p := &LLOEVMStreamlinedReportWithContext{
		ConfigDigest:  digest[:],
		SeqNr:         seqNr,
		PackedPayload: payload,
	}
	b, err := proto.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report with context; %w", err)
	}
	return b, nil
}

// Equivalent to abi.encodePacked(digest, len(report), report, len(sigs), sigs...)
// bytes32 config digest
// packed uint16 len report
// packed bytes report
// packed uint8 num sigs
// packed bytes sigs
func (rc ReportCodecEVMStreamlined) packPayload(digest types.ConfigDigest, report []byte, sigs []ocr2types.AttributedOnchainSignature) ([]byte, error) {
	buf := new(bytes.Buffer)
	// Encode the digest
	if _, err := buf.Write(digest[:]); err != nil {
		return nil, fmt.Errorf("failed to encode config digest; %w", err)
	}

	// Encode the sequence number
	if len(report) > math.MaxUint16 {
		return nil, fmt.Errorf("report length exceeds maximum size; got %d, max %d", len(report), math.MaxUint16)
	}
	lenReportBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenReportBytes, uint16(len(report))) //nolint:gosec // G115 false positive, its checked above
	if _, err := buf.Write(lenReportBytes); err != nil {
		return nil, fmt.Errorf("failed to encode report length; %w", err)
	}

	if _, err := buf.Write(report); err != nil {
		return nil, fmt.Errorf("failed to encode report; %w", err)
	}

	// Encode the number of signers
	if len(sigs) > math.MaxUint8 {
		return nil, fmt.Errorf("expected at most 255 signers; got %d", len(sigs))
	}
	numSigners := uint8(len(sigs)) //nolint:gosec // G115 false positive, it's checked right above and we don't expect sigs to mutate
	if err := buf.WriteByte(numSigners); err != nil {
		return nil, fmt.Errorf("failed to encode number of signers; %w", err)
	}

	// Encode the signatures
	for _, sig := range sigs {
		if len(sig.Signature) != 65 {
			return nil, fmt.Errorf("expected signature length of 65 bytes; got %d", len(sig.Signature))
		}
		if _, err := buf.Write(sig.Signature); err != nil {
			return nil, fmt.Errorf("failed to encode signature; %w", err)
		}
	}
	return buf.Bytes(), nil
}

type ReportFormatEVMStreamlinedOpts struct {
	// FeedID is optional. If unspecified, the uint32 channel ID will be used
	// instead to identify the payload.
	FeedID *common.Hash `json:"feedID,omitempty"`
	// ABI defines the encoding of the payload. Each element maps to exactly
	// one stream (although sub-arrays may be specified for streams that
	// produce a composite data type).
	//
	// EXAMPLE
	//
	// [{"multiplier":"10000","type":"uint192"}, ...]
	//
	// See definition of ABIEncoder struct for more details.
	//
	// The total number of streams must be n, where n is the number of
	// top-level elements in this ABI array
	ABI []ABIEncoder `json:"abi"`
}

func (r *ReportFormatEVMStreamlinedOpts) Decode(opts []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(opts))
	decoder.DisallowUnknownFields() // Error on unrecognized fields
	return decoder.Decode(r)
}

func (r *ReportFormatEVMStreamlinedOpts) Encode() ([]byte, error) {
	return json.Marshal(r)
}
