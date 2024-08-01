// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
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
	"encoding/binary"
	"reflect"
	"strings"
)

type fieldTag struct {
	SizeOf          string
	Skip            bool
	Order           binary.ByteOrder
	Option          bool
	COption         bool
	BinaryExtension bool

	IsBorshEnum bool
}

func isIn(s string, candidates ...string) bool {
	for _, c := range candidates {
		if s == c {
			return true
		}
	}
	return false
}

func parseFieldTag(tag reflect.StructTag) *fieldTag {
	t := &fieldTag{
		Order: defaultByteOrder,
	}
	tagStr := tag.Get("bin")
	for _, s := range strings.Split(tagStr, " ") {
		if strings.HasPrefix(s, "sizeof=") {
			tmp := strings.SplitN(s, "=", 2)
			t.SizeOf = tmp[1]
		} else if s == "big" {
			t.Order = binary.BigEndian
		} else if s == "little" {
			t.Order = binary.LittleEndian
		} else if isIn(s, "optional", "option") {
			t.Option = true
		} else if isIn(s, "coption") {
			t.COption = true
		} else if s == "binary_extension" {
			t.BinaryExtension = true
		} else if isIn(s, "-", "skip") {
			t.Skip = true
		} else if isIn(s, "enum") {
			t.IsBorshEnum = true
		}
	}

	// TODO: parse other borsh tags
	if strings.TrimSpace(tag.Get("borsh_skip")) == "true" {
		t.Skip = true
	}
	if strings.TrimSpace(tag.Get("borsh_enum")) == "true" {
		t.IsBorshEnum = true
	}
	return t
}
