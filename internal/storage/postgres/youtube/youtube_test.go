package youtube_test

import (
	"context"
	"os"
	"reflect"
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

			gotYoutube_video, err := tt.y.GetYoutubeVideo(tt.args.ctx, tt.args.youtube_id)
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
				gotYoutube, err := tt.y.GetYoutube(tt.args.ctx, tt.args.youtube.YouTube.YoutubeID)
				if err != nil {
					t.Errorf("GetYoutube err = %v", err)
					return
				}

				channel, err := tt.y.GetChannelByVideoID(tt.args.ctx, tt.args.youtube.YouTube.YoutubeID)
				if err != nil {
					t.Errorf("GetChannelByVideoID err = %v", err)
				} else {
					gotYoutube.Channel = &channel
				}

				format, err := tt.y.GetFormat(tt.args.ctx, tt.args.youtube.YouTube.YoutubeID)
				if err != nil {
					t.Errorf("GetFormat err = %v", err)
				} else {
					gotYoutube.Format = format
				}

				version, err := tt.y.GetYtdlpVersion(tt.args.ctx, tt.args.youtube.YouTube.YoutubeID, tt.args.file_id)
				if err != nil {
					t.Errorf("GetYtdlpVersion err = %v", err)
				} else {
					gotYoutube.DlpVersion = version
				}

				if !cmp.Equal(*gotYoutube, *tt.args.youtube) {
					t.Errorf("got youtube metadata != want youtube metadata, %v", cmp.Diff(*gotYoutube, *tt.args.youtube))
				}
			}

		})
	}
}

func TestYoutubeRepository_GetYoutubeFileIDs(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}

	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()
	if err := youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt); err != nil {
		t.Fatalf("failed to insert mock youtube, %v", err)
	}

	type args struct {
		ctx        context.Context
		youtube_id entities.YoutubeVideoID
	}
	tests := []struct {
		name         string
		y            youtube.YoutubeRepository
		args         args
		wantFile_ids []entities.FileID
		wantErr      bool
	}{
		{"invalid youtube_id", *youtubeRepo, args{ctx: context.Background(), youtube_id: "abcdefgh"}, nil, true},
		{"valid youtube_id", *youtubeRepo, args{ctx: context.Background(), youtube_id: mockYt.YouTube.YoutubeID}, []entities.FileID{file_id}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile_ids, err := tt.y.GetYoutubeFileIDs(tt.args.ctx, tt.args.youtube_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetYoutubeFileIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile_ids, tt.wantFile_ids) {
				t.Errorf("YoutubeRepository.GetYoutubeFileIDs() = %v, want %v", gotFile_ids, tt.wantFile_ids)
			}
		})
	}
}

func TestYoutubeRepository_GetYtdlpVersion(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	err = youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt)

	type args struct {
		ctx        context.Context
		youtube_id entities.YoutubeVideoID
		file_id    entities.FileID
	}
	tests := []struct {
		name        string
		y           youtube.YoutubeRepository
		args        args
		wantVersion *entities.VideoYoutubeDlpVersion
		wantErr     bool
	}{
		{"invalid youtube_id", *youtubeRepo, args{context.Background(), "abcdef", file_id}, nil, true},
		{"valid youtube_id/invalid file_id", *youtubeRepo, args{context.Background(), mockYt.YouTube.YoutubeID, file_id + 999}, nil, true},
		{"valid youtube_id/valid file_id", *youtubeRepo, args{context.Background(), mockYt.YouTube.YoutubeID, file_id}, mock.NewYoutube().DlpVersion, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, err := tt.y.GetYtdlpVersion(tt.args.ctx, tt.args.youtube_id, tt.args.file_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetYtdlpVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVersion, tt.wantVersion) {
				t.Errorf("YoutubeRepository.GetYtdlpVersion() = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestYoutubeRepository_GetFormat(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	err = youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt)

	type args struct {
		ctx        context.Context
		youtube_id entities.YoutubeVideoID
	}
	tests := []struct {
		name       string
		y          youtube.YoutubeRepository
		args       args
		wantFormat *entities.VideoYoutubeFormat
		wantErr    bool
	}{
		{"invalid youtube_id", *youtubeRepo, args{context.Background(), "abcdef"}, nil, true},
		{"valid youtube_id", *youtubeRepo, args{context.Background(), mockYt.YouTube.YoutubeID}, mockYt.Format, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormat, err := tt.y.GetFormat(tt.args.ctx, tt.args.youtube_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFormat, tt.wantFormat) {
				t.Errorf("YoutubeRepository.GetFormat() = %v, want %v", gotFormat, tt.wantFormat)
			}
		})
	}
}

func TestYoutubeRepository_GetTitle(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	err = youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt)

	type args struct {
		ctx        context.Context
		youtube_id entities.YoutubeVideoID
	}
	tests := []struct {
		name      string
		y         youtube.YoutubeRepository
		args      args
		wantTitle string
		wantErr   bool
	}{
		{"title not found", *youtubeRepo, args{context.Background(), "abcdefg"}, "", true},
		{"valid id", *youtubeRepo, args{context.Background(), mockYt.YouTube.YoutubeID}, mockYt.Title, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, err := tt.y.GetTitle(tt.args.ctx, tt.args.youtube_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("YoutubeRepository.GetTitle() = %v, want %v", gotTitle, tt.wantTitle)
			}
		})
	}
}

func TestYoutubeRepository_GetChannelByVideoID(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	err = youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt)

	type args struct {
		ctx        context.Context
		youtube_id entities.YoutubeVideoID
	}
	tests := []struct {
		name        string
		y           youtube.YoutubeRepository
		args        args
		wantChannel entities.VideoYoutubeChannel
		wantErr     bool
	}{
		{"invalid id", *youtubeRepo, args{context.Background(), "abcdefg"}, entities.VideoYoutubeChannel{}, true},
		{"valid id", *youtubeRepo, args{context.Background(), mockYt.YouTube.YoutubeID}, *mockYt.Channel, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChannel, err := tt.y.GetChannelByVideoID(tt.args.ctx, tt.args.youtube_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetChannelByVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChannel, tt.wantChannel) {
				t.Errorf("YoutubeRepository.GetChannelByVideoID() = %v, want %v", gotChannel, tt.wantChannel)
			}
		})
	}
}

func TestYoutubeRepository_GetDescription(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	err = youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt)
	if err != nil {
		t.Fatalf("failed to insert mock youtube entry, %v", err)
	}
	type args struct {
		ctx        context.Context
		youtube_id entities.YoutubeVideoID
	}
	tests := []struct {
		name            string
		y               youtube.YoutubeRepository
		args            args
		wantDescription string
		wantErr         bool
	}{
		{"invalid id", *youtubeRepo, args{context.Background(), "abcdefg"}, "", true},
		{"valid id", *youtubeRepo, args{context.Background(), mockYt.YouTube.YoutubeID}, mockYt.Description, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDescription, err := tt.y.GetDescription(tt.args.ctx, tt.args.youtube_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDescription != tt.wantDescription {
				t.Errorf("YoutubeRepository.GetDescription() = %v, want %v", gotDescription, tt.wantDescription)
			}
		})
	}
}

func TestYoutubeRepository_GetChannelVideos(t *testing.T) {
	db := helper_test.NewDatabase()
	defer db.Close()
	youtubeRepo := youtube.NewYoutubeRepository(db)
	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	file_id := helperInsertFile(*fileRepo, t)
	mockYt := mock.NewYoutube()

	err = youtubeRepo.NewYoutube(context.Background(), file_id, &mockYt)
	if err != nil {
		t.Fatalf("failed to insert mock youtube entry, %v", err)
	}

	type args struct {
		ctx        context.Context
		channel_id entities.YoutubeChannelID
	}
	tests := []struct {
		name       string
		y          youtube.YoutubeRepository
		args       args
		wantVideos []entities.YoutubeVideoID
		wantErr    bool
	}{
		{"no videos", *youtubeRepo, args{ctx: context.Background(), channel_id: entities.UnknownYoutubeChannelID}, nil, true},
		{"one videos", *youtubeRepo, args{ctx: context.Background(), channel_id: mockYt.Channel.ChannelID}, []entities.YoutubeVideoID{mockYt.YouTube.YoutubeID}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVideos, err := tt.y.GetChannelVideos(tt.args.ctx, tt.args.channel_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("YoutubeRepository.GetChannelVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVideos, tt.wantVideos) {
				t.Errorf("YoutubeRepository.GetChannelVideos() = %v, want %v", gotVideos, tt.wantVideos)
			}
		})
	}
}
