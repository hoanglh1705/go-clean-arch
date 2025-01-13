package cachehelper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisHelper struct {
	client *redis.Client
}

func (h *redisHelper) GetTransaction(ctx context.Context, transactionID string) CacheTransactionExecution {
	txPipeline := h.client.TxPipeline()
	return &redisCacheTransaction{
		baseRedisCachePipeline: baseRedisCachePipeline{
			Pipeliner:     txPipeline,
			transactionID: transactionID,
		},
	}
}

func (h *redisHelper) GetPipeline(ctx context.Context, transactionID string) CachePipelineExecution {
	pipeline := h.client.Pipeline()
	return &redisCachePipeline{
		baseRedisCachePipeline: baseRedisCachePipeline{
			Pipeliner:     pipeline,
			transactionID: transactionID,
		},
	}
}

func (h *redisHelper) Exists(ctx context.Context, key string) (err error) {
	indicator, err := h.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if indicator == 0 {
		return redis.Nil
	}
	return nil
}

func (h *redisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (isSucces bool, err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	isSucces, err = h.client.SetNX(ctx, key, string(data), expiration).Result()
	if err != nil {
		return false, err
	}
	return isSucces, nil
}

func (h *redisHelper) Get(ctx context.Context, key string, value interface{}) (err error) {
	data, err := h.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = h.client.Set(ctx, key, string(data), expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) Del(ctx context.Context, key string) (err error) {
	_, err = h.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) Expire(ctx context.Context, key string, expiration time.Duration) (err error) {
	_, err = h.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) GetInterface(ctx context.Context, key string, value interface{}) (interface{}, error) {
	var err error
	data, err := h.client.Get(ctx, key).Result()
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
		outDataValue := reflect.ValueOf(outData)

		if reflect.Indirect(reflect.ValueOf(outDataValue)).IsZero() {
			return nil, errors.New("Get redis nill result")
		}
		if outDataValue.IsZero() {
			return outDataValue.Interface(), nil
		}
		return outDataValue.Elem().Interface(), nil
	}
	var outValue interface{} = outData
	if reflect.TypeOf(outData).ConvertibleTo(typeValue) {
		outValueConverted := reflect.ValueOf(outData).Convert(typeValue)
		outValue = outValueConverted.Interface()
	}
	return outValue, nil
}

func (h *redisHelper) DelMulti(ctx context.Context, keys ...string) error {
	var err error
	pipeline := h.client.TxPipeline()
	pipeline.Del(ctx, keys...)
	_, err = pipeline.Exec(ctx)
	return err
}

func (h *redisHelper) GetKeysByPattern(ctx context.Context, pattern string, cursor uint64, limit int64) ([]string, uint64, error) {
	return h.client.Scan(ctx, cursor, pattern, limit).Result()
}

func (h *redisHelper) SubscribeMessage(ctx context.Context, keySpace string, subscribeFunc SubscribeFunc) {
	subscribes := h.client.Subscribe(ctx, keySpace)
	messageChan := subscribes.Channel()
	for {
		message, ok := <-messageChan
		if ok {
			go func() {
				_ = subscribeFunc(CacheMessage{Message: *message})
			}()
		}
	}
}
func (h *redisHelper) PublishMessage(ctx context.Context, keySpace string, message interface{}) error {
	result := h.client.Publish(ctx, keySpace, message)
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

func (h *redisHelper) GetMulti(ctx context.Context, data interface{}, keys ...string) (result []interface{}, err error) {
	var (
		cmds []redis.Cmder
	)
	p := h.client.Pipeline()
	p.MGet(ctx, keys...)
	cmds, err = p.Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, cmd := range cmds {
		if slice, ok := cmd.(*redis.SliceCmd); ok {
			resultItem, err := slice.Result()
			if err != nil {
				return nil, err
			}
			if len(resultItem) == 0 {
				continue
			}

			// get first one
			result = append(result, resultItem...)
		}

	}
	return result, nil
}

func (h *redisHelper) RenameKey(ctx context.Context, oldkey, newkey string) error {
	var err error
	_, err = h.client.Rename(ctx, oldkey, newkey).Result()
	return err
}

func (h *redisHelper) GetStrLenght(ctx context.Context, key string) (int64, error) {
	return h.client.StrLen(ctx, key).Result()
}

func (h *redisHelper) GetType(ctx context.Context, key string) (string, error) {
	return h.client.Type(ctx, key).Result()
}

func (h *redisHelper) DebugObjectByKey(ctx context.Context, key string) (string, error) {
	return h.client.DebugObject(ctx, key).Result()
}

func (h *redisHelper) TimeExpire(ctx context.Context, key string) (time.Duration, error) {
	return h.client.TTL(ctx, key).Result()
}

func (h *redisHelper) HSet(ctx context.Context, key, mapKey string, mapValue interface{}, expiration time.Duration) (isSet bool, err error) {
	var (
		marshalValue []byte
		boolResult   *redis.BoolCmd
		result       *redis.IntCmd
	)
	if _, isString := mapValue.(string); !isString {
		marshalValue, err = json.Marshal(mapValue)
		if err != nil {
			return isSet, err
		}
	}

	result = h.client.HSet(ctx, key, mapKey, string(marshalValue))
	if result.Err() != nil {
		return isSet, err
	}
	if isSet0, err := result.Result(); isSet0 == 0 || err != nil {
		return isSet0 == 0, err
	}

	if expiration != time.Duration(0) {
		boolResult = h.client.Expire(ctx, key, expiration)
	}

	if isSet, err = boolResult.Result(); !isSet || err != nil {
		return isSet, err
	}
	return true, nil
}

func (h *redisHelper) HSetNX(ctx context.Context, key string, mapKey string, mapValue interface{}, expiration time.Duration) (isSet bool, err error) {
	var (
		marshalValue []byte
		boolResult   *redis.BoolCmd
	)
	if _, isString := mapValue.(string); !isString {
		marshalValue, err = json.Marshal(mapValue)
		if err != nil {
			return isSet, err
		}
	}

	boolResult = h.client.HSetNX(ctx, key, mapKey, string(marshalValue))
	if isSet, err = boolResult.Result(); !isSet || err != nil {
		return isSet, err
	}
	if expiration != time.Duration(0) {
		boolResult = h.client.Expire(ctx, key, expiration)
	}
	if isSet, err = boolResult.Result(); !isSet || err != nil {
		return isSet, err
	}
	return isSet, err
}

func (h *redisHelper) HGet(ctx context.Context, key, mapKey string) (value string, err error) {
	result := h.client.HGet(ctx, key, mapKey)
	if result.Err() != nil {
		return value, err
	}
	if value, err = result.Result(); err != nil {
		return value, err
	}
	return value, nil
}

func (h *redisHelper) HGetAll(ctx context.Context, key string, mapKeys []string) (values map[string]string, err error) {
	result := h.client.HGetAll(ctx, key)
	if result.Err() != nil {
		return values, err
	}
	if values, err = result.Result(); values == nil || err != nil {
		return values, err
	}
	return values, nil
}

func (h *redisHelper) HIncreaseBy(ctx context.Context, key, mapKey string, increase int64) (isIncreased bool, value string, err error) {
	result := h.client.HIncrBy(ctx, key, mapKey, increase)
	if result.Err() != nil {
		return isIncreased, value, err
	}
	var (
		valueInt int64
	)
	if valueInt, err = result.Result(); err != nil {
		return isIncreased, value, err
	}

	return true, strconv.FormatInt(valueInt, 10), nil
}

func (h *redisHelper) HMSet(ctx context.Context, key string, mapData map[string]interface{}, expiration time.Duration) (isSet bool, err error) {
	var (
		inputData    map[string]interface{} = make(map[string]interface{}, len(mapData))
		marshalValue []byte
		stringValue  string
		isString     bool
	)

	for key, value := range mapData {
		if stringValue, isString = value.(string); !isString {
			marshalValue, err = json.Marshal(value)
			if err != nil {
				return isSet, err
			}
		}
		if stringValue == "" {
			stringValue = string(marshalValue)
		}
		inputData[key] = string(stringValue)
	}
	result := h.client.HMSet(ctx, key, inputData)
	if result.Err() != nil {
		return isSet, err
	}
	if ok, err := result.Result(); !ok || err != nil {
		return isSet, err
	}
	if expiration != time.Duration(0) {
		boolResult := h.client.Expire(ctx, key, expiration)
		if isSet, err = boolResult.Result(); !isSet || err != nil {
			return isSet, err
		}
	}

	return true, nil
}

func (h *redisHelper) HMGet(ctx context.Context, key string, fields []string) (result map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	sliceResult := h.client.HMGet(ctx, key, fields...)
	if sliceResult.Err() != nil {
		return result, err
	}

	if results, err = sliceResult.Result(); err != nil {
		return result, err
	}

	result = make(map[string]interface{}, len(results))
	for index, item := range fields {
		result[item] = results[index]
	}
	return result, nil
}
