package controller

import (
	mongodbv1 "github.com/lwpk110/mongo-operator/api/v1"
)

func (r *MongoDBReconciler) createStatefulset(instance *mongodbv1.MongoDB) {
	// appv1.StatefulSet{
	// 	TypeMeta:   metav1.TypeMeta{},
	// 	ObjectMeta: metav1.ObjectMeta{},
	// 	Spec:       appv1.StatefulSetSpec{},
	// 	Status:     appv1.StatefulSetStatus{},
	// }
}

func (r *MongoDBReconciler) ensureStatefulset(instance *mongodbv1.MongoDB) {

}
