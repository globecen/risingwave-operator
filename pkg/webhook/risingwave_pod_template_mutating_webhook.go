/*
 * Copyright 2023 RisingWave Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package webhook

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	risingwavev1alpha1 "github.com/risingwavelabs/risingwave-operator/apis/risingwave/v1alpha1"
	"github.com/risingwavelabs/risingwave-operator/pkg/metrics"
)

// RisingWavePodTemplateMutatingWebhook is the mutating webhook for RisingWavePodTemplate.
type RisingWavePodTemplateMutatingWebhook struct{}

func (pm *RisingWavePodTemplateMutatingWebhook) setDefault(ctx context.Context, obj *risingwavev1alpha1.RisingWavePodTemplate) error {
	return nil
}

// Default implements webhook.CustomDefaulter.
func (pm *RisingWavePodTemplateMutatingWebhook) Default(ctx context.Context, obj runtime.Object) error {
	return pm.setDefault(ctx, obj.(*risingwavev1alpha1.RisingWavePodTemplate))
}

// NewRisingWavePodTemplateMutatingWebhook returns a new mutating webhook for RisingWavePodTemplate.
func NewRisingWavePodTemplateMutatingWebhook() webhook.CustomDefaulter {
	return metrics.NewMutatingWebhookMetricsRecorder(&RisingWavePodTemplateMutatingWebhook{})
}
