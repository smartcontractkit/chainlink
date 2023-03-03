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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// PodIOChaosSpec defines the desired state of IOChaos
type PodIOChaosSpec struct {
	// VolumeMountPath represents the target mount path
	// It must be a root of mount path now.
	// TODO: search the mount parent of any path automatically.
	// TODO: support multiple different volume mount path in one pod
	VolumeMountPath string `json:"volumeMountPath"`

	// TODO: support multiple different container to inject in one pod
	// +optional
	Container *string `json:"container,omitempty"`

	// Actions are a list of IOChaos actions
	// +optional
	Actions []IOChaosAction `json:"actions,omitempty"`
}

// IOChaosAction defines an possible action of IOChaos
type IOChaosAction struct {
	Type IOChaosType `json:"type"`

	Filter `json:",inline"`

	// Faults represents the fault to inject
	// +optional
	Faults []IoFault `json:"faults,omitempty"`

	// Latency represents the latency to inject
	// +optional
	Latency string `json:"latency,omitempty"`

	// AttrOverride represents the attribution to override
	// +optional
	*AttrOverrideSpec `json:",inline"`

	// MistakeSpec represents the mistake to inject
	// +optional
	*MistakeSpec `json:"mistake,omitempty"`

	// Source represents the source of current rules
	Source string `json:"source,omitempty"`
}

// IOChaosType represents the type of an IOChaos Action
type IOChaosType string

const (
	// IoLatency represents injecting latency for io operation
	IoLatency IOChaosType = "latency"

	// IoFaults represents injecting faults for io operation
	IoFaults IOChaosType = "fault"

	// IoAttrOverride represents replacing attribution for io operation
	IoAttrOverride IOChaosType = "attrOverride"

	// IoMistake represents injecting incorrect read or write for io operation
	IoMistake IOChaosType = "mistake"
)

// Filter represents a filter of IOChaos action, which will define the
// scope of an IOChaosAction
type Filter struct {
	// Path represents a glob of injecting path
	Path string `json:"path"`

	// Methods represents the method that the action will inject in
	// +optional
	Methods []IoMethod `json:"methods,omitempty"`

	// Percent represents the percent probability of injecting this action
	Percent int `json:"percent"`
}

// IoFault represents the fault to inject and their weight
type IoFault struct {
	Errno  uint32 `json:"errno"`
	Weight int32  `json:"weight"`
}

// AttrOverrideSpec represents an override of attribution
type AttrOverrideSpec struct {
	//+optional
	Ino *uint64 `json:"ino,omitempty"`
	//+optional
	Size *uint64 `json:"size,omitempty"`
	//+optional
	Blocks *uint64 `json:"blocks,omitempty"`
	//+optional
	Atime *Timespec `json:"atime,omitempty"`
	//+optional
	Mtime *Timespec `json:"mtime,omitempty"`
	//+optional
	Ctime *Timespec `json:"ctime,omitempty"`
	//+optional
	Kind *FileType `json:"kind,omitempty"`
	//+optional
	Perm *uint16 `json:"perm,omitempty"`
	//+optional
	Nlink *uint32 `json:"nlink,omitempty"`
	//+optional
	UID *uint32 `json:"uid,omitempty"`
	//+optional
	GID *uint32 `json:"gid,omitempty"`
	//+optional
	Rdev *uint32 `json:"rdev,omitempty"`
}

// MistakeSpec represents one type of mistake
type MistakeSpec struct {
	// Filling determines what is filled in the miskate data.
	// +optional
	// +kubebuilder:validation:Enum=zero;random
	Filling FillingType `json:"filling,omitempty"`

	// There will be [1, MaxOccurrences] segments of wrong data.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxOccurrences int64 `json:"maxOccurrences,omitempty"`

	// Max length of each wrong data segment in bytes
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxLength int64 `json:"maxLength,omitempty"`
}

// FillingType represents type of data is filled for incorrectness
type FillingType string

const (
	// All zero
	Zero FillingType = "zero"

	// Random octets
	Random FillingType = "random"
)

// Timespec represents a time
type Timespec struct {
	Sec  int64 `json:"sec"`
	Nsec int64 `json:"nsec"`
}

// FileType represents type of a file
type FileType string

const (
	NamedPipe   FileType = "namedPipe"
	CharDevice  FileType = "charDevice"
	BlockDevice FileType = "blockDevice"
	Directory   FileType = "directory"
	RegularFile FileType = "regularFile"
	TSymlink    FileType = "symlink"
	Socket      FileType = "socket"
)

type IoMethod string

const (
	LookUp      IoMethod = "lookup"
	Forget      IoMethod = "forget"
	GetAttr     IoMethod = "getattr"
	SetAttr     IoMethod = "setattr"
	ReadLink    IoMethod = "readlink"
	Mknod       IoMethod = "mknod"
	Mkdir       IoMethod = "mkdir"
	UnLink      IoMethod = "unlink"
	Rmdir       IoMethod = "rmdir"
	MSymlink    IoMethod = "symlink"
	Rename      IoMethod = "rename"
	Link        IoMethod = "link"
	Open        IoMethod = "open"
	Read        IoMethod = "read"
	Write       IoMethod = "write"
	Flush       IoMethod = "flush"
	Release     IoMethod = "release"
	Fsync       IoMethod = "fsync"
	Opendir     IoMethod = "opendir"
	Readdir     IoMethod = "readdir"
	Releasedir  IoMethod = "releasedir"
	Fsyncdir    IoMethod = "fsyncdir"
	Statfs      IoMethod = "statfs"
	SetXAttr    IoMethod = "setxattr"
	GetXAttr    IoMethod = "getxattr"
	ListXAttr   IoMethod = "listxattr"
	RemoveXAttr IoMethod = "removexattr"
	Access      IoMethod = "access"
	Create      IoMethod = "create"
	GetLk       IoMethod = "getlk"
	SetLk       IoMethod = "setlk"
	Bmap        IoMethod = "bmap"
)

// +chaos-mesh:base
// +chaos-mesh:webhook:enableUpdate
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// PodIOChaos is the Schema for the podiochaos API
type PodIOChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PodIOChaosSpec `json:"spec,omitempty"`

	//+optional
	Status PodIOChaosStatus `json:"status,omitempty"`
}

type PodIOChaosStatus struct {

	// Pid represents a running toda process id
	// +optional
	Pid int64 `json:"pid,omitempty"`

	// StartTime represents the start time of a toda process
	// +optional
	StartTime int64 `json:"startTime,omitempty"`

	// +optional
	FailedMessage string `json:"failedMessage,omitempty"`

	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true

// PodIOChaosList contains a list of PodIOChaos
type PodIOChaosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodIOChaos `json:"items"`
}
