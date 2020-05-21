package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AceApplicationSpec defines the desired state of AceApplication
type AceApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Replicas *int32            `json:"replicas"`
	BarLocation string      `json:"barLocation"`
	License string          `json:"license"`
	AceServerName string    `json:"aceServerName"`
	AceTruststorePwd string `json:"aceTruststorePwd"`
	AceKeystorePwd string   `json:"aceKeystorePwd"`
	LogFormat string        `json:"logFormat"`
	EnableMetrics string    `json:"enableMetrics"`
	AceBaseImage string     `json:"aceBaseImage"`
	AceImageTag string      `json:"aceImageTag"`
	CpuRequest string		`json:"cpuRequest"`
	CpuLimit string			`json:"cpuLimit"`
	MemoryRequest string	`json:"memoryRequest"`
	MemoryLimit string		`json:"memoryLimit"`
	ServerConfig string		`json:"serverConfig"`
	SetDBParms string		`json:"setDBParms"`
	Keystore string 		`json:"keystore"`
	Truststore string 		`json:"truststore"`
	Policy string			`json:"policy"`
	Webusers string 		`json:"webusers"`
	NodeSelectorLabels map[string]string		`json:"nodeSelectorLabels"`
	ServiceAccountName string `json:"serviceAccountName"`
}

// AceApplicationStatus defines the observed state of AceApplication
type AceApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AceApplication is the Schema for the aceapplications API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=aceapplications,scope=Namespaced
type AceApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AceApplicationSpec   `json:"spec,omitempty"`
	Status AceApplicationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AceApplicationList contains a list of AceApplication
type AceApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AceApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AceApplication{}, &AceApplicationList{})
}
