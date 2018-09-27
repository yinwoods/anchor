package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/kubernetes/client-go/kubernetes/scheme"
	"github.com/yinwoods/anchor/anchor/cmd"
	"k8s.io/api/core/v1"
)

func servicesListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	services, err := cmd.ServicesList("")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "services.tmpl", services)
}

func serviceInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	service, err := cmd.ServiceGet(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	serviceJSON, err := json.Marshal(&service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, service)
	data = append(data, string(serviceJSON))
	c.HTML(http.StatusOK, "service_info.tmpl", data)
}

func serviceCreateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	body := c.PostForm("body")
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(body), nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	service := obj.(*v1.Service)
	_, err = cmd.ServiceCreate(service.Namespace, service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/services")
}

func serviceDeleteHandler(c *gin.Context) {
	type Input struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}
	var input Input
	c.BindJSON(&input)
	err := cmd.ServiceDelete(input.Namespace, input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func serviceUpdateHandler(c *gin.Context) {
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
	service := obj.(*v1.Service)
	cmd.ServiceUpdate(service.Namespace, service)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
