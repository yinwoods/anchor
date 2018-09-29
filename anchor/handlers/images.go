// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
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

func imageDeleteHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	type Input struct {
		ImageID string `json:"mid"`
	}
	var input Input
	c.BindJSON(&input)

	glog.V(3).Infoln("imageID: ", input.ImageID)

	_, err = cmd.ImageDelete(input.ImageID)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
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

	mid := c.Param("mid")
	image, imageJSON, err := cmd.ImageGet(mid)
	if err != nil {
		glog.Error(c.Request.URL.Path, c.Request.Method, err)
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
