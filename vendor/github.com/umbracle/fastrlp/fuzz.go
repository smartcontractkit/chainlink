package fastrlp

import (
	"bytes"
	"reflect"

	fuzz "github.com/google/gofuzz"
)

type FuzzObject interface {
	Marshaler
	Unmarshaler
}

type FuzzError struct {
	Source, Target interface{}
}

func (f *FuzzError) Error() string {
	return "failed to encode fuzz object"
}

type FuzzOption func(f *Fuzzer)

func WithPostHook(fnc func(FuzzObject) error) FuzzOption {
	return func(f *Fuzzer) {
		f.postHook = fnc
	}
}

func WithDefaults(fnc func(FuzzObject)) FuzzOption {
	return func(f *Fuzzer) {
		f.defaults = append(f.defaults, fnc)
	}
}

func copyObj(obj interface{}) interface{} {
	return reflect.New(reflect.TypeOf(obj).Elem()).Interface()
}

type Fuzzer struct {
	*fuzz.Fuzzer
	defaults []func(FuzzObject)
	postHook func(FuzzObject) error
}

func (f *Fuzzer) applyDefaults(obj FuzzObject) FuzzObject {
	for _, fn := range f.defaults {
		fn(obj)
	}
	return obj
}

func Fuzz(num int, base FuzzObject, opts ...FuzzOption) error {
	f := &Fuzzer{
		Fuzzer:   fuzz.New(),
		defaults: []func(FuzzObject){},
	}
	for _, opt := range opts {
		opt(f)
	}

	fuzzImpl := func() error {
		// marshal object with the fuzzing
		obj := copyObj(base).(FuzzObject)
		f.Fuzz(obj)
		f.applyDefaults(obj)

		// unmarshal object
		obj2 := f.applyDefaults(copyObj(obj).(FuzzObject))

		data, err := obj.MarshalRLPTo(nil)
		if err != nil {
			return err
		}
		if err := obj2.UnmarshalRLP(data); err != nil {
			return err
		}

		// instead of relying on DeepEqual and issues with zero arrays and so on
		// we use the rlp marshal values to compare
		data2, err := obj2.MarshalRLPTo(nil)
		if err != nil {
			return err
		}
		if !bytes.Equal(data, data2) {
			return &FuzzError{Source: obj, Target: obj2}
		}
		if f.postHook != nil {
			if err := f.postHook(obj2); err != nil {
				return err
			}
		}
		return nil
	}

	for i := 0; i < num; i++ {
		if err := fuzzImpl(); err != nil {
			return err
		}
	}
	return nil
}
