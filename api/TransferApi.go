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

type TransferApi struct {
	BaseApi
	service.TransferService
	PaymentNo       string  `json:"payment_no" bson:"payment_no"`
	UserNo          string  `json:"user_no" bson:"user_no"`
	TransferAmounts float32 `json:"transfer_amounts" bson:"transfer_amount"`
	Pin             string  `json:"pin" bson:"pin"`
}

func (tfApi TransferApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, *sql.DB, context.Context) {
	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		tfApi.Error(w, err)
		return nil, nil, nil
	}
	db := postgresql.OpenConnection(config.IMAM)

	_, err = tfApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		tfApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, db, ctx

}
func (tfApi TransferApi) Getransfer(w http.ResponseWriter, r *http.Request) {
	mongoClient, _, ctx := tfApi.Connection(w, r)
	//_ = postgresql.OpenConnection(config.IMAM)
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	Modeltf := data.CreateTransfer(mongoClient)
	devices, err := Modeltf.ListTransfer(ctx)
	if err != nil {
		tfApi.Error(w, err)
		return
	} else {
		tfApi.Json(w, devices, http.StatusOK)
		return
	}
}

// func (tfApi TransferApi) UpdateTopUpApi(w http.ResponseWriter, r *http.Request)  {
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&tfApi)
// 	if err != nil {
// 		tfApi.Error(w, err)
// 		return
// 	}
// 	mongoClient , _ ,ctx := tfApi.Connection(w , r)
// 	defer mangodb.CloseMongoDb(ctx,mongoClient)

// 	modelUpdate := data.CreateTransfer(mongoClient)
// 	err = modelUpdate.Up(ctx,
// 		tfApi.TopUpID.Hex(),
// 		tfApi.BankCode,
// 		tfApi.BankName)
// 	if err != nil {
// 		tfApi.Error(w , err)
// 		return
// 	} else {
// 		tfApi.Json(w , tfApi, http.StatusOK)
// 		return
// 	}
// }

func (tfApi TransferApi) AddTransfer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tfApi)
	if err != nil {
		tfApi.Error(w, err)
		return
	}
	mongoClient, db, ctx := tfApi.Connection(w, r)
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	ModelRemoved := service.CreateTransferService(mongoClient, db)
	DataTransfer, _, _, err := ModelRemoved.Transfer(ctx,
		tfApi.PaymentNo,
		tfApi.UserNo,
		tfApi.TransferAmounts,
		tfApi.Pin)
	if err != nil {
		tfApi.Error(w, err)
		return
	} else {
		tfApi.Json(w, DataTransfer, http.StatusOK)
		return
	}

}
