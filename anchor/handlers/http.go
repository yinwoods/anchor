// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/gin-gonic/gin"
)

// ServerRun starts server
func ServerRun() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/public/", "public")

	r.GET("/", indexHandler)
	r.POST("/", loginInitHandler)
	r.GET("/containers", containersHandler)
	r.GET("/images", imagesHandler)
	r.GET("/networks", networksHandler)
	r.GET("/refrigerations", refrigerationsHandler)
	r.GET("/powersupplies", powersuppliesHandler)
	r.GET("/settings", settingsGetHandler)
	r.POST("/settings", settingsUpdateHandler)

	r.GET("/login", loginHandler)
	r.GET("/logout", logoutHandler)
	r.GET("/install", installHandler)

	r.GET("/api/containers", apiContainer)
	r.GET("/api/images", apiImages)
	r.GET("/api/networks", apiNetworks)
	r.GET("/api/refrigerations", apiRefgerations)
	r.GET("/api/powersupplies", apiPowerSupplies)

	r.Run(":8090")
}
