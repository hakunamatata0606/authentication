package config

type Config struct {
	TokenTimeout        int64
	RefreshTokenTimeout int64
}

const (
	defaultTokenTimeout        = 3600
	defaultRefreshTokenTimeout = 86400
)

var globalConfig = Config{
	TokenTimeout:        defaultTokenTimeout,
	RefreshTokenTimeout: defaultRefreshTokenTimeout,
}

func GetConfig() *Config {
	return &globalConfig
}

func init() {
	//load config

}
