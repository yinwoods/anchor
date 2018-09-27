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
	pods, err := cmd.PodsList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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
	pod, err := cmd.PodGet(namespace, name)
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
	c.HTML(http.StatusOK, "pod_info.tmpl", data)
}

func podsUpdateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
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
	pod := obj.(*v1.Pod)

	pod = cmd.PodUpdate(pod.Namespace, pod)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func podCreateHandler(c *gin.Context) {
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

	pod := obj.(*v1.Pod)
	_, err = cmd.PodCreate(pod.Namespace, pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/pods")
}

func podDeleteHandler(c *gin.Context) {
	type Input struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}
	var input Input
	c.BindJSON(&input)
	err := cmd.PodDelete(input.Namespace, input.Name)
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
