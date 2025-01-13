package redisleaderhelper

import (
	"context"
	"go-clean-arch/helper-libs/loghelper"
	"go-clean-arch/helper-libs/redisclienthelper"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type (
	RedisLeaderHelper struct {
		redisClusterClient *redis.ClusterClient
		redisClient        *redis.Client
		locker             *redislock.Client
		lock               *redislock.Lock
		isStarted          bool
		wasLeading         bool
		canLead            bool
		id                 string
		renewTimeoutId     *time.Timer
		electTimeoutId     *time.Timer
		ttl                int64
		wait               int64
		key                string
		becomeLeaderFunc   func()
		notElectedFunc     func()
		demotedFunc        func()
		errorFunc          func(err error)
	}

	RedisLeaderOptions struct {
		RedisClientHelper *redisclienthelper.RedisClientHelper
		Ttl               int64
		Wait              int64
		Key               string
		BecomeLeaderFunc  func()
		NotElectedFunc    func()
		DemotedFunc       func()
		ErrorFunc         func(err error)
	}
)

func NewRedisLeaderHelper(opts *RedisLeaderOptions) *RedisLeaderHelper {
	if opts.RedisClientHelper.ClusterClient == nil && opts.RedisClientHelper.Client == nil {
		loghelper.Logger.Panic("redis client must specific")
		return &RedisLeaderHelper{}
	}

	if opts.RedisClientHelper.ClusterClient != nil {
		locker := redislock.New(opts.RedisClientHelper.ClusterClient)
		return &RedisLeaderHelper{
			locker:             locker,
			redisClusterClient: opts.RedisClientHelper.ClusterClient,
			ttl:                opts.Ttl,
			wait:               opts.Wait,
			key:                opts.Key,
			becomeLeaderFunc:   opts.BecomeLeaderFunc,
			demotedFunc:        opts.DemotedFunc,
			errorFunc:          opts.ErrorFunc,
			notElectedFunc:     opts.NotElectedFunc,
		}
	}

	locker := redislock.New(opts.RedisClientHelper.Client)
	return &RedisLeaderHelper{
		locker:           locker,
		redisClient:      opts.RedisClientHelper.Client,
		ttl:              opts.Ttl,
		wait:             opts.Wait,
		key:              opts.Key,
		becomeLeaderFunc: opts.BecomeLeaderFunc,
		demotedFunc:      opts.DemotedFunc,
		errorFunc:        opts.ErrorFunc,
		notElectedFunc:   opts.NotElectedFunc,
	}
}

func (h *RedisLeaderHelper) elect() {
	var err error
	h.lock, err = h.locker.Obtain(
		context.Background(),
		h.key,
		time.Millisecond*time.Duration(h.ttl),
		&redislock.Options{
			Token: h.id,
		},
	)

	if err == redislock.ErrNotObtained {
		h.notifyNotElected()
		h.electTimeoutId = time.NewTimer(time.Millisecond * time.Duration(h.wait))
		go func() {
			<-h.electTimeoutId.C
			h.elect()
		}()
		return
	} else if err != nil {
		if h.isStarted {
			h.notifyError(err)
		}
		h.electTimeoutId = time.NewTimer(time.Millisecond * time.Duration(h.wait))
		go func() {
			<-h.electTimeoutId.C
			h.elect()
		}()
		return
	}

	h.notifyElected()
	h.wasLeading = true
	if !h.canLead {
		h.stop()
	}
	h.renewTimeoutId = time.NewTimer(time.Millisecond * time.Duration(h.ttl/2))
	go func() {
		<-h.renewTimeoutId.C
		h.renew()
	}()
}

func (h *RedisLeaderHelper) renew() {
	err := h.lock.Refresh(
		context.Background(),
		time.Millisecond*time.Duration(h.ttl),
		&redislock.Options{
			Token: h.id,
		},
	)
	if err == redislock.ErrNotObtained {
		if h.wasLeading {
			h.wasLeading = false
			h.notifyDemoted()
		}
		if h.renewTimeoutId != nil {
			h.renewTimeoutId.Stop()
		}
		h.electTimeoutId = time.NewTimer(time.Millisecond * time.Duration(h.wait))
		go func() {
			<-h.electTimeoutId.C
			h.elect()
		}()
		return
	} else if err != nil {
		if h.isStarted {
			h.notifyError(err)
		}
		h.electTimeoutId = time.NewTimer(time.Millisecond * time.Duration(h.wait))
		go func() {
			<-h.electTimeoutId.C
			h.elect()
		}()

		return
	}

	h.wasLeading = true
	h.renewTimeoutId = time.NewTimer(time.Millisecond * time.Duration(h.ttl/2))
	go func() {
		<-h.renewTimeoutId.C
		h.renew()
	}()
}

func (h *RedisLeaderHelper) stop() {
	h.canLead = false
	if h.renewTimeoutId != nil {
		h.renewTimeoutId.Stop()
	}
	if h.electTimeoutId != nil {
		h.electTimeoutId.Stop()
	}
	err := h.lock.Release(
		context.Background(),
	)
	if err != nil {
		if h.isStarted {
			h.notifyError(err)
		}
		return
	}

	h.notifyDemoted()
	h.wasLeading = false
}

func (h *RedisLeaderHelper) IsLeader() bool {
	if h.redisClusterClient != nil {
		return h.id == h.redisClusterClient.Get(context.Background(), h.key).String()
	}
	return h.id == h.redisClient.Get(context.Background(), h.key).String()
}

func (h *RedisLeaderHelper) notifyElected() {
	if h.becomeLeaderFunc != nil {
		h.becomeLeaderFunc()
	}
}

func (h *RedisLeaderHelper) notifyNotElected() {
	if h.notElectedFunc != nil {
		h.notElectedFunc()
	}
}

func (h *RedisLeaderHelper) notifyDemoted() {
	if h.demotedFunc != nil {
		h.demotedFunc()
	}
}

func (h *RedisLeaderHelper) notifyError(err error) {
	if h.errorFunc != nil {
		h.errorFunc(err)
	}
}

func (h *RedisLeaderHelper) Run(
	ctx context.Context,
) {
	h.isStarted = true
	h.canLead = true
	h.elect()
}

type RedisLocker struct {
	Redis *redis.Client
}

type RedisClusterLocker struct {
	Redis *redis.ClusterClient
}
