package database

import (
	"fmt"
	"qa/conf"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type DbType int

const (
	POSTGRES DbType = iota
)

type Database interface {
	Connect() (*sqlx.DB, error)
}

type Postgres struct {
	User, DbName, Host, Password string
	Port                         int
	SslMode                      bool
}

func NewDatabase(engine DbType, config *conf.Config) Database {
	if engine == POSTGRES {
		return &Postgres{
			User:     config.DB.User,
			Password: config.DB.Password,
			DbName:   config.DB.DbName,
			Host:     config.DB.Host,
			Port:     config.DB.Port,
			SslMode:  config.DB.SslMode,
		}
	}
	panic("DB engine not supported")
}

func (p *Postgres) connString() string {
	cs := fmt.Sprintf("user=%s dbname=%s host=%s port=%d password=%s", p.User, p.DbName, p.Host, p.Port, p.Password)

	if p.SslMode {
		cs += " sslmode=enable"
	} else {
		cs += " sslmode=disable"
	}

	return cs
}

func (p *Postgres) Connect() (*sqlx.DB, error) {
	connStr := p.connString()
	return sqlx.Open("postgres", connStr)
}
