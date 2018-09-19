package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

// ServiceClient wraps client to operate service apis
type ServiceClient struct {
	corev1.ServiceInterface
}

// ServicesListOutput used to interact with template
type ServicesListOutput struct {
	Name      string `json:"Name"`
	Namespace string `json:"Namespace"`
}

// CreateService used to create service
func (client ServiceClient) CreateService(service *v1.Service) (*v1.Service, error) {
	result, err := client.Create(service)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	glog.V(2).Infof("Created service %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

// UpdateService used to update service
func (client ServiceClient) UpdateService(serviceName string) {
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

// GetService used to get service by service name
func (client ServiceClient) GetService(namespace, serviceName string) (*v1.Service, error) {
	return GetServiceClient(namespace).Get(serviceName, metav1.GetOptions{})
}

// ServicesList list services
func (client ServiceClient) ServicesList() ([]ServicesListOutput, error) {
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

// DeleteService used to delete service
func (client ServiceClient) DeleteService(serviceName string) {
	// Delete Service
	deletePolicy := metav1.DeletePropagationForeground
	if err := client.Delete(serviceName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	glog.V(2).Infoln("Deleted service.")
}
