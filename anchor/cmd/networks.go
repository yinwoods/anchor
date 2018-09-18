// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// NetworksListOutput used to interact with template
type NetworksListOutput struct {
	Name        string `json:"Name"`
	Driver      string `json:"Driver"`
	Scope       string `json:"Scope"`
	CreatedTime string `json:"CreatedTime"`
}

// NetworksList is parses and splits networks from output
func NetworksList() ([]NetworksListOutput, error) {
	networks, err := DockerClient.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return nil, fmt.Errorf("Docker daemon is not running")
	}

	networksListOutput := []NetworksListOutput{}
	for _, network := range networks {
		networksListOutput = append(networksListOutput, NetworksListOutput{
			Name:        network.Name,
			Driver:      network.Driver,
			Scope:       network.Scope,
			CreatedTime: network.Created.Format("2006-01-02 15:04"),
		})
	}

	return networksListOutput, nil
}
