package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/ns3777k/gitlab-branch-tracker/pkg/notification"
	"github.com/ns3777k/gitlab-branch-tracker/pkg/watcher"
	"github.com/urfave/cli"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/gomail.v2"
)

func createGitlabClient(dsn string, ignoreCertVerification bool) (*gitlab.Client, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	token, ok := u.User.Password()
	if !ok {
		token = ""
	}

	transport := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreCertVerification, //nolint:gosec
			},
		},
	}

	gitlabClient := gitlab.NewClient(transport, token)
	err = gitlabClient.SetBaseURL(u.String())

	return gitlabClient, err
}

func createNotificator(dsn string, from string, to []string) (*notification.EmailNotificator, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, err
	}

	password, ok := u.User.Password()
	if !ok {
		password = ""
	}

	dialer := gomail.NewDialer(u.Hostname(), port, u.User.Username(), password)
	notificatorOptions := &notification.EmailNotificatorOptions{From: from, To: to}
	notificator := notification.NewEmailNotificator(dialer, notificatorOptions)

	return notificator, nil
}

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	cancel()
}

func WatchAction(c *cli.Context) error {
	logger := log.New()

	if c.Bool("debug") {
		logger.SetLevel(log.DebugLevel)
	}

	notificator, err := createNotificator(
		c.String("smtp-dsn"),
		c.String("smtp-from"),
		c.StringSlice("smtp-to"),
	)
	if err != nil {
		return err
	}

	gitlabClient, err := createGitlabClient(c.String("gitlab-dsn"), c.Bool("gitlab-ignore-cert-verification"))
	if err != nil {
		return err
	}

	excludes := make(map[string]struct{})

	for _, branch := range c.StringSlice("exclude-branch") {
		excludes[branch] = struct{}{}
	}

	interval := c.Duration("interval")
	logger.WithField("interval", interval.String()).Info("start watching")

	w := watcher.NewWatcher(gitlabClient, logger, notificator, &watcher.Options{
		WorkersCount:    c.Int("max-workers"),
		BranchExcludes:  excludes,
		MaxDaysThrottle: c.Int("max-days"),
	})

	rootCtx, shutdownFn := context.WithCancel(context.Background())
	defer shutdownFn()

	go handleSignals(shutdownFn)

	err = w.Watch(rootCtx, &watcher.WatchOptions{
		StartImmediately: c.Bool("start-immediately"),
		Interval:         interval,
		Namespaces:       c.StringSlice("namespace"),
	})

	return err
}
