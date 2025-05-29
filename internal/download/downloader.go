package download

import (
	"context"
	"io"

	"github.com/dtbead/wc-maps-archive/internal/entities"
)

type Youtube interface {
	Download(ctx context.Context, url string, output io.Writer) (youtube *entities.Youtube, extension string, err error)
}
