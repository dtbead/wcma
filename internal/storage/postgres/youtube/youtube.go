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

func (y YoutubeRepository) GetYoutubeVideo(ctx context.Context, youtube_id entities.YoutubeVideoID) (youtube_video *entities.YoutubeVideo, err error) {
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

	err = y.q.NewFileVideo(ctx, queries.NewFileVideoParams{
		FileID:     int64(file_id),
		Duration:   int32(youtube.YouTube.Duration),
		Width:      youtube.YouTube.Video.Width,
		Height:     youtube.YouTube.Video.Height,
		Fps:        sql.NullInt16{Int16: youtube.YouTube.Video.Fps, Valid: youtube.YouTube.Video.Fps > 0},
		VideoCodec: sql.NullString{String: youtube.YouTube.Video.VideoCodec, Valid: len(youtube.YouTube.Video.VideoCodec) >= 3},
		AudioCodec: sql.NullString{String: youtube.YouTube.Video.AudioCodec, Valid: len(youtube.YouTube.Video.AudioCodec) >= 3},
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

func (y YoutubeRepository) GetYoutube(ctx context.Context, youtube_id entities.YoutubeVideoID) (youtube *entities.Youtube, err error) {
	file_id, err := y.q.GetYoutubeFileID(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	if file_id == nil || len(file_id) < 1 {
		return nil, errors.New("found no file_id for given youtube_id")
	}

	youtube_video, err := y.q.GetYoutubeVideo(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	file_video, err := y.q.GetFileVideo(ctx, file_id[0])
	if err != nil {
		return nil, err
	}

	yt := &entities.Youtube{
		YouTube: entities.YoutubeVideo{
			YoutubeID: youtube_id,
			Video: entities.Video{
				VideoCodec: file_video.VideoCodec.String,
				AudioCodec: file_video.AudioCodec.String,
				Duration:   int(file_video.Duration),
				Width:      file_video.Width,
				Height:     file_video.Height,
				Fps:        file_video.Fps.Int16,
			},
			UploadDate:   youtube_video.UploadDate,
			Duration:     int(youtube_video.Duration),
			ViewCount:    int(youtube_video.ViewCount.Int32),
			LikeCount:    int(youtube_video.LikeCount.Int32),
			DislikeCount: int(youtube_video.DislikeCount.Int32),
			IsLive:       youtube_video.IsLive.Bool,
			IsRestricted: youtube_video.IsRestricted.Bool,
		},
	}

	title, _ := y.q.GetYoutubeTitle(ctx, youtube_id)
	if title != nil && len(title) > 0 {
		yt.Title = title[0]
	}

	description, _ := y.q.GetYoutubeDescription(ctx, youtube_id)
	if description != nil && len(description) > 0 {
		yt.Description = description[0]
	}

	return yt, nil
}

func (y YoutubeRepository) GetYoutubeFileIDs(ctx context.Context, youtube_id entities.YoutubeVideoID) (file_ids []entities.FileID, err error) {
	if !youtube_id.IsValid() {
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

func (y YoutubeRepository) GetYtdlpVersion(ctx context.Context, youtube_id entities.YoutubeVideoID, file_id entities.FileID) (version *entities.VideoYoutubeDlpVersion, err error) {
	res, err := y.q.GetYoutubeYtdlpVersion(ctx, queries.GetYoutubeYtdlpVersionParams{
		YoutubeID: youtube_id,
		FileID:    int64(file_id),
	})
	if err != nil {
		return nil, err
	}

	return &entities.VideoYoutubeDlpVersion{
		YoutubeID:      youtube_id,
		FileID:         entities.FileID(res.FileID),
		Version:        res.Version,
		ReleaseGitHead: res.ReleaseGitHead,
		Repository:     res.Repository,
	}, nil
}

func (y YoutubeRepository) GetFormat(ctx context.Context, youtube_id entities.YoutubeVideoID) (format *entities.VideoYoutubeFormat, err error) {
	if !youtube_id.IsValid() {
		return nil, errors.New("invalid youtube_id")
	}

	res, err := y.q.GetYoutubeVideoFormatByYoutubeID(ctx, youtube_id)
	if err != nil {
		return nil, err
	}

	if res == nil || len(res) < 1 {
		return nil, errors.New("no format found")
	}

	return &entities.VideoYoutubeFormat{
		YoutubeID: youtube_id,
		FileID:    entities.FileID(res[0].FileID),
		Format:    res[0].Format,
		FormatID:  res[0].FormatID,
	}, nil
}

func (y YoutubeRepository) GetTitle(ctx context.Context, youtube_id entities.YoutubeVideoID) (title string, err error) {
	if !youtube_id.IsValid() {
		return "", entities.ErrorInvalidYoutubeID
	}

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

	if len(desc) < 1 {
		return "", errors.New("no description found")
	}

	return desc[0], nil
}
