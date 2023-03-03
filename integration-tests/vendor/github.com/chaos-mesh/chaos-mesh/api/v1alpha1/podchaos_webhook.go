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
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// validateContainerNames validates the ContainerNames
func (in *PodChaosSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if in.Action == ContainerKillAction {
		if len(in.ContainerSelector.ContainerNames) == 0 {
			err := errors.Wrapf(errInvalidValue, "the name of container is required on %s action", in.Action)
			allErrs = append(allErrs, field.Invalid(path.Child("containerNames"), in.ContainerNames, err.Error()))
		}
	}
	return allErrs
}
