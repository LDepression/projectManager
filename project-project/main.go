/**
 * @Author: lenovo
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2023/07/11 15:03
 */

package main

import (
	"github.com/gin-gonic/gin"
	srv "projectManager/project-common"
	"projectManager/project-project/config"
	"projectManager/project-project/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	//grpc服务注册
	gc := router.RegisterGrpc()
	//grpc服务注册到etcd上
	router.RegisterEtcdServer()
	stop := func() {
		gc.Stop()
	}
	//初始化rpc的调用
	router.InitUserRpc()
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
