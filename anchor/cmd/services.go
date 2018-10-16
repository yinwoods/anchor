package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServicesListOutput used to interact with template
type ServicesListOutput struct {
	Name      string            `json:"Name"`
	Namespace string            `json:"Namespace"`
	Labels    map[string]string `json:"Labels"`
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
func ServiceUpdate(namespace string, service *v1.Service) *v1.Service {
	if namespace == "" {
		namespace = v1.NamespaceAll
	}

	client := GetServiceClient(namespace)

	service, err := client.Update(service)
	if err != nil {
		glog.Errorf("Update service failed: %v", err)
	}
	return service
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
