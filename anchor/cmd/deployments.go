package cmd

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentsListOutput used to interact with template
type DeploymentsListOutput struct {
	Name              string            `json:"Name"`
	Labels            map[string]string `json:"Labels"`
	Namespace         string            `json:"Namespace"`
	CreationTimestamp string            `json:"CreationTimestamp"`
}

// GetDeployment used to get deployment by deployment name
func GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	return GetDeploymentClient(namespace).Get(name, metav1.GetOptions{})
}

// DeploymentCreate used to create deployment
func DeploymentCreate(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetDeploymentClient(namespace)
	result, err := client.Create(deployment)
	if err != nil {
		glog.Errorf("Err=%s", err)
		return nil, err
	}
	glog.V(2).Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// DeploymentsList used to list deployment
func DeploymentsList(namespace string) ([]DeploymentsListOutput, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetDeploymentClient(namespace)
	deploymentsListOutput := []DeploymentsListOutput{}
	deployments, err := client.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("List deployments failed : %v", err)
	}
	for _, deployment := range deployments.Items {

		// 过滤k8s部署
		if strings.Contains(deployment.Name, "kube") {
			continue
		}

		deploymentsListOutput = append(deploymentsListOutput, DeploymentsListOutput{
			Name:              deployment.Name,
			Labels:            deployment.Labels,
			Namespace:         deployment.Namespace,
			CreationTimestamp: deployment.CreationTimestamp.Format("2006-01-02 15:04"),
		})
	}
	return deploymentsListOutput, nil
}

// DeploymentDelete used to delete deployment
func DeploymentDelete(namespace, name string) error {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetDeploymentClient(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	return client.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// DeploymentUpdate used to update pod
func DeploymentUpdate(namespace string, deployment *appsv1.Deployment) *appsv1.Deployment {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}

	client := GetDeploymentClient(namespace)

	// result.Spec.Ports[0].Port = 3030
	deployment, err := client.Update(deployment)
	if err != nil {
		glog.Errorf("Update pod failed: %v", err)
	}
	return deployment
}

func int32Ptr(i int32) *int32 { return &i }
