package k8sinterface

import (
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

//K8sInterface defines methods to interface with Kubernetes
type K8sInterface struct {
}

//JobInfo contains parameters to run a Kubernetes Job
type JobInfo struct {
	JobName             string
	Namespace           string
	Image               string
	KlessRepo           string
	EventHandlerName    string
	EventHandlerVersion string
	EventHandlerSource  string
	InterfaceSource     string
	EntrypointSource    string
	ContextSource       string
	RequestSource       string
	ResponseSource      string
}

//DeploymentInfo contains parameters to create a Kubernetes Deployment
type DeploymentInfo struct {
	Namespace               string
	Name                    string
	Version                 string
	Replicas                int32
	EventHandlerName        string
	EventHandlerImage       string
	EventHandlerPort        int32
	EventHandlerCPULimit    string
	EventHandlerMemoryLimit string
	FrontendName            string
	FrontendImage           string
	FrontendPort            int32
	FrontendCPULimit        string
	FrontendMemoryLimit     string
	FrontendEnvironmentVars map[string]string
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", 0)
}

//CreateNamespace creates a Kubernetes namespace
func (k *K8sInterface) CreateNamespace(namespace string) error {
	namespaceSpec := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespaces := c.Core().Namespaces()
	_, err = namespaces.Get(namespace, metav1.GetOptions{})
	switch {
	case err == nil:
		logger.Println("namespace exists")
	case errors.IsNotFound(err):
		_, err = namespaces.Create(namespaceSpec)
		if err != nil {
			return fmt.Errorf("failed to create namespace: %s", err)
		}
		logger.Println("namespace created")
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	return nil
}

//CreateDeployment creates a deployment
func (k *K8sInterface) CreateDeployment(d *DeploymentInfo) error {
	frontendEnvironment := []apiv1.EnvVar{
		apiv1.EnvVar{
			Name:  "KLESS_EVENT_HANDLER_NAME",
			Value: d.Name,
		},
		apiv1.EnvVar{
			Name:  "KLESS_EVENT_HANDLER_NAMESPACE",
			Value: d.Namespace,
		},
		apiv1.EnvVar{
			Name:  "KLESS_EVENT_HANDLER_VERSION",
			Value: d.Version,
		},
	}
	if nil != d.FrontendEnvironmentVars {
		for k, v := range d.FrontendEnvironmentVars {
			frontendEnvironment = append(frontendEnvironment, apiv1.EnvVar{Name: k, Value: v})
		}
	}

	klessID := uuid.NewV4().String()

	deploySpec := &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
			Labels:    map[string]string{"app": d.Name, "k8s-type": "kless-handler", "kless-id": klessID},
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: int32p(d.Replicas),
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1beta1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			RevisionHistoryLimit: int32p(10),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   d.Name,
					Labels: map[string]string{"app": d.Name, "k8s-type": "kless-handler", "kless-id": klessID},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:  d.EventHandlerName,
							Image: d.EventHandlerImage,
							Ports: []apiv1.ContainerPort{
								apiv1.ContainerPort{ContainerPort: d.EventHandlerPort, Protocol: apiv1.ProtocolTCP},
							},
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "KLESS_EVENT_HANDLER_NAME",
									Value: d.Name,
								},
								apiv1.EnvVar{
									Name:  "KLESS_EVENT_HANDLER_NAMESPACE",
									Value: d.Namespace,
								},
								apiv1.EnvVar{
									Name:  "KLESS_EVENT_HANDLER_VERSION",
									Value: d.Version,
								},
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse(d.EventHandlerCPULimit),
									apiv1.ResourceMemory: resource.MustParse(d.EventHandlerMemoryLimit),
								},
							},
							ImagePullPolicy: apiv1.PullIfNotPresent,
						},
						apiv1.Container{
							Name:  d.FrontendName,
							Image: d.FrontendImage,
							Ports: []apiv1.ContainerPort{
								apiv1.ContainerPort{ContainerPort: d.FrontendPort, Protocol: apiv1.ProtocolTCP},
							},
							Env: frontendEnvironment,
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse(d.FrontendCPULimit),
									apiv1.ResourceMemory: resource.MustParse(d.FrontendMemoryLimit),
								},
							},
							ImagePullPolicy: apiv1.PullAlways,
						},
					},
					RestartPolicy: apiv1.RestartPolicyAlways,
					DNSPolicy:     apiv1.DNSClusterFirst,
				},
			},
		},
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploy := c.Extensions().Deployments(d.Namespace)

	_, err = deploy.Update(deploySpec)
	switch {
	case err == nil:
		logger.Println("deployment controller updated")
	case !errors.IsNotFound(err):
		return fmt.Errorf("could not update deployment controller: %s", err)
	default:
		_, err = deploy.Create(deploySpec)
		if err != nil {
			return fmt.Errorf("could not create deployment controller: %s", err)
		}
		logger.Println("deployment controller created")
	}

	return nil
}

//CreateJob creates a job
func (k *K8sInterface) CreateJob(j *JobInfo) error {
	jobSpec := &batch.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.JobName,
			Namespace: j.Namespace,
		},
		Spec: batch.JobSpec{
			ActiveDeadlineSeconds: int64p(300),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: j.JobName,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:  j.JobName,
							Image: j.Image,
							VolumeMounts: []apiv1.VolumeMount{
								apiv1.VolumeMount{
									Name:      "dockersocket",
									MountPath: "/var/run/docker.sock",
								},
							},
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "KLESS_REPO",
									Value: j.KlessRepo,
								},
								apiv1.EnvVar{
									Name:  "KLESS_EVENT_HANDLER_NAME",
									Value: j.EventHandlerName,
								},
								apiv1.EnvVar{
									Name:  "KLESS_EVENT_HANDLER_VERSION",
									Value: j.EventHandlerVersion,
								},
								apiv1.EnvVar{
									Name:  "KLESS_EVENT_HANDLER_SOURCE",
									Value: j.EventHandlerSource,
								},
								apiv1.EnvVar{
									Name:  "KLESS_INTERFACE_SOURCE",
									Value: j.InterfaceSource,
								},
								apiv1.EnvVar{
									Name:  "KLESS_ENTRYPOINT_SOURCE",
									Value: j.EntrypointSource,
								},
								apiv1.EnvVar{
									Name:  "KLESS_CONTEXT_SOURCE",
									Value: j.ContextSource,
								},
								apiv1.EnvVar{
									Name:  "KLESS_REQUEST_SOURCE",
									Value: j.RequestSource,
								},
								apiv1.EnvVar{
									Name:  "KLESS_RESPONSE_SOURCE",
									Value: j.ResponseSource,
								},
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse("100m"),
									apiv1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
							ImagePullPolicy: apiv1.PullIfNotPresent,
						},
					},
					Volumes: []apiv1.Volume{
						apiv1.Volume{
							Name: "dockersocket",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/var/run/docker.sock",
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					DNSPolicy:     apiv1.DNSClusterFirst,
				},
			},
		},
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	_, err = c.Batch().Jobs(j.Namespace).Create(jobSpec)
	switch {
	case err == nil:
		logger.Println("job created")
	case errors.IsNotFound(err):
		return fmt.Errorf("not found error")
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	return nil
}

//CreateService creates a service
func (k *K8sInterface) CreateService(serviceName string, appName string, namespace string) error {
	serviceSpec := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeNodePort,
			Selector: map[string]string{"app": appName},
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Name:     "handler",
					Protocol: apiv1.ProtocolTCP,
					Port:     8080,
				},
				apiv1.ServicePort{
					Name:     "frontend",
					Protocol: apiv1.ProtocolTCP,
					Port:     3080,
				},
			},
		},
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	service := c.Core().Services(namespace)
	svc, err := service.Get(appName, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceSpec.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err = service.Update(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to update service: %s", err)
		}
		logger.Println("service updated")
	case errors.IsNotFound(err):
		_, err = service.Create(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to create service: %s", err)
		}
		logger.Println("service created")
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	return nil
}

//DeleteDeployment deletes a deployment
func (k *K8sInterface) DeleteDeployment(name string, namespace string) error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploy := c.Extensions().Deployments(namespace)

	deployment, err := deploy.Get(name, metav1.GetOptions{})
	switch {
	case err == nil:
		logger.Println("got deployment controller")
	case errors.IsNotFound(err):
		return fmt.Errorf("deployment controller not found")
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	klessID := deployment.ObjectMeta.Labels["kless-id"]
	fmt.Printf("Found kless id of deployment with name = %s in namespace %s is %s\n", name, namespace, klessID)

	err = deploy.Delete(name, &metav1.DeleteOptions{OrphanDependents: boolp(false)})
	switch {
	case err == nil:
		logger.Println("deployment controller deleted")
	case errors.IsNotFound(err):
		return fmt.Errorf("deployment controller not found")
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	replicaSet := c.Extensions().ReplicaSets(namespace)

	replicaSetList, err := replicaSet.List(metav1.ListOptions{})
	switch {
	case err == nil:
		logger.Println("got ReplicaSet list")
	case !errors.IsNotFound(err):
		return fmt.Errorf("could not list ReplicaSets: %s", err)
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	var replicaSetName string

	for _, r := range replicaSetList.Items {
		rID := r.ObjectMeta.Labels["kless-id"]
		if "" != rID && rID == klessID {
			replicaSetName = r.ObjectMeta.Name
		}
	}

	if "" != replicaSetName {
		fmt.Printf("Deleting ReplicaSet %s\n", replicaSetName)

		err = replicaSet.Delete(replicaSetName, &metav1.DeleteOptions{})
		switch {
		case err == nil:
			logger.Println("replica set deleted")
		case errors.IsNotFound(err):
			return fmt.Errorf("replica set not found")
		default:
			return fmt.Errorf("unexpected error: %s", err)
		}
	} else {
		logger.Println("replica set name not found")
	}

	return nil
}

//DeleteService deletes a service
func (k *K8sInterface) DeleteService(serviceName string, namespace string) error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	service := c.Core().Services(namespace)
	err = service.Delete(serviceName, &metav1.DeleteOptions{})
	switch {
	case err == nil:
		logger.Println("service deleted")
	case errors.IsNotFound(err):
		logger.Println("service not found")
	default:
		return fmt.Errorf("unexpected error: %s", err)
	}

	return nil
}

func int32p(i int32) *int32 {
	r := new(int32)
	*r = i
	return r
}

func int64p(i int64) *int64 {
	r := new(int64)
	*r = i
	return r
}

func boolp(val bool) *bool {
	b := new(bool)
	*b = val
	return b
}
