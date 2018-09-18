package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

// PodClient wraps client to operate pod apis
type PodClient struct {
	corev1.PodInterface
}

// PodsListOutput used to interact with template
type PodsListOutput struct {
	Name      string `json:"Name"`
	Namespace string `json:"Namespace"`
	HostIP    string `json:"HostIP"`
	PodIP     string `json:"PodIP"`
	StartTime string `json:"StartTime"`
	Phase     string `json:"Phase"`
}

// CreatePod used to create pod
func (client PodClient) CreatePod(pod *v1.Pod) (*v1.Pod, error) {
	glog.V(2).Infoln("Creating pod...")
	result, err := client.Create(pod)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created pod %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// PodUpdate used to update pod
func (client PodClient) PodUpdate(pod *v1.Pod) *v1.Pod {
	glog.V(2).Infoln("Updating pod...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := client.Get(pod.Name, metav1.GetOptions{})
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
	return pod
}

// GetPod used to get pod by pod name
func (client PodClient) GetPod(podNamespace, podName string) (*v1.Pod, error) {
	glog.V(2).Infoln("namespace: ", podNamespace)
	return GetPodClient(podNamespace).Get(podName, metav1.GetOptions{})
}

// PodsList used to list pod
func (client PodClient) PodsList() ([]PodsListOutput, error) {

	podsListOutput := []PodsListOutput{}
	pods, err := client.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("List pods failed : %v", err)
	}
	for _, pod := range pods.Items {
		podsListOutput = append(podsListOutput, PodsListOutput{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			HostIP:    pod.Status.HostIP,
			PodIP:     pod.Status.PodIP,
			StartTime: pod.Status.StartTime.Format("2006-01-02 15:04"),
			Phase:     string(pod.Status.Phase),
		})
	}
	return podsListOutput, nil
}

// DeletePod used to delete pod by podName
func (client PodClient) DeletePod(podName string) {
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
