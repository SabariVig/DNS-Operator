package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusType string
type Provider string
type RecordType string

const (
	// StatusType
	Active       StatusType = "Active"
	Provisioning StatusType = "Provisioning"
	Failed       StatusType = "Failed"

	// Provider
	Namecheap Provider = "Namecheap"

	// RecordType
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	MX    RecordType = "MX"
	CNAME RecordType = "CNAME"
)

// RecordSpec defines the desired state of Record
type RecordSpec struct {
	Provider Provider   `json:"provider"`
	Domain   string     `json:"domain"`
	Hostname string     `json:"hostname"`
	Address  string     `json:"address"`
	Type     RecordType `json:"type"`
	TTL      int        `json:"ttl,omitempty"`
}

// RecordStatus defines the observed state of Record
type RecordStatus struct {
	LastModified metav1.Time `json:"last_modified,omitempty"`
	Status       StatusType  `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Last Modified",type="date",JSONPath=`.status.last_modified`
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=`.status.status`

// Record is the Schema for the records API
type Record struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordSpec   `json:"spec,omitempty"`
	Status RecordStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordList contains a list of Record
type RecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Record `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Record{}, &RecordList{})
}
