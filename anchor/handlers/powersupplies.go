package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/liman/anchor/cmd"
)

func powersuppliesHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}
	powersupplies, err := cmd.PowerSuppliesList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "powersupplies.tmpl", powersupplies)
}
