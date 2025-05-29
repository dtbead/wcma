package youtube

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/dtbead/wc-maps-archive/internal/service"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type YoutubeService struct {
	YoutubeRepository storage.YoutubeRepository
	FileRepository    storage.FileRepository
}

func (y YoutubeService) NewYoutubeProject(ctx context.Context, video io.Reader, extension string, project_youtube service.ProjectYoutube) (err error) {
	switch {
	case project_youtube.Project == nil:
		return errors.New("given nil project")
	case project_youtube.Youtube == nil:
		return errors.New("given nil youtube")
	case !project_youtube.Youtube.YouTube.YoutubeID.IsValid():
		return errors.New("invalid YoutubeID")
	case project_youtube.Youtube.Title == "":
		return errors.New("no youtube title given")
	case extension == "" || len(extension) < 2 || len(extension) > 6:
		return errors.New("invalid extension")
	}

	extension = strings.TrimPrefix(extension, ".")

	file_id, err := y.FileRepository.NewFile(ctx, video, extension)
	if err != nil {
		return err
	}

	err = y.YoutubeRepository.NewYoutube(ctx, file_id, project_youtube.Youtube)
	if err != nil {
		return errors.Join(err, y.FileRepository.DeleteFile(ctx, file_id))
	}

	return nil
}

func (y YoutubeService) DownloadVideo(ctx context.Context, url string, downloader service.YoutubeDownloader) (err error) {
	f, err := y.FileRepository.NewTempFile(ctx)
	if err != nil {
		return err
	}
	defer f.Close()

	yt, ext, err := downloader.Download(ctx, url, f)
	if err != nil {
		return err
	}

	file_id, err := y.FileRepository.NewFile(ctx, f, ext)
	if err != nil {
		return err
	}

	err = y.YoutubeRepository.NewYoutube(ctx, file_id, yt)
	if err != nil {
		return errors.Join(err, y.FileRepository.DeleteFile(ctx, file_id))
	}

	return nil
}
