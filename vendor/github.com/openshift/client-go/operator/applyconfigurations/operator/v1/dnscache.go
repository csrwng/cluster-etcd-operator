// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSCacheApplyConfiguration represents a declarative configuration of the DNSCache type for use
// with apply.
type DNSCacheApplyConfiguration struct {
	PositiveTTL *metav1.Duration `json:"positiveTTL,omitempty"`
	NegativeTTL *metav1.Duration `json:"negativeTTL,omitempty"`
}

// DNSCacheApplyConfiguration constructs a declarative configuration of the DNSCache type for use with
// apply.
func DNSCache() *DNSCacheApplyConfiguration {
	return &DNSCacheApplyConfiguration{}
}

// WithPositiveTTL sets the PositiveTTL field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PositiveTTL field is set to the value of the last call.
func (b *DNSCacheApplyConfiguration) WithPositiveTTL(value metav1.Duration) *DNSCacheApplyConfiguration {
	b.PositiveTTL = &value
	return b
}

// WithNegativeTTL sets the NegativeTTL field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NegativeTTL field is set to the value of the last call.
func (b *DNSCacheApplyConfiguration) WithNegativeTTL(value metav1.Duration) *DNSCacheApplyConfiguration {
	b.NegativeTTL = &value
	return b
}
