package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func upsListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}
	ups, err := cmd.UPSList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "ups.tmpl", ups)
}

func upsCreateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	type Input struct {
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	var item cmd.UPSListOutput
	err = json.Unmarshal([]byte(input.Body), &item)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	_, err = cmd.UPSCreate(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/ups")
}

func upsUpdateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	type Input struct {
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	var item cmd.UPSListOutput
	err = json.Unmarshal([]byte(input.Body), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = cmd.UPSUpdate(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func upsDeleteHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	var input struct {
		ID string `json:"id"`
	}
	c.BindJSON(&input)

	err = cmd.UPSDelete(input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func upsInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	id := c.Param("id")
	ups, err := cmd.UPSGet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	upsJSON, err := json.Marshal(&ups)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, ups)
	data = append(data, string(upsJSON))
	c.HTML(http.StatusOK, "ups_info.tmpl", data)
}
