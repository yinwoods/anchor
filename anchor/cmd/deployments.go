package cmd

import (
	"fmt"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// DeploymentsListOutput used to interact with template
type DeploymentsListOutput struct {
	Name              string `json:"Name"`
	Namespace         string `json:"Namespace"`
	CreationTimestamp string `json:"CreationTimestamp"`
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
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// UpdateDeployment used to update deployment
func UpdateDeployment(namespace, name string) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetDeploymentClient(namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := client.Get(name, metav1.GetOptions{})
		if getErr != nil {
			glog.Errorf("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := client.Update(result)
		return updateErr
	})
	if retryErr != nil {
		glog.Errorf("Update deployment failed: %v", retryErr)
	}
	glog.V(2).Infoln("Updated deployment...")
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
		deploymentsListOutput = append(deploymentsListOutput, DeploymentsListOutput{
			Name:              deployment.Name,
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
