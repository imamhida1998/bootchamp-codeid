package service

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"bootchamp-codeid/helpers"
	"context"
	"log"
	"testing"
)

func TestVirtualAccount(t *testing.T) {
	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	// ModelTestVirtualAccount := data.CreateLoginNoSql(mongoClient)
	// err = ModelTestVirtualAccount.Truncate(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//log.Println(helpers.ToJson(ModelTestVirtualAccount))

	NoHp := "081317767015"
	Password := "123123"
	VirtualService := CreateVirtualAccount(mongoClient, db)
	token, err := VirtualService.Login(ctx, NoHp, Password)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Token: ", token)

	login, err := VirtualService.ParseToken(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("login=", helpers.ToJson(login))
}

//
//func TestVATransaction_CreateTransfer(t *testing.T) {
//	ctx := context.TODO()
//	mongoClient, err := mangodb.OpenMongoDb(ctx)
//	db := postgresql.OpenConnection(config.IMAM)
//	defer db.Close()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mangodb.CloseMongoDb(ctx, mongoClient)
//
//	/*
//	HistoryID          primitive.ObjectID `json:"history_id" bson:"history_id"`
//		NoVirtualAccount   string             `json:"no_virtual_account" bson:"no_virtual_account"`
//		NameVirtualAccount string             `json:"name_virtual_account" bson:"name_virtual_account"`
//		TransactionAmount  float32            `json:"transaction_amount" bson:"transaction_amount"`
//		TransactionType string	`json:"transaction_type" bson:"transaction_type"`
//		Description          string             `json:"description" bson:"description"`
//		TransactionID string	`json:"transaction_id" bson:"transaction_id"`
//	 */
//	NoVirtualAccount := ""
//	NameVirtualAccount := ""
//	TransactionAmount := float32(3000)
//	TransactionType := "TRANSFER"
//	ModelTransaction := data.CreateVirtualAccount(db)
//	VirtualService := CreateVirtualAccount(mongoClient, db)
//	token, err := VirtualService.CreateTransfer(ctx,
//		NoVirtualAccount,
//		NameVirtualAccount,
//		TransactionAmount,
//		TransactionType)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
