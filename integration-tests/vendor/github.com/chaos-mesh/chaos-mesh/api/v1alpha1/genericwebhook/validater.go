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

package genericwebhook

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type FieldValidator interface {
	Validate(root interface{}, path *field.Path) field.ErrorList
}

// Validate would walk through all the fields of target struct recursively, and validate the value with validator declared with struct tag "webhook".
//
// Parameter obj should be a pointer to a data struct.
//
// Validate should return an empty field.ErrorList if all the fields are valid, or return each field.Error for every invalid values.
func Validate(obj interface{}) field.ErrorList {
	// TODO: how to resolve invalid input, for example: obj is a pointer to pointer
	errorList := field.ErrorList{}

	root := obj
	walker := NewFieldWalker(obj, func(path *field.Path, obj interface{}, field *reflect.StructField) bool {
		webhookAttr := ""
		if field != nil {
			webhookAttr = field.Tag.Get("webhook")
		}
		attributes := strings.Split(webhookAttr, ",")

		webhook := ""
		nilable := false
		if len(attributes) > 0 {
			webhook = attributes[0]
		}
		if len(attributes) > 1 {
			nilable = attributes[1] == "nilable"
		}

		validator := getValidator(obj, webhook, nilable)
		if validator != nil {
			if err := validator.Validate(root, path); err != nil {
				errorList = append(errorList, err...)
			}
		}

		return true
	})
	walker.Walk()

	return errorList
}

func Aggregate(errs field.ErrorList) error {
	if errs == nil || len(errs) == 0 {
		return nil
	}
	return errors.New(errs.ToAggregate().Error())
}

func getValidator(obj interface{}, webhook string, nilable bool) FieldValidator {
	// There are two possible situations:
	// 1. The field is a value (int, string, normal struct, etc), and the obj is the reference of it.
	// 2. The field is a pointer to a value or a slice, then the obj is itself.

	val := reflect.ValueOf(obj)

	if validator, ok := obj.(FieldValidator); ok {
		if nilable || !val.IsZero() {
			return validator
		}
	}

	if webhook != "" {
		webhookImpl := webhooks[webhook]

		v := val.Convert(webhookImpl).Interface()
		if validator, ok := v.(FieldValidator); ok {
			if nilable || !val.IsZero() {
				return validator
			}
		}
	}

	return nil
}
