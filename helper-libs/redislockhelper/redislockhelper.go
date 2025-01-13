package redislockhelper

import (
	"context"
	"go-clean-arch/helper-libs/loghelper"
	"go-clean-arch/helper-libs/redisclienthelper"
	"time"

	"github.com/bsm/redislock"
)

type (
	RedisLockerOptions struct {
		RedisClientHelper *redisclienthelper.RedisClientHelper
	}
)

type (
	RedisLockHelper struct {
		locker *redislock.Client
	}
)

func NewRedisLockerHelper(opts *RedisLockerOptions) *RedisLockHelper {
	if opts.RedisClientHelper.ClusterClient == nil && opts.RedisClientHelper.Client == nil {
		loghelper.Logger.Panic("redis client must specific")
		return &RedisLockHelper{}
	}

	if opts.RedisClientHelper.ClusterClient != nil {
		locker := redislock.New(opts.RedisClientHelper.ClusterClient)
		return &RedisLockHelper{
			locker: locker,
		}
	}

	locker := redislock.New(opts.RedisClientHelper.Client)
	return &RedisLockHelper{
		locker: locker,
	}
}

func (h *RedisLockHelper) Obtain(
	ctx context.Context,
	lockKey string,
	timeoutInSeconds int64,
) (*redislock.Lock, error) {
	lock, err := h.locker.Obtain(ctx, lockKey, time.Second*time.Duration(timeoutInSeconds), nil)
	if err != nil {
		return lock, err
	}

	return lock, nil
}

func (h *RedisLockHelper) ObtainWithAutoRefresh(
	ctx context.Context,
	lockKey string,
	timeoutInSeconds int64,
	refreshIntervalInSeconds int64,
) (*redislock.Lock, error) {
	lock, err := h.locker.Obtain(ctx, lockKey, time.Second*time.Duration(timeoutInSeconds), nil)
	if err != nil {
		return lock, err
	}

	refreshCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(refreshIntervalInSeconds))
		defer ticker.Stop()

		for {
			select {
			case <-refreshCtx.Done():
				return
			case <-ticker.C:
				err := lock.Refresh(refreshCtx, time.Second*time.Duration(timeoutInSeconds), nil)
				if err != nil {
					loghelper.Logger.Errorf("failed to refresh lock, err: %v", err)
				}
			}
		}
	}()

	return lock, nil
}

func (h *RedisLockHelper) ObtainWithAutoRefreshCallback(
	ctx context.Context,
	lockKey string,
	timeoutInSeconds int64,
	refreshIntervalInSeconds int64,
	refreshCallback func(),
	refreshErrorCallback func(err error),
) (*redislock.Lock, error) {
	lock, err := h.locker.Obtain(ctx, lockKey, time.Second*time.Duration(timeoutInSeconds), nil)
	if err != nil {
		return lock, err
	}

	refreshCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(refreshIntervalInSeconds))
		defer ticker.Stop()

		for {
			select {
			case <-refreshCtx.Done():
				return
			case <-ticker.C:
				err := lock.Refresh(refreshCtx, time.Second*time.Duration(timeoutInSeconds), nil)
				if err != nil {
					if refreshErrorCallback != nil {
						refreshErrorCallback(err)
					}
				} else {
					if refreshCallback != nil {
						refreshCallback()
					}
				}
			}
		}
	}()

	return lock, nil
}
