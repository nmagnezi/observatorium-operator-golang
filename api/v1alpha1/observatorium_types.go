/*

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ObservatoriumSpec defines the desired state of Observatorium
type ObservatoriumSpec struct {
	// Thanos Spec
	Thanos ThanosSpec `json:"thanos"`
}

type ThanosReceiveController struct {
	// Thanos receive controller Image name
	Image *string `json:"image"`
	// Tag describes the tag of Thanos receive controller to use.
	Tag *string `json:"tag,omitempty"`
	// Hashrings describes a list of Hashrings
	Hashrings []*Hashring `json:"hashrings,omitempty"`
}

type Hashring struct {
	// Thanos Hashring name
	Name *string `json:"name"`
	// Tenants describes a lists of tenants.
	Tenants []*string `json:"tenants,omitempty"`
	//
}

type ThanosSpec struct {
	// Thanos Image name
	Image *string `json:"image"`
	// Tag of Thanos sidecar container image to be deployed.
	Tag *string `json:"tag"`

	ThanosReceiveControllerSpec ThanosReceiveController `json:"thanosReceiveControllerSpec"`
	// Number of instances to deploy for a Thanos querier.
	QuerierReplicas *int32 `json:"querierReplicas,omitempty"`
	// Resources for Querier pods
	QuerierResources v1.ResourceRequirements `json:"querierResources,omitempty"`
	// Number of instances to deploy for a Thanos Store.
	StoreReplicas *int32 `json:"storeReplicas,omitempty"`
	// Resources for Store pods
	StoreResources v1.ResourceRequirements `json:"storeResources,omitempty"`
	// Number of instances to deploy for a Thanos Compactor.
	CompactorReplicas *int32 `json:"compactorReplicas,omitempty"`
	// Resources for Compactor pods
	CompactorResources v1.ResourceRequirements `json:"compactorResources,omitempty"`
	// Number of instances to deploy for a Thanos Receive.
	ReceiveReplicas *int32 `json:"receiveReplicas,omitempty"`
	// Resources for Receive pods
	ReceiveResources v1.ResourceRequirements `json:"receiveResources,omitempty"`
	// Receive Storage Class
	ReceiveStorageClass *string `json:"receiveStorageClass"`
	// Receive PVC size
	ReceivePVCSize *string `json:"receivePvcSize"`
	// Object Store Config Secret for Thanos
	ObjectStoreConfigSecret *string `json:"objectStoreConfigSecret"`
	// TODO: AWS secrets?
	// TODO: handle with THANOS_QUERIER_SVC_URL
	// TODO: Do we need a THANOS_RULER?
	// TODO: JAEGER
}

// ObservatoriumStatus defines the observed state of Observatorium
type ObservatoriumStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Observatorium is the Schema for the observatoria API
type Observatorium struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObservatoriumSpec   `json:"spec,omitempty"`
	Status ObservatoriumStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ObservatoriumList contains a list of Observatorium
type ObservatoriumList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Observatorium `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Observatorium{}, &ObservatoriumList{})
}
