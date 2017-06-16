package mysql

import (
	"database/sql"

	"github.com/fzerorubigd/services/mysql"
	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func newConnection(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}

func init() {
	mysql.RegisterConnectionFactory(newConnection)
}
