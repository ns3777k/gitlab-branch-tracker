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

type Watcher struct {
	gitlabClient *gitlab.Client
	logger       *log.Logger
	projects     chan *report.Project
	workers      []*worker
	wg           *sync.WaitGroup
	notificator  notification.Notificator
}

type Options struct {
	WorkersCount    int
	BranchExcludes  map[string]struct{}
	MaxDaysThrottle int
}

type WatchOptions struct {
	StartImmediately bool
	Interval         time.Duration
	Namespaces       []string
}

func NewWatcher(
	gitlabClient *gitlab.Client,
	logger *log.Logger,
	notificator notification.Notificator,
	options *Options,
) *Watcher {
	if options.WorkersCount == 0 {
		options.WorkersCount = 1
	}

	workers := make([]*worker, options.WorkersCount)
	for i := 0; i < options.WorkersCount; i++ {
		workers[i] = newWorker(gitlabClient, logger, notificator, &workerOptions{
			branchExcludes:  options.BranchExcludes,
			maxDaysThrottle: options.MaxDaysThrottle,
		})
	}

	return &Watcher{
		gitlabClient: gitlabClient,
		logger:       logger,
		projects:     make(chan *report.Project, options.WorkersCount),
		workers:      workers,
		notificator:  notificator,
		wg:           &sync.WaitGroup{},
	}
}

func (w *Watcher) addProjectsFromNamespace(ctx context.Context, namespace string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	group, _, err := w.gitlabClient.Groups.GetGroup(namespace, gitlab.WithContext(timeoutCtx))
	if err != nil {
		return err
	}

	for _, project := range group.Projects {
		reportProject := report.NewProject(project.Name, group.Name)
		logger := w.logger.WithField("project", reportProject.GetFullName())

		if project.Archived {
			logger.Debug("skipping archived project")
			continue
		}

		logger.Info("adding project")

		select {
		case <-ctx.Done():
			close(w.projects)
			return ctx.Err()
		default:
			w.projects <- reportProject
		}
	}

	return nil
}

func (w *Watcher) work(ctx context.Context) {
	for _, worker := range w.workers {
		w.wg.Add(1)
		go worker.work(ctx, w.projects, w.wg)
	}
}

func (w *Watcher) doWatch(ctx context.Context, namespace []string) {
	w.work(ctx)

	go func() {
		for _, namespace := range namespace {
			w.logger.WithField("namespace", namespace).Info("adding namespace")
			if err := w.addProjectsFromNamespace(ctx, namespace); err != nil {
				w.logger.WithError(err).Error("error adding namespace")
			}
		}
	}()

	w.wg.Wait()
	w.logger.Info("sending notification")

	if err := w.notificator.Notify(ctx); err != nil {
		w.logger.WithError(err).Error("notifying error")
	}
}

func (w *Watcher) Watch(ctx context.Context, options *WatchOptions) error {
	if options.StartImmediately {
		w.doWatch(ctx, options.Namespaces)
	}

	ticker := time.NewTicker(options.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.doWatch(ctx, options.Namespaces)
		case <-ctx.Done():
			w.wg.Wait()
			return ctx.Err()
		}
	}
}
