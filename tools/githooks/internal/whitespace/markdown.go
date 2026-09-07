package whitespace

import (
	"bytes"
)

func isAllWhitespace(data []byte) bool {
	for _, b := range data {
		if b != ' ' && b != '\t' && b != '\r' && b != '\n' {
			return false
		}
	}
	return true
}

// FixMarkdown processes markdown files with parity to pre-commit's trailing-whitespace fixer,
// preserving intentional CommonMark two-space hard line breaks for non-blank lines.
func FixMarkdown(content []byte) ([]byte, bool, error) {
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

		switch {
		case isAllWhitespace(rawLine):
			// Blank line: strip completely
			if len(rawLine) > 0 {
				changed = true
			}
			buf.Write(eol)

		case bytes.HasSuffix(rawLine, []byte("  ")):
			// Preserve 2 trailing spaces for markdown hard line break
			trimmed := bytes.TrimRight(rawLine[:len(rawLine)-2], " \t")
			if len(trimmed)+2 != len(rawLine) {
				changed = true
			}
			buf.Write(trimmed)
			buf.WriteString("  ")
			buf.Write(eol)

		default:
			trimmed := bytes.TrimRight(rawLine, " \t")
			if len(trimmed) != len(rawLine) {
				changed = true
			}
			buf.Write(trimmed)
			buf.Write(eol)
		}
	}

	if !changed {
		return content, false, nil
	}

	return buf.Bytes(), true, nil
}
