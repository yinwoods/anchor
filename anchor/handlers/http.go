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
	r.GET("/containers/:cid", containerInfoHandler)

	r.GET("/images", imagesListHandler)
	r.GET("/images/:mid", imageInfoHandler)

	r.GET("/networks", networksListHandler)
	r.GET("/refrigerations", refrigerationsListHandler)
	r.GET("/powersupplies", powersuppliesListHandler)

	r.GET("/pods", podsListHandler)
	r.POST("/pods", podCreateHandler)
	r.DELETE("/pods", podDeleteHandler)
	r.GET("/pods/:namespace/:name", podInfoHandler)
	r.PUT("/pods", podsUpdateHandler)

	r.GET("/nodes", nodesListHandler)
	r.GET("/nodes/:name", nodeInfoHandler)

	r.GET("/services", servicesListHandler)
	r.POST("/services", serviceCreateHandler)
	r.DELETE("/services", serviceDeleteHandler)
	r.GET("/services/:namespace/:name", serviceInfoHandler)

	r.GET("/deployments", deploymentsListHandler)
	r.POST("/deployments", deploymentCreateHandler)
	r.DELETE("/deployments", deploymentsDeleteHandler)
	r.GET("/deployments/:namespace/:name", deploymentInfoHandler)

	r.GET("/settings", settingsGetHandler)
	r.POST("/settings", settingsUpdateHandler)

	r.GET("/login", loginHandler)
	r.GET("/logout", logoutHandler)
	r.GET("/install", installHandler)

	r.GET("/api/containers", apiContainer)
	r.GET("/api/pods/:namespace/:name", apiPodInfo)
	r.GET("/api/images", apiImages)
	r.GET("/api/networks", apiNetworks)
	r.GET("/api/refrigerations", apiRefgerations)
	r.GET("/api/powersupplies", apiPowerSupplies)

	r.Run(":8090")
}
