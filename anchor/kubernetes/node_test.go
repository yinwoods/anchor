package kubernetes

import (
	"testing"
)

var nodeClient NodeClient

func TestListNode(*testing.T) {
	nodeClient.ListNode()
}
