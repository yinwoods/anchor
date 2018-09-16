// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/liman/anchor/util"
)

func settingsUpdateHandler(c *gin.Context) {
	pass := c.Request.FormValue("cpass")

	match := util.CheckPass(pass, userPassword)
	if !match {
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	nPass := c.Request.FormValue("npass")
	cNPass := c.Request.FormValue("cnpass")

	if nPass != cNPass {
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	match = util.CheckPass(nPass, userPassword)

	if match {
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	bNPass, err := util.HashPasswordAndSave(nPass)
	if err != nil {
		glog.Error(err)
		return
	}

	userPassword = string(bNPass)
	c.Redirect(http.StatusFound, "/logout")
	settingsGetHandler(c)
}

func settingsGetHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	version, err := util.Version()
	if err != nil {
		glog.Error(err)
	}

	var data []interface{}
	data = append(data, version)
	data = append(data, apiKey)

	c.HTML(http.StatusOK, "settings.tmpl", data)
}
