package ocr2key

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/box"

	"golang.org/x/crypto/curve25519"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.OffchainKeyring = &offchainKeyring{}

// offchainKeyring contains the secret keys needed for the OCR nodes to share secrets
// and perform aggregation.
//
// This is currently an ed25519 signing key and a separate encryption key.
//
// All its functions should be thread-safe.
type offchainKeyring struct {
	signingKey    func() ed25519.PrivateKey
	encryptionKey func() *[curve25519.ScalarSize]byte
}

func newOffchainKeyring(encryptionMaterial, signingMaterial io.Reader) (*offchainKeyring, error) {
	_, signingKey, err := ed25519.GenerateKey(signingMaterial)
	if err != nil {
		return nil, err
	}

	encryptionKey := [curve25519.ScalarSize]byte{}
	_, err = encryptionMaterial.Read(encryptionKey[:])
	if err != nil {
		return nil, err
	}

	ok := &offchainKeyring{
		signingKey:    func() ed25519.PrivateKey { return signingKey },
		encryptionKey: func() *[32]byte { return &encryptionKey },
	}
	_, err = ok.configEncryptionPublicKey()
	if err != nil {
		return nil, err
	}
	return ok, nil
}

// NaclBoxOpenAnonymous decrypts a message that was encrypted using the OCR2 Offchain public key
func (ok *offchainKeyring) NaclBoxOpenAnonymous(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) < box.Overhead {
		return nil, errors.New("ciphertext too short")
	}

	publicKey := [curve25519.PointSize]byte(ok.ConfigEncryptionPublicKey())

	decrypted, success := box.OpenAnonymous(nil, ciphertext, &publicKey, ok.encryptionKey())
	if !success {
		return nil, errors.New("decryption failed")
	}

	return decrypted, nil
}

// OffchainSign signs message using private key
func (ok *offchainKeyring) OffchainSign(msg []byte) (signature []byte, err error) {
	return ed25519.Sign(ok.signingKey(), msg), nil
}

// ConfigDiffieHellman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offchain key ring's encryption key.)
func (ok *offchainKeyring) ConfigDiffieHellman(point [curve25519.PointSize]byte) ([curve25519.PointSize]byte, error) {
	p, err := curve25519.X25519(ok.encryptionKey()[:], point[:])
	if err != nil {
		return [curve25519.PointSize]byte{}, err
	}
	sharedPoint := [ed25519.PublicKeySize]byte{}
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

// OffchainPublicKey returns the public component of this offchain keyring.
func (ok *offchainKeyring) OffchainPublicKey() ocrtypes.OffchainPublicKey {
	var offchainPubKey [ed25519.PublicKeySize]byte
	copy(offchainPubKey[:], ok.signingKey().Public().(ed25519.PublicKey)[:])
	return offchainPubKey
}

// ConfigEncryptionPublicKey returns config public key
func (ok *offchainKeyring) ConfigEncryptionPublicKey() ocrtypes.ConfigEncryptionPublicKey {
	cpk, _ := ok.configEncryptionPublicKey()
	return cpk
}

func (ok *offchainKeyring) configEncryptionPublicKey() (ocrtypes.ConfigEncryptionPublicKey, error) {
	rv, err := curve25519.X25519(ok.encryptionKey()[:], curve25519.Basepoint)
	if err != nil {
		return [curve25519.PointSize]byte{}, err
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed, nil
}

func (ok *offchainKeyring) marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, ok.signingKey())
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, ok.encryptionKey())
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (ok *offchainKeyring) unmarshal(in []byte) error {
	buffer := bytes.NewReader(in)
	signingKey := make(ed25519.PrivateKey, ed25519.PrivateKeySize)
	err := binary.Read(buffer, binary.LittleEndian, &signingKey)
	if err != nil {
		return err
	}
	ok.signingKey = func() ed25519.PrivateKey { return signingKey }
	encryptionKey := [curve25519.ScalarSize]byte{}
	err = binary.Read(buffer, binary.LittleEndian, &encryptionKey)
	if err != nil {
		return err
	}
	ok.encryptionKey = func() *[32]byte { return &encryptionKey }
	_, err = ok.configEncryptionPublicKey()
	return err
}
