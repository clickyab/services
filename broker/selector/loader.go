package selector

import (
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/broker/mock"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/safe"

	"github.com/clickyab/services/broker/rabbitmq"

	"context"

	"github.com/clickyab/services/initializer"
	"github.com/sirupsen/logrus"
)

type cfg struct {
}

func (cfg) Initialize() config.DescriptiveLayer {
	layer := config.NewDescriptiveLayer()
	layer.Add("application is in test mode and broker is not active", "services.broker.provider", "mock")
	return layer
}

func (cfg) Loaded() {
	provider := config.GetString("services.broker.provider")

	switch provider {
	case "mock":
		p := mock.GetChannelBroker()
		broker.SetActiveBroker(p)
		safe.GoRoutine(
			context.Background(),
			func() {
				ch := mock.GetChannel(10)
				for j := range ch {
					data, err := j.Encode()
					logrus.WithField("key", j.Key()).
						WithField("topic", j.Topic()).
						WithField("encode_err", err).
						Debug(string(data))
				}
			},
		)
	case "rabbitmq":
		initializer.Register(rabbitmq.NewRabbitMQInitializer(), 1)
		p := rabbitmq.NewRabbitBroker()
		broker.SetActiveBroker(p)
	default:
		logrus.Panicf("there is no valid broker configured , %s is not valid", provider)
	}
}

func init() {
	config.Register(&cfg{})
}
