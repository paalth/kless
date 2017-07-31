package servicehandler

import (
	"fmt"
	"log"
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
}

//EventHandlerFrontendInfo defines the frontend info
type EventHandlerFrontendInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Information map[string]string
}

func getJobName(e *EventHandlerInfo) string {
	return e.Namespace + "-" + e.Name + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}

func getServiceName(e *EventHandlerInfo) string {
	return e.Name + "-svc"
}

// CreateEventHandler builds and deploys an event handler
func (s *ServiceHandler) CreateEventHandler(e *EventHandlerInfo, f *EventHandlerFrontendInfo, frontendImageName string) error {
	fmt.Printf("Entering servicehandler.CreateEventHandler\n")

	k8s := &k.K8sInterface{}

	fmt.Printf("Creating namespace if it does not already exist\n")
	if err := k8s.CreateNamespace(e.Namespace); err != nil {
		log.Fatal(err)
	}

	klessRepo := "kless-registry.kless:5000"

	jobInfo := &k.JobInfo{
		JobName:             getJobName(e),
		Namespace:           e.Namespace,
		Image:               e.EventHandlerBuilderURL,
		KlessRepo:           klessRepo,
		EventHandlerName:    e.Name,
		EventHandlerVersion: e.Version,
		EventHandlerSource:  "kless-server.kless:8010/etcd?op=getsource&key=" + e.ID,
		InterfaceSource:     "kless-server.kless:8010/etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=KlessInterface",
		EntrypointSource:    "kless-server.kless:8010/etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=InvokeEventHandler",
		ContextSource:       "kless-server.kless:8010/etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=ContextSource",
		RequestSource:       "kless-server.kless:8010/etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=RequestSource",
		ResponseSource:      "kless-server.kless:8010/etcd?op=get&builder=" + e.EventHandlerBuilder + "&key=ResponseSource",
	}

	fmt.Printf("Creating job\n")
	if err := k8s.CreateJob(jobInfo); err != nil {
		log.Fatal(err)
	}

	eventHandlerImageName := klessRepo + "/" + e.Name + ":" + e.Version

	deploymentInfo := &k.DeploymentInfo{
		Namespace:               e.Namespace,
		Name:                    e.Name,
		Version:                 e.Version,
		Replicas:                1,
		EventHandlerName:        e.Name,
		EventHandlerImage:       eventHandlerImageName,
		EventHandlerPort:        8080,
		EventHandlerCPULimit:    "100m",
		EventHandlerMemoryLimit: "256Mi",
		FrontendName:            "frontend",
		FrontendImage:           frontendImageName,
		FrontendPort:            3080,
		FrontendCPULimit:        "100m",
		FrontendMemoryLimit:     "256Mi",
		FrontendEnvironmentVars: f.Information,
	}

	fmt.Printf("Creating deployment\n")
	if err := k8s.CreateDeployment(deploymentInfo); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating service\n")
	if err := k8s.CreateService(getServiceName(e), e.Name, e.Namespace); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Leaving servicehandler.CreateEventHandler\n")

	return nil
}

// DeleteEventHandler removes a deployed event handler
func (s *ServiceHandler) DeleteEventHandler(e *EventHandlerInfo) error {

	fmt.Printf("Entering servicehandler.DeleteEventHandler\n")

	k8s := &k.K8sInterface{}

	if err := k8s.DeleteDeployment(e.Name, e.Namespace); err != nil {
		log.Fatal(err)
	}

	// TODO: delete namespace if this is the last handler in the namespace
	// TODO: delete replica sets
	// TODO: delete running pods

	if err := k8s.DeleteService(getServiceName(e), e.Namespace); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Leaving servicehandler.DeleteEventHandler\n")

	return nil
}
