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
		return
	}

	services, err := cmd.K8SClient.ServiceClient.ServicesList()

	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "services.tmpl", services)
}

func serviceInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	service, err := cmd.K8SClient.GetService(namespace, name)
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

func servicesUpdateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
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
	cmd.K8SClient.PodClient.PodUpdate(pod)

	// TODO 改为查看pod详情页面
	pods, err := cmd.ContainersList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "pods.tmpl", pods)
}
