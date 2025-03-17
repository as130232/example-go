package utils

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	rdb        *redis.Client
	ctx        context.Context
	key        string        // 鎖名稱
	lockValue  string        // 辨識
	expiration time.Duration // 過期時間
}

type IRedisLock interface {
	Lock() bool
	Block(retryTime time.Duration) bool // 持续获取锁
	Release() bool
	ForceRelease()
}

// GetLock 生成鎖
func GetLock(rdb *redis.Client, key string, lockValue string, expiration time.Duration) IRedisLock {
	return &RedisLock{
		rdb,
		context.Background(),
		key,
		lockValue,
		expiration,
	}
}

// GetLockWithUuid 生成鎖時自動塞 uuid 到 lockValue
func GetLockWithUuid(rdb *redis.Client, key string, expiration time.Duration) IRedisLock {
	return &RedisLock{
		rdb,
		context.Background(),
		key,
		uuid.New().String(),
		expiration,
	}
}

// Block 阻塞1秒後，嘗試重新獲取鎖，持續指定時間
func (r *RedisLock) Block(retryTime time.Duration) bool {
	starting := time.Now().Unix()
	retrySeconds := int64(retryTime.Seconds())
	for {
		if !r.Lock() {
			time.Sleep(retryTime)
			if (time.Now().Unix() - retrySeconds) >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

// Lock 嘗試上鎖
func (r *RedisLock) Lock() bool {
	// SetNX：SET if Not EXISTS
	return r.rdb.SetNX(r.ctx, r.key, r.lockValue, r.expiration).Val()
}

// ForceRelease 強制釋放鎖
func (r *RedisLock) ForceRelease() {
	r.rdb.Del(r.ctx, r.key).Val()
}

// 釋放鎖 Lua 脚本，會檢查鎖定內容
// KEYS[1]: RedisLock.key
// ARGV[1]: RedisLock.lockValue
var releaseLockLuaScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`)

// Release 跑 lua 腳本釋放鎖
func (r *RedisLock) Release() bool {
	result := releaseLockLuaScript.Run(
		r.ctx, r.rdb,
		[]string{r.key}, // KEYS[1]
		r.lockValue,     // ARGV[1]
	).Val().(int64)
	return result != 0
}
