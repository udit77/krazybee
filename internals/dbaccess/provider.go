package dbaccess

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Provider interface {
	GetDBConn() *sql.DB
}