package util

import "github.com/spf13/viper"

type Config struct {
	DbDriver      string `mapstucture:"DB_DRIVER"`
	DbSource      string `mapstucture:"DB_SOURCE"`
	ServerAddress string `mapstucture:"ADDRESS_SERVER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
