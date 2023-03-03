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

// PhysicalMachineChaosAction represents the chaos action about physical machine.
type PhysicalMachineChaosAction string

var (
	PMStressCPUAction        PhysicalMachineChaosAction = "stress-cpu"
	PMStressMemAction        PhysicalMachineChaosAction = "stress-mem"
	PMDiskWritePayloadAction PhysicalMachineChaosAction = "disk-write-payload"
	PMDiskReadPayloadAction  PhysicalMachineChaosAction = "disk-read-payload"
	PMDiskFillAction         PhysicalMachineChaosAction = "disk-fill"
	PMNetworkCorruptAction   PhysicalMachineChaosAction = "network-corrupt"
	PMNetworkDuplicateAction PhysicalMachineChaosAction = "network-duplicate"
	PMNetworkLossAction      PhysicalMachineChaosAction = "network-loss"
	PMNetworkDelayAction     PhysicalMachineChaosAction = "network-delay"
	PMNetworkPartitionAction PhysicalMachineChaosAction = "network-partition"
	PMNetworkBandwidthAction PhysicalMachineChaosAction = "network-bandwidth"
	PMNetworkDNSAction       PhysicalMachineChaosAction = "network-dns"
	PMProcessAction          PhysicalMachineChaosAction = "process"
	PMJVMExceptionAction     PhysicalMachineChaosAction = "jvm-exception"
	PMJVMGCAction            PhysicalMachineChaosAction = "jvm-gc"
	PMJVMLatencyAction       PhysicalMachineChaosAction = "jvm-latency"
	PMJVMReturnAction        PhysicalMachineChaosAction = "jvm-return"
	PMJVMStressAction        PhysicalMachineChaosAction = "jvm-stress"
	PMJVMRuleDataAction      PhysicalMachineChaosAction = "jvm-rule-data"
	PMClockAction            PhysicalMachineChaosAction = "clock"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="action",type=string,JSONPath=`.spec.action`
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment

// PhysicalMachineChaos is the Schema for the physical machine chaos API
type PhysicalMachineChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a physical machine chaos experiment
	Spec PhysicalMachineChaosSpec `json:"spec"`

	// +optional
	// Most recently observed status of the chaos experiment
	Status PhysicalMachineChaosStatus `json:"status"`
}

// PhysicalMachineChaosSpec defines the desired state of PhysicalMachineChaos
type PhysicalMachineChaosSpec struct {
	// +kubebuilder:validation:Enum=stress-cpu;stress-mem;disk-read-payload;disk-write-payload;disk-fill;network-corrupt;network-duplicate;network-loss;network-delay;network-partition;network-dns;network-bandwidth;process;jvm-exception;jvm-gc;jvm-latency;jvm-return;jvm-stress;jvm-rule-data;clock
	Action PhysicalMachineChaosAction `json:"action"`

	PhysicalMachineSelector `json:",inline"`

	// ExpInfo string `json:"expInfo"`
	ExpInfo `json:",inline"`

	// Duration represents the duration of the chaos action
	// +optional
	Duration *string `json:"duration,omitempty" webhook:"Duration"`
}

// PhysicalMachineChaosStatus defines the observed state of PhysicalMachineChaos
type PhysicalMachineChaosStatus struct {
	ChaosStatus `json:",inline"`
}

func (obj *PhysicalMachineChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.PhysicalMachineSelector,
	}
}

type PhysicalMachineSelector struct {
	// DEPRECATED: Use Selector instead.
	// Only one of Address and Selector could be specified.
	// +optional
	Address []string `json:"address,omitempty"`

	// Selector is used to select physical machines that are used to inject chaos action.
	// +optional
	Selector PhysicalMachineSelectorSpec `json:"selector"`

	// Mode defines the mode to run chaos action.
	// Supported mode: one / all / fixed / fixed-percent / random-max-percent
	// +kubebuilder:validation:Enum=one;all;fixed;fixed-percent;random-max-percent
	Mode SelectorMode `json:"mode"`

	// Value is required when the mode is set to `FixedMode` / `FixedPercentMode` / `RandomMaxPercentMode`.
	// If `FixedMode`, provide an integer of physical machines to do chaos action.
	// If `FixedPercentMode`, provide a number from 0-100 to specify the percent of physical machines the server can do chaos action.
	// IF `RandomMaxPercentMode`,  provide a number from 0-100 to specify the max percent of pods to do chaos action
	// +optional
	Value string `json:"value,omitempty"`
}

// PhysicalMachineSelectorSpec defines some selectors to select objects.
// If the all selectors are empty, all objects will be used in chaos experiment.
type PhysicalMachineSelectorSpec struct {
	GenericSelectorSpec `json:",inline"`

	// PhysicalMachines is a map of string keys and a set values that used to select physical machines.
	// The key defines the namespace which physical machine belong,
	// and each value is a set of physical machine names.
	// +optional
	PhysicalMachines map[string][]string `json:"physicalMachines,omitempty"`
}

func (spec *PhysicalMachineSelectorSpec) Empty() bool {
	if spec == nil {
		return true
	}
	if len(spec.AnnotationSelectors) != 0 || len(spec.FieldSelectors) != 0 || len(spec.LabelSelectors) != 0 ||
		len(spec.Namespaces) != 0 || len(spec.PhysicalMachines) != 0 || len(spec.ExpressionSelectors) != 0 {
		return false
	}
	return true
}

type ExpInfo struct {
	// the experiment ID
	// +optional
	UID string `json:"uid,omitempty" swaggerignore:"true"`

	// the subAction, generate automatically
	// +optional
	Action string `json:"action,omitempty" swaggerignore:"true"`

	// +ui:form:when=action=='stress-cpu'
	// +optional
	StressCPU *StressCPUSpec `json:"stress-cpu,omitempty"`

	// +ui:form:when=action=='stress-mem'
	// +optional
	StressMemory *StressMemorySpec `json:"stress-mem,omitempty"`

	// +ui:form:when=action=='disk-read-payload'
	// +optional
	DiskReadPayload *DiskPayloadSpec `json:"disk-read-payload,omitempty"`

	// +ui:form:when=action=='disk-write-payload'
	// +optional
	DiskWritePayload *DiskPayloadSpec `json:"disk-write-payload,omitempty"`

	// +ui:form:when=action=='disk-fill'
	// +optional
	DiskFill *DiskFillSpec `json:"disk-fill,omitempty"`

	// +ui:form:when=action=='network-corrupt'
	// +optional
	NetworkCorrupt *NetworkCorruptSpec `json:"network-corrupt,omitempty"`

	// +ui:form:when=action=='network-duplicate'
	// +optional
	NetworkDuplicate *NetworkDuplicateSpec `json:"network-duplicate,omitempty"`

	// +ui:form:when=action=='network-loss'
	// +optional
	NetworkLoss *NetworkLossSpec `json:"network-loss,omitempty"`

	// +ui:form:when=action=='network-delay'
	// +optional
	NetworkDelay *NetworkDelaySpec `json:"network-delay,omitempty"`

	// +ui:form:when=action=='network-partition'
	// +optional
	NetworkPartition *NetworkPartitionSpec `json:"network-partition,omitempty"`

	// +ui:form:when=action=='network-dns'
	// +optional
	NetworkDNS *NetworkDNSSpec `json:"network-dns,omitempty"`

	// +ui:form:when=action=='network-bandwidth'
	// +optional
	NetworkBandwidth *NetworkBandwidthSpec `json:"network-bandwidth,omitempty"`

	// +ui:form:when=action=='process'
	// +optional
	Process *ProcessSpec `json:"process,omitempty"`

	// +ui:form:when=action=='jvm-exception'
	// +optional
	JVMException *JVMExceptionSpec `json:"jvm-exception,omitempty"`

	// +ui:form:when=action=='jvm-gc'
	// +optional
	JVMGC *JVMGCSpec `json:"jvm-gc,omitempty"`

	// +ui:form:when=action=='jvm-latency'
	// +optional
	JVMLatency *JVMLatencySpec `json:"jvm-latency,omitempty"`

	// +ui:form:when=action=='jvm-return'
	// +optional
	JVMReturn *JVMReturnSpec `json:"jvm-return,omitempty"`

	// +ui:form:when=action=='jvm-stress'
	// +optional
	JVMStress *JVMStressSpec `json:"jvm-stress,omitempty"`

	// +ui:form:when=action=='jvm-rule-data'
	// +optional
	JVMRuleData *JVMRuleDataSpec `json:"jvm-rule-data,omitempty"`

	// +ui:form:when=action=='clock'
	// +optional
	Clock *ClockSpec `json:"clock,omitempty"`
}

type StressCPUSpec struct {
	// specifies P percent loading per CPU worker. 0 is effectively a sleep (no load) and 100 is full loading.
	Load int `json:"load,omitempty"`
	// specifies N workers to apply the stressor.
	Workers int `json:"workers,omitempty"`
	// extend stress-ng options
	Options string `json:"options,omitempty"`
}

type StressMemorySpec struct {
	// specifies N bytes consumed per vm worker, default is the total available memory.
	// One can specify the size as % of total available memory or in units of B, KB/KiB, MB/MiB, GB/GiB, TB/TiB..
	Size string `json:"size,omitempty"`
	// extend stress-ng options
	Options string `json:"options,omitempty"`
}

type DiskFileSpec struct {
	// specifies how many units of data will write into the file path. support unit: c=1, w=2, b=512, kB=1000,
	// K=1024, MB=1000*1000, M=1024*1024, GB=1000*1000*1000, G=1024*1024*1024 BYTES. example : 1M | 512kB
	Size string `json:"size,omitempty"`
	// specifies the location to fill data in. if path not provided,
	// payload will read/write from/into a temp file, temp file will be deleted after writing
	Path string `json:"path,omitempty"`
}

type DiskPayloadSpec struct {
	DiskFileSpec `json:",inline"`

	// specifies the number of process work on writing, default 1, only 1-255 is valid value
	PayloadProcessNum uint8 `json:"payload-process-num,omitempty"`
}

type DiskFillSpec struct {
	DiskFileSpec `json:",inline"`

	// fill disk by fallocate
	FillByFallocate bool `json:"fill-by-fallocate,omitempty"`
}

type NetworkCommonSpec struct {
	// correlation is percentage (10 is 10%)
	Correlation string `json:"correlation,omitempty"`
	// the network interface to impact
	Device string `json:"device,omitempty"`
	// only impact egress traffic from these source ports, use a ',' to separate or to indicate the range, such as 80, 8001:8010.
	// it can only be used in conjunction with -p tcp or -p udp
	SourcePort string `json:"source-port,omitempty"`
	// only impact egress traffic to these destination ports, use a ',' to separate or to indicate the range, such as 80, 8001:8010.
	// it can only be used in conjunction with -p tcp or -p udp
	EgressPort string `json:"egress-port,omitempty"`
	// only impact egress traffic to these IP addresses
	IPAddress string `json:"ip-address,omitempty"`
	// only impact traffic using this IP protocol, supported: tcp, udp, icmp, all
	IPProtocol string `json:"ip-protocol,omitempty"`
	// only impact traffic to these hostnames
	Hostname string `json:"hostname,omitempty"`
}

type NetworkCorruptSpec struct {
	NetworkCommonSpec `json:",inline"`

	// percentage of packets to corrupt (10 is 10%)
	Percent string `json:"percent,omitempty"`
}

type NetworkDuplicateSpec struct {
	NetworkCommonSpec `json:",inline"`

	// percentage of packets to duplicate (10 is 10%)
	Percent string `json:"percent,omitempty"`
}

type NetworkLossSpec struct {
	NetworkCommonSpec `json:",inline"`

	// percentage of packets to loss (10 is 10%)
	Percent string `json:"percent,omitempty"`
}

type NetworkDelaySpec struct {
	NetworkCommonSpec `json:",inline"`

	// jitter time, time units: ns, us (or µs), ms, s, m, h.
	Jitter string `json:"jitter,omitempty"`
	// delay egress time, time units: ns, us (or µs), ms, s, m, h.
	Latency string `json:"latency,omitempty"`
}

type NetworkPartitionSpec struct {
	// the network interface to impact
	Device string `json:"device,omitempty"`
	// only impact traffic to these hostnames
	Hostname string `json:"hostname,omitempty"`
	// only impact egress traffic to these IP addresses
	IPAddress string `json:"ip-address,omitempty"`
	// specifies the partition direction, values can be 'from', 'to'.
	// 'from' means packets coming from the 'IPAddress' or 'Hostname' and going to your server,
	// 'to' means packets originating from your server and going to the 'IPAddress' or 'Hostname'.
	Direction string `json:"direction,omitempty"`
	// only impact egress traffic to these IP addresses
	IPProtocol string `json:"ip-protocol,omitempty"`
	// only the packet which match the tcp flag can be accepted, others will be dropped.
	// only set when the IPProtocol is tcp, used for partition.
	AcceptTCPFlags string `json:"accept-tcp-flags,omitempty"`
}

type NetworkDNSSpec struct {
	// update the DNS server in /etc/resolv.conf with this value
	DNSServer string `json:"dns-server,omitempty"`
	// map specified host to this IP address
	DNSIp string `json:"dns-ip,omitempty"`
	// map this host to specified IP
	DNSDomainName string `json:"dns-domain-name,omitempty"`
}

type NetworkBandwidthSpec struct {
	Rate string `json:"rate"`
	// +kubebuilder:validation:Minimum=1
	Limit uint32 `json:"limit"`
	// +kubebuilder:validation:Minimum=1
	Buffer uint32 `json:"buffer"`

	Peakrate *uint64 `json:"peakrate,omitempty"`
	Minburst *uint32 `json:"minburst,omitempty"`

	Device    string `json:"device,omitempty"`
	IPAddress string `json:"ip-address,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
}

type ProcessSpec struct {
	// the process name or the process ID
	Process string `json:"process,omitempty"`
	// the signal number to send
	Signal int `json:"signal,omitempty"`
}

type JVMCommonSpec struct {
	// the port of agent server
	Port int `json:"port,omitempty"`

	// the pid of Java process which need to attach
	Pid int `json:"pid,omitempty"`
}

type JVMClassMethodSpec struct {
	// Java class
	Class string `json:"class,omitempty"`

	// the method in Java class
	Method string `json:"method,omitempty"`
}

type JVMExceptionSpec struct {
	JVMCommonSpec      `json:",inline"`
	JVMClassMethodSpec `json:",inline"`

	// the exception which needs to throw for action `exception`
	ThrowException string `json:"exception,omitempty"`
}

type JVMGCSpec struct {
	JVMCommonSpec `json:",inline"`
}

type JVMLatencySpec struct {
	JVMCommonSpec      `json:",inline"`
	JVMClassMethodSpec `json:",inline"`

	// the latency duration for action 'latency', unit ms
	LatencyDuration int `json:"latency,omitempty"`
}

type JVMReturnSpec struct {
	JVMCommonSpec      `json:",inline"`
	JVMClassMethodSpec `json:",inline"`

	// the return value for action 'return'
	ReturnValue string `json:"value,omitempty"`
}

type JVMStressSpec struct {
	JVMCommonSpec `json:",inline"`

	// the CPU core number need to use, only set it when action is stress
	CPUCount int `json:"cpu-count,omitempty"`

	// the memory type need to locate, only set it when action is stress, the value can be 'stack' or 'heap'
	MemoryType string `json:"mem-type,omitempty"`
}

type JVMRuleDataSpec struct {
	JVMCommonSpec `json:",inline"`

	// RuleData used to save the rule file's data, will use it when recover
	RuleData string `json:"rule-data,omitempty"`
}

type ClockSpec struct {
	// the pid of target program.
	Pid int `json:"pid,omitempty"`
	// specifies the length of time offset.
	TimeOffset string `json:"time-offset,omitempty"`
	// the identifier of the particular clock on which to act.
	// More clock description in linux kernel can be found in man page of clock_getres, clock_gettime, clock_settime.
	// Muti clock ids should be split with ","
	ClockIdsSlice string `json:"clock-ids-slice,omitempty"`
}
