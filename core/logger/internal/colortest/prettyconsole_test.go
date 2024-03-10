package colortest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func init() {
	logger.InitColor(true)
}

func TestPrettyConsole_Write(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{
			"debug",
			`{"ts":1523537728, "level":"debug", "msg":"top level", "details":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[32m[DEBUG] \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \x1b[32mdetails\x1b[0m=nuances \n",
			false,
		},
		{
			"info",
			`{"ts":1523537728.7260377, "level":"info", "msg":"top level"}`,
			"2018-04-12T12:55:28Z \x1b[37m[INFO]  \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"warn",
			`{"ts":1523537728, "level":"warn", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[33m[WARN]  \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"error",
			`{"ts":1523537728, "level":"error", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[31m[ERROR] \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"critical",
			`{"ts":1523537728, "level":"crit", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[91m[CRIT]  \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"panic",
			`{"ts":1523537728, "level":"panic", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[91m[PANIC] \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"fatal",
			`{"ts":1523537728, "level":"fatal", "msg":"top level", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[91m[FATAL] \x1b[0mtop level                                          \x1b[34m\x1b[0m                        \n",
			false,
		},
		{
			"control",
			`{"ts":1523537728, "level":"fatal", "msg":"\u0008\t\n\r\u000b\u000c\ufffd\ufffd", "hash":"nuances"}`,
			"2018-04-12T12:55:28Z \x1b[91m[FATAL] \x1b[0m\\b\t\n\r\\v\\f��                                        \x1b[34m\x1b[0m                        \n",
			false,
		},
		{"broken", `{"broken":}`, `{}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &testReader{}
			pc := logger.PrettyConsole{Sink: tr}
			_, err := pc.Write([]byte(tt.input))

			if tt.wantError {
				assert.Error(t, err)
			} else {
				t.Log(tr.Written)
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
