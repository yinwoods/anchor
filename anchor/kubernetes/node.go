package kubernetes

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// NodeClient wraps client to operate node apis
type NodeClient struct {
	corev1.NodeInterface
}

// GetNode used to get node by node name
func (client NodeClient) GetNode(nodename string) (*v1.Node, error) {
	return client.Get(nodename, metav1.GetOptions{})
}

// ListNode used to list node
func (client NodeClient) ListNode() {

	list, err := client.List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("List node failed: %v", err)
	}
	for _, node := range list.Items {
		glog.V(2).Infof(" * %s \n", node.Name)
	}
}
