package cmd

import (
	"flag"
	"math/rand"
	"path/filepath"

	docker "github.com/docker/docker/client"
	"github.com/golang/glog"
	"github.com/kubernetes/client-go/util/homedir"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// DockerClient wraps docker client
var DockerClient *docker.Client

// K8SClient wraps kubernetes client
var K8SClient KubernetesClient

var ups []PowerSuppliesListOutput
var ref []RefrigerationsListOutput

func init() {
	var err error
	DockerClient, err = docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	K8SClient = GetK8SClient()

	for i := 0; i < randRange(3, 5); i++ {
		ups = append(ups, randomPowerSupplyList(i))
	}
	for i := 0; i < randRange(3, 5); i++ {
		ref = append(ref, randomRefgerationList(i))
	}

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
		glog.Errorf("Err=%s", err)
		panic(err)
	}
	K8SClient.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Err=%s", err)
		panic(err)
	}
	return K8SClient
}

// SetNamespace used to set k8s cliet namespace
func (client *KubernetesClient) SetNamespace(namespace string) {
	client.Namespace = namespace
}

// GetDeploymentClient used to get Deployment Client by namespace
func GetDeploymentClient(namespace string) appsv1.DeploymentInterface {
	return K8SClient.Clientset.Apps().Deployments(namespace)
}

// GetNodeClient used to get nodes client
func GetNodeClient() corev1.NodeInterface {
	return K8SClient.Clientset.CoreV1().Nodes()
}

// GetServiceClient used to get service client
func GetServiceClient(namespace string) corev1.ServiceInterface {
	return K8SClient.Clientset.CoreV1().Services(namespace)
}

// GetPodClient used to get pod client
func GetPodClient(namespace string) corev1.PodInterface {
	return K8SClient.Clientset.CoreV1().Pods(namespace)
}

// GetPodTemplateClient used to get pod template client
func GetPodTemplateClient(namespace string) corev1.PodTemplateInterface {
	return K8SClient.Clientset.CoreV1().PodTemplates(namespace)
}
