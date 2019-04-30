package main

import (
	"github.com/ns3777k/gitlab-branch-tracker/cmd"
	"github.com/urfave/cli"
	"log"
	"os"
)

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = "gitlab-branch-tracker"
	app.Commands = []cli.Command{
		{
			Name: "watch",
			Usage: "start watching for branches",
			Action: cmd.WatchAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "gitlab-dsn",
					Usage: "gitlab dsn",
					EnvVar: "GITLAB_DSN",
				},
			},
		},
	}

	return app
}

func main() {
	if err := createApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
