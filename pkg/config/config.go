package config

import "github.com/spf13/viper"

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type AppConfig struct {
}

type HTTPConfig struct {
}

type DatabaseConfig struct {
}

type LoggerConfig struct {
	Level string
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
