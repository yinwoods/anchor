package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type sysInfo struct {
	APIVersion string `json:"apiVersion"`
	Total      struct {
		Allocatable struct {
			CPU     string `json:"cpu"`
			Storage string `json:"ephemeral-storage"`
			Memory  string `json:"memory"`
			Pods    string `json:"pods"`
		} `json:"allocatable"`
		Capacity struct {
			CPU     string `json:"cpu"`
			Storage string `json:"ephemeral-storage"`
			Memory  string `json:"memory"`
			Pods    string `json:"pods"`
		} `json:"capacity"`
	} `json:"total"`
	Items []struct {
		Status struct {
			Addresses []struct {
				Address string `json:"address"`
				Type    string `json:"type"`
			} `json:"addresses"`
			Allocatable struct {
				CPU     string `json:"cpu"`
				Storage string `json:"ephemeral-storage"`
				Memory  string `json:"memory"`
				Pods    string `json:"pods"`
			} `json:"allocatable"`
			Capacity struct {
				CPU     string `json:"cpu"`
				Storage string `json:"ephemeral-storage"`
				Memory  string `json:"memory"`
				Pods    string `json:"pods"`
			} `json:"capacity"`
			NodeInfo struct {
				Architecture            string `json:"architecture"`
				ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
				KernelVersion           string `json:"kernelVersion"`
				OperatingSystem         string `json:"operatingSystem"`
				OsImage                 string `json:"osImage"`
			} `json:"nodeInfo"`
		} `json:"status"`
	} `json:"items"`
	Kind string `json:"kind"`
}

func isDotOrZero(ch rune) bool {
	return ch == '0' || ch == '.'
}

func apiTokensHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"currentTime":     time.Now(),
		"availableTokens": rateLimiter.Available(),
	})
}

func apiSysInfoHandler(c *gin.Context) {
	resp, err := httpGet(apiNodesURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var info sysInfo
	err = json.Unmarshal(resp, &info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for index, item := range info.Items {

		allocatable := &item.Status.Allocatable
		capacity := &item.Status.Capacity

		info.Total.Allocatable.CPU, _ = util.StringAdd(info.Total.Allocatable.CPU, allocatable.CPU)
		info.Total.Allocatable.Pods, _ = util.StringAdd(info.Total.Allocatable.Pods, allocatable.Pods)

		info.Total.Allocatable.Memory, _ = util.StringAdd(info.Total.Allocatable.Memory, allocatable.Memory[:len(allocatable.Memory)-2])
		info.Total.Allocatable.Storage, _ = util.StringAdd(info.Total.Allocatable.Storage, allocatable.Storage[:len(allocatable.Storage)-2])

		info.Total.Capacity.CPU, _ = util.StringAdd(info.Total.Capacity.CPU, capacity.CPU)
		info.Total.Capacity.Pods, _ = util.StringAdd(info.Total.Capacity.Pods, capacity.Pods)
		info.Total.Capacity.Memory, _ = util.StringAdd(info.Total.Capacity.Memory, capacity.Memory[:len(capacity.Memory)-2])
		info.Total.Capacity.Storage, _ = util.StringAdd(info.Total.Capacity.Storage, capacity.Storage[:len(capacity.Storage)-2])

		info.Items[index].Status.Allocatable.Storage += "Ki"
	}

	// only trim when it contains dot
	if strings.Contains(info.Total.Allocatable.CPU, ".") {
		info.Total.Allocatable.CPU = strings.TrimRightFunc(info.Total.Allocatable.CPU, isDotOrZero)
	}
	if strings.Contains(info.Total.Allocatable.Memory, ".") {
		info.Total.Allocatable.Memory = strings.TrimRightFunc(info.Total.Allocatable.Memory, isDotOrZero) + "Ki"
	}
	if strings.Contains(info.Total.Allocatable.Storage, ".") {
		info.Total.Allocatable.Storage = strings.TrimRightFunc(info.Total.Allocatable.Storage, isDotOrZero) + "Ki"
	}
	if strings.Contains(info.Total.Allocatable.Pods, ".") {
		info.Total.Allocatable.Pods = strings.TrimRightFunc(info.Total.Allocatable.Pods, isDotOrZero)
	}

	if strings.Contains(info.Total.Capacity.CPU, ".") {
		info.Total.Capacity.CPU = strings.TrimRightFunc(info.Total.Capacity.CPU, isDotOrZero)
	}
	if strings.Contains(info.Total.Capacity.Memory, ".") {
		info.Total.Capacity.Memory = strings.TrimRightFunc(info.Total.Capacity.Memory, isDotOrZero) + "Ki"
	}
	if strings.Contains(info.Total.Capacity.Storage, ".") {
		info.Total.Capacity.Storage = strings.TrimRightFunc(info.Total.Capacity.Storage, isDotOrZero) + "Ki"
	}
	if strings.Contains(info.Total.Capacity.Pods, ".") {
		info.Total.Capacity.Pods = strings.TrimRightFunc(info.Total.Capacity.Pods, isDotOrZero)
	}

	c.JSON(http.StatusOK, info)
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
