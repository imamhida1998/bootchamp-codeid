package data

import (
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"

	//	"bootchamp-codeid/helpers"
	//	"log"
	"testing"
)

func TestBankAccount(t *testing.T) {
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()

	ModelBankAccount := CreateBankAccount(db)
	// err := ModelBankAccount.Migrate()
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }
	// err = ModelBankAccount.ResetData()
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	bank_account_no := "BCA01"
	bank_account_owner := "Imam"
	saldo := float32(50000)
	//created := time.Now().UTC().Format(time.RFC1123)

	err := ModelBankAccount.InsertBankAccount(
		bank_account_no,
		bank_account_owner,
		saldo)
	if err != nil {
		t.Fatal(err.Error())
	}
	// databank, err := ModelBankAccount.ReadBankAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("OnInserted", helpers.ToJson(databank))

	// saldo = float32(10000)

	// err = ModelBankAccount.UpdateBankAccount(
	// 	bank_account_no,
	// 	bank_account_owner,
	// 	saldo)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }
	// databank, err = ModelBankAccount.ReadBankAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//log.Println("OnUpdated", helpers.ToJson(databank))

}
