package postgres

import (
	"database/sql"
	_ "embed"

	"github.com/dtbead/wc-maps-archive/internal/storage"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/file"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/project"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/youtube"
)

//go:embed schema.sql
var Schema string

func NewRepository(db *sql.DB, base_directory string) (*storage.Repository, error) {
	f, err := file.NewFileRepository(db, base_directory)
	if err != nil {
		return nil, err
	}

	return &storage.Repository{
		Project: project.NewProjectRepository(db),
		Youtube: youtube.NewYoutubeRepository(db),
		File:    f,
	}, nil
}
