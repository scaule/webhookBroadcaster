package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"webhookBroadcaster/model"
	"webhookBroadcaster/user"
	"webhookBroadcaster/webhook"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key, secret, _ := r.BasicAuth()
	if key == "" || secret == "" {
		return
	}
	//Check if user exist with this secret
	repository, err := user.GetRepository()
	var res model.ResponseResult
	if err != nil {
		res.Error = "Error"
		json.NewEncoder(w).Encode(res)
		return
	}
	userFromDb, err := repository.FindUserByKeyAndSecret(key, secret)
	if err != nil || userFromDb == nil {
		return
	}

	webhookRepository := webhook.GetBrokerRepository()
	body, _ := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	webhook := webhook.Webhook{}
	webhook.User = userFromDb.Email
	webhook.Data = string(body)
	webhook.Type = vars["type"]
	webhookRepository.Produce(&webhook)
}
func ConsumeHandler(w http.ResponseWriter, r *http.Request) {
	webhookRepository := webhook.GetBrokerRepository()
	webhookRepository.Consume()
	println("test")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser user.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &newUser)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	repository, err := user.GetRepository()
	if err != nil {
		res.Error = "Error"
		json.NewEncoder(w).Encode(res)
		return
	}
	err = repository.Create(&newUser)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Registration Successful"
	json.NewEncoder(w).Encode(res)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var newUser user.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &newUser)
	if err != nil {
		log.Fatal(err)
	}

	var res model.ResponseResult
	repository, err := user.NewMongoRepository()
	if err != nil {
		res.Error = "Error"
		json.NewEncoder(w).Encode(res)
		return
	}
	token, err := repository.Login(&newUser)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(token)
}
