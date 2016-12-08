package couchdb

func JoinAllClusterAction(pass string, nodes []string) {
	template := &Template{}
	coord := nodes[0]

	for i := 1; i < len(nodes); i++ {
		node := nodes[i]

		template.SetTemplate(&EnableClusterAction{
			NodePassHolder: NodePassHolder{
				pass: pass,
				node: node,
			},
		})
		template.Do(coord)

		template.SetTemplate(&AddNodeAction{
			NodePassHolder: NodePassHolder{
				pass: pass,
				node: node,
			},
		})
		template.Do(coord)
	}

	template.SetTemplate(&FinishClusterAction{
		NodePassHolder: NodePassHolder{
			pass: pass,
		},
	})
	template.Do(coord)

	CheckCluster(coord, pass)
}
