package api

import (
	"bootchamp-codeid/data"
	"bootchamp-codeid/db/mangodb"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestLogin_test(t *testing.T) {

	ctx := context.TODO()
	mongoClient, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mangodb.CloseMongoDb(ctx, mongoClient)

	LoginApi := LoginApi{}
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginApi.Login).Methods("POST")

	no_hp := "081317767015"
	password := "123123"

	request, err := http.NewRequest("POST", "/login", nil)
	request.SetBasicAuth(no_hp, password)
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, request)
	if err != nil {
		t.Fatal(err)
	}
	resp := recoder.Body.Bytes()
	log.Println(string(resp))

	newLogin := data.LoginNoSql{}
	err = json.Unmarshal(resp, &newLogin)
	if err != nil {
		t.Fatal(err)
	}
}
