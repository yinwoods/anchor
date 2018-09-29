package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
	"github.com/yinwoods/anchor/anchor/util"
)

var (
	apiKey = util.GeneratePassword(32)
)

func apiContainer(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	containers, err := cmd.ContainersList()
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": containers,
	})
}

func apiImages(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	images, err := cmd.ImagesList()
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": images,
	})
}

func apiNetworks(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	networks, err := cmd.NetworksList()
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": networks,
	})
}

func apiRefgerations(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	refrigerations, err := cmd.RefrigerationsList()
	if err != nil {

		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": refrigerations,
	})
}

func apiPowerSupplies(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	powerSupplies, err := cmd.PowerSuppliesList()
	if err != nil {

		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": powerSupplies,
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
