package youtube

import (
	"context"
	"errors"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type YoutubeService struct {
	YoutubeRepository storage.YoutubeRepository
}

func NewService(YoutubeRepo storage.YoutubeRepository) *YoutubeService {
	return &YoutubeService{YoutubeRepo}
}

func IsValidYoutube(y *entities.Youtube) error {
	if y == nil {
		return errors.New("given nil youtube")
	}

	switch {
	case !y.YouTube.YoutubeID.IsValid():
		return entities.ErrorInvalidYoutubeID
	case y.YouTube.UploadDate.IsZero():
		return errors.New("invalid UploadDate")
	case y.YouTube.Video.Width < 15 || y.YouTube.Video.Height < 15:
		return errors.New("invalid video width/height")
	case y.YouTube.Duration < 1:
		return errors.New("invalid youtube duration")
	case y.YouTube.Video.Duration < 1:
		return errors.New("invalid video duration")
	case y.YouTube.Video.Duration != y.YouTube.Duration:
		return errors.New("video duration != youtube video duration")
	case y.YouTube.Video.Fps < 1:
		return errors.New("invalid video fps")
	case y.YouTube.ViewCount < 0:
		return errors.New("invalid view count")
	case y.YouTube.DislikeCount < 0:
		return errors.New("invalid dislike count")
	case y.YouTube.LikeCount < 0:
		return errors.New("invalid like count")
	}

	if y.Channel != nil && !y.Channel.ChannelID.IsValid() {
		return entities.ErrorInvalidYoutubeChannelID
	}

	return nil
}

func (y YoutubeService) NewYoutube(ctx context.Context, file_id entities.FileID, youtube *entities.Youtube) (err error) {
	err = IsValidYoutube(youtube)
	if err != nil {
		return err
	}

	if file_id < 1 {
		return errors.New("invalid file_id")
	}

	return y.YoutubeRepository.NewYoutube(ctx, file_id, youtube)
}

func (y YoutubeService) GetYoutube(ctx context.Context, youtube_id entities.YoutubeVideoID) (youtube *entities.Youtube, err error) {
	if !youtube_id.IsValid() {
		return nil, entities.ErrorInvalidYoutubeID
	}

	return y.YoutubeRepository.GetYoutube(ctx, youtube_id)
}

func (y YoutubeService) GetYoutubeFileIDs(ctx context.Context, youtube_id entities.YoutubeVideoID) (file_ids []entities.FileID, err error) {
	if !youtube_id.IsValid() {
		return nil, entities.ErrorInvalidYoutubeID
	}

	return y.YoutubeRepository.GetYoutubeFileIDs(ctx, youtube_id)
}

func (y YoutubeService) GetTitle(ctx context.Context, youtube_id entities.YoutubeVideoID) (title string, err error) {
	if !youtube_id.IsValid() {
		return "", entities.ErrorInvalidYoutubeID
	}
	return y.YoutubeRepository.GetTitle(ctx, youtube_id)
}
func (y YoutubeService) GetDescription(ctx context.Context, youtube_id entities.YoutubeVideoID) (description string, err error) {
	if !youtube_id.IsValid() {
		return "", entities.ErrorInvalidYoutubeID
	}
	return y.YoutubeRepository.GetDescription(ctx, youtube_id)
}
func (y YoutubeService) GetChannelVideos(ctx context.Context, channel_id entities.YoutubeChannelID) (videos []entities.YoutubeVideoID, err error) {
	if !channel_id.IsValid() {
		return nil, entities.ErrorInvalidYoutubeChannelID
	}
	return y.YoutubeRepository.GetChannelVideos(ctx, channel_id)
}
func (y YoutubeService) GetYoutubeVideo(ctx context.Context, youtube_id entities.YoutubeVideoID) (video entities.YoutubeVideo, err error) {
	if !youtube_id.IsValid() {
		return entities.YoutubeVideo{}, entities.ErrorInvalidYoutubeID
	}

	v, err := y.YoutubeRepository.GetYoutubeVideo(ctx, youtube_id)
	if err != nil {
		return entities.YoutubeVideo{}, err
	}

	return *v, nil

}
