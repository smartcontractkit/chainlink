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

package text

import (
	"reflect"
	"strings"
)

type fieldTag struct {
	Linear     bool
	Skip       bool
	Label      string
	NoTypeName bool
}

func parseFieldTag(tag reflect.StructTag) *fieldTag {
	t := &fieldTag{}
	tagStr := tag.Get("text")
	if tagStr == "" {
		return t
	}
	for _, s := range strings.Split(tagStr, ",") {
		if strings.HasPrefix(s, "linear") {
			t.Linear = true
		} else if strings.HasPrefix(s, "notype") {
			t.NoTypeName = true
		} else if s == "-" {
			t.Skip = true
		} else {
			t.Label = s
		}
	}
	return t
}
