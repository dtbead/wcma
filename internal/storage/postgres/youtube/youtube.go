package youtube

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/helper"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/queries"
)

type YoutubeRepository struct {
	db *sql.DB
	q  *queries.Queries
}

func NewYoutubeRepository(db *sql.DB) *YoutubeRepository {
	return &YoutubeRepository{
		db: db,
		q:  queries.New(db),
	}
}

func (y YoutubeRepository) NewYoutubeVideo(ctx context.Context, file_id entities.FileID, youtube_video *entities.YoutubeVideo) (err error) {
	if youtube_video == nil {
		return errors.New("nil youtube_video given")
	}

	tx, err := y.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	y.q = queries.New(tx)

	err = y.q.NewYoutube(ctx, queries.NewYoutubeParams{
		ID:           youtube_video.YoutubeID,
		UploadDate:   youtube_video.UploadDate,
		Duration:     int32(youtube_video.Duration),
		ViewCount:    sql.NullInt32{Int32: int32(youtube_video.ViewCount), Valid: youtube_video.ViewCount > 0},
		LikeCount:    sql.NullInt32{Int32: int32(youtube_video.LikeCount), Valid: youtube_video.LikeCount > 0},
		DislikeCount: sql.NullInt32{Int32: int32(youtube_video.DislikeCount), Valid: youtube_video.DislikeCount > 0},
		IsLive:       sql.NullBool{Valid: true, Bool: youtube_video.IsLive},
		IsRestricted: sql.NullBool{Valid: true, Bool: youtube_video.IsRestricted},
	})
	if err != nil {
		return err
	}

	err = y.q.AssignYoutubeFileID(ctx, queries.AssignYoutubeFileIDParams{
		YoutubeID: youtube_video.YoutubeID,
		FileID:    int64(file_id),
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (y YoutubeRepository) GetYoutube(ctx context.Context, youtube_id entities.YoutubeVideoID) (youtube_video *entities.YoutubeVideo, err error) {
	res, err := y.q.GetYoutubeVideo(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	yt := entities.YoutubeVideo{
		YoutubeID:    youtube_id,
		UploadDate:   res.UploadDate,
		Duration:     int(res.Duration),
		IsLive:       res.IsLive.Bool,
		IsRestricted: res.IsRestricted.Bool,
	}

	switch {
	case res.ViewCount.Valid:
		yt.ViewCount = int(res.ViewCount.Int32)
		fallthrough
	case res.LikeCount.Valid:
		yt.LikeCount = int(res.LikeCount.Int32)
		fallthrough
	case res.DislikeCount.Valid:
		yt.DislikeCount = int(res.DislikeCount.Int32)
	}

	return &yt, nil
}

func (y YoutubeRepository) NewYoutube(ctx context.Context, file_id entities.FileID, youtube *entities.Youtube) (err error) {
	if youtube == nil {
		return errors.New("nil youtube given")
	}

	tx, err := y.db.Begin()
	if err != nil {
		return err
	}
	y.q = y.q.WithTx(tx)
	defer tx.Rollback()

	err = y.q.NewYoutube(ctx, queries.NewYoutubeParams{
		ID:           youtube.YouTube.YoutubeID,
		UploadDate:   youtube.YouTube.UploadDate,
		Duration:     int32(youtube.YouTube.Duration),
		ViewCount:    sql.NullInt32{Int32: int32(youtube.YouTube.ViewCount), Valid: youtube.YouTube.ViewCount > 0},
		LikeCount:    sql.NullInt32{Int32: int32(youtube.YouTube.LikeCount), Valid: youtube.YouTube.LikeCount > 0},
		DislikeCount: sql.NullInt32{Int32: int32(youtube.YouTube.LikeCount), Valid: youtube.YouTube.DislikeCount > 0},
		IsLive:       sql.NullBool{Valid: true, Bool: youtube.YouTube.IsLive},
		IsRestricted: sql.NullBool{Valid: true, Bool: youtube.YouTube.IsRestricted},
	})
	if err != nil {
		return err
	}

	err = y.q.AssignYoutubeTitle(ctx, queries.AssignYoutubeTitleParams{
		YoutubeID: youtube.YouTube.YoutubeID,
		Title:     youtube.Title,
		TitleMd5:  helper.GetMD5HashFromString(youtube.Title),
	})
	if err != nil {
		return err
	}

	err = y.q.AssignYoutubeDescription(ctx, queries.AssignYoutubeDescriptionParams{
		YoutubeID:      youtube.YouTube.YoutubeID,
		Description:    youtube.Description,
		DescriptionMd5: helper.GetMD5HashFromString(youtube.Description),
	})
	if err != nil {
		return err
	}

	err = y.q.AssignYoutubeFileID(ctx, queries.AssignYoutubeFileIDParams{YoutubeID: youtube.YouTube.YoutubeID, FileID: int64(file_id)})
	if err != nil {
		return err
	}

	if youtube.Channel != nil {
		err = y.q.NewYoutubeChannel(ctx, youtube.Channel.ChannelID)
		if err != nil {
			return err
		}

		err = y.q.NewYoutubeChannelVideo(ctx, queries.NewYoutubeChannelVideoParams{
			ChannelID: youtube.Channel.ChannelID,
			YoutubeID: youtube.YouTube.YoutubeID,
		})
		if err != nil {
			return err
		}

		err = y.q.NewYoutubeChannelUploaderID(ctx, queries.NewYoutubeChannelUploaderIDParams{
			ChannelID:  youtube.Channel.ChannelID,
			UploaderID: youtube.Channel.UploaderID,
		})
		if err != nil {
			return err
		}

		err = y.q.NewYoutubeChannelUploaderName(ctx, queries.NewYoutubeChannelUploaderNameParams{
			ChannelID: youtube.Channel.ChannelID,
			Uploader:  youtube.Channel.Uploader,
		})
		if err != nil {
			return err
		}
	}

	if youtube.Format != nil {
		err = y.q.NewYoutubeFormat(ctx, queries.NewYoutubeFormatParams{
			YoutubeID: youtube.YouTube.YoutubeID,
			FileID:    int64(file_id),
			Format:    youtube.Format.Format,
			FormatID:  youtube.Format.FormatID,
		})
		if err != nil {
			return err
		}
	}

	if youtube.DlpVersion != nil {
		err = y.q.NewYoutubeYtdlpVersion(ctx, queries.NewYoutubeYtdlpVersionParams{
			FileID:         int64(file_id),
			YoutubeID:      youtube.YouTube.YoutubeID,
			Repository:     youtube.DlpVersion.Repository,
			ReleaseGitHead: youtube.DlpVersion.ReleaseGitHead,
			Version:        youtube.DlpVersion.Version,
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (y YoutubeRepository) GetYoutubeFileIDs(ctx context.Context, youtube_id entities.YoutubeVideoID) (file_ids []entities.FileID, err error) {
	if youtube_id == "" {
		return nil, errors.New("invalid youtube_id given")
	}

	res, err := y.q.GetYoutubeFileID(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	file_ids = make([]entities.FileID, 0, len(res))
	for _, v := range res {
		file_ids = append(file_ids, entities.FileID(v))
	}

	return file_ids, nil
}

func (y YoutubeRepository) GetYoutubeFull(ctx context.Context, youtube_id entities.YoutubeVideoID) (youtube *entities.Youtube, err error) {
	if youtube_id == "" {
		return nil, errors.New("invalid youtube_id given")
	}

	yt, err := y.GetYoutube(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	title, err := y.GetTitle(ctx, youtube_id) // TODO: ignore error on empty row
	if err != nil {
		return nil, err
	}

	description, err := y.GetDescription(ctx, youtube_id) // TODO: ignore error on empty row
	if err != nil {
		return nil, err
	}

	channel, err := y.GetChannelByVideoID(ctx, youtube_id) // TODO: ignore error on empty row
	if err != nil {
		return nil, err
	}

	format, err := y.GetFormat(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	version, err := y.GetYtdlpVersion(ctx, youtube_id, format.FileID)
	if err != nil {
		return nil, err
	}

	youtube = new(entities.Youtube)
	youtube.Title = title
	youtube.Description = description
	youtube.Channel = &channel
	youtube.YouTube = *yt
	youtube.Format = &format
	youtube.DlpVersion = &version

	return youtube, nil
}

func (y YoutubeRepository) GetProjectType(ctx context.Context, youtube_id entities.YoutubeVideoID) (project entities.ProjectType, err error) {
	res, err := y.q.GetProjectTypeByYoutubeID(ctx, youtube_id)
	if err != nil {
		return entities.ProjectTypeUnknown, err
	}

	p, err := entities.NewProjectType(string(res))
	if err != nil {
		return entities.ProjectTypeUnknown, err
	}

	return p, nil
}

func (y YoutubeRepository) GetYtdlpVersion(ctx context.Context, youtube_id entities.YoutubeVideoID, file_id entities.FileID) (version entities.VideoYoutubeDlpVersion, err error) {
	res, err := y.q.GetYoutubeYtdlpVersion(ctx, queries.GetYoutubeYtdlpVersionParams{
		YoutubeID: youtube_id,
		FileID:    int64(file_id),
	})
	if err != nil {
		return entities.VideoYoutubeDlpVersion{}, err
	}

	return entities.VideoYoutubeDlpVersion{
		YoutubeID:      youtube_id,
		FileID:         entities.FileID(res.FileID),
		Version:        res.Version,
		ReleaseGitHead: res.ReleaseGitHead,
		Repository:     res.Repository,
	}, nil
}

func (y YoutubeRepository) GetFormat(ctx context.Context, youtube_id entities.YoutubeVideoID) (format entities.VideoYoutubeFormat, err error) {
	res, err := y.q.GetYoutubeVideoFormatByYoutubeID(ctx, youtube_id)
	if err != nil {
		return entities.VideoYoutubeFormat{}, err
	}

	return entities.VideoYoutubeFormat{
		YoutubeID: youtube_id,
		FileID:    entities.FileID(res[0].FileID),
		Format:    res[0].Format,
		FormatID:  res[0].FormatID,
	}, nil
}

func (y YoutubeRepository) GetTitle(ctx context.Context, youtube_id entities.YoutubeVideoID) (title string, err error) {
	t, err := y.q.GetYoutubeTitle(ctx, youtube_id)
	if err != nil {
		return "", err
	}

	if len(t) < 1 {
		return "", errors.New("no title found")
	}

	return t[0], nil
}

func (y YoutubeRepository) GetChannelByVideoID(ctx context.Context, youtube_id entities.YoutubeVideoID) (channel entities.VideoYoutubeChannel, err error) {
	c, err := y.q.GetYoutubeChannelByID(ctx, youtube_id)
	if err != nil {
		return entities.VideoYoutubeChannel{}, err
	}

	return entities.VideoYoutubeChannel{
		ChannelID:  entities.YoutubeChannelID(c.ChannelID.(string)),
		Uploader:   c.UploaderName,
		UploaderID: c.UploaderID,
	}, nil
}

func (y YoutubeRepository) GetDescription(ctx context.Context, youtube_id entities.YoutubeVideoID) (description string, err error) {
	desc, err := y.q.GetYoutubeDescription(ctx, youtube_id)
	if err != nil {
		return "", err
	}

	return desc[0], nil
}
