package repo

import (
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/templari/shire-client/core/model"
	"github.com/templari/shire-client/util"
)

var db *sqlx.DB

func InitDB() {
	var err error
	db, err = sqlx.Connect("mysql", os.Getenv("DB_URL"))
	// db.SetMaxOpenConns(1000) // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(0) // 0, connections are reused forever.
	if err != nil {
		util.Logger.Fatal(err)
	}
}

func CreateMessage(message model.Message) error {
	sql, args, _ := squirrel.Insert("message").
		Columns("from", "to", "time", "content", "groupId").
		Values(message.From, message.To, message.Time, message.Content, message.GroupId).ToSql()
	_, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

type GetMessagesArgs struct {
	From    int
	To      int
	GroupId int
}

func GetMessages(arg GetMessagesArgs) ([]model.Message, error) {
	messages := []model.Message{}
	query := squirrel.Select("*").From("message")
	if arg.GroupId != 0 {
		query = query.Where(squirrel.Eq{"groupId": arg.GroupId})
	} else if arg.From != 0 && arg.To != 0 {
		query = query.
			Where(squirrel.Eq{"from": arg.From}).
			Where(squirrel.Eq{"to": arg.To})
	}
	sql, args, _ := query.ToSql()
	if err := db.Select(&messages, sql, args...); err != nil {
		return nil, err
	}
	return messages, nil
}

func GetUser() ([]model.Message, error) {
	messages := []model.Message{}
	if err := db.Select(&messages, "SELECT * FROM message"); err != nil {
		return nil, err
	}
	return messages, nil
}
