// Copyright (C) 2019 ScyllaDB
// Use of this source code is governed by a ALv2-style
// license that can be found in the LICENSE file.

package reflectx

import (
	"fmt"
	"unicode"
)

// CamelToSnakeASCII converts camel case strings to snake case. For performance
// reasons it only works with ASCII strings.
func CamelToSnakeASCII(s string) string {
	buf := []byte(s)
	out := make([]byte, 0, len(buf)+3)

	l := len(buf)
	for i := 0; i < l; i++ {
		if !(allowedBindRune(buf[i]) || buf[i] == '_') {
			panic(fmt.Sprint("not allowed name ", s))
		}

		b := rune(buf[i])

		if unicode.IsUpper(b) {
			if i > 0 && buf[i-1] != '_' && (unicode.IsLower(rune(buf[i-1])) || (i+1 < l && unicode.IsLower(rune(buf[i+1])))) {
				out = append(out, '_')
			}
			b = unicode.ToLower(b)
		}

		out = append(out, byte(b))
	}

	return string(out)
}

func allowedBindRune(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
