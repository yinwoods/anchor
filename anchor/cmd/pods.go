package cmd

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
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
	StartTime string `json:"StartTime"`
	Phase     string `json:"Phase"`
}

type grandChildren struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type children struct {
	Name     string          `json:"name"`
	Children []grandChildren `json:"children"`
}

// PodsContainersJSON used to store pods and containers
type PodsContainersJSON struct {
	Name     string     `json:"name"`
	Children []children `json:"children"`
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
		glog.Errorf("Err=%s", err)
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

	// result.Spec.Ports[0].Port = 3030
	pod, err := client.Update(pod)
	if err != nil {
		glog.Errorf("Update pod failed: %v", err)
	}
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
			StartTime: pod.CreationTimestamp.Format("2006-01-02 15:04"),
			Phase:     string(pod.Status.Phase),
		})
	}
	return podsListOutput, nil
}

// PodDelete used to delete pod by podName
func PodDelete(namespace, name string) error {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetPodClient(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	return client.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// PodContainersList list pods and containers
func PodContainersList(namespace string) (PodsContainersJSON, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetPodClient(namespace)
	pods, err := client.List(metav1.ListOptions{})
	if err != nil {
		return PodsContainersJSON{}, fmt.Errorf("List pods failed : %v", err)
	}

	result := PodsContainersJSON{
		Name:     "容器组",
		Children: []children{},
	}

	podsChildren := []children{}

	for _, pod := range pods.Items {
		containers := []grandChildren{}
		for _, status := range pod.Status.ContainerStatuses {
			containerID := strings.Split(status.ContainerID, "//")[1]
			container, err := ContainerGet(containerID)
			if err != nil {
				glog.Error("container id not found ", containerID)
				continue
			}
			containers = append(containers, grandChildren{
				Name: container.Name,
				Size: 100,
			})
		}
		podsChildren = append(podsChildren, children{
			Name:     pod.Name,
			Children: containers,
		})
	}
	result.Children = podsChildren
	return result, nil
}
