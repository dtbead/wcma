package youtube

import (
	"context"
	"io"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type YoutubeService struct {
	YoutubeRepository storage.YoutubeRepository
}

func NewService(YoutubeRepo storage.YoutubeRepository) *YoutubeService {
	return &YoutubeService{YoutubeRepo}
}

func (y YoutubeService) NewYoutube(ctx context.Context, video io.Reader, project_youtube *entities.ProjectYoutube) (err error) {
	panic("unimplemented")
}
