package ocr3impls

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

var _ ocr3types.OnchainKeyring[models.ReportMetadata] = &onchainKeyringV3Wrapper[models.ReportMetadata]{}

type onchainKeyringV3Wrapper[RI any] struct {
	core types.OnchainKeyring
}

func NewOnchainKeyring[RI any](keyring types.OnchainKeyring) *onchainKeyringV3Wrapper[RI] {
	return &onchainKeyringV3Wrapper[RI]{
		core: keyring,
	}
}

func (w *onchainKeyringV3Wrapper[RI]) PublicKey() types.OnchainPublicKey {
	return w.core.PublicKey()
}

func (w *onchainKeyringV3Wrapper[RI]) MaxSignatureLength() int {
	return w.core.MaxSignatureLength()
}

func (w *onchainKeyringV3Wrapper[RI]) Sign(digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[RI]) (signature []byte, err error) {
	epoch, round := uint64ToUint32AndUint8(seqNr)
	rCtx := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        epoch,
			Round:        round,
		},
	}

	return w.core.Sign(rCtx, r.Report)
}

func (w *onchainKeyringV3Wrapper[RI]) Verify(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[RI], signature []byte) bool {
	epoch, round := uint64ToUint32AndUint8(seqNr)
	rCtx := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        epoch,
			Round:        round,
		},
	}

	return w.core.Verify(key, rCtx, r.Report, signature)
}

func uint64ToUint32AndUint8(x uint64) (uint32, uint8) {
	return uint32(x >> 32), uint8(x)
}
