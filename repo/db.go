package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/templari/shire-client/model"
	"github.com/templari/shire-client/util"

	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func InitDB(id int) {
	var err error
	db, err = sqlx.Connect("sqlite3", fmt.Sprint("shire", id, ".db?cache=shared&mode=rwc"))
	if err != nil {
		log.Println(err)
		util.Logger.Fatal(err)
	}
	_, err = createTable()
	log.Println(err)
	// db.SetMaxOpenConns(1000) // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(0) // 0, connections are reused forever.
	// if err != nil {
	// 	util.Logger.Fatal(err)
	// }
}

func createTable() (sql.Result, error) {
	schema := `CREATE TABLE group_user  (
  	groupId int NOT NULL,
  	userId int NOT NULL,
  	PRIMARY KEY (groupId, userId)
	);

	CREATE TABLE message  (
  	"from" int NOT NULL DEFAULT 0,
  	"to" int NOT NULL DEFAULT 0,
  	content varchar(255) NOT NULL DEFAULT "",
  	groupId int NOT NULL DEFAULT 0,
  	time datetime NULL
	);

	CREATE TABLE user  (
  	id int NOT NULL ,
  	name varchar(255) NOT NULL DEFAULT "",
  	address varchar(255) NOT NULL DEFAULT "",
  	port int NOT NULL DEFAULT 0,
  	rpcPort int NOT NULL DEFAULT 0,
  	createdAt datetime NOT NULL,
  	updatedAt datetime NOT NULL,
  	PRIMARY KEY (id)
	);`
	// execute a query on the server
	return db.Exec(schema)
}

func SaveMessage(message *model.Message) error {
	sql, args, _ := squirrel.Insert("message").Columns("'from'", "'to'", "content", "groupId", "time").Values(
		message.From, message.To, message.Content, message.GroupId, message.Time).ToSql()
	_, err := db.Exec(sql, args...)
	return err
}

func GetAllMessage() ([]model.Message, error) {
	messageList := []model.Message{}
	sql, _, _ := squirrel.Select("*").From("message").ToSql()
	if err := db.Select(&messageList, sql); err != nil {
		return []model.Message{}, err
	}
	return messageList, nil
}

func SetDB(dbx *sqlx.DB) {
	db = dbx
}

func CloseDB() {
	db.Close()
}
