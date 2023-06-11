/*
Copyright 2023 xihua.

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
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	testappv1 "webserver-operator/api/v1"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=testapp.github.com,resources=redis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=testapp.github.com,resources=redis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=testapp.github.com,resources=redis/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Redis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	redis := &testappv1.Redis{}
	// 获取历史数据
	if err := r.Client.Get(ctx, req.NamespacedName, redis); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if redis.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	deployment := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deployment); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		// 1. 不存在，则创建
		// 1-1. 创建 Deployment
		deployment = newDeployment(redis)
		if err := r.Client.Create(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
		fmt.Println("r.Client.Create Deployment success")

		// 1-2. 创建 Service
		svc := newService(redis)
		if err := r.Client.Create(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
		fmt.Println("r.Client.Create Service success")
		// 学习 Finalizers，OwnerReference
		/*redis.Finalizers = append(redis.Finalizers, deployment.Name)
		if err := r.Client.Update(ctx, redis); err != nil {
			return ctrl.Result{}, err
		}*/
	} else {
		oldSpec := &testappv1.Redis{}
		if err := json.Unmarshal([]byte(redis.Annotations["spec"]), oldSpec); err != nil {
			return ctrl.Result{}, err
		}
		// 2. 对比更新
		if !reflect.DeepEqual(redis.Spec, *oldSpec) {
			// 2-1. 更新 Deployment 资源
			newDeployment := newDeployment(redis)
			currDeployment := &appsv1.Deployment{}
			if err := r.Client.Get(ctx, req.NamespacedName, currDeployment); err != nil {
				return ctrl.Result{}, err
			}
			currDeployment.Spec = newDeployment.Spec
			if err := r.Client.Update(ctx, currDeployment); err != nil {
				return ctrl.Result{}, err
			}
			fmt.Println("r.Client.Create Update Deployment success")

			// 2-2. 更新 Service 资源
			newService := newService(redis)
			currService := &corev1.Service{}
			if err := r.Client.Get(ctx, req.NamespacedName, currService); err != nil {
				return ctrl.Result{}, err
			}

			currIP := currService.Spec.ClusterIP
			currService.Spec = newService.Spec
			currService.Spec.ClusterIP = currIP
			if err := r.Client.Update(ctx, currService); err != nil {
				return ctrl.Result{}, err
			}
			fmt.Println("r.Client.Create Update Service success")
		}
	}

	// 3. 关联 Annotations
	data, _ := json.Marshal(redis.Spec)
	if redis.Annotations != nil {
		redis.Annotations["spec"] = string(data)
	} else {
		redis.Annotations = map[string]string{"spec": string(data)}
	}
	if err := r.Client.Update(ctx, redis); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func newContainer(app *testappv1.Redis) []corev1.Container {
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range app.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	return []corev1.Container{
		{
			Name:            app.Name,
			Image:           app.Spec.Image,
			Resources:       app.Spec.Resources,
			Ports:           containerPorts,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             app.Spec.Envs,
		},
	}
}

// 获取一个新的deployment
func newDeployment(app *testappv1.Redis) *appsv1.Deployment {
	labels := map[string]string{"app": app.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   testappv1.GroupVersion.Group,
					Version: testappv1.GroupVersion.Version,
					Kind:    app.Kind,
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec:       corev1.PodSpec{Containers: newContainer(app)},
			},
		},
	}
}

// 获取一个新的service
func newService(app *testappv1.Redis) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   testappv1.GroupVersion.Group,
					Version: testappv1.GroupVersion.Version,
					Kind:    app.Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			// Ports: app.Spec.Ports,
			Ports: []corev1.ServicePort{{
				Name:     "http",
				Port:     app.Spec.Ports[0].TargetPort.IntVal,
				NodePort: app.Spec.Ports[0].NodePort,
			}},
			Selector: map[string]string{
				"app": app.Name,
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testappv1.Redis{}).
		Complete(r)
}
