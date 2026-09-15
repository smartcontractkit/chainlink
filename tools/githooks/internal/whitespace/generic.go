package whitespace

import (
	"bytes"
)

// FixGeneric trims trailing whitespace and whitespace-only lines for general text/code files.
func FixGeneric(content []byte) ([]byte, bool, error) {
	if len(content) == 0 {
		return content, false, nil
	}

	var (
		buf     bytes.Buffer
		changed bool
	)

	buf.Grow(len(content))

	lines := bytes.SplitAfter(content, []byte("\n"))
	for _, line := range lines {
		lineLen := len(line)
		if lineLen == 0 {
			continue
		}

		eolLen := 0
		if bytes.HasSuffix(line, []byte("\r\n")) {
			eolLen = 2
		} else if bytes.HasSuffix(line, []byte("\n")) {
			eolLen = 1
		}

		rawLine := line[:lineLen-eolLen]
		eol := line[lineLen-eolLen:]

		trimmed := bytes.TrimRight(rawLine, " \t")
		if len(trimmed) != len(rawLine) {
			changed = true
		}

		buf.Write(trimmed)
		buf.Write(eol)
	}

	if !changed {
		return content, false, nil
	}

	return buf.Bytes(), true, nil
}
