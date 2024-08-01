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

package text

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"strings"
	"sync"
)

var DisableColors = false

func S(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func Sf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func Ln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

// Lnsf is alias of fmt.Sprintln(fmt.Sprintf())
func Lnsf(format string, a ...interface{}) string {
	return Ln(Sf(format, a...))
}

// LnsfI is alias of fmt.Sprintln(fmt.Sprintf())
func LnsfI(indent int, format string, a ...interface{}) string {
	return Ln(Sf(strings.Repeat("	", indent)+format, a...))
}

// CC concats strings
func CC(elems ...string) string {
	return strings.Join(elems, "")
}

func Black(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 0, 0, 0)
}

func White(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 255, 255, 255)
}

func BlackBG(str string) string {
	if DisableColors {
		return str
	}
	return BgString(str, 0, 0, 0)
}

func WhiteBG(str string) string {
	if DisableColors {
		return str
	}
	return Black(BgString(str, 255, 255, 255))
}

func Lime(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 252, 255, 43)
}

func LimeBG(str string) string {
	if DisableColors {
		return str
	}
	return Black(BgString(str, 252, 255, 43))
}

func Yellow(str string) string {
	if DisableColors {
		return str
	}
	return BlackBG(FgString(str, 255, 255, 0))
}

func YellowBG(str string) string {
	return Black(BgString(str, 255, 255, 0))
}

func Orange(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 255, 165, 0)
}

func OrangeBG(str string) string {
	if DisableColors {
		return str
	}
	return Black(BgString(str, 255, 165, 0))
}

func Red(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 255, 0, 0)
}

func RedBG(str string) string {
	if DisableColors {
		return str
	}
	return White(BgString(str, 220, 20, 60))
}

// light blue?
func Shakespeare(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 82, 179, 217)
}

func ShakespeareBG(str string) string {
	if DisableColors {
		return str
	}
	return White(BgString(str, 82, 179, 217))
}

func Purple(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 255, 0, 255)
}

func PurpleBG(str string) string {
	if DisableColors {
		return str
	}
	return Black(BgString(str, 255, 0, 255))
}

func Indigo(str string) string {
	if DisableColors {
		return str
	}
	return FgString(str, 75, 0, 130)
}

func IndigoBG(str string) string {
	if DisableColors {
		return str
	}
	return BgString(str, 75, 0, 130)
}

func Bold(str string) string {
	if DisableColors {
		return str
	}
	return foreachLine(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[1m%s\033[0m", line)
	})
}

type sf func(int, string) string

func foreachLine(str string, transform sf) (out string) {
	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}
	return
}

func HighlightRedBG(str, substr string) string {
	return HighlightAnyCase(str, substr, RedBG)
}

func HighlightLimeBG(str, substr string) string {
	return HighlightAnyCase(str, substr, LimeBG)
}

func HighlightAnyCase(str, substr string, colorer func(string) string) string {
	substr = strings.ToLower(substr)
	str = strings.ToLower(str)

	hiSubstr := colorer(substr)
	return strings.Replace(str, substr, hiSubstr, -1)
}

func StringToColor(str string) func(string) string {
	hs := HashString(str)
	r, g, b, _ := calcColor(hs)

	bgColor := WhiteBG
	if IsLight(r, g, b) {
		bgColor = BlackBG
	}
	return func(str string) string {
		return bgColor(FgString(str, uint8(r), uint8(g), uint8(b)))
	}
}

func StringToColorBG(str string) func(string) string {
	hs := HashString(str)
	r, g, b, _ := calcColor(hs)

	textColor := White
	if IsLight(r, g, b) {
		textColor = Black
	}
	return func(str string) string {
		return textColor(BgString(str, uint8(r), uint8(g), uint8(b)))
	}
}

func Colorize(str string) string {
	if DisableColors {
		return str
	}
	colorizer := StringToColor(str)
	return colorizer(str)
}

func ColorizeBG(str string) string {
	if DisableColors {
		return str
	}
	colorizer := StringToColorBG(str)
	return colorizer(str)
}

func calcColor(color uint64) (red, green, blue, alpha uint64) {
	alpha = color & 0xFF
	blue = (color >> 8) & 0xFF
	green = (color >> 16) & 0xFF
	red = (color >> 24) & 0xFF

	return red, green, blue, alpha
}

// IsLight returns whether the color is perceived to be a light color
func IsLight(rr, gg, bb uint64) bool {
	r := float64(rr)
	g := float64(gg)
	b := float64(bb)

	hsp := math.Sqrt(0.299*math.Pow(r, 2) + 0.587*math.Pow(g, 2) + 0.114*math.Pow(b, 2))

	return hsp > 130
}

var hasherPool *sync.Pool

func init() {
	hasherPool = &sync.Pool{
		New: func() interface{} {
			return fnv.New64a()
		},
	}
}

func HashString(s string) uint64 {
	h := hasherPool.Get().(hash.Hash64)
	defer hasherPool.Put(h)
	h.Reset()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	return h.Sum64()
}
