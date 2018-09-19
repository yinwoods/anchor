package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func deploymentsListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	namespace := c.Request.URL.Query().Get("namespace")
	cmd.K8SClient.SetNamespace(namespace)
	pods, err := cmd.K8SClient.DeploymentClient.DeploymentsList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "deployments.tmpl", pods)
}

func deploymentInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	pod, err := cmd.K8SClient.GetDeployment(namespace, name)
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
