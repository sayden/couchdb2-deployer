package couchdb

// ClusterRequest serves to make all neccesary communications to join a cluster
type ClusterRequest struct {
	Action                string `json:"action"`
	BindAddress           string `json:"bind_address,omitempty"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	Port                  int    `json:"port"`
	RemoteNode            string `json:"remote_node,omitempty"`
	RemoteCurrentUser     string `json:"remote_current_user,omitempty"`
	RemoteCurrentPassword string `json:"remote_current_password,omitempty"`
	Host                  string `json:"host,omitempty"`
}

// ClusterResponse is returned when calling "_cluster_setup" endpoint
type ClusterResponse struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// MembershipResponse is returned by a CouchDB 2 when asking for the memberships of a cluster
type MembershipResponse struct {
	AllNodes     []string `json:"all_nodes"`
	ClusterNodes []string `json:"cluster_nodes"`
}

func newClusterRequestPrototype() *ClusterRequest {
	return &ClusterRequest{
		Username: "admin",
		Port:     5984,
	}
}
