package mock

import (
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
)

// NewYoutube returns a valid Youtube entity that would be expected from a typical youtube-dlp extraction.
func NewYoutube() entities.Youtube {
	return entities.Youtube{
		YouTube: entities.YoutubeVideo{
			YoutubeID:    "y_wo8pyoxyk",
			UploadDate:   time.Unix(1745372790, 0).UTC().Round(time.Second),
			Duration:     7,
			ViewCount:    326,
			LikeCount:    26,
			IsLive:       false,
			IsRestricted: false,
			Video: entities.Video{
				VideoCodec: "avc1.64001f",
				AudioCodec: "opus",
				Duration:   7,
				Width:      976,
				Height:     720,
				Fps:        30,
			},
		},
		Channel: &entities.VideoYoutubeChannel{
			ChannelID:  "UCKhKck7AoDI-H8PktMnZi0Q",
			UploaderID: "@glassfirestar",
			Uploader:   "Rusty",
		},
		Format: &entities.VideoYoutubeFormat{
			YoutubeID: "y_wo8pyoxyk",
			FileID:    1,
			Format:    "136 - 976x720 (720p)+251 - audio only (medium)",
			FormatID:  "136+251",
		},
		DlpVersion: &entities.VideoYoutubeDlpVersion{
			FileID:         1,
			YoutubeID:      "y_wo8pyoxyk",
			Version:        "2025.04.30",
			ReleaseGitHead: "505b400795af557bdcfd9d4fa7e9133b26ef431c",
			Repository:     "yt-dlp/yt-dlp",
		},

		Title:       "does he know about the dore",
		Description: "desc",
	}
}
