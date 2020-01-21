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

type ReceiveController struct {
	Replicas *int32 `json:"replicas,omitempty"`
	// Thanos receive controller Image name
	Image *string `json:"image"`
	// Tag describes the tag of Thanos receive controller to use.
	Tag *string `json:"tag,omitempty"`
	// Hashrings describes a list of Hashrings
	Hashrings []*Hashring `json:"hashrings,omitempty"`
	// Resources for component pods
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

type ThanosPersistentSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
	// Resources for component pods
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Receive Storage Class
	StorageClass *string `json:"storageClass"`
	// Receive PVC size
	PVCSize *string `json:"pvcSize"`
}

type QuerierCacheSpec struct {
	// Thanos receive controller Image name
	Image *string `json:"image"`
	// ConfigMap describes the Configuration of Querier Cache.
	ConfigMap *string `json:"config-map"`
	// Number of Querier Cache replicas.
	Replicas *int32 `json:"replicas,omitempty"`
	// Resources for Receive pods
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Receive Storage Class
}

type ThanosComponentSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
	// Resources for component pods
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

type Hashring struct {
	// Thanos Hashring name
	Name *string `json:"name"`
	// Tenants describes a lists of tenants.
	Tenants []*string `json:"tenants,omitempty"`
}

type ThanosSpec struct {
	// Thanos Image name
	Image *string `json:"image"`
	// Tag of Thanos sidecar container image to be deployed.
	Tag *string `json:"tag,omitempty"`
	// Thanos Receive Controller Spec
	ReceiveControllerSpec ReceiveController `json:"receiveController"`
	// Thanos ThanosPersistentSpec
	Receive ThanosPersistentSpec `json:"receive"`
	// Thanos QuerierSpec
	Querier ThanosComponentSpec `json:"querier"`
	// Thanos QuerierCache
	QuerierCache QuerierCacheSpec `json:"querier-cache"`
	// Thanos StoreSpec
	Store ThanosPersistentSpec `json:"store"`
	// Thanos CompactorSpec
	Compactor ThanosComponentSpec `json:"compactor"`
	// Thanos RulerSpec
	Ruler ThanosComponentSpec `json:"ruler"`
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
