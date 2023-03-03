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
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="action",type=string,JSONPath=`.spec.action`
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment
// +chaos-mesh:oneshot=in.Spec.Action==NodeReset

// GCPChaos is the Schema for the gcpchaos API
type GCPChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPChaosSpec   `json:"spec"`
	Status GCPChaosStatus `json:"status,omitempty"`
}

var _ InnerObjectWithCustomStatus = (*GCPChaos)(nil)
var _ InnerObjectWithSelector = (*GCPChaos)(nil)
var _ InnerObject = (*GCPChaos)(nil)

// GCPChaosAction represents the chaos action about gcp.
type GCPChaosAction string

const (
	// NodeStop represents the chaos action of stopping the node.
	NodeStop GCPChaosAction = "node-stop"
	// NodeReset represents the chaos action of resetting the node.
	NodeReset GCPChaosAction = "node-reset"
	// DiskLoss represents the chaos action of detaching the disk.
	DiskLoss GCPChaosAction = "disk-loss"
)

// GCPChaosSpec is the content of the specification for a GCPChaos
type GCPChaosSpec struct {
	// Action defines the specific gcp chaos action.
	// Supported action: node-stop / node-reset / disk-loss
	// Default action: node-stop
	// +kubebuilder:validation:Enum=node-stop;node-reset;disk-loss
	Action GCPChaosAction `json:"action"`

	// Duration represents the duration of the chaos action.
	// +optional
	Duration *string `json:"duration,omitempty" webhook:"Duration"`

	// SecretName defines the name of kubernetes secret. It is used for GCP credentials.
	// +optional
	SecretName *string `json:"secretName,omitempty"`

	GCPSelector `json:",inline"`
}

type GCPSelector struct {
	// Project defines the ID of gcp project.
	Project string `json:"project"`

	// Zone defines the zone of gcp project.
	Zone string `json:"zone"`

	// Instance defines the name of the instance
	Instance string `json:"instance"`

	// The device name of disks to detach.
	// Needed in disk-loss.
	// +ui:form:when=action=='disk-loss'
	// +optional
	DeviceNames []string `json:"deviceNames,omitempty" webhook:"GCPDeviceNames,nilable"`
}

func (obj *GCPChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.GCPSelector,
	}
}

func (selector *GCPSelector) Id() string {
	// TODO: handle the error here
	// or ignore it is enough ?
	json, _ := json.Marshal(selector)

	return string(json)
}

// GCPChaosStatus represents the status of a GCPChaos
type GCPChaosStatus struct {
	ChaosStatus `json:",inline"`

	// The attached disk info strings.
	// Needed in disk-loss.
	AttachedDisksStrings []string `json:"attachedDiskStrings,omitempty"`
}

func (obj *GCPChaos) GetCustomStatus() interface{} {
	return &obj.Status.AttachedDisksStrings
}
