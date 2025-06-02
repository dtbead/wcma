package service

import (
	"context"
	"io"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/service/file"
	"github.com/dtbead/wc-maps-archive/internal/service/project"
	"github.com/dtbead/wc-maps-archive/internal/service/youtube"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type Service struct {
	ProjectService ProjectService
	FileService    FileService
	YoutubeService YoutubeService
}

func NewService(repositories *storage.Repository) *Service {
	return &Service{
		ProjectService: project.NewService(repositories.Project),
		FileService:    file.NewService(repositories.File),
		YoutubeService: youtube.NewService(repositories.Youtube),
	}
}

type ProjectService interface {
	NewProject(ctx context.Context, project *entities.ProjectImport) (uuid entities.ProjectUUID, err error)
}

type FileService interface {
	NewFile(ctx context.Context, file io.Reader, extension string) (file_id entities.FileID, err error)
	DeleteFile(ctx context.Context, file_id entities.FileID) (err error)
}

type YoutubeService interface {
	NewYoutube(ctx context.Context, video io.Reader, project_youtube *entities.ProjectYoutube) (err error)
}
