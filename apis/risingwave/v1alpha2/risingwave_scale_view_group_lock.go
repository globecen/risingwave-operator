// Copyright 2023 RisingWave Labs
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

package v1alpha2

import "k8s.io/apimachinery/pkg/types"

// RisingWaveScaleViewNodeGroupLock stands for a record of a locked node group by some RisingWaveScaleView object.
type RisingWaveScaleViewNodeGroupLock struct {
	// Name of the node group.
	Name string `json:"name,omitempty"`

	// Replicas of the node group. The replicas here means the current allowed value for the following
	// updates on the `replicas` field of the corresponding group.
	Replicas int32 `json:"replicas,omitempty"`
}

// RisingWaveScaleViewReference is the reference to a RisingWaveScaleView. It also includes the generation field
// to indicate on which generation that the current RisingWaveScaleView was observed by the controller.
type RisingWaveScaleViewReference struct {
	// Name of the RisingWaveScaleView object.
	Name string `json:"name,omitempty"`

	// UID of the RisingWaveScaleView object.
	UID types.UID `json:"uid,omitempty"`

	// ObservedGeneration of the RisingWaveScaleView object that last observed.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// RisingWaveScaleViewLock stands for a lock to some specified node groups of one component owned by the RisingWaveScaleView.
// The lock aims to prevent bidirectional updates on the `replicas` of these node groups when a user wants to run auto-scaling
// staff on these node groups. The webhooks of the RisingWave resource will reject any update that is not allowed by the
// lock record to guarantee the semantic of the lock.
type RisingWaveScaleViewLock struct {
	// Reference to the RisingWaveScaleView object.
	Reference RisingWaveScaleViewReference `json:",inline"`

	// Component that the RisingWaveScaleView targets.
	Component string `json:"component,omitempty"`

	// Locks that owned by the RisingWaveScaleView.
	// +listType=map
	// +listMapKey=name
	Locks []RisingWaveScaleViewNodeGroupLock `json:"locks,omitempty"`
}
