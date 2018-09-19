package cmd

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// NodeClient wraps client to operate node apis
type NodeClient struct {
	corev1.NodeInterface
}

// NodesListOutput used to interact with template
type NodesListOutput struct {
	Name              string           `json:"Name"`
	Addresses         []v1.NodeAddress `json:"Addresses"`
	CreationTimestamp string           `json:"CreationTimestamp"`
}

// GetNode used to get node by node name
func (client NodeClient) GetNode(nodename string) (*v1.Node, error) {
	return client.Get(nodename, metav1.GetOptions{})
}

// NodesList used to list node
func (client NodeClient) NodesList() ([]NodesListOutput, error) {
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
