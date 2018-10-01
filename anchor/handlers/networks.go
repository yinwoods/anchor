// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func networksListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	networks, err := cmd.NetworksList()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "networks.tmpl", networks)
}

func networkCreateHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	type Input struct {
		Body string `json:"body"`
	}
	var input Input
	c.BindJSON(&input)

	var network struct {
		Name   string `json:"name"`
		Driver string `json:"driver"`
	}
	json.Unmarshal([]byte(input.Body), &network)

	_, err = cmd.NetworkCreate(network.Name, network.Driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func networkDeletehandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	var input struct {
		ID string `json:"id"`
	}
	c.BindJSON(&input)
	err = cmd.NetworkDelete(input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func networkInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	nid := c.Param("nid")
	network, err := cmd.NetworkGet(nid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	networkJSON, err := json.Marshal(&network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []interface{}
	data = append(data, network)
	data = append(data, string(networkJSON))
	c.HTML(http.StatusOK, "network_info.tmpl", data)
}
