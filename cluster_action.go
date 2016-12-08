package couchdb

type ClusterActionPerformer interface {
	GetRequest() *ClusterRequest
	HandleResponse(*ClusterResponse)
	SetNode(string)
	SetPass(string)
}

type NodePassHolder struct {
	node, pass string
}

func (n *NodePassHolder) SetNode(no string) {
	n.node = no
}

func (n *NodePassHolder) SetPass(p string) {
	n.pass = p
}
