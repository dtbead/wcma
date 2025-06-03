package file

import (
	"context"
	"io"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage"
)

type FileService struct {
	FileRepo *storage.FileRepository
}

func NewService(FileRepo storage.FileRepository) *FileService {
	return &FileService{FileRepo: &FileRepo}
}

func (f FileService) NewFile(ctx context.Context, file io.Reader, extension string) (file_id entities.FileID, err error) {
	panic("unimplemented")
}
func (f FileService) DeleteFile(ctx context.Context, file_id entities.FileID) (err error) {
	panic("unimplemented")
}

func (f FileService) NewTempFile(ctx context.Context) (file io.ReadWriteCloser, err error) {
	return f.NewTempFile(ctx)
}
