package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

// DeploymentClient wraps client to operate deployment apis
type DeploymentClient struct {
	appsv1.DeploymentInterface
}

// DeploymentsListOutput used to interact with template
type DeploymentsListOutput struct {
	Name              string `json:"Name"`
	Namespace         string `json:"Namespace"`
	CreationTimestamp string `json:"CreationTimestamp"`
}

// GetDeployment used to get deployment by deployment name
func (client DeploymentClient) GetDeployment(namespace, name string) (*v1.Deployment, error) {
	return GetDeploymentClient(namespace).Get(name, metav1.GetOptions{})
}

// CreateDeployment used to create deployment
func (client DeploymentClient) CreateDeployment(deployment *v1.Deployment) (*v1.Deployment, error) {
	result, err := client.Create(deployment)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// UpdateDeployment used to update deployment
func (client DeploymentClient) UpdateDeployment(deploymentName string) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := client.Get(deploymentName, metav1.GetOptions{})
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
func (client DeploymentClient) DeploymentsList() ([]DeploymentsListOutput, error) {
	deploymentsListOutput := []DeploymentsListOutput{}
	deployments, err := client.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("List pods failed : %v", err)
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

// DeleteDeployment used to delete deployment
func (client DeploymentClient) DeleteDeployment(deploymentName string) {
	// Delete Deployment
	deletePolicy := metav1.DeletePropagationForeground
	if err := client.Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	glog.V(2).Infoln("Deleted deployment.")
}

func int32Ptr(i int32) *int32 { return &i }
