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

package controller

import (
	"context"

	"github.com/go-logr/logr"
	mongodbv1 "github.com/lwpk110/mongo-operator/api/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// MongoDBReconciler reconciles a MongoDB object
type MongoDBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

func (r *MongoDBReconciler) UpdateStatus(ctx context.Context, instance *mongodbv1.MongoDB) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return r.Status().Update(ctx, instance)
		//return r.Status().Patch(ctx, instance, client.MergeFrom(instance))
	})

	if retryErr != nil {
		r.Log.Error(retryErr, "Failed to update vfm status after retries")
		return retryErr
	}

	r.Log.V(1).Info("Successfully patched object status")
	return nil
}

//+kubebuilder:rbac:groups=mongodb.steven.com,resources=mongodbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mongodb.steven.com,resources=mongodbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mongodb.steven.com,resources=mongodbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MongoDB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MongoDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	instance := &mongodbv1.MongoDB{}
	r.Log.Info("MongoDB reconcile starting...")
	// 抓取资源对象，判断资源是否存在
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if client.IgnoreNotFound(err) != nil { // 错误不为空，检查是否为NotFound错误
			r.Log.Error(err, "Get mongodb resource err")
			return ctrl.Result{}, err // 非 NotFound 异常，返回实际异常
		}
		r.Log.Info("MongoDB resource NotFound, may be deleted")
		return ctrl.Result{}, nil // NotFound 异常，返回空结果
	}
	// TODO(user): your logic here
	// 获取资源状态，如果状态不对，初始化，更新
	readCondition := apimeta.FindStatusCondition(instance.Status.Conditions, mongodbv1.ConditionTypeProgressing)
	if readCondition == nil || readCondition.ObservedGeneration != instance.GetGeneration() {
		instance.InitStatusCondtions()
		if err := r.UpdateStatus(ctx, instance); err != nil {
			r.Log.Error(err, "Update resource status error when init status stage")
			return ctrl.Result{}, err
		}
	}

	r.Log.Info("MongoDB resource found:", instance.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mongodbv1.MongoDB{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}
