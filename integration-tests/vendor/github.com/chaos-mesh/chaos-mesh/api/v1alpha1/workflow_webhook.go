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
	"sort"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	gw "github.com/chaos-mesh/chaos-mesh/api/v1alpha1/genericwebhook"
)

// log is for logging in this package.
var workflowlog = logf.Log.WithName("workflow-resource")

var _ webhook.Validator = &Workflow{}

func (in *Workflow) ValidateCreate() error {
	var allErrs field.ErrorList
	specPath := field.NewPath("spec")
	allErrs = append(allErrs, entryMustExists(specPath.Child("entry"), in.Spec.Entry, in.Spec.Templates)...)
	allErrs = append(allErrs, validateTemplates(specPath.Child("templates"), in.Spec.Templates)...)
	if len(allErrs) > 0 {
		return errors.New(allErrs.ToAggregate().Error())
	}
	return nil
}

func (in *Workflow) ValidateUpdate(old runtime.Object) error {
	return in.ValidateCreate()
}

func (in *Workflow) ValidateDelete() error {
	return nil
}

func entryMustExists(path *field.Path, entry string, templates []Template) field.ErrorList {
	var result field.ErrorList
	// name is required
	if len(entry) == 0 {
		result = append(result, field.Required(path, "the entry of workflow is required"))
	}
	founded := false
	for _, item := range templates {
		if item.Name == entry {
			founded = true
			break
		}
	}
	if !founded {
		result = append(result, field.Invalid(path, entry, fmt.Sprintf("can not find a template with name %s", entry)))
	}
	return result
}

func validateTemplates(path *field.Path, templates []Template) field.ErrorList {
	var result field.ErrorList
	if len(templates) == 0 {
		result = append(result, field.Invalid(path, templates, "templates in workflow could not be empty"))
		return result
	}
	var allNames []string
	for _, template := range templates {
		allNames = append(allNames, template.Name)
	}
	result = append(result, namesCouldNotBeDuplicated(path, allNames)...)

	for i, item := range templates {
		itemPath := path.Index(i)
		result = append(result, validateTemplate(itemPath, item, templates)...)
	}
	return result
}

func validateTemplate(path *field.Path, template Template, allTemplates []Template) field.ErrorList {
	var result field.ErrorList
	// name is required
	if len(template.Name) == 0 {
		result = append(result, field.Required(path.Child("name"), "name of template is required"))
	}

	// name must be restricted with DNS-1123
	errs := validation.IsDNS1123Subdomain(template.Name)
	if len(errs) > 0 {
		result = append(result, field.Invalid(path.Child("name"), template.Name, fmt.Sprintf("field name must be DNS-1123 subdomain, %s", errs)))
	}

	// template name could not be duplicated

	switch templateType := template.Type; {
	case templateType == TypeSuspend:
		if template.Deadline == nil || len(*template.Deadline) == 0 {
			result = append(result, field.Invalid(path.Child("deadline"), template.Deadline, "deadline in template with type Suspend could not be empty"))
		}
		result = append(result, shouldBeNoTask(path, template)...)
		result = append(result, shouldBeNoChildren(path, template)...)
		result = append(result, shouldBeNoConditionalBranches(path, template)...)
		result = append(result, shouldBeNoEmbedChaos(path, template)...)
		result = append(result, shouldBeNoSchedule(path, template)...)
	case templateType == TypeSerial, templateType == TypeParallel:
		for i, item := range template.Children {
			result = append(result, templateMustExists(item, path.Child("children").Index(i), allTemplates)...)
		}
		result = append(result, shouldBeNoTask(path, template)...)
		result = append(result, shouldBeNoConditionalBranches(path, template)...)
		result = append(result, shouldBeNoEmbedChaos(path, template)...)
		result = append(result, shouldBeNoSchedule(path, template)...)
	case templateType == TypeSchedule:
		result = append(result, shouldBeNoTask(path, template)...)
		result = append(result, shouldBeNoChildren(path, template)...)
		result = append(result, shouldBeNoConditionalBranches(path, template)...)
		result = append(result, shouldBeNoEmbedChaos(path, template)...)
	case templateType == TypeTask:
		result = append(result, shouldBeNoChildren(path, template)...)
		result = append(result, shouldBeNoEmbedChaos(path, template)...)
		result = append(result, shouldBeNoSchedule(path, template)...)
	case IsChaosTemplateType(templateType):
		result = append(result, shouldNotSetupDurationInTheChaos(path, template)...)

		result = append(result, shouldBeNoTask(path, template)...)
		result = append(result, shouldBeNoChildren(path, template)...)
		result = append(result, shouldBeNoConditionalBranches(path, template)...)
		result = append(result, shouldBeNoSchedule(path, template)...)

		result = append(result, template.EmbedChaos.Validate(path, string(templateType))...)
	default:
		result = append(result, field.Invalid(path.Child("templateType"), template.Type, fmt.Sprintf("unrecognized template type: %s", template.Type)))
	}

	return result
}

func namesCouldNotBeDuplicated(templatesPath *field.Path, names []string) field.ErrorList {
	nameCounter := make(map[string]int)
	for _, name := range names {
		if count, ok := nameCounter[name]; ok {
			nameCounter[name] = count + 1
		} else {
			nameCounter[name] = 1
		}
	}
	var duplicatedNames []string
	for name, count := range nameCounter {
		if count > 1 {
			duplicatedNames = append(duplicatedNames, name)
		}
	}
	sort.Strings(duplicatedNames)
	if len(duplicatedNames) > 0 {
		return field.ErrorList{
			field.Invalid(templatesPath, "", fmt.Sprintf("template name must be unique, duplicated names: %s", duplicatedNames)),
		}
	}
	return nil
}

func templateMustExists(templateName string, path *field.Path, template []Template) field.ErrorList {
	var result field.ErrorList

	founded := false
	for _, item := range template {
		if item.Name == templateName {
			founded = true
			break
		}
	}

	if !founded {
		err := field.Invalid(path, templateName, fmt.Sprintf("can not find a template with name %s", templateName))
		result = append(result, err)
	}
	return result
}

func shouldNotSetupDurationInTheChaos(path *field.Path, template Template) field.ErrorList {
	var result field.ErrorList

	if template.EmbedChaos == nil {
		result = append(result, field.Invalid(path.Child(string(template.Type)), nil, fmt.Sprintf("the value of chaos %s is required", template.Type)))
	}

	spec := reflect.ValueOf(template.EmbedChaos).Elem().FieldByName(string(template.Type))
	if !spec.IsValid() || spec.IsNil() {
		result = append(result, field.Invalid(path.Child(string(template.Type)),
			nil,
			fmt.Sprintf("parse workflow field error: missing chaos spec %s", template.Type)))
		return result
	}
	if commonSpec, ok := spec.Interface().(ContainsDuration); !ok {
		result = append(result, field.Invalid(path, "", fmt.Sprintf("Chaos: %s does not implement CommonSpec", template.Type)))
	} else {
		duration, err := commonSpec.GetDuration()
		if err != nil {
			result = append(result, field.Invalid(path, "", err.Error()))
			return result
		}
		if duration != nil {
			result = append(result, field.Invalid(path, duration, "should not define duration in chaos when using Workflow, use Template#Deadline instead."))
		}
	}
	return result
}

func shouldBeNoTask(path *field.Path, template Template) field.ErrorList {
	if template.Task != nil {
		return field.ErrorList{
			field.Invalid(path, template.Task, "this template should not contain Task"),
		}
	}
	return nil
}

func shouldBeNoChildren(path *field.Path, template Template) field.ErrorList {
	if len(template.Children) > 0 {
		return field.ErrorList{
			field.Invalid(path, template.Children, "this template should not contain Children"),
		}
	}
	return nil
}

func shouldBeNoConditionalBranches(path *field.Path, template Template) field.ErrorList {
	if len(template.ConditionalBranches) > 0 {
		return field.ErrorList{
			field.Invalid(path, template.ConditionalBranches, "this template should not contain ConditionalBranches"),
		}
	}
	return nil
}

func shouldBeNoEmbedChaos(path *field.Path, template Template) field.ErrorList {
	// TODO: we could improve that with code generation in the future
	if template.EmbedChaos != nil {
		return field.ErrorList{
			field.Invalid(path, template.EmbedChaos, "this template should not contain any Chaos"),
		}
	}
	return nil
}

func shouldBeNoSchedule(path *field.Path, template Template) field.ErrorList {
	if template.Schedule != nil {
		return field.ErrorList{
			field.Invalid(path, template.Schedule, "this template should not contain Schedule"),
		}
	}
	return nil
}

func (in *Workflow) Default() {
	gw.Default(in)
}
