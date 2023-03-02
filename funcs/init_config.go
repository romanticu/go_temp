package funcs

import "github.com/spf13/viper"

type Config struct {
	Mysql MysqlConfig
}

type MysqlConfig struct {
	Host string
	DBName string
	Password string
	User string
	Port string
}

func InitConfig(configName string, configPaths []string, config interface{}) error {
	vp := viper.New()
	vp.SetConfigName(configName)
	for _, path := range configPaths {
		vp.AddConfigPath(path)
	}
	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	err := vp.Unmarshal(config)
	if err != nil {
		return nil
	}

	return nil
}