package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Config struct {
	Settings Settings
	CtiAPI   CtiAPI
	HttpAPI  HttpAPI
	DB       DB
	Secret   string
	TokenTTL time.Duration
}

type Settings struct {
	LogLevel string
	UseAuth  bool
	TokenTTL int64
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

var ServerConfig *Config

func NewConfig() *Config {
	return &Config{}
}

func Init() {
	ServerConfig = NewConfig()
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

	c.Secret = "your-256-bit-secret"

	c.TokenTTL = time.Duration(c.Settings.TokenTTL) * time.Minute

	log.Debug().Msg("config file parsed - ok")

	return nil
}
