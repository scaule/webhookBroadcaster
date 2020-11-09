package webhook

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"errors"
)

const MONGO_DB_KEY = "MONGO_DB"
const MONGO_HOST_KEY = "MONGO_HOST"

const DB_DEFAULT = "ifhttt"
const MONGO_HOST_DEFAULT = "mongodb://localhost:27017"

const COLLECTION = "webhooks"

type mongoRepo struct {
	collection *mongo.Collection
}

//NewMongoRepository create new repository
func NewMongoRepository() (WebhookDbRepository, error) {
	//GET DB INFO
	db, ok := viper.Get(MONGO_HOST_KEY).(string)
	if !ok {
		db = DB_DEFAULT
	}
	host, ok := viper.Get(MONGO_DB_KEY).(string)
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

	return &mongoRepo{
		collection: c,
	}, nil
}

func (r *mongoRepo) create(webhook *Webhook) (error) {
	_, err := r.collection.InsertOne(context.TODO(), webhook)

	if err != nil {
		return errors.New("error While Creating User, Try Again")
	}

	return nil
}