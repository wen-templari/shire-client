package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/templari/shire-client/util"

	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func InitDB(id int) {
	var err error
	db, err = sqlx.Connect("sqlite3", fmt.Sprint("file:shire", id, ".db?cache=shared&mode=rwc"))
	if err != nil {
		log.Error(err)
		// util.Logger.Fatal(err)
	}
	// db.SetMaxOpenConns(1000) // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(0) // 0, connections are reused forever.
	if err != nil {
		util.Logger.Fatal(err)
	}
}

func SetDB(dbx *sqlx.DB) {
	db = dbx
}

func CloseDB() {
	db.Close()
}
