package servicehandler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	k "github.com/paalth/kless/pkg/k8sinterface"
)

//ServiceHandler manages event handlers
type ServiceHandler struct {
}

//EventHandlerInfo defines the info passed to the job that deploys an event handler
type EventHandlerInfo struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Namespace              string `json:"namespace"`
	Version                string `json:"version"`
	EventHandlerBuilder    string `json:"eventhandlerbuilder"`
	EventHandlerBuilderURL string `json:"eventhandlerbuilderurl"`
	Frontend               string `json:"frontend"`
	DependenciesURL        string `json:"dependenciesurl"`
	Comment                string `json:"comment"`
}

//EventHandlerFrontendInfo defines the frontend info
type EventHandlerFrontendInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Information map[string]string
	Comment     string `json:"comment"`
}

func getJobName(e *EventHandlerInfo) string {
	return e.Namespace + "-" + e.Name + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}

//GetServiceName returns the service name for a handler
func (s *ServiceHandler) GetServiceName(e *EventHandlerInfo) string {
	return e.Name + "-svc"
}

// GetKlessServer returns the hostname:port of the Kless server
func (s *ServiceHandler) GetKlessServer() string {
	return "kless-server." + os.Getenv("SERVER_NAMESPACE") + ":8010"
}

// GetKlessRegistry returns the hostname:port of the Kless registry
func (s *ServiceHandler) GetKlessRegistry() string {
	return "kless-registry." + os.Getenv("SERVER_NAMESPACE") + ":5000"
}

// GetClusterIngressWildcard returns the DNS wildcard to use for ingress resources
func (s *ServiceHandler) GetClusterIngressWildcard() string {
	return os.Getenv("INGRESS_DNS_WILDCARD")
}

// BuildEventHandler starts the build of an event handler
func (s *ServiceHandler) BuildEventHandler(e *EventHandlerInfo) error {
	fmt.Printf("Entering servicehandler.BuildEventHandler\n")

	k8s := &k.K8sInterface{}

	fmt.Printf("Creating namespace if it does not already exist\n")
	if err := k8s.CreateNamespace(e.Namespace); err != nil {
		log.Fatal(err)
	}

	klessServer := s.GetKlessServer()
	klessRegistry := s.GetKlessRegistry()

	jobInfo := &k.JobInfo{
		JobName:             getJobName(e),
		Namespace:           e.Namespace,
		Image:               e.EventHandlerBuilderURL,
		KlessServer:         klessServer,
		KlessRegistry:       klessRegistry,
		EventHandlerID:      e.ID,
		EventHandlerName:    e.Name,
		EventHandlerVersion: e.Version,
		EventHandlerSource:  "etcd?op=getsource&key=" + e.ID,
		InterfaceSource:     "etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=KlessInterface",
		EntrypointSource:    "etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=InvokeEventHandler",
		ContextSource:       "etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=ContextSource",
		RequestSource:       "etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=RequestSource",
		ResponseSource:      "etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=ResponseSource",
		DependenciesURL:     e.DependenciesURL,
	}

	fmt.Printf("Creating job\n")
	if err := k8s.CreateJob(jobInfo); err != nil {
		log.Fatal(err)
	}

	return nil
}

// DeployEventHandler deploys an event handler after a successful build
func (s *ServiceHandler) DeployEventHandler(e *EventHandlerInfo, f *EventHandlerFrontendInfo, frontendImageName string) error {
	fmt.Printf("Entering servicehandler.DeployEventHandler\n")

	k8s := &k.K8sInterface{}

	klessRegistry := s.GetKlessRegistry()

	eventHandlerImageName := klessRegistry + "/" + e.Name + ":" + e.Version

	frontendPortNumber, err := strconv.Atoi(f.Information["KLESS_FRONTEND_HTTP_PORT"])
	if err != nil {
		log.Fatal(err)
	}

	eventHandlerPort := int32(8080)
	frontendPort := int32(frontendPortNumber)

	deploymentInfo := &k.DeploymentInfo{
		Namespace:               e.Namespace,
		Name:                    e.Name,
		Version:                 e.Version,
		Replicas:                1,
		EventHandlerName:        e.Name,
		EventHandlerImage:       eventHandlerImageName,
		EventHandlerPort:        eventHandlerPort,
		EventHandlerCPULimit:    "100m",
		EventHandlerMemoryLimit: "256Mi",
		FrontendName:            "frontend",
		FrontendImage:           frontendImageName,
		FrontendPort:            frontendPort,
		FrontendCPULimit:        "100m",
		FrontendMemoryLimit:     "256Mi",
		FrontendEnvironmentVars: f.Information,
	}

	fmt.Printf("Creating deployment\n")
	if err := k8s.CreateDeployment(deploymentInfo); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating service\n")
	if err := k8s.CreateService(s.GetServiceName(e), e.Name, e.Namespace, eventHandlerPort, frontendPort); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating ingress\n")
	ingressHostname := e.Name + "." + s.GetClusterIngressWildcard()
	ingressPath := "/"
	if err := k8s.CreateIngress(e.Name, e.Namespace, ingressHostname, ingressPath, s.GetServiceName(e), frontendPort); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Leaving servicehandler.CreateEventHandler\n")

	return nil
}

// DeleteEventHandler removes a deployed event handler
func (s *ServiceHandler) DeleteEventHandler(e *EventHandlerInfo) error {

	fmt.Printf("Entering servicehandler.DeleteEventHandler\n")

	k8s := &k.K8sInterface{}

	// Only printing messages on failure to delete objects as
	// they could have been manually deleted without our knowledge

	fmt.Printf("Deleting deployment\n")
	if err := k8s.DeleteDeployment(e.Name, e.Namespace); err != nil {
		fmt.Printf("Caught error when attempting to delete Deployment\n")
	}

	// TODO: delete namespace if this is the last handler in the namespace
	// TODO: delete replica sets
	// TODO: delete running pods

	fmt.Printf("Deleting service\n")
	if err := k8s.DeleteService(s.GetServiceName(e), e.Namespace); err != nil {
		fmt.Printf("Caught error when attempting to delete Service\n")
	}

	fmt.Printf("Deleting ingress\n")
	if err := k8s.DeleteIngress(e.Name, e.Namespace); err != nil {
		fmt.Printf("Caught error when attempting to delete Ingress\n")
	}

	fmt.Printf("Leaving servicehandler.DeleteEventHandler\n")

	return nil
}
