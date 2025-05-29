package project

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/queries"
)

type ProjectRepository struct {
	db *sql.DB
	q  *queries.Queries
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
		q:  queries.New(db),
	}
}

func isEmptyTime(t time.Time) bool {
	return t.IsZero()
}

// DateArchived is ignored implicitly.
func (p ProjectRepository) NewProject(ctx context.Context, project *entities.Project) (uuid entities.ProjectUUID, err error) {
	if project == nil {
		return entities.InvalidProjectUUID, errors.New("nil project")
	}

	s, err := p.q.NewProject(ctx, queries.NewProjectParams{
		Uuid:          project.UUID,
		Type:          queries.Projecttype(project.ProjectType.ToString()),
		DateAnnounced: sql.NullTime{Time: project.DateAnnounced.UTC().Truncate(time.Second), Valid: !isEmptyTime(project.DateAnnounced.UTC())},
		DateCompleted: sql.NullTime{Time: project.DateCompleted.UTC().Truncate(time.Second), Valid: !isEmptyTime(project.DateCompleted.UTC())},
	})
	if err != nil {
		return entities.InvalidProjectUUID, err
	}

	return entities.ProjectUUID(s), nil

}

func (p ProjectRepository) DeleteProject(ctx context.Context, uuid entities.ProjectUUID) (err error) {
	return p.q.DeleteProjectByUUID(ctx, string(uuid))
}

func (p ProjectRepository) GetProject(ctx context.Context, uuid entities.ProjectUUID) (project *entities.Project, err error) {
	res, err := p.q.GetProjectByUUID(ctx, string(uuid))
	if err != nil {
		return nil, err
	}

	project_type, err := entities.NewProjectType(string(res.Type))
	if err != nil {
		return nil, err
	}

	file_ids, err := p.GetProjectVideos(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &entities.Project{
		UUID:          res.Uuid,
		ProjectType:   project_type,
		DateAnnounced: res.DateAnnounced.Time,
		DateArchived:  res.DateArchived,
		DateCompleted: res.DateCompleted.Time,
		FileIDs:       file_ids,
	}, nil
}
func (p ProjectRepository) AssignProjectFile(ctx context.Context, uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	return p.q.AssignProjectFile(ctx, queries.AssignProjectFileParams{Uuid: string(uuid), FileID: int64(file_id)})
}
func (p ProjectRepository) UnassignProjectVideo(ctx context.Context, uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	return p.q.UnassignProjectFile(ctx, int64(file_id))
}
func (p ProjectRepository) GetProjectVideos(ctx context.Context, uuid entities.ProjectUUID) (file_ids []entities.FileID, err error) {
	res, err := p.q.GetProjectFile(ctx, string(uuid))
	if err != nil {
		return nil, err
	}

	// this is really dumb, golang...
	file_ids = make([]entities.FileID, 0, len(res))
	for _, v := range res {
		file_ids = append(file_ids, entities.FileID(v))
	}

	return file_ids, nil
}
