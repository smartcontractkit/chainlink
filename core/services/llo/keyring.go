package llo

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type LLOOnchainKeyring ocr3types.OnchainKeyring[llotypes.ReportInfo]

var _ LLOOnchainKeyring = &onchainKeyring{}

type Key interface {
	Sign3(digest ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error)
	Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool
	PublicKey() ocrtypes.OnchainPublicKey
	MaxSignatureLength() int
}

type onchainKeyring struct {
	lggr logger.Logger
	keys map[llotypes.ReportFormat]Key
}

func NewOnchainKeyring(lggr logger.Logger, keys map[llotypes.ReportFormat]Key) LLOOnchainKeyring {
	return &onchainKeyring{
		lggr.Named("OnchainKeyring"), keys,
	}
}

func (okr *onchainKeyring) PublicKey() types.OnchainPublicKey {
	// All public keys combined
	var pk []byte
	for _, k := range okr.keys {
		pk = append(pk, k.PublicKey()...)
	}
	return pk
}

func (okr *onchainKeyring) MaxSignatureLength() (n int) {
	// Needs to be max of all chain sigs
	for _, k := range okr.keys {
		n += k.MaxSignatureLength()
	}
	return
}

func (okr *onchainKeyring) Sign(digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo]) (signature []byte, err error) {
	rf := r.Info.ReportFormat
	// HACK: sign/verify JSON payloads with EVM keys for now, this makes
	// debugging and testing easier
	if rf == llotypes.ReportFormatJSON {
		rf = llotypes.ReportFormatEVM
	}
	if key, exists := okr.keys[rf]; exists {
		return key.Sign3(digest, seqNr, r.Report)
	}
	return nil, fmt.Errorf("Sign failed; unsupported report format: %q", r.Info.ReportFormat)
}

func (okr *onchainKeyring) Verify(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool {
	rf := r.Info.ReportFormat
	// HACK: sign/verify JSON payloads with EVM keys for now, this makes
	// debugging and testing easier
	if rf == llotypes.ReportFormatJSON {
		rf = llotypes.ReportFormatEVM
	}
	if verifier, exists := okr.keys[rf]; exists {
		return verifier.Verify3(key, digest, seqNr, r.Report, signature)
	}
	okr.lggr.Errorf("Verify failed; unsupported report format: %q", r.Info.ReportFormat)
	return false
}
