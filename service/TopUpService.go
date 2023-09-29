package service

import (
	"bootchamp-codeid/data"
	"context"
	"database/sql"
	"encoding/base64"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type TopUpService struct {
	BaseService
	BankCode         string  `json:"bank_code"`
	NoVirtualAccount string  `json:"no_virtual_account"`
	TopUpAmount      float32 `json:"top_up_amount"`
	Pin              string  `json:"pin"`
}

func CreateTopUpService(mongoClient *mongo.Client, db *sql.DB) TopUpService {
	topup := TopUpService{}
	topup.MongoClient = mongoClient
	topup.DB = db
	return topup
}
func (topup TopUpService) TopUp(ctx context.Context, BankCode string, NoVirtualAccount string, TopUpAmount float32, pin string) (*data.VATransaction, *data.BankTransaction, *data.TopUp, error) {
	databank := data.CreateBankAccount(topup.DB)
	userBank, err := databank.FindOne(BankCode)
	if err != nil {
		return nil, nil, nil, err
	}
	datavirtual := data.CreateVirtualAccount(topup.DB)
	userAccount, err := datavirtual.FindOne(NoVirtualAccount)
	if err != nil {
		return nil, nil, nil, err
	}
	epin := base64.StdEncoding.EncodeToString([]byte(pin))
	if BankCode == userBank.BankAccount_No && NoVirtualAccount == userAccount.Account_no && epin == userAccount.Pin {
		topupNoSql := data.CreateTopUp(topup.MongoClient)
		topups, err := topupNoSql.AddTopUp(ctx,
			BankCode,
			userBank.BankAccount_Owner,
			userAccount.Account_no,
			userAccount.AccountName,
			TopUpAmount)
		if err != nil {
			return nil, nil, nil, err
		}
		accountsaldo := userAccount.Saldo
		banksaldo := userBank.Saldo
		accountsaldo += TopUpAmount
		banksaldo += TopUpAmount

		err = datavirtual.UpdateVirtualAccount(
			userAccount.Account_no,
			userAccount.No_hp,
			userAccount.Email,
			userAccount.AccountName,
			userAccount.Pin,
			userAccount.Password,
			accountsaldo)

		if err != nil {
			return nil, nil, nil, err
		}
		err = databank.UpdateBankAccount(
			userBank.BankAccount_No,
			userBank.BankAccount_Owner,
			banksaldo)
		if err != nil {
			return nil, nil, nil, err
		}

		description := BankCode + "/" + userBank.BankAccount_Owner
		VirtualAccountransactioNoSql := data.CreateVaTransaction(topup.MongoClient)
		virtualtraction, err := VirtualAccountransactioNoSql.AddVirtualAccountTransaction(ctx,
			NoVirtualAccount,
			datavirtual.AccountName,
			TopUpAmount,
			data.TOPUP,
			topups.TopUpID.Hex(),
			description)
		if err != nil {
			return nil, nil, nil, err
		}

		/*

			TransactionID primitive.ObjectID `json:"transaction_id" bson:"transaction_id"`
			BankAccount_No string`json:"NoBankAccount" bson:"NoBankAccount"`
			BankAccount_Name  string`json:"NameBankAccount" bson:"NameBankAccount"`
			TransactionAmount float32 `json:"transaction_amount" bson:"transaction_amount"`
			UpdateAt string `json:"update_at" bson:"update_at"`*/
		reference := NoVirtualAccount + " | " + userAccount.AccountName
		Banktransactionnosql := data.CreateBankTransaction(topup.MongoClient)
		virtualbanktransaction, err := Banktransactionnosql.InsertBankTransaction(ctx,
			databank.BankAccount_No,
			databank.BankAccount_Owner,
			TopUpAmount,
			reference,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		return virtualtraction, virtualbanktransaction, topups, err

	} else {
		log.Println("Data Tidak Ditemukan")
		return nil, nil, nil, err
	}
	//encryptedToken, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	//if err != nil {
	//	return nil,nil,nil,err
	//}
	//encodedToken := base64.StdEncoding.EncodeToString(encryptedToken)
}
