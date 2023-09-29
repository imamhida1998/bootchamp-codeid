package api

import (
	"bootchamp-codeid/data"
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"bootchamp-codeid/service"
	"context"
	"database/sql"
	"encoding/json"

	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentApi struct {
	BaseApi
	service.PaymentService
	PayID string `json:"pay_id" bson:"pay_id"`
}

func (PayApi PaymentApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		PayApi.Error(w, err)
		return nil, nil, nil
	}

	db := postgresql.OpenConnection(config.IMAM)

	_, err = PayApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		PayApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}

func (PayApi PaymentApi) CreatePaymentApi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&PayApi)
	if err != nil {
		PayApi.Error(w, err)
		return
	}
	//	ctx := context.Background()
	mongoClient, ctx, _ := PayApi.Connection(w, r)
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	ModelCreate := data.CreatePay(mongoClient)
	Create, err := ModelCreate.AddPayment(ctx,
		PayApi.MerchantNoVirtualAccount,
		PayApi.MerchantNameVirtualAccount,
		PayApi.SourceNoVirtualAccount,
		PayApi.SourceNameVirtualAccount,
		PayApi.PayAmount)
	if err != nil {
		PayApi.Error(w, err)
		return
	} else {
		PayApi.Json(w, Create, http.StatusOK)
		return
	}

}

func (PayApi PaymentApi) GetPayment(w http.ResponseWriter, r *http.Request) {
	mongoClient, ctx, _ := PayApi.Connection(w, r)
	//_ = postgresql.OpenConnection(config.IMAM)
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	deviceModel := data.CreatePay(mongoClient)
	devices, err := deviceModel.ListPayment(ctx)
	if err != nil {
		PayApi.Error(w, err)
		return
	} else {
		PayApi.Json(w, devices, http.StatusOK)
		return
	}
}

/*
pay_id,

	merchant_va_account_no,
	merchant_va_account_name,
	src_va_account_no,
	src_va_account_name,
	pay_amount,
	created_at,
*/
func (PayApi PaymentApi) UpdatePaymentApi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&PayApi)
	if err != nil {
		PayApi.Error(w, err)
		return
	}
	mongoClient, ctx, _ := PayApi.Connection(w, r)
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	modelUpdate := data.CreatePay(mongoClient)
	err = modelUpdate.Update(ctx,
		PayApi.PayID,
		PayApi.MerchantNoVirtualAccount,
		PayApi.MerchantNameVirtualAccount,
		PayApi.SourceNoVirtualAccount,
		PayApi.SourceNameVirtualAccount,
		PayApi.PayAmount)
	if err != nil {
		PayApi.Error(w, err)
		return
	} else {
		PayApi.Json(w, PayApi, http.StatusOK)
		return
	}
}

func (PayApi PaymentApi) RemovePaymentApi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&PayApi)
	if err != nil {
		PayApi.Error(w, err)
		return
	}
	mongoClient, ctx, _ := PayApi.Connection(w, r)
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	ModelRemoved := data.CreatePay(mongoClient)
	Removed, err := ModelRemoved.Delete(ctx, PayApi.PayID)
	if err != nil {
		PayApi.Error(w, err)
		return
	} else {
		PayApi.Json(w, Removed, http.StatusOK)
		return
	}

}
