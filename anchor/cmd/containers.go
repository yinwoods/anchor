// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/yinwoods/anchor/anchor/util"
	"golang.org/x/net/context"
)

const (
	dockerContainerURL = "http://localhost:8089/api/container"
)

// ContainerCreateConfig wraps container config 、host config and networking config
type ContainerCreateConfig struct {
	ContainerConfig  container.Config
	HostConfig       container.HostConfig
	NetworkingConfig network.NetworkingConfig
}

// ContainersListOutput used to interact with template
type ContainersListOutput struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	CreatedTime string `json:"CreatedTime"`
	State       string `json:"State"`
}

// ContainersList used to list containers
func ContainersList() ([]ContainersListOutput, error) {

	resp, _ := util.HTTPGet(dockerContainerURL)
	var containers []types.Container
	json.Unmarshal(resp, &containers)
	containersListOutput := []ContainersListOutput{}
	for _, container := range containers {
		// 过滤包含kube字段的容器
		if strings.Contains(container.Names[0], "kube") {
			continue
		}
		containerName := container.Names[0]
		// 截取k8s_
		if strings.HasPrefix(container.Names[0], "/k8s_") {
			containerName = container.Names[0][5:]
		}
		containersListOutput = append(containersListOutput, ContainersListOutput{
			ID:          container.ID,
			Name:        containerName,
			CreatedTime: time.Unix(container.Created, 0).Format("2006-01-02 15:04"),
			State:       container.State,
		})
	}
	return containersListOutput, nil
}

// ContainerGet returns docker inspect information
func ContainerGet(id string) (types.ContainerJSON, error) {
	resp, _ := util.HTTPGet(fmt.Sprintf("%s/%s/json", dockerContainerURL, id))
	var container types.ContainerJSON
	json.Unmarshal(resp, &container)
	return container, nil
}

// ContainerCreate create a container
func ContainerCreate(config ContainerCreateConfig) (container.ContainerCreateCreatedBody, error) {
	return DockerClient.ContainerCreate(context.Background(), &config.ContainerConfig, &config.HostConfig, &config.NetworkingConfig, "")
}

// ContainerUpdate updates an container
func ContainerUpdate(id, config string) error {
	_, err := util.HTTPPost(fmt.Sprintf("%s/%s", dockerContainerURL, id), config)
	return err
}

// ContainerDelete delete a container
func ContainerDelete(id string) error {
	_, err := util.HTTPDelete(fmt.Sprintf("%s/%s", dockerContainerURL, id))
	return err
}
