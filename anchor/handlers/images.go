// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func imagesListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	images, err := cmd.ImagesList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "images.tmpl", images)
}

func imageCreateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	type Input struct {
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	_, err = cmd.ImageCreate(input.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func imageDeleteHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	var input struct {
		ID string `json:"id"`
	}
	c.BindJSON(&input)

	if len(input.ID) < 7 {
		err = fmt.Errorf("image id error")
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	err = cmd.ImageDelete(input.ID[7:])
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func imageInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	id := c.Param("id")

	if len(id) < 7 {
		err = fmt.Errorf("image id error")
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	image, err := cmd.ImageGet(id[7:])
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	imageJSON, err := json.Marshal(&image)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, image)
	data = append(data, string(imageJSON))
	c.HTML(http.StatusOK, "image_info.tmpl", data)
}
