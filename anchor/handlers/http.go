// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/ratelimit"
)

const (
	apiURLPrefix = "http://localhost:8001/api/v1/namespaces/"
	apiV1Prefix  = "http://localhost:8001/apis/apps/v1/namespaces/"
	apiNodesURL  = "http://localhost:8080/api/v1/nodes"
)

var rateLimiter = ratelimit.NewBucketWithRate(1, 10)

func inc(x int) int {
	return x + 1
}

// RateLimit limit user request rate
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.V(3).Infoln("try to visit ", c.Request.URL)
		duration := rateLimiter.Take(5)
		if duration > 0 {
			glog.V(3).Infof("wait %f seconds...\n", duration.Seconds())
		}
		rateLimiter.Wait(5)
		glog.V(3).Infoln("available tokens: ", rateLimiter.Available())
		c.Next()
	}
}

// ServerRun starts server
func ServerRun() {
	r := gin.New()
	r.SetFuncMap(template.FuncMap{
		"inc": inc,
	})
	// recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	// TODO
	// r.Use(RateLimit())

	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/public", http.Dir("public"))

	r.GET("/", indexHandler)
	r.POST("/", loginInitHandler)

	r.GET("/containers", containersListHandler)
	r.POST("/containers", containerCreateHandler)
	r.DELETE("/containers", containerDeleteHandler)
	r.PUT("/containers", containerUpdateHandler)
	r.GET("/containers/:id", containerInfoHandler)

	r.GET("/images", imagesListHandler)
	r.POST("/images", imageCreateHandler)
	r.DELETE("/images", imageDeleteHandler)
	r.GET("/images/:id", imageInfoHandler)

	r.GET("/networks", networksListHandler)
	r.POST("/networks", networkCreateHandler)
	r.DELETE("/networks", networkDeletehandler)
	r.GET("/networks/:nid", networkInfoHandler)

	r.GET("/ups", upsListHandler)
	r.POST("/ups", upsCreateHandler)
	r.DELETE("/ups", upsDeleteHandler)
	r.PUT("/ups", upsUpdateHandler)
	r.GET("/ups/:id", upsInfoHandler)

	r.GET("/refs", refsListHandler)
	r.POST("/refs", refCreateHandler)
	r.DELETE("/refs", refDeleteHandler)
	r.PUT("/refs", refUpdateHandler)
	r.GET("/refs/:id", refInfoHandler)
	r.GET("/pods", podsListHandler)
	r.POST("/pods", podCreateHandler)
	r.DELETE("/pods", podDeleteHandler)
	r.PUT("/pods", podUpdateHandler)
	r.GET("/pods/:namespace/:name", podInfoHandler)

	r.GET("/nodes", nodesListHandler)
	r.GET("/nodes/:name", nodeInfoHandler)

	r.GET("/services", servicesListHandler)
	r.POST("/services", serviceCreateHandler)
	r.DELETE("/services", serviceDeleteHandler)
	r.PUT("/services", serviceUpdateHandler)
	r.GET("/services/:namespace/:name", serviceInfoHandler)

	r.GET("/deployments", deploymentsListHandler)
	r.POST("/deployments", deploymentCreateHandler)
	r.DELETE("/deployments", deploymentDeleteHandler)
	r.PUT("/deployments", deploymentUpdateHandler)
	r.GET("/deployments/:namespace/:name", deploymentInfoHandler)

	r.GET("/settings", settingsGetHandler)
	r.POST("/settings", settingsUpdateHandler)

	r.GET("/login", loginHandler)
	r.GET("/logout", logoutHandler)
	r.GET("/install", installHandler)

	r.GET("/graph", graphHandler)

	r.GET("/api/tokens", apiTokensHandler)
	r.GET("/api/sysinfo", apiSysInfoHandler)
	r.GET("/api/graph", apiGraphInfo)
	r.GET("/api/containers/:id/", apiContainerInfo)
	r.GET("/api/images/:id/", apiImageInfo)
	r.GET("/api/networks/:id/", apiNetworkInfo)
	r.GET("/api/ups/:id", apiUPSInfo)
	r.GET("/api/refs/:id", apiREFInfo)
	r.GET("/api/pods/:namespace/:name", apiPodInfo)
	r.GET("/api/services/:namespace/:name", apiServiceInfo)
	r.GET("/api/deployments/:namespace/:name", apiDeploymentInfo)

	r.Run(":8090")
}
