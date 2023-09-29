package service

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"context"
	"log"
	"testing"
)

func TestTransfers(t *testing.T) {
	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	NoVirtualAccount := "03131314000001"      // pemberi // imam
	payVirtualAccount := "081317767015000001" // penerima
	PaymentAmount := float32(14000)
	pin := "123456"
	VirtualService := CreateTransferService(mongoClient, db)
	Payment, UserTransaction, PayTransaction, err := VirtualService.Transfer(ctx, payVirtualAccount, NoVirtualAccount, PaymentAmount, pin)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Penerima :", Payment)
	log.Println("Pemberi : ", UserTransaction)
	log.Println("PaymentTransaction : ", PayTransaction)

}
