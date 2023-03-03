package schnorr

import (
	"bytes"
	"fmt"

	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/ocr2vrf/internal/util"
)

type Suite interface {
	kyber.Group
	kyber.Random
	kyber.XOFFactory
}

func Sign(s Suite, private kyber.Scalar, msg []byte) ([]byte, error) {
	var g kyber.Group = s

	k := g.Scalar().Pick(s.RandomStream())
	R := g.Point().Mul(k, nil)

	public := g.Point().Mul(private, nil)
	h, err := hash(s, public, R, msg)
	if err != nil {
		return nil, err
	}

	xh := g.Scalar().Mul(private, h)
	S := g.Scalar().Add(k, xh)

	var b bytes.Buffer
	if _, err := R.MarshalTo(&b); err != nil {
		return nil, err
	}
	if _, err := S.MarshalTo(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func VerifyWithChecks(g Suite, pub, msg, sig []byte) error {
	type scalarCanCheckCanonical interface {
		IsCanonical(b []byte) bool
	}

	type pointCanCheckCanonicalAndSmallOrder interface {
		HasSmallOrder() bool
		IsCanonical(b []byte) bool
	}

	R := g.Point()
	s := g.Scalar()
	pointSize := R.MarshalSize()
	scalarSize := s.MarshalSize()
	sigSize := scalarSize + pointSize
	if len(sig) != sigSize {
		return fmt.Errorf(
			"schnorr: signature of invalid length %d instead of %d",
			len(sig),
			sigSize,
		)
	}
	if err := R.UnmarshalBinary(sig[:pointSize]); err != nil {
		return util.WrapErrorf(err, "could not unmarshal R (0x%x)", sig[:pointSize])
	}
	if p, ok := R.(pointCanCheckCanonicalAndSmallOrder); ok {
		if !p.IsCanonical(sig[:pointSize]) {
			return fmt.Errorf("the point R is not canonical")
		}
		if p.HasSmallOrder() {
			return fmt.Errorf("the point R has small order")
		}
	}
	sc, ok := g.Scalar().(scalarCanCheckCanonical)
	if ok && !sc.IsCanonical(sig[pointSize:]) {
		return fmt.Errorf(
			"signature is not canonical: 0x%x %d",
			sig[pointSize:],
			len(sig[pointSize:]),
		)
	}
	if err := s.UnmarshalBinary(sig[pointSize:]); err != nil {
		return util.WrapError(err, "could not unmarshal S")
	}

	public := g.Point()
	err := public.UnmarshalBinary(pub)
	if err != nil {
		return util.WrapError(err, "schnorr: error unmarshalling public key")
	}
	if p, ok := public.(pointCanCheckCanonicalAndSmallOrder); ok {
		if !p.IsCanonical(pub) {
			return fmt.Errorf("public key is not canonical")
		}
		if p.HasSmallOrder() {
			return fmt.Errorf("public key has small order")
		}
	}

	h, err := hash(g, public, R, msg)
	if err != nil {
		return util.WrapError(err, "could not compute hash")
	}

	S := g.Point().Mul(s, nil)

	Ah := g.Point().Mul(h, public)
	RAs := g.Point().Add(R, Ah)

	if !S.Equal(RAs) {
		return fmt.Errorf("schnorr: invalid signature")
	}

	return nil

}

func Verify(g Suite, public kyber.Point, msg, sig []byte) error {
	PBuf, err := public.MarshalBinary()
	if err != nil {
		return util.WrapError(err, "could not marshal point for sig verification")
	}
	if err = g.Point().UnmarshalBinary(PBuf); err != nil {
		panic(err)
	}
	if err != nil {
		return util.WrapError(err, "error unmarshalling public key")
	}
	return VerifyWithChecks(g, PBuf, msg, sig)
}

func hash(g Suite, public, r kyber.Point, msg []byte) (kyber.Scalar, error) {
	h := g.XOF(nil)
	if _, err := r.MarshalTo(h); err != nil {
		return nil, err
	}
	if _, err := public.MarshalTo(h); err != nil {
		return nil, err
	}
	if _, err := h.Write(msg); err != nil {
		return nil, err
	}
	return g.Scalar().Pick(h), nil
}
