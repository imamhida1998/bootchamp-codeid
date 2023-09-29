package data

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/helpers"
	"context"
	"fmt"
	"log"
	"testing"
)

func TestTransfer_AddTransfer(t *testing.T) {
	ctx := context.TODO()
	conn, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, conn)

	ModelTransfer := CreateTransfer(conn)
	err = ModelTransfer.Truncate(ctx)
	if err != nil {
		t.Fatal(err)

	}
	/*
		TransferID primitive.ObjectID`json:"transfer_id" bson:"transfer_id"`
			SourceNoVirtualAccount string `json:"source_no_virtual_account" bson:"source_no_virtual_account"`
			SourceNameVirtualAccount string `json:"source_name_virtual_account" bson:"source_name_virtual_account"`
			DesNoVirtualAccount string `json:"des_no_virtual_account" bson:"des_no_virtual_account"`
			DescNameVirtualAccount string `json:"desc_name_virtual_account" bson:"desc_name_virtual_account"`
			TransferAmount float32 `json:"transfer_amount" bson:"transfer_amount"`
			CreateAt string `json:"create_at" json:"create_at"`
	*/
	SourceNoVirtualAccount := "081317767016000001"
	SourceNameVirtualAccount := "Yudis"
	DesNoVirtualAccount := "081317767015000001"
	DescNameVirtualAccount := "Imam"
	TransferAmount := float32(1500)

	CekTransfer, err := ModelTransfer.AddTransfer(ctx,
		SourceNoVirtualAccount,
		SourceNameVirtualAccount,
		DesNoVirtualAccount,
		DescNameVirtualAccount,
		TransferAmount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
	ListTransfer, err := ModelTransfer.ListTransfer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("ListTransfer :", helpers.ToJson(ListTransfer))

	//err = ModelTransfer.FindOne(ctx , CekTransfer.TransferID.Hex())
	//if err != nil {
	//	t.Fatal(err)
	//	t.Fatal(err)
	//}
	log.Println("CekTransfer :", helpers.ToJson(CekTransfer))

}
