package api

import (
	"bootchamp-codeid/data"
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"context"
	"database/sql"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type BankAccountApi struct {
	BaseApi
	data.BankAccount
	BankNo string `json:"bank_no"bson:"bank_no"`
}

func (BankApi BankAccountApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, *sql.DB, context.Context) {
	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()

	_, err = BankApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		BankApi.Error(w, err)
		return nil, nil, nil
	}
	return mongoClient, db, ctx

}

func (BankApi BankAccountApi) GetBankAccount(w http.ResponseWriter, r *http.Request) {
	_, db, _ := BankApi.Connection(w, r)
	defer db.Close()

	ModelBank := data.CreateBankAccount(db)
	devices, err := ModelBank.ReadBankAccount()
	if err != nil {
		BankApi.Error(w, err)
		return
	} else {
		BankApi.Json(w, devices, http.StatusOK)
		return
	}
}

// func (BankApi BankAccountApi) TopupBank(w http.ResponseWriter, r *http.Request) {

// 	_,db, _ := BankApi.Connection(w,r)
// 	defer db.Close()

// 	ModelBank := data.CreateBankAccount(db)
// 	devices ,err := ModelBank.UpdateBankAccount(BankApi.BankAccount_No,BankApi.BankAccount_Owner,BankApi.Saldo)
// 	 if err != nil {
// 		 BankApi.Error(w, err)
// 		 return
// 	 } else {
// 		 BankApi.Json(w, devices,http.StatusOK)
// 	 }

// }
