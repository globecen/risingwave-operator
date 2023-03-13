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

// RisingWaveStateStoreBackendType is the type for the state store backends.
type RisingWaveStateStoreBackendType string

// All valid state store backend types.
const (
	RisingWaveStateStoreBackendTypeMemory       RisingWaveStateStoreBackendType = "Memory"
	RisingWaveStateStoreBackendTypeMinIO        RisingWaveStateStoreBackendType = "MinIO"
	RisingWaveStateStoreBackendTypeS3           RisingWaveStateStoreBackendType = "S3"
	RisingWaveStateStoreBackendTypeS3Compatible RisingWaveStateStoreBackendType = "S3c"
	RisingWaveStateStoreBackendTypeHDFS         RisingWaveStateStoreBackendType = "HDFS"
	RisingWaveStateStoreBackendTypeWebHDFS      RisingWaveStateStoreBackendType = "WebHDFS"
	RisingWaveStateStoreBackendTypeGCS          RisingWaveStateStoreBackendType = "GCS"
)

// RisingWaveS3Credentials is the reference and keys selector to the AWS access credentials stored in a local secret.
type RisingWaveS3Credentials struct {
	// The name of the secret in the pod's namespace to select from.
	SecretName string `json:"secretName"`

	// AccessKeyRef is the key of the secret to be the access key. Must be a valid secret key. Defaults to "AccessKey".
	// +kubebuilder:default=AccessKey
	AccessKeyRef *string `json:"accessKeyRef,omitempty"`

	// SecretAccessKeyRef is the key of the secret to be the secret access key. Must be a valid secret key.
	// Defaults to "SecretAccessKey".
	// +kubebuilder:default=SecretAccessKey
	SecretAccessKeyRef *string `json:"secretAccessKeyRef,omitempty"`
}

// RisingWaveStateStoreBackendS3 is the collection of parameters for the S3 backend state store.
type RisingWaveStateStoreBackendS3 struct {
	// RisingWaveS3Credentials is the credentials provider from a Secret.
	RisingWaveS3Credentials `json:"credentials"`

	// Region is the region of the S3 bucket. Defaults to us-east-1.
	// +kubebuilder:default=us-east-1
	Region string `json:"region,omitempty"`

	// Bucket is the name of the S3 bucket.
	Bucket string `json:"bucket"`

	// DataDirectory is the prefix of the objects stored. Leave it to empty if you would like to override this option
	// with external configuration source like ConfigMap or Secret.
	DataDirectory string `json:"dataDirectory,omitempty"`
}

// RisingWaveStateStoreBackendS3C is the collection of parameters for the S3-compatible backend state store.
type RisingWaveStateStoreBackendS3C struct {
	// RisingWaveS3Credentials is the credentials provider from a Secret.
	RisingWaveS3Credentials `json:"credentials"`

	// Endpoint is the endpoint of the S3-compatible service. It could start with the "http://" or the "https://"
	// scheme prefix. But if no scheme is found, the default scheme is HTTPS.
	// Besides, One should decide if the endpoint should be within a virtual-hosted style. For more information about
	// the virtual host mode, please refer to https://docs.aws.amazon.com/AmazonS3/latest/userguide/VirtualHosting.html.
	Endpoint string `json:"endpoint"`

	// Region is the region of the S3-compatible bucket. Defaults to empty.
	// +optional
	Region string `json:"region,omitempty"`

	// Bucket is the name of the S3-compatible bucket.
	Bucket string `json:"bucket"`

	// DataDirectory is the prefix of the objects stored. Leave it to empty if you would like to override this option
	// with external configuration source like ConfigMap or Secret.
	DataDirectory string `json:"dataDirectory,omitempty"`
}

// RisingWaveMinIOCredentials is the reference and keys selector to the MinIO access credentials stored in a local secret.
type RisingWaveMinIOCredentials struct {
	// The name of the secret in the pod's namespace to select from.
	SecretName string `json:"secretName"`

	// UsernameKeyRef is the key of the secret to be the username. Must be a valid secret key. Defaults to "Username".
	// +kubebuilder:default=Username
	UsernameKeyRef *string `json:"usernameKeyRef,omitempty"`

	// PasswordKeyRef is the key of the secret to be the password. Must be a valid secret key.
	// Defaults to "Password".
	// +kubebuilder:default=Username
	PasswordKeyRef *string `json:"passwordKeyRef,omitempty"`
}

// RisingWaveStateStoreBackendMinIO is the collection of parameters for the MinIO backend state store.
type RisingWaveStateStoreBackendMinIO struct {
	// RisingWaveMinIOCredentials is the credentials provider from a Secret.
	RisingWaveMinIOCredentials `json:"credentials"`

	// Endpoint is the endpoint of the MinIO service. It should not contain any scheme prefix.
	Endpoint string `json:"endpoint"`

	// Bucket is the name of the MinIO bucket.
	Bucket string `json:"bucket"`
}

// RisingWaveStateStoreBackendHDFS is the collection of parameters for the HDFS backend state store.
type RisingWaveStateStoreBackendHDFS struct {
	// NameNode of the HDFS service.
	// +kubebuilder:validation:Required
	NameNode string `json:"nameNode,omitempty"`

	// Root of the working directory on the HDFS.
	// +kubebuilder:validation:Required
	Root string `json:"root,omitempty"`
}

// RisingWaveGCSCredentials is the reference and keys selector to the GCS access credentials stored in a local secret.
type RisingWaveGCSCredentials struct {
	// UseWorkloadIdentity indicates to use workload identity to access the GCS service. If this is enabled, secret is not required, and ADC is used.
	UseWorkloadIdentity bool `json:"useWorkloadIdentity,omitempty"`

	// The name of the secret in the pod's namespace to select from.
	// +optional
	SecretName string `json:"secretName,omitempty"`

	// ServiceAccountCredentialsKeyRef is the key of the secret to be the service account credentials. Must be a valid secret key. Defaults to "Username".
	// +kubebuilder:default=ServiceAccountCredentials
	// +optional
	ServiceAccountCredentialsKeyRef *string `json:"serviceAccountCredentialsKeyRef,omitempty"`
}

// RisingWaveStateStoreBackendGCS is the collection of parameters for the GCS backend state store.
type RisingWaveStateStoreBackendGCS struct {
	// RisingWaveGCSCredentials is the credentials provider from a Secret.
	RisingWaveGCSCredentials `json:"credentials,omitempty"`

	// Bucket of the GCS bucket service.
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`

	// Working directory root of the GCS bucket
	// +kubebuilder:validation:Required
	Root string `json:"root"`
}

// RisingWaveStateStoreBackend is the collection of parameters for the state store that RisingWave uses. Note that one
// and only one of the first-level fields could be set.
type RisingWaveStateStoreBackend struct {
	// Memory determines whether RisingWave uses a memory-based state store. Keep in mind that the memory
	// backend is only for test purposes and should not be used in production. Defaults to false.
	Memory *bool `json:"memory,omitempty"`

	// MinIO determines whether RisingWave uses a MinIO-backed state store and the parameters for accessing the MinIO.
	MinIO *RisingWaveStateStoreBackendMinIO `json:"minio,omitempty"`

	// S3 determines whether RisingWave uses a S3-backed state store and the parameters for accessing the S3.
	S3 *RisingWaveStateStoreBackendS3 `json:"s3,omitempty"`

	// S3C determines whether RisingWave uses a S3-compatible service-backed state store and the parameters for
	// accessing the S3-compatible service.
	S3C *RisingWaveStateStoreBackendS3C `json:"s3c,omitempty"`

	// HDFS determines whether RisingWave uses a HDFS service-backed state store and the parameters for
	// accessing the HDFS service.
	HDFS *RisingWaveStateStoreBackendHDFS `json:"hdfs,omitempty"`

	// GCS determines whether RisingWave uses a GCS service-backed state store and the parameters for
	// accessing the GCS service.
	GCS *RisingWaveStateStoreBackendGCS `json:"gcs,omitempty"`
}

// RisingWaveStateStoreStatus is the status of the state store.
type RisingWaveStateStoreStatus struct {
	// Backend type of the state store.
	Backend RisingWaveStateStoreBackendType `json:"backend,omitempty"`
}
