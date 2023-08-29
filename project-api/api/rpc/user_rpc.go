/**
 * @Author: lenovo
 * @Description:
 * @File:  user_rpc
 * @Version: 1.0.0
 * @Date: 2023/08/23 0:07
 */

package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"projectManager/project-common/discovery"
	"projectManager/project-common/logs"
	"projectManager/project-grpc/user/login"
	"projectManager/project-user/config"
)

var LoginServiceClient login.LoginServiceClient

func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.Dial("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = login.NewLoginServiceClient(conn)
}
