/**
 * @Author: lenovo
 * @Description:
 * @File:  router
 * @Version: 1.0.0
 * @Date: 2023/07/16 22:48
 */

package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"projectManager/project-api/router"
)

//Router

type RouterUser struct {
}

func init() {
	log.Println("INIT USER ROUTER")
	ru := &RouterUser{}
	router.Register(ru)
}

func (u *RouterUser) Route(r *gin.Engine) {
	//初始化grpc的连接
	InitRpcUserClient()
	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/register", h.register)
}
