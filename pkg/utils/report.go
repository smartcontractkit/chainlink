package utils

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// RawReportContext is a copy of evmutil.RawReportContext to avoid importing go-ethereum.
// github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil#RawReportContext
func RawReportContext(repctx types.ReportContext) [3][32]byte {
	rawRepctx := [3][32]byte{}
	copy(rawRepctx[0][:], repctx.ConfigDigest[:])
	binary.BigEndian.PutUint32(rawRepctx[1][32-5:32-1], repctx.Epoch)
	rawRepctx[1][31] = repctx.Round
	rawRepctx[2] = repctx.ExtraHash
	return rawRepctx
}

// HashReport returns a report digest using SHA256 hash.
func HashReport(ctx types.ReportContext, r types.Report) ([]byte, error) {
	rawCtx := RawReportContext(ctx)
	buf := sha256.New()
	for _, v := range [][]byte{r[:], rawCtx[0][:], rawCtx[1][:], rawCtx[2][:]} {
		if _, err := buf.Write(v); err != nil {
			return []byte{}, err
		}
	}

	return buf.Sum(nil), nil
}
