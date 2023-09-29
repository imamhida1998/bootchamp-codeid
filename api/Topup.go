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

type TopupApi struct {
	BaseApi
	service.TopUpService
}

func (TopUpApi TopupApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		TopUpApi.Error(w, err)
		return nil, nil, nil
	}

	db := postgresql.OpenConnection(config.IMAM)

	_, err = TopUpApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		TopUpApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}
func (TopUpApi TopupApi) GetTopUp(w http.ResponseWriter, r *http.Request) {
	mongoClient, ctx, db := TopUpApi.Connection(w, r)
	//_ = postgresql.OpenConnection(config.IMAM)
	defer mangodb.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	ModelTopup := data.CreateTopUp(mongoClient)
	devices, err := ModelTopup.ListTopUps(ctx)
	if err != nil {
		TopUpApi.Error(w, err)
		return
	} else {
		TopUpApi.Json(w, devices, http.StatusOK)
		return
	}
}

//func (TopUpApi TopupApi) CreateTopUpApi(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	err := decoder.Decode(&TopUpApi)
//	if err != nil {
//		TopUpApi.Error(w, err)
//		return
//	}
//	mongoClient, _, ctx := TopUpApi.Connection(w, r)
//	defer mangodb.CloseMongoDb(ctx, mongoClient)
//
//	ModelCreate := data.CreateTopUp(mongoClient)
//	Created, err := ModelCreate.AddTopUp(ctx,
//		TopUpApi.BankCode,
//		TopUpApi.BankName,
//		TopUpApi.NoVirtualAccount,
//		TopUpApi.NameVirtualAccount,
//		TopUpApi.TopUpAmount)
//	if err != nil {
//		TopUpApi.Error(w, err)
//		return
//	} else {
//		TopUpApi.Json(w, Created, http.StatusOK)
//		return
//	}
//
//}

func (TopUpApi TopupApi) TopUpVirtualAccount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&TopUpApi)
	if err != nil {
		TopUpApi.Error(w, err)
		return
	}
	mongoClient, ctx, db := TopUpApi.Connection(w, r)

	//_ = postgresql.OpenConnection(config.IMAM)
	defer mangodb.CloseMongoDb(ctx, mongoClient)
	defer db.Close()
	ctx = context.TODO()
	ModelTopUpVA := service.CreateTopUpService(mongoClient, db)
	_, _, TopUpVA, err := ModelTopUpVA.TopUp(ctx, TopUpApi.BankCode, TopUpApi.NoVirtualAccount, TopUpApi.TopUpAmount, TopUpApi.Pin)
	if err != nil {
		TopUpApi.Error(w, err)
	} else {
		TopUpApi.Json(w, TopUpVA, http.StatusOK)
		return
	}

}

//func (TopUpApi TopupApi) UpdateTopUpApi(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	err := decoder.Decode(&TopUpApi)
//	if err != nil {
//		TopUpApi.Error(w, err)
//		return
//	}
//	mongoClient, _, ctx := TopUpApi.Connection(w, r)
//	defer mangodb.CloseMongoDb(ctx, mongoClient)
//
//	modelUpdate := service.CreateTopUp(mongoClient)
//	err = modelUpdate.(ctx,
//		TopUpApi.TopUpID.Hex(),
//		TopUpApi.BankCode,
//		TopUpApi.BankName)
//	if err != nil {
//		TopUpApi.Error(w, err)
//		return
//	} else {
//		TopUpApi.Json(w, TopUpApi, http.StatusOK)
//		return
//	}
//}
//
//func (TopUpApi TopupApi) RemoveTopUpApi(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	err := decoder.Decode(&TopUpApi)
//	if err != nil {
//		TopUpApi.Error(w, err)
//		return
//	}
//	mongoClient, _, ctx := TopUpApi.Connection(w, r)
//	defer mangodb.CloseMongoDb(ctx, mongoClient)
//
//	ModelRemoved := data.CreateTopUp(mongoClient)
//	Removed, err := ModelRemoved.Delete(ctx, TopUpApi.TopUpID.Hex())
//	if err != nil {
//		TopUpApi.Error(w, err)
//		return
//	} else {
//		TopUpApi.Json(w, Removed, http.StatusOK)
//		return
//	}
//
//}
