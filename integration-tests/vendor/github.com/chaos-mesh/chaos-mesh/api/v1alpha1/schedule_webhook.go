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

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var schedulelog = logf.Log.WithName("schedule-resource")

// +kubebuilder:webhook:path=/mutate-chaos-mesh-org-v1alpha1-schedule,mutating=true,failurePolicy=fail,groups=chaos-mesh.org,resources=schedule,verbs=create;update,versions=v1alpha1,name=mschedule.kb.io

var _ webhook.Defaulter = &Schedule{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Schedule) Default() {
	schedulelog.Info("default", "name", in.Name)
	in.Spec.ConcurrencyPolicy.Default()
}

func (in *ConcurrencyPolicy) Default() {
	if *in == "" {
		*in = ForbidConcurrent
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-chaos-mesh-org-v1alpha1-schedule,mutating=false,failurePolicy=fail,groups=chaos-mesh.org,resources=schedule,versions=v1alpha1,name=vschedule.kb.io

var _ webhook.Validator = &Schedule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Schedule) ValidateCreate() error {
	schedulelog.Info("validate create", "name", in.Name)
	return in.Validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Schedule) ValidateUpdate(old runtime.Object) error {
	schedulelog.Info("validate update", "name", in.Name)
	return in.Validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Schedule) ValidateDelete() error {
	schedulelog.Info("validate delete", "name", in.Name)
	return nil
}

// Validate validates chaos object
func (in *Schedule) Validate() error {
	allErrs := in.Spec.Validate()
	if len(allErrs) > 0 {
		return errors.New(allErrs.ToAggregate().Error())
	}
	return nil
}

func (in *ScheduleSpec) Validate() field.ErrorList {
	specField := field.NewPath("spec")
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, in.validateSchedule(specField.Child("schedule"))...)
	allErrs = append(allErrs, in.validateChaos(specField)...)
	return allErrs
}

// validateSchedule validates the cron
func (in *ScheduleSpec) validateSchedule(schedule *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	_, err := StandardCronParser.Parse(in.Schedule)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(schedule,
			in.Schedule,
			fmt.Sprintf("parse schedule field error:%s", err)))
	}

	return allErrs
}

// validateChaos validates the chaos
func (in *ScheduleSpec) validateChaos(chaos *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if in.Type != ScheduleTypeWorkflow {
		allErrs = append(allErrs, in.EmbedChaos.Validate(chaos, string(in.Type))...)
	}
	return allErrs
}
