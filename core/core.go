package core

import (
	"github.com/jmoiron/sqlx"
	"github.com/templari/shire-client/core/model"
)

type Core struct {
	InfoServerAddress string
	db                *sqlx.DB
	user              model.User
	token             string
}

func MakeCore(address string, db *sqlx.DB) *Core {
	return &Core{
		InfoServerAddress: address,
		db:                db,
	}
}
