package couchdb

import "fmt"

type EnableClusterAction struct {
	NodePassHolder
}

func (e *EnableClusterAction) GetRequest() *ClusterRequest{
	enabler := newClusterRequestPrototype()
	enabler.Action = "enable_cluster"
	enabler.BindAddress = "0.0.0.0"
	enabler.RemoteNode = e.node
	enabler.RemoteCurrentUser = "admin"
	enabler.RemoteCurrentPassword = e.pass
	enabler.Password = e.pass

	return enabler
}

func (e *EnableClusterAction) HandleResponse(ok *ClusterResponse) {
	if ok.Error != "" {
		fmt.Printf("ERROR: %#v\n", ok)
	} else {
		fmt.Printf("OK ")
	}
}
