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

// PodCreate used to create pod
func PodCreate(namespace string, pod *v1.Pod) (*v1.Pod, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetPodClient(namespace)
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
func PodUpdate(namespace string, pod *v1.Pod) *v1.Pod {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetPodClient(namespace)
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

// PodGet used to get pod by pod name
func PodGet(namespace, name string) (*v1.Pod, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	return GetPodClient(namespace).Get(name, metav1.GetOptions{})
}

// PodsList used to list pod
func PodsList(namespace string) ([]PodsListOutput, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}

	client := GetPodClient(namespace)
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
func DeletePod(namespace, name string) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetPodClient(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := client.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	glog.V(2).Infoln("Deleted pod.")
}
