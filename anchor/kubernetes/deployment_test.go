package kubernetes

import (
	"testing"

	"github.com/golang/glog"
	"k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var deploymentName = "demo-deployment"
var deploymentClient DeploymentClient

func TestCreateDeployment(*testing.T) {
	var json = `{
    "apiVersion":"apps/v1",
    "kind":"Deployment",
    "metadata":{
        "name":"demo-deployment"
    },
    "spec":{
        "replicas":2,
		"selector": {
        	"matchLabels": {
      			"app": "my-nginx"
  			}
  		},
        "template":{
            "metadata":{
                "labels":{
                    "app":"my-nginx"
                }
            },
            "spec":{
                "containers":[
                    {
                        "image":"nginx",
                        "name":"my-nginx",
                        "ports":[
                            {
                                "containerPort":80
                            }
                        ]
                    }
                ]
            }
        }
    }
}
`
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(json), nil, nil)
	if err != nil {
		glog.V(2).Infof("%#v", err)
	}

	deployment := obj.(*v1.Deployment)
	deploymentClient.CreateDeployment(apiv1.NamespaceDefault, deployment)
}

func TestListDeployment(*testing.T) {
	deploymentClient.ListDeployment(apiv1.NamespaceDefault)
}

func TestUpdateDeployment(*testing.T) {
	deploymentClient.UpdateDeployment(deploymentName, apiv1.NamespaceDefault)
}

func TestDeleteDeployment(*testing.T) {
	deploymentClient.DeleteDeployment(deploymentName, apiv1.NamespaceDefault)
}
