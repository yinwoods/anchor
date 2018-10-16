package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
	"github.com/yinwoods/anchor/anchor/util"
)

var (
	apiKey = util.GeneratePassword(32)
)

func apiGraphInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := util.GenerateGraph()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, result)

}

func apiContainerUpdateConfigInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	cid := c.Param("cid")

	containerJSON, err := cmd.ContainerGet(cid)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	config := container.UpdateConfig{
		Resources:     containerJSON.HostConfig.Resources,
		RestartPolicy: containerJSON.HostConfig.RestartPolicy,
	}
	configJSON, err := json.Marshal(config)

	var out bytes.Buffer
	json.Indent(&out, []byte(configJSON), "", "  ")
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": out.String(),
		"cid":    containerJSON.ID,
	})
}

func apiImageInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	mid := c.Param("mid")

	_, imageJSON, err := cmd.ImageGet(mid)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(imageJSON), "", "  ")

	c.JSON(http.StatusOK, gin.H{
		"result": out.String(),
	})
}

func apiPodInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	namespace := c.Param("namespace")
	name := c.Param("name")

	url := apiURLPrefix + fmt.Sprintf("%s/pods/%s", namespace, name)
	podJSON, err := httpGet(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(podJSON), "", "  ")

	c.JSON(http.StatusOK, gin.H{
		"result": out.String(),
	})
}

func apiServiceInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	namespace := c.Param("namespace")
	name := c.Param("name")

	url := apiURLPrefix + fmt.Sprintf("%s/services/%s", namespace, name)
	podJSON, err := httpGet(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(podJSON), "", "  ")

	c.JSON(http.StatusOK, gin.H{
		"result": out.String(),
	})
}

func apiDeploymentInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	namespace := c.Param("namespace")
	name := c.Param("name")

	url := apiV1Prefix + fmt.Sprintf("%s/deployments/%s", namespace, name)
	podJSON, err := httpGet(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(podJSON), "", "  ")

	c.JSON(http.StatusOK, gin.H{
		"result": out.String(),
	})
}
