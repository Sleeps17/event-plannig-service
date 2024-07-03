package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	configEnv = "CONFIG_PATH"
)

type Config struct {
	Env     string         `yaml:"env"`
	Storage PostgresConfig `yaml:"storage"`
	Server  GrpcConfig     `yaml:"server"`
}

type GrpcConfig struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Database string        `yaml:"database"`
	Timeout  time.Duration `yaml:"timeout"`
}

func (c *PostgresConfig) Url() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func MustLoad() *Config {
	configPath := os.Getenv(configEnv)
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("failed read config: %v", err))
	}

	return &cfg
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
