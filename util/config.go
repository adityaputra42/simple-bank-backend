package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment         string        `mapstructure:"ENVIRONMENT"`
	AllowedOrigins      []string      `mapstructure:"ALLOWED_ORIGINS"`
	DbDriver            string        `mapstructure:"DB_DRIVER"`
	DbSource            string        `mapstructure:"DB_SOURCE"`
	MigrationUrl        string        `mapstructure:"MIGRATION_URL"`
	HttpServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress   string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	RedisAddress        string        `mapstructure:"REDIS_ADDRESS"`
	TokenSymetricKey    string        `mapstructure:"TOKEN_SYMENTRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	EmailSenderName     string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
