package ciphertext

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"

	"go.dedis.ch/kyber/v3/sign/anon"
)

const (
	SUITE_DESC_LEN                  = "suite description length"
	CRYPTOGRAPHIC_SUITE_DESCRIPTION = "cryptographic suite description"
	RECEIVER_INDEX                  = "receiver index"
	ENCRYPTION_KEY                  = "encryption key"
	SHARE_PROOF_LEN                 = "share proof length"
	SHARE_PROOF                     = "proof that cipher text encodes share"
	NUM_PAIRS                       = "number of ciphertext pairs"
)

func (c *cipherText) marshal() (m []byte, err error) {
	rv := new(bytes.Buffer)

	suiteName := []byte(c.suite.String())
	if len(suiteName) > math.MaxUint8 {
		return nil, errors.Errorf("suite name too long")
	}

	if err2 := hw(rv, []byte{uint8(len(suiteName))}, SUITE_DESC_LEN); err2 != nil {
		return nil, err2
	}

	if err2 := hw(rv, suiteName, CRYPTOGRAPHIC_SUITE_DESCRIPTION); err2 != nil {
		return nil, err2
	}

	if err2 := hw(rv, c.receiver.Marshal(), RECEIVER_INDEX); err2 != nil {
		return nil, err2
	}

	n, err := c.encryptionKey.MarshalTo(rv)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal "+ENCRYPTION_KEY)
	}
	if n != c.suite.PointLen() {
		return nil, errors.Errorf("could not marshal " + ENCRYPTION_KEY +
			"; wrote wrong number of bytes for it")
	}

	if len(c.encodesShareProof) > math.MaxUint16 {
		return nil, errors.Errorf("encodesShareProof too long")
	}
	proofLen := make([]byte, 2)
	binary.BigEndian.PutUint16(proofLen, uint16(len(c.encodesShareProof)))
	if err := hw(rv, proofLen, SHARE_PROOF_LEN); err != nil {
		return nil, err
	}

	if err := hw(rv, c.encodesShareProof, SHARE_PROOF); err != nil {
		return nil, err
	}

	numPairs := len(c.cipherText)
	if numPairs > math.MaxUint16 {
		return nil, errors.Errorf("too many pairs to marshal")
	}
	bigEndNp := make([]byte, 2)
	binary.BigEndian.PutUint16(bigEndNp, uint16(numPairs))
	if err := hw(rv, bigEndNp, NUM_PAIRS); err != nil {
		return nil, err
	}

	ctl := elGamalBitPairMarshalLength(c.suite)
	for _, ct := range c.cipherText {
		ctm, err := ct.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal ciphertext")
		}
		if len(ctm) != ctl {
			return nil, errors.Errorf("elGamalBitPair marshaled to wrong length")
		}
		if err := hw(rv, ctm, "ciphertext bit pair"); err != nil {
			return nil, err
		}
	}

	return rv.Bytes(), nil
}

func hw(rv io.Writer, d []byte, errmsg string) (err error) {
	errmsg = "could not marshal " + errmsg
	n, err := rv.Write(d)
	if err != nil {
		err = errors.Wrapf(err, errmsg)
	} else if n != len(d) {
		err = errors.Errorf(errmsg + ": failed to write all bytes")
	}
	return err
}

func unmarshal(suite anon.Suite, byteStream io.Reader) (c *cipherText, err error) {
	c = &cipherText{suite: suite}

	var strLen [1]byte
	if err2 := hr(byteStream, strLen[:], SUITE_DESC_LEN); err2 != nil {
		return nil, err2
	}
	str := make([]byte, strLen[0])
	if err2 := hr(byteStream, str, CRYPTOGRAPHIC_SUITE_DESCRIPTION); err2 != nil {
		return nil, err2
	}
	if string(str) != suite.String() {
		return nil, errors.Errorf(`wrong suite for unmarshalling: need "%s", got "%s"`, suite, str)
	}

	idx := make([]byte, player_idx.MarshalLen)
	if err2 := hr(byteStream, idx, RECEIVER_INDEX); err2 != nil {
		return nil, err2
	}
	c.receiver, _, err = player_idx.Unmarshal(idx)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ciphertext's "+RECEIVER_INDEX)
	}

	c.encryptionKey = suite.Point()
	n, err := c.encryptionKey.UnmarshalFrom(byteStream)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal "+ENCRYPTION_KEY)
	}
	if n != suite.PointLen() {
		return nil, errors.Errorf("could not unmarshal " + ENCRYPTION_KEY)
	}

	proofLen := make([]byte, 2)
	if err2 := hr(byteStream, proofLen[:], SHARE_PROOF_LEN); err2 != nil {
		return nil, err2
	}

	c.encodesShareProof = make([]byte, binary.BigEndian.Uint16(proofLen))
	if err2 := hr(byteStream, c.encodesShareProof, SHARE_PROOF); err2 != nil {
		return nil, err2
	}

	var rawNumPairs [2]byte
	if err2 := hr(byteStream, rawNumPairs[:], NUM_PAIRS); err2 != nil {
		return nil, err2
	}
	numPairs := binary.BigEndian.Uint16(rawNumPairs[:])

	ctm := make([]byte, elGamalBitPairMarshalLength(suite))
	c.cipherText = make([]*elGamalBitPair, numPairs)
	for bpIdx := uint16(0); bpIdx < numPairs; bpIdx++ {
		if err2 := hr(byteStream, ctm[:], "ciphertext bit pair"); err2 != nil {
			return nil, err2
		}
		c.cipherText[bpIdx], err = unmarshalElGamalBitPair(suite, ctm)
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal cipher text")
		}
	}

	return c, nil
}

func hr(byteStream io.Reader, dst []byte, errmsg string) (rv error) {
	errmsg = fmt.Sprint("could not read ", errmsg, " for unmarshalling")
	n, err := byteStream.Read(dst)
	if err != nil {
		return errors.Wrap(err, errmsg)
	}
	if n != len(dst) {
		return errors.Errorf(errmsg)
	}
	return nil
}
