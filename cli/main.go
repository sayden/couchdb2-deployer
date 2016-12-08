package main

import (
	"fmt"
	"os"

	"github.com/micro/cli"
	"github.com/sayden/couchdb2_orchestrator"
)

func main() {
	app := cli.NewApp()
	var password string

	getCLIHandler := func(t couchdb.ClusterActionPerformer) func(c *cli.Context) {
		return func(c *cli.Context) {
			if c.NArg() <= 1 {
				fmt.Println("Not enough arguments, you need to pass at least two" +
					"nodes: one to use as coordinator and one to join")
				return
			}

			t.SetNode(c.Args().Get(1))
			t.SetPass(password)

			template := couchdb.Template{
				Template: t,
			}

			template.Do(c.Args().First())
		}
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "password, p",
			Usage:       "Password to access each node. Default username is 'admin'",
			EnvVar:      "COUCHDB_PASSWORD",
			Value:       "password",
			Destination: &password,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "enable",
			Aliases: []string{"e"},
			Usage:   "Enable clustering in a specific node",
			Action:  getCLIHandler(&couchdb.EnableClusterAction{}),
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Adds a node to a cluster",
			Action:  getCLIHandler(&couchdb.AddNodeAction{}),
		},
		{
			Name:    "finish",
			Aliases: []string{"a"},
			Usage:   "Adds a node to a cluster",
			Action:  func (c *cli.Context){
				if c.NArg() < 1 {
					fmt.Println("Not enough arguments, you need to pass at least two" +
						"nodes: one to use as coordinator and one to join")
					return
				}

				t := &couchdb.FinishClusterAction{}

				t.SetNode(c.Args().Get(1))
				t.SetPass(password)

				template := couchdb.Template{
					Template: t,
				}

				template.Do(c.Args().First())
			},
		},
		{
			Name:    "join_all",
			Aliases: []string{"j"},
			Usage:   "Joins all specified nodes to a cluster",
			Action: func(c *cli.Context) {
				if c.NArg() <= 1 {
					fmt.Println("Not enough arguments, you need to pass at least two" +
						"nodes: one to use as coordinator and one to join")
					return
				}

				couchdb.JoinAllClusterAction(password, c.Args())
			},
		},
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Check nodes in cluster",
			Action: func(c *cli.Context) {
				if c.NArg() < 1 {
					fmt.Println("Not enough arguments, you need to pass at least two" +
						"nodes: one to use as coordinator and one to join")
					return
				}
				couchdb.CheckCluster(c.Args().First(), password)
			},
		},
	}

	app.Name = "CouchDB 2 Orchestrator"
	app.Usage = "Common usage would be to pass a '-c IP' for a coordination node then 'join_all' as an action and a list of IP's to conform cluster"

	app.Run(os.Args)
}
