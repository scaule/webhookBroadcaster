package webhook

import "github.com/spf13/viper"

const WEBHOOK_REPO_TYPE_KEY = "WEBHOOK_REPO_TYPE"

func GetBrokerRepository() WebhookBrokerRepository {

	repoType, ok := viper.Get(WEBHOOK_REPO_TYPE_KEY).(string)
	if !ok {
		repoType = "default"
	}
	//Atm we have only kafka Broker
	switch repoType {
	case "mongodb":
		return NewKafkaRepository()
	}

	//Return mongodb repo as default
	return NewKafkaRepository()
}

func GetDBRepository() (WebhookDbRepository, error) {

	repoType, ok := viper.Get(WEBHOOK_REPO_TYPE_KEY).(string)
	if !ok {
		repoType = "default"
	}
	//Atm we have only kafka Broker
	switch repoType {
	case "mongodb":
		return NewMongoRepository()
	}

	//Return mongodb repo as default
	return NewMongoRepository()
}
