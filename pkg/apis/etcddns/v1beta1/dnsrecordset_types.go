/*
Copyright 2018 The Kubernetes authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSRecordSetSpec represents a set of DNS records that must be created
// and managed on a DNS provider
type DNSRecordSetSpec struct {
	// Endpoints holds a list of DNS records to be synchronized with the
	// cloud provider by the external-dns controller
	// +optional
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
}

// DNSRecordSetStatus defines the observed state of DNSRecordSet
type DNSRecordSetStatus struct {
	// ObservedGeneration is the generation observed by the external-dns controller
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// TTL is a value defining the TTL of a DNS record
type TTL int64

// Targets is a representation of a list of targets for an endpoint
type Targets []string

// ProviderSpecific holds configuration which is specific to individual DNS providers
type ProviderSpecific map[string]string

// Labels store metadata related to the endpoint
// it is then stored in a persistent storage via serialization
type Labels map[string]string

// Endpoint represents a DNS record in a DNSRecordSet
type Endpoint struct {
	// The hostname of the DNS record
	DNSName string `json:"dnsName,omitempty"`
	// The targets the DNS record points to
	Targets Targets `json:"targets,omitempty"`
	// RecordType type of record, e.g. CNAME, A, SRV, TXT etc
	RecordType string `json:"recordType,omitempty"`
	// TTL for the record
	RecordTTL TTL `json:"recordTTL,omitempty"`
	// Labels stores labels defined for the Endpoint
	// +optional
	Labels Labels `json:"labels,omitempty"`
	// ProviderSpecific stores provider specific config
	// +optional
	ProviderSpecific ProviderSpecific `json:"providerSpecific,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSRecordSet is the Schema for the dnsrecordsets API
// +k8s:openapi-gen=true
// +kubebuilder:resource
// +kubebuilder:subresource:status
type DNSRecordSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSRecordSetSpec   `json:"spec,omitempty"`
	Status DNSRecordSetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSRecordSetList contains a list of DNSRecordSet
// +kubebuilder:resource
type DNSRecordSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSRecordSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSRecordSet{}, &DNSRecordSetList{})
}
