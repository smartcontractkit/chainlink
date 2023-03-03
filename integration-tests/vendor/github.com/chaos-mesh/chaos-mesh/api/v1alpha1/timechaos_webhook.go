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
	"time"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1/genericwebhook"
)

type ClockIds []string

// DefaultClockIds will set default value for empty ClockIds fields
func (in *ClockIds) Default(root interface{}, field *reflect.StructField) {
	// in cannot be nil
	if *in == nil || len(*in) == 0 {
		*in = []string{"CLOCK_REALTIME"}
	}
}

type TimeOffset string

func (in *TimeOffset) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	_, err := time.ParseDuration(string(*in))
	if err != nil {
		allErrs = append(allErrs, field.Invalid(path,
			in,
			fmt.Sprintf("parse timeOffset field error:%s", err)))
	}

	return allErrs
}

func init() {
	genericwebhook.Register("ClockIds", reflect.PtrTo(reflect.TypeOf(ClockIds{})))
	genericwebhook.Register("TimeOffset", reflect.PtrTo(reflect.TypeOf(TimeOffset(""))))
}
