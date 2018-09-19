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

func podsListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	namespace := c.Request.URL.Query().Get("namespace")
	cmd.K8SClient.SetNamespace(namespace)
	pods, err := cmd.K8SClient.PodClient.PodsList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "pods.tmpl", pods)
}

func podInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	pod, err := cmd.K8SClient.GetPod(namespace, name)
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
	c.HTML(http.StatusOK, "pod_info.tmpl", data)
}

func podsUpdateHandler(c *gin.Context) {
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
