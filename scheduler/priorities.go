package main

import "github.com/golang/glog"

// MaxPriority means highest score which node can get
const MaxPriority = 10

func balancedResourceScore(requested, allocatable ResourceUsage) float64 {
	cFraction := fractionOfCapacity(requested.CPU, allocatable.CPU)
	mFraction := fractionOfCapacity(requested.Memory, allocatable.Memory)
	pFraction := fractionOfCapacity(requested.Pod, allocatable.Pod)

	if cFraction >= 1 || mFraction >= 1 || pFraction >= 1 {
		return 0
	}

	return getBalancedResourceScore(cFraction, mFraction, pFraction)
}

func leastRequestedScore(requested, allocatable ResourceUsage) float64 {
	cRatio := getLeastRequestedScore(requested.CPU, allocatable.CPU)
	mRatio := getLeastRequestedScore(requested.Memory, allocatable.Memory)
	pRatio := getLeastRequestedScore(requested.Pod, allocatable.Pod)

	return (cRatio + mRatio + pRatio) / 3
}

func priorities(pod *Pod, nodes []*Node) (*Node, error) {

	var bestNode *Node
	nodeScore := make(map[*Node]float64)

	// 获取所有节点
	nodeList, err := getNodes()
	if err != nil {
		glog.Error("failed to get nodes")
		glog.Fatal(err)
	}

	// 获取所有pod
	podList, err := getPods()
	if err != nil {
		glog.Error("failed to get pods")
		glog.Fatal(err)
	}

	requested := requestedResource(pod)
	used := usedResource(nodeList, podList)

	for _, node := range nodeList.Items {
		nodeScore[node] = 0
	}

	for _, node := range nodeList.Items {

		allocatable := allocatableResource(node, used)
		nodeScore[node] += balancedResourceScore(requested, allocatable)
		nodeScore[node] += leastRequestedScore(requested, allocatable)
		nodeScore[node] /= 2
	}

	printNodeScores(nodeScore)

	var maxScore float64
	for node, score := range nodeScore {
		if score > maxScore {
			maxScore = score
			bestNode = node
		}
	}
	return bestNode, nil
}
