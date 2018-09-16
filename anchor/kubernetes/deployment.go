package kubernetes

import (
	"github.com/golang/glog"
	"k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

// DeploymentClient wraps client to operate deployment apis
type DeploymentClient struct {
	appsv1.DeploymentInterface
}

// CreateDeployment used to create deployment
func (client DeploymentClient) CreateDeployment(namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {

	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	glog.V(2).Infoln("Creating deployment...")
	result, err := client.Create(deployment)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// UpdateDeployment used to update deployment
func (client DeploymentClient) UpdateDeployment(deploymentName, namespace string) {

	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	glog.V(2).Infoln("Updating deployment...")

	// Update Deployment
	//    You have two options to Update() this Deployment:
	//
	//    1. Modify the "deployment" variable and call: Update(deployment).
	//       This works like the "kubectl replace" command and it overwrites/loses changes
	//       made by other clients between you Create() and Update() the object.
	//    2. Modify the "result" returned by Get() and retry Update(result) until
	//       you no longer get a conflict error. This way, you can preserve changes made
	//       by other clients between Create() and Update(). This is implemented below
	//			 using the retry utility package included with client-go. (RECOMMENDED)
	//
	// More Info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#concurrency-control-and-consistency

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

// ListDeployment used to list deployment
func (client DeploymentClient) ListDeployment(namespace string) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	glog.V(2).Infof("Listing deployments in namespace %q:\n", namespace)

	// List Deployments
	list, err := client.List(metav1.ListOptions{})
	if err != nil {
		glog.Error(err)
		return
	}
	for _, d := range list.Items {
		glog.V(2).Infof(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

// DeleteDeployment used to delete deployment
func (client DeploymentClient) DeleteDeployment(deploymentName, namespace string) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	glog.V(2).Infoln("Deleting deployment...")

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
