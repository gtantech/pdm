package store

import (
	"context"
	"errors"
	"uuid"

	"github.com/gtantech/pdm/examples/app/internal/db"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (db.Project, error)
	GetByName(ctx context.Context, search string) ([]db.Project, error)
	GetProjects(ctx context.Context) ([]db.Project, error)
	Create(ctx context.Context, params CreateProjectParams) (db.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateProjectParams struct {
	DisplayName string
}

type projectStore struct {
	queries *db.Queries
}

// GetProjects implements [ProjectStore].
func (p *projectStore) GetProjects(ctx context.Context) ([]db.Project, error) {
	return p.queries.FindAllProjects(ctx)
}

// GetByName implements [ProjectStore].
func (p *projectStore) GetByName(ctx context.Context, search string) ([]db.Project, error) {
	return p.queries.FindProjectByName(ctx, search)
}

// Create implements [ProjectStore].
func (p *projectStore) Create(ctx context.Context, params CreateProjectParams) (db.Project, error) {
	return p.queries.InsertProject(ctx, db.InsertProjectParams{ID: uuid.NewV7().String(), DispName: params.DisplayName})
}

// Delete implements [ProjectStore].
func (p *projectStore) Delete(ctx context.Context, id uuid.UUID) error {
	return p.queries.DeleteProject(ctx, id.String())
}

// GetByID implements [ProjectStore].
func (p *projectStore) GetByID(ctx context.Context, id uuid.UUID) (db.Project, error) {
	return p.queries.FindProjectById(ctx, id.String())
}

func NewProjectStore(queries *db.Queries) *projectStore {
	return &projectStore{
		queries: queries,
	}
}

var _ ProjectStore = (*projectStore)(nil) //ensures projectStore implements ProjectStore at compile time
