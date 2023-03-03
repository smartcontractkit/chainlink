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
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1/genericwebhook"
)

const (
	// DefaultJitter defines default value for jitter
	DefaultJitter = "0ms"

	// DefaultCorrelation defines default value for correlation
	DefaultCorrelation = "0"
)

func (in *Direction) Default(root interface{}, field *reflect.StructField) {
	if *in == "" {
		*in = To
	}
}

type Rate string

// validateBandwidth validates the bandwidth
func (in *Rate) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// in cannot be nil
	_, err := ConvertUnitToBytes(string(*in))

	if err != nil {
		allErrs = append(allErrs,
			field.Invalid(path, in,
				fmt.Sprintf("parse rate field error:%s", err)))
	}
	return allErrs
}

func ConvertUnitToBytes(nu string) (uint64, error) {
	// normalize input
	s := strings.ToLower(strings.TrimSpace(nu))

	for i, u := range []string{"tbps", "gbps", "mbps", "kbps", "bps"} {
		if strings.HasSuffix(s, u) {
			ts := strings.TrimSuffix(s, u)
			s := strings.TrimSpace(ts)

			n, err := strconv.ParseUint(s, 10, 64)

			if err != nil {
				return 0, err
			}

			// convert unit to bytes
			for j := 4 - i; j > 0; j-- {
				n = n * 1024
			}

			return n, nil
		}
	}

	return 0, errors.New("invalid unit")
}

// ValidateTargets validates externalTargets and Targets
func (in *NetworkChaosSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if in.Action == PartitionAction {
		return nil
	}

	if (in.Direction == From || in.Direction == Both) &&
		in.ExternalTargets != nil && in.Action != PartitionAction {
		allErrs = append(allErrs,
			field.Invalid(path.Child("direction"), in.Direction,
				"external targets cannot be used with `from` and `both` direction in netem action yet"))
	}

	if (in.Direction == From || in.Direction == Both) && in.Target == nil {
		if in.Action != PartitionAction {
			allErrs = append(allErrs,
				field.Invalid(path.Child("direction"), in.Direction,
					"`from` and `both` direction cannot be used when targets is empty in netem action"))
		} else if in.ExternalTargets == nil {
			allErrs = append(allErrs,
				field.Invalid(path.Child("direction"), in.Direction,
					"`from` and `both` direction cannot be used when targets and external targets are both empty"))
		}
	}

	// TODO: validate externalTargets are in ip or domain form
	return allErrs
}

func init() {
	genericwebhook.Register("Rate", reflect.PtrTo(reflect.TypeOf(Rate(""))))
}
