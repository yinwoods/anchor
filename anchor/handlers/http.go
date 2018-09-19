// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"html/template"

	"github.com/gin-gonic/gin"
)

func inc(x int) int {
	return x + 1
}

// ServerRun starts server
func ServerRun() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"inc": inc,
	})

	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/public", http.Dir("public"))

	r.GET("/", indexHandler)
	r.POST("/", loginInitHandler)

	r.GET("/containers", containersListHandler)
	r.GET("/images", imagesListHandler)
	r.GET("/networks", networksListHandler)
	r.GET("/refrigerations", refrigerationsListHandler)
	r.GET("/powersupplies", powersuppliesListHandler)

	r.GET("/pods", podsListHandler)
	r.GET("/pods/:namespace/:name", podInfoHandler)

	// TODO
	r.POST("/pods", podsUpdateHandler)

	r.GET("/nodes", nodesListHandler)
	r.GET("/nodes/:name", nodeInfoHandler)

	r.GET("/services", servicesListHandler)
	r.GET("/services/:namespace/:name", serviceInfoHandler)

	r.GET("/deployments", deploymentsListHandler)
	r.GET("/deployments/:namespace/:name", deploymentInfoHandler)

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
