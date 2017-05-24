package couchdb

type Node struct {
	host, pass, user string
	port             int
}

func (n *Node) SetHost(no string) {
	n.host = no
}

func (n *Node) SetPass(p string) {
	n.pass = p
}

func (n *Node) SetUser(u string) {
	n.user = u
}

func (n *Node) SetPort(p int) {
	n.port = p
}
