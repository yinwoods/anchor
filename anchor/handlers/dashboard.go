// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
	"github.com/yinwoods/anchor/anchor/util"
)

var (
	username = "root"
)

func loginInitHandler(c *gin.Context) {
	inputUser, ok := c.GetPostForm("inputUser")
	if !ok {
		glog.Error("failed to login")
	}
	inputPassword, ok := c.GetPostForm("inputPassword")
	if !ok {
		glog.Error("failed to login")
	}
	match := util.CheckPass(inputPassword, userPassword)

	if inputUser == username && match {
		cookie := &http.Cookie{
			Name:    "session",
			Value:   cookieValue,
			Path:    "/",
			Expires: time.Now().AddDate(2, 0, 0),
			MaxAge:  0,
		}
		http.SetCookie(c.Writer, cookie)
	}
	indexHandler(c)
}

func indexHandler(c *gin.Context) {

	err := parseSessionCookie(c)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL, err)
		return
	}

	infos, err := cmd.DashboardList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL, err)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", infos)
}
