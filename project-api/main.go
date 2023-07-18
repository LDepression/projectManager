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
	_ "projectManager/project-api/api"
	"projectManager/project-api/config"
	"projectManager/project-api/router"
	srv "projectManager/project-common"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}
