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

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RisingWaveNodeConfigurationConfigMapSource refers to a ConfigMap where the RisingWave configuration is stored.
type RisingWaveNodeConfigurationConfigMapSource struct {
	// Name determines the ConfigMap to provide the configs RisingWave requests. It will be mounted on the Pods
	// directly. It the ConfigMap isn't provided, the controller will use empty value as the configs.
	// +optional
	Name string `json:"name,omitempty"`

	// Key to the configuration file. Defaults to `risingwave.toml`.
	// +kubebuilder:default=risingwave.toml
	// +optional
	Key string `json:"key,omitempty"`

	// Optional determines if the key must exist in the ConfigMap. Defaults to false.
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

// RisingWaveNodeConfigurationSecretSource refers to a Secret where the RisingWave configuration is stored.
type RisingWaveNodeConfigurationSecretSource struct {
	// Name determines the Secret to provide the configs RisingWave requests. It will be mounted on the Pods
	// directly. It the Secret isn't provided, the controller will use empty value as the configs.
	// +optional
	Name string `json:"name,omitempty"`

	// Key to the configuration file. Defaults to `risingwave.toml`.
	// +kubebuilder:default=risingwave.toml
	// +optional
	Key string `json:"key,omitempty"`

	// Optional determines if the key must exist in the Secret. Defaults to false.
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

// RisingWaveNodeConfiguration determines where the configurations are from, either ConfigMap, Secret, or raw string.
type RisingWaveNodeConfiguration struct {
	// ConfigMap where the `risingwave.toml` locates.
	ConfigMap *RisingWaveNodeConfigurationConfigMapSource `json:"configMap,omitempty"`

	// Secret where the `risingwave.toml` locates.
	Secret *RisingWaveNodeConfigurationSecretSource `json:"secret,omitempty"`

	// Value of the `risingwave.toml`.
	Value *string `json:"value,omitempty"`
}

// RisingWaveSpec is the specification of a RisingWave object.
type RisingWaveSpec struct {
	// UseKruiseWorkloads determines which workload objects are used to run the nodes. It defaults to false, which means
	// the builtin workloads such as Deployment and StatefulSet are used.
	// +optional
	UseKruiseWorkloads *bool `json:"useKruiseWorkloads,omitempty"`

	// SyncPrometheusServiceMonitor determines if the default ServiceMonitor from the Prometheus Operator will be synced
	// during the reconciliation. Defaults to false.
	// +optional
	SyncPrometheusServiceMonitor *bool `json:"syncPrometheusServiceMonitor,omitempty"`

	// FrontendServiceType determines the service type of the frontend service. Defaults to ClusterIP.
	// +optional
	// +kubebuilder:default=ClusterIP
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	FrontendServiceType corev1.ServiceType `json:"frontendServiceType,omitempty"`

	// AdditionalFrontendServiceMetadata tells the operator to add the specified metadata onto the frontend Service.
	// Note that the system reserved labels and annotations are not valid and will be rejected by the webhook.
	AdditionalFrontendServiceMetadata PartialObjectMeta `json:"additionalFrontendServiceMetadata,omitempty"`

	// MetaStore determines which backend the meta store will use and the parameters for it. Defaults to memory.
	// But keep in mind that memory backend is not recommended in production.
	// +kubebuilder:default={memory: true}
	MetaStore RisingWaveMetaStoreBackend `json:"metaStore,omitempty"`

	// StateStore determines which backend the state store will use and the parameters for it. Defaults to memory.
	// But keep in mind that memory backend is not recommended in production.
	// +kubebuilder:default={memory: true}
	StateStore RisingWaveStateStoreBackend `json:"stateStore,omitempty"`

	// Image of the RisingWave.
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// PodTemplate determines how Pods of the RisingWave should be run. Leave it to empty if no advanced configurations
	// should be applied to the Pods. Note that setting the image field here isn't allowed.
	// +optional
	PodTemplate RisingWaveNodePodTemplate `json:"podTemplate,omitempty"`

	// Configuration determines the configuration to be used for the RisingWave nodes.
	Configuration RisingWaveNodeConfiguration `json:"configuration,omitempty"`

	// MetaComponent determines how and how many meta Pods should be run by setting the node groups inside.
	MetaComponent RisingWaveComponent `json:"metaComponent,omitempty"`

	// FrontendComponent determines how and how many frontend Pods should be run by setting the node groups inside.
	FrontendComponent RisingWaveComponent `json:"frontendComponent,omitempty"`

	// ComputeComponent determines how and how many compute Pods should be run by setting the node groups inside.
	ComputeComponent RisingWaveComponent `json:"computeComponent,omitempty"`

	// CompactorComponent determines how and how many compactor Pods should be run by setting the node groups inside.
	CompactorComponent RisingWaveComponent `json:"compactorComponent,omitempty"`
}

type RisingWaveConditionType string

const (
	// RisingWaveAvailable means the RisingWave is available for streaming and serving, i.e., all compute nodes
	// are running and ready, at lease one of the meta nodes are running and ready, and at lease one of the frontend
	// nodes are running and ready.
	RisingWaveAvailable RisingWaveConditionType = "Available"

	// RisingWaveProgressing means the RisingWave is progressing. Progress for a RisingWave is considered when any of
	// the underlying workloads are progressing.
	RisingWaveProgressing RisingWaveConditionType = "Progressing"

	// RisingWaveFailure means the RisingWave fails to progress.
	RisingWaveFailure RisingWaveConditionType = "Failure"
)

// RisingWaveCondition shows the current status of a condition of the RisingWave.
type RisingWaveCondition struct {
	// Type of RisingWave condition.
	Type RisingWaveConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=DeploymentConditionType"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,7,opt,name=lastTransitionTime"`

	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// A human-readable message indicating details about the transition.
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// RisingWaveStatus is the status of the RisingWave.
type RisingWaveStatus struct {
	// The generation observed by the RisingWave controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Image tag of the global image under the .spec.
	ImageTag string `json:"imageTag,omitempty"`

	// Represents the latest available observations of a RisingWave's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []RisingWaveCondition `json:"conditions,omitempty"`

	// Status of the meta store.
	MetaStore RisingWaveMetaStoreStatus `json:"metaStore,omitempty"`

	// Status of the state store.
	StateStore RisingWaveStateStoreStatus `json:"stateStore,omitempty"`

	// Status of the meta component.
	MetaComponent RisingWaveComponentStatus `json:"metaComponent,omitempty"`

	// Status of the frontend component.
	FrontendComponent RisingWaveComponentStatus `json:"frontendComponent,omitempty"`

	// Status of the compute component.
	ComputeComponent RisingWaveComponentStatus `json:"computeComponent,omitempty"`

	// Status of the compactor component.
	CompactorComponent RisingWaveComponentStatus `json:"compactorComponent,omitempty"`

	// Lock records maintained by the RisingWaveScaleView controller.
	// +listType=map
	// +listMapKey=name
	ScaleViewLocks []RisingWaveScaleViewLock `json:"scaleViewLocks,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// RisingWave is a RisingWave deployment consists of a collection of Pods. Generally, it contains the meta, frontend,
// compute, and compactor nodes.
type RisingWave struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RisingWaveSpec   `json:"spec,omitempty"`
	Status RisingWaveStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RisingWaveList contains a list of the RisingWaves.
type RisingWaveList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RisingWave `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RisingWave{}, &RisingWaveList{})
}
