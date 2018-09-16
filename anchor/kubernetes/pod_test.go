package kubernetes

import (
	"testing"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var podName = "demo-pod"
var podClient PodClient

func TestCreatePod(*testing.T) {
	var json = `{
   "id":"nginx-mysql",
   "kind": "Pod",
   "apiVersion": "v1",
   "metadata": {
      "name": "demo-pod",
      "labels": {
         "name": "nginx-mysql"
      }
   },
   "spec": {
      "containers": [
         {
            "name": "nginx",
            "image": "nginx",
            "ports": [
               {
                  "hostPort": 85,
                  "containerPort": 80
               }
            ]
         },
         {
            "name": "mysql",
            "image": "mysql",
            "ports": [
               {
                  "hostPort": 3306,
                  "containerPort": 3306
               }
            ]
         }
      ]
   }
}`
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(json), nil, nil)
	if err != nil {
		glog.V(2).Infof("%#v", err)
	}

	pod := obj.(*v1.Pod)
	podClient.CreatePod(namespace, pod)
}

func TestListPod(*testing.T) {
	podClient.ListPod(namespace)
}

func TestUpdatePod(*testing.T) {
	podClient.UpdatePod(podName, namespace)
}

func TestDeletePod(*testing.T) {
	podClient.DeletePod(podName, namespace)
}
