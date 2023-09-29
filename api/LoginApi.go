package api

import (
	"bootchamp-codeid/data"
	"bootchamp-codeid/db/mangodb"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"bootchamp-codeid/service"
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

type LoginApi struct {
	BaseApi
	data.LoginNoSql
	NoHp string `json:"no_hp" bson:"no_hp"`
}

func (loginApi LoginApi) Login(w http.ResponseWriter, r *http.Request) {

	authorization := r.Header.Get("Authorization")
	authorization = strings.Replace(authorization, "Basic ", "", -1)
	auth, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		loginApi.Error(w, err)
		return
	}
	auth1 := strings.Split(string(auth), ":")

	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()
	if err != nil {
		loginApi.Error(w, err)
		return
	}

	defer mangodb.CloseMongoDb(ctx, mongoClient)

	loginService := service.CreateVirtualAccount(mongoClient, db)
	encodedToken, err := loginService.Login(ctx, auth1[0], auth1[1])
	if err != nil {
		loginApi.Error(w, err)
		return
	} else {

		data := map[string]string{
			"token": encodedToken,
		}
		loginApi.Json(w, data, http.StatusOK)
		return
	}

}
