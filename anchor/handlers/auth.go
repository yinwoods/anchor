// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/util"
)

var (
	cookieValue  = util.GeneratePassword(140)
	userPassword = util.ReadPassword()
)

func parseSessionCookie(c *gin.Context) error {
	if userPassword == "" {
		c.Redirect(http.StatusFound, "/install")
		return fmt.Errorf("System not Installed")
	}

	cookie, err := c.Request.Cookie("session")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "session",
			Value: "",
			Path:  "/",
		}
		http.SetCookie(c.Writer, cookie)
		c.Redirect(http.StatusFound, "/login")
		return fmt.Errorf("No Cookie")
	}

	if cookie.Value != cookieValue {
		c.Redirect(http.StatusFound, "/login")
		return fmt.Errorf("Invalidated Cookie Value")
	}

	return nil
}

func loginHandler(c *gin.Context) {
	if userPassword == "" {
		c.Redirect(http.StatusFound, "/install")
	}

	if c.Request.URL.Path != "/login" {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path)
		c.Redirect(http.StatusFound, "/")
		return
	}

	cookie, err := c.Cookie("session")
	if err != nil {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path)
		c.Redirect(http.StatusFound, "/")
		return
	}

	if cookie == cookieValue {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path)
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func logoutHandler(c *gin.Context) {
	if c.Request.URL.Path != "/logout" {
		glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path)
		c.Redirect(http.StatusFound, "/")
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: "",
		Path:  "/",
	}

	http.SetCookie(c.Writer, cookie)
	glog.V(2).Infoln(c.Request.Method, c.Request.URL.Path)
	c.Redirect(http.StatusFound, "/")
}

func installHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		if userPassword == "" {
			inputPassword := c.Request.FormValue("inputPassword")
			hash, err := util.HashPasswordAndSave(inputPassword)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			userPassword = hash
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	if userPassword != "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "install.tmpl", nil)
}
