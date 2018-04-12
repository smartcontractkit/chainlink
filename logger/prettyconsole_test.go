package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyConsole_Write(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{
			"headline",
			`{"ts":1523537728.7260377, "level":"info", "msg":"top level"}`,
			"2018-04-12T12:55:28Z [INFO]  top level  \n",
			false,
		},
		{
			"details",
			`{"ts":1523537728, "level":"debug", "msg":"top level", "details":"nuances"}`,
			"2018-04-12T12:55:28Z [DEBUG] top level  \ndetails=nuances \n",
			false,
		},
		{
			"blacklist",
			`{"ts":1523537728, "level":"warn", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z [WARN]  top level  \n",
			false,
		},
		{"error", `{"broken":}`, `{}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &testReader{}
			pc := PrettyConsole{tr}
			_, err := pc.Write([]byte(tt.input))

			if tt.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, tr.Written)
			}
		})
	}
}

type testReader struct {
	Written string
}

func (*testReader) Sync() error  { return nil }
func (*testReader) Close() error { return nil }

func (tr *testReader) Write(b []byte) (int, error) {
	tr.Written = string(b)
	return 0, nil
}
