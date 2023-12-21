package evm

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/automation"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

var _ ocr3types.OnchainKeyring[automation.AutomationReportInfo] = &onchainKeyringV3Wrapper{}

type onchainKeyringV3Wrapper struct {
	core types.OnchainKeyring
}

func NewOnchainKeyringV3Wrapper(keyring types.OnchainKeyring) *onchainKeyringV3Wrapper {
	return &onchainKeyringV3Wrapper{
		core: keyring,
	}
}

func (w *onchainKeyringV3Wrapper) PublicKey() types.OnchainPublicKey {
	return w.core.PublicKey()
}

func (w *onchainKeyringV3Wrapper) MaxSignatureLength() int {
	return w.core.MaxSignatureLength()
}

func (w *onchainKeyringV3Wrapper) Sign(digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[automation.AutomationReportInfo]) (signature []byte, err error) {
	rCtx := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
		},
	}

	return w.core.Sign(rCtx, r.Report)
}

func (w *onchainKeyringV3Wrapper) Verify(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[automation.AutomationReportInfo], signature []byte) bool {
	rCtx := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
		},
	}

	return w.core.Verify(key, rCtx, r.Report, signature)
}
