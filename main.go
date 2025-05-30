package main

import (
	"context"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/dtbead/wc-maps-archive/internal/download/ytdlp"
	"github.com/dtbead/wc-maps-archive/internal/service"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var wc_main_pg string = `postgres://postgres:password@localhost:6663/wc_main`
var wc_test_pg string = `postgres://postgres:password@localhost:6663/postgres`

func main() {
	db, err := sql.Open("pgx", wc_main_pg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	storage, err := postgres.NewRepository(db, "c:/users/holly/desktop/wc-test")
	if err != nil {
		log.Fatal(err)
	}
	service := service.NewService(storage)
	ytdownloader := ytdlp.NewYtdlp(nil)

	err = service.YoutubeService.DownloadVideo(context.Background(), "https://www.youtube.com/watch?v=_RMk8qfoAAo", ytdownloader)
	if err != nil {
		log.Fatal(err)
	}

}
