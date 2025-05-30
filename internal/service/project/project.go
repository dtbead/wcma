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
