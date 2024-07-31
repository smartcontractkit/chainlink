package ocr3impls

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

var _ ocr3types.OnchainKeyring[models.Report] = &onchainKeyringV3Wrapper[models.Report]{}

type onchainKeyringV3Wrapper[RI MultichainMeta] struct {
	core types.OnchainKeyring
	lggr logger.Logger
}

func NewOnchainKeyring[RI MultichainMeta](keyring types.OnchainKeyring, lggr logger.Logger) *onchainKeyringV3Wrapper[RI] {
	return &onchainKeyringV3Wrapper[RI]{
		core: keyring,
		lggr: lggr.Named("MultichainOnchainKeyring"),
	}
}

func (w *onchainKeyringV3Wrapper[RI]) PublicKey() types.OnchainPublicKey {
	return w.core.PublicKey()
}

func (w *onchainKeyringV3Wrapper[RI]) MaxSignatureLength() int {
	return w.core.MaxSignatureLength()
}

func (w *onchainKeyringV3Wrapper[RI]) Sign(_ types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[RI]) (signature []byte, err error) {
	// the provided config digest is from the master chain
	// we need to sign with the config digest for the chain that we are transmitting to
	// which may be any chain.
	configDigest := r.Info.GetDestinationConfigDigest()

	epoch, round := uint64ToUint32AndUint8(seqNr)
	rCtx := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}

	w.lggr.Debugw("signing report", "configDigest", configDigest.Hex(), "seqNr", seqNr, "report", hexutil.Encode(r.Report))

	return w.core.Sign(rCtx, r.Report)
}

func (w *onchainKeyringV3Wrapper[RI]) Verify(key types.OnchainPublicKey, _ types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[RI], signature []byte) bool {
	// the provided config digest is from the master chain
	// we need to sign with the config digest for the chain that we are transmitting to
	// which may be any chain.
	configDigest := r.Info.GetDestinationConfigDigest()

	epoch, round := uint64ToUint32AndUint8(seqNr)
	rCtx := types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}

	w.lggr.Debugw("verifying report", "configDigest", configDigest.Hex(), "seqNr", seqNr, "report", hexutil.Encode(r.Report))

	return w.core.Verify(key, rCtx, r.Report, signature)
}

func uint64ToUint32AndUint8(x uint64) (uint32, uint8) {
	return uint32(x >> 32), uint8(x)
}
