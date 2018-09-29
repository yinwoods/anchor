package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func containersListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}
	containers, err := cmd.ContainersList()
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "containers.tmpl", containers)
}

func containerInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	cid := c.Param("cid")
	container, err := cmd.ContainerGet(cid)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	containerJSON, err := json.Marshal(&container)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, container)
	data = append(data, string(containerJSON))
	c.HTML(http.StatusOK, "container_info.tmpl", data)
}

func containerCreateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	container := c.PostForm("container")
	host := c.PostForm("host")
	network := c.PostForm("network")

	var config cmd.ContainerCreateConfig

	json.Unmarshal([]byte(container), &config.ContainerConfig)
	json.Unmarshal([]byte(host), &config.HostConfig)
	json.Unmarshal([]byte(network), &config.NetworkingConfig)

	_, err = cmd.ContainerCreate(config)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/containers")
}

func containerUpdateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	type Input struct {
		Cid  string `json:"cid"`
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	var updateConfig container.UpdateConfig
	json.Unmarshal([]byte(input.Body), &updateConfig)

	_, err = cmd.ContainerUpdate(input.Cid, updateConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func containerDeleteHandler(c *gin.Context) {
	type Input struct {
		Cid string `json:"cid"`
	}
	var input Input
	c.BindJSON(&input)
	err := cmd.ContainerDelete(input.Cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
