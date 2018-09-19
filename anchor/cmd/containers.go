// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainersListOutput used to interact with template
type ContainersListOutput struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	CreatedTime string `json:"CreatedTime"`
	State       string `json:"State"`
}

// ContainersList used to list containers
func ContainersList() ([]ContainersListOutput, error) {
	containers, err := DockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
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
