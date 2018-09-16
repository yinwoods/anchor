package kubernetes

var namespace = "kube-system"

func init() {
	client := GetClient(namespace)
	nodeClient = client.GetNodeClient()
	deploymentClient = client.GetDeploymentClient(namespace)
	serviceClient = client.GetServiceClient(namespace)
	podClient = client.GetPodClient(namespace)
}
