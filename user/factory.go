package user

import "github.com/spf13/viper"

const USER_REPO_TYPE_KEY = "USER_REPO_TYPE"

func GetRepository() (UserRepository, error) {

	repoType, ok := viper.Get(USER_REPO_TYPE_KEY).(string)
	if !ok {
		repoType = "default"
	}
	//Atm we have only Mongodb database
	switch repoType {
	case "mongodb":
		return NewMongoRepository()
	}

	//Return mongodb repo as default
	return NewMongoRepository()
}
