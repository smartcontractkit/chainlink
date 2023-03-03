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
	"strconv"

	"github.com/docker/go-units"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1/genericwebhook"
)

// Validate validates the scheduler and duration
func (in *StressChaosSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	if len(in.StressngStressors) == 0 && in.Stressors == nil {
		return field.ErrorList{
			field.Invalid(path, in, "missing stressors"),
		}
	}
	return nil
}

// Validate validates whether the Stressors are all well defined
func (in *Stressors) Validate(root interface{}, path *field.Path) field.ErrorList {
	if in == nil {
		return nil
	}

	if in.MemoryStressor == nil && in.CPUStressor == nil {
		return field.ErrorList{
			field.Invalid(path, in, "missing stressors"),
		}
	}
	return nil
}

type Bytes string

func (in *Bytes) Validate(root interface{}, path *field.Path) field.ErrorList {
	packError := func(err error) field.ErrorList {
		return field.ErrorList{
			field.Invalid(path, in, fmt.Sprintf("incorrect bytes format: %s", err.Error())),
		}
	}

	// in cannot be nil
	size := *in
	length := len(size)
	if length == 0 {
		return nil
	}

	var err error
	if size[length-1] == '%' {
		var percent int
		percent, err = strconv.Atoi(string(size)[:length-1])
		if err != nil {
			return packError(err)
		}
		if percent > 100 || percent < 0 {
			err = errors.New("illegal proportion")
			return packError(err)
		}
	} else {
		_, err = units.FromHumanSize(string(size))
		if err != nil {
			return packError(err)
		}
	}

	return nil
}

// Validate validates whether the Stressor is well defined
func (in *Stressor) Validate(parent *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	if in.Workers <= 0 {
		errs = append(errs, field.Invalid(parent, in, "workers should always be positive"))
	}
	return errs
}

func init() {
	genericwebhook.Register("Bytes", reflect.PtrTo(reflect.TypeOf(Bytes(""))))
}
