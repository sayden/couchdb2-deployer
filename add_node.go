package couchdb

import "fmt"

type AddNodeAction struct {
	NodePassHolder
}

func (a *AddNodeAction) GetRequest() *ClusterRequest {
	enabler := newClusterRequestPrototype()
	enabler.Action = "add_node"
	enabler.Host = a.node
	enabler.Password = a.pass

	return enabler
}

func (a *AddNodeAction) HandleResponse(resp *ClusterResponse) {
	if resp.Error == "" {
		fmt.Printf("Node '%s' joined successfully\n", a.node)
	} else {
		fmt.Printf("%#v\n", resp)
	}
}
