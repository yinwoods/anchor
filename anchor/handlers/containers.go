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
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
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

	id := c.Param("id")
	container, err := cmd.ContainerGet(id)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	containerJSON, err := json.Marshal(&container)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
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
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
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
		ID   string `json:"id"`
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	var updateConfig container.UpdateConfig
	json.Unmarshal([]byte(input.Body), &updateConfig)

	_, err = cmd.ContainerUpdate(input.ID, updateConfig)
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
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	var input struct {
		ID string `json:"id"`
	}
	c.BindJSON(&input)
	err = cmd.ContainerDelete(input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
