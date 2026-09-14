package store

import (
	"database/sql"
	"errors"
	"time"
	"uuid"

	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/examples/app/models"
)

var ErrorNotFound = errors.New("Record Not Found")

type Store interface {
	CreateProject(project *models.Project) (*models.Project, error)
	CreateActivity(activity *models.Activity) (*models.Activity, error)
	CreateDependency(dependency *models.Dependency) (*models.Dependency, error)
	GetProjects() ([]models.Project, error)
	GetActivitiesByProject(project models.ProjectID) ([]models.Activity, error)
	GetDependenciesByProject(project models.ProjectID) ([]models.Dependency, error)
	DeleteProject(project models.ProjectID) error
	DeleteActivity(activity models.ActivityID) error
	DeleteDependency(dependency models.DependencyID) error
	FindProjectByName(search string) ([]models.Project, error)
	FindActivityByNameAndProject(search string, project models.ProjectID) ([]models.Activity, error)
}

type Storage struct {
	db *sql.DB
}

// CreateActivity implements [Store].
func (s *Storage) CreateActivity(activity *models.Activity) (*models.Activity, error) {
	_, err := s.db.Exec("INSERT INTO activities (id, project_id, disp_name, duration) VALUES (?, ?, ?, ?)", activity.ID(), activity.ProjectID.ID(), activity.DispName, activity.Duration())
	if err != nil {
		return nil, err
	}

	return activity, nil
}

// CreateDependency implements [Store].
func (s *Storage) CreateDependency(dependency *models.Dependency) (*models.Dependency, error) {
	_, err := s.db.Exec("INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) VALUES (?, ?, ?, ?, ?)",
		dependency.ID(), dependency.ProjectID.ID(), dependency.RelationshipType, dependency.Predecessor.ID(), dependency.Successor.ID())
	if err != nil {
		return nil, err
	}

	return dependency, nil
}

// CreateProject implements [Store].
func (s *Storage) CreateProject(project *models.Project) (*models.Project, error) {
	_, err := s.db.Exec("INSERT INTO projects (id, disp_name) VALUES (?, ?)", project.ID(), project.DispName)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteActivity implements [Store].
func (s *Storage) DeleteActivity(activity models.ActivityID) error {
	result, err := s.db.Exec("DELETE FROM activities WHERE id = ?", activity.ID())
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return ErrorNotFound
	}
	return err
}

// DeleteDependency implements [Store].
func (s *Storage) DeleteDependency(dependency models.DependencyID) error {
	result, err := s.db.Exec("DELETE FROM dependencies WHERE id = ?", dependency.ID())
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return ErrorNotFound
	}
	return err
}

// DeleteProject implements [Store].
func (s *Storage) DeleteProject(project models.ProjectID) error {
	result, err := s.db.Exec("DELETE FROM projects WHERE id = ?", project.ID())
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return ErrorNotFound
	}
	return err
}

// FindActivityByNameAndProject implements [Store].
func (s *Storage) FindActivityByNameAndProject(search string, project models.ProjectID) ([]models.Activity, error) {
	rows, err := s.db.Query("SELECT * FROM activities WHERE disp_name LIKE ? AND project_id = ?", "%"+search+"%", project.ID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		car, err := scanActivity(rows)
		if err != nil {
			return nil, err
		}
		activities = append(activities, car)
	}

	return activities, nil
}

func scanActivity(row *sql.Rows) (models.Activity, error) {
	var activityID uuid.UUID
	var activityProjectID uuid.UUID
	var activityDispName string
	var activityDuration time.Duration
	err := row.Scan(&activityID, &activityProjectID, &activityDispName, &activityDuration)
	activity := *models.NewActivity(activityDispName, activityDuration, models.NewProjectID(activityProjectID))
	activity.ActivityID = models.NewActivityID(activityID)
	return activity, err
}

// FindProjectByName implements [Store].
func (s *Storage) FindProjectByName(search string) ([]models.Project, error) {
	rows, err := s.db.Query("SELECT * FROM projects WHERE disp_name LIKE ?", "%"+search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func scanProject(row *sql.Rows) (models.Project, error) {
	var projectID uuid.UUID
	var projectDispName string
	err := row.Scan(&projectID, &projectDispName)
	project := *models.NewProject(projectDispName)
	project.ProjectID = models.NewProjectID(projectID)
	return project, err
}

// GetActivitiesByProject implements [Store].
func (s *Storage) GetActivitiesByProject(project models.ProjectID) ([]models.Activity, error) {
	rows, err := s.db.Query("SELECT * FROM activities WHERE project_id = ?", project.ID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		activity, err := scanActivity(rows)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// GetDependenciesByProject implements [Store].
func (s *Storage) GetDependenciesByProject(project models.ProjectID) ([]models.Dependency, error) {
	rows, err := s.db.Query("SELECT * FROM dependencies WHERE project_id = ?", project.ID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dependencies []models.Dependency
	for rows.Next() {
		dependency, err := scanDependency(rows)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func scanDependency(row *sql.Rows) (models.Dependency, error) {
	var dependencyID uuid.UUID
	var dependencyProjectID uuid.UUID
	var dependencyRelationship enums.RelationshipType
	var predecessorActivityId uuid.UUID
	var successorActivityId uuid.UUID
	err := row.Scan(&dependencyID, &dependencyProjectID, &dependencyRelationship, &predecessorActivityId, &successorActivityId)
	dependency := *models.NewDependency(models.NewActivityID(predecessorActivityId), models.NewActivityID(successorActivityId), dependencyRelationship, models.NewProjectID(dependencyProjectID))
	dependency.DependencyID = models.NewDependencyID(dependencyID)
	return dependency, err
}

// GetProjects implements [Store].
func (s *Storage) GetProjects() ([]models.Project, error) {
	rows, err := s.db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		car, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, car)
	}

	return projects, nil
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

var _ Store = (*Storage)(nil) //ensures Storage implements Store at compile time
