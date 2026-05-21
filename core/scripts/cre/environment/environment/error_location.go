package environment

import (
	stderrors "errors"
	"fmt"
	"regexp"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

var (
	goStackFileLine     = regexp.MustCompile(`^\t(.*?\.go):(\d+)\s+`)
	pkgErrFramePathLine = regexp.MustCompile(`.+\.go:\d+$`)
)

// errorLocationForTracking returns a file:line suitable for DX telemetry.
// When panicStack is non-empty (e.g. from debug.Stack()), it is parsed first so the
// reported location reflects the panic site. Otherwise, or when parsing yields nothing,
// it falls back to the innermost github.com/pkg/errors stack attached on the unwrap chain.
func errorLocationForTracking(err error, panicStack []byte) string {
	if loc := innermostUserFrameFromGoStack(panicStack); loc != "" {
		return loc
	}
	return errorLocationFromWrappedErrors(err)
}

func innermostUserFrameFromGoStack(stack []byte) string {
	if len(stack) == 0 {
		return ""
	}
	for _, line := range strings.Split(string(stack), "\n") {
		m := goStackFileLine.FindStringSubmatch(line)
		if len(m) != 3 {
			continue
		}
		path := m[1]
		lineNo := m[2]
		if skipRuntimeStackPath(path) {
			continue
		}
		return path + ":" + lineNo
	}
	return ""
}

func skipRuntimeStackPath(path string) bool {
	return strings.Contains(path, "/runtime/") ||
		strings.Contains(path, `\runtime\`) ||
		strings.Contains(path, "/reflect/") ||
		strings.Contains(path, `\reflect\`)
}

type stackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

// errorLocationFromWrappedErrors returns the call site (file:line) from the innermost
// error in the unwrap chain that carries a non-empty pkg/errors stack trace.
func errorLocationFromWrappedErrors(err error) string {
	if err == nil {
		return ""
	}
	var innermost pkgerrors.StackTrace
	for e := err; e != nil; e = stderrors.Unwrap(e) {
		var st stackTracer
		if !stderrors.As(e, &st) {
			continue
		}
		frames := st.StackTrace()
		if len(frames) == 0 {
			continue
		}
		innermost = frames
	}
	if len(innermost) == 0 {
		return ""
	}
	return formatPkgErrorsFrameFileLine(innermost[0])
}

// formatPkgErrorsFrameFileLine returns absolute (or module) path and line only.
// fmt.Sprintf("%%+v", frame) also prints the qualified function name on a separate line.
func formatPkgErrorsFrameFileLine(fr pkgerrors.Frame) string {
	raw := fmt.Sprintf("%+v", fr)
	for _, ln := range strings.Split(raw, "\n") {
		t := strings.TrimSpace(ln)
		if pkgErrFramePathLine.MatchString(t) {
			return t
		}
	}
	return strings.TrimSpace(fmt.Sprintf("%v", fr))
}
