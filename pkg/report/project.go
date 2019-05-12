package report

type Branch struct {
	Name      string
	AgeInDays int
	Committer string
}

type Project struct {
	Name      string
	GroupName string
	Branches  []*Branch
}

func NewProject(name string, group string) *Project {
	return &Project{
		Name:      name,
		GroupName: group,
		Branches:  make([]*Branch, 0),
	}
}

func (p *Project) GetFullName() string {
	return p.GroupName + "/" + p.Name
}

func (p *Project) AddBranch(name string, age int, committer string) {
	p.Branches = append(p.Branches, &Branch{Name: name, AgeInDays: age, Committer: committer})
}
