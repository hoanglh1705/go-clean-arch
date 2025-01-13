package cachehelper

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type clusterRedisHelper struct {
	clusterClient *redis.ClusterClient
}

func (h *clusterRedisHelper) GetTransaction(ctx context.Context, transactionID string) CacheTransactionExecution {
	txPipeline := h.clusterClient.TxPipeline()
	return &redisCacheTransaction{
		baseRedisCachePipeline: baseRedisCachePipeline{
			Pipeliner:     txPipeline,
			transactionID: transactionID,
		},
	}
}

func (h *clusterRedisHelper) GetPipeline(ctx context.Context, transactionID string) CachePipelineExecution {
	pipeline := h.clusterClient.Pipeline()
	return &redisCachePipeline{
		baseRedisCachePipeline: baseRedisCachePipeline{
			Pipeliner:     pipeline,
			transactionID: transactionID,
		},
	}
}

func (h *clusterRedisHelper) Exists(ctx context.Context, key string) (err error) {
	indicator, err := h.clusterClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if indicator == 0 {
		return redis.Nil
	}
	return nil
}

func (h *clusterRedisHelper) Get(ctx context.Context, key string, value interface{}) (err error) {
	data, err := h.clusterClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = h.clusterClient.Set(ctx, key, string(data), expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (isSuccess bool, err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	_, err = h.clusterClient.SetNX(ctx, key, string(data), expiration).Result()
	if err != nil {
		return false, err
	}
	return isSuccess, nil
}

func (h *clusterRedisHelper) Del(ctx context.Context, key string) (err error) {
	_, err = h.clusterClient.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) (err error) {
	_, err = h.clusterClient.Expire(ctx, key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) GetInterface(ctx context.Context, key string, value interface{}) (interface{}, error) {
	data, err := h.clusterClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	typeValue := reflect.TypeOf(value)
	kind := typeValue.Kind()

	var outData interface{}
	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outData = reflect.New(typeValue).Interface()
	default:
		outData = reflect.Zero(typeValue).Interface()
	}
	err = json.Unmarshal([]byte(data), &outData)
	if err != nil {
		return nil, err
	}

	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		return reflect.ValueOf(outData).Elem().Interface(), nil
	}
	var outValue interface{} = outData
	if reflect.TypeOf(outData).ConvertibleTo(typeValue) {
		outValueConverted := reflect.ValueOf(outData).Convert(typeValue)
		outValue = outValueConverted.Interface()
	}
	return outValue, nil
}

func (h *clusterRedisHelper) DelMulti(ctx context.Context, keys ...string) error {
	var err error
	return err
}

func (h *clusterRedisHelper) GetKeysByPattern(ctx context.Context, pattern string, cursor uint64, limit int64) ([]string, uint64, error) {
	var err error
	var keys []string
	return keys, 0, err
}

func (h *clusterRedisHelper) SubscribeMessage(ctx context.Context, keySpace string, subscribeFunc SubscribeFunc) {
	subscribes := h.clusterClient.Subscribe(ctx, keySpace)
	messageChan := subscribes.Channel()
	for {
		message, ok := <-messageChan
		if ok {
			go func() {
				func() {
					_ = subscribeFunc(CacheMessage{Message: *message})
				}()
			}()
		}
	}
}

func (h *clusterRedisHelper) PublishMessage(ctx context.Context, keySpace string, message interface{}) error {
	result := h.clusterClient.Publish(ctx, keySpace, message)
	var out int64
	var err error
	if out, err = result.Result(); err != nil {
		return err
	}
	if out == 0 {
		return fmt.Errorf("published message with response:  %v", out)
	}
	return nil
}

func (h *clusterRedisHelper) GetMulti(ctx context.Context, data interface{}, keys ...string) ([]interface{}, error) {
	return nil, nil
}

func (h *clusterRedisHelper) RenameKey(ctx context.Context, oldkey, newkey string) error {
	return nil
}

func (h *clusterRedisHelper) GetStrLenght(ctx context.Context, key string) (int64, error) {
	return h.clusterClient.StrLen(ctx, key).Result()
}

func (h *clusterRedisHelper) GetType(ctx context.Context, key string) (string, error) {
	return h.clusterClient.Type(ctx, key).Result()
}

func (h *clusterRedisHelper) DebugObjectByKey(ctx context.Context, key string) (string, error) {
	return h.clusterClient.DebugObject(ctx, key).Result()
}

func (h *clusterRedisHelper) TimeExpire(ctx context.Context, key string) (time.Duration, error) {
	return h.clusterClient.TTL(ctx, key).Result()
}

func (h *clusterRedisHelper) HSet(ctx context.Context, key, mapKey string, mapValue interface{}, expiration time.Duration) (isSet bool, err error) {
	return isSet, err
}

func (h *clusterRedisHelper) HSetNX(ctx context.Context, key string, mapKey string, mapValue interface{}, expiration time.Duration) (isSet bool, err error) {
	return isSet, err
}

func (h *clusterRedisHelper) HGet(ctx context.Context, key, mapKey string) (value string, err error) {
	return value, err
}

func (h *clusterRedisHelper) HGetAll(ctx context.Context, key string, mapKeys []string) (values map[string]string, err error) {
	return values, err
}

func (h *clusterRedisHelper) HIncreaseBy(ctx context.Context, key, mapKey string, increase int64) (isIncreased bool, value string, err error) {
	return isIncreased, value, err
}

func (h *clusterRedisHelper) HMSet(ctx context.Context, key string, mapData map[string]interface{}, expiration time.Duration) (isSet bool, err error) {
	return isSet, err
}

func (h *clusterRedisHelper) HMGet(ctx context.Context, key string, fields []string) (result map[string]interface{}, err error) {
	return result, err
}
