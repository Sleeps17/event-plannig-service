package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	configEnv = "CONFIG_PATH"
)

type Config struct {
	Env       string            `yaml:"env"`
	Secret    string            `yaml:"secret"`
	Server    ServerConfig      `yaml:"server"`
	Auth      GrpcServiceConfig `yaml:"auth"`
	Employees GrpcServiceConfig `yaml:"employees"`
	Events    GrpcServiceConfig `yaml:"events"`
	Rooms     GrpcServiceConfig `yaml:"rooms"`
}

type ServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type GrpcServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func MustLoad() *Config {
	configPath := os.Getenv(configEnv)
	if configPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &configPath); err != nil {
		panic(err)
	}

	return &cfg
}
