/**
 * @Author: lenovo
 * @Description:
 * @File:  rpc
 * @Version: 1.0.0
 * @Date: 2023/07/16 22:53
 */

package project

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"projectManager/project-common/discovery"
	"projectManager/project-common/logs"
	"projectManager/project-grpc/project"
	task "projectManager/project-grpc/task"
	"projectManager/project-user/config"
)

var ProjectServiceClient project.ProjectServiceClient
var TaskServiceClient task.TaskServiceClient

func InitRpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.Dial("etcd:///project", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ProjectServiceClient = project.NewProjectServiceClient(conn)
	TaskServiceClient = task.NewTaskServiceClient(conn)
}
