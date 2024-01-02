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
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeProgressing string = "Progressing"
	ConditionTypeReconcile   string = "Reconcile"
	ConditionTypeAvailable   string = "Available"

	ConditionReasonPreparing           string = "Preparing"
	ConditionReasonRunning             string = "Running"
	ConditionReasonConfig              string = "Config"
	ConditionReasonReconcilePVC        string = "ReconcilePVC"
	ConditionReasonReconcileService    string = "ReconcileService"
	ConditionReasonReconcileIngress    string = "ReconcileIngress"
	ConditionReasonReconcileDeployment string = "ReconcileDeployment"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MongoDBSpec defines the desired state of MongoDB
type MongoDBSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MongoDB. Edit mongodb_types.go to remove/update
	ReplicaCount int32  `json:"replicaCount"`
	StorageSize  string `json:"storageSize"`
}

type Node struct {
	HostName string `json:"hostName,omitempty"`
	IP       string `json:"ip,omitempty"`
	Status   string `json:"status,omitempty"`
}

// MongoDBStatus defines the observed state of MongoDB
type MongoDBStatus struct {
	// Phase string `json:"phase"`
	Conditions []metav1.Condition `json:"condition"`
	Nodes      []Node             `json:"nodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MongoDB is the Schema for the mongodbs API
type MongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoDBSpec   `json:"spec,omitempty"`
	Status MongoDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MongoDBList contains a list of MongoDB
type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDB{}, &MongoDBList{})
}
func (m *MongoDB) SetStatusCondition(condition metav1.Condition) {
	// if the condition already exists, update it
	existingCondition := apimeta.FindStatusCondition(m.Status.Conditions, condition.Type)
	if existingCondition == nil {
		condition.ObservedGeneration = m.GetGeneration()
		condition.LastTransitionTime = metav1.Now()
		m.Status.Conditions = append(m.Status.Conditions, condition)
	} else if existingCondition.Status != condition.Status || existingCondition.Reason != condition.Reason || existingCondition.Message != condition.Message {
		existingCondition.Status = condition.Status
		existingCondition.Reason = condition.Reason
		existingCondition.Message = condition.Message
		existingCondition.ObservedGeneration = m.GetGeneration()
		existingCondition.LastTransitionTime = metav1.Now()
	}
}

func (m *MongoDB) InitStatusCondtions() {
	m.Status.Conditions = []metav1.Condition{}
	m.SetStatusCondition(metav1.Condition{
		Type:               ConditionTypeProgressing,
		Status:             metav1.ConditionTrue,
		ObservedGeneration: m.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonPreparing,
		Message:            "Mongodb is preparing",
	})
	m.SetStatusCondition(metav1.Condition{
		Type:               ConditionTypeAvailable,
		Status:             metav1.ConditionFalse,
		ObservedGeneration: m.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonPreparing,
		Message:            "Mongodb is preparing",
	})
}
