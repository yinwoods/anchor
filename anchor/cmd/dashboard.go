// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

// DashboardList is parses and splits docker stat values from output
func DashboardList() ([]interface{}, error) {
	infos, err := Client.Info(context.Background())

	if err != nil {
		return nil, fmt.Errorf("Docker daemon is not running %s", err)
	}

	var dashboard []interface{}

	dashboard = append(dashboard, infos.Containers)
	dashboard = append(dashboard, infos.Name)
	dashboard = append(dashboard, infos.ServerVersion)
	dashboard = append(dashboard, infos.NCPU)
	dashboard = append(dashboard, infos.MemTotal)

	images, err := ImagesList()
	if err != nil {
		return nil, fmt.Errorf("Docker daemon is not running %s", err)
	}

	networks, err := NetworksList()
	if err != nil {
		return nil, fmt.Errorf("Docker daemon is not running %s", err)
	}

	dashboard = append(dashboard, strconv.Itoa(len(images)))
	dashboard = append(dashboard, strconv.Itoa(len(networks)))

	refrigerations, err := RefrigerationsList()
	if err != nil {
		return nil, fmt.Errorf("List refrigerations error")
	}
	dashboard = append(dashboard, len(refrigerations))

	powerSupplies, err := PowerSuppliesList()
	if err != nil {
		return nil, fmt.Errorf("List powerSupplies error")
	}
	dashboard = append(dashboard, len(powerSupplies))

	intMemory := infos.MemTotal
	floatMemory := float64(intMemory)
	GibMemory := ((floatMemory / 1024) / 1024) / 1024
	dashboard[4] = strconv.FormatFloat(GibMemory, 'f', 2, 64)

	dashboard[1] = strings.Title(dashboard[1].(string))

	return dashboard, nil
}