package user

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"os/user"
	"strings"
	"time"
)

const DB_KEY = "MONGO_DB"
const HOST_KEY = "MONGO_HOST"

const DB_DEFAULT = "ifhttt"
const HOST_DEFAULT = "mongodb://localhost:27017"

const COLLECTION = "users"

type repo struct {
	collection *mongo.Collection
}

//NewMongoRepository create new repository
func NewMongoRepository() (UserRepository, error) {
	//GET DB INFO
	db, ok := viper.Get(DB_KEY).(string)
	if !ok {
		db = DB_DEFAULT
	}
	host, ok := viper.Get(HOST_KEY).(string)
	if !ok {
		host = HOST_DEFAULT
	}
	//INIT DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	c := client.Database(db).Collection(COLLECTION)

	return &repo{
		collection: c,
	}, nil
}

func (r *repo) UserExistByEmail(email string) (bool, error) {
	var result user.User
	err := r.collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, err
		}
	}

	return true, nil
}

func (r *repo) UserExistByKeyAndSecret(key string, secret string) (bool, error) {

	var result user.User
	err := r.collection.FindOne(context.TODO(), bson.D{{"accesskey", key}, {"secret", secret}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, err
		}
	}

	return true, nil
}

func (r *repo) FindUserByKeyAndSecret(key string, secret string) (*User, error) {

	var user User
	err := r.collection.FindOne(context.TODO(), bson.D{{"accesskey", key}, {"secret", secret}}).Decode(&user)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, err
		}
	}

	return &user, nil
}

func (r *repo) Login(newUser *User) (string, error) {

	var user User
	err := r.collection.FindOne(context.TODO(), bson.D{{"email", newUser.Email}}).Decode(&user)

	if err != nil {
		return "", errors.New("Invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUser.Password))

	if err != nil {
		return "", errors.New("Invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	return tokenString, nil
}

func (r *repo) Create(user *User) error {
	//check if user already exist
	userExist, _ := r.UserExistByEmail(user.Email)

	if !userExist {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

		if err != nil {
			return errors.New("error While Creating User, Try Again")
		}

		user.Password = string(hash)
		//Generate secret
		user.Secret = uuid.New().String()
		user.AccessKey = uuid.New().String()
		user.AccessKey = strings.Replace(user.AccessKey, "-", "", -1)
		_, err = r.collection.InsertOne(context.TODO(), user)
		if err != nil {
			return errors.New("error While Creating User, Try Again")
		}

		return nil
	}

	return errors.New("user already exist")
}

func (r *repo) Delete(email string) error {
	//try to delete user
	_, err := r.collection.DeleteOne(context.TODO(), bson.D{{"email", email}})

	return err
}
