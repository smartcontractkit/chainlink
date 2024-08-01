// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <https://unlicense.org>

package verkle

import (
	"encoding/hex"
	"sync"

	"github.com/crate-crypto/go-ipa/ipa"
)

// EmptyCodeHashPoint is a cached point that is used to represent an empty code hash.
// This value is initialized once in GetConfig().
var (
	EmptyCodeHashPoint           Point
	EmptyCodeHashFirstHalfValue  Fr
	EmptyCodeHashSecondHalfValue Fr
)

const (
	CodeHashVectorPosition     = 3 // Defined by the spec.
	EmptyCodeHashFirstHalfIdx  = CodeHashVectorPosition * 2
	EmptyCodeHashSecondHalfIdx = EmptyCodeHashFirstHalfIdx + 1
)

var (
	FrZero Fr
	FrOne  Fr

	cfg     *Config
	onceCfg sync.Once
)

func init() {
	FrZero.SetZero()
	FrOne.SetOne()
}

type IPAConfig struct {
	conf *ipa.IPAConfig
}

type Config = IPAConfig

func GetConfig() *Config {
	onceCfg.Do(func() {
		conf, err := ipa.NewIPASettings()
		if err != nil {
			panic(err)
		}
		cfg = &IPAConfig{conf: conf}

		// Initialize the empty code cached values.
		emptyHashCode, _ := hex.DecodeString("c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
		values := make([][]byte, NodeWidth)
		values[CodeHashVectorPosition] = emptyHashCode
		var c1poly [NodeWidth]Fr
		if _, err := fillSuffixTreePoly(c1poly[:], values[:NodeWidth/2]); err != nil {
			panic(err)
		}
		EmptyCodeHashPoint = *cfg.CommitToPoly(c1poly[:], 0)
		EmptyCodeHashFirstHalfValue = c1poly[EmptyCodeHashFirstHalfIdx]
		EmptyCodeHashSecondHalfValue = c1poly[EmptyCodeHashSecondHalfIdx]
	})
	return cfg
}

func (conf *IPAConfig) CommitToPoly(poly []Fr, _ int) *Point {
	ret := conf.conf.Commit(poly)
	return &ret
}
