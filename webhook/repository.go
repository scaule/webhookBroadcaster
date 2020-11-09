package webhook

/* Broker interface*/
type Reader interface {
	Consume() (*Webhook, error)
}

type Writer interface {
	Produce(webhook *Webhook) (bool, error)
}

//Repository repository interface
type WebhookBrokerRepository interface {
	Reader
	Writer
}

/* DB interface */
type DBReader interface {
}

type DBWriter interface {
	create(webhook *Webhook) error
}

//Repository repository interface
type WebhookDbRepository interface {
	DBReader
	DBWriter
}
