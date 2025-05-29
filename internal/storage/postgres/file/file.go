package file

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"os"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	file_helper "github.com/dtbead/wc-maps-archive/internal/helper/file"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/queries"
)

type FileRepository struct {
	db            *sql.DB
	q             *queries.Queries
	baseDirectory string
}

func NewFileRepository(db *sql.DB, base_directory string) (*FileRepository, error) {
	if base_directory == "" {
		return nil, errors.New("invalid base_directory string")
	}

	return &FileRepository{
		db:            db,
		q:             queries.New(db),
		baseDirectory: file_helper.SanitizePath(base_directory),
	}, nil
}

// entities.File.Hashes gets ignored implicitly, with NewFile using file io.Reader to generate hashes. TODO: add check whether
// file already exists in database, and in folder
func (f FileRepository) NewFile(ctx context.Context, file io.Reader, extension string) (file_id entities.FileID, err error) {
	if file == nil {
		return entities.InvalidFileID, entities.ErrorInvalidFilePtr
	}

	hashes, read, err := file_helper.GetHash(file)
	if err != nil {
		return entities.InvalidFileID, err
	}

	if read < 16 {
		return entities.InvalidFileID, errors.New("read less than 16 bytes from file io.Reader")
	}

	path_relative := file_helper.BuildPath(hashes.SHA256, extension)
	path_absolute := f.baseDirectory + "/" + path_relative
	if file_helper.DoesPathExist(path_absolute) {
		return entities.InvalidFileID, errors.New("path already exists")
	}

	r, err := f.q.NewFile(ctx, queries.NewFileParams{
		Path:      path_relative,
		Extension: extension,
		Md5:       hashes.MD5,
		Sha1:      hashes.SHA1,
		Sha256:    hashes.SHA256,
		Filesize:  read,
	})
	if err != nil {
		return entities.InvalidFileID, err
	}

	file_helper.ResetFileSeek(file)
	err = file_helper.Copy(path_absolute, file)
	if err != nil {
		return entities.InvalidFileID, err
	}

	return entities.FileID(r), nil
}

func (f FileRepository) DeleteFile(ctx context.Context, file_id entities.FileID) (err error) {
	meta, err := f.GetFile(ctx, file_id)
	if err != nil {
		return err
	}

	err = os.Remove(f.baseDirectory + "/" + meta.PathRelative)
	if err != nil {
		return err
	}

	err = f.q.DeleteFileByID(ctx, int64(file_id))
	if err != nil {
		return err
	}

	return nil
}

func (f FileRepository) GetFile(ctx context.Context, file_id entities.FileID) (file_metadata *entities.File, err error) {
	meta, err := f.q.GetFileByID(ctx, int64(file_id))
	if err != nil {
		return nil, err
	}

	return &entities.File{
		PathRelative: meta.Path,
		PathAbsolute: f.baseDirectory + "/" + meta.Path,
		Extension:    meta.Extension,
		Size:         meta.Filesize,
		Hashes: entities.Hashes{
			MD5:    meta.Md5,
			SHA1:   meta.Sha1,
			SHA256: meta.Sha256,
		},
	}, nil
}

func (f FileRepository) NewTempFile(ctx context.Context) (file io.ReadWriteCloser, err error) {
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}

	return tempFile{f: tmp}, nil
}

type tempFile struct {
	f *os.File
}

func (t tempFile) Read(p []byte) (n int, err error) {
	return t.f.Read(p)
}

func (t tempFile) Write(p []byte) (n int, err error) {
	return t.f.Write(p)
}

func (t tempFile) Close() error {
	err := t.f.Close()
	return errors.Join(err, os.Remove(t.f.Name()))
}
