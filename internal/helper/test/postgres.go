package test

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DATA-DOG/go-txdb"
	"github.com/dtbead/wc-maps-archive/internal/storage/postgres"
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

func init() {
	// creates custom database driver which doesn't commit any transactions to postgres/pgx
	txdb.Register("txdb1", "pgx", dsn)
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
