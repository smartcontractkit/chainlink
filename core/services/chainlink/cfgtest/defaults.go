package cfgtest

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

// DocDefaultsOnly reads only the default values from a docs TOML file and decodes in to cfg.
// Fields without defaults will set to zero values.
func DocDefaultsOnly(r io.Reader, cfg any) error {
	pr, pw := io.Pipe()
	defer pr.Close()
	go writeDefaults(r, pw)
	if err := decodeDefaults(pr, cfg); err != nil {
		return err
	}
	// replace niled examples with zero values.
	nilToZero(reflect.ValueOf(cfg))
	return nil
}

// decodeDefaults decodes to defaults from r.
func decodeDefaults(r io.Reader, cfg any) error {
	d := toml.NewDecoder(r).DisallowUnknownFields()
	if err := d.Decode(cfg); err != nil {
		return errors.Wrapf(err, "failed to decode default core configuration")
	}
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
	k := val.Kind()
	switch k { //nolint:exhaustive
	case reflect.Map, reflect.Slice:
		return
	case reflect.Ptr:
		if val.IsNil() {
			t := val.Type().Elem()
			val.Set(reflect.New(t))
		}
	}
	if k == reflect.Ptr {
		if val.Type().Implements(textUnmarshalerType) {
			return // skip values unmarshaled from strings
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
		return
	case reflect.Slice:
		return
	default:
		return
	}
}
