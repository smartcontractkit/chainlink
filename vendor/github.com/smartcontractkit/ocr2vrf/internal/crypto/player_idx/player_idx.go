package player_idx

import (
	"reflect"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
)

type PlayerIdx struct{ idx Int }

type Int = uint8

var MaxPlayer = -Int(one)

func PlayerIdxs(n Int) (rv []*PlayerIdx, err error) {
	if n < 1 {
		return nil, errors.Errorf("%d is too few players", n)
	}
	for i := Int(1); i <= n; i++ {
		rv = append(rv, &PlayerIdx{i})
	}
	return rv, nil
}

func (pi PlayerIdx) Eval(f *share.PriPoly) kyber.Scalar {
	pi.mustBeNonZero()
	return f.Eval(int(pi.idx - 1)).V
}

func (pi PlayerIdx) EvalPoint(f *share.PubPoly) kyber.Point {
	pi.mustBeNonZero()
	return f.Eval(int(pi.idx - 1)).V
}

func (pi PlayerIdx) Index(a interface{}) interface{} {
	pi.mustBeNonZero()
	return reflect.ValueOf(a).Index(int(pi.idx) - 1).Interface()
}

func (pi PlayerIdx) Set(a, e interface{}) {
	pi.mustBeNonZero()
	reflect.ValueOf(a).Index(int(pi.idx) - 1).Set(reflect.ValueOf(e))
}

func (pi PlayerIdx) PriShare(sk kyber.Scalar) share.PriShare {
	pi.mustBeNonZero()
	return share.PriShare{int(pi.idx) - 1, sk}
}

func (pi PlayerIdx) PubShare(p kyber.Point) share.PubShare {
	pi.mustBeNonZero()
	return share.PubShare{int(pi.idx) - 1, p}
}
