package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/cmd"
)

func nodesListHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	nodes, err := cmd.K8SClient.NodeClient.NodesList()
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	c.HTML(http.StatusOK, "nodes.tmpl", nodes)
}

func nodeInfoHandler(c *gin.Context) {
	err := parseSessionCookie(c)
	if err != nil {
		return
	}

	nodeName := c.Param("name")
	node, err := cmd.K8SClient.GetNode(nodeName)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}
	nodeJSON, err := json.Marshal(&node)
	if err != nil {
		glog.Error(c.Request.Method, c.Request.URL.Path, err)
		return
	}

	var data []interface{}
	data = append(data, node)
	data = append(data, string(nodeJSON))
	c.HTML(http.StatusOK, "node_info.tmpl", data)
}
