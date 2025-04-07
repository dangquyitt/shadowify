package config

import "github.com/spf13/viper"

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Database DatabaseConfig `mapstructure:"database"`
	I18n     I18nConfig     `mapstructure:"i18n"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Youtube  YoutubeConfig  `mapstructure:"youtube"`
}

type AppConfig struct {
}

type HTTPConfig struct {
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type I18nConfig struct {
	BundleDir string
}

type LoggerConfig struct {
	Level string
}

type YoutubeConfig struct {
	APIKey string
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
