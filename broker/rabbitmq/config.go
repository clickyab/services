package rabbitmq

import (
	"time"

	"github.com/clickyab/services/config"
)

var (
	dsn        = config.RegisterString("services.amqp.dsn", "amqp://server:bita123@127.0.0.1:5672/cy", "amqp dsn")
	exchange   = config.RegisterString("services.amqp.exchange", "cy", "amqp exchange to publish into")
	publisher  = config.RegisterInt("services.amqp.publisher", 30, "amqp publisher to publish into")
	connection = config.RegisterInt("services.amqp.connection.count", 5, "connection count to prevent the limited connection")
	confirmLen = config.RegisterInt("services.amqp.confirm_len", 200, "amqp confirm channel len")
	debug      = config.RegisterBoolean("services.amqp.debug", false, "amqp debug mode")
	tryLimit   = config.RegisterDuration("services.amqp.try_limit", time.Second * 5, "the limit to incremental try wait")
)
