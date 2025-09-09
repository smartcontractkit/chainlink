package ccipaptos

import (
	"context"
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with ccip::offramp version 1.6.0
type CommitPluginCodecV1 struct{}

var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	s := &bcs.Serializer{}
	bcs.SerializeSequenceWithFunction(report.PriceUpdates.TokenPriceUpdates, s, func(s *bcs.Serializer, item cciptypes.TokenPrice) {
		sourceToken := aptos.AccountAddress{}
		err := sourceToken.ParseStringRelaxed(string(item.TokenID))
		if err != nil {
			s.SetError(fmt.Errorf("failed to parse source token address: %w", err))
			return
		}
		s.Struct(&sourceToken)
		if item.Price.IsEmpty() {
			s.SetError(fmt.Errorf("empty price for token: %s", item.TokenID))
			return
		}

		s.U256(*item.Price.Int)
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize TokenPriceUpdates: %w", s.Error())
	}
	bcs.SerializeSequenceWithFunction(report.PriceUpdates.GasPriceUpdates, s, func(s *bcs.Serializer, item cciptypes.GasPriceChain) {
		s.U64(uint64(item.ChainSel))
		if item.GasPrice.IsEmpty() {
			s.SetError(fmt.Errorf("empty gas price for chain: %d", item.ChainSel))
			return
		}
		s.U256(*item.GasPrice.Int)
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize GasPriceUpdates: %w", s.Error())
	}
	bcs.SerializeSequenceWithFunction(report.BlessedMerkleRoots, s, func(s *bcs.Serializer, item cciptypes.MerkleRootChain) {
		s.U64(uint64(item.ChainSel))
		s.WriteBytes(item.OnRampAddress[:])
		s.U64(uint64(item.SeqNumsRange.Start()))
		s.U64(uint64(item.SeqNumsRange.End()))
		s.FixedBytes(item.MerkleRoot[:])
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize BlessedMerkleRoots: %w", s.Error())
	}
	bcs.SerializeSequenceWithFunction(report.UnblessedMerkleRoots, s, func(s *bcs.Serializer, item cciptypes.MerkleRootChain) {
		s.U64(uint64(item.ChainSel))
		s.WriteBytes(item.OnRampAddress[:])
		s.U64(uint64(item.SeqNumsRange.Start()))
		s.U64(uint64(item.SeqNumsRange.End()))
		s.FixedBytes(item.MerkleRoot[:])
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize UnblessedMerkleRoots: %w", s.Error())
	}
	bcs.SerializeSequenceWithFunction(report.RMNSignatures, s, func(s *bcs.Serializer, item cciptypes.RMNECDSASignature) {
		s.FixedBytes(item.R[:])
		s.FixedBytes(item.S[:])
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize RMNSignatures: %w", s.Error())
	}

	return s.ToBytes(), nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, data []byte) (cciptypes.CommitPluginReport, error) {
	des := bcs.NewDeserializer(data)
	report := cciptypes.CommitPluginReport{}

	report.PriceUpdates.TokenPriceUpdates = bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *cciptypes.TokenPrice) {
		var sourceToken aptos.AccountAddress
		des.Struct(&sourceToken)
		if des.Error() != nil {
			return
		}
		// StringLong() instead of String() to standardize encoding across all addresses,
		// since String() shortens system addresses.
		item.TokenID = cciptypes.UnknownEncodedAddress(sourceToken.StringLong())
		price := des.U256()
		if des.Error() != nil {
			return
		}

		// we need this clause because the zero token price test fails otherwise:
		// -      abs: (big.nat) <nil>
		// +      abs: (big.nat) {
		// +      }
		// the reason is because big.NewInt(0) ends up not setting the `abs` field at all, while big.NewInt().SetBytes(..) will
		// set the `abs` value to 0.
		// ref: https://cs.opensource.google/go/go/+/master:src/math/big/int.go;drc=432fd9c60fac4485d0473173171206f1ef558829;l=85
		if price.Sign() == 0 {
			item.Price = cciptypes.NewBigInt(big.NewInt(0))
		} else {
			item.Price = cciptypes.NewBigInt(&price)
		}
	})

	if des.Error() != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to deserialize TokenPriceUpdates: %w", des.Error())
	}

	report.PriceUpdates.GasPriceUpdates = bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *cciptypes.GasPriceChain) {
		item.ChainSel = cciptypes.ChainSelector(des.U64())
		if des.Error() != nil {
			return
		}
		gasPrice := des.U256()
		if des.Error() != nil {
			return
		}
		if gasPrice.Sign() == 0 {
			item.GasPrice = cciptypes.NewBigInt(big.NewInt(0))
		} else {
			item.GasPrice = cciptypes.NewBigInt(&gasPrice)
		}
	})
	if des.Error() != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to deserialize GasPriceUpdates: %w", des.Error())
	}

	deserializeMerkleRootChain := func(des *bcs.Deserializer, item *cciptypes.MerkleRootChain) {
		item.ChainSel = cciptypes.ChainSelector(des.U64())
		if des.Error() != nil {
			return
		}
		onRampAddrBytes := des.ReadBytes()
		if des.Error() != nil {
			return
		}
		item.OnRampAddress = onRampAddrBytes
		startSeqNum := des.U64()
		if des.Error() != nil {
			return
		}
		endSeqNum := des.U64()
		if des.Error() != nil {
			return
		}
		item.SeqNumsRange = cciptypes.NewSeqNumRange(cciptypes.SeqNum(startSeqNum), cciptypes.SeqNum(endSeqNum))
		des.ReadFixedBytesInto(item.MerkleRoot[:])
		if des.Error() != nil {
			return
		}
	}

	report.BlessedMerkleRoots = bcs.DeserializeSequenceWithFunction(des, deserializeMerkleRootChain)
	if des.Error() != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to deserialize BlessedMerkleRoots: %w", des.Error())
	}

	report.UnblessedMerkleRoots = bcs.DeserializeSequenceWithFunction(des, deserializeMerkleRootChain)
	if des.Error() != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to deserialize UnblessedMerkleRoots: %w", des.Error())
	}

	report.RMNSignatures = bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *cciptypes.RMNECDSASignature) {
		des.ReadFixedBytesInto(item.R[:])
		if des.Error() != nil {
			return
		}
		des.ReadFixedBytesInto(item.S[:])
		if des.Error() != nil {
			return
		}
	})
	if des.Error() != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to deserialize RMNSignatures: %w", des.Error())
	}

	if des.Remaining() > 0 {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("unexpected remaining bytes after decoding: %d", des.Remaining())
	}

	return report, nil
}
