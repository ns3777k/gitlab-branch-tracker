package notification

import (
	"context"

	"github.com/ns3777k/gitlab-branch-tracker/pkg/report"
)

type Notificator interface {
	AddProject(project *report.Project)
	Notify(ctx context.Context) error
}
