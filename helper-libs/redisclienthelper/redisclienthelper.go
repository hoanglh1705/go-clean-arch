package redisclienthelper

import (
	"context"
	"crypto/tls"
	"go-clean-arch/helper-libs/tlshelper"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	RedisClientHelper struct {
		Client        *redis.Client
		ClusterClient *redis.ClusterClient
	}

	RedisConfigOptions struct {
		Addrs              []string
		Password           string
		DB                 int
		UseTls             bool
		TlsConfig          *tlshelper.TlsClientOptions
		CaCert             string
		Cert               string
		Key                string
		CaCertFile         string
		CertFile           string
		KeyFile            string
		InsecureSkipVerify bool
	}
)

func NewRedisClientHelper(cfg *RedisConfigOptions) *RedisClientHelper {
	if len(cfg.Addrs) > 1 {
		clusterClient, err := initRedisCluster(cfg)
		if err != nil {
			zap.S().Panic("Failed to init redis cluster", zap.Error(err))
		}
		return &RedisClientHelper{
			ClusterClient: clusterClient,
		}
	}
	client, err := initRedis(cfg)
	if err != nil {
		zap.S().Panic("Failed to init redis", zap.Error(err))
	}
	return &RedisClientHelper{
		Client: client,
	}
}

func initRedisCluster(cfg *RedisConfigOptions) (*redis.ClusterClient, error) {
	var tlsConfig *tls.Config
	var err error
	options := &redis.ClusterOptions{
		Addrs:    []string{cfg.Addrs[0]},
		Password: cfg.Password,
	}

	if cfg.CaCertFile != "" && cfg.CertFile != "" && cfg.KeyFile != "" {
		tlsConfig, err = tlshelper.NewClientTLSConfigFromFile(
			cfg.CertFile,
			cfg.KeyFile,
			cfg.CaCertFile,
			cfg.InsecureSkipVerify,
		)
		if err != nil {
			zap.S().Errorf("Failed to init TLS: %v", err)
		}
	}

	if tlsConfig != nil {
		options.TLSConfig = tlsConfig
	}

	clusterClient := redis.NewClusterClient(options)
	_, err = clusterClient.Ping(context.Background()).Result()
	return clusterClient, err
}

func initRedis(cfg *RedisConfigOptions) (*redis.Client, error) {
	var tlsConfig *tls.Config
	var err error
	options := &redis.Options{
		Addr:     cfg.Addrs[0],
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.CaCertFile != "" && cfg.CertFile != "" && cfg.KeyFile != "" {
		tlsConfig, err = tlshelper.NewClientTLSConfigFromFile(
			cfg.CertFile,
			cfg.KeyFile,
			cfg.CaCertFile,
			cfg.InsecureSkipVerify,
		)
		if err != nil {
			zap.S().Errorf("Failed to init TLS: %v", err)
		}
	}

	if tlsConfig != nil {
		options.TLSConfig = tlsConfig
	}

	client := redis.NewClient(options)
	_, err = client.Ping(context.Background()).Result()
	return client, err
}
