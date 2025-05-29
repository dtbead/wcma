package artist

import (
	"context"
	"database/sql"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/queries"
)

type ArtistRepository struct {
	db *sql.DB
	q  *queries.Queries
}

func NewArtistRepository(db *sql.DB) *ArtistRepository {
	return &ArtistRepository{
		db: db,
		q:  queries.New(db),
	}
}

func (a ArtistRepository) NewArtist(ctx context.Context, artist_name string) (err error) {
	panic("unimplemented")
}

func (a ArtistRepository) DeleteArtist(ctx context.Context, artist_name string) (err error) {
	panic("unimplemented")
}

func (a ArtistRepository) AssignArtistChannel(ctx context.Context, artist_name string, channel_id entities.YoutubeChannelID) (err error) {
	panic("unimplemented")
}

func (a ArtistRepository) UnassignArtistChannel(ctx context.Context, artist_name string, channel_id entities.YoutubeChannelID) (err error) {
	panic("unimplemented")
}
