package mercury

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func buildSampleReport() []byte {
	feedID := [32]byte{'f', 'o', 'o'}
	timestamp := uint32(42)
	blockNumber := uint64(142)
	median := big.NewInt(242)

	b, err := reportTypes.Pack(feedID, timestamp, blockNumber, median)
	if err != nil {
		panic(err)
	}
	return b
}

func buildSamplePayload() []byte {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range sampleSigs {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := evmutil.RawReportContext(sampleReportContext)
	payload, err := payloadTypes.Pack(rawReportCtx, []byte(sampleReport), rs, ss, vs)
	if err != nil {
		panic(err)
	}
	return payload
}

var (
	sampleReport        = buildSampleReport()
	sampleReportHex     = hexutil.Encode(sampleReport)
	sampleFromAccount   = utils.RandomAddress()
	sig2                = ocrtypes.AttributedOnchainSignature{Signature: mustDecodeBase64("kbeuRczizOJCxBzj7MUAFpz3yl2WRM6K/f0ieEBvA+oTFUaKslbQey10krumVjzAvlvKxMfyZo0WkOgNyfF6xwE="), Signer: 2}
	sig3                = ocrtypes.AttributedOnchainSignature{Signature: mustDecodeBase64("9jz4b6Dh2WhXxQ97a6/S9UNjSfrEi9016XKTrfN0mLQFDiNuws23x7Z4n+6g0sqKH/hnxx1VukWUH/ohtw83/wE="), Signer: 3}
	sampleSigs          = []ocrtypes.AttributedOnchainSignature{sig2, sig3}
	sampleReportContext = ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: mustHexToConfigDigest("0x0001fc30092226b37f6924b464e16a54a7978a9a524519a73403af64d487dc45"),
			Epoch:        6,
			Round:        28,
		},
		ExtraHash: [32]uint8{27, 144, 106, 73, 166, 228, 123, 166, 179, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114},
	}
	samplePayload       = buildSamplePayload()
	samplePayloadHex    = hexutil.Encode(samplePayload)
	sampleMercuryReport = fmt.Sprintf(`{"Payload":"%s","FromAccount":"%s"}`, samplePayloadHex, strings.ToLower(sampleFromAccount.Hex()))
)

func mustDecodeBase64(s string) (b []byte) {
	var err error
	b, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustHexToConfigDigest(s string) (cd ocrtypes.ConfigDigest) {
	b := hexutil.MustDecode(s)
	var err error
	cd, err = ocrtypes.BytesToConfigDigest(b)
	if err != nil {
		panic(err)
	}
	return
}
