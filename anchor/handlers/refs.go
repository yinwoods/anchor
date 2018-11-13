package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func refsListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}
	refs, err := cmd.REFsList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "refs.tmpl", refs)
}

func refCreateHandler(c *gin.Context) {
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

	var item cmd.REFsListOutput
	err = json.Unmarshal([]byte(input.Body), &item)
	_, err = cmd.REFsCreate(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/refs")
}

func refUpdateHandler(c *gin.Context) {
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

	var item cmd.REFsListOutput
	err = json.Unmarshal([]byte(input.Body), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = cmd.REFUpdate(item)
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

func refDeleteHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	var input struct {
		ID string `json:"id"`
	}
	c.BindJSON(&input)

	err = cmd.REFDelete(input.ID)
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

func refInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		return
	}

	id := c.Param("id")
	ref, err := cmd.REFGet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	refJSON, err := json.Marshal(&ref)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, ref)
	data = append(data, string(refJSON))
	c.HTML(http.StatusOK, "ref_info.tmpl", data)
}
