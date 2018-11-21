// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
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

	type Network struct {
		Name          string `json:"Name"`
		networkCreate types.NetworkCreate
	}
	var network Network

	json.Unmarshal([]byte(input.Body), &network)

	_, err = cmd.NetworkCreate(network.Name, network.networkCreate)
	if err != nil {
		glog.Errorf("URL=%s; Method=%s; Err=%s", c.Request.URL.Path, c.Request.Method, err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/networks")
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
		return
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
	network, _, err := cmd.NetworkGet(nid)
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
