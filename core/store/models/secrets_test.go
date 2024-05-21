package models

import (
	"encoding"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
	type secret interface {
		fmt.Stringer
		encoding.TextMarshaler
	}
	for _, v := range []secret{
		Secret("secret"),
		MustSecretURL("http://secret.url"),
	} {
		t.Run(fmt.Sprintf("%T", v), func(t *testing.T) {
			assert.Equal(t, redacted, v.String())
			got, err := v.MarshalText()
			if assert.NoError(t, err) {
				assert.Equal(t, redacted, string(got))
			}
			assert.Equal(t, redacted, fmt.Sprint(v))
			assert.Equal(t, redacted, fmt.Sprintf("%s", v)) //nolint:gosimple
			assert.Equal(t, redacted, fmt.Sprintf("%v", v))
			assert.Equal(t, redacted, fmt.Sprintf("%#v", v))
			got, err = json.Marshal(v)
			if assert.NoError(t, err) {
				assert.Equal(t, fmt.Sprintf(`"%s"`, redacted), string(got))
			}
		})
	}
}
