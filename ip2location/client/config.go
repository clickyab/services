package client

import (
	"github.com/clickyab/services/config"
)

var (
	ip2lserver string
)

type cfgInitializer struct {
}

func (ci *cfgInitializer) Initialize() config.DescriptiveLayer {
	l := config.NewDescriptiveLayer()
	l.Add("IP to location service port", "services.ip2location.client.endpoint", "127.0.0.1:8190")
	return l
}

func (ci *cfgInitializer) Loaded() {
	ip2lserver = config.GetStringDefault("services.ip2location.client.endpoint", "127.0.0.1:8190")
}

func init() {
	config.Register(&cfgInitializer{})
}
