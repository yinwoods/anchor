// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"golang.org/x/net/context"
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
	containers, err := DockerClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("Docker daemon is not running")
	}
	containersListOutput := []ContainersListOutput{}
	for _, container := range containers {
		containersListOutput = append(containersListOutput, ContainersListOutput{
			ID:          container.ID,
			Name:        container.Names[0],
			CreatedTime: time.Unix(container.Created, 0).Format("2006-01-02 15:04"),
			State:       container.State,
		})
	}
	return containersListOutput, nil
}

// ContainerGet returns docker inspect information
func ContainerGet(cid string) (types.ContainerJSON, error) {
	return DockerClient.ContainerInspect(context.Background(), cid)
}

// ContainerCreate create a container
func ContainerCreate(config ContainerCreateConfig) (container.ContainerCreateCreatedBody, error) {
	return DockerClient.ContainerCreate(context.Background(), &config.ContainerConfig, &config.HostConfig, &config.NetworkingConfig, "")
}

// ContainerUpdate updates an container
func ContainerUpdate(cid string, config container.UpdateConfig) (container.ContainerUpdateOKBody, error) {
	return DockerClient.ContainerUpdate(context.Background(), cid, config)
}

// ContainerDelete delete a container
func ContainerDelete(cid string) error {
	return DockerClient.ContainerRemove(context.Background(), cid, types.ContainerRemoveOptions{})
}
