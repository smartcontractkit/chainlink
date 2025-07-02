package keyring

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// DummyEVMOnchainKeyring is not intended to be used in production.
type DummyEVMOnchainKeyring struct {
	PrivateKey ecdsa.PrivateKey
}

var _ ocr3types.OnchainKeyring[por.ChainSelector] = &DummyEVMOnchainKeyring{}

func (ring *DummyEVMOnchainKeyring) MaxSignatureLength() int {
	return 65
}

func (ring *DummyEVMOnchainKeyring) Sign(configDigest types.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[por.ChainSelector]) (signature []byte, err error) {
	sigData := crypto.Keccak256(reportWithInfo.Report)
	sigData = append(sigData, configDigest[:]...)
	sigData = binary.BigEndian.AppendUint64(sigData, seqNr)
	return crypto.Sign(crypto.Keccak256(sigData), &ring.PrivateKey)
}

func (ring *DummyEVMOnchainKeyring) PublicKey() types.OnchainPublicKey {
	address := crypto.PubkeyToAddress(ring.PrivateKey.PublicKey)
	return address[:]
}

func (ring *DummyEVMOnchainKeyring) Verify(pubkey types.OnchainPublicKey, configDigest types.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[por.ChainSelector], sig []byte) bool {
	sigData := crypto.Keccak256(reportWithInfo.Report)
	sigData = append(sigData, configDigest[:]...)
	sigData = binary.BigEndian.AppendUint64(sigData, seqNr)
	hash := crypto.Keccak256(sigData)
	authorPubkey, err := crypto.SigToPub(hash, sig)

	// fmt.Printf("author pubkey %x\n", authorPubkey)
	if err != nil {
		// fmt.Printf("error while doing SigToPub: %v\n", err)
		return false
	}
	authorAddress := crypto.PubkeyToAddress(*authorPubkey)
	// fmt.Printf("author address %x\n", authorAddress)
	// fmt.Printf("expected address %x\n", common.BytesToAddress(pubkey))
	return bytes.Equal(pubkey[:], authorAddress[:])
}
