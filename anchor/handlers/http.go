// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"io/ioutil"
	"net/http"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const (
	apiURLPrefix = "http://localhost:8001/api/v1/namespaces/"
	apiV1Prefix  = "http://localhost:8001/apis/apps/v1/namespaces/"
)

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		glog.Errorf("get url: %s error: %s", url, err.Error())
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

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
	r.POST("/containers", containerCreateHandler)
	r.DELETE("/containers", containerDeleteHandler)
	r.PUT("/containers", containerUpdateHandler)
	r.GET("/containers/:cid", containerInfoHandler)

	r.GET("/images", imagesListHandler)
	r.POST("/images", imageCreateHandler)
	r.DELETE("/images", imageDeleteHandler)
	r.GET("/images/:mid", imageInfoHandler)

	r.GET("/networks", networksListHandler)
	r.POST("/networks", networkCreateHandler)
	r.DELETE("/networks", networkDeletehandler)
	r.GET("/networks/:nid", networkInfoHandler)

	r.GET("/refrigerations", refrigerationsListHandler)
	r.GET("/powersupplies", powersuppliesListHandler)

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

	r.GET("/api/graph", apiGraphInfo)
	r.GET("/api/containers/:cid/", apiContainerUpdateConfigInfo)
	r.GET("/api/images/:mid/", apiImageInfo)
	r.GET("/api/pods/:namespace/:name", apiPodInfo)
	r.GET("/api/services/:namespace/:name", apiServiceInfo)
	r.GET("/api/deployments/:namespace/:name", apiDeploymentInfo)

	r.Run(":8090")
}
