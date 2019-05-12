package notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sync"

	"github.com/ns3777k/gitlab-branch-tracker/pkg/report"
	"gopkg.in/gomail.v2"
)

type EmailNotificator struct {
	projectsLock *sync.Mutex
	projects     []*report.Project
	options      *EmailNotificatorOptions
	mailer       Mailer
}

type EmailNotificatorOptions struct {
	From string
	To   []string
}

type Mailer interface {
	DialAndSend(m ...*gomail.Message) error
}

func NewEmailNotificator(mailer Mailer, options *EmailNotificatorOptions) *EmailNotificator {
	return &EmailNotificator{
		options:      options,
		projects:     make([]*report.Project, 0),
		projectsLock: &sync.Mutex{},
		mailer:       mailer,
	}
}

func (e *EmailNotificator) template() string {
	return `<html><body>
<h2>Gitlab: Left Branches Report</h2>
{{range $index, $project := .}}
<h3>{{$project.GetFullName}}</h3>
<ul>
{{range $index, $branch := .Branches}}
<li>{{$branch.Name}} {{$branch.AgeInDays}} days ago by {{$branch.Committer}}</li>
{{end}}
</ul>
{{end}}
</body></html>`
}

func (e *EmailNotificator) createMessage(content fmt.Stringer) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":       {e.options.From},
		"To":         e.options.To,
		"Subject":    {"Gitlab: Left Branches Report"},
		"Importance": {"high"},
	})
	m.SetBody("text/html", content.String())

	return m
}

func (e *EmailNotificator) hasProjects() bool {
	return len(e.projects) > 0
}

func (e *EmailNotificator) AddProject(project *report.Project) {
	e.projectsLock.Lock()
	e.projects = append(e.projects, project)
	e.projectsLock.Unlock()
}

func (e *EmailNotificator) Notify(ctx context.Context) error {
	if len(e.projects) == 0 {
		return nil
	}

	defer func() {
		e.projectsLock.Lock()
		e.projects = nil
		e.projectsLock.Unlock()
	}()

	buffer := new(bytes.Buffer)
	t := template.Must(template.New("").Parse(e.template()))
	if err := t.Execute(buffer, e.projects); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return e.mailer.DialAndSend(e.createMessage(buffer))
	}
}
