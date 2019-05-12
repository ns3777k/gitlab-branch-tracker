package notification

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ns3777k/gitlab-branch-tracker/pkg/report"
	"gopkg.in/gomail.v2"
)

type mailerMock struct {
	dialFn func(messages []*gomail.Message) error
}

func (m *mailerMock) DialAndSend(message ...*gomail.Message) error {
	return m.dialFn(message)
}

func TestEmailNotificator_EmptyProjects(t *testing.T) {
	m := &mailerMock{dialFn: func(message []*gomail.Message) error {
		return errors.New("return right away if there are no projects")
	}}
	n := NewEmailNotificator(m, &EmailNotificatorOptions{To: []string{"to"}, From: "from"})
	if err := n.Notify(context.TODO()); err != nil {
		t.Error(err)
	}
}

func TestEmailNotificator_ValidHeaders(t *testing.T) {
	project := report.NewProject("gitlab-branch-tracker", "ns3777k")
	project.AddBranch("tests", 255, "ns3777k")

	m := &mailerMock{dialFn: func(message []*gomail.Message) error {
		from := message[0].GetHeader("From")
		to := message[0].GetHeader("To")

		if from[0] != "from" {
			return fmt.Errorf("from address mismatch. expected: from, got: %s", from[0])
		}

		if to[0] != "to" {
			return fmt.Errorf("to address mismatch. expected: to, got: %s", to[0])
		}

		return nil
	}}

	n := NewEmailNotificator(m, &EmailNotificatorOptions{To: []string{"to"}, From: "from"})
	n.AddProject(project)

	if err := n.Notify(context.TODO()); err != nil {
		t.Error(err)
	}
}

func TestEmailNotificator_CleanupProjects(t *testing.T) {
	project := report.NewProject("gitlab-branch-tracker", "ns3777k")
	project.AddBranch("tests", 255, "ns3777k")

	m := &mailerMock{dialFn: func(message []*gomail.Message) error {
		return nil
	}}
	n := NewEmailNotificator(m, &EmailNotificatorOptions{To: []string{"to"}, From: "from"})
	n.AddProject(project)

	if err := n.Notify(context.TODO()); err != nil {
		t.Error(err)
	}

	if n.hasProjects() {
		t.Error("notificator must clean projects slice after sending")
	}
}
