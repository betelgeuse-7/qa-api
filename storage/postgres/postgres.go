package postgres

import (
	"errors"
	"fmt"
	"os"

	"github.com/betelgeuse-7/qa/config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	ERROR_UNIQUE_VIOLATION pq.ErrorCode = pq.ErrorCode("23505")
	MAX_INT_VAL            int32        = 2147483647
)

type Postgres struct {
	cfg      *config.ConfigRelationalDB
	password string
	Db       *sqlx.DB
}

func New(cfg *config.ConfigRelationalDB) (*Postgres, error) {
	p := &Postgres{cfg: cfg}
	pwd, err := getPostgresPwdFromEnv()
	if err != nil {
		return nil, err
	}
	p.password = pwd
	return p, nil
}

func (p *Postgres) Connect() error {
	dsn := p.makeConnStr()
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}
	p.Db = db
	return nil
}

func (p *Postgres) makeConnStr() string {
	str := fmt.Sprintf("dbname=%s host=%s user=%s port=%d sslmode=%s password=%s", p.cfg.DbName, p.cfg.Host, p.cfg.User, p.cfg.Port, p.cfg.Ssl, p.password)
	return str
}

func getPostgresPwdFromEnv() (string, error) {
	pwd := os.Getenv("POSTGRES_PWD")
	if len(pwd) == 0 {
		return "", errors.New("environment variable 'POSTGRES_PWD' is not set")
	}
	return pwd, nil
}
