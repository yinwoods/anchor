// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/liman/anchor/cmd"
)

func networksHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	networks, err := cmd.NetworksList()

	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	c.HTML(http.StatusOK, "networks.tmpl", networks)
}
