/**
 * @Author: lenovo
 * @Description:
 * @File:  router
 * @Version: 1.0.0
 * @Date: 2023/07/16 22:48
 */

package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"projectManager/project-api/api/middl"
	"projectManager/project-api/router"
)

//Router

type RouterProject struct {
}

func init() {
	log.Println("init Project router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (u *RouterProject) Route(r *gin.Engine) {
	//初始化grpc的连接
	InitRpcProjectClient()
	h := New()
	group := r.Group("/project")
	group.Use(middl.TokenVerify())
	group.POST("/index", h.index)
	group.POST("/project/selfList", h.myProjectList)
	group.POST("/project", h.myProjectList)
	group.POST("/project_template", h.projectTemplate)
}
