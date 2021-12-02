package ocr2key

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"io"

	"github.com/smartcontractkit/chainlink/core/logger"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
)

var _ ocrtypes.OffchainKeyring = &OffchainKeyring{}

// OffchainKeyring contains the secret keys needed for the OCR nodes to share secrets
// and perform aggregation.
//
// This is currently an ed25519 signing key and a separate encryption key.
//
// All its functions should be thread-safe.
type OffchainKeyring struct {
	signingKey    ed25519.PrivateKey
	encryptionKey [curve25519.ScalarSize]byte
}

func newOffchainKeyring(encryptionMaterial, signingMaterial io.Reader) (*OffchainKeyring, error) {
	_, signingKey, err := ed25519.GenerateKey(signingMaterial)
	if err != nil {
		return nil, err
	}

	encryptionKey := [curve25519.ScalarSize]byte{}
	_, err = encryptionMaterial.Read(encryptionKey[:])
	if err != nil {
		return nil, err
	}

	ok := &OffchainKeyring{
		signingKey:    signingKey,
		encryptionKey: encryptionKey,
	}
	return ok, nil
}

func (ok *OffchainKeyring) OffchainSign(msg []byte) (signature []byte, err error) {
	return ed25519.Sign(ed25519.PrivateKey(ok.signingKey), msg), nil
}

func (ok *OffchainKeyring) ConfigDiffieHellman(point [curve25519.PointSize]byte) ([curve25519.PointSize]byte, error) {
	p, err := curve25519.X25519(ok.encryptionKey[:], point[:])
	if err != nil {
		return [curve25519.PointSize]byte{}, err
	}
	sharedPoint := [ed25519.PublicKeySize]byte{}
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

func (ok *OffchainKeyring) OffchainPublicKey() ocrtypes.OffchainPublicKey {
	return ocrtypes.OffchainPublicKey(ed25519.PrivateKey(ok.signingKey).Public().(ed25519.PublicKey))
}

func (ok *OffchainKeyring) ConfigEncryptionPublicKey() ocrtypes.ConfigEncryptionPublicKey {
	rv, err := curve25519.X25519(ok.encryptionKey[:], curve25519.Basepoint)
	if err != nil {
		logger.Errorw("Failure computing public key", "error", err)
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
}

func (ok *OffchainKeyring) marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, ok.signingKey)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, ok.encryptionKey)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (ok *OffchainKeyring) unmarshal(in []byte) error {
	buffer := bytes.NewReader(in)
	ok.signingKey = make(ed25519.PrivateKey, ed25519.PrivateKeySize)
	err := binary.Read(buffer, binary.LittleEndian, &ok.signingKey)
	if err != nil {
		return err
	}
	ok.encryptionKey = [curve25519.ScalarSize]byte{}
	err = binary.Read(buffer, binary.LittleEndian, &ok.encryptionKey)
	if err != nil {
		return err
	}
	return nil
}
