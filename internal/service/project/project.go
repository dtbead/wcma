package project

import (
	"context"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type ProjectService struct {
	ProjectRepo *storage.ProjectRepository
}

func NewService(ProjectRepo storage.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepo: &ProjectRepo}
}

func (p ProjectService) NewProject(ctx context.Context, project *entities.ProjectImport) (uuid entities.ProjectUUID, err error) {
	panic("unimplemented")
}

func (p ProjectService) DeleteProject(ctx context.Context, project_uuid entities.ProjectUUID) (err error) {
	panic("unimplemented")
}

func (p ProjectService) AssignFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error) {
	panic("unimplemented")
}

func (p ProjectService) AssignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error) {
	panic("unimplemented")
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

func (p ProjectService) UnassignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error) {
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
