package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/liman/anchor/cmd"
)

func refrigerationsHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}
	refrigerations, err := cmd.RefrigerationsList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "refrigerations.tmpl", refrigerations)
}
