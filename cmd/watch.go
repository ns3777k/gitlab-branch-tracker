package cmd

import (
	"github.com/ns3777k/gitlab-branch-tracker/pkg"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

func WatchAction(c *cli.Context) error {
	dsn := c.String("gitlab-dsn")
	logger := log.New(os.Stdout, "> ", log.LstdFlags)
	watcher := pkg.NewWatcher(nil, dsn, logger)

	watcher.Start(time.Minute * 5)

	return nil
}
