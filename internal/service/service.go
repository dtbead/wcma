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
	DeleteProject(ctx context.Context, project_uuid entities.ProjectUUID) (err error)
	AssignFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error)
	AssignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error)
	AssignPrimaryFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error)
	UnassignPrimaryFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error)
	UnassignFile(ctx context.Context, project_uuid entities.ProjectUUID, file_id entities.FileID) (err error)
	UnassignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error)
	SetProjectType(ctx context.Context, project_uuid entities.ProjectUUID, project_type entities.ProjectType) (err error)
	GetProject(ctx context.Context, project_uuid entities.ProjectUUID) (project entities.Project, err error)
	GetProjectYoutube(ctx context.Context, project_uuid entities.ProjectUUID) (youtube_ids []entities.YoutubeVideoID, err error)
}

type FileService interface {
	NewFile(ctx context.Context, file io.Reader, extension string) (file_id entities.FileID, err error)
	DeleteFile(ctx context.Context, file_id entities.FileID) (err error)
	NewTempFile(ctx context.Context) (file io.ReadWriteCloser, err error)
	GetHash(ctx context.Context, file_id entities.FileID) (err error, hashes entities.Hashes)
	GetReader(ctx context.Context, file_id entities.FileID) (file io.ReadCloser, err error)
	GetFileRelationship(ctx context.Context, file_id entities.FileID) (relationships entities.FileRelationship, err error)
}

type YoutubeService interface {
	NewYoutube(ctx context.Context, file_id entities.FileID, youtube *entities.Youtube) (err error)
	GetYoutubeFileIDs(ctx context.Context, youtube_id entities.YoutubeVideoID) (file_ids []entities.FileID, err error)
	GetTitle(ctx context.Context, youtube_id entities.YoutubeVideoID) (title string, err error)
	GetDescription(ctx context.Context, youtube_id entities.YoutubeVideoID) (description string, err error)
	GetChannelVideos(ctx context.Context, channel_id entities.YoutubeChannelID) (videos []entities.YoutubeVideoID, err error)
	GetYoutubeVideo(ctx context.Context, youtube_id entities.YoutubeVideoID) (video entities.YoutubeVideo, err error)
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
