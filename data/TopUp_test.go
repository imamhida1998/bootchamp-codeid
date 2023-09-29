package data

import (
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/helpers"
	"context"
	"log"
	"testing"
)

func Test_topup(t *testing.T) {

	ctx := context.TODO()

	conn, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, conn)
	//
	bankcode := "BCA-01"
	bankname := "Bank Central Anjay"
	NoVirtualAccount := "081317767015000001"
	NameVirtualAccount := "Imam"
	TopUpAmount := float32(50000)
	topUpNoSql := CreateTopUp(conn)
	// err = topUpNoSql.Truncate(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	user, err := topUpNoSql.AddTopUp(ctx, bankcode, bankname, NoVirtualAccount, NameVirtualAccount, TopUpAmount)

	if err != nil {
		t.Fatal(err)
	}

	users, err := topUpNoSql.ListTopUps(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("listUserGroups=", helpers.ToJson(users))

	err, _ = topUpNoSql.FindOne(ctx, user.TopUpID.Hex())
	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(topUpNoSql))

	bankcode = "BCA-01"
	bankname = "Bank Central Asiaa"

	err = topUpNoSql.Update(ctx, user.TopUpID.Hex(),
		bankcode, bankname)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(topUpNoSql))

	if topUpNoSql.BankCode != bankcode {
		t.Fatal("Expected=", bankcode, "actual=", topUpNoSql.BankCode)
	}

	// err = topUpNoSql.Delete(ctx, topUpNoSql.ID.Hex())
	// if err != nil {
	// 	t.Fatal(err)
	// }

}
