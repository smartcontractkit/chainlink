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
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	return k.o.Sign(ocr2key.OCR3ReportContext(digest, seqNr), r.Report)
}

func (k *OCR3OnchainKeyringAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	return k.o.Verify(opk, ocr2key.OCR3ReportContext(digest, seqNr), ri.Report, signature)
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
	return c.ct.Transmit(ctx, ocr2key.OCR3ReportContext(digest, seqNr), r.Report, signatures)
}

func (c *OCR3ContractTransmitterAdapter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.ct.FromAccount(ctx)
}

var _ ocr3types.OnchainKeyring[[]byte] = (*OCR3OnchainKeyringMultiChainAdapter)(nil)
var _ ocr3types.OnchainKeyring2[[]byte] = (*OCR3OnchainKeyringMultiChainAdapter)(nil)

// The multichain onchain public key encoding lives in chainlink-common, beside the
// key bundles whose public keys it encodes. An oracle joining a configuration and
// whoever wrote that configuration have to agree on it byte for byte, and they do
// not all share this package - a standalone capability signing on a node's behalf
// reads the same encoding without importing core. These remain as the names core
// reaches it by.
var (
	MarshalMultichainKeyBundle   = ocr2key.MarshalMultichainKeyBundle
	MarshalMultichainPublicKey   = ocr2key.MarshalMultichainPublicKey
	UnmarshalMultichainPublicKey = ocr2key.UnmarshalMultichainPublicKey
)

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
	keys, err := UnmarshalMultichainPublicKey(k)
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
	return kb.Sign(ocr2key.OCR3ReportContext(digest, seqNr), r.Report)
}

func (a *OCR3OnchainKeyringMultiChainAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	kbName, kb, err := a.getKeyBundleFromInfo(ri.Info)
	if err != nil {
		a.lggr.Warnf("verify: failed to get key bundle from report info: %v", err)
		return false
	}
	keys, err := UnmarshalMultichainPublicKey(opk)
	if err != nil {
		a.lggr.Warnf("verify: failed to unmarshal public keys: %v", err)
		return false
	}
	publicKey, ok := keys[kbName]
	if !ok {
		a.lggr.Warnf("verify: publicKey not found: %v", kbName)
		return false
	}
	return kb.Verify(publicKey, ocr2key.OCR3ReportContext(digest, seqNr), ri.Report, signature)
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
