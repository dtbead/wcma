package ytdlp

import (
	"encoding/json"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
)

func newMetadata(ytdlpJSON []byte) (*metadata, error) {
	m := new(metadata)
	err := json.Unmarshal(ytdlpJSON, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

type metadata struct {
	Id            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Fulltitle     string  `json:"fulltitle"`
	Channel_id    string  `json:"channel_id"`
	Uploader      string  `json:"uploader"`
	Uploader_id   string  `json:"uploader_id"`
	Upload_date   string  `json:"upload_date"`
	Availability  string  `json:"availability"`
	Format        string  `json:"format"`
	Format_id     string  `json:"format_id"`
	Extension     string  `json:"ext"`
	Resolution    string  `json:"resolution"`
	Vcodec        string  `json:"vcodec"`
	Acodec        string  `json:"acodec"`
	AudioBitrate  float64 `json:"abr"`
	Vbr           float64 `json:"vbr"`
	Upload_epoch  int     `json:"timestamp"`
	Duration      int     `json:"duration"`
	View_count    int     `json:"view_count"`
	Comment_count int     `json:"comment_count"`
	Like_count    int     `json:"like_count"`
	Dislike_count int     `json:"dislike_count"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Fps           float64 `json:"fps"`
	Is_live       bool    `json:"is_live"`
	Age_limit     int     `json:"age_limit"`
	Version       struct {
		Version          string `json:"version"`
		Release_git_head string `json:"release_git_head"`
		Repository       string `json:"repository"`
	} `json:"_version"`
}

func (m metadata) ToYoutubeEntity() entities.Youtube {
	yt := entities.Youtube{
		YouTube: entities.YoutubeVideo{
			YoutubeID:    entities.YoutubeVideoID(m.Id),
			UploadDate:   time.Unix(int64(m.Upload_epoch), 0),
			Duration:     m.Duration,
			ViewCount:    m.View_count,
			LikeCount:    m.Like_count,
			DislikeCount: m.Dislike_count,
			IsLive:       m.Is_live,
			IsRestricted: m.Age_limit > 0,
			Video: entities.Video{
				VideoCodec: m.Vcodec,
				AudioCodec: m.Acodec,
				Fps:        int16(m.Fps),
				Width:      int16(m.Width),
				Height:     int16(m.Height),
			},
		},
		Channel: &entities.VideoYoutubeChannel{
			ChannelID:  entities.YoutubeChannelID(m.Channel_id),
			UploaderID: m.Uploader_id,
			Uploader:   m.Uploader,
		},
		Format: &entities.VideoYoutubeFormat{
			YoutubeID: entities.YoutubeVideoID(m.Id),
			FileID:    entities.InvalidFileID,
			Format:    m.Format,
			FormatID:  m.Format_id,
		},
		DlpVersion: &entities.VideoYoutubeDlpVersion{
			YoutubeID:      entities.YoutubeVideoID(m.Id),
			Version:        m.Version.Version,
			ReleaseGitHead: m.Version.Release_git_head,
			Repository:     m.Version.Repository,
		},
		Title:       m.Title,
		Description: m.Description,
	}

	return yt
}
