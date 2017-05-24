package couchdb

import "fmt"

type FinishAction struct {
	Node
}

func (c *FinishAction) GetRequest() *ClusterRequest {
	enabler := newClusterRequestPrototype(c.Node)
	enabler.Action = "finish_cluster"

	return enabler
}

func (c *FinishAction) HandleResponse(r *ClusterResponse) {
	if r.Error == "" {
		fmt.Printf("\nCluster created\n\n")
	} else {
		fmt.Printf("ERROR: %#v\n", r)
	}
}
