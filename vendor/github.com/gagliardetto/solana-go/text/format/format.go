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

package format

import (
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/text"
	. "github.com/gagliardetto/solana-go/text"
)

func Program(name string, programID solana.PublicKey) string {
	return IndigoBG("Program") + ": " + Bold(name) + " " + text.ColorizeBG(programID.String())
}

func Instruction(name string) string {
	return Purple(Bold("Instruction")) + ": " + Bold(name)
}

func Param(name string, value interface{}) string {
	return Sf(
		Shakespeare(name)+": %s",
		strings.TrimSpace(
			prefixEachLineExceptFirst(
				strings.Repeat(" ", len(name)+2),
				strings.TrimSpace(spew.Sdump(value)),
			),
		),
	)
}

func Account(name string, pubKey solana.PublicKey) string {
	return Shakespeare(name) + ": " + text.ColorizeBG(pubKey.String())
}

func Meta(name string, meta *solana.AccountMeta) string {
	if meta == nil {
		return Shakespeare(name) + ": " + "<nil>"
	}
	out := Shakespeare(name) + ": " + text.ColorizeBG(meta.PublicKey.String())
	out += " ["
	if meta.IsWritable {
		out += "WRITE"
	}
	if meta.IsSigner {
		if meta.IsWritable {
			out += ", "
		}
		out += "SIGN"
	}
	out += "] "
	return out
}

func prefixEachLineExceptFirst(prefix string, s string) string {
	return foreachLine(s,
		func(i int, line string) string {
			if i == 0 {
				return Lime(line) + "\n"
			}
			return prefix + Lime(line) + "\n"
		})
}

type sf func(int, string) string

func foreachLine(str string, transform sf) (out string) {
	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}
	return
}
