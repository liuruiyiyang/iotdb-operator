package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IoTDBSpec defines the desired state of IoTDB
// +k8s:openapi-gen=true
type IoTDBSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html


	Sharding []TimeSegment `json:"sharding"`
	// ImagePullPolicy defines how the image is pulled.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy"`
	// IotdbImage is the image to use for the IoTDB Pods.
	IotdbImage string `json:"iotdbImage"`
	// StorageMode can be EmptyDir, HostPath, NFS
	StorageMode string `json:"storageMode"`
	// HostPath is the local path to store data
	HostPath string `json:"hostPath"`
	// VolumeClaimTemplates defines the StorageClass
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates"`
	// Replica number of each read and write instance
	ReplicaNum int `json:"replicaNum"`
	// IkrImage is the image to use for the IKR Pods.
	IkrImage string `json:"ikrImage"`
	// IkrSize is the size of the IKR instance
	IkrSize int `json:"ikrSize"`

}

// +k8s:openapi-gen=true
type TimeSegment struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Time string of time segment start time
	ScaleTime string `json:"scaleTime"`
	// The schema segment instance number in this time segment
	SchemaSegmentNum int `json:"schemaSegmentNum"`
}

// IoTDBStatus defines the observed state of IoTDB
// +k8s:openapi-gen=true
type IoTDBStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of the IoTDB pods
	Nodes []string `json:"nodes"`
	TimeSegmentNum int `json:"timeSegmentNum"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IoTDB is the Schema for the iotdbs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type IoTDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IoTDBSpec   `json:"spec,omitempty"`
	Status IoTDBStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IoTDBList contains a list of IoTDB
type IoTDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IoTDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IoTDB{}, &IoTDBList{})
}
