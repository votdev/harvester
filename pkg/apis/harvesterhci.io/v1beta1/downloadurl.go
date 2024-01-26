package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	INIT_STATE        = "INIT"
	DOWNLOADING_STATE = "DOWNLOADING"
	DONE_STATE        = "DONE"
)

type DownloadURLStatus struct {
	// +kubebuilder:default:=INIT
	Status string `json:"status,omitempty" default:"INIT"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:shortName=durl;durls,scope=Namespaced
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=`.url`
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.status`
// +kubebuilder:subresource:status

type DownloadURL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	URL    string            `json:"url,omitempty"`
	Status DownloadURLStatus `json:"status,omitempty"`
}
