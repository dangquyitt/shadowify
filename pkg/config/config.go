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
