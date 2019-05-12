package watcher

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ns3777k/gitlab-branch-tracker/pkg/notification"
	"github.com/ns3777k/gitlab-branch-tracker/pkg/report"
	"github.com/xanzy/go-gitlab"
)

type worker struct {
	gitlabClient *gitlab.Client
	notificator  notification.Notificator
	logger       *log.Logger
	options      *workerOptions
}

type workerOptions struct {
	branchExcludes  map[string]struct{}
	maxDaysThrottle int
}

func newWorker(
	gitlabClient *gitlab.Client,
	logger *log.Logger,
	notificator notification.Notificator,
	options *workerOptions,
) *worker {
	return &worker{
		gitlabClient: gitlabClient,
		notificator:  notificator,
		logger:       logger,
		options:      options,
	}
}

func (w *worker) workProject(ctx context.Context, project *report.Project) error {
	page := 1
	fullName := project.GetFullName()

	for {
		listOptions := &gitlab.ListBranchesOptions{Page: page}
		timeoutContext, timeoutCancel := context.WithTimeout(ctx, time.Second*10)

		branches, response, err := w.gitlabClient.Branches.ListBranches(
			fullName, listOptions, gitlab.WithContext(timeoutContext))
		if err != nil {
			timeoutCancel()
			return err
		}

		for _, branch := range branches {
			diffDays := time.Since(*branch.Commit.CommittedDate).Round(time.Hour*24).Hours() / 24
			if diffDays < float64(w.options.maxDaysThrottle) {
				continue
			}

			logger := w.logger.WithField("project", fullName).WithField("branch", branch.Name)

			if _, ok := w.options.branchExcludes[branch.Name]; ok {
				logger.Debug("skipping branch")
				continue
			}

			logger.Info("adding branch to report")
			project.AddBranch(branch.Name, int(diffDays), branch.Commit.CommitterName)
		}

		if response.CurrentPage == response.TotalPages {
			timeoutCancel()
			break
		}

		page = response.CurrentPage + 1
		timeoutCancel()
	}

	return nil
}

func (w *worker) work(ctx context.Context, projects chan *report.Project, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case project, ok := <-projects:
			if !ok {
				return
			}
			if err := w.workProject(ctx, project); err != nil {
				w.logger.WithField("project", project.GetFullName()).
					WithError(err).Error("failed to handle project")
				continue
			}

			if len(project.Branches) > 0 {
				w.notificator.AddProject(project)
			}

		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				w.logger.WithError(err).Error("worker context done")
			}
			return

		case <-time.After(time.Second * 10):
			w.logger.Info("worker exiting due to timeout")
			return
		}
	}
}
