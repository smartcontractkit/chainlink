//
// Copyright (c) 2016-2020 The Aurora Authors. All rights reserved.
// This program is free software. It comes without any warranty,
// to the extent permitted by applicable law. You can redistribute
// it and/or modify it under the terms of the Unlicense. See LICENSE
// file for more details or see below.
//

//
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
// For more information, please refer to <http://unlicense.org/>
//

package aurora

import (
	"fmt"
)

// Sprintf allows to use Value as format. For example
//
//    v := Sprintf(Red("total: +3.5f points"), Blue(3.14))
//
// In this case "total:" and "points" will be red, but
// 3.14 will be blue. But, in another example
//
//    v := Sprintf(Red("total: +3.5f points"), 3.14)
//
// full string will be red. And no way to clear 3.14 to
// default format and color
func Sprintf(format interface{}, args ...interface{}) string {
	switch ft := format.(type) {
	case string:
		return fmt.Sprintf(ft, args...)
	case Value:
		for i, v := range args {
			if val, ok := v.(Value); ok {
				args[i] = val.setTail(ft.Color())
				continue
			}
		}
		return fmt.Sprintf(ft.String(), args...)
	}
	// unknown type of format (we hope it's a string)
	return fmt.Sprintf(fmt.Sprint(format), args...)
}
