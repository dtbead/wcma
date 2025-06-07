package testing

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	"github.com/dtbead/wc-maps-archive/internal/entities"
	mock_file "github.com/dtbead/wc-maps-archive/internal/helper/testing/mock/file"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConnection struct {
	Ip, User, Password, DbName, Port string
}

var DefaultConnection = PostgresConnection{
	Ip:       "localhost",
	User:     "postgres",
	Password: "password",
	DbName:   "wc_staging",
	Port:     "6664",
}

// NewDsn returns a new database connection url string
func NewDsn(pgConnOpt PostgresConnection) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		pgConnOpt.User, pgConnOpt.Password, pgConnOpt.Ip, pgConnOpt.Port, pgConnOpt.DbName)
}

// isInitialized determines whether the database driver "txdb1" has already been registered.
var isInitialized bool

// NewDatabase returns a new database which doesn't commit any transactions. If no pgConnOpt is given, then
// the DefaultConnection variable will be used instead.
func NewDatabase(pgConnOpt *PostgresConnection) *sql.DB {
	if pgConnOpt == nil {
		pgConnOpt = &DefaultConnection
	}

	dsn := NewDsn(*pgConnOpt)

	// creates custom database driver which doesn't commit any transactions to postgres/pgx
	if !isInitialized {
		txdb.Register("txdb1", "pgx", dsn)
		isInitialized = true
	}

	db, err := sql.Open("txdb1", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to create database for testing, %v", err))
	}

	// initialize database schema
	_, err = db.Exec(postgres.Schema)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize testing database schema, %v", err))
	}

	return db
}

// InsertFakFile creates a new and random database entry into the 'file' table. amount specifies how much
// fake files are to be generated.
// The given sql.DB is expected to have a matching table schema in `...\internal\storage\postgres\schema.sql`.
func InsertFakeFile(db *sql.DB, amount int) (file_ids []entities.FileID, err error) {
	/*
		CREATE TABLE "file" (
			"id" BIGINT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
			"path" TEXT NOT NULL UNIQUE,
			"extension" TEXT NOT NULL CHECK (length(extension) >= 3 AND length(extension) <= 6),
			"md5" BYTEA NOT NULL UNIQUE CHECK (length(md5) = 16),
			"sha1" BYTEA NOT NULL UNIQUE CHECK (length(sha1) = 20),
			"sha256" BYTEA NOT NULL UNIQUE CHECK (length(sha256) = 32),
			"filesize" BIGINT NOT NULL CHECK (filesize >= 16),
			PRIMARY KEY("id")
		);

		CREATE TABLE "file_video" (
			"file_id" BIGINT NOT NULL UNIQUE,
			"duration" INTEGER NOT NULL CHECK (duration >= 0),
			"width" SMALLINT NOT NULL CHECK (width >= 0),
			"height" SMALLINT NOT NULL CHECK (height >= 0),
			"fps" SMALLINT CHECK (fps >= 0),
			"video_codec" TEXT,
			"audio_codec" TEXT,
			FOREIGN KEY ("file_id") REFERENCES file("id")
			ON UPDATE CASCADE ON DELETE CASCADE
		);
	*/
	stmt, err := db.Prepare("INSERT INTO file (path, extension, md5, sha1, sha256, filesize) VALUES ($1, $2, $3, $4, $5) RETURNING id;")
	if err != nil {
		return nil, err
	}

	file_ids = []entities.FileID{}
	var file_id int64

	for range amount {
		mockFile := mock_file.RandomFile()
		res, err := stmt.Query(mockFile.PathRelative, mockFile.Extension, mockFile.Hashes.MD5, mockFile.Hashes.SHA1, mockFile.Hashes.SHA256, mockFile.Size)
		if err != nil {
			return nil, err
		}

		for res.Next() {
			err = res.Scan(&file_id)
			if err != nil {
				return nil, err
			}
		}

		file_ids = append(file_ids, entities.FileID(file_id))
	}

	return file_ids, nil
}
