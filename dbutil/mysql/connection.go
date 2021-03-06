package mysql

import (
	"database/sql"
	"fmt"

	"github.com/levinholsety/common-go/dbutil"
)

// Connection represents a connection of MySQL.
type Connection struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

var _ interface {
	dbutil.Connection
	fmt.Stringer
} = (*Connection)(nil)

// Open opens MySQL connection and creates a sql.DB instance for using.
// Should import _ "github.com/go-sql-driver/mysql".
func (p *Connection) Open() (*sql.DB, error) {
	return sql.Open("mysql", p.String())
}

func (p *Connection) String() string {
	host := p.Host
	if len(host) == 0 {
		host = "localhost"
	}
	port := p.Port
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", p.User, p.Password, host, port, p.Database)
}
