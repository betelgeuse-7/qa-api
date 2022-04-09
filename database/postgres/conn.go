package postgres

import (
	"fmt"
	"qa/conf"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	User, DbName, Host, Password string
	Port                         int
	SslMode                      bool
}

func NewPostgres(config *conf.Config) *Postgres {
	return &Postgres{
		User:     config.DB.User,
		Password: config.DB.Password,
		DbName:   config.DB.DbName,
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		SslMode:  config.DB.SslMode,
	}
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
