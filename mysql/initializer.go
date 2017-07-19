package mysql

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"

	"fmt"

	"strings"

	"math/rand"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/healthz"
	"github.com/clickyab/services/safe"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbmap  []*gorp.DbMap
	wdbmap  *gorp.DbMap
	rdb     []*sql.DB
	wdb     *sql.DB
	once    = sync.Once{}
	all     []Initializer
	factory func(string) (*sql.DB, error)
)

type initMysql struct {
}

type gorpLogger struct {
}

func (g gorpLogger) Printf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

// Initialize the modules, its safe to call this as many time as you want.
func (in *initMysql) Initialize(ctx context.Context) {
	once.Do(func() {
		var err error
		assert.NotNil(factory)

		uris := strings.Split(rdsnSlice.String(), ",")
		err = fetchReadConnection(uris)
		assert.Nil(err)

		wdb, err = factory(wdsn.String())
		assert.Nil(err)

		wdb.SetMaxIdleConns(maxIdleConnection.Int())
		wdb.SetMaxOpenConns(maxConnection.Int())

		err = wdb.Ping()
		assert.Nil(err)

		wdbmap = &gorp.DbMap{Db: wdb, Dialect: gorp.MySQLDialect{}}

		if develMode.Bool() {
			logger := gorpLogger{}
			wdbmap.TraceOn("[wdb]", logger)
		} else {
			wdbmap.TraceOff()
		}

		for i := range all {
			all[i].Initialize()
		}
		healthz.Register(in)
		logrus.Debug("mysql is ready.")
		go func() {
			c := ctx.Done()
			assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
			<-c
			for _, i := range rdb {
				assert.Nil(i.Close())
			}
			assert.Nil(wdb.Close())
			logrus.Debug("mysql finalized.")
		}()

		go updateRdbMap(ctx, uris)
	})
}

func updateRdbMap(ctx context.Context, uris []string) {
	ticker := time.NewTicker(rdbUpdateCD.Duration())
	for {
		select {
		case <-ticker.C:
			safe.Try(func() error {
				return fetchReadConnection(uris)
			}, maxRdbRetry.Duration())

		case <-ctx.Done():
			break
		}

	}
}

func fetchReadConnection(uris []string) error {
	rdbmap = []*gorp.DbMap{}
	rdb = []*sql.DB{}
	for _, i := range uris {
		r, err := factory(i)
		if err != nil {
			continue
		}

		r.SetMaxIdleConns(maxIdleConnection.Int())
		r.SetMaxOpenConns(maxConnection.Int())

		err = r.Ping()
		if err != nil {
			continue
		}

		currentRdbmap := &gorp.DbMap{Db: r, Dialect: gorp.MySQLDialect{}}

		if develMode.Bool() {
			logger := gorpLogger{}
			currentRdbmap.TraceOn("[rdb]", logger)
		} else {
			currentRdbmap.TraceOff()
		}

		rdbmap = append(rdbmap, currentRdbmap)
		rdb = append(rdb, r)
	}
	if len(rdbmap) < 1 {
		return errors.New("all read databases are down")
	}
	return nil
}

// Healthy return true if the databases are ok and ready for ping
func (in *initMysql) Healthy(context.Context) error {
	var rErr error
	for _, i := range rdb {
		rErr = i.Ping()
		if rErr != nil {
			break
		}
	}

	wErr := wdb.Ping()

	if rErr != nil || wErr != nil {
		return fmt.Errorf("mysql PING failed, read error was %s and write error was %s", rErr, wErr)
	}

	return nil
}

// Manager is a base manager for transaction model
type Manager struct {
	tx     *gorp.Transaction
	rdbmap []*gorp.DbMap
	rdb    []*sql.DB
	wdbmap *gorp.DbMap
	wdb    *sql.DB

	transaction bool
}

// InTransaction return true if this manager s in transaction
func (m *Manager) InTransaction() bool {
	return m.transaction
}

// Begin is for begin transaction
func (m *Manager) Begin() error {
	var err error
	if m.transaction {
		logrus.Panic("already in transaction")
	}
	m.sureDbMap()
	m.tx, err = m.wdbmap.Begin()
	if err == nil {
		m.transaction = true
	}
	return err
}

// Commit is for committing transaction. panic if transaction is not started
func (m *Manager) Commit() error {
	if !m.transaction {
		logrus.Panic("not in transaction")
	}
	err := m.tx.Commit()
	if err != nil {
		return err
	}
	m.tx = nil
	m.transaction = false
	return nil
}

// Rollback is for RollBack transaction. panic if transaction is not started
func (m *Manager) Rollback() error {
	if !m.transaction {
		logrus.Panic("Not in transaction")
	}
	err := m.tx.Rollback()

	if err != nil {
		return err
	}

	m.transaction = false
	return nil
}

func (m *Manager) sureDbMap() {
	if m.rdbmap == nil || m.wdbmap == nil {
		m.rdbmap = rdbmap
		m.wdbmap = wdbmap
	}
}

// GetRDbMap is for getting the current dbmap
func (m *Manager) GetRDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	m.sureDbMap()
	index := rand.Intn(len(rdbmap))
	return m.rdbmap[index]
}

// GetRSQLDB return the raw connection to database
func (m *Manager) GetRSQLDB() *sql.DB {
	if m.rdb == nil {
		m.rdb = rdb
	}

	index := rand.Intn(len(rdb))
	return m.rdb[index]
}

// GetWDbMap is for getting the current dbmap
func (m *Manager) GetWDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	m.sureDbMap()
	return m.wdbmap
}

// GetWSQLDB return the raw connection to database
func (m *Manager) GetWSQLDB() *sql.DB {
	if m.wdb == nil {
		m.wdb = wdb
	}

	return m.wdb
}

// GetProperDBMap try to get the current writer for development mode
func (m *Manager) GetProperDBMap() gorp.SqlExecutor {
	if develMode.Bool() {
		return m.GetWDbMap()
	}
	return m.GetRDbMap()
}

// Hijack try to hijack into a transaction
func (m *Manager) Hijack(ts gorp.SqlExecutor) error {
	if m.transaction {
		return errors.New("already in transaction")
	}
	t, ok := ts.(*gorp.Transaction)
	if !ok {
		return errors.New("there is no transaction to hijack")
	}

	m.transaction = true
	m.tx = t

	return nil
}

// AddTable registers the given interface type with gorp. The table name
// will be given the name of the TypeOf(i).  You must call this function,
// or AddTableWithName, for any struct type you wish to persist with
// the given DbMap.
//
// This operation is idempotent. If i's type is already mapped, the
// existing *TableMap is returned
func (m *Manager) AddTable(i interface{}) *gorp.TableMap {
	m.sureDbMap()
	return m.wdbmap.AddTable(i)
}

// AddTableWithName has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	m.sureDbMap()
	return m.wdbmap.AddTableWithName(i, name)
}

// AddTableWithNameAndSchema has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithNameAndSchema(i interface{}, schema string, name string) *gorp.TableMap {
	m.sureDbMap()
	return m.wdbmap.AddTableWithNameAndSchema(i, schema, name)
}

// TruncateTables try to truncate tables , useful for tests
func (m *Manager) TruncateTables(tbl string) error {
	m.sureDbMap()
	q := "TRUNCATE " + tbl
	_, err := m.wdbmap.Exec(q)
	return err
}

// Register a new initMysql module
func Register(m ...Initializer) {
	all = append(all, m...)
}

// RegisterConnectionFactory register a connection factory for sql connection.
func RegisterConnectionFactory(f func(string) (*sql.DB, error)) {
	factory = f
}

func init() {
	initializer.Register(&initMysql{}, 0)
}
