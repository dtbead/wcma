package service

import (
	"context"
	"errors"
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
	NewTempFile(ctx context.Context) (file io.ReadWriteCloser, err error)
}

type YoutubeService interface {
	NewYoutube(ctx context.Context, file_id entities.FileID, youtube *entities.Youtube) (err error)
	GetYoutubeFileIDs(ctx context.Context, youtube_id entities.YoutubeVideoID) (file_ids []entities.FileID, err error)
}

func (s Service) DownloadYoutube(ctx context.Context, url string, downloader entities.YoutubeDownloader) (err error) {
	tmp, err := s.FileService.NewTempFile(ctx)
	if err != nil {
		return err
	}
	defer tmp.Close()

	yt, ext, err := downloader.Download(ctx, url, tmp)
	if err != nil {
		return err
	}

	file_id, err := s.FileService.NewFile(ctx, tmp, ext)
	if err != nil {
		return err
	}

	err = s.YoutubeService.NewYoutube(ctx, file_id, yt)
	if err != nil {
		return errors.Join(err, s.FileService.DeleteFile(ctx, file_id))
	}

	return nil
}
