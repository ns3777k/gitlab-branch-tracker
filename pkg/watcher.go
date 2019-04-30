package pkg

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Watcher struct {
	logger *log.Logger
	dsn string
	httpClient *http.Client
}

func NewWatcher(httpClient *http.Client, dsn string, logger *log.Logger) *Watcher {
	return &Watcher{
		logger: logger,
		dsn: dsn,
		httpClient: httpClient,
	}
}

func (w *Watcher) createGitlabClient() (*gitlab.Client, error) {
	u, err := url.Parse(w.dsn)
	if err != nil {
		return nil, err
	}

	token, ok := u.User.Password()
	if !ok {
		token = ""
	}

	gitlabClient := gitlab.NewClient(w.httpClient, token)
	err = gitlabClient.SetBaseURL(u.String())

	return gitlabClient, err
}

func (w *Watcher) Start(interval time.Duration) error {
	gitlabClient, err := w.createGitlabClient()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			branches, _, err := gitlabClient.Branches.ListBranches("group/repo", nil)
			w.logger.Print(err)
			for _, branch := range branches {
				w.logger.Print(branch.Name)
			}
		}
	}

	return nil
}
