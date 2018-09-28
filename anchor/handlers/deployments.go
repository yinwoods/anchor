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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	namespace := c.Request.URL.Query().Get("namespace")
	deployments, err := cmd.DeploymentsList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "deployments.tmpl", deployments)
}

func deploymentDeleteHandler(c *gin.Context) {
	type Input struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}

	var input Input
	c.BindJSON(&input)
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
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	type Input struct {
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(input.Body), nil, nil)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	deployment := obj.(*appsv1.Deployment)
	_, err = cmd.DeploymentCreate(deployment.Namespace, deployment)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/deployments")
}

func deploymentInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	pod, err := cmd.GetDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	podJSON, err := json.Marshal(&pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, pod)
	data = append(data, string(podJSON))
	c.HTML(http.StatusOK, "deployment_info.tmpl", data)
}

func deploymentUpdateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	type Input struct {
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(input.Body), nil, nil)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	deployment := obj.(*appsv1.Deployment)
	cmd.DeploymentUpdate(deployment.Namespace, deployment)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
