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

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +chaos-mesh:base
// +chaos-mesh:webhook:enableUpdate

// PodNetworkChaos is the Schema for the PodNetworkChaos API
type PodNetworkChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a pod chaos experiment
	Spec PodNetworkChaosSpec `json:"spec"`

	// +optional
	// Most recently observed status of the chaos experiment about pods
	Status PodNetworkChaosStatus `json:"status"`
}

// PodNetworkChaosSpec defines the desired state of PodNetworkChaos
type PodNetworkChaosSpec struct {
	// The ipset on the pod
	// +optional
	IPSets []RawIPSet `json:"ipsets,omitempty"`

	// The iptables rules on the pod
	// +optional
	Iptables []RawIptables `json:"iptables,omitempty"`

	// The tc rules on the pod
	// +optional
	TrafficControls []RawTrafficControl `json:"tcs,omitempty"`
}

// RawIPSet represents an ipset on specific pod
type RawIPSet struct {
	// The name of ipset
	Name string `json:"name"`

	// The contents of ipset
	Cidrs []string `json:"cidrs"`

	// The name and namespace of the source network chaos
	RawRuleSource `json:",inline"`
}

// ChainDirection represents the direction of chain
type ChainDirection string

const (
	// Input means this chain is linked with INPUT chain
	Input ChainDirection = "input"

	// Output means this chain is linked with OUTPUT chain
	Output ChainDirection = "output"
)

// RawIptables represents the iptables rules on specific pod
type RawIptables struct {
	// The name of iptables chain
	Name string `json:"name"`

	// The name of related ipset
	// +optional
	// +nullable
	IPSets []string `json:"ipsets,omitempty"`

	// The block direction of this iptables rule
	Direction ChainDirection `json:"direction"`

	// Device represents the network device to be affected.
	// +optional
	Device string `json:"device,omitempty"`

	RawRuleSource `json:",inline"`
}

// TcType the type of traffic control
type TcType string

const (
	// Netem represents netem traffic control
	Netem TcType = "netem"

	// Bandwidth represents bandwidth shape traffic control
	Bandwidth TcType = "bandwidth"
)

// RawTrafficControl represents the traffic control chaos on specific pod
type RawTrafficControl struct {
	// The type of traffic control
	Type TcType `json:"type"`

	TcParameter `json:",inline"`

	// The name of target ipset
	// +optional
	IPSet string `json:"ipset,omitempty"`

	// The name and namespace of the source network chaos
	Source string `json:"source"`

	// Device represents the network device to be affected.
	// +optional
	Device string `json:"device,omitempty"`
}

// TcParameter represents the parameters for a traffic control chaos
type TcParameter struct {
	// Delay represents the detail about delay action
	// +optional
	Delay *DelaySpec `json:"delay,omitempty"`

	// Loss represents the detail about loss action
	// +optional
	Loss *LossSpec `json:"loss,omitempty"`

	// DuplicateSpec represents the detail about loss action
	// +optional
	Duplicate *DuplicateSpec `json:"duplicate,omitempty"`

	// Corrupt represents the detail about corrupt action
	// +optional
	Corrupt *CorruptSpec `json:"corrupt,omitempty"`

	// Bandwidth represents the detail about bandwidth control action
	// +optional
	Bandwidth *BandwidthSpec `json:"bandwidth,omitempty"`
}

// RawRuleSource represents the name and namespace of the source network chaos
type RawRuleSource struct {
	Source string `json:"source"`
}

// PodNetworkChaosStatus defines the observed state of PodNetworkChaos
type PodNetworkChaosStatus struct {
	FailedMessage string `json:"failedMessage,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true

// PodNetworkChaosList contains a list of PodNetworkChaos
type PodNetworkChaosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodNetworkChaos `json:"items"`
}
