package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/kubernetes/client-go/kubernetes/scheme"
	"github.com/yinwoods/anchor/anchor/cmd"
	appsv1 "k8s.io/api/apps/v1"
)

func deploymentsListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	namespace := c.Request.URL.Query().Get("namespace")
	deployments, err := cmd.DeploymentsList(namespace)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "deployments.tmpl", deployments)
}

func deploymentsDeleteHandler(c *gin.Context) {
	type Input struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}

	var input Input
	c.BindJSON(&input)
	glog.V(3).Infoln("namespace: ", input.Namespace)
	glog.V(3).Infoln("name: ", input.Name)
	err := cmd.DeploymentDelete(input.Namespace, input.Name)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"status": "fail",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func deploymentCreateHandler(c *gin.Context) {
	// TODO 创建成功信息返回
	err := parseSessionCookie(c)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	body := c.PostForm("body")
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(body), nil, nil)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	deployment := obj.(*appsv1.Deployment)
	_, err = cmd.DeploymentCreate(deployment.Namespace, deployment)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	c.Redirect(http.StatusFound, "/deployments")
}

func deploymentInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	pod, err := cmd.GetDeployment(namespace, name)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	podJSON, err := json.Marshal(&pod)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	var data []interface{}
	data = append(data, pod)
	data = append(data, string(podJSON))
	c.HTML(http.StatusOK, "deployment_info.tmpl", data)
}
