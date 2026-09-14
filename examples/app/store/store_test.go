package store

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/examples/app/models"
)

func TestCreateGetProject(t *testing.T) {
	db, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}

	store := NewStore(db)
	p1 := models.NewProject("test project")
	p2 := models.NewProject("test project2")
	projectList := []*models.Project{p1, p2}
	store.CreateProject(p1)
	store.CreateProject(p2)
	projects, err := store.GetProjects()
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	if got, want := len(projects), 2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	for i, p := range projects {
		if got, want := p.DispName, projectList[i].DispName; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := p.ID(), projectList[i].ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestCreateGetActivity(t *testing.T) {
	db, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}

	store := NewStore(db)
	project := models.NewProject("test project")
	store.CreateProject(project)
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	a1 := models.NewActivity("task 1", 5*time.Minute, project)
	a2 := models.NewActivity("task 2", 4*time.Minute, project)

	activityList := []*models.Activity{a1, a2}

	store.CreateActivity(a1)
	store.CreateActivity(a2)

	activities, err := store.GetActivitiesByProject(project)
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	if got, want := len(activities), 2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	for i, a := range activities {
		if got, want := a.DispName, activityList[i].DispName; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := a.ID(), activityList[i].ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := a.Duration(), activityList[i].Duration(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := a.ProjectID.ID(), activityList[i].ProjectID.ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestCreateGetDependency(t *testing.T) {
	db, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}

	store := NewStore(db)
	project := models.NewProject("test project")
	store.CreateProject(project)
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	a1 := models.NewActivity("task 1", 5*time.Minute, project)
	a2 := models.NewActivity("task 2", 4*time.Minute, project)
	a3 := models.NewActivity("task 3", 2*time.Minute, project)

	store.CreateActivity(a1)
	store.CreateActivity(a2)
	store.CreateActivity(a3)

	d1 := models.NewDependency(a1, a2, enums.FS, project)
	d2 := models.NewDependency(a3, a2, enums.FS, project)

	dependencyList := []*models.Dependency{d1, d2}
	store.CreateDependency(d1)
	store.CreateDependency(d2)

	dependencies, err := store.GetDependenciesByProject(project)
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	if got, want := len(dependencies), 2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	for i, d := range dependencies {
		if got, want := d.ID(), dependencyList[i].ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := d.ProjectID.ID(), dependencyList[i].ProjectID.ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := d.RelationshipType, dependencyList[i].RelationshipType; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := d.Predecessor.ID(), dependencyList[i].Predecessor.ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := d.Successor.ID(), dependencyList[i].Successor.ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestDeleteProject(t *testing.T) {
	db, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}

	store := NewStore(db)
	p1 := models.NewProject("test project")
	store.CreateProject(p1)
	projects, err := store.GetProjects()
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	if got, want := len(projects), 1; got != want {
		t.Errorf("expected GetProjects to return slice of length %v, got %v", want, got)
	}
	for i, p := range projects {
		if got, want := p.DispName, p1.DispName; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := p.ID(), p1.ID(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if i > 0 {
			t.Errorf("expected only one row returned")
		}
	}
	err = store.DeleteProject(p1)
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	projects, err = store.GetProjects()

	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	if got, want := len(projects), 0; got != want {
		t.Errorf("expected GetProjects to return slice of length %v, got %v", want, got)
	}
}
