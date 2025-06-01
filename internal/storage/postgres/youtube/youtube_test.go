package youtube_test

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/helper/mock"
	helper_test "github.com/dtbead/wc-maps-archive/internal/helper/test"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/file"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/youtube"
	"github.com/google/go-cmp/cmp"
)

// helperInsertFile imports a test video file into FileRepository and returns a corrasponding file_id.
// helperInsertFile will implicicly call t.Fatal if any error occurs when trying to import a file.
func helperInsertFile(fileRepo file.FileRepository, t *testing.T) (file_id entities.FileID) {
	t.Helper()

	f := helperOpenTestFile(t)
	defer f.Close()

	file_id, err := fileRepo.NewFile(context.Background(), f, "mkv")
	if err != nil {
		t.Fatalf("failed to insert file to db, %v", err)
	}

	return file_id
}

func TestYoutubeRepository_NewYoutubeVideo(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}

	file_id := helperInsertFile(*fileRepo, t)

	type args struct {
		ctx           context.Context
		file_id       entities.FileID
		youtube_video *entities.YoutubeVideo
	}
	tests := []struct {
		name    string
		y       youtube.YoutubeRepository
		args    args
		wantErr bool
	}{
		{"nil youtube video", *youtubeRepo, args{context.Background(), file_id, nil}, true},
		{"invalid/missing file_id insert", *youtubeRepo, args{context.Background(), file_id + 9999, &entities.YoutubeVideo{
			YoutubeID:    "y_wo8pyoxyk",
			UploadDate:   time.Unix(1745372790, 0),
			Duration:     7,
			ViewCount:    326,
			LikeCount:    26,
			DislikeCount: 0,
			IsLive:       false,
			IsRestricted: false,
		}}, true},
		{"valid insert", *youtubeRepo, args{context.Background(), file_id, &entities.YoutubeVideo{
			YoutubeID:    "y_wo8pyoxyk",
			UploadDate:   time.Unix(1745372790, 0),
			Duration:     7,
			ViewCount:    326,
			LikeCount:    26,
			DislikeCount: 0,
			IsLive:       false,
			IsRestricted: false,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.y.NewYoutubeVideo(tt.args.ctx, tt.args.file_id, tt.args.youtube_video)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.NewYoutube() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.args.youtube_video != nil && err == nil {
				file_ids, err := tt.y.GetYoutubeFileIDs(tt.args.ctx, tt.args.youtube_video.YoutubeID)
				if err != nil {
					t.Errorf("failed to get file_ids from %s,  %v", tt.args.youtube_video.YoutubeID, err)
				}

				if !slices.Contains(file_ids, tt.args.file_id) {
					t.Errorf("expect file_id '%d', got %v", tt.args.file_id, file_ids)
				}
			}
		})
	}
}

func TestYoutubeRepository_GetYoutubeVideo(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}

	file_id := helperInsertFile(*fileRepo, t)

	yt := entities.YoutubeVideo{
		YoutubeID:    "y_wo8pyoxyk",
		UploadDate:   time.Unix(1745372790, 0).UTC().Round(time.Second),
		Duration:     7,
		ViewCount:    326,
		LikeCount:    26,
		IsLive:       false,
		IsRestricted: false,
	}

	type args struct {
		ctx        context.Context
		file_id    entities.FileID
		youtube_id entities.YoutubeVideoID
	}
	tests := []struct {
		name              string
		y                 youtube.YoutubeRepository
		args              args
		wantYoutube_video *entities.YoutubeVideo
		wantErr           bool
	}{
		{"youtube_id not found", *youtubeRepo, args{context.Background(), file_id, "abcdefghj"}, nil, true},
		{"valid youtube_id", *youtubeRepo, args{context.Background(), file_id, "y_wo8pyoxyk"}, &yt, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantYoutube_video != nil {
				err := tt.y.NewYoutubeVideo(tt.args.ctx, tt.args.file_id, tt.wantYoutube_video)
				if (err != nil) != tt.wantErr {
					t.Errorf("YoutubeRepository.GetYoutube() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			gotYoutube_video, err := tt.y.GetYoutube(tt.args.ctx, tt.args.youtube_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetYoutube() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(gotYoutube_video, tt.wantYoutube_video) {
				t.Errorf("YoutubeRepository.GetYoutube() diff = %s", cmp.Diff(gotYoutube_video, tt.wantYoutube_video))
			}
		})
	}
}

func TestYoutubeRepository_NewYoutube(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}

	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	type args struct {
		ctx     context.Context
		file_id entities.FileID
		youtube *entities.Youtube
	}
	tests := []struct {
		name    string
		y       youtube.YoutubeRepository
		args    args
		wantErr bool
	}{
		{"nil youtube", *youtubeRepo, args{}, true},
		{"invalid file_id", *youtubeRepo, args{context.Background(), 0, &mockYt}, true},
		{"valid insert", *youtubeRepo, args{context.Background(), file_id, &mockYt}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.y.NewYoutube(tt.args.ctx, tt.args.file_id, tt.args.youtube); (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.NewYoutubeFull() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				yt, err := tt.y.GetYoutubeFull(tt.args.ctx, tt.args.youtube.YouTube.YoutubeID)
				if err != nil {
					t.Errorf("failed to get youtube metadata from db, %v", err)
					return
				}

				if !cmp.Equal(*yt, *tt.args.youtube) {
					t.Errorf("got youtube metadata != want youtube metadata, %v", cmp.Diff(*yt, *tt.args.youtube))
				}
			}

		})
	}
}

var mockYt = &entities.Youtube{
	YouTube: entities.YoutubeVideo{
		YoutubeID:    "y_wo8pyoxyk",
		UploadDate:   time.Unix(1745372790, 0).UTC().Round(time.Second),
		Duration:     7,
		ViewCount:    326,
		LikeCount:    26,
		IsLive:       false,
		IsRestricted: false,
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

// helperOpenTestFile opens a new video file for testing. helperOpenTestFile will implicicly call t.Fatalf if it fails to open a test file.
// Caller is responsible for closing the file when they're finished.
func helperOpenTestFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err)
	}
	return f
}
