package redisClient

import (
	"context"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
	"strings"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/errorType"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"
)

type IRedisClient interface {
	Get(actionId, key string) *redis.StringCmd
	Set(actionId, key string, value any, expiration time.Duration) *redis.StatusCmd
	SetNx(actionId, key string, value any, expiration time.Duration) *redis.BoolCmd
	Expire(actionId, key string, expiration time.Duration) *redis.BoolCmd
	Exists(actionId string, keys ...string) *redis.IntCmd
	Incr(actionId, key string) *redis.IntCmd
	Del(actionId string, keys ...string) *redis.IntCmd
	HGet(actionId, key, field string) *redis.StringCmd
	HGetAll(actionId, key string) *redis.MapStringStringCmd
	HSet(actionId, key string, values map[string]string) *redis.IntCmd
	HDel(actionId, key string, fields ...string) *redis.IntCmd
	HIncrBy(actionId, key, field string, incr int64) *redis.IntCmd
	MGet(actionId string, keys ...string) *redis.SliceCmd
	SMembers(actionId, key string) *redis.StringSliceCmd
	SMembersMap(actionId, key string) *redis.StringStructMapCmd
	SAdd(actionId, key string, members ...any) *redis.IntCmd
	SRem(actionId, key string, members ...any) *redis.IntCmd
	SCard(actionId, key string) *redis.IntCmd
	Scan(actionId string, cursor uint64, match string, count int64) *redis.ScanCmd
	Publish(c context.Context, channel string, message any) *redis.IntCmd
	ReceiveMessage(c context.Context, pubsub *redis.PubSub, callback func(c context.Context, message *redis.Message))
}

type BaseRedisClient struct {
	name        string
	redisClient *redis.Client
	isLog       bool
}

func NewBaseRedisClient(name string, redisClient *redis.Client) *BaseRedisClient {
	return &BaseRedisClient{name: name, redisClient: redisClient, isLog: false}
}

func NewLogBaseRedisClient(name string, redisClient *redis.Client, isLog bool) *BaseRedisClient {
	return &BaseRedisClient{name: name, redisClient: redisClient, isLog: isLog}
}

func (b *BaseRedisClient) Get(actionId, key string) *redis.StringCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Get", actionId)
	logMessage[logKey.RedisKey] = key

	result := b.redisClient.Get(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Set(actionId, key string, value any, expiration time.Duration) *redis.StatusCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Set", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisValue] = value
	logMessage[logKey.RedisExpiration] = expiration.String()

	result := b.redisClient.Set(context.Background(), key, value, expiration)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) SetNx(actionId, key string, value any, expiration time.Duration) *redis.BoolCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "SetNx", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisValue] = value
	logMessage[logKey.RedisExpiration] = expiration.String()

	result := b.redisClient.SetNX(context.Background(), key, value, expiration)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Expire(actionId, key string, expiration time.Duration) *redis.BoolCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Expire", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisExpiration] = expiration.String()

	result := b.redisClient.Expire(context.Background(), key, expiration)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Exists(actionId string, keys ...string) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Exists", actionId)
	logMessage[logKey.RedisKey] = fmt.Sprintf("%+v", keys)

	result := b.redisClient.Exists(context.Background(), keys...)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Incr(actionId, key string) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Incr", actionId)
	logMessage[logKey.RedisKey] = key

	result := b.redisClient.Incr(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Del(actionId string, keys ...string) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Del", actionId)
	logMessage[logKey.RedisKey] = fmt.Sprintf("%+v", keys)

	result := b.redisClient.Del(context.Background(), keys...)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) HGet(actionId, key, field string) *redis.StringCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "HGet", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisField] = field

	result := b.redisClient.HGet(context.Background(), key, field)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) HGetAll(actionId, key string) *redis.MapStringStringCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "HGetAll", actionId)
	logMessage[logKey.RedisKey] = key

	result := b.redisClient.HGetAll(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) HSet(actionId, key string, values map[string]string) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "HSet", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisValue] = fmt.Sprintf("%+v", values)

	result := b.redisClient.HSet(context.Background(), key, values)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) HDel(actionId, key string, fields ...string) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "HDel", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisField] = fmt.Sprintf("%+v", fields)

	result := b.redisClient.HDel(context.Background(), key, fields...)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) HIncrBy(actionId, key, field string, incr int64) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "HIncrBy", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisField] = field
	logMessage[logKey.RedisValue] = incr

	result := b.redisClient.HIncrBy(context.Background(), key, field, incr)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) MGet(actionId string, keys ...string) *redis.SliceCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "MGet", actionId)
	logMessage[logKey.RedisKey] = fmt.Sprintf("%+v", keys)

	result := b.redisClient.MGet(context.Background(), keys...)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) SMembers(actionId, key string) *redis.StringSliceCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "SMembers", actionId)
	logMessage[logKey.RedisKey] = key

	result := b.redisClient.SMembers(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) SMembersMap(actionId, key string) *redis.StringStructMapCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "SMembersMap", actionId)
	logMessage[logKey.RedisKey] = key

	result := b.redisClient.SMembersMap(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) SAdd(actionId, key string, members ...any) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "SAdd", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisValue] = fmt.Sprintf("%+v", members)

	result := b.redisClient.SAdd(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) SRem(actionId, key string, members ...any) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "SRem", actionId)
	logMessage[logKey.RedisKey] = key
	logMessage[logKey.RedisValue] = fmt.Sprintf("%+v", members)

	result := b.redisClient.SRem(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) SCard(actionId, key string) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "SCard", actionId)
	logMessage[logKey.RedisKey] = key

	result := b.redisClient.SCard(context.Background(), key)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Scan(actionId string, cursor uint64, match string, count int64) *redis.ScanCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Scan", actionId)
	logMessage[logKey.RedisCursor] = cursor
	logMessage[logKey.RedisMatch] = match
	logMessage[logKey.RedisCount] = count

	result := b.redisClient.Scan(context.Background(), cursor, match, count)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) Publish(c context.Context, channel string, message any) *redis.IntCmd {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	b.logBegin(logMessage, now, "Publish", utils.GetActionId(c))
	logMessage[logKey.RedisChannel] = channel
	logMessage[logKey.RedisMessage] = fmt.Sprintf("%+v", message)

	result := b.redisClient.Publish(c, channel, message)

	b.logEnd(logMessage, now)

	return result
}

func (b *BaseRedisClient) ReceiveMessage(c context.Context, pubsub *redis.PubSub, callback func(c context.Context, message *redis.Message)) {
	logMessage := map[string]any{}
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			b.handleErrorRecover(r, logMessage, now)
			return
		}
	}()

	message, err := pubsub.ReceiveMessage(c)
	now = time.Now()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			time.Sleep(5 * time.Second)
		} else if strings.Contains(err.Error(), "use of closed network") || strings.Contains(err.Error(), "EOF") {
			b.logEnd(logMessage, now)
			return
		}

		panic(err)
	}

	b.logBegin(logMessage, now, "ReceiveMessage", "")

	logMessage[logKey.RedisChannel] = message.Channel
	if message != nil {
		utils.SetActionLog(c, contextKey.ActionLogs, logMessage)
		callback(c, message)
	}

	b.logEnd(logMessage, now)
}

func (b *BaseRedisClient) logBegin(logMessage map[string]any, now time.Time, redisMethod string, actionId string) {
	logMessage[logKey.Type] = "redis-client"
	logMessage[logKey.ServerTime] = now.String()
	if actionId != "" {
		logMessage[logKey.Id] = actionId
	} else {
		logMessage[logKey.Id] = uuid.NewString()
	}
	logMessage[logKey.ServiceName] = global.AppName
	logMessage[logKey.HostName] = global.ServerConfig.HostName
	logMessage[logKey.RedisClientName] = b.name
	logMessage[logKey.RedisMethod] = redisMethod
}

func (b *BaseRedisClient) logEnd(logMessage map[string]any, now time.Time) {
	timeUsed := time.Since(now)
	logMessage[logKey.TimeUsed] = timeUsed.String()
	logMessage[logKey.TimeUsedNano] = timeUsed
	if timeUsed.Seconds() > 2 {
		logMessage[logKey.SlowRedis] = true
		utils.OutputLog(logMessage)
	} else if b.isLog {
		utils.OutputLog(logMessage)
	}
}

func (b *BaseRedisClient) handleErrorRecover(r any, logMessage map[string]any, now time.Time) {
	logMessage[logKey.ErrorType] = errorType.InternalServerError
	logMessage[logKey.ErrorMessage] = fmt.Sprintf("%+v", r)
	logMessage[logKey.StackTrace] = string(debug.Stack())
	timeUsed := time.Since(now)
	logMessage[logKey.TimeUsed] = timeUsed.String()
	logMessage[logKey.TimeUsedNano] = timeUsed
	if timeUsed.Seconds() > 2 {
		logMessage[logKey.SlowRedis] = true
		utils.OutputLog(logMessage)
	} else if b.isLog {
		utils.OutputLog(logMessage)
	}
}
