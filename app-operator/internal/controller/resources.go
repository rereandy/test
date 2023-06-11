// controllers/resource.go

package controller

import (
	v1 "app-operator/api/v1"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	ElasticWebCommonLabelKey = "app"
)

const (
	// APP_NAME deployment 中 App 标签名
	APP_NAME = "elastic-app"
	// CONTAINER_PORT 容器的端口号
	// 这里有个坑点，原本想着可以自定义端口，但是最后发现访问不了
	// 所以这句可有可无，Nginx的默认http端口不就80么(狗头保命)
	CONTAINER_PORT = 80
	// CPU_REQUEST 单个POD的CPU资源申请
	CPU_REQUEST = "100m"
	// CPU_LIMIT 单个POD的CPU资源上限
	CPU_LIMIT = "100m"
	// MEM_REQUEST 单个POD的内存资源申请
	MEM_REQUEST = "512Mi"
	// MEM_LIMIT 单个POD的内存资源上限
	MEM_LIMIT = "512Mi"
)

// 根据总QPS以及单个POD的QPS，计算需要多少个Pod
func getExpectReplicas(elasticWeb *v1.ElasticWeb) int32 {
	// 单个pod的QPS
	singlePodQPS := *elasticWeb.Spec.SinglePodsQPS
	// 期望的总QPS
	totalQPS := *elasticWeb.Spec.TotalQPS
	// 需要创建的副本数
	replicas := totalQPS / singlePodQPS

	if totalQPS%singlePodQPS != 0 {
		replicas += 1
	}
	return replicas
}

// CreateServiceIfNotExists  创建service
func CreateServiceIfNotExists(ctx context.Context, r *ElasticWebReconciler, elasticWeb *v1.ElasticWeb, req ctrl.Request) error {
	logger := log.FromContext(ctx)
	logger.WithValues("func", "createService")
	svc := &corev1.Service{}

	svc.Name = elasticWeb.Name
	svc.Namespace = elasticWeb.Namespace

	svc.Spec = corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:     "http",
				Port:     CONTAINER_PORT,
				NodePort: *elasticWeb.Spec.Port,
			},
		},
		Type: corev1.ServiceTypeNodePort,
		Selector: map[string]string{
			ElasticWebCommonLabelKey: APP_NAME,
		},
	}

	// 设置关联关系
	logger.Info("set reference")
	if err := controllerutil.SetControllerReference(elasticWeb, svc, r.Scheme); err != nil {
		logger.Error(err, "SetControllerReference error")
		return err
	}

	logger.Info("start create service")
	if err := r.Create(ctx, svc); err != nil {
		logger.Error(err, "create service error")
		return err
	}

	return nil
}

// CreateDeployment 创建deployment
func CreateDeployment(ctx context.Context, r *ElasticWebReconciler, elasticWeb *v1.ElasticWeb) error {
	logger := log.FromContext(ctx)
	logger.WithValues("func", "createDeploy")

	// 计算期待pod的数量
	expectReplicas := getExpectReplicas(elasticWeb)
	logger.Info(fmt.Sprintf("expectReplicas [%d]", expectReplicas))

	deploy := &appsv1.Deployment{}

	deploy.Labels = map[string]string{
		ElasticWebCommonLabelKey: APP_NAME,
	}

	deploy.Name = elasticWeb.Name
	deploy.Namespace = elasticWeb.Namespace

	deploy.Spec = appsv1.DeploymentSpec{
		Replicas: pointer.Int32Ptr(expectReplicas),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				ElasticWebCommonLabelKey: APP_NAME,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					ElasticWebCommonLabelKey: APP_NAME,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  APP_NAME,
						Image: elasticWeb.Spec.Image,
						Ports: []corev1.ContainerPort{
							{
								Name:          "http",
								ContainerPort: CONTAINER_PORT,
								Protocol:      corev1.ProtocolSCTP,
							},
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(CPU_LIMIT),
								corev1.ResourceMemory: resource.MustParse(MEM_LIMIT),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(CPU_REQUEST),
								corev1.ResourceMemory: resource.MustParse(MEM_REQUEST),
							},
						},
					},
				},
			},
		},
	}

	// 建立关联,删除web后会将deploy一起删除
	logger.Info("set reference")
	if err := controllerutil.SetControllerReference(elasticWeb, deploy, r.Scheme); err != nil {
		logger.Error(err, "SetControllerReference error")
		return err
	}

	// 创建Deployment
	logger.Info("start create deploy")
	if err := r.Create(ctx, deploy); err != nil {
		logger.Error(err, "create deploy error")
		return err
	}

	logger.Info("create deploy success")
	return nil
}

func UpdateStatus(ctx context.Context, r *ElasticWebReconciler, elasticWeb *v1.ElasticWeb) error {
	logger := log.FromContext(ctx)
	logger.WithValues("func", "updateStatus")

	// 单个pod的QPS
	singlePodQPS := *elasticWeb.Spec.SinglePodsQPS

	// pod 总数
	replicas := getExpectReplicas(elasticWeb)

	// 当pod创建完成后，当前系统的QPS为： 单个pod的QPS * pod总数
	// 如果没有初始化，则需要先初始化
	if nil == elasticWeb.Status.RealQPS {
		elasticWeb.Status.RealQPS = new(int32)
	}

	*elasticWeb.Status.RealQPS = singlePodQPS * replicas
	logger.Info(fmt.Sprintf("singlePodQPS [%d],replicas [%d],realQPS[%d]", singlePodQPS, replicas, *elasticWeb.Status.RealQPS))

	if err := r.Update(ctx, elasticWeb); err != nil {
		logger.Error(err, "update instance error")
		return err
	}
	return nil
}
