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
	return nl, l.UnwrapTo(&nl)
}

func (l *List) UnwrapTo(to any) error {
	val := reflect.ValueOf(to)
	if val.Kind() != reflect.Pointer {
		return fmt.Errorf("cannot unwrap to non-pointer type %T", to)
	}

	if val.IsNil() {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	ptrVal := reflect.Indirect(val)
	switch ptrVal.Kind() {
	case reflect.Slice:
		newList := reflect.MakeSlice(ptrVal.Type(), len(l.Underlying), len(l.Underlying))
		for i, el := range l.Underlying {
			newElm := newList.Index(i)
			if newElm.Kind() == reflect.Pointer {
				newElm.Set(reflect.New(newElm.Type().Elem()))
			} else {
				newElm = newElm.Addr()
			}

			if el == nil {
				continue
			}
			if err := el.UnwrapTo(newElm.Interface()); err != nil {
				return err
			}
		}
		reflect.Indirect(val).Set(newList)
		return nil
	default:
		dl := []any{}
		err := l.UnwrapTo(&dl)
		if err != nil {
			return err
		}

		if reflect.TypeOf(dl).AssignableTo(ptrVal.Type()) {
			ptrVal.Set(reflect.ValueOf(dl))
			return nil
		}

		return fmt.Errorf("cannot unwrap to type %T", to)
	}
}
