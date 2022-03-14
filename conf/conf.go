package conf

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const DEFAULT_PORT = 8000

type Config struct {
	Port int
	DB   struct {
		Port                         int
		User, DbName, Host, Password string
		SslMode                      bool
	}
}

func NewConfig(dbSslMode bool) (*Config, error) {
	godotenv.Load()

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	return &Config{
		Port: DEFAULT_PORT,
		DB: struct {
			Port     int
			User     string
			DbName   string
			Host     string
			Password string
			SslMode  bool
		}{
			Port:     dbPort,
			User:     dbUser,
			DbName:   dbName,
			Password: dbPassword,
			Host:     dbHost,
			SslMode:  dbSslMode,
		},
	}, nil
}

func (c *Config) SetPort(port int) {
	c.Port = port
}
