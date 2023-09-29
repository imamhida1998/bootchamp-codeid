package api

import (
	"bootchamp-codeid/data"
	"bootchamp-codeid/db/postgresql"
	"bootchamp-codeid/db/postgresql/config"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestVirtualAccountApi_PostVirtualAccount(t *testing.T) {
	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()

	router := mux.NewRouter()
	virtualaccountApi := VirtualAccountApi{}
	router.HandleFunc("/virtualaccount", virtualaccountApi.PostVirtualAccount).Methods("POST")

	virtualAccount := data.VirtualAccount{}
	virtualAccount.AccountName = "Imam"

	body, err := json.Marshal(&virtualAccount)
	if err != nil {
		t.Fatal(err)
	}
	request, er := http.NewRequest("POST", "/virtualaccount", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(er)
	}
	token := "JDJhJDEwJE5PeVBUNTdhTy9kRjVmYUcwMER4eC5BTzI4cEtjeGF0S1RzMzhVWWNxMVhCeGNRUXpLN3RX"
	request.Header.Add("Authorization", token)
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, request)

	resp := recoder.Body.Bytes()
	log.Println(string(resp))
	if recoder.Code != http.StatusOK {
		t.Fatal(string(resp))
	}
	NewVirtualAccount := data.VirtualAccount{}
	err = json.Unmarshal(resp, &NewVirtualAccount)
	if err != nil {
		t.Fatal(err)
	}
	if NewVirtualAccount.AccountName != virtualaccountApi.AccountName {
		t.Fatal("Expected =", virtualaccountApi.AccountName, "/nActual", NewVirtualAccount.AccountName)
	}

}
func TestVirtualAccountApi_GetVirtualAccount(t *testing.T) {

	db := postgresql.OpenConnection(config.IMAM)
	defer db.Close()
	router := mux.NewRouter()
	branchApi := VirtualAccountApi{}
	router.HandleFunc("/branches", branchApi.GetVirtualAccount).Methods("GET")

	request, err := http.NewRequest("GET", "/branches", nil)
	token := "JDJhJDEwJE5PeVBUNTdhTy9kRjVmYUcwMER4eC5BTzI4cEtjeGF0S1RzMzhVWWNxMVhCeGNRUXpLN3RX"
	request.Header.Add("Authorization", token)
	recoder := httptest.NewRecorder()
	router.ServeHTTP(recoder, request)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	resp := recorder.Body.Bytes()
	log.Println(string(resp))
}

func TestVirtualAcount_Update(t *testing.T) {

}
