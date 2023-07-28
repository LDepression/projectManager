/**
 * @Author: lenovo
 * @Description:
 * @File:  cache
 * @Version: 1.0.0
 * @Date: 2023/07/15 19:16
 */

package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
