package whitespace

import (
	"strings"
	"unicode"
)

// findPythonStringSpans locates byte offset ranges of multiline triple-quoted strings
// (and single-quoted strings) in Python source code.
func findPythonStringSpans(content []byte) []byteSpan {
	var (
		spans []byteSpan
		i     = 0
		n     = len(content)
	)

	for i < n {
		b := content[i]

		// Comment: skip to end of line
		if b == '#' {
			for i < n && content[i] != '\n' {
				i++
			}
			continue
		}

		// Check for string prefix (r, u, f, b, fr, rf, br, rb, etc.)
		prefixLen := 0
		if isPythonStringPrefixChar(b) {
			pEnd := i
			for pEnd < n && isPythonStringPrefixChar(content[pEnd]) {
				pEnd++
			}
			prefixCandidate := strings.ToLower(string(content[i:pEnd]))
			if isValidPythonPrefix(prefixCandidate) && pEnd < n && (content[pEnd] == '\'' || content[pEnd] == '"') {
				prefixLen = pEnd - i
			}
		}

		quoteStart := i + prefixLen
		if quoteStart < n && (content[quoteStart] == '\'' || content[quoteStart] == '"') {
			quoteChar := content[quoteStart]
			// Check for triple-quote
			if quoteStart+2 < n && content[quoteStart+1] == quoteChar && content[quoteStart+2] == quoteChar {
				// Triple-quoted string
				startOffset := i
				i = quoteStart + 3
				for i < n {
					if content[i] == '\\' {
						i += 2 // skip escaped character
						continue
					}
					if i+2 < n && content[i] == quoteChar && content[i+1] == quoteChar && content[i+2] == quoteChar {
						i += 3
						break
					}
					i++
				}
				spans = append(spans, byteSpan{start: startOffset, end: i})
				continue
			}

			// Single-quoted string
			startOffset := i
			i = quoteStart + 1
			for i < n && content[i] != '\n' {
				if content[i] == '\\' {
					i += 2
					continue
				}
				if content[i] == quoteChar {
					i++
					break
				}
				i++
			}
			spans = append(spans, byteSpan{start: startOffset, end: i})
			continue
		}

		i++
	}

	return spans
}

func isPythonStringPrefixChar(b byte) bool {
	c := unicode.ToLower(rune(b))
	return c == 'r' || c == 'u' || c == 'f' || c == 'b'
}

func isValidPythonPrefix(s string) bool {
	switch s {
	case "r", "u", "f", "b", "rf", "fr", "rb", "br":
		return true
	default:
		return false
	}
}

// FixPython trims trailing whitespace from Python files while strictly preserving
// all whitespace inside triple-quoted multiline strings.
func FixPython(content []byte) ([]byte, bool, error) {
	if len(content) == 0 {
		return content, false, nil
	}

	spans := findPythonStringSpans(content)
	return processLinesWithSpans(content, spans)
}
