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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSEndpointSpec defines the desired state of DNSEndpoint
type DNSEndpointSpec struct {
	// Endpoints is the list of DNS records to create/update
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
}

// DNSEndpointStatus defines the observed state of DNSEndpoint
type DNSEndpointStatus struct {
	// ObservedGeneration is the generation observed by the external-dns controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// TTL is the time to live of DNS records
type TTL int64

// Targets is the set of values associated with a DNS record
type Targets []string

// Labels is a set of labels associated with a DNS record
type EndpointLabels map[string]string

// ProviderSpecific contains cloud provider specific configuration for a DNS record
type ProviderSpecific map[string]string

// Endpoint represents a single DNS record
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
	Labels EndpointLabels `json:"labels,omitempty"`
	// ProviderSpecific stores provider specific config
	// +optional
	ProviderSpecific ProviderSpecific `json:"providerSpecific,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSEndpoint is the Schema for the dnsendpoints API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type DNSEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSEndpointSpec   `json:"spec,omitempty"`
	Status DNSEndpointStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSEndpointList contains a list of DNSEndpoint
type DNSEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSEndpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSEndpoint{}, &DNSEndpointList{})
}
