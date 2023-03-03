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
	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KindSchedule = "Schedule"

// +kubebuilder:object:root=true

// Schedule is the cronly schedule object
type Schedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ScheduleSpec `json:"spec"`

	// +optional
	Status ScheduleStatus `json:"status"`
}

type ConcurrencyPolicy string

var (
	ForbidConcurrent ConcurrencyPolicy = "Forbid"
	AllowConcurrent  ConcurrencyPolicy = "Allow"
)

func (c ConcurrencyPolicy) IsForbid() bool {
	return c == ForbidConcurrent || c == ""
}

func (c ConcurrencyPolicy) IsAllow() bool {
	return c == AllowConcurrent
}

// ScheduleSpec is the specification of a schedule object
type ScheduleSpec struct {
	Schedule string `json:"schedule"`

	// +optional
	// +nullable
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:ExclusiveMinimum=true
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds"`

	// +optional
	// +kubebuilder:default=Forbid
	// +kubebuilder:validation:Enum=Forbid;Allow
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	HistoryLimit int `json:"historyLimit,omitempty"`

	// TODO: use a custom type, as `TemplateType` contains other possible values
	Type ScheduleTemplateType `json:"type"`

	ScheduleItem `json:",inline"`
}

// ScheduleStatus is the status of a schedule object
type ScheduleStatus struct {
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// +optional
	// +nullable
	LastScheduleTime metav1.Time `json:"time,omitempty"`
}

type ScheduleTemplateType string

func (in *Schedule) IsPaused() bool {
	if in.Annotations == nil || in.Annotations[PauseAnnotationKey] != "true" {
		return false
	}
	return true
}

// +kubebuilder:object:root=true

// ScheduleList contains a list of Schedule
type ScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Schedule `json:"items"`
}

func (in *ScheduleList) GetItems() []GenericChaos {
	var result []GenericChaos
	for _, item := range in.Items {
		item := item
		result = append(result, &item)
	}
	return result
}

func (in *ScheduleList) DeepCopyList() GenericChaosList {
	return in.DeepCopy()
}

var StandardCronParser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

func init() {
	SchemeBuilder.Register(&Schedule{})
	SchemeBuilder.Register(&ScheduleList{})
}
