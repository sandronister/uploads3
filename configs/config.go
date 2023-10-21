package configs

import "github.com/spf13/viper"

type awsConf struct {
	AcessKey  string `mapstructure:"ACESS_KEY"`
	SecretKey string `mapstructure:"SECRET_KEY"`
	Bucket    string `mapstructure:"BUCKET"`
}

func GetConfig(path string) (*awsConf, error) {
	var cfg *awsConf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
