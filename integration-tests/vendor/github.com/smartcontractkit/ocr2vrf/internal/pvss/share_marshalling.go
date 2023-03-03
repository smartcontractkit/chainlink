package pvss

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/ciphertext"
)

func (s *share) marshal() ([]byte, error) {
	rv := make([][]byte, 3)
	cursor := 0

	cm, err := s.cipherText.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal ciphertext of share")
	}
	if len(cm) > math.MaxUint16 {
		return nil, errors.Errorf("marshaled ciphertext too long")
	}
	ctLen := make([]byte, 2)
	binary.BigEndian.PutUint16(ctLen, uint16(len(cm)))
	rv[cursor] = append(ctLen, cm...)
	cursor++

	if rv[cursor], err = marshalKyberPointWithLen(s.encryptionKey); err != nil {
		return nil, errors.Wrap(err, "could not marshal encryptionKey")
	}
	cursor++

	if rv[cursor], err = marshalKyberPointWithLen(s.subKeyTranslation); err != nil {
		return nil, errors.Wrap(err, "could not marshal subKeyTranslation")
	}
	cursor++

	if cursor != len(rv) {
		panic(errors.Errorf("marshal fields out of alignment"))
	}

	return bytes.Join(rv, nil), nil
}

func marshalKyberPointWithLen(p kyber.Point) ([]byte, error) {
	pm, err := p.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal key of share")
	}
	if len(pm) > math.MaxUint8 {
		return nil, errors.Errorf("marshalled point too long")
	}
	return bytes.Join([][]byte{{uint8(len(pm))}, pm}, nil), nil
}

func unmarshal(
	group anon.Suite, translationGroup kyber.Group, data []byte, ss *ShareSet,
) (*share, []byte, error) {

	cipherTextB, data, err := readLenPrefixedBytes(data, 2)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal ciphertext of share")
	}
	cipherText, err := ciphertext.Unmarshal(group, bytes.NewBuffer(cipherTextB))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal ciphertext of share")
	}

	encryptionKeyB, data, err := readLenPrefixedBytes(data, 1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal encryptionKey of share")
	}
	encryptionKey := group.Point()
	if err2 := encryptionKey.UnmarshalBinary(encryptionKeyB); err2 != nil {
		return nil, nil, errors.Wrap(err2, "could not unmarshal encryptionKey of share")
	}

	subKeyTranslationB, data, err := readLenPrefixedBytes(data, 1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal subKeyTranslation of share")
	}
	subKeyTranslation := translationGroup.Point()
	if err2 := subKeyTranslation.UnmarshalBinary(subKeyTranslationB); err2 != nil {
		return nil, nil, errors.Wrap(err2, "could not unmarshal subKeyTranslation of share")
	}

	return &share{cipherText, encryptionKey, subKeyTranslation, ss}, data, nil
}

func readLenPrefixedBytes(data []byte, prefixLen uint8) (read, remains []byte, err error) {
	if len(data) < int(prefixLen) {
		return nil, nil, errors.Errorf("marshalled data too short for length prefix")
	}
	length := 0
	for i := uint8(0); i < prefixLen; i++ {
		length <<= 8
		length += int(data[i])
	}
	readEnd := int(prefixLen) + length
	if len(data) < readEnd {
		return nil, nil, errors.Errorf("read length longer than marshal data")
	}
	read = data[prefixLen:readEnd]
	return read, data[readEnd:], nil
}
