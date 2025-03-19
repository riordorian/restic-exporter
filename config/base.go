package config

import (
	v "github.com/spf13/viper"
	"grpc/internal/shared/interfaces"
	"log"
)

func InitConfig() interfaces.ConfigProviderInterface {
	viper := v.GetViper()
	viper.AddConfigPath("../")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
		// TODO: Panic handling()?
		panic(err.Error())
	}

	viper.SetDefault("BASE_PATH", "/mnt/nfs/")
	viper.SetDefault("METRIC_COLLECT_INTERVAL_SECONDS", "60")
	viper.SetDefault("RESTIC_PASSWORD_COMMAND", "pass")

	return viper
}
