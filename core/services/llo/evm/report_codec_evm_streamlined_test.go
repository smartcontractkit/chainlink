package evm

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

func TestReportCodecEVMStreamlined(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	codec := ReportCodecEVMStreamlined{}

	t.Run("Encode", func(t *testing.T) {
		t.Run("one value, without feed ID - fits into one evm word", func(t *testing.T) {
			cd := llotypes.ChannelDefinition{
				ReportFormat: 42,
				Opts:         []byte(`{"abi":[{"type":"int128"}]}`),
			}
			payload, err := codec.Encode(ctx, llo.Report{
				ChannelID:             1,
				ValidAfterNanoseconds: 1234567890,
				Values: []llo.StreamValue{
					llo.ToDecimal(decimal.NewFromFloat(1123455935.123)),
				},
			}, cd)
			require.NoError(t, err)
			require.Len(t, payload, 32)
			// Report Format
			assert.Equal(t, "0000002a", hex.EncodeToString(payload[:4]))
			// Channel ID
			assert.Equal(t, "00000001", hex.EncodeToString(payload[4:8]))
			// Timestamp
			assert.Equal(t, "00000000499602d2", hex.EncodeToString(payload[8:16]))
			// Value
			assert.Equal(t, "00000000000000000000000042f693bf", hex.EncodeToString(payload[16:]))
		})
		t.Run("one value, with feed ID - fits into two evm words", func(t *testing.T) {
			feedID := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
			cd := llotypes.ChannelDefinition{
				Opts: []byte(fmt.Sprintf(`{"abi":[{"type":"int192"}], "feedID":"0x%s"}`, feedID)),
			}
			payload, err := codec.Encode(ctx, llo.Report{
				ChannelID:             1,
				ValidAfterNanoseconds: 1234567890,
				Values: []llo.StreamValue{
					llo.ToDecimal(decimal.NewFromFloat(1123455935.123)),
				},
			}, cd)
			require.NoError(t, err)
			require.Len(t, payload, 64)
			assert.Equal(t, feedID, hex.EncodeToString(payload[:32]))                                             // feed id
			assert.Equal(t, "00000000499602d2", hex.EncodeToString(payload[32:40]))                               // timestamp
			assert.Equal(t, "000000000000000000000000000000000000000042f693bf", hex.EncodeToString(payload[40:])) // value
		})
	})
	t.Run("Verify", func(t *testing.T) {
		t.Run("with valid opts", func(t *testing.T) {
			err := codec.Verify(ctx, llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 123, Aggregator: llotypes.AggregatorMedian},
				},
				Opts: []byte(`{"abi":[{"type":"int160"}]}`),
			})
			require.NoError(t, err)
			err = codec.Verify(ctx, llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 123, Aggregator: llotypes.AggregatorMedian},
				},
				Opts: []byte(`{"feedID":"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef","abi":[{"type":"int160"}]}`),
			})
			require.NoError(t, err)
			t.Run("with invalid opts", func(t *testing.T) {
				err := codec.Verify(ctx, llotypes.ChannelDefinition{
					Opts: []byte(`{"abi":[{"type":"int160"}]}`),
				})
				require.EqualError(t, err, "ABI length mismatch; expected: 0, got: 1")
				err = codec.Verify(ctx, llotypes.ChannelDefinition{
					Streams: []llotypes.Stream{
						{StreamID: 123, Aggregator: llotypes.AggregatorMedian},
					},
				})
				require.EqualError(t, err, "failed to decode opts; got: ''; EOF")
				err = codec.Verify(ctx, llotypes.ChannelDefinition{
					Opts: []byte(`{"feedID":"0xinvalid","abi":[{"type":"int160"}]}`),
				})
				require.EqualError(t, err, "failed to decode opts; got: '{\"feedID\":\"0xinvalid\",\"abi\":[{\"type\":\"int160\"}]}'; json: cannot unmarshal hex string of odd length into Go struct field ReportFormatEVMStreamlinedOpts.feedID of type common.Hash")
			})
		})
	})
	t.Run("Pack", func(t *testing.T) {
		validSig := make([]byte, 65)
		for i := range validSig {
			validSig[i] = byte(i)
		}

		t.Run("valid", func(t *testing.T) {
			cdBytes := make([]byte, 32)
			for i := range cdBytes {
				cdBytes[i] = byte(i)
			}
			cd, err := ocr2types.BytesToConfigDigest(cdBytes)
			require.NoError(t, err)
			packed, err := codec.Pack(
				cd,
				math.MaxUint32,
				[]byte("report"),
				[]ocr2types.AttributedOnchainSignature{{Signer: 1, Signature: validSig}},
			)
			require.NoError(t, err)
			p := &LLOEVMStreamlinedReportWithContext{}
			err = proto.Unmarshal(packed, p)
			require.NoError(t, err)
			assert.Equal(t, cdBytes, p.ConfigDigest)
			assert.Equal(t, uint64(math.MaxUint32), p.SeqNr)
			assert.Equal(t, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"+ // digest
				"0006"+ // uint16 report len
				"7265706f7274"+ // report
				"01"+ // num sigs
				"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40", // 65 byte sig
				hex.EncodeToString(p.PackedPayload))
		})
		t.Run("invalid", func(t *testing.T) {
			// wrong sig length
			_, err := codec.Pack(
				ocr2types.ConfigDigest{0x01, 0x02, 0x03},
				123,
				[]byte("report"),
				[]ocr2types.AttributedOnchainSignature{{Signer: 1, Signature: []byte("sig")}},
			)
			require.EqualError(t, err, "failed to encode report payload; expected signature length of 65 bytes; got 3")
			// too many signers
			tooManySigners := make([]ocr2types.AttributedOnchainSignature, 256)
			_, err = codec.Pack(
				ocr2types.ConfigDigest{0x01, 0x02, 0x03},
				123,
				[]byte("report"),
				tooManySigners,
			)
			require.EqualError(t, err, "failed to encode report payload; expected at most 255 signers; got 256")
			// report too large
			largeReport := make([]byte, math.MaxUint16+1)
			_, err = codec.Pack(
				ocr2types.ConfigDigest{0x01, 0x02, 0x03},
				123,
				largeReport,
				[]ocr2types.AttributedOnchainSignature{{Signer: 1, Signature: validSig}},
			)
			require.EqualError(t, err, "failed to encode report payload; report length exceeds maximum size; got 65536, max 65535")
		})
	})
}
