package point_translation

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
)

type PairingTranslation struct{ pairing.Suite }

func (t *PairingTranslation) TranslateKey(s kyber.Scalar) (kyber.Point, error) {
	if reflect.TypeOf(s) != reflect.TypeOf(t.G1().Scalar()) {
		return nil, errors.Errorf("need scalar of type %T, got %T", t.G1().Scalar(), s)
	}
	return t.G2().Point().Mul(s, nil), nil
}

func (t *PairingTranslation) VerifyTranslation(pk1, pk2 kyber.Point) error {
	if reflect.TypeOf(pk1) != reflect.TypeOf(t.G1().Point()) {
		return fmt.Errorf("point for translation must be on G1, not %T", pk1)
	}
	if reflect.TypeOf(pk2) != reflect.TypeOf(t.G2().Point()) {
		return fmt.Errorf("translation must be on G2, not %T", pk2)
	}
	g1, g2 := t.G1().Point().Base(), t.G2().Point().Base()

	if !t.Pair(pk1, g2).Equal(t.Pair(g1, pk2)) {
		return errors.Errorf("putative G₂ public key has wrong discrete log")
	}
	return nil
}

func (t *PairingTranslation) Name() string {
	return fmt.Sprintf("translator from %s to %s", t.G1().String(), t.G2().String())
}

func (t *PairingTranslation) TargetGroup(
	sourceGroup kyber.Group,
) (targetGroup kyber.Group, err error) {
	if sourceGroup.String() != t.G1().String() {
		return nil, errors.Errorf("attempt to get target group from wrong source group")
	}
	return t.G2(), nil
}
