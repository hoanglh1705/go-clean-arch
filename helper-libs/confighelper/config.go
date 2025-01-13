package confighelper

import (
	"bytes"
	"strings"

	"github.com/spf13/viper"
)

func Load(cfgMap any, defaultConfig []byte) error {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return err
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfgMap)
	if err != nil {
		return err
	}
	return nil
}

type (
	ServerConfig struct {
		RateLimit           int64  `mapstructure:"rate_limit"`
		UseTls              bool   `mapstructure:"use_tls"`
		TlsCertFile         string `mapstructure:"tls_cert_file"`
		TlsKeyFile          string `mapstructure:"tls_key_file"`
		TlsRootCACertFile   string `mapstructure:"tls_rootca_cert_file"`
		TlsCertBase64       string `mapstructure:"tls_cert_base64"`
		TlsKeyBase64        string `mapstructure:"tls_key_base64"`
		TlsRootCACertBase64 string `mapstructure:"tls_rootca_cert_base64"`
		InsecureSkipVerify  bool   `mapstructure:"insecure_skip_verify"`
		UseCORs             bool   `mapstructure:"use_cors"`
	}

	DefaultMerchantConfig struct {
		Timezone string `mapstructure:"timezone"`
	}

	BasicAuth struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	ClientAuth struct {
		ClientId     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		Token        string `mapstructure:"token"`
	}

	SqlDatabaseConfig struct {
		Addrs               []string `mapstructure:"addrs"`
		Host                string   `mapstructure:"host"`
		Port                int      `mapstructure:"port"`
		Username            string   `mapstructure:"username"`
		Password            string   `mapstructure:"password"`
		Database            string   `mapstructure:"database"`
		Schema              string   `mapstructure:"schema"`
		UseTls              bool     `mapstructure:"use_tls"`
		TlsMode             string   `mapstructure:"tls_mode"`
		TlsRootCACertFile   string   `mapstructure:"tls_rootca_cert_file"`
		TlsKeyFile          string   `mapstructure:"tls_key_file"`
		TlsCertFile         string   `mapstructure:"tls_cert_file"`
		TlsRootCACertBase64 string   `mapstructure:"tls_rootca_cert_base64"`
		TlsKeyBase64        string   `mapstructure:"tls_key_base64"`
		TlsCertBase64       string   `mapstructure:"tls_cert_base64"`
		InsecureSkipVerify  bool     `mapstructure:"insecure_skip_verify"`
		MaxOpenConns        int      `mapstructure:"max_open_conns"`
		LoggingEnabled      bool     `mapstructure:"logging_enabled"`
		UseLoggingDb        bool     `mapstructure:"use_logging_db"`
		AutoMigration       bool     `mapstructure:"auto_migration"`
	}

	RedisCacheConfig struct {
		Addrs               []string `mapstructure:"addrs"`
		Host                string   `mapstructure:"host"`
		Port                int      `mapstructure:"port"`
		Username            string   `mapstructure:"username"`
		Password            string   `mapstructure:"password"`
		UseTls              bool     `mapstructure:"use_tls"`
		TlsRootCACertFile   string   `mapstructure:"tls_rootca_cert_file"`
		TlsKeyFile          string   `mapstructure:"tls_key_file"`
		TlsCertFile         string   `mapstructure:"tls_cert_file"`
		TlsRootCACertBase64 string   `mapstructure:"tls_rootca_cert_base64"`
		TlsKeyBase64        string   `mapstructure:"tls_key_base64"`
		TlsCertBase64       string   `mapstructure:"tls_cert_base64"`
		InsecureSkipVerify  bool     `mapstructure:"insecure_skip_verify"`
	}

	GrpcClientConfig struct {
		Addrs               []string `mapstructure:"addrs"`
		Host                string   `mapstructure:"host"`
		Port                int      `mapstructure:"port"`
		Username            string   `mapstructure:"username"`
		Password            string   `mapstructure:"password"`
		UseTls              bool     `mapstructure:"use_tls"`
		TlsRootCACertFile   string   `mapstructure:"tls_rootca_cert_file"`
		TlsKeyFile          string   `mapstructure:"tls_key_file"`
		TlsCertFile         string   `mapstructure:"tls_cert_file"`
		TlsRootCACertBase64 string   `mapstructure:"tls_rootca_cert_base64"`
		TlsKeyBase64        string   `mapstructure:"tls_key_base64"`
		TlsCertBase64       string   `mapstructure:"tls_cert_base64"`
		InsecureSkipVerify  bool     `mapstructure:"insecure_skip_verify"`
	}

	HttpClientConfig struct {
		Addrs               string `mapstructure:"addrs"`
		Host                string `mapstructure:"host"`
		Port                string `mapstructure:"port"`
		Username            string `mapstructure:"username"`
		Password            string `mapstructure:"password"`
		UseHttp2            bool   `mapstructure:"use_http2"`
		UseTls              bool   `mapstructure:"use_tls"`
		TlsCertFile         string `mapstructure:"tls_cert_file"`
		TlsKeyFile          string `mapstructure:"tls_key_file"`
		TlsRootCACertFile   string `mapstructure:"tls_rootca_cert_file"`
		TlsCertBase64       string `mapstructure:"tls_cert_base64"`
		TlsKeyBase64        string `mapstructure:"tls_key_base64"`
		TlsRootCACertBase64 string `mapstructure:"tls_rootca_cert_base64"`
		InsecureSkipVerify  bool   `mapstructure:"insecure_skip_verify"`
	}

	TlsConfig struct {
		UseTls             bool   `mapstructure:"use_tls"`
		CertFile           string `mapstructure:"cert_file"`
		KeyFile            string `mapstructure:"key_file"`
		RootCACertFile     string `mapstructure:"rootca_cert_file"`
		CertBase64         string `mapstructure:"cert_base64"`
		KeyBase64          string `mapstructure:"key_base64"`
		RootCACertBase64   string `mapstructure:"rootca_cert_base64"`
		InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
	}

	WorkflowConfig struct {
		Topic            string `mapstructure:"topic"`
		RequestTimeout   int64  `mapstructure:"request_timeout"`
		ScheduleInterval int64  `mapstructure:"schedule_interval"`
		ScheduleCron     int64  `mapstructure:"schedule_cron"`
		MaxAttempt       int    `mapstructure:"max_attempt"`
		MaxConcurrent    int    `mapstructure:"max_concurrent"`
	}
)
