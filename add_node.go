package couchdb

import "fmt"

type AddNodeAction struct {
	Node
}

func (a *AddNodeAction) GetRequest() *ClusterRequest {
	enabler := newClusterRequestPrototype(a.Node)
	enabler.Action = "add_node"

	return enabler
}

func (a *AddNodeAction) HandleResponse(resp *ClusterResponse) {
	if resp.Error == "" {
		fmt.Printf("Node '%s' joined successfully\n", a.host)
	} else {
		fmt.Printf("ERROR Adding node: %#v\n%#v\n", resp, a.Node)
	}
}
