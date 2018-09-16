package kubernetes

import (
	"flag"
	"path/filepath"

	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client wraps kubernetes Clientset
type Client struct {
	Clientset *kubernetes.Clientset
	NodeClient
	DeploymentClient
	ServiceClient
	PodClient
}

// GetClient used to get client
func GetClient(namespace string) Client {
	var client = Client{}
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
	client.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	client.DeploymentClient = client.GetDeploymentClient(namespace)
	client.NodeClient = client.GetNodeClient()
	client.PodClient = client.GetPodClient(namespace)
	client.ServiceClient = client.GetServiceClient(namespace)
	return client
}

// GetDeploymentClient used to get Deployment Client by namespace
func (client Client) GetDeploymentClient(namespace string) DeploymentClient {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	return DeploymentClient{client.Clientset.Apps().Deployments(namespace)}
}

// GetNodeClient used to get nodes client
func (client Client) GetNodeClient() NodeClient {
	return NodeClient{client.Clientset.CoreV1().Nodes()}
}

// GetServiceClient used to get service client
func (client Client) GetServiceClient(namespace string) ServiceClient {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	return ServiceClient{client.Clientset.CoreV1().Services(namespace)}
}

// GetPodClient used to get pod client
func (client Client) GetPodClient(namespace string) PodClient {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	return PodClient{client.Clientset.CoreV1().Pods(namespace)}
}
