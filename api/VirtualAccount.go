package api

import (
	"bootchamp-codeid/data"
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type VirtualAccountApi struct {
	BaseApi
	data.VirtualAccount
	VirtualAccountNo string `json:"virtual_account_no" bson:"virtual_account_no"`
}

func (virtualaccountApi VirtualAccountApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, *sql.DB, context.Context) {
	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()

	_, err = virtualaccountApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		virtualaccountApi.Error(w, err)
		return nil, nil, nil
	}
	return mongoClient, db, ctx

}

func (virtualaccountApi VirtualAccountApi) GetVirtualAccount(w http.ResponseWriter, r *http.Request) {
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()

	branchModel := data.CreateVirtualAccount(db)
	branchs, err := branchModel.ReadVirtualAccount()
	if err != nil {
		virtualaccountApi.Error(w, err)
		return
	} else {
		virtualaccountApi.Json(w, branchs, http.StatusOK)
		return
	}
}

/*
Account_no  string  `json:"NoAccount" bson:"account___no"`
	No_hp       string  `json:"NoHandphone" bson:"no___hp"`
	Email       string  `json:"Email" bson:"email"`
	seq_no      int     `json:"seq_no" bson:"seq_no"`
	AccountName string  `json:"AccountName" bson:"account_name"`
	Pin         string  `json:"PIN" bson:"pin" bson:"pin"`
	Password    string  `json:"password" bson:"password"`
	Saldo       float32 `json:"Saldo" bson:"saldo"`
	CreatedAt   string  `json:"CreatedAt" bson:"created_at"`
	//	CreatedAt   string  `json:"CreatedAt" bson:"created_at"`
	UpdatedAt *string `json:"UpdatedAt" bson:"updated_at"`
*/

func (virtualaccountApi VirtualAccountApi) PostVirtualAccount(w http.ResponseWriter, r *http.Request) {
	_, db, _ := virtualaccountApi.Connection(w, r)
	defer db.Close()

	PostModel := data.CreateVirtualAccount(db)
	Post, err := PostModel.InsertVirtualAccount(
		virtualaccountApi.Account_no,
		virtualaccountApi.No_hp,
		virtualaccountApi.AccountName,
		virtualaccountApi.Email,
		virtualaccountApi.Pin,
		virtualaccountApi.Password,
		virtualaccountApi.Saldo)
	if err != nil {
		return
	} else {
		virtualaccountApi.Json(w, Post, http.StatusOK)
	}

}

func (virtualaccountApi VirtualAccountApi) UpdateVirtual(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&virtualaccountApi)
	if err != nil {
		virtualaccountApi.Error(w, err)
		return
	}
	_, db, _ := virtualaccountApi.Connection(w, r)
	defer db.Close()

	ModelUpdate := data.CreateVirtualAccount(db)
	err = virtualaccountApi.UpdateVirtualAccount(
		virtualaccountApi.Account_no,
		virtualaccountApi.No_hp,
		virtualaccountApi.AccountName,
		virtualaccountApi.Email,
		virtualaccountApi.Pin,
		virtualaccountApi.Password,
		virtualaccountApi.Saldo)
	if err != nil {
		virtualaccountApi.Error(w, err)
		return
	} else {
		virtualaccountApi.Json(w, ModelUpdate, http.StatusOK)
	}

}

func (virtualaccountApi VirtualAccountApi) DeleteVirtualAccount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&virtualaccountApi)
	if err != nil {
		virtualaccountApi.Error(w, err)
		return
	}
	_, db, _ := virtualaccountApi.Connection(w, r)
	defer db.Close()

	ModelDelete := data.CreateVirtualAccount(db)
	Deletes, err := ModelDelete.Remove(virtualaccountApi.Account_no)
	if err != nil {
		virtualaccountApi.Error(w, err)
		return
	} else {
		virtualaccountApi.Json(w, Deletes, http.StatusOK)
	}
}
