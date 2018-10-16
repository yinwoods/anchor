package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yinwoods/anchor/anchor/util"
)

func graphHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}
	_, err = util.GenerateGraph()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.HTML(http.StatusOK, "graph.tmpl", nil)
}
