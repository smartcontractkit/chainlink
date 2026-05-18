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
