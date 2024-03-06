package values

import (
	"fmt"
	"reflect"

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

func (l *List) proto() *pb.Value {
	v := []*pb.Value{}
	for _, e := range l.Underlying {
		v = append(v, Proto(e))
	}
	return pb.NewListValue(v)
}

func (l *List) Unwrap() (any, error) {
	nl := []any{}
	for _, v := range l.Underlying {
		uv, err := Unwrap(v)
		if err != nil {
			return nil, err
		}

		nl = append(nl, uv)
	}

	return nl, nil
}

func (l *List) UnwrapTo(to any) error {
	val := reflect.ValueOf(to)
	if val.Kind() != reflect.Ptr && reflect.Indirect(val).Kind() != reflect.Slice {
		return fmt.Errorf("cannot unwrap to type %T", to)
	}

	if val.IsNil() {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	newList := reflect.New(reflect.Indirect(val).Type()).Elem()
	for _, el := range l.Underlying {
		newEl := reflect.New(reflect.Indirect(val).Type().Elem()).Elem()
		ptrEl := newEl.Addr().Interface()
		err := el.UnwrapTo(ptrEl)
		if err != nil {
			return err
		}
		newList = reflect.Append(newList, reflect.Indirect(reflect.ValueOf(ptrEl)))
	}

	reflect.Indirect(val).Set(newList)
	return nil
}
