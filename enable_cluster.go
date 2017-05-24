package couchdb

import "fmt"

type EnableHostAction struct {
	Node
}

func (e *EnableHostAction) GetRequest() *ClusterRequest{
	enabler := newClusterRequestPrototype(e.Node)
	enabler.Action = "enable_cluster"
	enabler.BindAddress = "0.0.0.0"
	enabler.RemoteNode = e.host

	return enabler
}

func (e *EnableHostAction) HandleResponse(ok *ClusterResponse) {
	if ok.Error != "" {
		fmt.Printf("ERROR: %#v\n %#v", ok, e.Node)
	} else {
		fmt.Printf("OK ")
	}
}
