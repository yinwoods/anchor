package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// GenerateGraph generate graph data and write into file
func GenerateGraph() (PodsContainersJSON, error) {
	result, err := PodContainersList("")
	if err != nil {
		return PodsContainersJSON{}, err
	}

	b, err := json.Marshal(result)
	if err != nil {
		return PodsContainersJSON{}, err
	}

	// pod -> container
	path, _ := os.Getwd()
	ioutil.WriteFile(path+"/public/data/flare.json", b, 0644)
	return result, nil
}
