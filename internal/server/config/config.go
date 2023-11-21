package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	Settings Settings
	CtiAPI   CtiAPI
	HttpAPI  HttpAPI
	DB       DB
}

type Settings struct {
	LogLevel string
	UseAuth  bool
}

type CtiAPI struct {
	Scheme string
	Path   string
	Host   string
	Port   string
}

type HttpAPI struct {
	Scheme string
	Host   string
	Port   string
}

type DB struct {
	Host     string
	Port     string
	DBname   string
	User     string
	Password string
	DBConn   string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ReadConfigFile() error {
	fileName := "config.toml"

	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	if err := toml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	c.DB.DBConn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBname)

	log.Debug().Msg("config file parsed - ok")

	return nil
}
