package config

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Database DatabaseConfig
	HTTP     HTTPConfig
}

type DatabaseConfig struct {
	InMemory bool
	DSN      string
}

type HTTPConfig struct {
	Host         string
	Port         string
	HTTPTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
}

type LoggerConf struct {
	Level string
	File  string
}

// NewConfig read configs and return Config
// read from flag --config if exists
// else find config.yaml file in:
// - /etc/calendar
// - $HOME/.calendar
// - $PWD/configs
// - current dir.
func NewConfig() *Config {
	config := &Config{
		Logger: LoggerConf{
			Level: "INFO",
		},
		HTTP: HTTPConfig{
			Host:         "127.0.0.1",
			Port:         "5000",
			HTTPTimeout:  time.Second * 60,
			WriteTimeout: time.Second * 60,
			IdleTimeout:  time.Second * 60,
			ReadTimeout:  time.Second * 60,
		},
		Database: DatabaseConfig{
			InMemory: true,
		},
	}
	cfgFile := viper.GetString("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("http", map[string]string{"Host": "127.0.0.1", "Port": "5000"})

	if cfgFile != "" {
		// Если указан конфиг, то читаем только его
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log.Fatalf("Config %s not found", cfgFile)
		}
		viper.AddConfigPath(filepath.Dir(cfgFile))
		viper.SetConfigName(filepath.Base(cfgFile))
	} else {
		// Если не указан, то ищем конфиги в "default" каталогах
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/calendar")
		viper.AddConfigPath("$HOME/.calendar")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
