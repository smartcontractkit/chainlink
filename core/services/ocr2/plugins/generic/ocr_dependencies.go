package generic

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/ocr/keyring"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

type OCRKeyringsParams struct {
	Logger                 logger.Logger
	KB                     ocr2key.KeyBundle
	OnchainSigningStrategy job.OnchainSigningStrategy
	OcrKeystore            keystore.OCR2
}

// NewOCRKeyrings returns the OCR keys the node lends a capability that runs its own OCR oracle: those the oracle
// factory gives the oracles it makes.
func NewOCRKeyrings(params OCRKeyringsParams) (keyring.Keyrings, error) {
	keyBundles := map[string]ocr2key.KeyBundle{}
	for name, kbID := range params.OnchainSigningStrategy.Config {
		kb, err := params.OcrKeystore.Get(kbID)
		if err != nil {
			return keyring.Keyrings{}, fmt.Errorf("failed to get ocr key for key bundle ID '%s': %w", kbID, err)
		}

		keyBundles[name] = kb
	}

	onchainKeyring, err := ocrcommon.NewOCR3OnchainKeyringMultiChainAdapter(keyBundles, params.Logger)
	if err != nil {
		return keyring.Keyrings{}, fmt.Errorf("failed to instantiate onchain keyring with multi chain adapter: %w", err)
	}

	return keyring.Keyrings{OffchainKeyring: params.KB, OnchainKeyring: keyBundleKeyring{keyring: onchainKeyring}}, nil
}

// NewOCRNetwork returns the networking the node lends a capability that runs its own OCR oracle: its peer's. The peer
// wrapper only has a peer once it is started, so the peer is read when asked.
func NewOCRNetwork(peerWrapper *ocrcommon.SingletonPeerWrapper) func(context.Context) (ocrtypes.BinaryNetworkEndpointFactory, error) {
	return func(context.Context) (ocrtypes.BinaryNetworkEndpointFactory, error) {
		if !peerWrapper.IsStarted() {
			return nil, errors.New("peer wrapper not started")
		}

		return peerWrapper.Peer2, nil
	}
}

// keyBundleKeyring signs with the key bundle a report's info names, such as evm. The multi-chain adapter reads the
// name from a report's info as a protobuf Struct, so it is encoded as one.
type keyBundleKeyring struct {
	keyring ocr3types.OnchainKeyring[[]byte]
}

func (k keyBundleKeyring) PublicKey() ocrtypes.OnchainPublicKey { return k.keyring.PublicKey() }

func (k keyBundleKeyring) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, report ocr3types.ReportWithInfo[string]) ([]byte, error) {
	r, err := withKeyBundleInfo(report)
	if err != nil {
		return nil, err
	}

	return k.keyring.Sign(digest, seqNr, r)
}

func (k keyBundleKeyring) Verify(key ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64,
	report ocr3types.ReportWithInfo[string], signature []byte) bool {
	r, err := withKeyBundleInfo(report)
	return err == nil && k.keyring.Verify(key, digest, seqNr, r, signature)
}

func (k keyBundleKeyring) MaxSignatureLength() int { return k.keyring.MaxSignatureLength() }

func withKeyBundleInfo(report ocr3types.ReportWithInfo[string]) (ocr3types.ReportWithInfo[[]byte], error) {
	info, err := structpb.NewStruct(map[string]any{"keyBundleName": report.Info})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, err
	}

	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(info)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to encode report info: %w", err)
	}

	return ocr3types.ReportWithInfo[[]byte]{Report: report.Report, Info: b}, nil
}
