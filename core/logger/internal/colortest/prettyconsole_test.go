package colortest

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.SetColor(true)
}

func TestPrettyConsole_Write(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{
			"headline",
			`{"ts":1523537728.7260377, "level":"info", "msg":"top level"}`,
			"2018-04-12T12:55:28Z \x1b[37m[INFO]  \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"details",
			`{"ts":1523537728, "level":"debug", "msg":"top level", "details":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[32m[DEBUG] \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \x1b[32mdetails\x1b[0m=nuances \n",
			false,
		},
		{
			"blacklist",
			`{"ts":1523537728, "level":"warn", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[33m[WARN]  \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{"error", `{"broken":}`, `{}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &testReader{}
			pc := logger.PrettyConsole{tr}
			_, err := pc.Write([]byte(tt.input))

			if tt.wantError {
				assert.Error(t, err)
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
