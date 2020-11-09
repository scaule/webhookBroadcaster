package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const TOPIC_KEY = "KAFKA_TOPIC"
const HOST_KEY = "KAFKA_HOST"

const TOPIC_DEFAULT = "webhook"
const HOST_DEFAULT = "localhost:9092"

type repo struct {
	topic string
	host  string
}

//NewKafkaRepository create new repository
func NewKafkaRepository() WebhookBrokerRepository {
	topic, ok := viper.Get(TOPIC_KEY).(string)
	if !ok {
		topic = TOPIC_DEFAULT
	}
	host, ok := viper.Get(TOPIC_KEY).(string)
	if !ok {
		host = HOST_DEFAULT
	}

	return &repo{
		topic,
		host,
	}
}

func (r *repo) Consume() (*Webhook, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": r.host,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{r.topic, "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			var webhook Webhook
			if err := json.Unmarshal(msg.Value, &webhook); err != nil {
				panic(err)
			}
			println("test")
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()

	return nil, nil
}

func (r *repo) Produce(webhook *Webhook) (bool, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": r.host})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	//Transform webhook to json
	jWebhook, err := json.Marshal(webhook)
	if err != nil {
		return false, errors.New("Can't format webhook as json")
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &r.topic, Partition: kafka.PartitionAny},
		Value:          []byte(jWebhook),
	}, nil)
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

	//TODO: ADD CIRCUIT BREAKER HERE TO PREVENT ERROR ON KAFKA PUSH
	return true, nil
}
