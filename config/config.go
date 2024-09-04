package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type TokenConfig struct {
	TokenTimeout        int64  `yaml:"token_timeout"`
	RefreshTokenTimeout int64  `yaml:"refresh_token_timeout"`
	Secret              string `yaml:"secret"`
}

type DbConfig struct {
	Driver string `yaml:"driver"`
	Addr   string `yaml:"addr"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Token      TokenConfig `yaml:"token"`
	Db         DbConfig    `yaml:"db"`
	Redis      RedisConfig `yaml:"redis"`
	ServerAddr string      `yaml:"server_addr"`
}

const (
	envFile                    = "env.yaml"
	defaultTokenTimeout        = 3600
	defaultRefreshTokenTimeout = 86400
)

var globalConfig = Config{
	Token: TokenConfig{
		TokenTimeout:        defaultTokenTimeout,
		RefreshTokenTimeout: defaultRefreshTokenTimeout,
		Secret:              "secret",
	},
	ServerAddr: "localhost:8080",
}

func GetConfig() *Config {
	return &globalConfig
}

func init() {
	//load config
	gp := os.Getenv("ROOT_DIR")

	configFile, err := os.ReadFile(filepath.Join(gp, envFile))
	if err != nil {
		log.Fatalf("fail to read env file - %s", err.Error())
	}
	err = yaml.Unmarshal(configFile, &globalConfig)
	if err != nil {
		log.Fatalf("fail to unmarshal config - %s", err.Error())
	}

}
