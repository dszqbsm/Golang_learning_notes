package main

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/go-sql-driver/mysql"
)

type Connector interface {
	Connect(context.Context) (Conn, error)
	Driver() Driver
}

func OpenDB(c driver.Connector) *DB {

}

func main() {
	connector, err := mysql.NewConnector(&mysql.Config{
		User:      "root",
		Passwd:    "123456",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "test",
		ParseTime: true,
	})

	db := sql.OpenDB(connector)
}
