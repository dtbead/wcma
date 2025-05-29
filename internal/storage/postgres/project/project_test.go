package project_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/helper"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/file"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres/project"
	"github.com/google/go-cmp/cmp"
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
	txdb.Register("txdb1", "pgx", dsn)
	os.Exit(m.Run())
}

func NewDatabase() *sql.DB {
	db, err := sql.Open("txdb1", dsn)
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

func TestProjectRepository_NewProject(t *testing.T) {
	db := NewDatabase()
	defer db.Close()
	projectRepo := project.NewProjectRepository(db)

	type args struct {
		ctx     context.Context
		project *entities.Project
	}
	tests := []struct {
		name     string
		p        project.ProjectRepository
		args     args
		wantUuid entities.ProjectUUID
		wantErr  bool
	}{
		{"valid project insert", *projectRepo, args{context.Background(), &entities.Project{
			UUID:          "r4ruYKZQT3XBtpxPjr6s9k",
			FileIDs:       nil,
			ProjectType:   entities.ProjectAnimatedMusicVideo,
			DateAnnounced: time.Now(),
			DateCompleted: time.Now(),
			DateArchived:  time.Time{},
		}}, "r4ruYKZQT3XBtpxPjr6s9k", false},
		{"nil project", *projectRepo, args{context.Background(), &entities.Project{}}, entities.InvalidProjectUUID, true},
		{"empty project uuid", *projectRepo, args{context.Background(), &entities.Project{}}, entities.InvalidProjectUUID, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUuid, err := tt.p.NewProject(tt.args.ctx, tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectRepository.NewProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUuid, tt.wantUuid) {
				t.Errorf("ProjectRepository.NewProject() = %v, want %v", gotUuid, tt.wantUuid)
			}
		})
	}
}

func TestProjectRepository_DeleteProject(t *testing.T) {
	db := NewDatabase()
	defer db.Close()
	projectRepo := project.NewProjectRepository(db)

	p := entities.Project{
		UUID:          helper.RandomUUID(),
		ProjectType:   entities.ProjectAnimatedMusicVideo,
		DateAnnounced: time.Now(),
		DateCompleted: time.Now(),
	}

	uuid, err := projectRepo.NewProject(context.Background(), &p)
	if err != nil {
		t.Fatalf("failed to create mock project, %v", err.Error())
	}

	type args struct {
		ctx  context.Context
		uuid entities.ProjectUUID
	}
	tests := []struct {
		name    string
		p       project.ProjectRepository
		args    args
		wantErr bool
	}{
		{"empty uuid", *projectRepo, args{context.Background(), ""}, false},
		{"valid delete", *projectRepo, args{context.Background(), uuid}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.DeleteProject(tt.args.ctx, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ProjectRepository.DeleteProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProjectRepository_GetProject(t *testing.T) {
	db := NewDatabase()
	defer db.Close()
	projectRepo := project.NewProjectRepository(db)

	p := entities.Project{
		UUID:          helper.RandomUUID(),
		FileIDs:       []entities.FileID{},
		ProjectType:   entities.ProjectAnimatedMusicVideo,
		DateAnnounced: time.Now().UTC().Round(time.Second),
		DateCompleted: time.Now().UTC().Round(time.Second),
	}

	uuid, err := projectRepo.NewProject(context.Background(), &p)
	if err != nil {
		t.Fatalf("failed to create mock project, %v", err.Error())
	}

	type args struct {
		ctx  context.Context
		uuid entities.ProjectUUID
	}
	tests := []struct {
		name        string
		p           project.ProjectRepository
		args        args
		wantProject *entities.Project
		wantErr     bool
	}{
		{"empty project", *projectRepo, args{context.Background(), ""}, nil, true},
		{"valid project", *projectRepo, args{context.Background(), uuid}, &p, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProject, err := tt.p.GetProject(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectRepository.GetProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotProject != nil {
				// DateArchived is ignored when inserting to database. Just set it the same so cmp.Diff won't
				// complain.
				tt.wantProject.DateArchived = gotProject.DateArchived

				if !cmp.Equal(gotProject, tt.wantProject, cmp.AllowUnexported(entities.Project{})) {
					t.Errorf("got diff, %v", cmp.Diff(*gotProject, *tt.wantProject, cmp.AllowUnexported(entities.Project{})))
				}
			}
		})
	}
}

func TestProjectRepository_AssignProjectFile(t *testing.T) {
	db := NewDatabase()
	defer db.Close()
	projectRepo := project.NewProjectRepository(db)

	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repository, %v", err.Error())
	}

	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err.Error())
	}

	file_id, err := fileRepo.NewFile(context.Background(), f, "mkv")
	if err != nil {
		t.Fatalf("failed to store test file in file repo, %v", err.Error())
	}

	p := entities.Project{
		UUID:          helper.RandomUUID(),
		ProjectType:   entities.ProjectAnimatedMusicVideo,
		DateAnnounced: time.Now(),
		DateCompleted: time.Now(),
	}

	uuid, err := projectRepo.NewProject(context.Background(), &p)
	if err != nil {
		t.Fatalf("failed to create mock project, %v", err.Error())
	}

	type args struct {
		ctx     context.Context
		uuid    entities.ProjectUUID
		file_id entities.FileID
	}
	tests := []struct {
		name    string
		p       project.ProjectRepository
		args    args
		wantErr bool
	}{
		{"valid assign", *projectRepo, args{context.Background(), uuid, file_id}, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.AssignProjectFile(tt.args.ctx, tt.args.uuid, tt.args.file_id); (err != nil) != tt.wantErr {
				t.Errorf("ProjectRepository.AssignProjectFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotFileIds, err := tt.p.GetProjectVideos(tt.args.ctx, tt.args.uuid)
			if err != nil {
				t.Errorf("failed to get project videos, %v", err)
			}

			if !slices.Contains(gotFileIds, tt.args.file_id) {
				t.Errorf("got fileIds %v, want %v", gotFileIds, tt.args.file_id)
			}
		})
	}
}

func TestProjectRepository_AssignProjectFileEmptyUUID(t *testing.T) {
	db := NewDatabase()
	defer db.Close()
	projectRepo := project.NewProjectRepository(db)

	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repository, %v", err.Error())
	}

	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err.Error())
	}

	file_id, err := fileRepo.NewFile(context.Background(), f, "mkv")
	if err != nil {
		t.Fatalf("failed to store test file in file repo, %v", err.Error())
	}

	_, err = projectRepo.NewProject(context.Background(), &entities.Project{
		UUID:          helper.RandomUUID(),
		ProjectType:   entities.ProjectAnimatedMusicVideo,
		DateAnnounced: time.Now(),
		DateCompleted: time.Now(),
	})

	if err != nil {
		t.Fatalf("failed to create mock project, %v", err.Error())
	}

	err = projectRepo.AssignProjectFile(context.Background(), "", file_id)
	if err == nil {
		t.Errorf("ProjectRepository.AssignProjectFile() error = %v, wantErr %v", err, true)
		return
	}

}

func TestProjectRepository_UnassignProjectVideoInvalidFileID(t *testing.T) {
	db := NewDatabase()
	defer db.Close()
	projectRepo := project.NewProjectRepository(db)

	fileRepo, err := file.NewFileRepository(db, t.TempDir())
	if err != nil {
		t.Fatalf("failed to create file repository, %v", err.Error())
	}

	f, err := os.Open("testdata/y_wo8pyoxyk.mkv")
	if err != nil {
		t.Fatalf("failed to open test file, %v", err.Error())
	}

	file_id, err := fileRepo.NewFile(context.Background(), f, "mkv")
	if err != nil {
		t.Fatalf("failed to store test file in file repo, %v", err.Error())
	}

	uuid, err := projectRepo.NewProject(context.Background(), &entities.Project{
		UUID:          helper.RandomUUID(),
		ProjectType:   entities.ProjectAnimatedMusicVideo,
		DateAnnounced: time.Now(),
		DateCompleted: time.Now(),
	})

	if err != nil {
		t.Fatalf("failed to create mock project, %v", err.Error())
	}

	err = projectRepo.AssignProjectFile(context.Background(), uuid, file_id+999)
	if err == nil {
		t.Errorf("ProjectRepository.AssignProjectFile() error = %v, wantErr %v", err, true)
		return
	}
}
