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

// RisingWaveMetaStoreBackendType is the type for the meta store backends.
type RisingWaveMetaStoreBackendType string

// All valid meta store backend types.
const (
	RisingWaveMetaStoreBackendTypeMemory RisingWaveMetaStoreBackendType = "Memory"
	RisingWaveMetaStoreBackendTypeEtcd   RisingWaveMetaStoreBackendType = "Etcd"
)

// RisingWaveEtcdCredentials is the reference and keys selector to the etcd access credentials stored in a local secret.
type RisingWaveEtcdCredentials struct {
	// The name of the secret in the pod's namespace to select from.
	SecretName string `json:"secretName"`

	// UsernameKeyRef is the key of the secret to be the username. Must be a valid secret key. Defaults to "Username".
	UsernameKeyRef *string `json:"usernameKeyRef,omitempty"`

	// PasswordKeyRef is the key of the secret to be the password. Must be a valid secret key.
	// Defaults to "Password".
	PasswordKeyRef *string `json:"passwordKeyRef,omitempty"`
}

// RisingWaveMetaStoreBackendEtcd is the collection of parameters for the etcd backend meta store.
type RisingWaveMetaStoreBackendEtcd struct {
	// RisingWaveEtcdCredentials is the credentials provider from a Secret. It could be optional to mean that
	// the etcd service could be accessed without authentication.
	// +optional
	*RisingWaveEtcdCredentials `json:"credentials,omitempty"`

	// Endpoints are the endpoints of the etcd service, separated with comma. A valid endpoint should not contain any scheme prefix.
	Endpoints string `json:"endpoints"`
}

// RisingWaveMetaStoreBackend is the collection of parameters for the meta store that RisingWave uses. Note that one
// and only one of the first-level fields could be set.
type RisingWaveMetaStoreBackend struct {
	// Memory determines whether RisingWave uses a memory-based meta store. Keep in mind that the memory
	// backend is only for test purposes and should not be used in production. Defaults to false.
	Memory *bool `json:"memory,omitempty"`

	// Etcd determine whether RisingWave uses the etcd-backed meta store and the parameters for accessing the etcd.
	Etcd *RisingWaveMetaStoreBackendEtcd `json:"etcd,omitempty"`
}

// RisingWaveMetaStoreStatus is the status of the meta store.
type RisingWaveMetaStoreStatus struct {
	// Backend type of the meta store.
	Backend RisingWaveMetaStoreBackendType `json:"backend,omitempty"`
}
