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

package v1alpha1

import (
	"strings"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/risingwavelabs/risingwave-operator/apis/risingwave/v1alpha2"
)

const (
	aliyunOSSEndpoint         = "${BUCKET}.oss-${REGION}.aliyuncs.com"
	internalAliyunOSSEndpoint = "${BUCKET}.oss-${REGION}-internal.aliyuncs.com"
)

var _ conversion.Convertible = (*RisingWave)(nil)

func buildV1alpha2RisingWaveNodePodTemplate(src *RisingWaveComponentGroupTemplate) v1alpha2.RisingWaveNodePodTemplate {
	if src == nil {
		return v1alpha2.RisingWaveNodePodTemplate{}
	}

	var dst v1alpha2.RisingWaveNodePodTemplate

	// Convert the image.
	dst.Image = src.Image

	// Convert the image pull policy.
	dst.ImagePullPolicy = src.ImagePullPolicy

	// Convert the ImagePullSecrets.
	dst.ImagePullSecrets = lo.Map(src.ImagePullSecrets, func(s string, _ int) corev1.LocalObjectReference {
		return corev1.LocalObjectReference{Name: s}
	})

	// Convert the resources.
	dst.Resources = src.Resources

	// Convert the node selector.
	dst.NodeSelector = src.NodeSelector

	// Convert the tolerations.
	dst.Tolerations = src.Tolerations

	// Convert the affinity.
	dst.Affinity = src.Affinity

	// Convert the pod annotations.
	dst.Annotations = src.Metadata.Annotations

	// Convert the pod labels.
	dst.Labels = src.Metadata.Labels

	// Convert the pod security context.
	dst.SecurityContext = src.SecurityContext

	// Convert the pod image pull secrets.
	dst.ImagePullSecrets = lo.Map(src.ImagePullSecrets, func(s string, _ int) corev1.LocalObjectReference {
		return corev1.LocalObjectReference{Name: s}
	})

	// Convert the pod priority class name.
	dst.PriorityClassName = src.PriorityClassName

	// Convert the pod termination grace period.
	dst.TerminationGracePeriodSeconds = src.TerminationGracePeriodSeconds

	// Convert the pod DNS config.
	dst.DNSConfig = src.DNSConfig

	// Convert the envs.
	dst.Env = src.Env

	// Convert the env from.
	dst.EnvFrom = src.EnvFrom

	return dst
}

func buildV1alpha2RisingWaveNodeGroup(src *RisingWaveComponentGroup) v1alpha2.RisingWaveNodeGroup {
	var dst v1alpha2.RisingWaveNodeGroup

	// Convert the name.
	dst.Name = src.Name

	// Convert the replicas.
	dst.Replicas = src.Replicas

	// Convert the upgrade strategy.
	dst.UpgradeStrategy = v1alpha2.RisingWaveNodeGroupUpgradeStrategy{
		Type: v1alpha2.RisingWaveNodeGroupUpgradeStrategyType(src.UpgradeStrategy.Type),
		RollingUpdate: lo.If(src.UpgradeStrategy.RollingUpdate != nil, (*v1alpha2.RisingWaveNodeGroupRollingUpdate)(nil)).
			Else(&v1alpha2.RisingWaveNodeGroupRollingUpdate{
				MaxUnavailable: src.UpgradeStrategy.RollingUpdate.MaxUnavailable,
				Partition:      src.UpgradeStrategy.RollingUpdate.Partition,
				MaxSurge:       src.UpgradeStrategy.RollingUpdate.MaxSurge,
			}),
		InPlaceUpdateStrategy: src.UpgradeStrategy.InPlaceUpdateStrategy,
	}

	// Convert the pod template.
	dst.Template = buildV1alpha2RisingWaveNodePodTemplate(src.RisingWaveComponentGroupTemplate)

	return dst
}

func buildV1alpha2RisingWaveComponent(replicasOfDefaultGroup int32, restartAt *metav1.Time, groups []RisingWaveComponentGroup) v1alpha2.RisingWaveComponent {
	var dst v1alpha2.RisingWaveComponent

	// Default log level to INFO.
	dst.LogLevel = "INFO"

	// Convert the node groups.
	dst.NodeGroups = lo.Map(groups, func(group RisingWaveComponentGroup, _ int) v1alpha2.RisingWaveNodeGroup {
		ng := buildV1alpha2RisingWaveNodeGroup(&group)
		ng.RestartAt = restartAt
		return ng
	})

	// Default group.
	dst.NodeGroups = append(dst.NodeGroups, v1alpha2.RisingWaveNodeGroup{
		Name:      "",
		Replicas:  replicasOfDefaultGroup,
		RestartAt: restartAt,
	})

	return dst
}

func buildV1alpha2RisingWaveComponentStatus(src ComponentReplicasStatus) v1alpha2.RisingWaveComponentStatus {
	var dst v1alpha2.RisingWaveComponentStatus

	// Convert the total replicas.
	dst.Total = v1alpha2.WorkloadReplicaStatus{
		Replicas:            src.Target,
		ReadyReplicas:       src.Running,
		AvailableReplicas:   src.Running,
		UpdatedReplicas:     src.Running,
		UnavailableReplicas: 0,
	}

	// Convert the node groups.
	dst.NodeGroups = lo.Map(src.Groups, func(group ComponentGroupReplicasStatus, _ int) v1alpha2.RisingWaveNodeGroupStatus {
		return v1alpha2.RisingWaveNodeGroupStatus{
			Name: group.Name,
			WorkloadReplicaStatus: v1alpha2.WorkloadReplicaStatus{
				Replicas:            group.Target,
				ReadyReplicas:       group.Running,
				AvailableReplicas:   group.Running,
				UpdatedReplicas:     group.Running,
				UnavailableReplicas: 0,
			},
			Exists: group.Exists,
		}
	})

	return dst
}

// ConvertTo converts this RisingWave to the Hub version (v1alpha2).
//
//goland:noinspection GoReceiverNames
func (src *RisingWave) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.RisingWave)

	dst.ObjectMeta = src.ObjectMeta

	// Convert the spec.
	dst.Spec.UseKruiseWorkloads = src.Spec.EnableOpenKruise
	dst.Spec.SyncPrometheusServiceMonitor = src.Spec.EnableDefaultServiceMonitor
	dst.Spec.FrontendServiceType = src.Spec.Global.ServiceType
	dst.Spec.AdditionalFrontendServiceMetadata.Labels = src.Spec.Global.ServiceMeta.Labels
	dst.Spec.AdditionalFrontendServiceMetadata.Annotations = src.Spec.Global.ServiceMeta.Annotations
	dst.Spec.MetaStore = func() v1alpha2.RisingWaveMetaStoreBackend {
		srcMetaStore := src.Spec.Storages.Meta
		metaStore := v1alpha2.RisingWaveMetaStoreBackend{}
		metaStore.Memory = srcMetaStore.Memory
		if src.Spec.Storages.Meta.Etcd != nil {
			metaStore.Etcd = &v1alpha2.RisingWaveMetaStoreBackendEtcd{
				RisingWaveEtcdCredentials: lo.If(srcMetaStore.Etcd.Secret == "", (*v1alpha2.RisingWaveEtcdCredentials)(nil)).
					Else(&v1alpha2.RisingWaveEtcdCredentials{
						SecretName:     srcMetaStore.Etcd.Secret,
						UsernameKeyRef: pointer.String("username"),
						PasswordKeyRef: pointer.String("password"),
					}),
				Endpoints: srcMetaStore.Etcd.Endpoint,
			}
		}
		return metaStore
	}()
	dst.Spec.StateStore = func() v1alpha2.RisingWaveStateStoreBackend {
		srcStateStore := src.Spec.Storages.Object
		stateStore := v1alpha2.RisingWaveStateStoreBackend{}
		stateStore.Memory = srcStateStore.Memory

		if srcStateStore.MinIO != nil {
			stateStore.MinIO = &v1alpha2.RisingWaveStateStoreBackendMinIO{
				RisingWaveMinIOCredentials: v1alpha2.RisingWaveMinIOCredentials{
					SecretName:     srcStateStore.MinIO.Secret,
					UsernameKeyRef: pointer.String("username"),
					PasswordKeyRef: pointer.String("password"),
				},
				Endpoint: srcStateStore.MinIO.Endpoint,
				Bucket:   srcStateStore.MinIO.Bucket,
			}
		}

		if srcStateStore.S3 != nil {
			if srcStateStore.S3.Endpoint == "" {
				stateStore.S3 = &v1alpha2.RisingWaveStateStoreBackendS3{
					RisingWaveS3Credentials: v1alpha2.RisingWaveS3Credentials{
						SecretName:         srcStateStore.S3.Secret,
						AccessKeyRef:       pointer.String("AccessKeyID"),
						SecretAccessKeyRef: pointer.String("SecretAccessKey"),
					},
					Region: stateStore.S3.Region,
					Bucket: stateStore.S3.Bucket,
				}
			} else {
				endpoint := srcStateStore.S3.Endpoint
				if srcStateStore.S3.VirtualHostedStyle {
					if strings.HasPrefix(endpoint, "https://") {
						endpoint = "https://${BUCKET}" + endpoint[len("https://"):]
					} else {
						endpoint = "${BUCKET}" + endpoint
					}
				}

				stateStore.S3C = &v1alpha2.RisingWaveStateStoreBackendS3C{
					RisingWaveS3Credentials: v1alpha2.RisingWaveS3Credentials{
						SecretName:         srcStateStore.S3.Secret,
						AccessKeyRef:       pointer.String("AccessKeyID"),
						SecretAccessKeyRef: pointer.String("SecretAccessKey"),
					},
					Endpoint: endpoint,
					Region:   srcStateStore.S3.Region,
					Bucket:   srcStateStore.S3.Bucket,
				}
			}
		}

		if srcStateStore.AliyunOSS != nil {
			stateStore.S3C = &v1alpha2.RisingWaveStateStoreBackendS3C{
				RisingWaveS3Credentials: v1alpha2.RisingWaveS3Credentials{
					SecretName:         srcStateStore.AliyunOSS.Secret,
					AccessKeyRef:       pointer.String("AccessKeyID"),
					SecretAccessKeyRef: pointer.String("SecretAccessKey"),
				},
				Endpoint: lo.If(srcStateStore.AliyunOSS.InternalEndpoint, internalAliyunOSSEndpoint).Else(aliyunOSSEndpoint),
				Region:   srcStateStore.AliyunOSS.Region,
				Bucket:   srcStateStore.AliyunOSS.Bucket,
			}
		}

		if srcStateStore.HDFS != nil {
			stateStore.HDFS = &v1alpha2.RisingWaveStateStoreBackendHDFS{
				NameNode: srcStateStore.HDFS.NameNode,
				Root:     srcStateStore.HDFS.Root,
			}
		}

		if srcStateStore.GCS != nil {
			stateStore.GCS = &v1alpha2.RisingWaveStateStoreBackendGCS{
				RisingWaveGCSCredentials: v1alpha2.RisingWaveGCSCredentials{
					UseWorkloadIdentity:             srcStateStore.GCS.UseWorkloadIdentity,
					SecretName:                      srcStateStore.GCS.Secret,
					ServiceAccountCredentialsKeyRef: pointer.String("ServiceAccountCredentials"),
				},
				Bucket: srcStateStore.GCS.Bucket,
				Root:   srcStateStore.GCS.Root,
			}
		}

		return stateStore
	}()
	dst.Spec.Image = src.Spec.Global.Image
	dst.Spec.PodTemplate = buildV1alpha2RisingWaveNodePodTemplate(&src.Spec.Global.RisingWaveComponentGroupTemplate)
	dst.Spec.Configuration = v1alpha2.RisingWaveNodeConfiguration{
		ConfigMap: lo.If(src.Spec.Configuration.ConfigMap == nil, (*v1alpha2.RisingWaveNodeConfigurationConfigMapSource)(nil)).
			Else(&v1alpha2.RisingWaveNodeConfigurationConfigMapSource{
				Name:     src.Spec.Configuration.ConfigMap.Name,
				Key:      src.Spec.Configuration.ConfigMap.Key,
				Optional: src.Spec.Configuration.ConfigMap.Optional,
			}),
	}
	dst.Spec.MetaComponent = buildV1alpha2RisingWaveComponent(src.Spec.Global.Replicas.Meta,
		src.Spec.Components.Meta.RestartAt, src.Spec.Components.Meta.Groups)
	dst.Spec.ComputeComponent = buildV1alpha2RisingWaveComponent(src.Spec.Global.Replicas.Compute,
		src.Spec.Components.Compute.RestartAt, lo.Map(src.Spec.Components.Compute.Groups, func(g RisingWaveComputeGroup, _ int) RisingWaveComponentGroup {
			return RisingWaveComponentGroup{
				Name:                             g.Name,
				Replicas:                         g.Replicas,
				RisingWaveComponentGroupTemplate: lo.If(g.RisingWaveComputeGroupTemplate == nil, (*RisingWaveComponentGroupTemplate)(nil)).Else(&g.RisingWaveComponentGroupTemplate),
			}
		}))
	// Build a map of volumeMounts for each node group.
	volumeMounts := make(map[string][]corev1.VolumeMount)
	for _, ng := range src.Spec.Components.Compute.Groups {
		volumeMounts[ng.Name] = ng.VolumeMounts
	}
	for _, ng := range dst.Spec.ComputeComponent.NodeGroups {
		ng.VolumeClaimTemplates = lo.Map(src.Spec.Storages.PVCTemplates, func(src PersistentVolumeClaim, _ int) v1alpha2.PersistentVolumeClaim {
			return v1alpha2.PersistentVolumeClaim{
				PersistentVolumeClaimPartialObjectMeta: v1alpha2.PersistentVolumeClaimPartialObjectMeta{
					Name:        src.PersistentVolumeClaimPartialObjectMeta.Name,
					Labels:      src.PersistentVolumeClaimPartialObjectMeta.Labels,
					Annotations: src.PersistentVolumeClaimPartialObjectMeta.Annotations,
					Finalizers:  src.PersistentVolumeClaimPartialObjectMeta.Finalizers,
				},
				Spec: src.Spec,
			}
		})
		ng.Template.VolumeMounts = volumeMounts[ng.Name]
	}
	dst.Spec.FrontendComponent = buildV1alpha2RisingWaveComponent(src.Spec.Global.Replicas.Frontend,
		src.Spec.Components.Frontend.RestartAt, src.Spec.Components.Frontend.Groups)
	dst.Spec.CompactorComponent = buildV1alpha2RisingWaveComponent(src.Spec.Global.Replicas.Compactor,
		src.Spec.Components.Compactor.RestartAt, src.Spec.Components.Compactor.Groups)

	// Convert the status.
	dst.Status.ObservedGeneration = src.Status.ObservedGeneration
	dst.Status.Conditions = func() []v1alpha2.RisingWaveCondition {
		conditions := make([]v1alpha2.RisingWaveCondition, len(src.Status.Conditions))
		for i, cond := range src.Status.Conditions {
			conditions[i] = v1alpha2.RisingWaveCondition{
				Type:               v1alpha2.RisingWaveConditionType(cond.Type),
				Status:             corev1.ConditionStatus(cond.Status),
				LastTransitionTime: cond.LastTransitionTime,
				Reason:             cond.Reason,
				Message:            cond.Message,
			}
		}
		return conditions
	}()
	dst.Status.ImageTag = src.Status.Version
	dst.Status.MetaStore = v1alpha2.RisingWaveMetaStoreStatus{
		Backend: lo.Switch[MetaStorageType, v1alpha2.RisingWaveMetaStoreBackendType](src.Status.Storages.Meta.Type).
			Case(MetaStorageTypeEtcd, v1alpha2.RisingWaveMetaStoreBackendTypeEtcd).
			Case(MetaStorageTypeMemory, v1alpha2.RisingWaveMetaStoreBackendTypeMemory).
			Default(""),
	}
	dst.Status.StateStore = v1alpha2.RisingWaveStateStoreStatus{
		Backend: lo.Switch[ObjectStorageType, v1alpha2.RisingWaveStateStoreBackendType](src.Status.Storages.Object.Type).
			Case(ObjectStorageTypeMemory, v1alpha2.RisingWaveStateStoreBackendTypeMemory).
			Case(ObjectStorageTypeMinIO, v1alpha2.RisingWaveStateStoreBackendTypeMinIO).
			Case(ObjectStorageTypeS3, v1alpha2.RisingWaveStateStoreBackendTypeS3).
			Case(ObjectStorageTypeGCS, v1alpha2.RisingWaveStateStoreBackendTypeGCS).
			Case(ObjectStorageTypeAliyunOSS, v1alpha2.RisingWaveStateStoreBackendTypeS3Compatible).
			Case(ObjectStorageTypeHDFS, v1alpha2.RisingWaveStateStoreBackendTypeHDFS).
			Default(""),
	}
	dst.Status.MetaComponent = buildV1alpha2RisingWaveComponentStatus(src.Status.ComponentReplicas.Meta)
	dst.Status.ComputeComponent = buildV1alpha2RisingWaveComponentStatus(src.Status.ComponentReplicas.Compute)
	dst.Status.CompactorComponent = buildV1alpha2RisingWaveComponentStatus(src.Status.ComponentReplicas.Compactor)
	dst.Status.FrontendComponent = buildV1alpha2RisingWaveComponentStatus(src.Status.ComponentReplicas.Frontend)
	dst.Status.ScaleViewLocks = lo.Map(src.Status.ScaleViews, func(l RisingWaveScaleViewLock, _ int) v1alpha2.RisingWaveScaleViewLock {
		return v1alpha2.RisingWaveScaleViewLock{
			Reference: v1alpha2.RisingWaveScaleViewReference{
				Name:               l.Name,
				UID:                l.UID,
				ObservedGeneration: l.Generation,
			},
			Component: l.Component,
			Locks: lo.Map(l.GroupLocks, func(gl RisingWaveScaleViewLockGroupLock, _ int) v1alpha2.RisingWaveScaleViewNodeGroupLock {
				return v1alpha2.RisingWaveScaleViewNodeGroupLock{
					Name:     gl.Name,
					Replicas: gl.Replicas,
				}
			}),
		}
	})

	return nil
}

// ConvertFrom converts this RisingWave from the Hub version (v1alpha2).
//
//goland:noinspection GoReceiverNames
func (dst *RisingWave) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.RisingWave)

	dst.ObjectMeta = src.ObjectMeta

	// Convert the spec.

	// Convert the status.

	return nil
}
