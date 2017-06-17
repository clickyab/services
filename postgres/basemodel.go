package postgres

import (
	"context"
	"database/sql"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/postgres/model"
	_ "github.com/lib/pq"
	gorp "gopkg.in/gorp.v2" // Make sure postgres is included in any build
)

var (
	dbmap *gorp.DbMap
	db    *sql.DB
	once  = sync.Once{}
	all   []initializer.Simple
)

// Hooker interface :))))) You have a dirty mind.
type Hooker interface {
	// AddHook is called after initialize only if the manager implement it
	AddHook()
}

type gorpLogger struct {
}

type modelsInitializer struct {
}

func (g gorpLogger) Printf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

// Initialize the modules, its safe to call this as many time as you want.
func (modelsInitializer) Initialize(ctx context.Context) {
	var err error
	db, err = sql.Open("postgres", cfg.DSN)
	assert.Nil(err)

	db.SetMaxIdleConns(cfg.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.MaxConnection)
	err = db.Ping()
	assert.Nil(err)

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	if cfg.DevelMode {
		logger := gorpLogger{}
		dbmap.TraceOn("[DB]", logger)
	} else {
		dbmap.TraceOff()
	}
	model.Initialize(db, dbmap)
	for i := range all {
		all[i].Initialize()

	}
	// If they are hooker call them.
	for i := range all {
		if h, ok := all[i].(Hooker); ok {
			h.AddHook()
		}
	}
	go func() {
		c := ctx.Done()
		if c == nil {
			return
		}
		<-c
		logrus.Debug("postgres finalized")
	}()
	logrus.Debug("postgres is ready")
}

// Register a new initializer module
func Register(m ...initializer.Simple) {
	all = append(all, m...)
}

func init() {
	initializer.Register(&modelsInitializer{}, 0)
}
