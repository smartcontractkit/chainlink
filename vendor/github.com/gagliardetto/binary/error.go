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

import "reflect"

// An InvalidDecoderError describes an invalid argument passed to Decoder.
// (The argument to Decoder must be a non-nil pointer.)
type InvalidDecoderError struct {
	Type reflect.Type
}

func (e *InvalidDecoderError) Error() string {
	if e.Type == nil {
		return "decoder: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "decoder: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "decoder: Decode(nil " + e.Type.String() + ")"
}
