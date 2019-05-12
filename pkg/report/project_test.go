package report

import (
	"testing"
)

func TestProjectGetFullName(t *testing.T) {
	expected := "test-group/test-project"
	project := NewProject("test-project", "test-group")

	if project.GetFullName() != expected {
		t.Errorf("get full name error. expected: %s, got %s", expected, project.GetFullName())
	}
}

func TestProjectAddBranch(t *testing.T) {
	project := NewProject("test-project", "test-group")
	project.AddBranch("test", 255, "qa")

	if len(project.Branches) != 1 {
		t.Errorf("invalid branches length. expected: 1, got: %d", len(project.Branches))
	}
}
