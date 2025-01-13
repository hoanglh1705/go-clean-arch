package cachehelper

import (
	"context"
	"go-clean-arch/helper-libs/loghelper"
	"go-clean-arch/helper-libs/redisclienthelper"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheMessage struct {
	redis.Message
}

type SubscribeFunc func(CacheMessage) error

type (
	CacheConfigOptions struct {
		RedisClientHelper *redisclienthelper.RedisClientHelper
	}
)

// CacheHelper is helper of Cache
type CacheHelper interface {
	Exists(ctx context.Context, key string) error
	Get(ctx context.Context, key string, value interface{}) error
	GetInterface(ctx context.Context, key string, value interface{}) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	DelMulti(ctx context.Context, keys ...string) error
	GetKeysByPattern(ctx context.Context, pattern string, cursor uint64, limit int64) ([]string, uint64, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	SubscribeMessage(ctx context.Context, keySpace string, subscribeFunc SubscribeFunc)
	PublishMessage(ctx context.Context, keySpace string, message interface{}) error
	GetMulti(ctx context.Context, data interface{}, keys ...string) ([]interface{}, error)
	//
	RenameKey(ctx context.Context, oldKey, newKey string) error
	GetStrLenght(ctx context.Context, key string) (int64, error)
	GetType(ctx context.Context, key string) (string, error)
	DebugObjectByKey(ctx context.Context, key string) (string, error)
	TimeExpire(ctx context.Context, key string) (time.Duration, error) // return second
	HSet(ctx context.Context, key, mapKey string, mapValue interface{}, expiration time.Duration) (bool, error)
	HSetNX(ctx context.Context, key string, mapKey string, mapValue interface{}, expiration time.Duration) (bool, error)
	HGet(ctx context.Context, key, mapKey string) (string, error)
	HGetAll(ctx context.Context, key string, mapKeys []string) (map[string]string, error)
	HIncreaseBy(ctx context.Context, key, mapKey string, increase int64) (bool, string, error)
	HMSet(ctx context.Context, key string, mapData map[string]interface{}, expiration time.Duration) (bool, error)
	HMGet(ctx context.Context, key string, fields []string) (map[string]interface{}, error)
}

type CacheHelperEnhancement interface {
	CacheHelper
	GetTransaction(ctx context.Context, transactionID string) CacheTransactionExecution
	GetPipeline(ctx context.Context, transactionID string) CachePipelineExecution
}

type CacheCommandType string

const (
	CacheCommandTypeGet                    CacheCommandType = "CacheCommandTypeGet"
	CacheCommandTypeGetInterface           CacheCommandType = "CacheCommandTypeGetInterface"
	CacheCommandTypeAddMemberWithScore     CacheCommandType = "CacheCommandTypeAddMemberWithScore"
	CacheCommandTypeGetMembersWithScore    CacheCommandType = "CacheCommandTypeGetMembersWithScore"
	CacheCommandTypeRemoveMembersWithScore CacheCommandType = "CacheCommandTypeRemoveMembersWithScore"
	CacheCommandTypeExpire                 CacheCommandType = "CacheCommandTypeExpire"
	CacheCommandTypeSetNX                  CacheCommandType = "CacheCommandTypeSetNX"
	CacheCommandTypeIncrease               CacheCommandType = "CacheCommandTypeIncrease"
	CacheCommandTypeDel                    CacheCommandType = "CacheCommandTypeDel"
)

type (
	CacheLazyExecute interface {
		Exec(context.Context) ([]CachePipelineResult, error)
		Discard(context.Context) error
	}

	CacheMutilCommandBuilder interface {
		BuildCommand(ctx context.Context, cacheCommandType CacheCommandType, data ...interface{}) error
		GetCommands(context.Context) (CacheLazyExecute, error)
	}

	CacheTransactionExecution interface {
		CacheMutilCommandBuilder
	}

	CachePipelineExecution interface {
		CacheMutilCommandBuilder
	}
	redisCacheTransaction struct {
		baseRedisCachePipeline
	}

	redisCachePipeline struct {
		baseRedisCachePipeline
	}

	baseRedisCachePipeline struct {
		redis.Pipeliner
		transactionID string
	}

	CachePipelineResult struct {
		Result []interface{}
		Err    error
	}

	RedisZSliceResult struct {
		redis.Z
	}
)

// CacheOption represents cache option
type CacheOption struct {
	Key   string
	Value interface{}
}

// NewCacheHelper creates an instance
func NewCacheHelper(opts *CacheConfigOptions) CacheHelper {
	if opts.RedisClientHelper.ClusterClient == nil && opts.RedisClientHelper.Client == nil {
		loghelper.Logger.Panic("redis client must specific")
		return &redisHelper{}
	}

	if opts.RedisClientHelper.ClusterClient != nil {
		return &clusterRedisHelper{
			clusterClient: opts.RedisClientHelper.ClusterClient,
		}
	}

	return &redisHelper{
		client: opts.RedisClientHelper.Client,
	}
}
