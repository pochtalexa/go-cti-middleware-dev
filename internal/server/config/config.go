package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"

	"github.com/pochtalexa/go-cti-middleware/internal/server/flags"
)

type Config struct {
	Settings Settings
	CtiAPI   CtiAPI
	HttpAPI  HttpAPI
	DB       DB
	Secret   string
	TokenTTL time.Duration
	WsConn   *websocket.Conn
}

type Settings struct {
	LogLevel string
	UseAuth  bool
	TokenTTL int64
	LogPath  string
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
	const op = "config.ReadConfigFile"
	fileName := flags.CfgPath

	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = toml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	c.DB.DBConn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBname)

	// TODO: хранить токен во вне
	c.Secret = "your-256-bit-secret"

	c.TokenTTL = time.Duration(c.Settings.TokenTTL) * time.Minute

	log.Debug().Msg("config file parsed - ok")

	return nil
}
