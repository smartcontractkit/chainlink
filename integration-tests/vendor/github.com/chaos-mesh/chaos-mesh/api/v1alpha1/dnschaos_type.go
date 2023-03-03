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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DNSChaosAction represents the chaos action about DNS.
type DNSChaosAction string

const (
	// ErrorAction represents get error when send DNS request.
	ErrorAction DNSChaosAction = "error"

	// RandomAction represents get random IP when send DNS request.
	RandomAction DNSChaosAction = "random"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="action",type=string,JSONPath=`.spec.action`
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment

// DNSChaos is the Schema for the networkchaos API
type DNSChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a pod chaos experiment
	Spec DNSChaosSpec `json:"spec"`

	// +optional
	// Most recently observed status of the chaos experiment about pods
	Status DNSChaosStatus `json:"status"`
}

var _ InnerObjectWithSelector = (*DNSChaos)(nil)
var _ InnerObject = (*DNSChaos)(nil)

// DNSChaosSpec defines the desired state of DNSChaos
type DNSChaosSpec struct {
	// Action defines the specific DNS chaos action.
	// Supported action: error, random
	// Default action: error
	// +kubebuilder:validation:Enum=error;random
	Action DNSChaosAction `json:"action"`

	ContainerSelector `json:",inline"`

	// Duration represents the duration of the chaos action
	Duration *string `json:"duration,omitempty" webhook:"Duration"`

	// Choose which domain names to take effect, support the placeholder ? and wildcard *, or the Specified domain name.
	// Note:
	//      1. The wildcard * must be at the end of the string. For example, chaos-*.org is invalid.
	//      2. if the patterns is empty, will take effect on all the domain names.
	// For example:
	// 		The value is ["google.com", "github.*", "chaos-mes?.org"],
	// 		will take effect on "google.com", "github.com" and "chaos-mesh.org"
	// +optional
	DomainNamePatterns []string `json:"patterns"`
}

// DNSChaosStatus defines the observed state of DNSChaos
type DNSChaosStatus struct {
	ChaosStatus `json:",inline"`
}

func (obj *DNSChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.ContainerSelector,
	}
}
