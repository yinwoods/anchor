package handlers

import (
	"net/http"

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
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "containers.tmpl", containers)
}
