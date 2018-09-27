package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func refrigerationsListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}
	refrigerations, err := cmd.RefrigerationsList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "refrigerations.tmpl", refrigerations)
}
