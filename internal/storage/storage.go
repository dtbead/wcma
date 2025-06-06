package storage

import (
	"context"
	"io"

	"github.com/dtbead/wc-maps-archive/internal/entities"
)

type FileRepository interface {
	NewFile(ctx context.Context, file io.Reader, extension string) (file_id entities.FileID, err error)
	DeleteFile(ctx context.Context, file_id entities.FileID) (err error)
	GetFile(ctx context.Context, file_id entities.FileID) (file_metadata *entities.File, err error)
	NewTempFile(ctx context.Context) (file io.ReadWriteCloser, err error)
}

type ProjectRepository interface {
	NewProject(ctx context.Context, project *entities.Project) (uuid entities.ProjectUUID, err error)
	DeleteProject(ctx context.Context, uuid entities.ProjectUUID) (err error)
	GetProject(ctx context.Context, uuid entities.ProjectUUID) (project *entities.Project, err error)
	AssignProjectFile(ctx context.Context, uuid entities.ProjectUUID, file_id entities.FileID) (err error)
	UnassignProjectVideo(ctx context.Context, uuid entities.ProjectUUID, file_id entities.FileID) (err error)
	GetProjectVideos(ctx context.Context, uuid entities.ProjectUUID) (file_ids []entities.FileID, err error)
	AssignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error)
	UnassignYoutube(ctx context.Context, project_uuid entities.ProjectUUID, youtube_id entities.YoutubeVideoID) (err error)
}

type YoutubeRepository interface {
	NewYoutubeVideo(ctx context.Context, file_id entities.FileID, youtube_video *entities.YoutubeVideo) (err error)
	GetYoutubeVideo(ctx context.Context, youtube_id entities.YoutubeVideoID) (video *entities.YoutubeVideo, err error)
	NewYoutube(ctx context.Context, file_id entities.FileID, youtube *entities.Youtube) (err error)
	GetYoutube(ctx context.Context, youtube_id entities.YoutubeVideoID) (youtube *entities.Youtube, err error)
	GetTitle(ctx context.Context, youtube_id entities.YoutubeVideoID) (title string, err error)
	GetDescription(ctx context.Context, youtube_id entities.YoutubeVideoID) (description string, err error)
	GetChannelVideos(ctx context.Context, channel_id entities.YoutubeChannelID) (videos []entities.YoutubeVideoID, err error)
	GetYoutubeFileIDs(ctx context.Context, youtube_id entities.YoutubeVideoID) (file_ids []entities.FileID, err error)
}

type VideoRepository interface {
	NewVideo(ctx context.Context, youtube_video *entities.Video) (err error)
}

type Repository struct {
	Project ProjectRepository
	Youtube YoutubeRepository
	File    FileRepository
}
