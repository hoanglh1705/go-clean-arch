package cachehelper

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *baseRedisCachePipeline) Exec(ctx context.Context) (result []CachePipelineResult, err error) {
	var (
		outputResult []redis.Cmder
	)
	outputResult, err = r.Pipeliner.Exec(ctx)
	if err != nil {
		return nil, err
	}

	result = make([]CachePipelineResult, len(outputResult))
	for index, item := range outputResult {
		switch v := item.(type) {
		case *redis.ZSliceCmd:
			resutlReturn := make([]interface{}, len(v.Val()))
			for i, val := range v.Val() {
				resutlReturn[i] = RedisZSliceResult{
					Z: val,
				}
			}
			result[index] = CachePipelineResult{
				Result: resutlReturn,
				Err:    item.Err(),
			}
		case *redis.StringSliceCmd:
			resutlReturn := make([]interface{}, len(v.Val()))
			for i, val := range v.Val() {
				resutlReturn[i] = val
			}
			result[index] = CachePipelineResult{
				Result: resutlReturn,
				Err:    item.Err(),
			}
		case *redis.StringCmd:
			result[index] = CachePipelineResult{
				Result: []interface{}{v.Val()},
				Err:    item.Err(),
			}
		case *redis.IntCmd:
			result[index] = CachePipelineResult{
				Result: []interface{}{v.Val()},
				Err:    item.Err(),
			}
		default:
			result[index] = CachePipelineResult{
				Result: item.Args(),
				Err:    item.Err(),
			}
		}
	}
	return result, nil
}

func (r *baseRedisCachePipeline) Discard(ctx context.Context) error {
	r.Pipeliner.Discard()
	return nil
}

func (r *baseRedisCachePipeline) BuildCommand(ctx context.Context, cacheCommandType CacheCommandType, data ...interface{}) (err error) {

	if len(data) == 0 {
		return errors.New("missing data to process")
	}
	var (
		keyCache string = data[0].(string)
		cmd      redis.Cmder
	)
	switch cacheCommandType {
	case CacheCommandTypeGetInterface:
		cmd = r.Pipeliner.Get(ctx, keyCache)
	case CacheCommandTypeAddMemberWithScore:
		member := redis.Z{
			Member: data[1].(string),
			Score:  data[2].(float64),
		}
		cmd = r.Pipeliner.ZAdd(ctx, keyCache, member)
	case CacheCommandTypeGetMembersWithScore:
		min, _ := strconv.ParseInt(data[1].(string), 10, 64)
		max, _ := strconv.ParseInt(data[2].(string), 10, 64)
		cmd = r.Pipeliner.ZRangeWithScores(ctx, keyCache, min, max)
	case CacheCommandTypeRemoveMembersWithScore:
		var (
			min string = data[1].(string)
			max string = data[2].(string)
		)
		cmd = r.Pipeliner.ZRemRangeByScore(ctx, keyCache, min, max)

	case CacheCommandTypeExpire:
		var (
			ttl      uint32        = data[1].(uint32)
			duration time.Duration = data[2].(time.Duration)
		)
		cmd = r.Pipeliner.Expire(ctx, keyCache, time.Duration(ttl)*(duration))
	case CacheCommandTypeSetNX:
		var (
			value    interface{}   = data[1]
			ttl      uint32        = data[2].(uint32)
			duration time.Duration = data[3].(time.Duration)
		)

		cmd = r.Pipeliner.SetNX(ctx, keyCache, value, time.Duration(ttl)*(duration))
	case CacheCommandTypeIncrease:
		cmd = r.Pipeliner.Incr(ctx, keyCache)
	case CacheCommandTypeDel:
		cmd = r.Pipeliner.Del(ctx, keyCache)
	default:
		return errors.New("not found any matched command type to process")
	}
	if cmd != nil && cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *baseRedisCachePipeline) GetCommands(context.Context) (CacheLazyExecute, error) {
	return r, nil
}
