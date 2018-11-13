package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
	"github.com/yinwoods/anchor/anchor/util"
)

var (
	apiKey = util.GeneratePassword(32)
)

func apiTokensHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"currentTime":     time.Now(),
		"availableTokens": rateLimiter.Available(),
	})
}

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

func apiContainerInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	id := c.Param("id")

	containerJSON, err := cmd.ContainerGet(id)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
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
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": out.String(),
		"id":     containerJSON.ID,
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
	id := c.Param("id")

	_, imageJSON, err := cmd.ImageGet(id)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
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

func apiUPSInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	id := c.Param("id")
	ups, err := cmd.UPSGet(id)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	configJSON, err := json.Marshal(ups)
	var out bytes.Buffer
	json.Indent(&out, []byte(configJSON), "", "  ")
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     ups.ID,
		"result": out.String(),
	})
}

func apiREFInfo(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	id := c.Param("id")
	ref, err := cmd.REFGet(id)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	configJSON, err := json.Marshal(ref)
	var out bytes.Buffer
	json.Indent(&out, []byte(configJSON), "", "  ")
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     ref.ID,
		"result": out.String(),
	})
}
