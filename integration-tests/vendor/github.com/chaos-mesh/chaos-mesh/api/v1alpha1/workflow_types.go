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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=wf
// +kubebuilder:subresource:status
// +chaos-mesh:experiment
type Workflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a workflow
	Spec WorkflowSpec `json:"spec"`

	// +optional
	// Most recently observed status of the workflow
	Status WorkflowStatus `json:"status"`
}

func (in *Workflow) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

const KindWorkflow = "Workflow"

type WorkflowSpec struct {
	Entry     string     `json:"entry"`
	Templates []Template `json:"templates"`
}

type WorkflowStatus struct {
	// +optional
	EntryNode *string `json:"entryNode,omitempty"`
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// +optional
	EndTime *metav1.Time `json:"endTime,omitempty"`
	// Represents the latest available observations of a workflow's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []WorkflowCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type WorkflowConditionType string

const (
	WorkflowConditionAccomplished WorkflowConditionType = "Accomplished"
	WorkflowConditionScheduled    WorkflowConditionType = "Scheduled"
)

type WorkflowCondition struct {
	Type      WorkflowConditionType  `json:"type"`
	Status    corev1.ConditionStatus `json:"status"`
	Reason    string                 `json:"reason"`
	StartTime *metav1.Time           `json:"startTime,omitempty"`
}

type TemplateType string

const (
	TypeTask     TemplateType = "Task"
	TypeSerial   TemplateType = "Serial"
	TypeParallel TemplateType = "Parallel"
	TypeSuspend  TemplateType = "Suspend"
	TypeSchedule TemplateType = "Schedule"
)

func IsChaosTemplateType(target TemplateType) bool {
	return contains(allChaosTemplateType, target)
}

func contains(arr []TemplateType, target TemplateType) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

type Template struct {
	Name string       `json:"name"`
	Type TemplateType `json:"templateType"`
	// +optional
	Deadline *string `json:"deadline,omitempty"`
	// Task describes the behavior of the custom task. Only used when Type is TypeTask.
	// +optional
	Task *Task `json:"task,omitempty"`
	// Children describes the children steps of serial or parallel node. Only used when Type is TypeSerial or TypeParallel.
	// +optional
	Children []string `json:"children,omitempty"`
	// ConditionalBranches describes the conditional branches of custom tasks. Only used when Type is TypeTask.
	// +optional
	ConditionalBranches []ConditionalBranch `json:"conditionalBranches,omitempty"`
	// EmbedChaos describe the chaos to be injected with chaos nodes. Only used when Type is Type<Something>Chaos.
	// +optional
	*EmbedChaos `json:",inline"`
	// Schedule describe the Schedule(describing scheduled chaos) to be injected with chaos nodes. Only used when Type is TypeSchedule.
	// +optional
	Schedule *ChaosOnlyScheduleSpec `json:"schedule,omitempty"`
}

// ChaosOnlyScheduleSpec is very similar with ScheduleSpec, but it could not schedule Workflow
// because we could not resolve nested CRD now
type ChaosOnlyScheduleSpec struct {
	Schedule string `json:"schedule"`

	// +optional
	// +nullable
	// +kubebuilder:validation:Minimum=0
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds"`

	// +optional
	// +kubebuilder:validation:Enum=Forbid;Allow
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	HistoryLimit int `json:"historyLimit,omitempty"`

	// TODO: use a custom type, as `TemplateType` contains other possible values
	Type ScheduleTemplateType `json:"type"`

	EmbedChaos `json:",inline"`
}

type Task struct {
	// Container is the main container image to run in the pod
	Container *corev1.Container `json:"container,omitempty"`

	// Volumes is a list of volumes that can be mounted by containers in a template.
	// +patchStrategy=merge
	// +patchMergeKey=name
	Volumes []corev1.Volume `json:"volumes,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// TODO: maybe we could specify parameters in other ways, like loading context from file
}

// +kubebuilder:object:root=true
type WorkflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workflow `json:"items"`
}

func (in *WorkflowList) GetItems() []GenericChaos {
	var result []GenericChaos
	for _, item := range in.Items {
		item := item
		result = append(result, &item)
	}
	return result
}

// TODO: refactor: not so accurate
func (in *WorkflowList) DeepCopyList() GenericChaosList {
	return in.DeepCopy()
}

func init() {
	SchemeBuilder.Register(&Workflow{}, &WorkflowList{})
}

func FetchChaosByTemplateType(templateType TemplateType) (runtime.Object, error) {
	if kind, ok := all.kinds[string(templateType)]; ok {
		return kind.SpawnObject(), nil
	}
	return nil, errors.Wrapf(errInvalidValue, "unknown template type %s", templateType)
}
