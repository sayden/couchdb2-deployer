package main

import (
	"fmt"
	"os"

	"github.com/sayden/couchdb2_deployer"
    "github.com/urfave/cli"
)

var password, admin string
var port int
var sendFinishCluster bool

func main() {
	app := cli.NewApp()

	getSingleActionCLIHandler := func(t couchdb.ClusterActionPerformer) func(c *cli.Context) {
		return func(c *cli.Context) {
			if c.NArg() < 1 {
				fmt.Println("Not enough arguments, you need to pass at least two" +
					"nodes: one to use as coordinator and one to join")
				return
			}

			t.SetHost(c.Args().Get(0))
			t.SetPass(password)
			t.SetUser(admin)
			t.SetPort(port)

			var template couchdb.TemplateExecutor = &couchdb.CouchDbHTTP{
				ActionPerformer: t,
			}

			if err := template.Do(c.String("coord")); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "password, p",
			Usage:       "Password to access each node. Default username is 'admin'",
			EnvVar:      "COUCHDB_PASSWORD",
			Destination: &password,
		},
		cli.StringFlag{
			Name:        "user, u",
			Usage:       "User to access each node. Default username is 'admin'",
			EnvVar:      "COUCHDB_USER",
			Destination: &admin,
		},
		cli.IntFlag{
			Name:        "http",
			Usage:       "HTTP CouchDB comm port",
			Value:       5984,
			Destination: &port,
		},
		cli.BoolFlag{
			Name:        "finish, f",
			Usage:       "Executes finish function (or not)",
			Destination: &sendFinishCluster,
		},
		cli.StringFlag{
			Name:   "coord, c",
			Usage:  "Coordination node",
			EnvVar: "COORD",
			Value:  "couchdb-0",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "enable",
			Aliases: []string{"e"},
			Usage:   "Enable clustering in a specific node",
			Action:  getSingleActionCLIHandler(&couchdb.EnableHostAction{}),
		},
		{
			Name:    "remove",
			Aliases: []string{"a"},
			Usage:   "Removes a node from cluster",
			Action: func(c *cli.Context) {
				if c.NArg() < 1 {
					fmt.Println("Not enough arguments")
					os.Exit(1)
				}

				n := couchdb.Node{}
				n.SetPort(port)
				n.SetHost(c.Args().Get(1))
				n.SetPass(password)
				n.SetUser(admin)

				if err := couchdb.Remove(c.String("coord"), n); err != nil {
					fmt.Println(err)
					os.Exit(1)

				}
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Adds a node to a cluster",
			Action:  getSingleActionCLIHandler(&couchdb.AddNodeAction{}),
		},
		{
			Name:    "add_retry",
			Aliases: []string{"a"},
			Usage: "Adds a node to a cluster providing an array of nodes to try joining. " +
				"Specify a --node (-n) flag and a list of args of nodes to try to join. You" +
				"can also specify a number of retries with --retries (-r)",
			Action: func(c *cli.Context) {
				if c.NArg() < 1 {
					fmt.Println("Not enough arguments, you need to pass at least one" +
						"node to use as coordinator")
					return
				}

				n := couchdb.Node{}

				n.SetUser(admin)
				n.SetPass(password)
				n.SetPort(port)
				n.SetHost(c.String("node"))

				var template couchdb.TemplateExecutor = &couchdb.RetryTemplate{
					Coordinators:  c.Args(),
					CandidateNode: n,
					Wait:          c.Int("wait"),
				}

				if err := template.Do(c.String("node")); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			},
			Flags: []cli.Flag{
				cli.StringFlag{

					Name:  "node, n",
					Usage: "Node that you want to add to a cluster",
				},
				cli.IntFlag{
					Name:  "wait, w",
					Usage: "Seconds before trying connection again",
					Value: 2,
				},
			},
		},
		{
			Name:    "finish",
			Aliases: []string{"a"},
			Usage:   "Signals that the cluster formation has finished",
			Action: func(c *cli.Context) {
				t := &couchdb.FinishAction{}

				t.SetPass(password)
				t.SetPort(port)

				template := couchdb.CouchDbHTTP{
					ActionPerformer: t,
				}

				if err := template.Do(c.String("coord")); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
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

				err := couchdb.JoinAllClusterAction(admin, password, c.String("coord"), c.Args(), port)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				//if sendFinishCluster {
				//	template := &couchdb.CouchDbHTTP{}
				//	node := couchdb.Node{}
				//	node.SetPass(password)
				//	node.SetUser(admin)
				//	node.SetHost(c.String("coord"))
				//	node.SetPort(port)
				//	template.SetActionPerformer(&couchdb.FinishAction{
				//		Node: node,
				//	})
				//	template.Do(c.String("coord"))
				//}
			},
		},
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Check nodes in cluster",
			Action: func(c *cli.Context) {
				couchdb.CheckCluster(c.String("coord"), admin, password, port)
			},
		},
	}

	app.Name = "CouchDB 2 Orchestrator"
	app.Usage = "Common usage would be to pass a '-c IP' for a coordination node then 'join_all' as an action and a list of IP's to conform cluster"

	app.Run(os.Args)
}
