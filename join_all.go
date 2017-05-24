package couchdb

func JoinAllClusterAction(user, pass, coord string, nodes []string, port int) error {
	template := &CouchDbHTTP{}
	n := Node{
		user: user,
		pass: pass,
		port: port,
	}

	for _, node := range nodes {
		n.host = node

		template.SetActionPerformer(&EnableHostAction{n})
		_ = template.Do(coord)

		template.SetActionPerformer(&AddNodeAction{n})

		err := template.Do(coord)
		if err != nil {
			return err
		}
	}

	CheckCluster(coord, user, pass, port)

	return nil
}
