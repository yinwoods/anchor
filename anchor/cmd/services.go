package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// ServicesListOutput used to interact with template
type ServicesListOutput struct {
	Name      string `json:"Name"`
	Namespace string `json:"Namespace"`
}

// ServiceCreate used to create service
func ServiceCreate(namespace string, service *v1.Service) (*v1.Service, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetServiceClient(namespace)
	result, err := client.Create(service)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created service %q.\n", result.Name)
	return result, nil
}

// ServiceUpdate used to update service
func ServiceUpdate(namespace, serviceName string) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetServiceClient(namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := client.Get(serviceName, metav1.GetOptions{})
		if getErr != nil {
			glog.Errorf("Failed to get latest version of service: %v", getErr)
			return getErr
		}

		result.Spec.Ports[0].Port = 3030
		_, updateErr := client.Update(result)
		return updateErr
	})
	if retryErr != nil {
		glog.Errorf("Update service failed: %v", retryErr)
	}
	glog.V(2).Infoln("Updated service...")
}

// ServiceGet used to get service by service name
func ServiceGet(namespace, serviceName string) (*v1.Service, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	return GetServiceClient(namespace).Get(serviceName, metav1.GetOptions{})
}

// ServicesList list services
func ServicesList(namespace string) ([]ServicesListOutput, error) {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetServiceClient(namespace)
	servicesListOutput := []ServicesListOutput{}
	services, err := client.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("List services failed : %v", err)
	}
	for _, service := range services.Items {
		servicesListOutput = append(servicesListOutput, ServicesListOutput{
			Name:      service.Name,
			Namespace: service.Namespace,
		})
	}
	return servicesListOutput, nil
}

// ServiceDelete used to delete service
func ServiceDelete(namespace, name string) error {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}
	client := GetServiceClient(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := client.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	return err
}
