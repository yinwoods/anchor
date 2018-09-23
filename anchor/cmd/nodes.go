package cmd

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodesListOutput used to interact with template
type NodesListOutput struct {
	Name              string           `json:"Name"`
	Addresses         []v1.NodeAddress `json:"Addresses"`
	CreationTimestamp string           `json:"CreationTimestamp"`
}

// GetNode used to get node by node name
func GetNode(nodename string) (*v1.Node, error) {
	client := GetNodeClient()
	return client.Get(nodename, metav1.GetOptions{})
}

// NodesList used to list node
func NodesList() ([]NodesListOutput, error) {
	client := GetNodeClient()
	nodesListOutput := []NodesListOutput{}
	nodes, err := client.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("List nodes failed: %v", err)
	}
	for _, node := range nodes.Items {
		nodesListOutput = append(nodesListOutput, NodesListOutput{
			Name:              node.Name,
			Addresses:         node.Status.Addresses,
			CreationTimestamp: node.CreationTimestamp.Format("2006-01-02 15:04"),
		})
	}
	return nodesListOutput, nil
}
