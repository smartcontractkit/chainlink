package values

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type List struct {
	Underlying []Value
}

func NewList(l []any) (*List, error) {
	lv := []Value{}
	for _, v := range l {
		ev, err := Wrap(v)
		if err != nil {
			return nil, err
		}

		lv = append(lv, ev)
	}
	return &List{Underlying: lv}, nil
}

func (l *List) Proto() (*pb.Value, error) {
	v := []*pb.Value{}
	for _, e := range l.Underlying {
		pe, err := e.Proto()
		if err != nil {
			return nil, err
		}

		v = append(v, pe)
	}
	return pb.NewListValue(v)
}

func (l *List) Unwrap() (any, error) {
	nl := []any{}
	for _, v := range l.Underlying {
		uv, err := v.Unwrap()
		if err != nil {
			return nil, err
		}

		nl = append(nl, uv)
	}

	return nl, nil
}
