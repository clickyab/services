package postgres

import (
	"os"
	"time"

	"github.com/clickyab/services/config"
)

type configInitializer struct {
	DSN               string
	MaxConnection     int
	MaxIdleConnection int

	DevelMode bool
}

var retryMax = config.RegisterDuration("services.postgres.max_retry_connection", time.Minute, "max time app should fallback to get mysql connection")
var startupMigration = config.RegisterBoolean("services.postgres.startup_migration", false, "do a migration on startup")

var cfg configInitializer

func (cfg *configInitializer) Initialize() config.DescriptiveLayer {

	d := config.NewDescriptiveLayer()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:bita123@localhost/malooch"
	}
	dsn += "?sslmode=disable"
	d.Add("postgres DSN", "services.postgres.dsn", dsn)
	d.Add("postgres maximum connection", "services.postgres.max_connection", 150)
	d.Add("postgres max idle connection", "services.postgres.max_idle_connection", 10)

	return d
}

func (cfg *configInitializer) Loaded() {
	cfg.DSN = config.GetString("services.postgres.dsn")
	cfg.MaxConnection = config.GetInt("services.postgres.max_connection")
	cfg.MaxIdleConnection = config.GetInt("services.postgres.max_idle_connection")
	cfg.DevelMode = config.GetBool("core.devel_mode")
}

func init() {
	config.Register(&cfg)
}
