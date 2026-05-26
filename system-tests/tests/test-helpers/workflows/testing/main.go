package testing

import (
	"fmt"
	"log/slog"
	"strings"
)

type T struct {
	*slog.Logger
}

func (t *T) Errorf(format string, args ...interface{}) {
	// if the log was produced by require/assert we need to split it, as engine does not allow logs longer than 1k bytes
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok && strings.Contains(msg, "Error:") && strings.Contains(msg, "Error Trace:") {
			var out []string
			for _, line := range strings.Split(msg, "\n") {
				if strings.Contains(line, "Error Trace") {
					continue
				}

				out = append(out, line)
			}

			t.Logger.Error(strings.Join(out, ";"))
			return
		}
	}
	t.Logger.Error(fmt.Sprintf(format, args...))
}

func (t *T) FailNow() {
	panic("Test failed. Panic to stop execution")
}
