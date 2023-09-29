package data

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/helpers"
	"context"
	"log"
	"testing"
)

func TestPay_AddPayment(t *testing.T) {
	ctx := context.TODO()

	conn, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, conn)
	/*
	 MerchantNoVirtualAccount   string             `json:"merchant_no_virtual_account" bson:"merchant_no_virtual_account"`
	 	MerchantNameVirtualAccount string             `json:"merchant_name_virtual_account" bson:"merchant_name_virtual_account"`
	 	SourceNoVirtualAccount     string             `json:"source_no_virtual_account" bson:"source_no_virtual_account"`
	 	SourceNameVirtualAccount   string             `json:"source_name_virtual_account" bson:"source_name_virtual_account"`
	 	PayAmount                  float32            `json:"pay_amount" bson:"pay_amount"`
	 	CreatedAt                  string             `json:"created_at" bson:"created_at"`

	*/
	MerchantNoVirtualAccount := "NO1"
	MerchantNameVirtualAccount := "GATAU"
	SourceNoVirtualAccount := "ACCOUNT-01"
	SourceNameVirtualAccount := "Imam"
	PayAmount := float32(1000000)

	ModelPay := CreatePay(conn)
	err = ModelPay.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	AddPayment, err := ModelPay.AddPayment(ctx,
		MerchantNoVirtualAccount,
		MerchantNameVirtualAccount,
		SourceNoVirtualAccount,
		SourceNameVirtualAccount,
		PayAmount)
	if err != nil {
		t.Fatal(err)
	}
	listPayment, err := ModelPay.ListPayment(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("ListPayment=", helpers.ToJson(listPayment))

	err = ModelPay.FindOne(ctx, AddPayment.PayID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	log.Println(err)
}
