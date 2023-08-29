/**
 * @Author: lenovo
 * @Description:
 * @File:  interceptor
 * @Version: 1.0.0
 * @Date: 2023/08/02 15:34
 */

package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"projectManager/project-grpc/user/login"
	"projectManager/project-user/internal/dao"
	"projectManager/project-user/internal/repo"
)

type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

func New() *CacheInterceptor {
	cacheMap := make(map[string]any)
	cacheMap["/login.LoginService/Login"] = &login.LoginMessage{}
	return &CacheInterceptor{
		cache:    dao.Rc,
		cacheMap: cacheMap,
	}
}

func (c *CacheInterceptor) Cache() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		fmt.Println("进去.........")
		resp, err = handler(ctx, req)
		fmt.Println("退出........")
		//先查询是否有缓存,如果是有的话 直接返回,没有的话,先请求,然后写入缓存中去
		return
	})
}
