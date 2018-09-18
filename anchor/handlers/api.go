// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
	"github.com/yinwoods/anchor/anchor/util"
)

var (
	apiKey = util.GeneratePassword(32)
)

func apiAuth(c *gin.Context) error {
	if c.Request.Method != "GET" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"ok":     "false",
			"result": "METHOD_NOT_ALLOWED",
		})
	}

	params := c.Params
	key, ok := params.Get("key")

	if !ok || len(key) < 1 {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound, gin.H{
			"ok":     "false",
			"result": "api_KEY_NOT_FOUND",
		})
	}

	if string(key[0]) != apiKey {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound, gin.H{
			"ok":     "false",
			"result": "api_KEY_INVALID",
		})
	}

	return nil
}

func apiContainer(c *gin.Context) {
	err := apiAuth(c)
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	containers, err := cmd.ContainersList()
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": containers,
	})
}

func apiImages(c *gin.Context) {
	err := apiAuth(c)
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	images, err := cmd.ImagesList()
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": images,
	})
}

func apiNetworks(c *gin.Context) {
	err := apiAuth(c)
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	networks, err := cmd.NetworksList()
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": networks,
	})
}

func apiRefgerations(c *gin.Context) {
	err := apiAuth(c)
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	refrigerations, err := cmd.RefrigerationsList()
	if err != nil {

		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": refrigerations,
	})
}

func apiPowerSupplies(c *gin.Context) {
	err := apiAuth(c)
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	powerSupplies, err := cmd.PowerSuppliesList()
	if err != nil {

		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"ok":     "true",
		"result": powerSupplies,
	})
}
