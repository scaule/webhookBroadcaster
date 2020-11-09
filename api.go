package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"webhookBroadcaster/controller"
)

func main() {
	//Init Viper with env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/register", controller.RegisterHandler).
		Methods("POST")
	router.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")
	router.HandleFunc("/webhook/{type}", controller.WebhookHandler).
		Methods("GET")
	router.HandleFunc("/consume", controller.ConsumeHandler).
		Methods("GET")
	log.Fatal(http.ListenAndServe(":8090", router))
}
