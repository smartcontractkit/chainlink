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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	LabelControlledBy = "chaos-mesh.org/controlled-by"
	LabelWorkflow     = "chaos-mesh.org/workflow"
)

const KindWorkflowNode = "WorkflowNode"

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=wfn
// +kubebuilder:subresource:status
type WorkflowNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a node of workflow
	Spec WorkflowNodeSpec `json:"spec"`

	// +optional
	// Most recently observed status of the workflow node
	Status WorkflowNodeStatus `json:"status"`
}

type WorkflowNodeSpec struct {
	TemplateName string       `json:"templateName"`
	WorkflowName string       `json:"workflowName"`
	Type         TemplateType `json:"type"`
	StartTime    *metav1.Time `json:"startTime"`
	// +optional
	Deadline *metav1.Time `json:"deadline,omitempty"`
	// +optional
	Task *Task `json:"task,omitempty"`
	// +optional
	Children []string `json:"children,omitempty"`
	// +optional
	ConditionalBranches []ConditionalBranch `json:"conditionalBranches,omitempty"`
	// +optional
	*EmbedChaos `json:",inline,omitempty"`
	// +optional
	Schedule *ScheduleSpec `json:"schedule,omitempty"`
}

type WorkflowNodeStatus struct {

	// ChaosResource refs to the real chaos CR object.
	// +optional
	ChaosResource *corev1.TypedLocalObjectReference `json:"chaosResource,omitempty"`

	// ConditionalBranchesStatus records the evaluation result of each ConditionalBranch
	// +optional
	ConditionalBranchesStatus *ConditionalBranchesStatus `json:"conditionalBranchesStatus,omitempty"`

	// ActiveChildren means the created children node
	// +optional
	ActiveChildren []corev1.LocalObjectReference `json:"activeChildren,omitempty"`

	// Children is necessary for representing the order when replicated child template references by parent template.
	// +optional
	FinishedChildren []corev1.LocalObjectReference `json:"finishedChildren,omitempty"`

	// Represents the latest available observations of a workflow node's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []WorkflowNodeCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type ConditionalBranch struct {
	// Target is the name of other template, if expression is evaluated as true, this template will be spawned.
	Target string `json:"target"`
	// Expression is the expression for this conditional branch, expected type of result is boolean. If expression is empty, this branch will always be selected/the template will be spawned.
	// +optional
	Expression string `json:"expression,omitempty"`
}

type ConditionalBranchesStatus struct {
	// +optional
	Branches []ConditionalBranchStatus `json:"branches"`
	// +optional
	Context []string `json:"context"`
}

type ConditionalBranchStatus struct {
	Target           string                 `json:"target"`
	EvaluationResult corev1.ConditionStatus `json:"evaluationResult"`
}

type WorkflowNodeConditionType string

const (
	ConditionAccomplished   WorkflowNodeConditionType = "Accomplished"
	ConditionDeadlineExceed WorkflowNodeConditionType = "DeadlineExceed"
	ConditionChaosInjected  WorkflowNodeConditionType = "ChaosInjected"
)

type WorkflowNodeCondition struct {
	Type   WorkflowNodeConditionType `json:"type"`
	Status corev1.ConditionStatus    `json:"status"`
	Reason string                    `json:"reason"`
}

// +kubebuilder:object:root=true
type WorkflowNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkflowNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkflowNode{}, &WorkflowNodeList{})
}

// Reasons
const (
	EntryCreated                string = "EntryCreated"
	InvalidEntry                string = "InvalidEntry"
	WorkflowAccomplished        string = "WorkflowAccomplished"
	NodeAccomplished            string = "NodeAccomplished"
	NodesCreated                string = "NodesCreated"
	NodeDeadlineExceed          string = "NodeDeadlineExceed"
	NodeDeadlineNotExceed       string = "NodeDeadlineNotExceed"
	NodeDeadlineOmitted         string = "NodeDeadlineOmitted"
	ParentNodeDeadlineExceed    string = "ParentNodeDeadlineExceed"
	ChaosCRCreated              string = "ChaosCRCreated"
	ChaosCRCreateFailed         string = "ChaosCRCreateFailed"
	ChaosCRDeleted              string = "ChaosCRDeleted"
	ChaosCRDeleteFailed         string = "ChaosCRDeleteFailed"
	ChaosCRNotExists            string = "ChaosCRNotExists"
	TaskPodSpawned              string = "TaskPodSpawned"
	TaskPodSpawnFailed          string = "TaskPodSpawnFailed"
	TaskPodPodCompleted         string = "TaskPodPodCompleted"
	ConditionalBranchesSelected string = "ConditionalBranchesSelected"
	RerunBySpecChanged          string = "RerunBySpecChanged"
)

// GenericChaosList only use to list GenericChaos by certain EmbedChaos
// +kubebuilder:object:generate=false
type GenericChaosList interface {
	runtime.Object
	metav1.ListInterface
	GetItems() []GenericChaos
	DeepCopyList() GenericChaosList
}

// GenericChaos could be a place holder for any kubernetes Kind
// +kubebuilder:object:generate=false
type GenericChaos interface {
	runtime.Object
	metav1.Object
}
