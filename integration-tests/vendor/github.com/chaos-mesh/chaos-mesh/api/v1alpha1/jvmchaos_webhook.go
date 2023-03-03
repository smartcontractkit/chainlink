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
)

const DefaultJVMAgentPort int32 = 9277

func (in *JVMChaosSpec) Default(root interface{}, field *reflect.StructField) {
	if in == nil {
		return
	}

	jvmChaos := root.(*JVMChaos)
	if len(in.Name) == 0 {
		in.Name = jvmChaos.Name
	}

	if in.Port == 0 {
		in.Port = DefaultJVMAgentPort
	}
}

func (in *JVMChaosSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	switch in.Action {
	case JVMStressAction:
		if in.CPUCount == 0 && len(in.MemoryType) == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "must set one of cpu-count and mem-type when action is 'stress'"))
		}

		if in.CPUCount > 0 && len(in.MemoryType) > 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "inject stress on both CPU and memory is not support now"))
		}

		if len(in.MemoryType) != 0 {
			if in.MemoryType != "stack" && in.MemoryType != "heap" {
				allErrs = append(allErrs, field.Invalid(path, in, "value should be 'stack' or 'heap'"))
			}
		}
	case JVMGCAction:
		// do nothing
	case JVMExceptionAction, JVMReturnAction, JVMLatencyAction:
		if len(in.Class) == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "class not provided"))
		}

		if len(in.Method) == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "method not provided"))
		}
		if in.Action == JVMExceptionAction && len(in.ThrowException) == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "exception not provided"))
		} else if in.Action == JVMReturnAction && len(in.ReturnValue) == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "value not provided"))
		} else if in.Action == JVMLatencyAction && in.LatencyDuration == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "latency not provided"))
		}

	case JVMRuleDataAction:
		if len(in.RuleData) == 0 {
			allErrs = append(allErrs, field.Invalid(path, in, "rule data not provide"))
		}
	case "":
		allErrs = append(allErrs, field.Invalid(path, in, "action not provided"))
	default:
		allErrs = append(allErrs, field.Invalid(path, in, fmt.Sprintf("action %s not supported, action can be 'latency', 'exception', 'return', 'stress', 'gc' or 'ruleData'", in.Action)))
	}

	return allErrs
}
