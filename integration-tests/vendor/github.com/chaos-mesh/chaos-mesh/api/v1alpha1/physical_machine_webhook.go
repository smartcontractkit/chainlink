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
	"net/url"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var physicalMachineLog = logf.Log.WithName("physical-machine-resource")

func (in *PhysicalMachineSpec) Default(root interface{}, field *reflect.StructField) {
	if in == nil {
		return
	}

	physicalMachineLog.Info("default", "address", in.Address)
	// add http prefix for address
	if len(in.Address) > 0 && !strings.HasPrefix(in.Address, "http") {
		in.Address = fmt.Sprintf("http://%s", in.Address)
	}
}

func (in *PhysicalMachineSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// make sure address is not empty
	if len(in.Address) == 0 {
		allErrs = append(allErrs,
			field.Invalid(path.Child("address"), in.Address, "the address is required"))
	}

	if _, err := url.ParseRequestURI(in.Address); err != nil {
		allErrs = append(allErrs,
			field.Invalid(path.Child("address"), in.Address, fmt.Sprintf("the address is invalid, %s", err)))
	}
	return allErrs
}
