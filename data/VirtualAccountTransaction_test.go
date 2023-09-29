package data

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/helpers"
	"context"
	"log"
	"testing"
)

func TestVATrasaction(t *testing.T) {

	ctx := context.TODO()
	conn, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, conn)

	/*
		HistoryID primitive.ObjectID `json:"history_id" bson:"history_id"`
			NoVirtualAccount string `json:"no_virtual_account" bson:"no_virtual_account"`
			NameVirtualAccount string `json:"name_virtual_account" bson:"name_virtual_account"`
			TransactionAmount float32 `json:"transaction_amount" bson:"transaction_amount"`
			reference string `json:"reference" bson:"reference"`
	*/
	NoVirtualAccount := "asdad000001"
	NameVirtualAccount := "Imam"
	TransactionAmount := float32(0)
	TransactionType := "PAYMENT"
	reference := ""
	TransactionID := ""

	ModelVAtrasaction := CreateVaTransaction(conn)
	err = ModelVAtrasaction.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	CekVATrasaction, err := ModelVAtrasaction.AddVirtualAccountTransaction(ctx,
		NoVirtualAccount,
		NameVirtualAccount,
		TransactionAmount,
		TransactionType,
		reference,
		TransactionID)
	if err != nil {
		t.Fatal(err)
	}

	VAtrasactions, err := ModelVAtrasaction.ListVirtualAccountTransaction(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("List Virtual Trasaction", helpers.ToJson(VAtrasactions))

	err, _ = ModelVAtrasaction.FindOene(ctx, CekVATrasaction.HistoryID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	log.Println("FoundVirtualAccount Transaction :", helpers.ToJson(CekVATrasaction))

}
