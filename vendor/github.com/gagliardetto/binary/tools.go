// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bin

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatByteSlice formats the given byte slice into a readable format.
func FormatByteSlice(buf []byte) string {
	elems := make([]string, 0)
	for _, v := range buf {
		elems = append(elems, strconv.Itoa(int(v)))
	}

	return "{" + strings.Join(elems, ", ") + "}" + fmt.Sprintf("(len=%v)", len(elems))
}

func FormatDiscriminator(disc [8]byte) string {
	elems := make([]string, 0)
	for _, v := range disc {
		elems = append(elems, strconv.Itoa(int(v)))
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

type WriteByWrite struct {
	writes [][]byte
	name   string
}

func NewWriteByWrite(name string) *WriteByWrite {
	return &WriteByWrite{
		name: name,
	}
}

func (rec *WriteByWrite) Write(b []byte) (int, error) {
	rec.writes = append(rec.writes, b)
	return len(b), nil
}

func (rec *WriteByWrite) Bytes() []byte {
	out := make([]byte, 0)
	for _, v := range rec.writes {
		out = append(out, v...)
	}
	return out
}

func (rec WriteByWrite) String() string {
	builder := new(strings.Builder)
	if rec.name != "" {
		builder.WriteString(rec.name + ":\n")
	}
	for index, v := range rec.writes {
		builder.WriteString(fmt.Sprintf("- %v: %s\n", index, FormatByteSlice(v)))
	}
	return builder.String()
}

// IsByteSlice returns true if the provided element is a []byte.
func IsByteSlice(v interface{}) bool {
	_, ok := v.([]byte)
	return ok
}
