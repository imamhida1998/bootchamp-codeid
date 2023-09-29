package main

import (
	"bootchamp-codeid/api"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome server"))
	}).Methods("GET")

	LoginApi := api.LoginApi{}
	router.HandleFunc("/login", LoginApi.Login).Methods("POST")

	VirtualAccountAPI := api.VirtualAccountApi{}
	router.HandleFunc("/virtualaccount", VirtualAccountAPI.GetVirtualAccount).Methods("GET")

	BankAPI := api.BankAccountApi{}
	router.HandleFunc("/bankaccout", BankAPI.GetBankAccount).Methods("GET")

	PaymentAPI := api.PaymentApi{}
	router.HandleFunc("/payment", PaymentAPI.GetPayment).Methods("GET")
	router.HandleFunc("/payment", PaymentAPI.CreatePaymentApi).Methods("POST")

	TopUpAPI := api.TopupApi{}
	router.HandleFunc("/topup", TopUpAPI.TopUpVirtualAccount).Methods("POST")

	TransferAPI := api.TransferApi{}
	router.HandleFunc("/transfer", TransferAPI.Getransfer).Methods("GET")
	router.HandleFunc("/transfer", TransferAPI.AddTransfer).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("Server listening to 127.0.0.1:8005")
	err := http.ListenAndServe("127.0.0.1:8005", handlers.CORS(headers, methods, origins)(router))
	if err != nil {
		log.Fatal(err)
	}
}
