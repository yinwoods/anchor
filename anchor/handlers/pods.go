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

	podNamespace := c.Param("podNamespace")
	podName := c.Param("podname")
	pod, err := cmd.K8SClient.GetPod(podNamespace, podName)
	if err != nil {
		glog.Error("get pod by name error: ", err)
		c.HTML(http.StatusBadRequest, "pod_info.tmpl", v1.Pod{})
	}
	podJSON, err := json.Marshal(&pod)
	if err != nil {
		glog.Error("get pod by name error: ", err)
		c.HTML(http.StatusBadRequest, "pod_info.tmpl", v1.Pod{})
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
		glog.Error("get pods update body error")
		c.HTML(http.StatusBadRequest, "pods.tmpl", cmd.ContainersListOutput{})
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(body), nil, nil)
	if err != nil {
		glog.Errorf("decode pods update body error: %#v", err)
		c.HTML(http.StatusBadRequest, "pods.tmpl", cmd.ContainersListOutput{})
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
