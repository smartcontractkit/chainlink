package text

// This code is a modified version of the rgbterm package:
// see https://github.com/aybabtme/rgbterm/blob/master/LICENSE
// Original license:
//
// The MIT License

// Copyright (c) 2014, Antoine Grondin.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// > Except for functions HSLtoRGB and RGBtoHSL, under a BSD 2 clause license:

// Copyright (c) 2012 Rodrigo Moraes. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//      * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//      * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//      * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

var (
	before   = []byte("\033[")
	after    = []byte("m")
	reset    = []byte("\033[0;00m")
	fgcolors = fgTermRGB[16:232]
	bgcolors = bgTermRGB[16:232]
)

// FgString colorizes the foreground of the input with the terminal color
// that matches the closest the RGB color.
//
// This is simply a helper for Bytes.
func FgString(in string, r, g, b uint8) string {
	return string(FgBytes([]byte(in), r, g, b))
}

// BgString colorizes the background of the input with the terminal color
// that matches the closest the RGB color.
//
// This is simply a helper for Bytes.
func BgString(in string, r, g, b uint8) string {
	return string(BgBytes([]byte(in), r, g, b))
}

// Bytes colorizes the foreground with the terminal color that matches
// the closest the RGB color.
func FgBytes(in []byte, r, g, b uint8) []byte {
	return colorize(color_(r, g, b, true), in)
}

// BgBytes colorizes the background of the input with the terminal color
// that matches the closest the RGB color.
func BgBytes(in []byte, r, g, b uint8) []byte {
	return colorize(color_(r, g, b, false), in)
}

func colorize(color, in []byte) []byte {
	return append(append(append(append(before, color...), after...), in...), reset...)
}

func color_(r, g, b uint8, foreground bool) []byte {
	// if all colors are equal, it might be in the grayscale range
	if r == g && g == b {
		color, ok := grayscale(r, foreground)
		if ok {
			return color
		}
	}

	// the general case approximates RGB by using the closest color.
	r6 := ((uint16(r) * 5) / 255)
	g6 := ((uint16(g) * 5) / 255)
	b6 := ((uint16(b) * 5) / 255)
	i := 36*r6 + 6*g6 + b6
	if foreground {
		return fgcolors[i]
	} else {
		return bgcolors[i]
	}
}

func grayscale(scale uint8, foreground bool) ([]byte, bool) {
	var source [256][]byte

	if foreground {
		source = fgTermRGB
	} else {
		source = bgTermRGB
	}

	switch scale {
	case 0x08:
		return source[232], true
	case 0x12:
		return source[233], true
	case 0x1c:
		return source[234], true
	case 0x26:
		return source[235], true
	case 0x30:
		return source[236], true
	case 0x3a:
		return source[237], true
	case 0x44:
		return source[238], true
	case 0x4e:
		return source[239], true
	case 0x58:
		return source[240], true
	case 0x62:
		return source[241], true
	case 0x6c:
		return source[242], true
	case 0x76:
		return source[243], true
	case 0x80:
		return source[244], true
	case 0x8a:
		return source[245], true
	case 0x94:
		return source[246], true
	case 0x9e:
		return source[247], true
	case 0xa8:
		return source[248], true
	case 0xb2:
		return source[249], true
	case 0xbc:
		return source[250], true
	case 0xc6:
		return source[251], true
	case 0xd0:
		return source[252], true
	case 0xda:
		return source[253], true
	case 0xe4:
		return source[254], true
	case 0xee:
		return source[255], true
	}
	return nil, false
}
