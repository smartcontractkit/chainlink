package whitespace

import (
	"bytes"
	"go/scanner"
	"go/token"
)

type byteSpan struct {
	start int
	end   int
}

func (s byteSpan) contains(pos int) bool {
	return pos >= s.start && pos < s.end
}

// FixGo trims trailing whitespace from Go files while strictly preserving
// all whitespace inside string literals (especially raw backtick strings).
func FixGo(content []byte) ([]byte, bool, error) {
	if len(content) == 0 {
		return content, false, nil
	}

	fset := token.NewFileSet()
	f := fset.AddFile("", fset.Base(), len(content))

	var s scanner.Scanner
	s.Init(f, content, nil, scanner.ScanComments)

	var stringSpans []byteSpan

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.STRING {
			offset := f.Offset(pos)
			stringSpans = append(stringSpans, byteSpan{
				start: offset,
				end:   offset + len(lit),
			})
		}
	}

	return processLinesWithSpans(content, stringSpans)
}

func processLinesWithSpans(content []byte, protectedSpans []byteSpan) ([]byte, bool, error) {
	var (
		buf     bytes.Buffer
		offset  int
		changed bool
	)

	buf.Grow(len(content))

	lines := bytes.SplitAfter(content, []byte("\n"))
	for _, line := range lines {
		lineLen := len(line)
		if lineLen == 0 {
			continue
		}

		lineStart := offset
		offset += lineLen

		eolLen := 0
		if bytes.HasSuffix(line, []byte("\r\n")) {
			eolLen = 2
		} else if bytes.HasSuffix(line, []byte("\n")) {
			eolLen = 1
		}

		rawLine := line[:lineLen-eolLen]
		eol := line[lineLen-eolLen:]

		// Find trailing whitespace
		trimmedLine := bytes.TrimRight(rawLine, " \t")
		wsLen := len(rawLine) - len(trimmedLine)

		if wsLen == 0 {
			buf.Write(line)
			continue
		}

		// Check if trailing whitespace falls inside a protected span
		wsStartOffset := lineStart + len(trimmedLine)
		isProtected := false
		for _, span := range protectedSpans {
			if span.contains(wsStartOffset) {
				isProtected = true
				break
			}
		}

		if isProtected {
			buf.Write(line)
			continue
		}

		changed = true
		buf.Write(trimmedLine)
		buf.Write(eol)
	}

	if !changed {
		return content, false, nil
	}

	return buf.Bytes(), true, nil
}
