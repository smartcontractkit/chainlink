package point_translation

import (
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
)

type TrivialTranslation struct{ base kyber.Point }

var _ PubKeyTranslation = (*TrivialTranslation)(nil)

func NewTrivialTranslation(base kyber.Point) *TrivialTranslation {
	return &TrivialTranslation{base}
}

func (t *TrivialTranslation) TranslateKey(share kyber.Scalar) (kyber.Point, error) {
	return t.base.Clone().Mul(share, nil), nil
}

func (t *TrivialTranslation) VerifyTranslation(pk1, pk2 kyber.Point) error {
	if pk1.Equal(pk2) {
		return nil
	}
	return errors.Errorf("putative translated points are not equal")
}

func (t *TrivialTranslation) Name() string {
	return "trivial translator"
}

func (t *TrivialTranslation) TargetGroup(
	sourceGroup kyber.Group,
) (targetGroup kyber.Group, err error) {
	return sourceGroup, nil
}
