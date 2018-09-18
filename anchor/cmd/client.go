package cmd

import (
	"flag"
	"math/rand"
	"path/filepath"
	"time"

	docker "github.com/docker/docker/client"
	"github.com/golang/glog"
	"github.com/kubernetes/client-go/util/homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// DockerClient wraps docker client
var DockerClient *docker.Client

// K8SClient wraps kubernetes client
var K8SClient KubernetesClient

func init() {
	var err error
	DockerClient, err = docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	K8SClient = GetK8SClient()
}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func randString(str ...string) string {
	return str[rand.Intn(len(str))]
}

// KubernetesClient wraps kubernetes Clientset
type KubernetesClient struct {
	Namespace string
	Clientset *kubernetes.Clientset
	NodeClient
	DeploymentClient
	ServiceClient
	PodClient
}

// GetK8SClient used to get client
func GetK8SClient() KubernetesClient {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	K8SClient.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	K8SClient.DeploymentClient = GetDeploymentClient(K8SClient.Namespace)
	K8SClient.NodeClient = GetNodeClient()
	K8SClient.PodClient = GetPodClient(K8SClient.Namespace)
	K8SClient.ServiceClient = GetServiceClient(K8SClient.Namespace)
	return K8SClient
}

// SetNamespace used to set k8s cliet namespace
func (client *KubernetesClient) SetNamespace(namespace string) {
	client.Namespace = namespace
	client.DeploymentClient = GetDeploymentClient(namespace)
	client.NodeClient = GetNodeClient()
	client.PodClient = GetPodClient(namespace)
	client.ServiceClient = GetServiceClient(namespace)
}

// GetDeploymentClient used to get Deployment Client by namespace
func GetDeploymentClient(namespace string) DeploymentClient {
	return DeploymentClient{K8SClient.Clientset.Apps().Deployments(namespace)}
}

// GetNodeClient used to get nodes client
func GetNodeClient() NodeClient {
	return NodeClient{K8SClient.Clientset.CoreV1().Nodes()}
}

// GetServiceClient used to get service client
func GetServiceClient(namespace string) ServiceClient {
	return ServiceClient{K8SClient.Clientset.CoreV1().Services(namespace)}
}

// GetPodClient used to get pod client
func GetPodClient(namespace string) PodClient {
	return PodClient{K8SClient.Clientset.CoreV1().Pods(namespace)}
}
