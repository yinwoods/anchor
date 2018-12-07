package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func nodesListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	nodes, err := cmd.NodesList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "nodes.tmpl", nodes)
}
