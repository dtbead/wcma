package mock_youtube

import (
	"errors"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/helper"
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

var commonWidthHeight [][2]int = [][2]int{
	{3840, 2160},
	{2048, 1080},
	{1920, 1080},
	{720, 480},
	{480, 360},
	{360, 144},
}

var commonVideoCodec []string = []string{
	"vp9",
	"avc",
	"avc1.64001F",
	"avc1.64001f",
	"avc1.640028",
	"avc1.640028",
}

var commonAudioCodec []string = []string{
	"mp4a.40.2",
	"opus",
}

var commonFPS []int16 = []int16{
	12, 24, 30, 120, 144, 240,
}

// RandomYoutube returns a valid and randomly generated Youtube entity. file_id assigns a file_id to the Youtube entity.
func RandomYoutube(file_id entities.FileID) entities.Youtube {
	duration := helper.RandomInt(1, 300)
	youtube_id := helper.RandomYoutubeID()
	width_height, _ := randomMultiDimensionalSlice(commonWidthHeight)
	video_codec, _ := helper.RandomSlice(commonVideoCodec)
	audio_codec, _ := helper.RandomSlice(commonAudioCodec)
	fps, _ := helper.RandomSlice(commonFPS)

	yt := entities.Youtube{
		YouTube: entities.YoutubeVideo{
			YoutubeID:    youtube_id,
			UploadDate:   helper.RandomTime(time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC), time.Now()),
			Duration:     duration,
			ViewCount:    helper.RandomInt(1, 9999999),
			LikeCount:    helper.RandomInt(1, 999999),
			DislikeCount: helper.RandomInt(1, 99999),
			IsLive:       helper.RandomBool(),
			IsRestricted: helper.RandomBool(),
			Video: entities.Video{
				VideoCodec: video_codec,
				AudioCodec: audio_codec,
				Duration:   duration,
				Width:      int16(width_height[0]),
				Height:     int16(width_height[1]),
				Fps:        fps,
			},
		},
		Channel: &entities.VideoYoutubeChannel{
			ChannelID:  helper.RandomYoutubeChannelID(),
			Uploader:   helper.RandomString(helper.RandomInt(3, 20)),
			UploaderID: "@" + helper.RandomString(helper.RandomInt(6, 20)),
		},
		Format: &entities.VideoYoutubeFormat{
			YoutubeID: youtube_id,
			FileID:    file_id,
			Format:    "136 - 976x720 (720p)+251 - audio only (medium)", // make consist with width/height
			FormatID:  "136+251",                                        // make consist with width/height
		},
		DlpVersion: &entities.VideoYoutubeDlpVersion{
			FileID:         file_id,
			YoutubeID:      youtube_id,
			Version:        "2025.04.30", // TODO: add randomization later when new versions are released
			ReleaseGitHead: "505b400795af557bdcfd9d4fa7e9133b26ef431c",
			Repository:     "yt-dlp/yt-dlp",
		},
		Title:       helper.RandomString(helper.RandomInt(6, 32)),
		Description: helper.RandomString(helper.RandomInt(6, 300)),
	}

	return yt
}

func randomMultiDimensionalSlice(s [][2]int) ([2]int, error) {
	if s == nil || len(s) < 1 {
		return [2]int{}, errors.New("no slice given")
	}

	r := helper.RandomInt(0, len(s))
	return s[r], nil
}
