package couchdb

type ClusterActionPerformer interface {
	GetRequest() *ClusterRequest
	HandleResponse(*ClusterResponse)
	SetHost(string)
	SetPass(string)
	SetUser(string)
	SetPort(int)
}