package cfgtest

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// DocDefaultsOnly reads only the default values from a docs TOML file and decodes in to cfg.
// Fields without defaults will set to zero values.
func DocDefaultsOnly(r io.Reader, cfg any, decode func(io.Reader, any) error) error {
	pr, pw := io.Pipe()
	defer pr.Close()
	go writeDefaults(r, pw)
	if err := decode(pr, cfg); err != nil {
		return errors.Wrapf(err, "failed to decode default core configuration")
	}
	// replace niled examples with zero values.
	nilToZero(reflect.ValueOf(cfg))
	return nil
}

// writeDefaults writes default lines from defaultsTOML to w.
func writeDefaults(r io.Reader, w *io.PipeWriter) {
	defer w.Close()
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		// Skip comments and examples (which become zero values)
		if strings.HasPrefix(t, "#") || strings.HasSuffix(t, "# Example") {
			continue
		}
		if _, err := io.WriteString(w, t); err != nil {
			w.CloseWithError(err)
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			w.CloseWithError(err)
		}
	}
	if err := s.Err(); err != nil {
		w.CloseWithError(fmt.Errorf("failed to scan core defaults: %v", err))
	}
}

func nilToZero(val reflect.Value) {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			t := val.Type().Elem()
			val.Set(reflect.New(t))
		}
		if val.Type().Implements(textUnmarshalerType) {
			return // don't descend inside - leave whole zero value
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct:
		if val.Type().Implements(textUnmarshalerType) {
			return // skip values unmarshaled from strings
		}
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			nilToZero(f)
		}
		return
	case reflect.Map:
		if !val.IsNil() {
			for _, k := range val.MapKeys() {
				nilToZero(val.MapIndex(k))
			}
		}
		return
	case reflect.Slice:
		if !val.IsNil() {
			for i := 0; i < val.Len(); i++ {
				nilToZero(val.Index(i))
			}
		}
		return
	default:
		return
	}
}
