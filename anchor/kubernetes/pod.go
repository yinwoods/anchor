package kubernetes

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

// PodClient wraps client to operate pod apis
type PodClient struct {
	corev1.PodInterface
}

// CreatePod used to create pod
func (client PodClient) CreatePod(namespace string, pod *v1.Pod) (*v1.Pod, error) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	glog.V(2).Infoln("Creating pod...")
	result, err := client.Create(pod)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created pod %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// UpdatePod used to update pod
func (client PodClient) UpdatePod(podName, namespace string) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	glog.V(2).Infoln("Updating pod...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := client.Get(podName, metav1.GetOptions{})
		if getErr != nil {
			glog.Errorf("Failed to get latest version of pod: %v", getErr)
			return getErr
		}

		// result.Spec.Ports[0].Port = 3030
		_, updateErr := client.Update(result)
		return updateErr
	})
	if retryErr != nil {
		glog.Errorf("Update pod failed: %v", retryErr)
	}
	glog.V(2).Infoln("Updated pod...")
}

// GetPod used to get pod by pod name
func (client PodClient) GetPod(podName string) (*v1.Pod, error) {
	return client.Get(podName, metav1.GetOptions{})
}

// ListPod used to list pod
func (client PodClient) ListPod(namespace string) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	glog.V(2).Infof("Listing pods in namespace: %q:\n", namespace)
	list, err := client.List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("List pod failed : %v", err)
	}
	for _, p := range list.Items {
		glog.V(2).Infof(" * %s\n", p.Name)
	}
}

// DeletePod used to delete pod by podName
func (client PodClient) DeletePod(podName, namespace string) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	glog.V(2).Infoln("Deleting pod...")

	// Delete Pod
	deletePolicy := metav1.DeletePropagationForeground
	if err := client.Delete(podName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	glog.V(2).Infoln("Deleted pod.")
}
