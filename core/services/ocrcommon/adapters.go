package ocrcommon

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var _ ocr3types.OnchainKeyring[[]byte] = (*OCR3OnchainKeyringAdapter)(nil)

type OCR3OnchainKeyringAdapter struct {
	o ocrtypes.OnchainKeyring
}

func NewOCR3OnchainKeyringAdapter(o ocrtypes.OnchainKeyring) *OCR3OnchainKeyringAdapter {
	return &OCR3OnchainKeyringAdapter{o}
}

func (k *OCR3OnchainKeyringAdapter) PublicKey() ocrtypes.OnchainPublicKey {
	return k.o.PublicKey()
}

func (k *OCR3OnchainKeyringAdapter) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte]) (signature []byte, err error) {
	return k.o.Sign(ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr), //nolint:gosec // seqNr maps to the uint32 epoch field
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report)
}

func (k *OCR3OnchainKeyringAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	return k.o.Verify(opk, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr), //nolint:gosec // seqNr maps to the uint32 epoch field
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, ri.Report, signature)
}

func (k *OCR3OnchainKeyringAdapter) MaxSignatureLength() int {
	return k.o.MaxSignatureLength()
}

var _ ocr3types.ContractTransmitter[[]byte] = (*OCR3ContractTransmitterAdapter)(nil)

type OCR3ContractTransmitterAdapter struct {
	ct ocrtypes.ContractTransmitter
}

func NewOCR3ContractTransmitterAdapter(ct ocrtypes.ContractTransmitter) *OCR3ContractTransmitterAdapter {
	return &OCR3ContractTransmitterAdapter{ct}
}

func (c *OCR3ContractTransmitterAdapter) Transmit(ctx context.Context, digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte], signatures []ocrtypes.AttributedOnchainSignature) error {
	return c.ct.Transmit(ctx, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr), //nolint:gosec // seqNr maps to the uint32 epoch field
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report, signatures)
}

func (c *OCR3ContractTransmitterAdapter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.ct.FromAccount(ctx)
}

var (
	_ ocr3types.OnchainKeyring[[]byte]  = (*OCR3OnchainKeyringMultiChainAdapter)(nil)
	_ ocr3types.OnchainKeyring2[[]byte] = (*OCR3OnchainKeyringMultiChainAdapter)(nil)
)

func MarshalMultichainKeyBundle(ost map[string]ocr2key.KeyBundle) (ocrtypes.OnchainPublicKey, error) {
	pubKeys := map[string]ocrtypes.OnchainPublicKey{}
	for k, b := range ost {
		pubKeys[k] = []byte(b.PublicKey())
	}
	return MarshalMultichainPublicKey(pubKeys)
}

func MarshalMultichainPublicKey(ost map[string]ocrtypes.OnchainPublicKey) (ocrtypes.OnchainPublicKey, error) {
	return ocr2key.MultichainPublicKey(ost).MarshalBinary()
}

type OCR3OnchainKeyringMultiChainAdapter struct {
	keyBundles map[string]ocr2key.KeyBundle
	publicKey  ocrtypes.OnchainPublicKey
	lggr       logger.Logger
}

func NewOCR3OnchainKeyringMultiChainAdapter(ost map[string]ocr2key.KeyBundle, lggr logger.Logger) (*OCR3OnchainKeyringMultiChainAdapter, error) {
	if len(ost) == 0 {
		return nil, errors.New("no key bundles provided")
	}
	publicKey, err := MarshalMultichainKeyBundle(ost)
	if err != nil {
		return nil, err
	}
	return &OCR3OnchainKeyringMultiChainAdapter{ost, publicKey, lggr}, nil
}

func (a *OCR3OnchainKeyringMultiChainAdapter) PublicKey() ocrtypes.OnchainPublicKey {
	return a.publicKey
}

// Has returns true if every chain public key in k matches this adapter (k may be a subset of supported chains).
func (a *OCR3OnchainKeyringMultiChainAdapter) Has(k ocrtypes.OnchainPublicKey) bool {
	if len(k) == 0 {
		return false
	}
	var keys ocr2key.MultichainPublicKey
	err := keys.UnmarshalBinary(k)
	if err != nil {
		return bytes.Equal(k, a.publicKey)
	}
	if len(keys) == 0 {
		return false
	}
	for chainName, pubKey := range keys {
		kb, ok := a.keyBundles[chainName]
		if !ok {
			return false
		}
		if !bytes.Equal(pubKey, kb.PublicKey()) {
			return false
		}
	}
	return true
}

func (a *OCR3OnchainKeyringMultiChainAdapter) DebugIdentifier() string {
	return fmt.Sprintf("%x", a.publicKey)
}

func (a *OCR3OnchainKeyringMultiChainAdapter) getKeyBundleFromInfo(info []byte) (string, ocr2key.KeyBundle, error) {
	unmarshalledInfo := new(structpb.Struct)
	err := proto.Unmarshal(info, unmarshalledInfo)
	if err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal report info: %w", err)
	}
	infoMap := unmarshalledInfo.AsMap()
	keyBundleName, ok := infoMap["keyBundleName"]
	if !ok {
		return "", nil, errors.New("keyBundleName not found in report info")
	}
	name, ok := keyBundleName.(string)
	if !ok {
		return "", nil, errors.New("keyBundleName is not a string")
	}
	kb, ok := a.keyBundles[name]
	if !ok {
		return "", nil, fmt.Errorf("keyBundle not found: %s", name)
	}
	return name, kb, nil
}

func (a *OCR3OnchainKeyringMultiChainAdapter) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte]) (signature []byte, err error) {
	_, kb, err := a.getKeyBundleFromInfo(r.Info)
	if err != nil {
		return nil, fmt.Errorf("sign: failed to get key bundle from report info: %w", err)
	}
	return kb.Sign(ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr), //nolint:gosec // seqNr maps to the uint32 epoch field
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report)
}

func (a *OCR3OnchainKeyringMultiChainAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	kbName, kb, err := a.getKeyBundleFromInfo(ri.Info)
	if err != nil {
		a.lggr.Warnf("verify: failed to get key bundle from report info: %v", err)
		return false
	}
	var keys ocr2key.MultichainPublicKey
	err = keys.UnmarshalBinary(opk)
	if err != nil {
		a.lggr.Warnf("verify: failed to unmarshal public keys: %v", err)
		return false
	}
	publicKey, ok := keys[kbName]
	if !ok {
		a.lggr.Warnf("verify: publicKey not found: %v", kbName)
		return false
	}
	return kb.Verify(publicKey, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr), //nolint:gosec // seqNr maps to the uint32 epoch field
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, ri.Report, signature)
}

func (a *OCR3OnchainKeyringMultiChainAdapter) MaxSignatureLength() int {
	maxLength := -1
	for _, kb := range a.keyBundles {
		l := kb.MaxSignatureLength()
		if l > maxLength {
			maxLength = l
		}
	}
	return maxLength
}
