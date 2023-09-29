package data

import (
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"fmt"
	"testing"
)

func TestVirtualAccount(t *testing.T) {
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()

	ModelVirtalAccount := CreateVirtualAccount(db)
	//err := ModelVirtalAccount.Migrate()
	//if err != nil {
	//	t.Fatal(err.Error())
	//}

	//err = ModelVirtalAccount.ResetData()
	//if err != nil {
	//	t.Fatal(err.Error())
	//}

	No_hp := "081317767015"

	PaddingMask := fmt.Sprintf("%06d", 1)

	Account_no := No_hp + PaddingMask

	Email := "Yudis@gmail.com"
	AccountName := "Imam"
	Pin := "123456"
	Password := "hello123"
	saldo := float32(30000)

	err, _ := ModelVirtalAccount.InsertVirtualAccount(Account_no,
		No_hp,
		Email,
		AccountName,
		Pin,
		Password,
		saldo)
	if err != nil {
		t.Fatal(err.Error())
	}
	// datavirtual, err := ModelVirtalAccount.ReadVirtualAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("OnInserted", helpers.ToJson(datavirtual))

	// // //No_hp = "081317767016"
	// // Email = "imam@gmail.com"
	// // AccountName = "Imam"
	// // Pin = "123455"
	// // //Password = "hello123"
	// // saldo = float32(4000)

	// // err = ModelVirtalAccount.UpdateVirtualAccount(
	// // 	Account_no,
	// // 	No_hp,
	// // 	Email,
	// // 	AccountName,
	// // 	Pin,
	// // 	Password,
	// // 	saldo)
	// // if err != nil {
	// // 	t.Fatal(err.Error())
	// // }

	// datahp, err := ModelVirtalAccount.FindByPhone(No_hp)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("Cek Hp", helpers.ToJson(datahp))

	// datavirtual, err = ModelVirtalAccount.ReadVirtualAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("OnUpdated", helpers.ToJson(datavirtual))

	// //	AccountName = "Imam Hidayat"
	// //	Pin = "123321"

	// // err = ModelVirtalAccount.UpdateAccountPin(
	// // 	Account_no,
	// // 	AccountName,
	// // 	Pin)
	// // if err != nil {
	// // 	t.Fatal(err.Error())
	// // }
	// datavirtual, err = ModelVirtalAccount.ReadVirtualAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("OnUpdated", helpers.ToJson(datavirtual))
	// //	fmt.Sprintf(`%08d`, 5)

	// // hapus data
	// // err = ModelVirtalAccount.Remove(Account_no)
	// // if err != nil {
	// // 	t.Fatal(err.Error())
	// // }
}
