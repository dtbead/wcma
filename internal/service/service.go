package service

import (
	"context"
	"io"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type Service struct {
	ProjectService storage.ProjectRepository
	FileService    storage.FileRepository
	YoutubeService storage.YoutubeRepository
}

func NewService(repositories *storage.Repository) *Service {
	return &Service{
		ProjectService: repositories.Project,
		FileService:    repositories.File,
		YoutubeService: repositories.Youtube,
	}
}

type Project struct {
	ProjectType                                entities.ProjectType
	DateAnnounced, DateCompleted, DateArchived time.Time
	FileIDs                                    []entities.FileID
}

type ProjectYoutube struct {
	Project *Project
	Youtube *entities.Youtube
}

type ProjectService interface {
	NewProject(ctx context.Context, project *Project) (uuid entities.ProjectUUID, err error)
}

type FileService interface {
	NewFile(ctx context.Context, file io.Reader, extension string) (file_id entities.FileID, err error)
	DeleteFile(ctx context.Context, file_id entities.FileID) (err error)
}

type YoutubeService interface {
	NewYoutube(ctx context.Context, video io.Reader, project_youtube ProjectYoutube) (err error)
}

type YoutubeDownloader interface {
	Download(ctx context.Context, url string, output io.Writer) (youtube *entities.Youtube, extension string, err error)
}
