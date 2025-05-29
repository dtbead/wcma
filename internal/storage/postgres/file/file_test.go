package file_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/dtbead/wc-maps-archive/internal/entities"
	file_helper "github.com/dtbead/wc-maps-archive/internal/helper/file"
	"github.com/google/go-cmp/cmp"

	"github.com/dtbead/wc-maps-archive/internal/storage/postgres"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConnection struct {
	Ip, User, Password, DbName, Port string
}

var pgConnOpt = PostgresConnection{
	Ip:       "localhost",
	User:     "postgres",
	Password: "password",
	DbName:   "wc_staging",
	Port:     "6664",
}

// database connection url
var dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
	pgConnOpt.User, pgConnOpt.Password, pgConnOpt.Ip, pgConnOpt.Port, pgConnOpt.DbName)

func TestMain(m *testing.M) {
	// creates custom database driver which doesn't commit any transactions to postgres/pgx
	txdb.Register("txdb", "pgx", dsn)
	os.Exit(m.Run())
}

func NewDatabase() *sql.DB {
	db, err := sql.Open("txdb", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// initialize database schema
	_, err = db.Exec(postgres.Schema)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestFileRepository_NewFile(t *testing.T) {
	db := NewDatabase()
	defer db.Close()

	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}

	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err)
	}

	type args struct {
		ctx       context.Context
		file      io.Reader
		extension string
	}
	tests := []struct {
		name        string
		f           file.FileRepository
		args        args
		wantFile_id entities.FileID
		wantErr     bool
	}{
		{"nil io.Reader", *fileRepo, args{ctx: context.Background(), file: nil, extension: "mkv"}, entities.InvalidFileID, true},
		{"valid file", *fileRepo, args{ctx: context.Background(), file: f, extension: "mkv"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile_id, err := tt.f.NewFile(tt.args.ctx, tt.args.file, tt.args.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepository.NewFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile_id, tt.wantFile_id) {
				t.Errorf("FileRepository.NewFile() = %v, want %v", gotFile_id, tt.wantFile_id)
			}
		})
	}
}

func TestFileRepository_DeleteFile(t *testing.T) {
	db := NewDatabase()
	defer db.Close()

	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}
	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err)
	}

	file_id, err := fileRepo.NewFile(context.Background(), f, "mkv")
	if err != nil {
		t.Fatalf("failed to insert test file, %v", err)
	}

	type args struct {
		ctx     context.Context
		file_id entities.FileID
	}
	tests := []struct {
		name    string
		f       file.FileRepository
		args    args
		wantErr bool
	}{
		{"invalid/missing file_id", *fileRepo, args{context.Background(), file_id + 1}, true},
		{"valid delete", *fileRepo, args{context.Background(), file_id}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.DeleteFile(tt.args.ctx, tt.args.file_id); (err != nil) != tt.wantErr {
				t.Errorf("FileRepository.DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileRepository_GetFile(t *testing.T) {
	db := NewDatabase()
	defer db.Close()

	tempDir := t.TempDir()
	fileRepo, err := file.NewFileRepository(db, tempDir)
	if err != nil {
		t.Fatalf("failed to create file repo, %v", err)
	}

	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err)
	}

	file_id, err := fileRepo.NewFile(context.Background(), f, "mkv")
	if err != nil {
		t.Fatalf("failed to insert test file, %v", err)
	}

	type args struct {
		ctx     context.Context
		file_id entities.FileID
	}
	tests := []struct {
		name              string
		f                 file.FileRepository
		args              args
		wantFile_metadata *entities.File
		wantErr           bool
	}{
		{"invalid/missing file_id", *fileRepo, args{context.Background(), file_id + 1}, nil, true},
		{"valid file_id", *fileRepo, args{context.Background(), file_id}, &entities.File{
			PathRelative: "6b/6bebd6bfc85e9840e6bb47e1f329b5453afd184fa7fe2d52da3cd46200062ddc.mkv",
			PathAbsolute: file_helper.SanitizePath(tempDir + "/" + "6b/6bebd6bfc85e9840e6bb47e1f329b5453afd184fa7fe2d52da3cd46200062ddc.mkv"),
			Extension:    "mkv",
			Size:         330070,
			Hashes: entities.Hashes{
				SHA256: file_helper.HexStringToByte("6bebd6bfc85e9840e6bb47e1f329b5453afd184fa7fe2d52da3cd46200062ddc"),
				SHA1:   file_helper.HexStringToByte("a51852be25a413051cfc8954380f64a3f0668478"),
				MD5:    file_helper.HexStringToByte("3e9cb26fede9bd75e406cbb7e6ff81e6"),
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile_metadata, err := tt.f.GetFile(tt.args.ctx, tt.args.file_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepository.GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotFile_metadata, tt.wantFile_metadata) {
				t.Errorf("got diff %s", cmp.Diff(gotFile_metadata, tt.wantFile_metadata))
			}
		})
	}
}
