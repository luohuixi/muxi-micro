package sql

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-redis/redis"
)

const CacheNotFound = redis.Nil

type CacheExecute struct {
	rdb      *redis.Client
	CacheTTL time.Duration
	SetTTl   time.Duration //写入缓存的允许最大时长
}

func ConnectCache(addr, password string, number int, CacheTTL, SetTTL time.Duration) (*CacheExecute, *redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       number,
	})

	return &CacheExecute{
		rdb,
		CacheTTL,
		SetTTL,
	}, rdb
}

// cache 管理主键缓存 和 非主键向主键映射缓存
func (c *CacheExecute) GetCache(cachestr string, ctx context.Context) (string, error) {
	cache, err := c.rdb.WithContext(ctx).Get(cachestr).Result()
	if err != nil {
		return "", err
	}
	return cache, nil
}

func (c *CacheExecute) SetCache(cachestr string, ctx context.Context, data any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.rdb.WithContext(ctx).Set(cachestr, val, c.CacheTTL).Err()
}

func (c *CacheExecute) DeleteCache(cachestr string, ctx context.Context) error {
	return c.rdb.WithContext(ctx).Del(cachestr).Err()
}
