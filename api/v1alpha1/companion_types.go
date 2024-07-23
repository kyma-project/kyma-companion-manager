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

// CompanionSpec defines the desired state of Companion.
type CompanionSpec struct {

	// Companion configuration.
	// +kubebuilder:default:={namespace:"companion",deploymentNamespace:"ai-core",configMapNames:{"companion-config"},secretsNames:{"companion-secrets"},replicas:1}
	CompanionConfig CompanionConfig `json:"companionConfig"`

	// HANA Cloud configuration.
	// +kubebuilder:default:={namespace:"hana-cloud"}
	HanaCloudConfig HanaCloudConfig `json:"hanaCloudConfig"`

	// Redis configuration.
	// +kubebuilder:default:={namespace:"redis"}
	RedisConfig RedisConfig `json:"redisConfig"`

	// Container port for the companion backend. Default value is 5000.
	// ContainerPort int32 `json:"containerPort"`

	// Annotations allows to add annotations to NATS.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels allows to add Labels to NATS.
	Labels map[string]string `json:"labels,omitempty"`
}

// CompanionConfig defines the configuration for the Companion.
type CompanionConfig struct {
	// Companion namespace where the companion backend will be deployed
	// and the related configMaps and secrets are already stored.
	// +kubebuilder:default:="ai-core"
	Namespace string `json:"namespace"`
	// ConfigMap names for the companion backend.
	ConfigMapNames []string `json:"configMapNames"`
	// Secret names for the companion backend.
	SecretsNames []string `json:"secretsNames"`

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

// hanaCloudConfig defines the configuration for the HANA Cloud.
type HanaCloudConfig struct {
	// HANA Cloud namespace where the HANA Cloud is deployed.
	// +kubebuilder:default:="hana-cloud"
	Namespace string `json:"namespace"`
}

// redisConfig defines the configuration for the Redis.
type RedisConfig struct {
	// Redis namespace where the Redis is deployed.
	// +kubebuilder:default:="redis"
	Namespace string `json:"namespace"`
}

// CompanionStatus defines the observed state of Companion.
type CompanionStatus struct {
	// Result of prerequisites validation.
	// NamespacesExist: Map of namespaces and their existence status.
	NamespacesExist map[string]bool `json:"namespacesExist"`
	// ConfigMapsExists: Map of ConfigMaps and their existence status.
	ConfigMapsExists map[string]bool `json:"configMapsExists"`
	// SecretsExists: Map of Secrets and their existence status.
	SecretsExists map[string]bool `json:"secretsExists"`
	// ConfigMapsData: Map of ConfigMaps and their data. (optional)
	ConfigMapsData map[string]map[string]string `json:"configMapsData,omitempty"`
	// SecretsData: Map of Secrets and their data. (optional)
	SecretsData map[string]map[string][]byte `json:"secretsData,omitempty"`

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
