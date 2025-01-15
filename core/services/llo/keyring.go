package llo

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/exp/maps"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
)

type LLOOnchainKeyring ocr3types.OnchainKeyring[llotypes.ReportInfo]

var _ LLOOnchainKeyring = &onchainKeyring{}

type Key interface {
	// Legacy Sign/Verify methods needed for v0.3 report compatibility
	// New keys can leave these stubbed
	Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error)
	Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool

	Sign3(digest ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error)
	Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool
	PublicKey() ocrtypes.OnchainPublicKey
	MaxSignatureLength() int
}

type onchainKeyring struct {
	lggr  logger.Logger
	keys  map[llotypes.ReportFormat]Key
	donID uint32
}

func NewOnchainKeyring(lggr logger.Logger, keys map[llotypes.ReportFormat]Key, donID uint32) LLOOnchainKeyring {
	return &onchainKeyring{
		lggr.Named("OnchainKeyring"), keys, donID,
	}
}

func (okr *onchainKeyring) PublicKey() types.OnchainPublicKey {
	// All unique public keys sorted in ascending order and combined into one
	// byte string
	onchainPublicKey := []byte{}

	keys := maps.Values(okr.keys)
	if len(keys) == 0 {
		return onchainPublicKey
	}
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].PublicKey(), keys[j].PublicKey()) < 0
	})

	onchainPublicKey = append(onchainPublicKey, keys[0].PublicKey()...)
	if len(keys) == 1 {
		return onchainPublicKey
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] != keys[i-1] {
			onchainPublicKey = append(onchainPublicKey, keys[i].PublicKey()...)
		}
	}

	return onchainPublicKey
}

func (okr *onchainKeyring) MaxSignatureLength() (n int) {
	// Needs to be max of all chain sigs
	for _, k := range okr.keys {
		n += k.MaxSignatureLength()
	}
	return
}

func (okr *onchainKeyring) Sign(digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo]) (signature []byte, err error) {
	switch r.Info.ReportFormat {
	case llotypes.ReportFormatEVMPremiumLegacy, llotypes.ReportFormatEVMABIEncodeUnpacked:
		rf := r.Info.ReportFormat
		if key, exists := okr.keys[rf]; exists {
			// NOTE: Must use legacy Sign method for compatibility with v0.3 report verification
			rc, err := evm.LegacyReportContext(digest, seqNr, okr.donID)
			if err != nil {
				return nil, fmt.Errorf("failed to get legacy report context: %w", err)
			}
			return key.Sign(rc, r.Report)
		}
	default:
		rf := r.Info.ReportFormat
		if key, exists := okr.keys[rf]; exists {
			return key.Sign3(digest, seqNr, r.Report)
		}
	}
	return nil, fmt.Errorf("Sign failed; unsupported report format: %q", r.Info.ReportFormat)
}

func (okr *onchainKeyring) Verify(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool {
	switch r.Info.ReportFormat {
	case llotypes.ReportFormatEVMPremiumLegacy, llotypes.ReportFormatEVMABIEncodeUnpacked:
		rf := r.Info.ReportFormat
		if verifier, exists := okr.keys[rf]; exists {
			// NOTE: Must use legacy Verify method for compatibility with v0.3 report verification
			rc, err := evm.LegacyReportContext(digest, seqNr, okr.donID)
			if err != nil {
				okr.lggr.Errorw("Verify failed; unable to get legacy report context", "err", err)
				return false
			}
			return verifier.Verify(key, rc, r.Report, signature)
		}
	default:
		rf := r.Info.ReportFormat
		if verifier, exists := okr.keys[rf]; exists {
			return verifier.Verify3(key, digest, seqNr, r.Report, signature)
		}
	}
	okr.lggr.Errorf("Verify failed; unsupported report format: %q", r.Info.ReportFormat)
	return false
}
