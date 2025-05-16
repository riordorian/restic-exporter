package config

import (
	v "github.com/spf13/viper"
	"log"
	"restic-exporter/internal/shared/interfaces"
)

func InitConfig() interfaces.ConfigProviderInterface {
	viper := v.GetViper()
	viper.AddConfigPath("../")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
		panic(err.Error())
	}

	//viper.SetDefault("BASE_PATH", "/mnt/nfs/")
	viper.SetDefault("BASE_PATH", "/Users/riordorian/Downloads/restic")
	viper.SetDefault("METRICS_COLLECTING_INTERVAL_SECONDS", "60")
	viper.SetDefault("EXPOSE_PORT", "8085")
	viper.SetDefault("RESTIC_PASSWORD_COMMAND", "pass")

	return viper
}
