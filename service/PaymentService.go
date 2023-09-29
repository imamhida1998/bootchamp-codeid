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

type PaymentService struct {
	BaseService
	data.Pay
}

func CreatePaymentService(mongoClient *mongo.Client, db *sql.DB) PaymentService {
	payment := PaymentService{}
	payment.MongoClient = mongoClient
	payment.DB = db
	return payment
}

/*
PayID                      string `json:"PayID" bson:"PayID"`
MerchantNoVirtualAccount   string             `json:"merchant_no_virtual_account" bson:"merchant_no_virtual_account"`
MerchantNameVirtualAccount string             `json:"merchant_name_virtual_account" bson:"merchant_name_virtual_account"`
SourceNoVirtualAccount     string             `json:"source_no_virtual_account" bson:"source_no_virtual_account"`
SourceNameVirtualAccount   string             `json:"source_name_virtual_account" bson:"source_name_virtual_account"`
PayAmount                  float32            `json:"pay_amount" bson:"pay_amount"`
CreatedAt                  string             `json:"created_at" bson:"created_at"`
*/
func (Payment PaymentService) Payments(ctx context.Context, payVirtualAccount string, NoVirtualAccount string, PaymentAmount float32, pin string) (*data.Pay, *data.VATransaction, *data.VATransaction, error) {

	datavirtual := data.CreateVirtualAccount(Payment.DB)
	userAccount, err := datavirtual.FindOne(NoVirtualAccount)
	if err != nil {
		return nil, nil, nil, err
	}
	payAccount, err := datavirtual.FindOne(payVirtualAccount)
	if err != nil {
		return nil, nil, nil, err
	}
	if userAccount.Saldo < PaymentAmount {
		fmt.Println("saldo kurang")

		return nil, nil, nil, err
	} else {
		/*

			TransferID               primitive.ObjectID `json:"transfer_id" bson:"transfer_id"`
			SourceNoVirtualAccount   string             `json:"source_no_virtual_account" bson:"source_no_virtual_account"`
			SourceNameVirtualAccount string             `json:"source_name_virtual_account" bson:"source_name_virtual_account"`
			DesNoVirtualAccount      string             `json:"des_no_virtual_account" bson:"des_no_virtual_account"`
			DescNameVirtualAccount   string             `json:"desc_name_virtual_account" bson:"desc_name_virtual_account"`
			TransferAmount           float32            `json:"transfer_amount" bson:"transfer_amount"`
			CreateAt                 string             `json:"create_at" bson:"create_at"`
		*/
		epin := base64.StdEncoding.EncodeToString([]byte(pin))
		if payVirtualAccount == payAccount.Account_no && NoVirtualAccount == userAccount.Account_no && epin == userAccount.Pin {
			log.Println("Note Transaction")
			paymentNoSql := data.CreatePay(Payment.MongoClient)
			payments, err := paymentNoSql.AddPayment(ctx,
				payAccount.Account_no,
				payAccount.AccountName,
				userAccount.Account_no,
				userAccount.AccountName,
				PaymentAmount)
			if err != nil {
				return nil, nil, nil, err

			}
			accountsaldo := userAccount.Saldo
			paymentsaldo := payAccount.Saldo
			accountsaldo -= PaymentAmount
			paymentsaldo += PaymentAmount

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
				payAccount.Account_no,
				payAccount.No_hp,
				payAccount.Email,
				payAccount.AccountName,
				payAccount.Pin,
				payAccount.Password,
				paymentsaldo)
			if err != nil {
				return nil, nil, nil, err
			}

			description := "pay to " + payAccount.Account_no + "/" + payAccount.AccountName
			VirtualAccountransactioNoSql := data.CreateVaTransaction(Payment.MongoClient)
			virtualtraction, err := VirtualAccountransactioNoSql.AddVirtualAccountTransaction(ctx,
				NoVirtualAccount,
				datavirtual.AccountName,
				PaymentAmount,
				data.PAYMENT,
				payments.PayID.Hex(),
				description)
			if err != nil {
				return nil, nil, nil, err
			}
			description = "pay to from " + userAccount.Account_no + "/" + userAccount.AccountName
			Virtualmerchanttransaction, err := VirtualAccountransactioNoSql.AddVirtualAccountTransaction(ctx,
				payVirtualAccount,
				payAccount.AccountName,
				PaymentAmount,
				data.PAYMENT,
				payments.PayID.Hex(),
				description)
			if err != nil {
				return nil, nil, nil, err
			}

			return payments, virtualtraction, Virtualmerchanttransaction, err
		} else {
			log.Println("Data Tidak DItemukan")
			return nil, nil, nil, err

		}

	}

}
