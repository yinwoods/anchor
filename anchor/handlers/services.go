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
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	services, err := cmd.ServicesList("")

	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "services.tmpl", services)
}

func serviceInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	service, err := cmd.ServiceGet(namespace, name)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	serviceJSON, err := json.Marshal(&service)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	var data []interface{}
	data = append(data, service)
	data = append(data, string(serviceJSON))
	c.HTML(http.StatusOK, "service_info.tmpl", data)
}

func serviceCreateHandler(c *gin.Context) {
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

	service := obj.(*v1.Service)
	_, err = cmd.ServiceCreate(service.Namespace, service)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
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
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"status": "fail",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func servicesUpdateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	body, ok := c.Params.Get("body")
	if !ok {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(body), nil, nil)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	pod := obj.(*v1.Pod)
	cmd.PodUpdate(pod.Namespace, pod)

	// TODO 改为查看pod详情页面
	pods, err := cmd.ContainersList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "pods.tmpl", pods)
}
