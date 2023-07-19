/**
 * @Author: lenovo
 * @Description:
 * @File:  router
 * @Version: 1.0.0
 * @Date: 2023/07/11 15:28
 */

package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"projectManager/project-common/discovery"
	"projectManager/project-common/logs"
	"projectManager/project-grpc/user/login"
	"projectManager/project-user/config"
	loginServiceV1 "projectManager/project-user/pkg/service/login.service.v1"
)

// 接口
type Router interface {
	Route(r *gin.Engine)
}

type RegisterRouter struct{}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type GrpcConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	//构建grpc的服务
	info := discovery.Server{
		Name:    config.C.GC.Name,
		Addr:    config.C.GC.Addr,
		Version: config.C.GC.Version,
		Weight:  config.C.GC.Weight, //用于负载均衡
	}
	r := discovery.NewRegister(config.C.EtcdConfig.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
func RegisterGrpc() *grpc.Server {
	c := GrpcConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			login.RegisterLoginServiceServer(g, loginServiceV1.New())
		},
	}
	s := grpc.NewServer()
	c.RegisterFunc(s)

	lis, err := net.Listen("tcp", config.C.GC.Addr)
	if err != nil {
		log.Println("can not listen")
	}
	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Println("")
		}
	}()
	return s
}
