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
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

const (
	StateReady      string = "Ready"
	StateError      string = "Error"
	StateProcessing string = "Processing"
	StateWarning    string = "Warning"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CompanionSpec defines the desired state of Companion.
type CompanionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Namespaces          []string `json:"namespaces"`
	DeploymentNamespace string   `json:"deploymentNamespace"`
}

// CompanionStatus defines the observed state of Companion.
type CompanionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Defines the overall state of the Companion custom resource.<br/>
	// - `Ready` when all the resources managed by the Kyma companion manager are deployed successfully and
	// the companion backend is ready.<br/>
	// - `Warning` if there is a user input misconfiguration.<br/>
	// - `Processing` if the resources managed by the Kyma companion manager are being created or updated.<br/>
	// - `Error` if an error occurred while reconciling the Companion custom resource.
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
