package project

import (
	"context"
	"errors"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/helper"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type ProjectService struct {
	ProjectRepo storage.ProjectRepository
}

func NewService(ProjectRepo storage.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepo: ProjectRepo}
}

func (p ProjectService) NewProject(ctx context.Context, project *entities.ProjectImport) (uuid entities.ProjectUUID, err error) {
	if project == nil {
		return entities.InvalidProjectUUID, errors.New("given nil project")
	}

	proj := &entities.Project{
		UUID:        helper.RandomUUID(),
		FileIDs:     project.FileIDs,
		ProjectType: project.ProjectType,
		// DateArchived: time.Now().UTC().Truncate(time.Second),
		// database layer handles this for us
	}
	if !project.DateCompleted.IsZero() {
		proj.DateAnnounced = project.DateAnnounced.UTC().Truncate(time.Second)
	}
	if !project.DateAnnounced.IsZero() {
		proj.DateCompleted = project.DateCompleted.UTC().Truncate(time.Second)
	}

	uuid, err = p.ProjectRepo.NewProject(ctx, proj)
	if err != nil {
		return entities.InvalidProjectUUID, err
	}

	return uuid, err
}

func (p ProjectService) DeleteProject(ctx context.Context, project_uuid entities.ProjectUUID) (err error) {
	return p.ProjectRepo.DeleteProject(ctx, project_uuid)
}

func (p ProjectService) AssignFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	return p.AssignFile(ctx, project_uuid, file_id)
}

func (p ProjectService) AssignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error) {
	return p.ProjectRepo.AssignYoutube(ctx, project_uuid, youtube_id)
}

func (p ProjectService) UnassignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error) {
	return p.ProjectRepo.UnassignYoutube(ctx, project_uuid, youtube_id)
}

func (p ProjectService) AssignPrimaryFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	panic("unimplemented")
}

func (p ProjectService) UnassignPrimaryFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	panic("unimplemented")
}

func (p ProjectService) UnassignFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	panic("unimplemented")
}

func (p ProjectService) SetProjectType(ctx context.Context, project_uuid entities.ProjectUUID, project_type entities.ProjectType) (err error) {
	panic("unimplemented")
}

func (p ProjectService) GetProject(ctx context.Context, project_uuid entities.ProjectUUID) (project entities.Project, err error) {
	panic("unimplemented")
}

func (p ProjectService) GetProjectYoutube(ctx context.Context, project_uuid entities.ProjectUUID) (youtube_ids []entities.YoutubeVideoID, err error) {
	panic("unimplemented")
}
