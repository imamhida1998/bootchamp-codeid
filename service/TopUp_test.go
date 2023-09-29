package service

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"context"
	"log"
	"testing"
)

func TestTopUp(t *testing.T) {
	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	bankcode := "BCA01"
	accountno := "081317767015000001"
	topupamount := float32(5000)
	pin := "123456"
	VirtualTopUp := CreateTopUpService(mongoClient, db)
	virtualtransaction, BankTransaction, TopUp, err := VirtualTopUp.TopUp(ctx, bankcode, accountno, topupamount, pin)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Bank Transaction: ", BankTransaction)
	log.Println("Virtual Transaction: ", virtualtransaction)
	log.Println("Topup: ", TopUp)

}
