package rabbitmq

/*
import (
	"github.com/clickyab/services/broker"
)

// RabbitInterface is a interface for amqp basic functions
type RabbitInterface interface {
	MakeConnections(count int) ([]string, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool) error
	MakePublisher(connectionKey string) (string, error)
	Publish(in broker.Job, pubKey string) error
	// FinalizeWait is a function to wait for all publication to finish. after calling this,
	// must not call the PublishEvent
	FinalizeWait()
	RegisterConsumer(consumer broker.Consumer, prefetchCount int) error
}
*/
