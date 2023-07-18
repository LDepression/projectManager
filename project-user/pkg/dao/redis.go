/**
 * @Author: lenovo
 * @Description:
 * @File:  redis
 * @Version: 1.0.0
 * @Date: 2023/07/15 19:19
 */

package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"projectManager/project-user/config"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func init() {
	rdb := redis.NewClient(config.C.ReadRedisConfig())
	Rc = &RedisCache{rdb: rdb}
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	return result, err
}
