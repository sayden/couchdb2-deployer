package couchdb

import "fmt"

type FinishClusterAction struct {
	NodePassHolder
}

func (c *FinishClusterAction) GetRequest() *ClusterRequest {
	enabler := newClusterRequestPrototype()
	enabler.Action = "finish_cluster"
	enabler.Password = c.pass

	return enabler
}

func (c *FinishClusterAction) HandleResponse(r *ClusterResponse) {
	if r.Error == "" {
		fmt.Printf("\nCluster created\n\n")
	} else {
		fmt.Printf("ERROR: %#v\n", r)
	}
}
