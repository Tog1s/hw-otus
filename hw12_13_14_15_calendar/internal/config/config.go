package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger  Logger
	Storage Storage
	DB      DB
	Server  Server
	Grpc    Grpc
}

type Logger struct {
	Level  string
	Output string
}

type Storage struct {
	Type string
}

type DB struct {
	Host     string
	Port     string
	DBName   string `yaml:"dbName"`
	User     string
	Password string
}

type Server struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
}

type Grpc struct {
	Host string
	Port string
}

func New(path string) (*Config, error) {
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading config file %s: %w", path, err)
	}

	config := &Config{}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("error while parse yaml file %s: %w", path, err)
	}
	return config, nil
}

//nolint:unused
func (db *DB) getDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.DBName,
	)
}
