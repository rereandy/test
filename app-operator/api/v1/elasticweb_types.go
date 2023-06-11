/*
Copyright 2023.

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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ElasticWebSpec defines the desired state of ElasticWeb
type ElasticWebSpec struct {
	Image string `json:"image"`
	Port  *int32 `json:"port"`
	// 单个pod的QPS上限
	SinglePodsQPS *int32 `json:"singlePodsQPS"`
	// 当前整个业务的QPS
	TotalQPS *int32 `json:"totalQPS,omitempty"`
}

type ElasticWebStatus struct {
	// 当前 Kubernetes 集群实际支持的总QPS
	RealQPS *int32 `json:"realQPS"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ElasticWeb is the Schema for the elasticwebs API
type ElasticWeb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticWebSpec   `json:"spec,omitempty"`
	Status ElasticWebStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ElasticWebList contains a list of ElasticWeb
type ElasticWebList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticWeb `json:"items"`
}

func (e *ElasticWeb) String() string {
	var realQPS string
	if nil == e.Status.RealQPS {
		realQPS = ""
	} else {
		realQPS = strconv.Itoa(int(*e.Status.RealQPS))
	}

	return fmt.Sprintf("Image [%s], Port [%d], SinglePodQPS [%d], TotalQPS [%d], RealQPS [%s]",
		e.Spec.Image,
		*e.Spec.Port,
		*e.Spec.SinglePodsQPS,
		*e.Spec.TotalQPS,
		realQPS)
}

func init() {
	SchemeBuilder.Register(&ElasticWeb{}, &ElasticWebList{})
}
