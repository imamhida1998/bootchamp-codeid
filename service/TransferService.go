package service

import (
	"bootchamp-codeid/data"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransferService struct {
	BaseService
}

func CreateTransferService(mongoClient *mongo.Client, db *sql.DB) TransferService {
	tf := TransferService{}
	tf.MongoClient = mongoClient
	tf.DB = db
	return tf
}

func (tf TransferService) Transfer(ctx context.Context, tfVirtualAccount string, NoVirtualAccount string, TransferAmount float32, pin string) (*data.Transfer, *data.VATransaction, *data.VATransaction, error) {

	datavirtual := data.CreateVirtualAccount(tf.DB)
	userAccount, err := datavirtual.FindOne(NoVirtualAccount)
	if err != nil {
		return nil, nil, nil, err
	}
	tfAccount, err := datavirtual.FindOne(tfVirtualAccount)
	if err != nil {
		return nil, nil, nil, err
	}
	if userAccount.Saldo < TransferAmount {
		fmt.Println("saldo kurang")

		return nil, nil, nil, err
	} else {
		epin := base64.StdEncoding.EncodeToString([]byte(pin))
		if tfVirtualAccount == tfAccount.Account_no && NoVirtualAccount == userAccount.Account_no && epin == userAccount.Pin {
			log.Println("Note Transaction")
			paymentNoSql := data.CreateTransfer(tf.MongoClient)
			transfers, err := paymentNoSql.AddTransfer(ctx,
				tfAccount.Account_no,
				tfAccount.AccountName,
				userAccount.Account_no,
				userAccount.AccountName,
				TransferAmount)
			if err != nil {
				return nil, nil, nil, err

			}
			accountsaldo := userAccount.Saldo
			tfsaldo := tfAccount.Saldo
			accountsaldo -= TransferAmount
			tfsaldo += TransferAmount

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
			err = datavirtual.UpdateVirtualAccount(
				tfAccount.Account_no,
				tfAccount.No_hp,
				tfAccount.Email,
				tfAccount.AccountName,
				tfAccount.Pin,
				tfAccount.Password,
				tfsaldo)
			if err != nil {
				return nil, nil, nil, err
			}

			description := "pay to " + tfAccount.Account_no + "/" + tfAccount.AccountName
			VirtualAccountransactioNoSql := data.CreateVaTransaction(tf.MongoClient)
			virtualtraction, err := VirtualAccountransactioNoSql.AddVirtualAccountTransaction(ctx,
				NoVirtualAccount,
				tfAccount.AccountName,
				TransferAmount,
				data.PAYMENT,
				transfers.TransferID.Hex(),
				description)
			if err != nil {
				return nil, nil, nil, err
			}
			description = "pay to from " + userAccount.Account_no + "/" + userAccount.AccountName
			Virtualmerchanttransaction, err := VirtualAccountransactioNoSql.AddVirtualAccountTransaction(ctx,
				tfVirtualAccount,
				tfAccount.AccountName,
				TransferAmount,
				data.PAYMENT,
				transfers.TransferID.Hex(),
				description)
			if err != nil {
				return nil, nil, nil, err
			}

			return transfers, virtualtraction, Virtualmerchanttransaction, err
		} else {
			log.Println("Data Tidak Ditemukan")
			return nil, nil, nil, err

		}

	}

}
