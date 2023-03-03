package dkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"testing"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"

	"go.dedis.ch/kyber/v3"

	"golang.org/x/crypto/pbkdf2"
)

type SecretShare struct {
	Idx   player_idx.PlayerIdx
	share kyber.Scalar
}

func (s *SecretShare) Mul(p kyber.Point) kyber.Point {
	return p.Clone().Mul(s.share, p)
}

var defaultPBKDF2NumberOfIterations uint32 = 150_000

var encryptionVersionNum uint8

const saltLen, nonceLen = 32, 12

const preambleLength = 1 + saltLen + 4 + nonceLen

func (s *SecretShare) Encrypt(
	passphrase []byte, itersTestingOnly ...uint32,
) ([]byte, error) {
	if len(itersTestingOnly) > 1 {
		return nil, errors.Errorf("at most one derived-key iteration value allowed")
	}
	var salt [saltLen]byte
	if _, err := io.ReadFull(rand.Reader, salt[:]); err != nil {
		return nil, errors.Wrap(err, "could not sample salt value for derived key "+
			"in encryption of secret share")
	}
	iter := defaultPBKDF2NumberOfIterations
	if len(itersTestingOnly) > 0 {
		iter = itersTestingOnly[0]
	}
	var iterBin [4]byte
	binary.BigEndian.PutUint32(iterBin[:], iter)

	var nonce [nonceLen]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, errors.Wrap(err, "could not sample nonce for encryption of "+
			"secret share")
	}
	plaintext, err := s.share.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not serialize secret share for "+
			"encryption")
	}
	preamble := bytes.Join(
		[][]byte{{encryptionVersionNum}, salt[:], iterBin[:], nonce[:]}, nil,
	)
	if len(preamble) != preambleLength {
		panic("wrong length for preamble")
	}
	gcm, err := newGCM(passphrase, salt, iter, "en")
	if err != nil {
		return nil, err
	}
	return gcm.Seal(preamble, nonce[:], plaintext, nil), nil
}

func (s *SecretShare) Decrypt(
	passphrase, ciphertext []byte, shareGroup kyber.Group,
) error {
	cursor := 0
	versionNum := ciphertext[0]
	cursor++
	if versionNum != encryptionVersionNum {
		return errors.Errorf(
			"don't know how to decrypt version %d ciphertexts", versionNum,
		)
	}
	var salt [saltLen]byte
	if n := copy(salt[:], ciphertext[cursor:cursor+saltLen]); n != saltLen {
		return errors.Errorf("failed to read entire salt")
	}
	cursor += saltLen
	var iterBin [4]byte
	if n := copy(iterBin[:], ciphertext[cursor:cursor+4]); n != 4 {
		return errors.Errorf("failed to read entire number of pbkdf2 iterations")
	}
	cursor += 4
	iter := binary.BigEndian.Uint32(iterBin[:])
	var nonce [nonceLen]byte
	if n := copy(nonce[:], ciphertext[cursor:cursor+nonceLen]); n != nonceLen {
		return errors.Errorf("failed to read entire nonce")
	}
	cursor += nonceLen
	if cursor != preambleLength {
		panic("wrong preamble length")
	}
	gcm, err := newGCM(passphrase, salt, iter, "de")
	if err != nil {
		return err
	}
	plaintext, err := gcm.Open(nil, nonce[:], ciphertext[preambleLength:], nil)
	if err != nil {
		return errors.Wrap(err, "could not decrypt secret share")
	}
	s.share = shareGroup.Scalar()
	return errors.Wrap(
		s.share.UnmarshalBinary(plaintext),
		"could not deserialize decrypted secret share",
	)
}

func (s SecretShare) Equal(os SecretShare) bool {
	return s.Idx.Equal(&os.Idx) && s.share.Equal(os.share)
}

func (s SecretShare) Clone() *SecretShare {
	return &SecretShare{s.Idx, s.share.Clone()}
}

func newGCM(
	passphrase []byte, salt [saltLen]byte, iter uint32, dir string,
) (cipher.AEAD, error) {
	errSuff := dir + "cryption of secret share"
	dk := pbkdf2.Key(passphrase, salt[:], int(iter), 32, sha256.New)
	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct block cipher for "+errSuff)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct GCM cipher for "+errSuff)
	}
	return gcm, nil
}

func XXXNewSecretShareTestingOnly(
	t *testing.T, i player_idx.PlayerIdx, s kyber.Scalar,
) *SecretShare {
	return &SecretShare{i, s}
}
