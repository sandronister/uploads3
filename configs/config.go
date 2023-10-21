package configs

import "github.com/spf13/viper"

type awsConf struct {
	Credentials string `mapstructure:"AWS_CREDENTIALS"`
	Password    string `mapstructure:"AWS_PASSWORD"`
	Bucket      string `mapstructure:"BUCKET_AWS"`
}

func GetConfig(path string) (*awsConf, error) {
	var cfg *awsConf
	viper.SetConfigName("APP_S3")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
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
