package handlers

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gin-gonic/gin"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func powersuppliesListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}
	powersupplies, err := cmd.PowerSuppliesList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "powersupplies.tmpl", powersupplies)
}
