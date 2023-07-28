/**
 * @Author: lenovo
 * @Description:
 * @File:  middl
 * @Version: 1.0.0
 * @Date: 2023/07/21 17:13
 */

package middl

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"projectManager/project-api/api/user"
	common "projectManager/project-common"
	"projectManager/project-common/errs"
	"projectManager/project-grpc/user/login"
	"time"
)

func TokenVerify() func(c *gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		//1. 从header中获取token
		token := c.GetHeader("Authorization")
		//2.调用user服务进行token认证
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		rsp, err := user.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token})
		if err != nil {
			code, msg := errs.HandleGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort()
			return
		}
		//3. 处理结果
		c.Set("memberID", rsp.Member.Id)
		c.Set("memberName", rsp.Member.Name)
		c.Set("organizationCode", rsp.Member.OrganizationCode)
		c.Next()
	}
}
