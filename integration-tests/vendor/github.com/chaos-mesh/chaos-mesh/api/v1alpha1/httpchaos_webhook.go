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
	"net/http"
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1/genericwebhook"
)

type Port int32

func (in *Port) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// in cannot be zero or negative
	if *in <= 0 {
		allErrs = append(allErrs, field.Invalid(path, in, fmt.Sprintf("port %d is not supported", *in)))
	}
	return allErrs
}

type HTTPMethod string

func (in *HTTPMethod) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if in != nil && root.(*HTTPChaos).Spec.Target == PodHttpRequest {
		switch *in {
		case http.MethodGet:
		case http.MethodPost:
		case http.MethodPut:
		case http.MethodDelete:
		case http.MethodPatch:
		case http.MethodHead:
		case http.MethodOptions:
		case http.MethodTrace:
		case http.MethodConnect:
		default:
			allErrs = append(allErrs, field.Invalid(path, in, fmt.Sprintf("method %s is not supported", *in)))
		}
	}
	return allErrs
}

func (in *PodHttpChaosTarget) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	switch *in {
	case PodHttpRequest:
	case PodHttpResponse:
	default:
		allErrs = append(allErrs, field.Invalid(path, in, fmt.Sprintf("target %s is not supported", *in)))
	}
	return allErrs
}

func init() {
	genericwebhook.Register("Port", reflect.PtrTo(reflect.TypeOf(Port(0))))
	genericwebhook.Register("HTTPMethod", reflect.PtrTo(reflect.TypeOf(HTTPMethod(""))))
	genericwebhook.Register("PodHttpChaosTarget", reflect.PtrTo(reflect.TypeOf(PodHttpChaosTarget(""))))
}
