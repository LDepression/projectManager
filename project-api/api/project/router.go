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
	group.POST("/save", h.projectSave)
	group.POST("/read", h.readProject)
	group.POST("/project/recycle", h.recycleProject)
	group.POST("/project/recovery", h.recoveryProject)
	group.POST("/project_collect/collect", h.collectProject)
	group.POST("/project/edit", h.editProject)

	t := NewTask()
	group.POST("/task_stages", t.taskStages)
	group.POST("/project_member/index", t.memberProjectList)
	group.POST("/task_stages/tasks", t.taskList)
	group.POST("/task/save", t.saveTask)
	group.POST("/task/sort", t.taskSort)
	group.POST("/task/selfList", t.myTaskList)
	group.POST("/task/read", t.readTask)
	group.POST("/task_member", t.listTaskMember)
	group.POST("/task/taskLog", t.taskLog)
	group.POST("/task/_taskWorkTimeList", t.taskWorkTimeList)
	group.POST("/file/uploadFiles", t.uploadFiles)

}
