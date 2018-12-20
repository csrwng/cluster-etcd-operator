/*
Copyright 2018 The Kubernetes Authors.

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

package v1

import (
	operatorsv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EtcdConfigSpec defines the desired state of EtcdConfig
type EtcdConfigSpec struct {
	// managementState indicates whether and how the operator should manage the component
	ManagementState operatorsv1.ManagementState `json:"managementState"`

	// dnsConfig specifies configuration for etcd DNS entry management
	DNSConfig DNSConfig `json:"dnsConfig"`
}

type DNSConfig struct {
	// logLevel is the level of logging for the external-dns controller
	// Valid values: debug, info, warning, error, fatal
	LogLevel string `json:"logLevel,omitempty"`

	// automaticUpdates indicates that the DNS entries should be automatically
	// updated based on the IP address of master Machine resources
	AutomaticUpdates bool `json:"automaticUpdates"`
}

// EtcdConfigStatus defines the observed state of EtcdConfig
type EtcdConfigStatus struct {
	operatorsv1.OperatorStatus `json:",inline"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdConfig is the Schema for the etcdconfigs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +genclient:nonNamespaced
type EtcdConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdConfigSpec   `json:"spec,omitempty"`
	Status EtcdConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdConfigList contains a list of EtcdConfig
type EtcdConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdConfig{}, &EtcdConfigList{})
}
