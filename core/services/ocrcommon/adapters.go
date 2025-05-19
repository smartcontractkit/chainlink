package ocrcommon

import (
	"bytes"
	"cmp"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"slices"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
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
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report)
}

func (k *OCR3OnchainKeyringAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	return k.o.Verify(opk, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
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
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report, signatures)
}

func (c *OCR3ContractTransmitterAdapter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.ct.FromAccount(ctx)
}

var _ ocr3types.OnchainKeyring[[]byte] = (*OCR3OnchainKeyringMultiChainAdapter)(nil)

func MarshalMultichainKeyBundle(ost map[string]ocr2key.KeyBundle) (ocrtypes.OnchainPublicKey, error) {
	pubKeys := map[string]ocrtypes.OnchainPublicKey{}
	for k, b := range ost {
		pubKeys[k] = []byte(b.PublicKey())
	}
	return MarshalMultichainPublicKey(pubKeys)
}

func MarshalMultichainPublicKey(ost map[string]ocrtypes.OnchainPublicKey) (ocrtypes.OnchainPublicKey, error) {
	var pubKeys [][]byte
	for k, pubKey := range ost {
		typ, err := chaintype.ChainType(k).Type()
		if err != nil {
			// skipping unknown key type
			continue
		}
		buf := new(bytes.Buffer)
		if err = binary.Write(buf, binary.LittleEndian, typ); err != nil {
			return nil, err
		}
		length := len(pubKey)
		if length < 0 || length > math.MaxUint16 {
			return nil, errors.New("pubKey doesn't fit into uint16")
		}
		if err = binary.Write(buf, binary.LittleEndian, uint16(length)); err != nil {
			return nil, err
		}
		_, _ = buf.Write(pubKey)
		pubKeys = append(pubKeys, buf.Bytes())
	}
	// sort keys based on encoded type to make encoding deterministic
	slices.SortFunc(pubKeys, func(a, b []byte) int { return cmp.Compare(a[0], b[0]) })
	return bytes.Join(pubKeys, nil), nil
}

func UnmarshalMultichainPublicKey(d []byte) (map[string]ocrtypes.OnchainPublicKey, error) {
	m := map[string]ocrtypes.OnchainPublicKey{}
	buf := bytes.NewReader(d)

	for {
		// type
		typ, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		// length
		var length uint16
		err = binary.Read(buf, binary.LittleEndian, &length)
		if err != nil {
			return nil, err
		}
		// value
		pubKey := make([]byte, length)
		n, err := buf.Read(pubKey)
		if err != nil {
			return nil, err
		}
		if n != int(length) {
			return nil, io.EOF
		}

		k, err := chaintype.NewChainType(typ)
		if err != nil {
			// skipping unknown key type
			continue
		}
		m[string(k)] = pubKey

		if buf.Len() == 0 {
			break
		}
	}

	return m, nil
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
			Epoch:        uint32(seqNr),
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
	return kb.Verify(publicKey, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
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
