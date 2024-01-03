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

// MongoDBSpec defines the desired state of MongoDB
type MongoDBSpec struct {
	Image string `json:"image"`
	// +kubebuilder:validation:Optionl
	// +kubebuilder:default:1
	Replicas     int32            `json:"replicas"`
	ReplicasName string           `json:"replicasName"`
	Persistence  *PersistenceSpec `json:"persistence"`
}

type PersistenceSpec struct {
	// +kubebuilder:validation:Optional
	StorageClass string `json:"storageClass"`
	Size         string `json:"size"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:ReadWriteOnce
	AccessMode string `json:"accessMode"`
}

// MongoDBStatus defines the observed state of MongoDB
type MongoDBStatus struct {
	// Phase string `json:"phase"`
	Conditions []metav1.Condition `json:"condition"`
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

func (m *MongoDB) InitStatusConditions() {
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
