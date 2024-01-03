package controller

import (
	mongodbv1 "github.com/lwpk110/mongo-operator/api/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *MongoDBReconciler) createStatefulSet(instance *mongodbv1.MongoDB) *appv1.StatefulSet {
	pvcName := instance.Name + "pvc"
	obj := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &instance.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: instance.Name + "pod-",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: instance.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "tcp",
									ContainerPort: 27017,
								},
							},
							Command: []string{
								"mongod",
								"--replSete",
								instance.Spec.ReplicasName,
								"--bind_ip",
								"0.0.0.0",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      pvcName,
									MountPath: "/data/db",
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: pvcName,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.PersistentVolumeAccessMode(instance.Spec.Persistence.AccessMode),
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse(instance.Spec.Persistence.Size),
							},
						},
						StorageClassName: &instance.Spec.Persistence.StorageClass,
					},
				},
			},
		},
		Status: appv1.StatefulSetStatus{},
	}
	err := controllerutil.SetControllerReference(instance, obj, r.Scheme)
	if err != nil {
		r.Log.Error(err, "StatefulSet set to controller reference error, name: %s", instance.Name)
		return nil //todo
	}
	return obj
}

func (r *MongoDBReconciler) ensureStatefulSet(instance *mongodbv1.MongoDB) {

}
