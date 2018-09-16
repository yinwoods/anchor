// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// Event is a report of an event somewhere in the cluster.
type Event struct {
	APIVersion     string          `json:"apiVersion,omitempty"`
	Count          int64           `json:"count,omitempty"`
	FirstTimestamp string          `json:"firstTimestamp"`
	LastTimestamp  string          `json:"lastTimestamp"`
	InvolvedObject ObjectReference `json:"involvedObject"`
	Kind           string          `json:"kind,omitempty"`
	Message        string          `json:"message,omitempty"`
	Metadata       Metadata        `json:"metadata"`
	Reason         string          `json:"reason,omitempty"`
	Source         EventSource     `json:"source,omitempty"`
	Type           string          `json:"type,omitempty"`
}

// EventSource contains information for an event.
type EventSource struct {
	Component string `json:"component,omitempty"`
	Host      string `json:"host,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify
// the referred object.
type ObjectReference struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	UID        string `json:"uid"`
}

// PodList is a list of Pods.
type PodList struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ListMetadata `json:"metadata"`
	Items      []Pod        `json:"items"`
}

// PodWatchEvent for pod watch event
type PodWatchEvent struct {
	Type   string `json:"type"`
	Object Pod    `json:"object"`
}

// Pod for pod
type Pod struct {
	Kind     string   `json:"kind,omitempty"`
	Metadata Metadata `json:"metadata"`
	Spec     PodSpec  `json:"spec"`
}

// PodSpec for pod spec
type PodSpec struct {
	NodeName   string      `json:"nodeName"`
	Containers []Container `json:"containers"`
}

// Container for container
type Container struct {
	Name      string               `json:"name"`
	Resources ResourceRequirements `json:"resources"`
}

// ResourceRequirements means resource requirements
type ResourceRequirements struct {
	Limits   ResourceList `json:"limits"`
	Requests ResourceList `json:"requests"`
}

// ResourceList key-value resource
type ResourceList map[string]string

// Binding used to bind pod with node
type Binding struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Target     Target   `json:"target"`
	Metadata   Metadata `json:"metadata"`
}

// Target means target
type Target struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

// NodeList wraps node
type NodeList struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []*Node
}

// Node means node
type Node struct {
	Metadata Metadata   `json:"metadata"`
	Status   NodeStatus `json:"status"`
}

// NodeStatus wraps node status
type NodeStatus struct {
	Capacity    ResourceList `json:"capacity"`
	Allocatable ResourceList `json:"allocatable"`
}

// ListMetadata wraps meta data
type ListMetadata struct {
	ResourceVersion string `json:"resourceVersion"`
}

// Metadata means metadata
type Metadata struct {
	Name            string            `json:"name"`
	GenerateName    string            `json:"generateName"`
	ResourceVersion string            `json:"resourceVersion"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
	UID             string            `json:"uid"`
}

// ResourceUsage means resource usage
type ResourceUsage struct {
	CPU    int64
	Memory int64
	Pod    int64
}
