package config

import "github.com/spf13/viper"

type Config struct {
	App           AppConfig           `mapstructure:"app"`
	HTTP          HTTPConfig          `mapstructure:"http"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Logger        LoggerConfig        `mapstructure:"logger"`
	Youtube       YoutubeConfig       `mapstructure:"youtube"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
}

type AppConfig struct {
}

type HTTPConfig struct {
}

type KeycloakConfig struct {
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Host         string `mapstructure:"host"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
}

type ElasticsearchConfig struct {
	URL string `mapstructure:"url"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type LoggerConfig struct {
	File    LoggerFileConfig    `mapstructure:"file"`
	Console LoggerConsoleConfig `mapstructure:"console"`
}

type LoggerFileConfig struct {
	Level      string `mapstructure:"level"`
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}

type LoggerConsoleConfig struct {
	Level   string `mapstructure:"level"`
	Enabled bool   `mapstructure:"enabled"`
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
