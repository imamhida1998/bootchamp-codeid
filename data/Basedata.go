package data

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

type BaseNoSql struct {
	Client *mongo.Client `json:"-" bson:"-"`
}

type BaseData struct {
	DB             *sql.DB `json:",omitempty"`
	Mangodb		   *mongo.Client `json:",omitempty"`
	Transaction    *sql.Tx `json:",omitempty"`
	UseTransaction bool    `json:",omitempty"`
}

func (baseData BaseData) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if baseData.UseTransaction {
		return baseData.Transaction.Exec(sql, args...)
	} else {
		return baseData.DB.Exec(sql, args...)
	}
}

func (baseData BaseData) QueryRow(sql string, args ...interface{}) *sql.Row {
	if baseData.UseTransaction {
		return baseData.Transaction.QueryRow(sql, args...)
	} else {
		return baseData.DB.QueryRow(sql, args...)
	}
}

func (baseData BaseData) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if baseData.UseTransaction {
		return baseData.Transaction.Query(sql, args...)
	} else {
		return baseData.DB.Query(sql, args...)
	}
}
