// Copyright 2021 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package v1alpha1

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1/genericwebhook"
)

type IOErrno uint32

func (in *IOErrno) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	obj := root.(*IOChaos)
	if obj.Spec.Action == IoFaults {
		// in cannot be nil
		if *in == 0 {
			allErrs = append(allErrs, field.Invalid(path, in,
				fmt.Sprintf("action %s: errno 0 is not supported", obj.Spec.Action)))
		}
	}
	return allErrs
}

func init() {
	genericwebhook.Register("IOErrno", reflect.PtrTo(reflect.TypeOf(IOErrno(0))))
}
