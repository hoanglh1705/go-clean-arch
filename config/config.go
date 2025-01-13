package config

import (
	"bytes"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var defaultConfig = []byte(`
app: card-integrate-proxy
env: dev
http_address: 8280
sensitive_fields:
  password: (?P<FIRST>[0-9]{6})(?P<MASK>[0-9]*)(?P<LAST>[0-9]{4})
basic_auth:
  username: admin
  password: 12345678
jwt_exclude_paths:
database:
  host: 0.0.0.0
  port: 5432
  username: pmtrade
  password: Abc12345
  database: pmtrade
  search_path: pst
  auto_migration: false
cache:
  host: 0.0.0.0
  port: 6379
  password: 12345678
`)

type (
	Config struct {
		Env             string            `mapstructure:"env"`
		App             string            `mapstructure:"app"`
		HttpAddress     uint32            `mapstructure:"http_address"`
		SensitiveFields map[string]string `mapstructure:"sensitive_fields"`
		JwtExcludePaths []string          `mapstructure:"jwt_exclude_paths"`
		Database        databaseConfig    `mapstructure:"database"`
		Cache           cacheConfig       `mapstructure:"cache"`
	}

	// kafkaConfig struct {
	// 	Brokers            string         `mapstructure:"brokers"`
	// 	Topics             string         `mapstructure:"topics"`
	// 	Username           string         `mapstructure:"username"`
	// 	Password           string         `mapstructure:"password"`
	// 	Name               string         `mapstructure:"name"`
	// 	WorkerCount        uint32         `mapstructure:"worker_count"`
	// 	TimeDeadlineFormat string         `mapstructure:"time_deadline_format"`
	// 	PriorityDuration   map[string]int `mapstructure:"priority_duraction"`
	// }

	databaseConfig struct {
		Host          string `mapstructure:"host"`
		Port          uint32 `mapstructure:"port"`
		Username      string `mapstructure:"username"`
		Password      string `mapstructure:"password"`
		Database      string `mapstructure:"database"`
		SearchPath    string `mapstructure:"search_path"`
		AutoMigration bool   `mapstructure:"auto_migration"`
	}

	cacheConfig struct {
		Host     string `mapstructure:"host"`
		Port     uint32 `mapstructure:"port"`
		Password string `mapstructure:"password"`
	}
)

func Load() (*Config, error) {
	var cfg = &Config{}
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return nil, err
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
