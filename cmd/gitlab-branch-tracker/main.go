package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

var version = "dev" //nolint:gochecknoglobals

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = "gitlab-branch-tracker"
	app.Version = version
	app.Usage = "start watching for branches"
	app.Action = WatchAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "gitlab-dsn",
			Usage:  "https://username:password@my.git/api/v4",
			EnvVar: "GITLAB_DSN",
		},
		cli.StringFlag{
			Name:   "smtp-dsn",
			Usage:  "smtp://username:password@host:port",
			EnvVar: "SMTP_DSN",
		},
		cli.StringFlag{
			Name:   "smtp-from",
			Usage:  "sender's email",
			EnvVar: "SMTP_FROM",
		},
		cli.StringSliceFlag{
			Name:   "smtp-to",
			Usage:  "emails to deliver to separated by comma",
			EnvVar: "SMTP_TO",
		},
		cli.StringSliceFlag{
			Name:   "namespace",
			Usage:  "namespace to get projects from",
			EnvVar: "NAMESPACE",
		},
		cli.DurationFlag{
			Name:   "interval",
			Value:  time.Hour * 24,
			Usage:  "interval between reports",
			EnvVar: "INTERVAL",
		},
		cli.BoolFlag{
			Name:   "start-immediately",
			Usage:  "start right away and schedule an interval after",
			EnvVar: "START_IMMEDIATELY",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "expose more verbose information",
			EnvVar: "DEBUG",
		},
		cli.IntFlag{
			Name:   "max-days",
			Usage:  "maximum amount of days to include branch into notification",
			Value:  90,
			EnvVar: "MAX_DAYS",
		},
		cli.IntFlag{
			Name:   "max-workers",
			Usage:  "how many projects to handle at a time",
			Value:  1,
			EnvVar: "MAX_WORKERS",
		},
		cli.StringSliceFlag{
			Name:   "exclude-branch",
			Usage:  "exclude branches like master",
			Value:  &cli.StringSlice{"master"},
			EnvVar: "EXCLUDE_BRANCH",
		},
	}

	return app
}

func main() {
	if err := createApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
