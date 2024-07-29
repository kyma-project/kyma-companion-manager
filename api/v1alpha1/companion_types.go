/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StateReady      string = "Ready"
	StateError      string = "Error"
	StateProcessing string = "Processing"
	StateDeleting   string = "Deleting"
	StateWarning    string = "Warning"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecretSpec defines the secret name and namespace.
type SecretSpec struct {
	// Secret name and namespace for the secret.
	// Name: Name of the secret.
	// Namespace: Namespace of the secret.
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// CompanionSpec defines the desired state of Companion.
type CompanionSpec struct {

	// AI Core configuration
	// +kubebuilder:default:={aicore:{secret: "ai-core/aicore"}}
	AICore AICoreConfig `json:"aicore"`

	// HANA Cloud configuration
	// +kubebuilder:default:={hanaCloud:{secret: "companion/hana-cloud"}}
	HanaCloud HanaConfig `json:"hanaCloud"`

	// Redis configuration
	// +kubebuilder:default:={redis:{secret: "companion/redis"}}
	Redis RedisConfig `json:"redis"`

	// Companion configuration.
	// +kubebuilder:default:={replicas:{min:1, max:3}, resources:{limits:{cpu:"4",memory:"4Gi"}, requests:{cpu:"500m",memory:"256Mi"}}}
	Companion CompanionConfig `json:"companion"`
}

// AICoreConfig defines the configuration for the AI Core.
type AICoreConfig struct {
	// Secret name and namespace for the AI Core.
	// +kubebuilder:default:= {secret:{"name": "aicore", "namespace": "ai-core"}}
	Secret SecretSpec `json:"secret"`
}

// HanaConfig defines the configuration for the Hana Cloud.
type HanaConfig struct {
	// Secret name and namespace for the Han Cloud.
	// +kubebuilder:default:= {secret:{"name": "companion", "namespace": "hana-cloud"}}
	Secret SecretSpec `json:"secret"`
}

// RedisConfig defines the configuration for the Redis.
type RedisConfig struct {
	// Secret name and namespace for the Redis.
	// +kubebuilder:default:= {secret:{"name": "companion", "namespace": "redis"}}
	Secret SecretSpec `json:"secret"`
}

// CompanionConfig defines the configuration for the Companion.
type CompanionConfig struct {
	// Secret name and namespace for the companion backend.
	// +kubebuilder:default:= {secret:{"name": "companion", "namespace": "ai-core"}}
	Secret SecretSpec `json:"secret"`
	// Number of replicas for the companion backend.
	// +kubebuilder:default:={min:1, max:1}
	Replicas ReplicasConfig `json:"replicas"`

	// Specify required resources and resource limits for the companion backend.
	// Example:
	// resources:
	//   limits:
	//     cpu: 1
	//     memory: 1Gi
	//   requests:
	//     cpu: 500m
	//     memory: 256Mi
	// +kubebuilder:default:={resources:{limits:{cpu:"4",memory:"4Gi"},requests:{cpu:"500m",memory:"256Mi"}}}
	// +kubebuilder:default:={limits:{cpu:"4",memory:"4Gi"}, requests:{cpu:"500m",memory:"256Mi"}}
	Resources kcorev1.ResourceRequirements `json:"resources,omitempty"`
}

// ReplicasConfig defines the min and max replicas
type ReplicasConfig struct {
	// +kubebuilder:validation:Minimum=1
	Min int `json:"min"`

	// +kubebuilder:validation:Minimum=1
	Max int `json:"max"`
}

// CompanionStatus defines the observed state of Companion.
type CompanionStatus struct {
	// Defines the overall state of the Companion custom resource.<br/>
	// - `Ready` when all the resources managed by the Kyma companion manager are deployed successfully and
	// the companion backend is ready.<br/>
	// - `Warning` if there is a user input misconfiguration.<br/>
	// - `Processing` if the resources managed by the Kyma companion manager are being created or updated.<br/>
	// - `Error` if an error occurred while reconciling the Companion custom resource.
	// - `Deleting` if the resources managed by the Kyma companion manager are being deleted.
	State string `json:"state"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Companion is the Schema for the companions API.
type Companion struct {
	kmetav1.TypeMeta   `json:",inline"`
	kmetav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompanionSpec   `json:"spec,omitempty"`
	Status CompanionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CompanionList contains a list of Companion.
type CompanionList struct {
	kmetav1.TypeMeta `json:",inline"`
	kmetav1.ListMeta `json:"metadata,omitempty"`
	Items            []Companion `json:"items"`
}

//nolint:gochecknoinits // scaffolded by kubebuilder.
func init() {
	SchemeBuilder.Register(&Companion{}, &CompanionList{})
}
