package rabbitmq

import (
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/random"

	"sync/atomic"

	"github.com/clickyab/services/safe"

	"context"

	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var prefetchCount = config.RegisterInt("services.broker.rabbitmq.prefetch", 100, "the prefetch count")

func (cn consumer) RegisterConsumer(consumer broker.Consumer) error {
	connRng = connRng.Next()
	conn := connRng.Value.(*ccn).Connection()
	c, err := conn.Channel()
	if err != nil {
		return err
	}
	err = c.ExchangeDeclare(
		exchange.String(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		return err
	}
	qu := consumer.Queue()
	if debug.Bool() {
		qu = "debug." + qu
	}
	q, err := c.QueueDeclare(qu, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// prefetch count
	// **WARNING**
	// If ignore this, then there is a problem with rabbit. prefetch all jobs for this worker then.
	// the next worker get nothing at all!
	// **WARNING**
	// TODO : limit on workers must match with this prefetch
	err = c.Qos(prefetchCount.Int(), 0, false)
	if err != nil {
		return err
	}

	topic := consumer.Topic()
	if debug.Bool() {
		topic = "debug." + topic
	}
	err = c.QueueBind(
		q.Name,            // queue name
		topic,             // routing key
		exchange.String(), // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	safe.ContinuesGoRoutine(kill, func(cnl context.CancelFunc) time.Duration {
		consumerTag := <-random.ID
		delivery, err := c.Consume(q.Name, consumerTag, false, false, false, false, nil)
		if err != nil {
			cnl()
			assert.Nil(err) // I know its somehow redundant.
			return 0
		}
		logrus.Debug("Worker started")
		cn.consume(kill, cnl, consumer.Consume(kill), c, delivery, consumerTag)
		return time.Second
	})
	return nil
}

func (consumer) consume(ctx context.Context, cnl context.CancelFunc, consumer chan<- broker.Delivery, c *amqp.Channel, delivery <-chan amqp.Delivery, consumerTag string) {
	atomic.SwapInt64(&hasConsumer, 1)
	done := ctx.Done()

	cErr := c.NotifyClose(make(chan *amqp.Error))
bigLoop:
	for {
		select {
		case job, ok := <-delivery:
			assert.True(ok, "[BUG] Channel is closed! why??")
			consumer <- &jsonDelivery{delivery: &job}
		case <-done:
			logrus.Debug("closing channel")
			// break the continues loop
			cnl()
			_ = c.Cancel(consumerTag, true)
			break bigLoop
		case e := <-cErr:
			logrus.Errorf("%T => %+v", *e, *e)
		}
	}
}
