package postgres

import (
	"os"

	"github.com/clickyab/services/config"
)

type configInitializer struct {
	DSN               string
	MaxConnection     int
	MaxIdleConnection int

	DevelMode bool
}

var cfg configInitializer

func (cfg *configInitializer) Initialize() config.DescriptiveLayer {

	d := config.NewDescriptiveLayer()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:bita123@localhost/malooch"
	}
	dsn += "?sslmode=disable"
	d.Add("postgres DSN", "service.postgres.dsn", dsn)
	d.Add("postgres maximum connection", "service.postgres.max_connection", 150)
	d.Add("postgres max idle connection", "service.postgres.max_idle_connection", 10)

	return d
}

func (cfg *configInitializer) Loaded() {
	cfg.DSN = config.GetString("service.postgres.dsn")
	cfg.MaxConnection = config.GetInt("service.postgres.max_connection")
	cfg.MaxIdleConnection = config.GetInt("service.postgres.max_idle_connection")
	cfg.DevelMode = config.GetBool("core.devel_mode")
}

func init() {
	config.Register(&cfg)
}
