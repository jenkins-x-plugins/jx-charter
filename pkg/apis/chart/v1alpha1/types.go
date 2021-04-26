package v1alpha1

import (
	"helm.sh/helm/v3/pkg/chart"
	rspb "helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=ch

// Chart contains the definition of a preview environment
type Chart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChartSpec    `json:"spec,omitempty"`
	Status *ChartStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// ChartList represents a list of pipeline options
type ChartList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Chart `json:"items"`
}

// +k8s:deepcopy-gen=false

// ChartSpec contains the chart metadata
type ChartSpec struct {
	chart.Metadata
}

// ChartStatus contains the chart status information
type ChartStatus struct {
	// FirstDeployed is when the release was first deployed.
	FirstDeployed metav1.Time `json:"firstDeployed,omitempty"`
	// LastDeployed is when the release was last deployed.
	LastDeployed metav1.Time `json:"lastDeployed,omitempty"`
	// Deleted tracks when this object was deleted.
	Deleted *metav1.Time `json:"deleted,omitempty"`
	// Description is human-friendly "log entry" about this release.
	Description string `json:"description,omitempty"`
	// Status is the current state of the release
	Status string `json:"status,omitempty"`
	// Contains the rendered templates/NOTES.txt if available
	Notes string `json:"notes,omitempty"`
}

// DeepCopy a custom deep copy function
func (in *ChartSpec) DeepCopy() *ChartSpec {
	out := *in
	return &out
}

// DeepCopyInto creates
func (in *ChartSpec) DeepCopyInto(c *ChartSpec) {
	*c = *in
}

// ToChartStatus converts the release info to a chart status
func ToChartStatus(info *rspb.Info) *ChartStatus {
	if info == nil {
		return nil
	}
	answer := &ChartStatus{
		FirstDeployed: metav1.Time{Time: info.FirstDeployed.Time},
		LastDeployed:  metav1.Time{Time: info.LastDeployed.Time},
		Description:   info.Description,
		Status:        string(info.Status),
		Notes:         info.Notes,
	}

	if !info.Deleted.IsZero() {
		t := metav1.Time{Time: info.Deleted.Time}
		answer.Deleted = &t
	}
	return answer
}
