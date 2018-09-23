package kubernetes

import (
	"testing"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var serviceName = "demo-service"
var serviceClient ServiceClient

func TestCreateService(*testing.T) {
	var json = `{
	"apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "labels": {
            "name": "nginxservice"
        },
        "name": "demo-service",
		"namespace": "default"
    },
    "spec": {
        "ports": [
            {
                "port": 80,
                "protocol": "TCP",
                "targetPort": 80
            }
        ],
        "selector": {
            "app": "nginx"
        },
        "type": "LoadBalancer"
    }
}`
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(json), nil, nil)
	if err != nil {
		glog.V(2).Infof("%#v", err)
	}

	service := obj.(*v1.Service)
	serviceClient.CreateService(namespace, service)
}

func TestListService(*testing.T) {
	serviceClient.ListService()
}

func TestUpdateService(*testing.T) {
	serviceClient.ServiceUpdate(serviceName, namespace)
}

func TestDeleteService(*testing.T) {
	serviceClient.DeleteService(serviceName, namespace)
}
