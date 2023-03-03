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
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:object:generate=false

// ChaosKindMap defines a map including all chaos kinds.
type chaosKindMap struct {
	sync.RWMutex
	kinds map[string]*ChaosKind
}

func (c *chaosKindMap) register(name string, kind *ChaosKind) {
	c.Lock()
	defer c.Unlock()
	c.kinds[name] = kind
}

// clone will build a new map with kinds, so if user add or delete entries of the map, the origin global map will not be affected.
func (c *chaosKindMap) clone() map[string]*ChaosKind {
	c.RLock()
	defer c.RUnlock()

	out := make(map[string]*ChaosKind)
	for key, kind := range c.kinds {
		out[key] = &ChaosKind{
			chaos: kind.chaos,
			list:  kind.list,
		}
	}

	return out
}

// AllKinds returns all chaos kinds, key is name of Kind, value is an accessor for spawning Object and List
func AllKinds() map[string]*ChaosKind {
	return all.clone()
}

func AllKindsIncludeScheduleAndWorkflow() map[string]*ChaosKind {
	all := chaosKindMap{
		kinds: all.clone(),
	}
	all.register(KindSchedule, &ChaosKind{
		chaos: &Schedule{},
		list:  &ScheduleList{},
	})
	all.register(KindWorkflow, &ChaosKind{
		chaos: &Workflow{},
		list:  &WorkflowList{},
	})

	return all.kinds
}

// all is a ChaosKindMap instance.
var all = &chaosKindMap{
	kinds: make(map[string]*ChaosKind),
}

// +kubebuilder:object:generate=false

// ChaosKind includes one kind of chaos and its list type
type ChaosKind struct {
	chaos client.Object
	list  GenericChaosList
}

// SpawnObject will deepcopy a clean struct for the acquired kind as placeholder
func (it *ChaosKind) SpawnObject() client.Object {
	return it.chaos.DeepCopyObject().(client.Object)
}

// SpawnList will deepcopy a clean list for the acquired kind of chaos as placeholder
func (it *ChaosKind) SpawnList() GenericChaosList {
	return it.list.DeepCopyList()
}

// AllKinds returns all chaos kinds.
func AllScheduleItemKinds() map[string]*ChaosKind {
	return allScheduleItem.clone()
}

// allScheduleItem is a ChaosKindMap instance.
var allScheduleItem = &chaosKindMap{
	kinds: make(map[string]*ChaosKind),
}
