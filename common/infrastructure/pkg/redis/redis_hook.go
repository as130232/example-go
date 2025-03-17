package redis

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type Hook interface {
	DialHook(next redis.DialHook) redis.DialHook
	ProcessHook(next redis.ProcessHook) redis.ProcessHook
	ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook
}

type DecimalCricketHook struct{}

func (DecimalCricketHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (DecimalCricketHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		now := time.Now()
		err := next(ctx, cmd)
		timeUsed := time.Since(now)
		timeUsedSeconds := timeUsed.Seconds()
		if timeUsedSeconds >= 3 {
			log.Printf("decimal cricket read-write slow redis command:%+v,timeUsedSeconds:%f", cmd, timeUsedSeconds)
		}
		return err
	}
}

func (DecimalCricketHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

type DecimalCricketReadOnlyHook struct{}

func (DecimalCricketReadOnlyHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (DecimalCricketReadOnlyHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		now := time.Now()
		err := next(ctx, cmd)
		timeUsed := time.Since(now)
		timeUsedSeconds := timeUsed.Seconds()
		if timeUsedSeconds >= 3 {
			log.Printf("decimal cricket read-only slow redis command:%+v,timeUsedSeconds:%f", cmd, timeUsedSeconds)
		}
		return err
	}
}

func (DecimalCricketReadOnlyHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

type LockHook struct{}

func (LockHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (LockHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		now := time.Now()
		err := next(ctx, cmd)
		timeUsed := time.Since(now)
		timeUsedSeconds := timeUsed.Seconds()
		if timeUsedSeconds >= 3 {
			log.Printf("lock read-write slow redis command:%+v,timeUsedSeconds:%f", cmd, timeUsedSeconds)
		}
		return err
	}
}

func (LockHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
