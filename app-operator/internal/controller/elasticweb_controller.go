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
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	elasticwebv1 "app-operator/api/v1"
)

// ElasticWebReconciler reconciles a ElasticWeb object
type ElasticWebReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=elasticweb.example.com,resources=elasticwebs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=elasticweb.example.com,resources=elasticwebs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=elasticweb.example.com,resources=elasticwebs/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *ElasticWebReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 相关接口都是传入这个参数当作 object
	instance := &elasticwebv1.ElasticWeb{}

	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info(fmt.Sprintf("instance:%s", instance.String()))

	// 获取deployment，pod删除的时候，相关的deployment,service自动删除
	deploy := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deploy); err != nil {
		if errors.IsNotFound(err) {
			// 如果没有查找到，则需要创建
			logger.Info("deploy not exists")
			// 判断qps的需求，如果qps没有需求，则啥都不做
			if *instance.Spec.TotalQPS < 1 {
				logger.Info("not need deployment")
				return ctrl.Result{}, nil
			}

			// 创建Deploy
			if err := CreateDeployment(ctx, r, instance); err != nil {
				return ctrl.Result{}, err
			}

			// 创建service
			if err = CreateServiceIfNotExists(ctx, r, instance, req); err != nil {
				return ctrl.Result{}, err
			}

			// 更新状态
			if err := UpdateStatus(ctx, r, instance); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get deploy")
		return ctrl.Result{}, err
	}

	// 根据单个Pod的QPS计算期望pod的副本
	expectReplicas := getExpectReplicas(instance)

	// 获取当前deployment实际的pod副本
	realReplicas := deploy.Spec.Replicas

	if expectReplicas == *realReplicas {
		logger.Info("not need to reconcile")
		return ctrl.Result{}, nil
	}

	// 重新赋值
	deploy.Spec.Replicas = &expectReplicas
	// 更新 deploy
	if err := r.Update(ctx, deploy); err != nil {
		logger.Error(err, "update deploy replicas error")
		return ctrl.Result{}, err
	}

	// 更新状态
	if err := UpdateStatus(ctx, r, instance); err != nil {
		logger.Error(err, "update status error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ElasticWebReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&elasticwebv1.ElasticWeb{}).
		Complete(r)
}
